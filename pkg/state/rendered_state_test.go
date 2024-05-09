package state

import (
	"main/pkg/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestToRenderedStateFilteredChain(t *testing.T) {
	t.Parallel()

	state := State{
		ChainInfos: map[string]*ChainInfo{
			"chain": {
				ProposalVotes: map[string]WalletVotes{},
			},
		},
	}

	renderedState := state.ToRenderedState()
	assert.Empty(t, renderedState.ChainInfos)
}

func TestToRenderedStateFilterProposalsNotInVoting(t *testing.T) {
	t.Parallel()

	state := State{
		ChainInfos: map[string]*ChainInfo{
			"chain": {
				ProposalVotes: map[string]WalletVotes{
					"proposal": {
						Proposal: types.Proposal{
							ID:     "proposal",
							Status: types.ProposalStatusPassed,
						},
					},
				},
			},
		},
	}

	renderedState := state.ToRenderedState()
	assert.Empty(t, renderedState.ChainInfos)
}

func TestToRenderedStateSortProposals(t *testing.T) {
	t.Parallel()

	chain := &types.Chain{Name: "chain"}
	state := NewState()
	state.SetProposal(chain, types.Proposal{ID: "15", Status: types.ProposalStatusVoting})
	state.SetProposal(chain, types.Proposal{ID: "231", Status: types.ProposalStatusVoting})
	state.SetProposal(chain, types.Proposal{ID: "2", Status: types.ProposalStatusVoting})

	renderedState := state.ToRenderedState()
	assert.Len(t, renderedState.ChainInfos, 1)

	renderedChain := renderedState.ChainInfos[0]
	assert.Len(t, renderedChain.ProposalVotes, 3)
	assert.Equal(t, "231", renderedChain.ProposalVotes[0].Proposal.ID)
	assert.Equal(t, "15", renderedChain.ProposalVotes[1].Proposal.ID)
	assert.Equal(t, "2", renderedChain.ProposalVotes[2].Proposal.ID)
}

func TestToRenderedStateSortProposalsInvalid(t *testing.T) {
	t.Parallel()

	chain := &types.Chain{Name: "chain"}
	state := NewState()
	state.SetProposal(chain, types.Proposal{ID: "1", Status: types.ProposalStatusVoting})
	state.SetProposal(chain, types.Proposal{ID: "a", Status: types.ProposalStatusVoting})
	state.SetProposal(chain, types.Proposal{ID: "b", Status: types.ProposalStatusVoting})

	renderedState := state.ToRenderedState()
	assert.Len(t, renderedState.ChainInfos, 1)

	renderedChain := renderedState.ChainInfos[0]
	assert.Len(t, renderedChain.ProposalVotes, 3)

	// we don't know the sorting
}

func TestToRenderedStateSortChains(t *testing.T) {
	t.Parallel()

	state := State{
		ChainInfos: map[string]*ChainInfo{
			"cosmos": {
				Chain: &types.Chain{Name: "cosmos"},
				ProposalVotes: map[string]WalletVotes{
					"proposal": {Proposal: types.Proposal{ID: "1", Status: types.ProposalStatusVoting}},
				},
			},
			"sentinel": {
				Chain: &types.Chain{Name: "sentinel"},
				ProposalVotes: map[string]WalletVotes{
					"proposal": {Proposal: types.Proposal{ID: "1", Status: types.ProposalStatusVoting}},
				},
			},
			"bitsong": {
				Chain: &types.Chain{Name: "bitsong"},
				ProposalVotes: map[string]WalletVotes{
					"proposal": {Proposal: types.Proposal{ID: "1", Status: types.ProposalStatusVoting}},
				},
			},
		},
	}

	renderedState := state.ToRenderedState()
	assert.Len(t, renderedState.ChainInfos, 3)
	assert.Equal(t, "bitsong", renderedState.ChainInfos[0].Chain.Name)
	assert.Equal(t, "cosmos", renderedState.ChainInfos[1].Chain.Name)
	assert.Equal(t, "sentinel", renderedState.ChainInfos[2].Chain.Name)
}

func TestToRenderedStateSortWallets(t *testing.T) {
	t.Parallel()

	chain := &types.Chain{Name: "chain"}
	proposal := types.Proposal{ID: "proposal", Status: types.ProposalStatusVoting}
	state := NewState()

	wallets := []string{"21", "352", "2"}

	for _, addr := range wallets {
		wallet := &types.Wallet{Address: addr}
		state.SetVote(chain, proposal, wallet, ProposalVote{Wallet: wallet})
	}

	renderedState := state.ToRenderedState()
	assert.Len(t, renderedState.ChainInfos, 1)

	renderedChain := renderedState.ChainInfos[0]
	assert.Len(t, renderedChain.ProposalVotes, 1)

	renderedVotes := renderedChain.ProposalVotes[0]
	assert.Len(t, renderedVotes.Votes, 3)

	assert.Equal(t, "2", renderedVotes.Votes[0].Wallet.Address)
	assert.Equal(t, "21", renderedVotes.Votes[1].Wallet.Address)
	assert.Equal(t, "352", renderedVotes.Votes[2].Wallet.Address)
}

func TestRenderedWalletVoteHasVoted(t *testing.T) {
	t.Parallel()

	assert.True(t, RenderedWalletVote{Vote: &types.Vote{}}.HasVoted())
	assert.False(t, RenderedWalletVote{Vote: &types.Vote{}, Error: &types.QueryError{}}.HasVoted())
	assert.False(t, RenderedWalletVote{Error: &types.QueryError{}}.HasVoted())
}

func TestRenderedWalletVoteIsError(t *testing.T) {
	t.Parallel()

	assert.False(t, RenderedWalletVote{}.IsError())
	assert.True(t, RenderedWalletVote{Error: &types.QueryError{}}.IsError())
}

func TestRenderedChainInfoHasError(t *testing.T) {
	t.Parallel()

	assert.False(t, RenderedChainInfo{}.HasProposalsError())
	assert.True(t, RenderedChainInfo{ProposalsError: &types.QueryError{}}.HasProposalsError())
}
