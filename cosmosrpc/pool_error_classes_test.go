package cosmosrpc

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	emissionstypes "github.com/allora-network/allora-chain/x/emissions/types"
	"github.com/rs/zerolog"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/allora-network/allora-sdk-go/gen/rest"
)

// restEndpoint starts a REST endpoint that answers every request with the given
// status. It returns the endpoint URL and a counter of the requests it
// received.
func restEndpoint(t *testing.T, httpStatus int) (string, *atomic.Int64) {
	t.Helper()

	var hits atomic.Int64
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		hits.Add(1)
		w.WriteHeader(httpStatus)
	}))
	t.Cleanup(server.Close)

	return server.URL, &hits
}

// TestPoolRESTStatusClasses checks status mapping, caller-visible gRPC codes,
// attempts and health together. A cooling-only assertion could pass even if
// these statuses reached the pool as uncoded errors.
func TestPoolRESTStatusClasses(t *testing.T) {
	for _, tc := range []struct {
		name string
		// httpStatus is what the first endpoint answers.
		httpStatus int
		// wantAttempts is how many endpoints the pool calls.
		wantAttempts int64
		// wantCooled is how many endpoints end up in the cooling pool.
		wantCooled int
		// wantCode is the gRPC status supplied to the caller.
		wantCode codes.Code
	}{
		// An unsupported method ends the call without a health penalty.
		{"501 unimplemented", http.StatusNotImplemented, 1, 0, codes.Unimplemented},
		// A provider's gateway rejects our credentials on every call until the
		// config is fixed, so the endpoint is persistently unusable to us.
		{"401 unauthenticated", http.StatusUnauthorized, 2, 2, codes.Unauthenticated},
		{"403 permission denied", http.StatusForbidden, 2, 2, codes.PermissionDenied},
		// The node or its proxy is rate limiting us.
		{"429 rate limited", http.StatusTooManyRequests, 2, 2, codes.ResourceExhausted},
		// The node is down or too slow.
		{"503 unavailable", http.StatusServiceUnavailable, 2, 2, codes.Unavailable},
		{"504 gateway timeout", http.StatusGatewayTimeout, 2, 2, codes.DeadlineExceeded},
		// Ambiguous statuses stay uncoded and count as endpoint faults.
		{"500 internal", http.StatusInternalServerError, 2, 2, codes.Unknown},
		{"502 bad gateway", http.StatusBadGateway, 2, 2, codes.Unknown},
		{"404 not found", http.StatusNotFound, 2, 2, codes.Unknown},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logger := zerolog.Nop()
			urlA, hitsA := restEndpoint(t, tc.httpStatus)
			urlB, hitsB := restEndpoint(t, tc.httpStatus)

			p := NewClientPool([]Client{
				rest.NewRESTClient(urlA, logger),
				rest.NewRESTClient(urlB, logger),
			}, logger)
			t.Cleanup(func() { _ = p.Close() })

			_, err := p.Emissions().IsTopicWorkerWhitelistEnabled(context.Background(),
				&emissionstypes.IsTopicWorkerWhitelistEnabledRequest{TopicId: 1})

			require.Error(t, err)
			require.Equal(t, tc.wantCode, status.Code(err), "caller status for HTTP %d", tc.httpStatus)
			require.Equal(t, tc.wantAttempts, hitsA.Load()+hitsB.Load(),
				"attempts for HTTP %d", tc.httpStatus)

			health := p.GetHealthStatus()
			require.Equal(t, tc.wantCooled, health["cooling_clients"],
				"cooling pool for HTTP %d", tc.httpStatus)
			require.Equal(t, 2-tc.wantCooled, health["active_clients"],
				"active pool for HTTP %d", tc.httpStatus)
		})
	}
}

