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
	_, found := state.GetVote("chain", "proposal", "wallet")
	assert.False(t, found, "Vote should not be presented!")

	state.SetVote(
		&types.Chain{Name: "chain"},
		types.Proposal{ID: "proposal"},
		&types.Wallet{Address: "wallet"},
		ProposalVote{
			Vote: &types.Vote{
				Options: []types.VoteOption{{Option: "YES", Weight: 1}},
			},
		},
	)

	vote, found2 := state.GetVote("chain", "proposal", "wallet")
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
					Vote: &types.Vote{
						Options: []types.VoteOption{{Option: "YES", Weight: 1}},
					},
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

func TestSetProposalLastHeight(t *testing.T) {
	t.Parallel()

	state := State{
		ChainInfos: map[string]*ChainInfo{
			"chain1": {ProposalsHeight: 123},
		},
	}

	assert.Equal(t, int64(123), state.GetLastProposalsHeight(&types.Chain{Name: "chain1"}))
	assert.Equal(t, int64(0), state.GetLastProposalsHeight(&types.Chain{Name: "chain2"}))

	state.SetChainProposalsHeight(&types.Chain{Name: "chain1"}, 456)
	state.SetChainProposalsHeight(&types.Chain{Name: "chain2"}, 789)

	assert.Equal(t, int64(456), state.GetLastProposalsHeight(&types.Chain{Name: "chain1"}))
	assert.Equal(t, int64(789), state.GetLastProposalsHeight(&types.Chain{Name: "chain2"}))
}

func TestGetVoteWithoutChainInfo(t *testing.T) {
	t.Parallel()

	state := State{}

	_, found := state.GetVote("chain", "proposal", "wallet")
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

	_, found := state.GetVote("chain", "proposal", "wallet")
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

	_, found := state.GetVote("chain", "proposal", "wallet")
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

	_, found := state.GetVote("chain", "proposal", "wallet")
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
									Options: []types.VoteOption{{Option: "YES", Weight: 1}},
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

func TestSetProposal(t *testing.T) {
	t.Parallel()

	state := NewState()
	proposal := types.Proposal{ID: "id"}
	chain := &types.Chain{Name: "chain"}
	assert.Empty(t, state.ChainInfos)

	state.SetProposal(chain, proposal)
	assert.NotEmpty(t, state.ChainInfos)

	chainInfo, ok := state.ChainInfos["chain"]
	assert.True(t, ok)
	assert.NotNil(t, chainInfo)

	proposalInfo, ok := chainInfo.ProposalVotes["id"]
	assert.True(t, ok)
	assert.NotNil(t, proposalInfo)
}

func TestGetProposalWithChainNotPresent(t *testing.T) {
	t.Parallel()

	state := State{
		ChainInfos: map[string]*ChainInfo{},
	}

	_, found := state.GetProposal("chain", "proposal")
	assert.False(t, found)
}

func TestGetProposalWithProposalNotPresent(t *testing.T) {
	t.Parallel()

	state := State{
		ChainInfos: map[string]*ChainInfo{
			"chain": {
				ProposalVotes: map[string]WalletVotes{},
			},
		},
	}

	_, found := state.GetProposal("chain", "proposal")
	assert.False(t, found)
}

func TestGetProposalWithProposalPresent(t *testing.T) {
	t.Parallel()

	state := State{
		ChainInfos: map[string]*ChainInfo{
			"chain": {
				ProposalVotes: map[string]WalletVotes{
					"proposal": {Proposal: types.Proposal{ID: "proposal"}},
				},
			},
		},
	}

	proposal, found := state.GetProposal("chain", "proposal")
	assert.True(t, found)
	assert.Equal(t, "proposal", proposal.ID)
}
