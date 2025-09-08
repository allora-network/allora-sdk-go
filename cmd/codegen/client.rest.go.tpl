package rest

import (
    "bytes"
    "context"
    "encoding/json/jsontext"
    "encoding/json/v2"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "reflect"
    "regexp"
    "strconv"
    "strings"
    "time"

    cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
    stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"
    "github.com/rs/zerolog"

    "github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// RESTClient implements the interfaces.Client interface using REST/JSON-RPC
type RESTClient struct {
    baseURL string
    logger  zerolog.Logger

    // Core HTTP client (shared)
    core *RESTClientCore

    // Module clients (generated)
    {{- range .Modules }}
    {{ .ModuleName }} *{{ .ModuleName | title }}RESTClient
    {{- end }}
}

var _ interfaces.Client = (*RESTClient)(nil)

// NewRESTClient creates a new REST aggregated client
// Constructor takes a base URL and a logger
func NewRESTClient(baseURL string, logger zerolog.Logger) *RESTClient {
    core := NewRESTClientCore(baseURL, logger)

    return &RESTClient{
        baseURL: baseURL,
        logger:  logger.With().Str("protocol", "json-rpc").Str("endpoint", baseURL).Logger(),
        core:    core,
        {{- range .Modules }}
        {{ .ModuleName }}: New{{ .ModuleName | title }}RESTClient(core, logger),
        {{- end }}
    }
}

// Close closes the client (no-op for HTTP clients)
func (c *RESTClient) Close() error {
    return c.core.Close()
}

{{- range .Modules }}
func (c *RESTClient) {{ .ModuleName | title }}() interfaces.{{ .ModuleName | title }}Client {
    return c.{{ .ModuleName }}
}
{{ end }}

func (c *RESTClient) GetEndpointURL() string {
    return c.baseURL
}

func (c *RESTClient) GetProtocol() string {
    return "json-rpc"
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
}

func NewRESTClientCore(baseURL string, logger zerolog.Logger) *RESTClientCore {
    return &RESTClientCore{
        baseURL: strings.TrimSuffix(baseURL, "/"),
        httpClient: &http.Client{
            Timeout: 30 * time.Second,
        },
        logger: logger,
    }
}

var cosmosUnmarshalers = json.JoinUnmarshalers(
    // Handle time.Duration from strings like "1814400s"
    json.UnmarshalFromFunc(func(dec *jsontext.Decoder, val *time.Duration) error {
        tok, err := dec.ReadToken()
        if err != nil {
            return err
        }
        duration, err := time.ParseDuration(tok.String())
        if err != nil {
            return fmt.Errorf("invalid duration format %q: %w", tok.String(), err)
        }
        *val = duration
        return nil
    }),

    json.UnmarshalFromFunc(func(dec *jsontext.Decoder, val *int64) error {
        tok, err := dec.ReadToken()
        if err != nil {
            return err
        }
        i, err := strconv.ParseInt(tok.String(), 10, 64)
        if err != nil {
            return fmt.Errorf("invalid int64 format %q: %w", tok.String(), err)
        }
        *val = i
        return nil
    }),

    json.UnmarshalFromFunc(func(dec *jsontext.Decoder, val *uint64) error {
        tok, err := dec.ReadToken()
        if err != nil {
            return err
        }
        i, err := strconv.ParseUint(tok.String(), 10, 64)
        if err != nil {
            return fmt.Errorf("invalid uint64 format %q: %w", tok.String(), err)
        }
        *val = i
        return nil
    }),

    // Handle BondStatus enum from strings like "BOND_STATUS_BONDED"
    json.UnmarshalFromFunc(func(dec *jsontext.Decoder, val *stakingtypes.BondStatus) error {
        tok, err := dec.ReadToken()
        if err != nil {
            return err
        }
        switch tok.String() {
        case "BOND_STATUS_UNSPECIFIED":
            *val = stakingtypes.Unspecified
        case "BOND_STATUS_UNBONDED":
            *val = stakingtypes.Unbonded
        case "BOND_STATUS_UNBONDING":
            *val = stakingtypes.Unbonding
        case "BOND_STATUS_BONDED":
            *val = stakingtypes.Bonded
        default:
            return fmt.Errorf("unknown bond status: %q", tok.String())
        }
        return nil
    }),
)

func (c *RESTClientCore) executeRequest(ctx context.Context, httpMethod, httpPath string, pathParams []string, queryParams []string, request interface{}, response interface{}) error {
    // Build the URL path by replacing path parameters
    finalPath := httpPath
    if len(pathParams) > 0 {
        requestValue := reflect.ValueOf(request).Elem()
        for _, paramName := range pathParams {
            paramValue := c.getFieldValue(requestValue, paramName)
            placeholder := "{" + paramName + "}"
            finalPath = strings.ReplaceAll(finalPath, placeholder, paramValue)
        }
    }

    // Build query parameters
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

    // Construct full URL
    fullURL := c.baseURL + finalPath
    if len(queryValues) > 0 {
        fullURL += "?" + queryValues.Encode()
    }

    c.logger.Debug().
        Str("method", httpMethod).
        Str("url", fullURL).
        Msg("Making JSON-RPC request")

    // Create HTTP request
    var reqBody io.Reader
    if httpMethod == "POST" && request != nil {
        jsonData, err := json.Marshal(request)
        if err != nil {
            return fmt.Errorf("failed to marshal request: %w", err)
        }
        reqBody = bytes.NewReader(jsonData)
    }

    req, err := http.NewRequestWithContext(ctx, httpMethod, fullURL, reqBody)
    if err != nil {
        return fmt.Errorf("failed to create request: %w", err)
    }

    if reqBody != nil {
        req.Header.Set("Content-Type", "application/json")
    }
    req.Header.Set("Accept", "application/json")

    // Execute HTTP request
    resp, err := c.httpClient.Do(req)
    if err != nil {
        return fmt.Errorf("HTTP request failed: %w", err)
    }
    defer resp.Body.Close()

    // Read response body
    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return fmt.Errorf("failed to read response: %w", err)
    }

    // Check HTTP status
    if resp.StatusCode != http.StatusOK {
        return fmt.Errorf("HTTP error %d: %s", resp.StatusCode, string(respBody))
    }

    if response != nil {
        if err := json.Unmarshal(respBody, response, json.WithUnmarshalers(cosmosUnmarshalers)); err != nil {
            return fmt.Errorf("failed to unmarshal response: %w", err)
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

func (c *RESTClientCore) Close() error {
    return nil
}

func (c *RESTClientCore) GetEndpointURL() string {
    return c.baseURL
}

func (c *RESTClientCore) GetProtocol() string {
    return "json-rpc"
}
