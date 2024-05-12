package fetchers

import (
	"main/pkg/fetchers/cosmos"
	"main/pkg/fetchers/neutron"
	"main/pkg/types"

	"github.com/rs/zerolog"
)

type Fetcher interface {
	GetAllProposals(prevHeight int64) ([]types.Proposal, int64, *types.QueryError)
	GetVote(proposal, voter string, prevHeight int64) (*types.Vote, int64, *types.QueryError)
	GetTallies() (types.ChainTallyInfos, error)

	GetChainParams() (*types.ChainWithVotingParams, []error)
}

func GetFetcher(chainConfig *types.Chain, logger *zerolog.Logger) Fetcher {
	if chainConfig.Type == "neutron" {
		return neutron.NewFetcher(chainConfig, logger)
	}

	return cosmos.NewRPC(chainConfig, logger)
}
