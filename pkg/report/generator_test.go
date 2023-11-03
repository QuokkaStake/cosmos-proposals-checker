package report

import (
	"main/pkg/events"
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
				ProposalsError: types.NewJSONError("test error"),
			},
		},
	}

	generator := NewReportGenerator(stateManager, logger.GetDefaultLogger(), configTypes.Chains{
		&configTypes.Chain{Name: "chain"},
	})

	report := generator.GenerateReport(oldState, newState)
	assert.Len(t, report.Entries, 1, "Expected to have 1 entry!")

	entry, ok := report.Entries[0].(events.ProposalsQueryErrorEvent)
	assert.True(t, ok, "Expected to have a proposal query error!")
	assert.Equal(t, "test error", entry.Error.Error(), "Error text mismatch!")
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
							ID: "proposal",
						},
						Votes: map[string]state.ProposalVote{
							"wallet": {
								Error: types.NewJSONError("test error"),
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
	assert.Len(t, report.Entries, 1, "Expected to have 1 entry!")

	entry, ok := report.Entries[0].(events.VoteQueryError)
	assert.True(t, ok, "Expected to have a vote query error!")
	assert.Equal(t, "proposal", entry.Proposal.ID, "Proposal ID mismatch!")
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
							ID: "proposal",
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
	assert.Len(t, report.Entries, 1, "Expected to have 1 entry!")

	entry, ok := report.Entries[0].(events.NotVotedEvent)
	assert.True(t, ok, "Expected to have not voted type!")
	assert.Equal(t, "proposal", entry.Proposal.ID, "Proposal ID mismatch!")
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
							ID: "proposal",
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
							ID: "proposal",
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
	assert.Len(t, report.Entries, 1, "Expected to have 1 entry!")

	entry, ok := report.Entries[0].(events.VotedEvent)
	assert.True(t, ok, "Expected to have voted type!")
	assert.Equal(t, "proposal", entry.Proposal.ID, "Proposal ID mismatch!")
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
							ID: "proposal",
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
							ID: "proposal",
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
	assert.Len(t, report.Entries, 1, "Expected to have 1 entry!")

	entry, ok := report.Entries[0].(events.RevotedEvent)
	assert.True(t, ok, "Expected to have revoted type!")
	assert.Equal(t, "proposal", entry.Proposal.ID, "Proposal ID mismatch!")
}
