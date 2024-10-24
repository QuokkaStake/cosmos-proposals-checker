package fetchers

import (
	"context"
	"errors"
	"main/pkg/types"
)

type TestFetcher struct {
	WithPassedProposals bool
	WithProposalsError  bool
	WithVote            bool
	WithVoteError       bool
	WithTallyError      bool
	WithTallyNotEmpty   bool
	WithParamsError     bool
}

func (f *TestFetcher) GetAllProposals(
	prevHeight int64,
	ctx context.Context,
) ([]types.Proposal, int64, *types.QueryError) {
	if f.WithProposalsError {
		return []types.Proposal{}, 123, &types.QueryError{
			QueryError: errors.New("error"),
		}
	}

	if f.WithPassedProposals {
		return []types.Proposal{
			{
				ID:     "1",
				Status: types.ProposalStatusPassed,
			},
		}, 123, nil
	}

	return []types.Proposal{
		{
			ID:     "1",
			Status: types.ProposalStatusVoting,
		},
	}, 123, nil
}

func (f *TestFetcher) GetVote(
	proposal, voter string,
	prevHeight int64,
	ctx context.Context,
) (*types.Vote, int64, *types.QueryError) {
	if f.WithVoteError {
		return nil, 456, &types.QueryError{
			QueryError: errors.New("vote query error"),
		}
	}

	if f.WithVote {
		return &types.Vote{
			ProposalID: "1",
			Voter:      "me",
			Options:    types.VoteOptions{},
		}, 456, nil
	}

	return nil, 456, nil
}

func (f *TestFetcher) GetTallies(ctx context.Context) (types.ChainTallyInfos, error) {
	if f.WithTallyError {
		return types.ChainTallyInfos{}, &types.QueryError{
			QueryError: errors.New("error"),
		}
	}

	if f.WithTallyNotEmpty {
		return types.ChainTallyInfos{
			Chain: &types.Chain{Name: "test"},
			TallyInfos: []types.TallyInfo{
				{
					Proposal: types.Proposal{ID: "id"},
					Tally:    types.Tally{},
				},
			},
		}, nil
	}

	return types.ChainTallyInfos{}, nil
}

func (f *TestFetcher) GetChainParams(ctx context.Context) (*types.ChainWithVotingParams, []error) {
	if f.WithParamsError {
		return &types.ChainWithVotingParams{}, []error{
			errors.New("test"),
		}
	}

	return &types.ChainWithVotingParams{
		Chain:  &types.Chain{Name: "test"},
		Params: []types.ChainParam{types.BoolParam{Value: true, Description: "param"}},
	}, []error{}
}
