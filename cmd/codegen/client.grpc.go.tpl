package grpc

import (
    "context"
    "crypto/tls"
    "strconv"
    "strings"
    "sync"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc/encoding"
    "google.golang.org/grpc/metadata"
    "github.com/rs/zerolog"
    sdkgrpc "github.com/cosmos/cosmos-sdk/types/grpc"
    "github.com/brynbellomy/go-utils/errors"

    cmtservice "github.com/cosmos/cosmos-sdk/client/grpc/cmtservice"
    cosmoscodec "github.com/cosmos/cosmos-sdk/codec"
    codectypes "github.com/cosmos/cosmos-sdk/codec/types"
    "github.com/cosmos/cosmos-sdk/std"
    banktypes "github.com/cosmos/cosmos-sdk/x/bank/types"
    distributiontypes "github.com/cosmos/cosmos-sdk/x/distribution/types"
    slashingtypes "github.com/cosmos/cosmos-sdk/x/slashing/types"
    stakingtypes "github.com/cosmos/cosmos-sdk/x/staking/types"

    emissionstypes "github.com/allora-network/allora-chain/x/emissions/types"
    minttypes "github.com/allora-network/allora-chain/x/mint/types"

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

var (
    grpcCodecOnce sync.Once
    grpcCodec     encoding.Codec
)

func buildGRPCCodec() encoding.Codec {
    grpcCodecOnce.Do(func() {
        registry := codectypes.NewInterfaceRegistry()
        registerFuncs := []func(codectypes.InterfaceRegistry){
            std.RegisterInterfaces,
            banktypes.RegisterInterfaces,
            stakingtypes.RegisterInterfaces,
            slashingtypes.RegisterInterfaces,
            distributiontypes.RegisterInterfaces,
            minttypes.RegisterInterfaces,
            emissionstypes.RegisterInterfaces,
        }
        for _, register := range registerFuncs {
            register(registry)
        }
        grpcCodec = cosmoscodec.NewProtoCodec(registry).GRPCCodec()
    })
    return grpcCodec
}

// NewGRPCClient creates a new gRPC aggregated client
func NewGRPCClient(cfg config.EndpointConfig, logger zerolog.Logger) (*GRPCClient, error) {
    if cfg.Protocol != "grpc" {
        return nil, errors.Errorf("unsupported protocol: %s, only 'grpc' is supported", cfg.Protocol)
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
        grpc.WithDefaultCallOptions(), //grpc.ForceCodec(buildGRPCCodec())),
    )
    if err != nil {
        return nil, errors.Errorf("failed to connect to %s: %w", address, err)
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

func (c *GRPCClient) GetProtocol() config.Protocol {
    return config.ProtocolGRPC
}

func (c *GRPCClient) GetEndpointURL() string {
    return c.endpointURL
}

{{- range .Modules }}
func (c *GRPCClient) {{ .ModuleName | title }}() interfaces.{{ .ModuleName | title }}Client {
    return c.{{ .ModuleName }}
}
{{- end }}

type callFn[In, Out any] func(ctx context.Context, i In, opts ...grpc.CallOption) (Out, error)

func queryWithHeight[In any, Out any](ctx context.Context, height int64, queryFn callFn[In, Out], in In) (Out, error) {
    var zero Out

    if height > 0 {
        ctx = metadata.AppendToOutgoingContext(ctx, sdkgrpc.GRPCBlockHeightHeader, strconv.FormatInt(height, 10))
    }

    resp, err := queryFn(ctx, in)
    if err != nil {
        return zero, err
    }
    return resp, nil
}

// Status implements the Status method required by ClientPoolManager
func (c *GRPCClient) Status(ctx context.Context) error {
    _, err := c.tendermint.GetSyncing(ctx, &cmtservice.GetSyncingRequest{})
    return err
}
