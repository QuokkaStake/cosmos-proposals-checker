package report

import (
	"context"
	"errors"
	databasePkg "main/pkg/database"
	"main/pkg/events"
	fetchersPkg "main/pkg/fetchers"
	loggerPkg "main/pkg/logger"
	"main/pkg/tracing"
	"main/pkg/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGeneratorNew(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	db := &databasePkg.StubDatabase{}
	chains := types.Chains{{Name: "chain"}}
	generator := NewReportNewGenerator(logger, chains, db, tracer)
	require.NotNil(t, generator)
}

func TestGeneratorProposalsError(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	db := &databasePkg.StubDatabase{}
	chains := types.Chains{{
		Name:    "chain",
		Wallets: []*types.Wallet{{Address: "address"}},
	}}
	generator := &Generator{
		Logger:   *logger,
		Chains:   chains,
		Database: db,
		Tracer:   tracer,
		Fetchers: map[string]fetchersPkg.Fetcher{
			"chain": &fetchersPkg.TestFetcher{WithProposalsError: true},
		},
	}

	report := generator.GenerateReport(context.Background())
	require.Len(t, report.Entries, 1)

	firstEntry, ok := report.Entries[0].(events.ProposalsQueryErrorEvent)
	require.True(t, ok)
	require.NotNil(t, firstEntry)
}

func TestGeneratorProposalsGetLastBlockHeightError(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	db := &databasePkg.StubDatabase{
		LastHeightQueryErrors: map[string]map[string]error{
			"chain": {
				"proposals": errors.New("custom error"),
			},
		},
	}
	chains := types.Chains{{
		Name:    "chain",
		Wallets: []*types.Wallet{{Address: "address"}},
	}}
	generator := &Generator{
		Logger:   *logger,
		Chains:   chains,
		Database: db,
		Tracer:   tracer,
		Fetchers: map[string]fetchersPkg.Fetcher{
			"chain": &fetchersPkg.TestFetcher{},
		},
	}

	report := generator.GenerateReport(context.Background())
	require.Len(t, report.Entries, 1)

	firstEntry, ok := report.Entries[0].(events.GenericErrorEvent)
	require.True(t, ok)
	require.NotNil(t, firstEntry)
	require.ErrorContains(t, firstEntry.Error, "custom error")
}

func TestGeneratorProposalGetLastBlockHeightError(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	db := &databasePkg.StubDatabase{
		GetProposalError:     errors.New("custom error"),
		LastHeightWriteError: errors.New("write error"),
	}
	chains := types.Chains{{
		Name:    "chain",
		Wallets: []*types.Wallet{{Address: "address"}},
	}}
	generator := &Generator{
		Logger:   *logger,
		Chains:   chains,
		Database: db,
		Tracer:   tracer,
		Fetchers: map[string]fetchersPkg.Fetcher{
			"chain": &fetchersPkg.TestFetcher{},
		},
	}

	report := generator.GenerateReport(context.Background())
	require.Len(t, report.Entries, 1)

	firstEntry, ok := report.Entries[0].(events.GenericErrorEvent)
	require.True(t, ok)
	require.NotNil(t, firstEntry)
	require.ErrorContains(t, firstEntry.Error, "custom error")
}

func TestGeneratorProposalNotInVoting(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	db := &databasePkg.StubDatabase{}
	chains := types.Chains{{
		Name:    "chain",
		Wallets: []*types.Wallet{{Address: "address"}},
	}}
	generator := &Generator{
		Logger:   *logger,
		Chains:   chains,
		Database: db,
		Tracer:   tracer,
		Fetchers: map[string]fetchersPkg.Fetcher{
			"chain": &fetchersPkg.TestFetcher{WithPassedProposals: true},
		},
	}

	report := generator.GenerateReport(context.Background())
	require.Empty(t, report.Entries)
}

func TestGeneratorProposalFinishedVoting(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	db := &databasePkg.StubDatabase{
		UpsertProposalError: errors.New("write error"),
		Proposals: map[string]map[string]*types.Proposal{
			"chain": {
				"1": &types.Proposal{Status: types.ProposalStatusVoting},
			},
		},
	}
	chains := types.Chains{{
		Name:    "chain",
		Wallets: []*types.Wallet{{Address: "address"}},
	}}
	generator := &Generator{
		Logger:   *logger,
		Chains:   chains,
		Database: db,
		Tracer:   tracer,
		Fetchers: map[string]fetchersPkg.Fetcher{
			"chain": &fetchersPkg.TestFetcher{WithPassedProposals: true},
		},
	}

	report := generator.GenerateReport(context.Background())
	require.Len(t, report.Entries, 1)
	firstEntry, ok := report.Entries[0].(events.FinishedVotingEvent)
	require.True(t, ok)
	require.NotNil(t, firstEntry)
}

func TestGeneratorProposalVoteLastHeightQueryError(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	db := &databasePkg.StubDatabase{
		LastHeightQueryErrors: map[string]map[string]error{
			"chain": {
				"proposal_1_vote_address": errors.New("custom error"),
			},
		},
	}
	chains := types.Chains{{
		Name:    "chain",
		Wallets: []*types.Wallet{{Address: "address"}},
	}}
	generator := &Generator{
		Logger:   *logger,
		Chains:   chains,
		Database: db,
		Tracer:   tracer,
		Fetchers: map[string]fetchersPkg.Fetcher{
			"chain": &fetchersPkg.TestFetcher{},
		},
	}

	report := generator.GenerateReport(context.Background())
	require.Len(t, report.Entries, 1)

	firstEntry, ok := report.Entries[0].(events.GenericErrorEvent)
	require.True(t, ok)
	require.NotNil(t, firstEntry)
	require.ErrorContains(t, firstEntry.Error, "custom error")
}

