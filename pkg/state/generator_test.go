package state

import (
	"main/pkg/fetchers"
	"main/pkg/logger"
	"main/pkg/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReportGeneratorNew(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	chain := &types.Chain{Name: "chain", Type: "cosmos"}
	chains := types.Chains{chain}

	generator := NewStateGenerator(log, chains)
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
			"chain": &fetchers.TestFetcher{},
		},
	}

	oldState := NewState()
	newState := generator.GetState(oldState)
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
	}

	oldVotes := map[string]WalletVotes{
		"1": {
			Proposal: types.Proposal{ID: "1"},
		},
	}

	oldState := NewState()
	oldState.SetChainProposalsHeight(chain, 15)
	oldState.SetChainVotes(chain, oldVotes)

	newState := NewState()
	generator.ProcessChain(chain, newState, oldState, fetcher)
	assert.Len(t, newState.ChainInfos, 1)

	newVotes, ok := newState.ChainInfos["chain"]
	assert.True(t, ok)
	assert.NotNil(t, newVotes)
	assert.NotNil(t, newVotes.ProposalsError)
	assert.Equal(t, int64(15), newVotes.ProposalsHeight)

	proposal, ok := newVotes.ProposalVotes["1"]
	assert.True(t, ok)
	assert.Equal(t, "1", proposal.Proposal.ID)
}
