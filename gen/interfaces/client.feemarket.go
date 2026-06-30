package interfaces

import (
	"context"

	feemarkettypes "github.com/skip-mev/feemarket/x/feemarket/types"

	"github.com/allora-network/allora-sdk-go/config"
)

type FeemarketClient interface {
	GasPrice(ctx context.Context, req *feemarkettypes.GasPriceRequest, opts ...config.CallOpt) (*feemarkettypes.GasPriceResponse, error)
	GasPrices(ctx context.Context, req *feemarkettypes.GasPricesRequest, opts ...config.CallOpt) (*feemarkettypes.GasPricesResponse, error)
	Params(ctx context.Context, req *feemarkettypes.ParamsRequest, opts ...config.CallOpt) (*feemarkettypes.ParamsResponse, error)
	State(ctx context.Context, req *feemarkettypes.StateRequest, opts ...config.CallOpt) (*feemarkettypes.StateResponse, error)
}