// TestPoolRESTRateLimitFailsOverToHealthyEndpoint exercises a limited first
// endpoint under either initial round-robin position. The caller gets the
// healthy endpoint's answer, with exactly one request to each endpoint.
func TestPoolRESTRateLimitFailsOverToHealthyEndpoint(t *testing.T) {
	logger := zerolog.Nop()

	limitedURL, limitedHits := restEndpoint(t, http.StatusTooManyRequests)

	var healthyHits atomic.Int64
	healthy := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		healthyHits.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"is_topic_worker_whitelist_enabled":true}`))
	}))
	t.Cleanup(healthy.Close)

	p := NewClientPool([]Client{
		rest.NewRESTClient(limitedURL, logger),
		rest.NewRESTClient(healthy.URL, logger),
	}, logger)
	t.Cleanup(func() { _ = p.Close() })

	// Advance the healthy endpoint when it would go first. The following call
	// always starts on the rate-limited endpoint and exercises failover.
	if p.GetHealthStatus()["current_client"] == healthy.URL {
		response, err := p.Emissions().IsTopicWorkerWhitelistEnabled(context.Background(),
			&emissionstypes.IsTopicWorkerWhitelistEnabledRequest{TopicId: 1})
		require.NoError(t, err)
		require.True(t, response.IsTopicWorkerWhitelistEnabled)
	}
	resp, err := p.Emissions().IsTopicWorkerWhitelistEnabled(context.Background(),
		&emissionstypes.IsTopicWorkerWhitelistEnabledRequest{TopicId: 1})

	require.NoError(t, err, "a rate-limited endpoint must not fail the request when another endpoint serves it")
	require.True(t, resp.IsTopicWorkerWhitelistEnabled)
	require.Equal(t, int64(1), limitedHits.Load())
	require.Equal(t, int64(1), healthyHits.Load())
}

// TestPoolRESTUnimplementedIsTerminal pins that a 501 ends the request on the
// first endpoint with the code the node meant, so the caller can tell "this
// node does not serve this method" from "this node is broken".
func TestPoolRESTUnimplementedIsTerminal(t *testing.T) {
	logger := zerolog.Nop()
	urlA, hitsA := restEndpoint(t, http.StatusNotImplemented)
	urlB, hitsB := restEndpoint(t, http.StatusNotImplemented)

	p := NewClientPool([]Client{
		rest.NewRESTClient(urlA, logger),
		rest.NewRESTClient(urlB, logger),
	}, logger)
	t.Cleanup(func() { _ = p.Close() })

	_, err := p.Emissions().IsTopicWorkerWhitelistEnabled(context.Background(),
		&emissionstypes.IsTopicWorkerWhitelistEnabledRequest{TopicId: 1})

	require.Error(t, err)
	require.Equal(t, codes.Unimplemented, status.Code(err), "got %v", err)
	require.Equal(t, int64(1), hitsA.Load()+hitsB.Load(),
		"an unimplemented method must not be retried on the other endpoints")
	require.Equal(t, 2, p.GetHealthStatus()["active_clients"],
		"a 501 says nothing about endpoint health")
}

// TestPoolRESTMessageErrorClasses contrasts a gateway-generated message-size
// error with a proxy-generated bare 429. The former must return immediately,
// while the latter rate-limits both endpoints and preserves ResourceExhausted.
func TestPoolRESTMessageErrorClasses(t *testing.T) {
	for _, tc := range []struct {
		name         string
		body         string
		wantAttempts int
		wantActive   int
	}{
		{"message too large", `{"code":8,"message":"grpc: received message larger than max (5000000 vs. 4194304)"}`, 1, 2},
		{"bare proxy rate limit", "", 2, 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			logger := zerolog.Nop()
			newEndpoint := func() (*httptest.Server, *atomic.Int64) {
				var hits atomic.Int64
				server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
					hits.Add(1)
					w.Header().Set("Content-Type", "application/json")
					w.WriteHeader(http.StatusTooManyRequests)
					_, _ = w.Write([]byte(tc.body))
				}))
				t.Cleanup(server.Close)
				return server, &hits
			}
			serverA, hitsA := newEndpoint()
			serverB, hitsB := newEndpoint()

			p := NewClientPool([]Client{
				rest.NewRESTClient(serverA.URL, logger),
				rest.NewRESTClient(serverB.URL, logger),
			}, logger)
			t.Cleanup(func() { _ = p.Close() })
			start := time.Now()

			_, err := p.Emissions().IsTopicWorkerWhitelistEnabled(context.Background(),
				&emissionstypes.IsTopicWorkerWhitelistEnabledRequest{TopicId: 1})

			require.Error(t, err)
			require.Equal(t, codes.ResourceExhausted, status.Code(err))
			require.Equal(t, tc.wantAttempts, int(hitsA.Load()+hitsB.Load()))
			require.Equal(t, tc.wantActive, p.GetHealthStatus()["active_clients"])
			require.Equal(t, 2-tc.wantActive, p.GetHealthStatus()["cooling_clients"])
			if tc.body != "" {
				require.Less(t, time.Since(start), time.Second,
					"a message-size answer must return without a backoff wait")
				require.Contains(t, status.Convert(err).Message(), "larger than max")
			} else {
				require.Contains(t, status.Convert(err).Message(), "HTTP error 429")
			}
		})
	}
}
