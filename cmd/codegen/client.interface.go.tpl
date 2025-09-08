package interfaces

import (
    "context"
)

// Client defines the interface for protocol clients that aggregate multiple blockchain modules.
type Client interface {
    // Connection management methods required by ClientPoolManager
    GetEndpointURL() string
    GetProtocol() string
    Status(ctx context.Context) error
    Close() error
    {{ range .Modules }}
    {{ .ModuleName | title }}() {{ .ModuleName | title }}Client
    {{- end }}
}
