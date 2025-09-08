package interfaces

import (
    "context"
)

// Client defines the interface for protocol clients that aggregate multiple blockchain modules.
type Client interface {
    Close() error

    {{ range .Modules }}
    {{ .ModuleName | title }}() {{ .ModuleName | title }}Client
    {{- end }}
}
