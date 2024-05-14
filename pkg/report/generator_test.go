package report

import (
	"context"
	"errors"
	"main/pkg/events"
	"main/pkg/fs"
	"main/pkg/tracing"
	"testing"

	"main/pkg/logger"
	"main/pkg/state"
	"main/pkg/types"

	"github.com/stretchr/testify/assert"
)

func TestReportGeneratorWithProposalError(t *testing.T) {
	t.Parallel()

	stateManager := state.NewStateManager("./state.json", &fs.TestFS{}, logger.GetNopLogger())

	oldState := state.NewState()
	newState := state.State{
		ChainInfos: map[string]*state.ChainInfo{
			"chain": {
				ProposalsError: &types.QueryError{
					QueryError: errors.New("test error"),
				},
			},
		},
	}

	generator := NewReportGenerator(stateManager, logger.GetNopLogger(), types.Chains{
		&types.Chain{Name: "chain"},
	}, tracing.InitNoopTracer())

	report := generator.GenerateReport(oldState, newState, context.Background())
	assert.Len(t, report.Entries, 1, "Expected to have 1 entry!")

	entry, ok := report.Entries[0].(events.ProposalsQueryErrorEvent)
	assert.True(t, ok, "Expected to have a proposal query error!")
	assert.Equal(t, "test error", entry.Error.QueryError.Error(), "Error text mismatch!")
}

func TestReportGeneratorWithVoteError(t *testing.T) {
	t.Parallel()

	stateManager := state.NewStateManager("./state.json", &fs.TestFS{}, logger.GetNopLogger())

	oldState := state.NewState()
	newState := state.State{
		ChainInfos: map[string]*state.ChainInfo{
			"chain": {
				ProposalVotes: map[string]state.WalletVotes{
					"proposal": {
						Proposal: types.Proposal{
							ID:     "proposal",
							Status: types.ProposalStatusVoting,
						},
						Votes: map[string]state.ProposalVote{
							"wallet": {
								Error: &types.QueryError{
									QueryError: errors.New("test error"),
								},
							},
						},
					},
				},
			},
		},
	}

	generator := NewReportGenerator(stateManager, logger.GetNopLogger(), types.Chains{
		&types.Chain{Name: "chain"},
	}, tracing.InitNoopTracer())

	report := generator.GenerateReport(oldState, newState, context.Background())
	assert.Len(t, report.Entries, 1, "Expected to have 1 entry!")

	entry, ok := report.Entries[0].(events.VoteQueryError)
	assert.True(t, ok, "Expected to have a vote query error!")
	assert.Equal(t, "proposal", entry.Proposal.ID, "Proposal ID mismatch!")
}

