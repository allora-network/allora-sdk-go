package pool

import (
	"context"
	"fmt"
	"math"
	"math/rand"
	"sort"
	"sync"
	"time"

	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// ClientPoolManager manages a pool of Client instances with health tracking and load balancing
// It reuses the health management algorithms from NodeManager but operates on client instances
type ClientPoolManager struct {
	mu                       sync.RWMutex
	active, cooling          []ClientInfo
	currentIndex             int // round-robin index for load distribution
	checkRate                time.Duration
	coolingThreshold         float64
	minActiveStreak          int
	startReactivatedSuccRate float64
	rateLimitDelayIncrease   time.Duration
	maxRateLimitDelay        time.Duration
	logger                   zerolog.Logger

	backMu     sync.Mutex
	backoff    map[string]*backoffState
	baseDelay  time.Duration
	maxDelay   time.Duration
	jitterFrac float64
}

// ClientInfo wraps an Client with health tracking metadata
type ClientInfo struct {
	Client         interfaces.Client
	MaxRetries     int
	successRate    float64
	latEWMA        float64 // exponential-weighted moving average in milliseconds
	healthStreak   int
	rateLimitDelay time.Duration
}

const (
	defaultClientCheckRate                = 10 * time.Second
	defaultClientCoolingThreshold         = 0.5
	defaultClientMinActiveStreak          = 3
	defaultClientStartReactivatedSuccRate = 0.8
	defaultClientRateLimitDelayIncrease   = 100 * time.Millisecond
	defaultClientMaxRateLimitDelay        = 2 * time.Second

	// Backoff constants
	defaultBaseDelay  = 100 * time.Millisecond
	defaultMaxDelay   = 30 * time.Second
	defaultJitterFrac = 0.1
)

type backoffState struct {
	failures int
	until    time.Time
}

// NewClientPoolManager creates a new client pool manager with the provided clients
func NewClientPoolManager(clients []interfaces.Client, logger zerolog.Logger) *ClientPoolManager {
	if len(clients) == 0 {
		panic("ClientPoolManager requires at least one client")
	}

	clientInfos := make([]ClientInfo, len(clients))
	for i, client := range clients {
		clientInfos[i] = ClientInfo{
			Client:      client,
			MaxRetries:  2,   // Default retry count (matches NodeManager LB=2)
			successRate: 1.0, // Start with perfect success rate
		}
	}

	// Use random starting offset to distribute initial load
	startOffset := 0
	if len(clients) > 1 {
		startOffset = int(time.Now().UnixNano()) % len(clients)
	}

	cpm := &ClientPoolManager{
		active:                   clientInfos,
		currentIndex:             startOffset,
		checkRate:                defaultClientCheckRate,
		coolingThreshold:         defaultClientCoolingThreshold,
		minActiveStreak:          defaultClientMinActiveStreak,
		startReactivatedSuccRate: defaultClientStartReactivatedSuccRate,
		rateLimitDelayIncrease:   defaultClientRateLimitDelayIncrease,
		maxRateLimitDelay:        defaultClientMaxRateLimitDelay,
		logger:                   logger.With().Str("component", "client_pool_manager").Logger(),

		// Initialize backoff management
		backoff:    make(map[string]*backoffState),
		baseDelay:  defaultBaseDelay,
		maxDelay:   defaultMaxDelay,
		jitterFrac: defaultJitterFrac,
	}

	cpm.sortActive()

	go cpm.healthLoop()

	return cpm
}

func (cpm *ClientPoolManager) Close() {
	cpm.mu.Lock()
	defer cpm.mu.Unlock()

	for i := range cpm.active {
		if err := cpm.active[i].Client.Close(); err != nil {
			cpm.logger.Error().Err(err).Str("client_url", cpm.active[i].Client.GetEndpointURL()).Msg("error closing active client")
		}
	}
	for i := range cpm.cooling {
		if err := cpm.cooling[i].Client.Close(); err != nil {
			cpm.logger.Error().Err(err).Str("client_url", cpm.cooling[i].Client.GetEndpointURL()).Msg("error closing cooling client")
		}
	}
}

// GetClient returns the next client using round-robin selection among active clients
// The skip function allows filtering clients based on custom criteria (e.g. backoff state)
func (cpm *ClientPoolManager) GetClient(skip func(interfaces.Client) bool) (interfaces.Client, bool) {
	if skip == nil {
		skip = func(interfaces.Client) bool { return false } // Default: don't skip anything
	}
	cpm.mu.Lock()
	defer cpm.mu.Unlock()

	if len(cpm.active) == 0 {
		// fallback to cooling pool
		for i := range cpm.cooling {
			client := cpm.cooling[i].Client
			if !skip(client) {
				// Track cooling pool usage
				cpm.logger.Warn().Str("client_url", client.GetEndpointURL()).Msg("using cooling pool client - no active clients available")
				return client, true
			}
		}
		return nil, false
	}

	// Round-robin through active clients
	startIdx := cpm.currentIndex
	for i := 0; i < len(cpm.active); i++ {
		idx := (startIdx + i) % len(cpm.active)
		client := cpm.active[idx].Client
		if !skip(client) {
			cpm.currentIndex = (idx + 1) % len(cpm.active)

			return client, true
		}
	}

	// fallback to cooling pool if all active clients skipped
	for i := range cpm.cooling {
		client := cpm.cooling[i].Client
		if !skip(client) {
			// Track cooling pool usage
			cpm.logger.Warn().Str("client_url", client.GetEndpointURL()).Msg("using cooling pool client - all active clients skipped")
			return client, true
		}
	}

	return nil, false
}

// ReportHealth reports the health status of a client operation
func (cpm *ClientPoolManager) ReportHealth(client interfaces.Client, tries int, latencyMS float64, success bool) {
	cpm.mu.Lock()
	defer cpm.mu.Unlock()

	if tries == 0 {
		tries = 1
	}

	clientURL := client.GetEndpointURL()
	var clientInfo *ClientInfo
	var foundInActive bool
	var activeIndex int

	// Find the client in active pool
	for i := range cpm.active {
		if cpm.active[i].Client.GetEndpointURL() == clientURL {
			clientInfo = &cpm.active[i]
			foundInActive = true
			activeIndex = i
			break
		}
	}

	// If not in active, check cooling pool
	if !foundInActive {
		for i := range cpm.cooling {
			if cpm.cooling[i].Client.GetEndpointURL() == clientURL {
				clientInfo = &cpm.cooling[i]
				break
			}
		}
	}

	if clientInfo == nil {
		cpm.logger.Error().Str("client_url", clientURL).Str("protocol", client.GetProtocol()).Msg("client not found in pools - this indicates a bug")
		return
	}

	// Update success rate (same algorithm as NodeManager)
	clientInfo.successRate = math.Pow(0.8, float64(tries-1))

	if success {
		const alpha = 0.3
		if clientInfo.latEWMA == 0 {
			clientInfo.latEWMA = latencyMS
		} else {
			clientInfo.latEWMA = alpha*latencyMS + (1-alpha)*clientInfo.latEWMA
		}
	} else {
		clientInfo.successRate = 0
	}

	// Handle cooling/activation logic
	if foundInActive {
		if clientInfo.successRate < cpm.coolingThreshold {
			cpm.coolClientByIndex(activeIndex)
		} else {
			cpm.sortActive()
		}
	} else {
		cpm.sortCooling()
	}
}

// UpdateRateLimitDelay increases the rate limit delay for a client
func (cpm *ClientPoolManager) UpdateRateLimitDelay(client interfaces.Client) {
	cpm.mu.Lock()
	defer cpm.mu.Unlock()

	clientURL := client.GetEndpointURL()

	// Check active pool
	for i := range cpm.active {
		if cpm.active[i].Client.GetEndpointURL() == clientURL {
			cpm.active[i].rateLimitDelay = cpm.increaseRateLimit(cpm.active[i].rateLimitDelay)
			return
		}
	}

	// Check cooling pool
	for i := range cpm.cooling {
		if cpm.cooling[i].Client.GetEndpointURL() == clientURL {
			cpm.cooling[i].rateLimitDelay = cpm.increaseRateLimit(cpm.cooling[i].rateLimitDelay)
			return
		}
	}
}

// Helper methods (reused from NodeManager)
func (cpm *ClientPoolManager) increaseRateLimit(current time.Duration) time.Duration {
	newDelay := current + cpm.rateLimitDelayIncrease
	if newDelay > cpm.maxRateLimitDelay {
		return cpm.maxRateLimitDelay
	}
	return newDelay
}

func (cpm *ClientPoolManager) coolClientByIndex(index int) {
	if index >= len(cpm.active) {
		return
	}
	client := cpm.active[index]
	cpm.active = append(cpm.active[:index], cpm.active[index+1:]...)
	cpm.cooling = append(cpm.cooling, client)
	cpm.sortCooling()

	// Adjust currentIndex if necessary
	if cpm.currentIndex >= len(cpm.active) && len(cpm.active) > 0 {
		cpm.currentIndex = 0
	}
}

func (cpm *ClientPoolManager) sortActive() {
	// Use the exact same sorting algorithm as NodeManager.sortPool()
	sortClientPool(cpm.active)
}

func (cpm *ClientPoolManager) sortCooling() {
	// Use the exact same sorting algorithm as NodeManager.sortPool()
	sortClientPool(cpm.cooling)
}

// sortClientPool uses the exact same sorting algorithm as NodeManager.sortPool()
// Higher successRate is better; lower latency is better
func sortClientPool(clients []ClientInfo) {
	sort.Slice(clients, func(i, j int) bool {
		// Higher successRate is better; lower latency is better
		ci, cj := clients[i], clients[j]

		if ci.successRate != cj.successRate {
			return ci.successRate > cj.successRate
		}
		li, lj := ci.latEWMA, cj.latEWMA
		if li == 0 {
			li = 9e9 // Same as NodeManager - treat zero latency as infinite
		}
		if lj == 0 {
			lj = 9e9
		}
		return li < lj
	})
}

// healthLoop continuously probes cooling clients to check if they should be reactivated
// This matches the behavior of NodeManager.healthLoop()
func (cpm *ClientPoolManager) healthLoop() {
	tk := time.NewTicker(cpm.checkRate)
	defer tk.Stop()

	for range tk.C {
		cpm.probeCooling()
	}
}

// probeCooling checks each cooling client and reactivates healthy ones
func (cpm *ClientPoolManager) probeCooling() {
	cpm.mu.Lock()
	defer cpm.mu.Unlock()

	activatingIndexes := make([]int, 0)

	for i := 0; i < len(cpm.cooling); i++ {
		clientInfo := &cpm.cooling[i]
		if cpm.pingClient(clientInfo.Client) {
			clientInfo.healthStreak++
			cpm.logger.Debug().Str("client_url", clientInfo.Client.GetEndpointURL()).Int("streak", clientInfo.healthStreak).Msg("client health check passed")
			if clientInfo.healthStreak >= cpm.minActiveStreak {
				clientInfo.healthStreak = 0
				clientInfo.successRate = cpm.startReactivatedSuccRate
				activatingIndexes = append(activatingIndexes, i)
			}
		} else {
			clientInfo.healthStreak = 0
			cpm.logger.Debug().Str("client_url", clientInfo.Client.GetEndpointURL()).Msg("client health check failed")
		}
	}

	// Activate clients in reverse order to maintain slice integrity
	for i := len(activatingIndexes) - 1; i >= 0; i-- {
		cpm.activateClientByIndex(activatingIndexes[i])
	}
}

// pingClient performs a health check on the client
// For JSON-RPC clients, it calls /status endpoint
// For gRPC clients, it calls the Status method
func (cpm *ClientPoolManager) pingClient(client interfaces.Client) bool {
	ctx, cancel := context.WithTimeout(context.Background(), 4*time.Second)
	defer cancel()

	// Use the client's Status method for health checking
	err := client.Status(ctx)
	if err != nil {
		cpm.logger.Debug().Err(err).Str("client_url", client.GetEndpointURL()).Msg("client health check failed")
		return false
	}

	cpm.logger.Debug().Str("client_url", client.GetEndpointURL()).Msg("client health check passed")
	return true
}

// activateClientByIndex moves a client from cooling to active pool
func (cpm *ClientPoolManager) activateClientByIndex(index int) {
	if index >= len(cpm.cooling) {
		return
	}
	client := cpm.cooling[index]
	cpm.cooling = append(cpm.cooling[:index], cpm.cooling[index+1:]...)
	cpm.active = append(cpm.active, client)
	cpm.sortActive()

	cpm.logger.Info().Str("client_url", client.Client.GetEndpointURL()).Str("protocol", client.Client.GetProtocol()).Msg("reactivated client")
}

// Backoff management methods (from RoundRobinTransport)

// IsClientInBackoff returns the remaining backoff duration for a client
func (cpm *ClientPoolManager) IsClientInBackoff(client interfaces.Client) time.Duration {
	return cpm.backoffWait(client.GetEndpointURL())
}

// ApplyBackoffPenalty applies a backoff penalty to a client
func (cpm *ClientPoolManager) ApplyBackoffPenalty(client interfaces.Client, base time.Duration) time.Duration {
	return cpm.computeAndSetBackoff(client.GetEndpointURL(), base)
}

// ClearBackoff clears the backoff for a client if the window has expired
func (cpm *ClientPoolManager) ClearBackoff(client interfaces.Client) {
	cpm.clearBackoff(client.GetEndpointURL())
}

// GetShortestBackoff returns the shortest backoff duration across all clients
func (cpm *ClientPoolManager) GetShortestBackoff() time.Duration {
	return cpm.shortestBackoff()
}

// Private backoff methods (adapted from RoundRobinTransport)

func (cpm *ClientPoolManager) backoffWait(clientURL string) time.Duration {
	cpm.backMu.Lock()
	defer cpm.backMu.Unlock()

	s, ok := cpm.backoff[clientURL]
	if !ok {
		return 0
	}
	if time.Now().Before(s.until) {
		return time.Until(s.until)
	}
	return 0
}

func (cpm *ClientPoolManager) computeAndSetBackoff(clientURL string, base time.Duration) time.Duration {
	cpm.backMu.Lock()
	defer cpm.backMu.Unlock()

	now := time.Now()

	if bs, ok := cpm.backoff[clientURL]; ok && now.Before(bs.until) {
		return bs.until.Sub(now) // already penalised
	}

	fails := 0
	if bs, ok := cpm.backoff[clientURL]; ok {
		fails = bs.failures // keep exponential counter
	}

	if fails > 30 { // 2^30 seconds is already over 34 years
		fails = 30
	}

	delay := min(base*time.Duration(1<<fails), cpm.maxDelay)

	jitter := float64(delay) * cpm.jitterFrac * (2*rand.Float64() - 1)
	delay = time.Duration(float64(delay) + jitter)

	cpm.backoff[clientURL] = &backoffState{
		failures: fails + 1,
		until:    now.Add(delay),
	}

	cpm.logger.Debug().Str("client_url", clientURL).Dur("delay", delay).Int("failures", fails+1).Msg("applied backoff penalty")
	return delay
}

func (cpm *ClientPoolManager) clearBackoff(clientURL string) {
	cpm.backMu.Lock()
	if bs, ok := cpm.backoff[clientURL]; ok && time.Now().After(bs.until) {
		delete(cpm.backoff, clientURL)
		cpm.logger.Debug().Str("client_url", clientURL).Msg("cleared backoff")
	}
	cpm.backMu.Unlock()
}

func (cpm *ClientPoolManager) shortestBackoff() time.Duration {
	cpm.backMu.Lock()
	defer cpm.backMu.Unlock()

	var minDur time.Duration
	for _, st := range cpm.backoff {
		if d := time.Until(st.until); d > 0 {
			if minDur == 0 || d < minDur {
				minDur = d
			}
		}
	}
	return minDur
}

// GetMaxRetries returns the MaxRetries value for a given client
func (cpm *ClientPoolManager) GetMaxRetries(client interfaces.Client) int {
	cpm.mu.RLock()
	defer cpm.mu.RUnlock()

	clientURL := client.GetEndpointURL()

	// Check active pool first
	for i := range cpm.active {
		if cpm.active[i].Client.GetEndpointURL() == clientURL {
			return cpm.active[i].MaxRetries
		}
	}

	// Check cooling pool
	for i := range cpm.cooling {
		if cpm.cooling[i].Client.GetEndpointURL() == clientURL {
			return cpm.cooling[i].MaxRetries
		}
	}

	return 2 // Default fallback
}

// GetRateLimitDelay returns the current rate limit delay for a client
func (cpm *ClientPoolManager) GetRateLimitDelay(client interfaces.Client) time.Duration {
	cpm.mu.RLock()
	defer cpm.mu.RUnlock()

	clientURL := client.GetEndpointURL()

	// Check active pool first
	for i := range cpm.active {
		if cpm.active[i].Client.GetEndpointURL() == clientURL {
			return cpm.active[i].rateLimitDelay
		}
	}

	// Check cooling pool
	for i := range cpm.cooling {
		if cpm.cooling[i].Client.GetEndpointURL() == clientURL {
			return cpm.cooling[i].rateLimitDelay
		}
	}

	return 0
}

// ExpDelay calculates exponential delay for retry attempt i
func (cpm *ClientPoolManager) ExpDelay(attempt int) time.Duration {
	d := float64(cpm.baseDelay) * math.Pow(2, float64(attempt))
	if d > float64(cpm.maxDelay) {
		d = float64(cpm.maxDelay)
	}
	jitter := d * cpm.jitterFrac * (2*rand.Float64() - 1)
	return time.Duration(d + jitter)
}

// GetClientWithBackoff returns a client while respecting backoff states
// This is a convenience method that combines GetClient with backoff checking
func (cpm *ClientPoolManager) GetClientWithBackoff() (interfaces.Client, bool) {
	return cpm.GetClient(func(client interfaces.Client) bool {
		// Skip clients that are in backoff
		return cpm.IsClientInBackoff(client) > 0
	})
}

// GetHealthStatus returns a summary of the client pool health
func (cpm *ClientPoolManager) GetHealthStatus() map[string]any {
	cpm.mu.RLock()
	defer cpm.mu.RUnlock()

	status := make(map[string]any)
	status["active_clients"] = len(cpm.active)
	status["cooling_clients"] = len(cpm.cooling)

	if len(cpm.active) > 0 {
		status["current_client"] = cpm.active[cpm.currentIndex].Client.GetEndpointURL()
		status["current_client_protocol"] = cpm.active[cpm.currentIndex].Client.GetProtocol()
	}

	activeSummary := make([]map[string]any, len(cpm.active))
	for i, clientInfo := range cpm.active {
		activeSummary[i] = map[string]any{
			"url":          clientInfo.Client.GetEndpointURL(),
			"protocol":     clientInfo.Client.GetProtocol(),
			"success_rate": clientInfo.successRate,
			"latency_ms":   clientInfo.latEWMA,
			"max_retries":  clientInfo.MaxRetries,
		}
	}
	status["active"] = activeSummary

	coolingSummary := make([]map[string]any, len(cpm.cooling))
	for i, clientInfo := range cpm.cooling {
		coolingSummary[i] = map[string]any{
			"url":           clientInfo.Client.GetEndpointURL(),
			"protocol":      clientInfo.Client.GetProtocol(),
			"success_rate":  clientInfo.successRate,
			"latency_ms":    clientInfo.latEWMA,
			"health_streak": clientInfo.healthStreak,
			"max_retries":   clientInfo.MaxRetries,
		}
	}
	status["cooling"] = coolingSummary

	return status
}

// ExecuteWithRetry executes a function with automatic retry and load balancing
// This generic function provides type-safe retry logic with full pool management integration.
//
// Type parameters:
//   - Result: The expected return type (e.g., *authtypes.QueryAccountResponse)
//
// Parameters:
//   - ctx: Context for the operation
//   - poolManager: The client pool manager for load balancing and health tracking
//   - logger: Logger for debugging and monitoring
//   - operation: Function that receives an Client and returns the result
//
// Returns the result with full type safety, or an error if all clients fail.
func ExecuteWithRetry[Result any](
	ctx context.Context,
	poolManager *ClientPoolManager,
	logger *zerolog.Logger,
	operation func(client interfaces.Client) (Result, error),
) (_ Result, err error) {
	start := time.Now()
	var lastErr error

	// Try up to the number of available clients
	maxAttempts := len(poolManager.active) + len(poolManager.cooling)
	if maxAttempts == 0 {
		var zero Result
		return zero, fmt.Errorf("no clients available")
	}

	for attempt := range maxAttempts {
		// Get the next available client
		aggregatedClient, available := poolManager.GetClientWithBackoff()
		if !available {
			// If no clients available due to backoff, wait for the shortest backoff
			if backoffDuration := poolManager.GetShortestBackoff(); backoffDuration > 0 {
				logger.Debug().Dur("backoff", backoffDuration).Msg("waiting for client backoff")
				select {
				case <-time.After(backoffDuration):
					continue
				case <-ctx.Done():
					var zero Result
					return zero, ctx.Err()
				}
			}
			break
		}

		// Execute the operation with the Client
		// The operation function will extract the specific module client it needs
		result, operationErr := operation(aggregatedClient)
		duration := time.Since(start)

		if operationErr != nil {
			lastErr = operationErr
			logger.Debug().
				Err(operationErr).
				Str("endpoint", aggregatedClient.GetEndpointURL()).
				Int("attempt", attempt+1).
				Dur("duration", duration).
				Msg("operation failed, trying next client")

			// Report failure and apply backoff penalty
			poolManager.ReportHealth(aggregatedClient, attempt+1, float64(duration.Milliseconds()), false)
			poolManager.ApplyBackoffPenalty(aggregatedClient, poolManager.ExpDelay(attempt))
			continue
		}

		// Success - report health and return result
		poolManager.ReportHealth(aggregatedClient, attempt+1, float64(duration.Milliseconds()), true)
		poolManager.ClearBackoff(aggregatedClient)

		logger.Debug().
			Str("endpoint", aggregatedClient.GetEndpointURL()).
			Int("attempt", attempt+1).
			Dur("duration", duration).
			Msg("operation succeeded")

		return result, nil
	}

	if lastErr != nil {
		var zero Result
		return zero, fmt.Errorf("all clients failed, last error: %w", lastErr)
	}
	var zero Result
	return zero, fmt.Errorf("no clients available")
}
