package state

import (
	"context"
	"main/pkg/fetchers"
	"main/pkg/logger"
	"main/pkg/tracing"
	"main/pkg/types"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestReportGeneratorNew(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	chain := &types.Chain{Name: "chain", Type: "cosmos"}
	chains := types.Chains{chain}

	generator := NewStateGenerator(log, tracing.InitNoopTracer(), chains)
	assert.NotNil(t, generator)
}

func TestReportGeneratorProcessChain(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	chain := &types.Chain{Name: "chain", Type: "cosmos"}
	chains := types.Chains{chain}

	generator := Generator{
		Logger: *log,
		Chains: chains,
		Fetchers: map[string]fetchers.Fetcher{
			"chain": &fetchers.TestFetcher{WithPassedProposals: true},
		},
		Tracer: tracing.InitNoopTracer(),
	}

	oldState := NewState()
	newState := generator.GetState(oldState, context.Background())
	assert.Len(t, newState.ChainInfos, 1)
}

func TestReportGeneratorProcessProposalsWithPassed(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	chain := &types.Chain{Name: "chain", Type: "cosmos"}
	chains := types.Chains{chain}

	generator := Generator{
		Logger: *log,
		Chains: chains,
		Fetchers: map[string]fetchers.Fetcher{
			"chain": &fetchers.TestFetcher{WithPassedProposals: true},
		},
		Tracer: tracing.InitNoopTracer(),
	}

	oldState := NewState()
	newState := generator.GetState(oldState, context.Background())
	assert.Len(t, newState.ChainInfos, 1)
}

func TestReportGeneratorProcessProposalWithError(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	chain := &types.Chain{Name: "chain", Type: "cosmos"}
	chains := types.Chains{chain}
	fetcher := &fetchers.TestFetcher{WithProposalsError: true}

	generator := Generator{
		Logger: *log,
		Chains: chains,
		Fetchers: map[string]fetchers.Fetcher{
			"chain": fetcher,
		},
		Tracer: tracing.InitNoopTracer(),
	}

	oldVotes := map[string]WalletVotes{
		"1": {
			Proposal: types.Proposal{ID: "1", Status: types.ProposalStatusVoting},
		},
	}

	oldState := NewState()
	oldState.SetChainProposalsHeight(chain, 15)
	oldState.SetChainVotes(chain, oldVotes)

	newState := NewState()
	generator.ProcessChain(chain, newState, oldState, fetcher, context.Background())
	assert.Len(t, newState.ChainInfos, 1)

	newVotes, ok := newState.ChainInfos["chain"]
	assert.True(t, ok)
	assert.NotNil(t, newVotes)
	require.Error(t, newVotes.ProposalsError)
	assert.Equal(t, int64(15), newVotes.ProposalsHeight)

	proposal, ok := newVotes.ProposalVotes["1"]
	assert.True(t, ok)
	assert.Equal(t, "1", proposal.Proposal.ID)
}

func TestReportGeneratorProcessProposalWithoutError(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	chain := &types.Chain{
		Name:    "chain",
		Type:    "cosmos",
		Wallets: []*types.Wallet{{Address: "me"}},
	}
	chains := types.Chains{chain}
	fetcher := &fetchers.TestFetcher{WithVote: true}

	generator := Generator{
		Logger: *log,
		Chains: chains,
		Fetchers: map[string]fetchers.Fetcher{
			"chain": fetcher,
		},
		Tracer: tracing.InitNoopTracer(),
	}

	oldVotes := map[string]WalletVotes{
		"1": {
			Proposal: types.Proposal{ID: "1", Status: types.ProposalStatusVoting},
		},
	}

	oldState := NewState()
	oldState.SetChainProposalsHeight(chain, 15)
	oldState.SetChainVotes(chain, oldVotes)

	newState := NewState()
	generator.ProcessChain(chain, newState, oldState, fetcher, context.Background())
	assert.Len(t, newState.ChainInfos, 1)

	newVotes, ok := newState.ChainInfos["chain"]
	assert.True(t, ok)
	assert.NotNil(t, newVotes)
	require.Error(t, newVotes.ProposalsError)
	assert.Equal(t, int64(123), newVotes.ProposalsHeight)

	proposal, ok := newVotes.ProposalVotes["1"]
	assert.True(t, ok)
	assert.Equal(t, "1", proposal.Proposal.ID)
}

func TestReportGeneratorProcessVoteWithError(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	wallet := &types.Wallet{Address: "me"}
	chain := &types.Chain{
		Name:    "chain",
		Type:    "cosmos",
		Wallets: []*types.Wallet{wallet},
	}
	chains := types.Chains{chain}
	fetcher := &fetchers.TestFetcher{WithVoteError: true}

	proposal := types.Proposal{ID: "1", Status: types.ProposalStatusVoting}
	generator := Generator{
		Logger: *log,
		Chains: chains,
		Fetchers: map[string]fetchers.Fetcher{
			"chain": fetcher,
		},
		Tracer: tracing.InitNoopTracer(),
	}

	oldVotes := map[string]WalletVotes{
		"1": {
			Proposal: proposal,
			Votes: map[string]ProposalVote{
				"me": {
					Vote:   &types.Vote{Voter: "not_me"},
					Height: 15,
				},
			},
		},
	}

	oldState := NewState()
	oldState.SetChainProposalsHeight(chain, 15)
	oldState.SetChainVotes(chain, oldVotes)

	newState := NewState()
	generator.ProcessProposalAndWallet(chain, proposal, fetcher, wallet, newState, oldState, context.Background())
	assert.Len(t, newState.ChainInfos, 1)

	newVotes, ok := newState.ChainInfos["chain"]
	assert.True(t, ok)
	assert.NotNil(t, newVotes)

	newProposal, ok := newVotes.ProposalVotes["1"]
	assert.True(t, ok)
	assert.Equal(t, "1", newProposal.Proposal.ID)

	newVote, ok := newProposal.Votes["me"]
	assert.True(t, ok)
	assert.Equal(t, int64(15), newVote.Height)
	require.Error(t, newVote.Error)
	assert.Equal(t, "not_me", newVote.Vote.Voter)
}

func TestReportGeneratorProcessVoteWithDisappearedVote(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	wallet := &types.Wallet{Address: "me"}
	chain := &types.Chain{
		Name:    "chain",
		Type:    "cosmos",
		Wallets: []*types.Wallet{wallet},
	}
	chains := types.Chains{chain}
	fetcher := &fetchers.TestFetcher{}

	proposal := types.Proposal{ID: "1", Status: types.ProposalStatusVoting}
	generator := Generator{
		Logger: *log,
		Chains: chains,
		Fetchers: map[string]fetchers.Fetcher{
			"chain": fetcher,
		},
		Tracer: tracing.InitNoopTracer(),
	}

	oldVotes := map[string]WalletVotes{
		"1": {
			Proposal: proposal,
			Votes: map[string]ProposalVote{
				"me": {
					Vote:   &types.Vote{Voter: "not_me"},
					Height: 15,
				},
			},
		},
	}

	oldState := NewState()
	oldState.SetChainProposalsHeight(chain, 15)
	oldState.SetChainVotes(chain, oldVotes)

	newState := NewState()
	generator.ProcessProposalAndWallet(chain, proposal, fetcher, wallet, newState, oldState, context.Background())
	assert.Len(t, newState.ChainInfos, 1)

	newVotes, ok := newState.ChainInfos["chain"]
	assert.True(t, ok)
	assert.NotNil(t, newVotes)

	newProposal, ok := newVotes.ProposalVotes["1"]
	assert.True(t, ok)
	assert.Equal(t, "1", newProposal.Proposal.ID)

	newVote, ok := newProposal.Votes["me"]
	assert.True(t, ok)
	assert.Equal(t, int64(456), newVote.Height)
	require.Error(t, newVote.Error)
	assert.Equal(t, "not_me", newVote.Vote.Voter)
}

func TestReportGeneratorProcessVoteWithOkVote(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	wallet := &types.Wallet{Address: "me"}
	chain := &types.Chain{
		Name:    "chain",
		Type:    "cosmos",
		Wallets: []*types.Wallet{wallet},
	}
	chains := types.Chains{chain}
	fetcher := &fetchers.TestFetcher{WithVote: true}

	proposal := types.Proposal{ID: "1", Status: types.ProposalStatusVoting}
	generator := Generator{
		Logger: *log,
		Chains: chains,
		Fetchers: map[string]fetchers.Fetcher{
			"chain": fetcher,
		},
		Tracer: tracing.InitNoopTracer(),
	}

	oldVotes := map[string]WalletVotes{
		"1": {
			Proposal: proposal,
			Votes: map[string]ProposalVote{
				"me": {
					Vote:   &types.Vote{Voter: "not_me"},
					Height: 15,
				},
			},
		},
	}

	oldState := NewState()
	oldState.SetChainProposalsHeight(chain, 15)
	oldState.SetChainVotes(chain, oldVotes)

	newState := NewState()
	generator.ProcessProposalAndWallet(chain, proposal, fetcher, wallet, newState, oldState, context.Background())
	assert.Len(t, newState.ChainInfos, 1)

	newVotes, ok := newState.ChainInfos["chain"]
	assert.True(t, ok)
	assert.NotNil(t, newVotes)

	newProposal, ok := newVotes.ProposalVotes["1"]
	assert.True(t, ok)
	assert.Equal(t, "1", newProposal.Proposal.ID)

	newVote, ok := newProposal.Votes["me"]
	assert.True(t, ok)
	assert.Equal(t, int64(456), newVote.Height)
	require.Error(t, newVote.Error)
	assert.Equal(t, "me", newVote.Vote.Voter)
}