func TestReportGeneratorWithProposalNotInVotingPeriod(t *testing.T) {
	t.Parallel()

	stateManager := state.NewStateManager("./state.json", &fs.TestFS{}, logger.GetNopLogger())

	oldState := state.State{
		ChainInfos: map[string]*state.ChainInfo{
			"chain": {
				ProposalVotes: map[string]state.WalletVotes{
					"proposal": {
						Proposal: types.Proposal{
							ID:     "proposal",
							Status: types.ProposalStatusPassed,
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
							ID:     "proposal",
							Status: types.ProposalStatusPassed,
						},
						Votes: map[string]state.ProposalVote{
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

	generator := NewReportGenerator(stateManager, logger.GetNopLogger(), types.Chains{
		&types.Chain{Name: "chain"},
	}, tracing.InitNoopTracer())

	report := generator.GenerateReport(oldState, newState, context.Background())
	assert.Empty(t, report.Entries)
}

func TestReportGeneratorWithNotVoted(t *testing.T) {
	t.Parallel()

	stateManager := state.NewStateManager("./state.json", &fs.TestFS{}, logger.GetNopLogger())

	oldState := state.NewState()
	newState := state.State{
		ChainInfos: map[string]*state.ChainInfo{
			"chain": {
				ProposalVotes: map[string]state.WalletVotes{
					"proposal": {
						Proposal: types.Proposal{
							ID:     "proposal",
							Status: types.ProposalStatusVoting,
						},
						Votes: map[string]state.ProposalVote{
							"wallet": {},
						},
					},
				},
			},
		},
	}

	generator := NewReportGenerator(stateManager, logger.GetNopLogger(), types.Chains{
		&types.Chain{Name: "chain"},
	}, tracing.InitNoopTracer())

	report := generator.GenerateReport(oldState, newState, context.Background())
	assert.Len(t, report.Entries, 1, "Expected to have 1 entry!")

	entry, ok := report.Entries[0].(events.NotVotedEvent)
	assert.True(t, ok, "Expected to have not voted type!")
	assert.Equal(t, "proposal", entry.Proposal.ID, "Proposal ID mismatch!")
}

func TestReportGeneratorWithVoted(t *testing.T) {
	t.Parallel()

	stateManager := state.NewStateManager("./state.json", &fs.TestFS{}, logger.GetNopLogger())

	oldState := state.State{
		ChainInfos: map[string]*state.ChainInfo{
			"chain": {
				ProposalVotes: map[string]state.WalletVotes{
					"proposal": {
						Proposal: types.Proposal{
							ID:     "proposal",
							Status: types.ProposalStatusVoting,
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
							ID:     "proposal",
							Status: types.ProposalStatusVoting,
						},
						Votes: map[string]state.ProposalVote{
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

	generator := NewReportGenerator(stateManager, logger.GetNopLogger(), types.Chains{
		&types.Chain{Name: "chain"},
	}, tracing.InitNoopTracer())

	report := generator.GenerateReport(oldState, newState, context.Background())
	assert.Len(t, report.Entries, 1, "Expected to have 1 entry!")

	entry, ok := report.Entries[0].(events.VotedEvent)
	assert.True(t, ok, "Expected to have voted type!")
	assert.Equal(t, "proposal", entry.Proposal.ID, "Proposal ID mismatch!")
}

func TestReportGeneratorWithRevoted(t *testing.T) {
	t.Parallel()

	stateManager := state.NewStateManager("./state.json", &fs.TestFS{}, logger.GetNopLogger())

	oldState := state.State{
		ChainInfos: map[string]*state.ChainInfo{
			"chain": {
				ProposalVotes: map[string]state.WalletVotes{
					"proposal": {
						Proposal: types.Proposal{
							ID:     "proposal",
							Status: types.ProposalStatusVoting,
						},
						Votes: map[string]state.ProposalVote{
							"wallet": {
								Vote: &types.Vote{
									Options: []types.VoteOption{{Option: "NO", Weight: 1}},
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
							ID:     "proposal",
							Status: types.ProposalStatusVoting,
						},
						Votes: map[string]state.ProposalVote{
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

	generator := NewReportGenerator(stateManager, logger.GetNopLogger(), types.Chains{
		&types.Chain{Name: "chain"},
	}, tracing.InitNoopTracer())

	report := generator.GenerateReport(oldState, newState, context.Background())
	assert.Len(t, report.Entries, 1, "Expected to have 1 entry!")

	entry, ok := report.Entries[0].(events.RevotedEvent)
	assert.True(t, ok, "Expected to have revoted type!")
	assert.Equal(t, "proposal", entry.Proposal.ID, "Proposal ID mismatch!")
}

func TestReportGeneratorWithFinishedVoting(t *testing.T) {
	t.Parallel()

	stateManager := state.NewStateManager("./state.json", &fs.TestFS{}, logger.GetNopLogger())

	oldState := state.State{
		ChainInfos: map[string]*state.ChainInfo{
			"chain": {
				ProposalVotes: map[string]state.WalletVotes{
					"proposal": {
						Proposal: types.Proposal{
							ID:     "proposal",
							Status: types.ProposalStatusVoting,
						},
						Votes: map[string]state.ProposalVote{
							"wallet": {
								Vote: &types.Vote{Options: types.VoteOptions{{Option: "Yes"}}},
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
							ID:     "proposal",
							Status: types.ProposalStatusPassed,
						},
						Votes: map[string]state.ProposalVote{},
					},
				},
			},
		},
	}

	generator := NewReportGenerator(stateManager, logger.GetNopLogger(), types.Chains{
		&types.Chain{Name: "chain"},
	}, tracing.InitNoopTracer())

	report := generator.GenerateReport(oldState, newState, context.Background())
	assert.Len(t, report.Entries, 1, "Expected to have 1 entry!")

	entry, ok := report.Entries[0].(events.FinishedVotingEvent)
	assert.True(t, ok, "Expected to have not voted type!")
	assert.Equal(t, "proposal", entry.Proposal.ID, "Proposal ID mismatch!")
}
