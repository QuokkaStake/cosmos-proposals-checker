package state

import (
	"errors"
	"testing"

	"main/pkg/types"

	"github.com/stretchr/testify/assert"
)

func TestSetVoteWithoutChainInfo(t *testing.T) {
	t.Parallel()

	state := NewState()
	_, _, found := state.GetVoteAndProposal("chain", "proposal", "wallet")
	assert.False(t, found, "Vote should not be presented!")

	state.SetVote(
		&types.Chain{Name: "chain"},
		types.Proposal{ID: "proposal"},
		&types.Wallet{Address: "wallet"},
		ProposalVote{
			Vote: &types.Vote{
				Option: "yep",
			},
		},
	)

	vote, _, found2 := state.GetVoteAndProposal("chain", "proposal", "wallet")
	assert.True(t, found2, "Vote should be presented!")
	assert.True(t, vote.HasVoted(), "Vote should be presented!")
	assert.False(t, vote.IsError(), "There should be no error!")
}

func TestSetProposalErrorWithoutChainInfo(t *testing.T) {
	t.Parallel()

	state := NewState()
	state.SetChainProposalsError(&types.Chain{Name: "test"}, &types.QueryError{
		QueryError: errors.New("test error"),
	})

	hasError2 := state.ChainInfos["test"].HasProposalsError()
	assert.True(t, hasError2, "Chain info should have a proposal error!")

	err := state.ChainInfos["test"].ProposalsError
	assert.Equal(t, "test error", err.QueryError.Error(), "Errors text should match!")
}

func TestSetVotes(t *testing.T) {
	t.Parallel()

	state := State{
		ChainInfos: map[string]*ChainInfo{
			"chain": {},
		},
	}

	votes := map[string]WalletVotes{
		"proposal": {
			Votes: map[string]ProposalVote{
				"wallet": {
					Vote: &types.Vote{Option: "YES"},
				},
			},
		},
	}

	state.SetChainVotes(&types.Chain{Name: "chain"}, votes)

	votes2 := state.ChainInfos["chain"].ProposalVotes
	assert.NotNil(t, votes2, "There should be votes!")
}

func TestSetProposalErrorWithChainInfo(t *testing.T) {
	t.Parallel()

	state := State{
		ChainInfos: map[string]*ChainInfo{
			"test": {},
		},
	}

	hasError := state.ChainInfos["test"].HasProposalsError()
	assert.False(t, hasError, "Chain info should not have a proposal error!")

	state.SetChainProposalsError(&types.Chain{Name: "test"}, &types.QueryError{
		QueryError: errors.New("test error"),
	})

	hasError2 := state.ChainInfos["test"].HasProposalsError()
	assert.True(t, hasError2, "Chain info should have a proposal error!")

	err := state.ChainInfos["test"].ProposalsError
	assert.Equal(t, "test error", err.QueryError.Error(), "Errors text should match!")
}

func TestGetVoteWithoutChainInfo(t *testing.T) {
	t.Parallel()

	state := State{}

	_, _, found := state.GetVoteAndProposal("chain", "proposal", "wallet")
	assert.False(t, found, "There should be no vote!")
}

func TestGetVoteWithoutProposalVotes(t *testing.T) {
	t.Parallel()

	state := State{
		ChainInfos: map[string]*ChainInfo{
			"chain": {
				ProposalVotes: map[string]WalletVotes{},
			},
		},
	}

	_, _, found := state.GetVoteAndProposal("chain", "proposal", "wallet")
	assert.False(t, found, "There should be no vote!")
}

func TestGetVoteWithWalletVoteNotPresent(t *testing.T) {
	t.Parallel()

	state := State{
		ChainInfos: map[string]*ChainInfo{
			"chain": {
				ProposalVotes: map[string]WalletVotes{
					"proposal": {},
				},
			},
		},
	}

	_, _, found := state.GetVoteAndProposal("chain", "proposal", "wallet")
	assert.False(t, found, "There should be no vote!")
}

func TestGetVoteWithWalletVotePresent(t *testing.T) {
	t.Parallel()

	state := State{
		ChainInfos: map[string]*ChainInfo{
			"chain": {
				ProposalVotes: map[string]WalletVotes{
					"proposal": {
						Votes: map[string]ProposalVote{
							"wallet": {},
						},
					},
				},
			},
		},
	}

	_, _, found := state.GetVoteAndProposal("chain", "proposal", "wallet")
	assert.True(t, found, "There should be a vote!")
}

func TestHasVotedWithoutChainInfo(t *testing.T) {
	t.Parallel()

	state := State{}

	voted := state.HasVoted("chain", "proposal", "wallet")
	assert.False(t, voted, "There should be no vote!")
}

func TestHasVotedWithChainInfo(t *testing.T) {
	t.Parallel()

	state := State{
		ChainInfos: map[string]*ChainInfo{
			"chain": {},
		},
	}

	voted := state.HasVoted("chain", "proposal", "wallet")
	assert.False(t, voted, "There should be no vote!")
}

func TestHasVotedWithWalletVoteIntoNotPresent(t *testing.T) {
	t.Parallel()

	state := State{
		ChainInfos: map[string]*ChainInfo{
			"chain": {
				ProposalVotes: map[string]WalletVotes{
					"proposal": {},
				},
			},
		},
	}

	voted := state.HasVoted("chain", "proposal", "wallet")
	assert.False(t, voted, "There should be no vote!")
}

func TestHasVotedWithWalletVoteInfoPresent(t *testing.T) {
	t.Parallel()

	state := State{
		ChainInfos: map[string]*ChainInfo{
			"chain": {
				ProposalVotes: map[string]WalletVotes{
					"proposal": {
						Votes: map[string]ProposalVote{
							"wallet": {},
						},
					},
				},
			},
		},
	}

	voted := state.HasVoted("chain", "proposal", "wallet")
	assert.False(t, voted, "There should be no vote!")
}

func TestHasVotedWithWalletVotePresent(t *testing.T) {
	t.Parallel()

	state := State{
		ChainInfos: map[string]*ChainInfo{
			"chain": {
				ProposalVotes: map[string]WalletVotes{
					"proposal": {
						Votes: map[string]ProposalVote{
							"wallet": {
								Vote: &types.Vote{
									Option: "YEP",
								},
							},
						},
					},
				},
			},
		},
	}

	voted := state.HasVoted("chain", "proposal", "wallet")
	assert.True(t, voted, "There should be a vote!")
}
