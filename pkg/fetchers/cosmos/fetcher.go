package cosmos

import (
	"main/pkg/fetchers/cosmos/responses"
	"main/pkg/http"
	"main/pkg/types"

	"github.com/rs/zerolog"
)

const PaginationLimit = 1000

type RPC struct {
	ChainConfig   *types.Chain
	ProposalsType string
	Client        *http.Client
	Logger        zerolog.Logger
}

func NewRPC(chainConfig *types.Chain, logger zerolog.Logger) *RPC {
	return &RPC{
		ChainConfig:   chainConfig,
		ProposalsType: chainConfig.ProposalsType,
		Logger:        logger.With().Str("component", "rpc").Logger(),
		Client:        http.NewClient(chainConfig.Name, chainConfig.LCDEndpoints, logger),
	}
}

func (rpc *RPC) GetAllProposals() ([]types.Proposal, *types.QueryError) {
	if rpc.ProposalsType == "v1" {
		return rpc.GetAllV1Proposals()
	}

	return rpc.GetAllV1beta1Proposals()
}

func (rpc *RPC) GetStakingPool() (*responses.PoolRPCResponse, *types.QueryError) {
	url := "/cosmos/staking/v1beta1/pool"

	var pool responses.PoolRPCResponse
	if errs := rpc.Client.Get(url, &pool); len(errs) > 0 {
		return nil, &types.QueryError{
			QueryError: nil,
			NodeErrors: errs,
		}
	}

	return &pool, nil
}
