package rest

import (
	"bytes"
	"context"
	"encoding/json/v2"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/brynbellomy/go-utils/errors"
	cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
	grpctypes "github.com/cosmos/cosmos-sdk/types/grpc"
	"github.com/rs/zerolog"

	"github.com/cosmos/gogoproto/jsonpb"
	"github.com/cosmos/gogoproto/proto"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/gen/interfaces"
	"github.com/allora-network/allora-sdk-go/metrics"
)

// RESTClient implements the interfaces.Client interface using REST/JSON-RPC
type RESTClient struct {
	baseURL string
	logger  zerolog.Logger

	core         *RESTClientCore
	auth         *AuthRESTClient
	mint         *MintRESTClient
	evidence     *EvidenceRESTClient
	staking      *StakingRESTClient
	distribution *DistributionRESTClient
	emissions    *EmissionsRESTClient
	params       *ParamsRESTClient
	feegrant     *FeegrantRESTClient
	tx           *TxRESTClient
	bank         *BankRESTClient
	slashing     *SlashingRESTClient
	node         *NodeRESTClient
	authz        *AuthzRESTClient
	consensus    *ConsensusRESTClient
	tendermint   *TendermintRESTClient
	gov          *GovRESTClient
}

var _ interfaces.Client = (*RESTClient)(nil)

// NewRESTClient creates a new REST aggregated client
// Constructor takes a base URL and a logger
func NewRESTClient(baseURL string, logger zerolog.Logger, opts ...RESTClientOption) *RESTClient {
	core := NewRESTClientCore(baseURL, logger, opts...)

	return &RESTClient{
		baseURL:      baseURL,
		logger:       logger.With().Str("protocol", "json-rpc").Str("endpoint", baseURL).Logger(),
		core:         core,
		auth:         NewAuthRESTClient(core, logger),
		mint:         NewMintRESTClient(core, logger),
		evidence:     NewEvidenceRESTClient(core, logger),
		staking:      NewStakingRESTClient(core, logger),
		distribution: NewDistributionRESTClient(core, logger),
		emissions:    NewEmissionsRESTClient(core, logger),
		params:       NewParamsRESTClient(core, logger),
		feegrant:     NewFeegrantRESTClient(core, logger),
		tx:           NewTxRESTClient(core, logger),
		bank:         NewBankRESTClient(core, logger),
		slashing:     NewSlashingRESTClient(core, logger),
		node:         NewNodeRESTClient(core, logger),
		authz:        NewAuthzRESTClient(core, logger),
		consensus:    NewConsensusRESTClient(core, logger),
		tendermint:   NewTendermintRESTClient(core, logger),
		gov:          NewGovRESTClient(core, logger),
	}
}

type RESTClientOption func(*RESTClientCore)

// WithMetrics wraps the HTTP client's transport with a Prometheus metrics collector.
// Safe to call multiple times; the wrapper will wrap the current transport.
func WithMetrics() RESTClientOption {
	return func(c *RESTClientCore) {
		c.httpClient.Transport = metrics.NewMetricsRoundTripper(c.httpClient.Transport)
	}
}

// WithTimeout sets the HTTP client timeout.
func WithTimeout(d time.Duration) RESTClientOption {
	return func(c *RESTClientCore) {
		c.httpClient.Timeout = d
	}
}

// Close closes the client
func (c *RESTClient) Close() error {
	return nil
}

func (c *RESTClient) GetEndpointURL() string {
	return c.baseURL
}

func (c *RESTClient) GetProtocol() config.Protocol {
	return config.ProtocolREST
}

func (c *RESTClient) Auth() interfaces.AuthClient {
	return c.auth
}

func (c *RESTClient) Mint() interfaces.MintClient {
	return c.mint
}

func (c *RESTClient) Evidence() interfaces.EvidenceClient {
	return c.evidence
}

func (c *RESTClient) Staking() interfaces.StakingClient {
	return c.staking
}

func (c *RESTClient) Distribution() interfaces.DistributionClient {
	return c.distribution
}

func (c *RESTClient) Emissions() interfaces.EmissionsClient {
	return c.emissions
}

func (c *RESTClient) Params() interfaces.ParamsClient {
	return c.params
}

func (c *RESTClient) Feegrant() interfaces.FeegrantClient {
	return c.feegrant
}

func (c *RESTClient) Tx() interfaces.TxClient {
	return c.tx
}

func (c *RESTClient) Bank() interfaces.BankClient {
	return c.bank
}

func (c *RESTClient) Slashing() interfaces.SlashingClient {
	return c.slashing
}

func (c *RESTClient) Node() interfaces.NodeClient {
	return c.node
}

func (c *RESTClient) Authz() interfaces.AuthzClient {
	return c.authz
}

func (c *RESTClient) Consensus() interfaces.ConsensusClient {
	return c.consensus
}

func (c *RESTClient) Tendermint() interfaces.TendermintClient {
	return c.tendermint
}

func (c *RESTClient) Gov() interfaces.GovClient {
	return c.gov
}

// Status implements a basic health check using the Tendermint service
func (c *RESTClient) Status(ctx context.Context) error {
	_, err := c.tendermint.GetSyncing(ctx, &cmtservice.GetSyncingRequest{})
	return err
}

