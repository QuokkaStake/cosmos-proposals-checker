package fetchers

import (
	"errors"
	"main/pkg/types"
)

type TestFetcher struct {
	WithProposals      bool
	WithProposalsError bool
}

func (f *TestFetcher) GetAllProposals(
	prevHeight int64,
) ([]types.Proposal, int64, *types.QueryError) {
	if f.WithProposalsError {
		return []types.Proposal{}, 123, &types.QueryError{
			QueryError: errors.New("error"),
		}
	}

	if f.WithProposals {
		return []types.Proposal{
			{
				ID: "1",
			},
		}, 123, nil
	}

	return []types.Proposal{}, 123, nil
}

func (f *TestFetcher) GetVote(
	proposal, voter string,
	prevHeight int64,
) (*types.Vote, int64, *types.QueryError) {
	return nil, 456, nil
}

func (f *TestFetcher) GetTallies() (types.ChainTallyInfos, error) {
	return types.ChainTallyInfos{}, nil
}

func (f *TestFetcher) GetChainParams() (*types.ChainWithVotingParams, []error) {
	return nil, []error{}
}
