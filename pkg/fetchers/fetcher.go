package fetchers

import (
	"main/pkg/fetchers/cosmos"
	"main/pkg/types"

	"github.com/rs/zerolog"
)

type Fetcher interface {
	GetAllProposals() ([]types.Proposal, *types.QueryError)
	GetVote(proposal, voter string) (*types.Vote, *types.QueryError)
	GetTallies() (types.ChainTallyInfos, error)

	GetChainParams() (*types.ChainWithVotingParams, []error)
}

func GetFetcher(chainConfig *types.Chain, logger zerolog.Logger) Fetcher {
	return cosmos.NewRPC(chainConfig, logger)
}
