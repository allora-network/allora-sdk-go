package interfaces

import (
    "github.com/allora-network/allora-sdk-go/pool"
)

type CosmosClientPool interface {
    Close() error

    {{ range .Modules }}
    {{ .ModuleName | title }}() {{ .ModuleName | title }}Client
    {{- end }}
}

type CosmosClient interface {
    CosmosClientPool
    pool.PoolParticipant
}

