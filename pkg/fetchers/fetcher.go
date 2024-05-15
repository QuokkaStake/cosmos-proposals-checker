package fetchers

import (
	"context"
	"main/pkg/fetchers/cosmos"
	"main/pkg/fetchers/neutron"
	"main/pkg/types"

	"go.opentelemetry.io/otel/trace"

	"github.com/rs/zerolog"
)

type Fetcher interface {
	GetAllProposals(prevHeight int64, ctx context.Context) ([]types.Proposal, int64, *types.QueryError)
	GetVote(proposal, voter string, prevHeight int64, ctx context.Context) (*types.Vote, int64, *types.QueryError)
	GetTallies(ctx context.Context) (types.ChainTallyInfos, error)

	GetChainParams(ctx context.Context) (*types.ChainWithVotingParams, []error)
}

func GetFetcher(
	chainConfig *types.Chain,
	logger *zerolog.Logger,
	tracer trace.Tracer,
) Fetcher {
	if chainConfig.Type == "neutron" {
		return neutron.NewFetcher(chainConfig, logger, tracer)
	}

	return cosmos.NewRPC(chainConfig, logger, tracer)
}
