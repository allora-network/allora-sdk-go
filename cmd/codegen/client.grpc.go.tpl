package grpc

import (
    "context"
    "crypto/tls"
    "fmt"
    "strings"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/credentials/insecure"
    "github.com/rs/zerolog"

    cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"

    "github.com/allora-network/allora-sdk-go/config"
    "github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// GRPCClient implements the Client interface using gRPC
type GRPCClient struct {
    endpointURL string
    conn        *grpc.ClientConn

    {{- range .Modules }}
    {{ .ModuleName }} *{{ .ModuleName | title }}GRPCClient
    {{- end }}
}

var _ interfaces.Client = (*GRPCClient)(nil)

// NewGRPCClient creates a new gRPC aggregated client
func NewGRPCClient(cfg config.EndpointConfig, logger zerolog.Logger) (*GRPCClient, error) {
    if cfg.Protocol != "grpc" {
        return nil, fmt.Errorf("unsupported protocol: %s, only 'grpc' is supported", cfg.Protocol)
    }

    // Parse the URL to determine if TLS is needed and get the address
    var creds credentials.TransportCredentials
    var address string

    if strings.HasPrefix(cfg.URL, "grpcs://") {
        // Secure gRPC with TLS
        address = strings.TrimPrefix(cfg.URL, "grpcs://")
        creds = credentials.NewTLS(&tls.Config{})
    } else if strings.HasPrefix(cfg.URL, "grpc://") {
        // Check if port 443 (typically TLS)
        address = strings.TrimPrefix(cfg.URL, "grpc://")
        if strings.Contains(address, ":443") {
            // Port 443 or .network domains typically require TLS
            creds = credentials.NewTLS(&tls.Config{})
        } else {
            creds = insecure.NewCredentials()
        }
    } else {
        // Default to the URL as-is, assume TLS for port 443
        address = cfg.URL
        if strings.Contains(address, ":443") {
            creds = credentials.NewTLS(&tls.Config{})
        } else {
            creds = insecure.NewCredentials()
        }
    }

    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    conn, err := grpc.DialContext(ctx, address,
        grpc.WithTransportCredentials(creds),
        grpc.WithBlock(),
    )
    if err != nil {
        return nil, fmt.Errorf("failed to connect to %s: %w", address, err)
    }

    client := &GRPCClient{
        endpointURL:      cfg.URL,
        conn:             conn,
        {{- range .Modules }}
        {{ .ModuleName }}: New{{ .ModuleName | title }}GRPCClient(conn, logger),
        {{- end }}
    }

    return client, nil
}

// Close closes the gRPC connection
func (c *GRPCClient) Close() error {
    return c.conn.Close()
}

{{- range .Modules }}
func (c *GRPCClient) {{ .ModuleName | title }}() interfaces.{{ .ModuleName | title }}Client {
    return c.{{ .ModuleName }}
}
{{- end }}

func (c *GRPCClient) GetEndpointURL() string {
    return c.endpointURL
}

func (c *GRPCClient) GetProtocol() string {
    return "grpc"
}

// Status implements the Status method required by ClientPoolManager
func (c *GRPCClient) Status(ctx context.Context) error {
    _, err := c.tendermint.GetSyncing(ctx, &cmtservice.GetSyncingRequest{})
    return err
}

