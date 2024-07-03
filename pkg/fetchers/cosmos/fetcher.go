package cosmos

import (
	"context"
	"main/pkg/fetchers/cosmos/responses"
	"main/pkg/http"
	"main/pkg/types"

	"go.opentelemetry.io/otel/trace"

	"github.com/rs/zerolog"
)

const PaginationLimit = 1000

type RPC struct {
	ChainConfig     *types.Chain
	ProposalsType   string
	Client          *http.Client
	Logger          zerolog.Logger
	PaginationLimit int
}

func NewRPC(chainConfig *types.Chain, logger *zerolog.Logger, tracer trace.Tracer) *RPC {
	return &RPC{
		ChainConfig:     chainConfig,
		ProposalsType:   chainConfig.ProposalsType,
		Logger:          logger.With().Str("component", "rpc").Logger(),
		Client:          http.NewClient(chainConfig.Name, chainConfig.LCDEndpoints, logger, tracer),
		PaginationLimit: PaginationLimit,
	}
}

func (rpc *RPC) GetAllProposals(prevHeight int64, ctx context.Context) ([]types.Proposal, int64, *types.QueryError) {
	if rpc.ProposalsType == "v1" {
		return rpc.GetAllV1Proposals(prevHeight, ctx)
	}

	return rpc.GetAllV1beta1Proposals(prevHeight, ctx)
}

func (rpc *RPC) GetStakingPool(ctx context.Context) (*responses.PoolRPCResponse, *types.QueryError) {
	url := "/cosmos/staking/v1beta1/pool"

	var pool responses.PoolRPCResponse
	if errs := rpc.Client.Get(url, &pool, ctx); len(errs) > 0 {
		return nil, &types.QueryError{
			QueryError: nil,
			NodeErrors: errs,
		}
	}

	return &pool, nil
}
