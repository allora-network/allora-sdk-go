package wrapper

import (
    "context"

    "github.com/rs/zerolog"
{{if eq .ModuleName "emissions" "mint"}}    {{ .PackageName }} "{{ .ImportPath }}"{{else}}  {{ .PackageName }} "{{ .ImportPath }}"{{end}}

    "github.com/allora-network/allora-sdk-go/config"
    "github.com/allora-network/allora-sdk-go/gen/interfaces"
    "github.com/allora-network/allora-sdk-go/pool"
)

// {{ .ModuleName | title }}ClientWrapper wraps the {{ .ModuleName }} module with pool management and retry logic
type {{ .ModuleName | title }}ClientWrapper struct {
    poolManager *pool.ClientPoolManager[interfaces.CosmosClient]
    logger      zerolog.Logger
}

// New{{ .ModuleName | title }}ClientWrapper creates a new {{ .ModuleName }} client wrapper
func New{{ .ModuleName | title }}ClientWrapper(poolManager *pool.ClientPoolManager[interfaces.CosmosClient], logger zerolog.Logger) *{{ .ModuleName | title }}ClientWrapper {
    return &{{ .ModuleName | title }}ClientWrapper{
        poolManager: poolManager,
        logger:      logger.With().Str("module", "{{ .ModuleName }}").Logger(),
    }
}

{{range .Methods}}{{if .Comment}}// {{ .Comment }}{{end}}
func (c *{{ $.ModuleName | title }}ClientWrapper) {{ .Name }}(ctx context.Context, req *{{ $.PackageName }}.{{ .RequestType }}, opts ...config.CallOpt) (*{{ $.PackageName }}.{{ .ResponseType }}, error) {
    return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client interfaces.CosmosClient) (*{{ $.PackageName }}.{{ .ResponseType }}, error) {
        return client.{{ $.ModuleName | title }}().{{ .Name }}(ctx, req, opts...)
    })
}

{{end}}
