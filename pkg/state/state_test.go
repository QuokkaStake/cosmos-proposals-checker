package state

import (
	"errors"
	"testing"

	configTypes "main/pkg/config/types"
	"main/pkg/types"

	"github.com/stretchr/testify/assert"
)

func TestSetVoteWithoutChainInfo(t *testing.T) {
	t.Parallel()

	state := NewState()
	_, _, found := state.GetVoteAndProposal("chain", "proposal", "wallet")
	assert.Equal(t, found, false, "Vote should not be presented!")

	state.SetVote(
		&configTypes.Chain{Name: "chain"},
		types.Proposal{ProposalID: "proposal"},
		&configTypes.Wallet{Address: "wallet"},
		ProposalVote{
			Vote: &types.Vote{
				Option: "yep",
			},
		},
	)

	vote, _, found2 := state.GetVoteAndProposal("chain", "proposal", "wallet")
	assert.Equal(t, found2, true, "Vote should be presented!")
	assert.Equal(t, vote.HasVoted(), true, "Vote should be presented!")
	assert.Equal(t, vote.IsError(), false, "There should be no error!")
}

func TestSetProposalErrorWithoutChainInfo(t *testing.T) {
	t.Parallel()

	state := NewState()
	state.SetChainProposalsError(&configTypes.Chain{Name: "test"}, errors.New("test error"))

	hasError2 := state.ChainInfos["test"].HasProposalsError()
	assert.Equal(t, hasError2, true, "Chain info should have a proposal error!")

	err := state.ChainInfos["test"].ProposalsError
	assert.Equal(t, err.Error(), "test error", "Errors text should match!")
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

	state.SetChainVotes(&configTypes.Chain{Name: "chain"}, votes)

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
	assert.Equal(t, hasError, false, "Chain info should not have a proposal error!")

	state.SetChainProposalsError(&configTypes.Chain{Name: "test"}, errors.New("test error"))

	hasError2 := state.ChainInfos["test"].HasProposalsError()
	assert.Equal(t, hasError2, true, "Chain info should have a proposal error!")

	err := state.ChainInfos["test"].ProposalsError
	assert.Equal(t, err.Error(), "test error", "Errors text should match!")
}

func TestGetVoteWithoutChainInfo(t *testing.T) {
	t.Parallel()

	state := State{}

	_, _, found := state.GetVoteAndProposal("chain", "proposal", "wallet")
	assert.Equal(t, found, false, "There should be no vote!")
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
	assert.Equal(t, found, false, "There should be no vote!")
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
	assert.Equal(t, found, false, "There should be no vote!")
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
	assert.Equal(t, found, true, "There should be a vote!")
}

func TestHasVotedWithoutChainInfo(t *testing.T) {
	t.Parallel()

	state := State{}

	voted := state.HasVoted("chain", "proposal", "wallet")
	assert.Equal(t, voted, false, "There should be no vote!")
}

func TestHasVotedWithChainInfo(t *testing.T) {
	t.Parallel()

	state := State{
		ChainInfos: map[string]*ChainInfo{
			"chain": {},
		},
	}

	voted := state.HasVoted("chain", "proposal", "wallet")
	assert.Equal(t, voted, false, "There should be no vote!")
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
	assert.Equal(t, voted, false, "There should be no vote!")
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
	assert.Equal(t, voted, false, "There should be no vote!")
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
	assert.Equal(t, voted, true, "There should be a vote!")
}
