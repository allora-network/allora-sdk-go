package grpc

import (
    "context"

    "github.com/brynbellomy/go-utils/errors"
    "github.com/rs/zerolog"
    "google.golang.org/grpc"

{{if eq .ModuleName "emissions" "mint"}}    {{ .PackageName }} "{{ .ImportPath }}"{{else}}    {{ .PackageName }} "{{ .ImportPath }}"{{end}}

    "github.com/allora-network/allora-sdk-go/gen/interfaces"
)

// {{ .ModuleName | title }}GRPCClient provides gRPC access to the {{ .ModuleName }} module
type {{ .ModuleName | title }}GRPCClient struct {
    client {{ .PackageName }}.{{ .ServiceName }}Client
    logger zerolog.Logger
}

var _ interfaces.{{ .ModuleName | title }}Client = (*{{ .ModuleName | title }}GRPCClient)(nil)

// New{{ .ModuleName | title }}GRPCClient creates a new {{ .ModuleName }} REST client
func New{{ .ModuleName | title }}GRPCClient(conn *grpc.ClientConn, logger zerolog.Logger) *{{ .ModuleName | title }}GRPCClient {
    return &{{ .ModuleName | title }}GRPCClient{
        client: {{ .PackageName }}.New{{ .ServiceName }}Client(conn),
        logger: logger.With().Str("module", "{{ .ModuleName }}").Str("protocol", "grpc").Logger(),
    }
}

{{range .Methods}}
func (c *{{ $.ModuleName | title }}GRPCClient) {{ .Name }}(ctx context.Context, req *{{ $.PackageName }}.{{ .RequestType }}, opts ...config.CallOpt) (*{{ $.PackageName }}.{{ .ResponseType }}, error) {
    callOpts := config.DefaultCallOpts()
    callOpts.Apply(opts...)

    resp, err := queryWithHeight(ctx, callOpts.Height, c.client.{{ .Name }}, req)
    if err != nil {
        return resp, errors.Wrap(err, "while calling {{ $.ModuleName | title }}GRPCClient.{{ .Name }}")
    }
    return resp, nil
}
{{end}}