func TestGeneratorProposalVoteFetchVoteError(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	db := &databasePkg.StubDatabase{}
	chains := types.Chains{{
		Name:    "chain",
		Wallets: []*types.Wallet{{Address: "address"}},
	}}
	generator := &Generator{
		Logger:   *logger,
		Chains:   chains,
		Database: db,
		Tracer:   tracer,
		Fetchers: map[string]fetchersPkg.Fetcher{
			"chain": &fetchersPkg.TestFetcher{WithVoteError: true},
		},
	}

	report := generator.GenerateReport(context.Background())
	require.Len(t, report.Entries, 1)

	firstEntry, ok := report.Entries[0].(events.VoteQueryError)
	require.True(t, ok)
	require.NotNil(t, firstEntry)
	require.ErrorContains(t, firstEntry.Error, "vote query error")
}

func TestGeneratorProposalVoteNotVoted(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	db := &databasePkg.StubDatabase{}
	chains := types.Chains{{
		Name:    "chain",
		Wallets: []*types.Wallet{{Address: "address"}},
	}}
	generator := &Generator{
		Logger:   *logger,
		Chains:   chains,
		Database: db,
		Tracer:   tracer,
		Fetchers: map[string]fetchersPkg.Fetcher{
			"chain": &fetchersPkg.TestFetcher{},
		},
	}

	report := generator.GenerateReport(context.Background())
	require.Len(t, report.Entries, 1)

	firstEntry, ok := report.Entries[0].(events.NotVotedEvent)
	require.True(t, ok)
	require.NotNil(t, firstEntry)
}

func TestGeneratorProposalVoteGetVoteError(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	db := &databasePkg.StubDatabase{
		GetVoteError: errors.New("get vote error"),
	}
	chains := types.Chains{{
		Name:    "chain",
		Wallets: []*types.Wallet{{Address: "address"}},
	}}
	generator := &Generator{
		Logger:   *logger,
		Chains:   chains,
		Database: db,
		Tracer:   tracer,
		Fetchers: map[string]fetchersPkg.Fetcher{
			"chain": &fetchersPkg.TestFetcher{WithVote: true},
		},
	}

	report := generator.GenerateReport(context.Background())
	require.Len(t, report.Entries, 1)

	firstEntry, ok := report.Entries[0].(events.GenericErrorEvent)
	require.True(t, ok)
	require.NotNil(t, firstEntry)
	require.ErrorContains(t, firstEntry.Error, "get vote error")
}

func TestGeneratorProposalVoteVoted(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	db := &databasePkg.StubDatabase{
		UpsertVoteError: errors.New("write error"),
	}
	chains := types.Chains{{
		Name:    "chain",
		Wallets: []*types.Wallet{{Address: "address"}},
	}}
	generator := &Generator{
		Logger:   *logger,
		Chains:   chains,
		Database: db,
		Tracer:   tracer,
		Fetchers: map[string]fetchersPkg.Fetcher{
			"chain": &fetchersPkg.TestFetcher{WithVote: true},
		},
	}

	report := generator.GenerateReport(context.Background())
	require.Len(t, report.Entries, 1)

	firstEntry, ok := report.Entries[0].(events.VotedEvent)
	require.True(t, ok)
	require.NotNil(t, firstEntry)
}

func TestGeneratorProposalVoteRevoted(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	db := &databasePkg.StubDatabase{
		Votes: map[string]map[string]map[string]*types.Vote{
			"chain": {
				"1": {
					"address": &types.Vote{
						Options: types.VoteOptions{
							{Option: "YES", Weight: 1},
						},
					},
				},
			},
		},
	}
	chains := types.Chains{{
		Name:    "chain",
		Wallets: []*types.Wallet{{Address: "address"}},
	}}
	generator := &Generator{
		Logger:   *logger,
		Chains:   chains,
		Database: db,
		Tracer:   tracer,
		Fetchers: map[string]fetchersPkg.Fetcher{
			"chain": &fetchersPkg.TestFetcher{WithVote: true},
		},
	}

	report := generator.GenerateReport(context.Background())
	require.Len(t, report.Entries, 1)

	firstEntry, ok := report.Entries[0].(events.RevotedEvent)
	require.True(t, ok)
	require.NotNil(t, firstEntry)
}

func TestGeneratorProposalVoteAlreadyVoted(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	db := &databasePkg.StubDatabase{
		LastHeightWriteError: errors.New("custom error"),
		Votes: map[string]map[string]map[string]*types.Vote{
			"chain": {
				"1": {
					"address": &types.Vote{
						Options: types.VoteOptions{},
					},
				},
			},
		},
	}
	chains := types.Chains{{
		Name:    "chain",
		Wallets: []*types.Wallet{{Address: "address"}},
	}}
	generator := &Generator{
		Logger:   *logger,
		Chains:   chains,
		Database: db,
		Tracer:   tracer,
		Fetchers: map[string]fetchersPkg.Fetcher{
			"chain": &fetchersPkg.TestFetcher{WithVote: true},
		},
	}

	report := generator.GenerateReport(context.Background())
	require.Empty(t, report.Entries)
}