type RESTClientCore struct {
	baseURL    string
	httpClient *http.Client
	logger     zerolog.Logger

	marshaler   jsonpb.Marshaler
	unmarshaler jsonpb.Unmarshaler
}

func NewRESTClientCore(baseURL string, logger zerolog.Logger, opts ...RESTClientOption) *RESTClientCore {
	c := &RESTClientCore{
		baseURL: strings.TrimSuffix(baseURL, "/"),
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		logger: logger.With().Str("protocol", "rest").Str("endpoint", baseURL).Logger(),
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func (c *RESTClientCore) executeRequest(
	ctx context.Context,
	httpMethod, httpPath string,
	pathParams []string,
	queryParams []string,
	request proto.Message,
	response proto.Message,
	height int64,
) error {
	finalPath := httpPath
	if len(pathParams) > 0 {
		requestValue := reflect.ValueOf(request).Elem()
		for _, paramName := range pathParams {
			paramValue := c.getFieldValue(requestValue, paramName)
			placeholder := "{" + paramName + "}"
			finalPath = strings.ReplaceAll(finalPath, placeholder, paramValue)
		}
	}

	queryValues := url.Values{}
	if len(queryParams) > 0 && request != nil {
		requestValue := reflect.ValueOf(request).Elem()
		for _, paramName := range queryParams {
			paramValue := c.getFieldValue(requestValue, paramName)
			if paramValue != "" {
				queryValues.Set(paramName, paramValue)
			}
		}
	}

	fullURL := c.baseURL + finalPath
	if len(queryValues) > 0 {
		fullURL += "?" + queryValues.Encode()
	}

	c.logger.Debug().
		Str("method", httpMethod).
		Str("url", fullURL).
		Msg("Making JSON-RPC request")

	hasBody := httpMethod == "POST" && request != nil

	var reqBody bytes.Buffer
	if hasBody {
		err := c.marshaler.Marshal(&reqBody, request)
		if err != nil {
			return errors.Wrapf(err, "failed to marshal request")
		}
	}

	var req *http.Request
	var err error
	if hasBody {
		req, err = http.NewRequestWithContext(ctx, httpMethod, fullURL, &reqBody)
		req.Header.Set("Content-Type", "application/json")
	} else {
		req, err = http.NewRequestWithContext(ctx, httpMethod, fullURL, nil)
	}
	if err != nil {
		return errors.Wrapf(err, "failed to create request")
	}

	req.Header.Set("Accept", "application/json")
	if height > 0 {
		req.Header.Set(grpctypes.GRPCBlockHeightHeader, strconv.FormatInt(height, 10))
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return errors.Wrapf(err, "HTTP request failed")
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return errors.Errorf("HTTP error %d: %s", resp.StatusCode, resp.Status)
	}

	if response != nil {
		if err := c.unmarshaler.Unmarshal(resp.Body, response); err != nil {
			return errors.Wrapf(err, "failed to unmarshal response")
		}
	}
	return nil
}

// getFieldValue extracts a field value from a protobuf message using reflection
// Supports nested field access using dot notation (e.g., "pagination.limit")
func (c *RESTClientCore) getFieldValue(v reflect.Value, fieldPath string) string {
	parts := strings.Split(fieldPath, ".")
	currentValue := v

	for _, part := range parts {
		// Handle different field naming conventions
		fieldName := c.findFieldByName(currentValue, part)
		if fieldName == "" {
			return ""
		}

		field := currentValue.FieldByName(fieldName)
		if !field.IsValid() {
			return ""
		}

		// If it's a pointer, dereference it
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				return ""
			}
			field = field.Elem()
		}

		currentValue = field
	}

	return c.formatFieldValue(currentValue)
}

// findFieldByName finds a field by name, handling different naming conventions
func (c *RESTClientCore) findFieldByName(v reflect.Value, name string) string {
	t := v.Type()

	// Try exact match first
	if _, exists := t.FieldByName(name); exists {
		return name
	}

	// Try PascalCase version
	pascalName := strings.Title(name)
	if _, exists := t.FieldByName(pascalName); exists {
		return pascalName
	}

	// Try converting snake_case to PascalCase
	re := regexp.MustCompile(`_([a-z])`)
	camelName := re.ReplaceAllStringFunc(name, func(match string) string {
		return strings.ToUpper(match[1:])
	})
	camelName = strings.Title(camelName)
	if _, exists := t.FieldByName(camelName); exists {
		return camelName
	}

	return ""
}

// formatFieldValue converts a reflect.Value to its string representation
func (c *RESTClientCore) formatFieldValue(v reflect.Value) string {
	if !v.IsValid() {
		return ""
	}

	switch v.Kind() {
	case reflect.String:
		return v.String()
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return fmt.Sprintf("%d", v.Int())
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return fmt.Sprintf("%d", v.Uint())
	case reflect.Float32, reflect.Float64:
		return fmt.Sprintf("%g", v.Float())
	case reflect.Bool:
		return fmt.Sprintf("%t", v.Bool())
	case reflect.Slice:
		if v.Type().Elem().Kind() == reflect.Uint8 {
			// Handle []byte
			return string(v.Bytes())
		}
		// For other slices, could implement comma-separated values
		return ""
	default:
		// For complex types, try JSON marshaling
		if v.CanInterface() {
			if jsonBytes, err := json.Marshal(v.Interface()); err == nil {
				return string(jsonBytes)
			}
		}
		return ""
	}
}
