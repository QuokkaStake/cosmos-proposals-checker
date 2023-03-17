package report

import (
	"testing"

	configTypes "main/pkg/config/types"
	"main/pkg/logger"
	"main/pkg/state"
	"main/pkg/types"

	"github.com/stretchr/testify/assert"
)

func TestReportGeneratorWithProposalError(t *testing.T) {
	t.Parallel()

	stateManager := state.NewStateManager("./state.json", logger.GetDefaultLogger())

	oldState := state.NewState()
	newState := state.State{
		ChainInfos: map[string]*state.ChainInfo{
			"chain": {
				ProposalsError: "test error",
			},
		},
	}

	generator := NewReportGenerator(stateManager, logger.GetDefaultLogger(), configTypes.Chains{
		&configTypes.Chain{Name: "chain"},
	})

	report := generator.GenerateReport(oldState, newState)
	assert.Equal(t, len(report.Entries), 1, "Expected to have 1 entry!")
	assert.Equal(t, report.Entries[0].Type, types.ProposalQueryError, "Expected to have a proposal query error!")
	assert.Equal(t, report.Entries[0].Value, "test error", "Error text mismatch!")
}

func TestReportGeneratorWithVoteError(t *testing.T) {
	t.Parallel()

	stateManager := state.NewStateManager("./state.json", logger.GetDefaultLogger())

	oldState := state.NewState()
	newState := state.State{
		ChainInfos: map[string]*state.ChainInfo{
			"chain": {
				ProposalVotes: map[string]state.WalletVotes{
					"proposal": {
						Proposal: types.Proposal{
							Content: &types.ProposalContent{},
						},
						Votes: map[string]state.ProposalVote{
							"wallet": {
								Error: "test error",
							},
						},
					},
				},
			},
		},
	}

	generator := NewReportGenerator(stateManager, logger.GetDefaultLogger(), configTypes.Chains{
		&configTypes.Chain{Name: "chain"},
	})

	report := generator.GenerateReport(oldState, newState)
	assert.Equal(t, len(report.Entries), 1, "Expected to have 1 entry!")
	assert.Equal(t, report.Entries[0].Type, types.VoteQueryError, "Expected to have a vote query error!")
	assert.Equal(t, report.Entries[0].ProposalID, "proposal", "Proposal ID mismatch!")
}

func TestReportGeneratorWithNotVoted(t *testing.T) {
	t.Parallel()

	stateManager := state.NewStateManager("./state.json", logger.GetDefaultLogger())

	oldState := state.NewState()
	newState := state.State{
		ChainInfos: map[string]*state.ChainInfo{
			"chain": {
				ProposalVotes: map[string]state.WalletVotes{
					"proposal": {
						Proposal: types.Proposal{
							Content: &types.ProposalContent{},
						},
						Votes: map[string]state.ProposalVote{
							"wallet": {},
						},
					},
				},
			},
		},
	}

	generator := NewReportGenerator(stateManager, logger.GetDefaultLogger(), configTypes.Chains{
		&configTypes.Chain{Name: "chain"},
	})

	report := generator.GenerateReport(oldState, newState)
	assert.Equal(t, len(report.Entries), 1, "Expected to have 1 entry!")
	assert.Equal(t, report.Entries[0].Type, types.NotVoted, "Expected to have not voted type!")
	assert.Equal(t, report.Entries[0].ProposalID, "proposal", "Proposal ID mismatch!")
}

func TestReportGeneratorWithVoted(t *testing.T) {
	t.Parallel()

	stateManager := state.NewStateManager("./state.json", logger.GetDefaultLogger())

	oldState := state.State{
		ChainInfos: map[string]*state.ChainInfo{
			"chain": {
				ProposalVotes: map[string]state.WalletVotes{
					"proposal": {
						Proposal: types.Proposal{
							Content: &types.ProposalContent{},
						},
						Votes: map[string]state.ProposalVote{
							"wallet": {},
						},
					},
				},
			},
		},
	}
	newState := state.State{
		ChainInfos: map[string]*state.ChainInfo{
			"chain": {
				ProposalVotes: map[string]state.WalletVotes{
					"proposal": {
						Proposal: types.Proposal{
							Content: &types.ProposalContent{},
						},
						Votes: map[string]state.ProposalVote{
							"wallet": {
								Vote: &types.Vote{
									Option: "YES",
								},
							},
						},
					},
				},
			},
		},
	}

	generator := NewReportGenerator(stateManager, logger.GetDefaultLogger(), configTypes.Chains{
		&configTypes.Chain{Name: "chain"},
	})

	report := generator.GenerateReport(oldState, newState)
	assert.Equal(t, len(report.Entries), 1, "Expected to have 1 entry!")
	assert.Equal(t, report.Entries[0].Type, types.Voted, "Expected to have voted type!")
	assert.Equal(t, report.Entries[0].ProposalID, "proposal", "Proposal ID mismatch!")
}

func TestReportGeneratorWithRevoted(t *testing.T) {
	t.Parallel()

	stateManager := state.NewStateManager("./state.json", logger.GetDefaultLogger())

	oldState := state.State{
		ChainInfos: map[string]*state.ChainInfo{
			"chain": {
				ProposalVotes: map[string]state.WalletVotes{
					"proposal": {
						Proposal: types.Proposal{
							Content: &types.ProposalContent{},
						},
						Votes: map[string]state.ProposalVote{
							"wallet": {
								Vote: &types.Vote{
									Option: "NO",
								},
							},
						},
					},
				},
			},
		},
	}
	newState := state.State{
		ChainInfos: map[string]*state.ChainInfo{
			"chain": {
				ProposalVotes: map[string]state.WalletVotes{
					"proposal": {
						Proposal: types.Proposal{
							Content: &types.ProposalContent{},
						},
						Votes: map[string]state.ProposalVote{
							"wallet": {
								Vote: &types.Vote{
									Option: "YES",
								},
							},
						},
					},
				},
			},
		},
	}

	generator := NewReportGenerator(stateManager, logger.GetDefaultLogger(), configTypes.Chains{
		&configTypes.Chain{Name: "chain"},
	})

	report := generator.GenerateReport(oldState, newState)
	assert.Equal(t, len(report.Entries), 1, "Expected to have 1 entry!")
	assert.Equal(t, report.Entries[0].Type, types.Revoted, "Expected to have revoted type!")
	assert.Equal(t, report.Entries[0].ProposalID, "proposal", "Proposal ID mismatch!")
}
