package wrapper

import (
	"context"

	authz "github.com/cosmos/cosmos-sdk/x/authz"
	"github.com/rs/zerolog"

	"github.com/allora-network/allora-sdk-go/config"
	"github.com/allora-network/allora-sdk-go/pool"
)

// AuthzClientWrapper wraps the authz module with pool management and retry logic
type AuthzClientWrapper struct {
	poolManager *pool.ClientPoolManager
	logger      zerolog.Logger
}

// NewAuthzClientWrapper creates a new authz client wrapper
func NewAuthzClientWrapper(poolManager *pool.ClientPoolManager, logger zerolog.Logger) *AuthzClientWrapper {
	return &AuthzClientWrapper{
		poolManager: poolManager,
		logger:      logger.With().Str("module", "authz").Logger(),
	}
}

func (c *AuthzClientWrapper) Grants(ctx context.Context, req *authz.QueryGrantsRequest, opts ...config.CallOpt) (*authz.QueryGrantsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*authz.QueryGrantsResponse, error) {
		return client.Authz().Grants(ctx, req, opts...)
	})
}

func (c *AuthzClientWrapper) GranterGrants(ctx context.Context, req *authz.QueryGranterGrantsRequest, opts ...config.CallOpt) (*authz.QueryGranterGrantsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*authz.QueryGranterGrantsResponse, error) {
		return client.Authz().GranterGrants(ctx, req, opts...)
	})
}

func (c *AuthzClientWrapper) GranteeGrants(ctx context.Context, req *authz.QueryGranteeGrantsRequest, opts ...config.CallOpt) (*authz.QueryGranteeGrantsResponse, error) {
	return pool.ExecuteWithRetry(ctx, c.poolManager, &c.logger, func(client pool.Client) (*authz.QueryGranteeGrantsResponse, error) {
		return client.Authz().GranteeGrants(ctx, req, opts...)
	})
}
