package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	stateManager = NewStateManager("./state.json", GetDefaultLogger())
)

func TestReportGeneratorWithProposalError(t *testing.T) {
	t.Parallel()

	oldState := NewState()
	newState := State{
		ChainInfos: map[string]*ChainInfo{
			"chain": {
				ProposalsError: "test error",
			},
		},
	}

	generator := NewReportGenerator(stateManager, GetDefaultLogger(), Chains{
		Chain{Name: "chain"},
	})

	report := generator.GenerateReport(oldState, newState)
	assert.Equal(t, len(report.Entries), 1, "Expected to have 1 entry!")
	assert.Equal(t, report.Entries[0].Type, ProposalQueryError, "Expected to have a proposal query error!")
	assert.Equal(t, report.Entries[0].Value, "test error", "Error text mismatch!")
}

func TestReportGeneratorWithVoteError(t *testing.T) {
	t.Parallel()

	oldState := NewState()
	newState := State{
		ChainInfos: map[string]*ChainInfo{
			"chain": {
				Proposals: map[string]Proposal{
					"proposal": {
						Content: &ProposalContent{},
					},
				},
				ProposalVotes: map[string]WalletVotes{
					"proposal": {
						Votes: map[string]ProposalVote{
							"wallet": {
								Error: "test error",
							},
						},
					},
				},
			},
		},
	}

	generator := NewReportGenerator(stateManager, GetDefaultLogger(), Chains{
		Chain{Name: "chain"},
	})

	report := generator.GenerateReport(oldState, newState)
	assert.Equal(t, len(report.Entries), 1, "Expected to have 1 entry!")
	assert.Equal(t, report.Entries[0].Type, VoteQueryError, "Expected to have a vote query error!")
	assert.Equal(t, report.Entries[0].ProposalID, "proposal", "Proposal ID mismatch!")
}

func TestReportGeneratorWithNotVoted(t *testing.T) {
	t.Parallel()

	oldState := NewState()
	newState := State{
		ChainInfos: map[string]*ChainInfo{
			"chain": {
				Proposals: map[string]Proposal{
					"proposal": {
						Content: &ProposalContent{},
					},
				},
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

	generator := NewReportGenerator(stateManager, GetDefaultLogger(), Chains{
		Chain{Name: "chain"},
	})

	report := generator.GenerateReport(oldState, newState)
	assert.Equal(t, len(report.Entries), 1, "Expected to have 1 entry!")
	assert.Equal(t, report.Entries[0].Type, NotVoted, "Expected to have not voted type!")
	assert.Equal(t, report.Entries[0].ProposalID, "proposal", "Proposal ID mismatch!")
}

func TestReportGeneratorWithVoted(t *testing.T) {
	t.Parallel()

	oldState := State{
		ChainInfos: map[string]*ChainInfo{
			"chain": {
				Proposals: map[string]Proposal{
					"proposal": {
						Content: &ProposalContent{},
					},
				},
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
	newState := State{
		ChainInfos: map[string]*ChainInfo{
			"chain": {
				Proposals: map[string]Proposal{
					"proposal": {
						Content: &ProposalContent{},
					},
				},
				ProposalVotes: map[string]WalletVotes{
					"proposal": {
						Votes: map[string]ProposalVote{
							"wallet": {
								Vote: &Vote{
									Option: "YES",
								},
							},
						},
					},
				},
			},
		},
	}

	generator := NewReportGenerator(stateManager, GetDefaultLogger(), Chains{
		Chain{Name: "chain"},
	})

	report := generator.GenerateReport(oldState, newState)
	assert.Equal(t, len(report.Entries), 1, "Expected to have 1 entry!")
	assert.Equal(t, report.Entries[0].Type, Voted, "Expected to have voted type!")
	assert.Equal(t, report.Entries[0].ProposalID, "proposal", "Proposal ID mismatch!")
}

func TestReportGeneratorWithRevoted(t *testing.T) {
	t.Parallel()

	oldState := State{
		ChainInfos: map[string]*ChainInfo{
			"chain": {
				Proposals: map[string]Proposal{
					"proposal": {
						Content: &ProposalContent{},
					},
				},
				ProposalVotes: map[string]WalletVotes{
					"proposal": {
						Votes: map[string]ProposalVote{
							"wallet": {
								Vote: &Vote{
									Option: "NO",
								},
							},
						},
					},
				},
			},
		},
	}
	newState := State{
		ChainInfos: map[string]*ChainInfo{
			"chain": {
				Proposals: map[string]Proposal{
					"proposal": {
						Content: &ProposalContent{},
					},
				},
				ProposalVotes: map[string]WalletVotes{
					"proposal": {
						Votes: map[string]ProposalVote{
							"wallet": {
								Vote: &Vote{
									Option: "YES",
								},
							},
						},
					},
				},
			},
		},
	}

	generator := NewReportGenerator(stateManager, GetDefaultLogger(), Chains{
		Chain{Name: "chain"},
	})

	report := generator.GenerateReport(oldState, newState)
	assert.Equal(t, len(report.Entries), 1, "Expected to have 1 entry!")
	assert.Equal(t, report.Entries[0].Type, Revoted, "Expected to have revoted type!")
	assert.Equal(t, report.Entries[0].ProposalID, "proposal", "Proposal ID mismatch!")
}
