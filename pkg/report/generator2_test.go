package report

import (
	"context"
	databasePkg "main/pkg/database"
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

func TestGeneratorProposalError(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	db := &databasePkg.StubDatabase{}
	chains := types.Chains{{
		Name:    "chain",
		Wallets: []*types.Wallet{{Address: "address"}},
	}}
	generator := &NewGenerator{
		Logger:   *logger,
		Chains:   chains,
		Database: db,
		Tracer:   tracer,
		Fetchers: map[string]fetchersPkg.Fetcher{
			"chain": &fetchersPkg.TestFetcher{WithProposalsError: true},
		},
	}

	report := generator.GenerateReport(context.Background())
	require.NotEmpty(t, report.Entries)
}
