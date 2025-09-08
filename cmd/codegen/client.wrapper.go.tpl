package wrapper

import (
    "github.com/rs/zerolog"

    "github.com/allora-network/allora-sdk-go/gen/interfaces"
    "github.com/allora-network/allora-sdk-go/pool"
)

type WrapperClient struct {
    {{- range .Modules }}
    {{ .ModuleName }} *{{ .ModuleName | title }}ClientWrapper
    {{- end }}
}

func NewWrapperClient(poolManager *pool.ClientPoolManager, logger zerolog.Logger) *WrapperClient {
    return &WrapperClient{
        {{- range .Modules }}
        {{ .ModuleName }}: New{{ .ModuleName | title }}ClientWrapper(poolManager, logger),
        {{- end}}
    }
}

func (c *WrapperClient) Close() error {
    return nil
}

{{- range .Modules }}
func (c *WrapperClient) {{ .ModuleName | title }}() interfaces.{{ .ModuleName | title }}Client {
    return c.{{ .ModuleName }}
}
{{ end }}
