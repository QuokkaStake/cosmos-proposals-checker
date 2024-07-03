package cosmos

import (
	"context"
	"errors"
	"main/assets"
	loggerPkg "main/pkg/logger"
	"main/pkg/tracing"
	"main/pkg/types"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest // disabled due to httpmock usage
func TestTallyAllFail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config := &types.Chain{
		Name:         "chain",
		LCDEndpoints: []string{"https://example.com"},
	}
	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()

	fetcher := NewRPC(config, logger, tracer)

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/gov/v1beta1/proposals?pagination.limit=1000&pagination.offset=0&pagination.count_total=1",
		httpmock.NewErrorResponder(errors.New("custom error")),
	)

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/staking/v1beta1/pool",
		httpmock.NewErrorResponder(errors.New("custom error")),
	)

	tallies, err := fetcher.GetTallies(
		context.Background(),
	)

	require.Error(t, err)
	require.Empty(t, tallies.TallyInfos)
	require.Nil(t, tallies.Chain)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestTallyTallyError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config := &types.Chain{
		Name:          "chain",
		LCDEndpoints:  []string{"https://example.com"},
		ProposalsType: "v1",
	}
	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()

	fetcher := NewRPC(config, logger, tracer)
	fetcher.PaginationLimit = 100

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/staking/v1beta1/pool",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("staking_pool.json")),
	)

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/gov/v1/proposals?pagination.limit=100&pagination.offset=0&pagination.count_total=1",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("proposals_v1_page1.json")),
	)

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/gov/v1/proposals?pagination.limit=100&pagination.offset=100&pagination.count_total=1",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("proposals_v1_page2.json")),
	)

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/gov/v1/proposals?pagination.limit=100&pagination.offset=200&pagination.count_total=1",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("proposals_v1_page3.json")),
	)

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/gov/v1beta1/proposals/936/tally",
		httpmock.NewErrorResponder(errors.New("custom error")),
	)

	tallies, err := fetcher.GetTallies(
		context.Background(),
	)

	require.Error(t, err)
	require.Empty(t, tallies.TallyInfos)
	require.Nil(t, tallies.Chain)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestTallyAllOk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config := &types.Chain{
		Name:          "chain",
		LCDEndpoints:  []string{"https://example.com"},
		ProposalsType: "v1",
	}
	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()

	fetcher := NewRPC(config, logger, tracer)
	fetcher.PaginationLimit = 100

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/staking/v1beta1/pool",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("staking_pool.json")),
	)

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/gov/v1/proposals?pagination.limit=100&pagination.offset=0&pagination.count_total=1",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("proposals_v1_page1.json")),
	)

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/gov/v1/proposals?pagination.limit=100&pagination.offset=100&pagination.count_total=1",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("proposals_v1_page2.json")),
	)

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/gov/v1/proposals?pagination.limit=100&pagination.offset=200&pagination.count_total=1",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("proposals_v1_page3.json")),
	)

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/gov/v1beta1/proposals/936/tally",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("tally.json")),
	)

	tallies, err := fetcher.GetTallies(
		context.Background(),
	)

	require.NoError(t, err)
	require.NotEmpty(t, tallies.TallyInfos)
	require.NotNil(t, tallies.Chain)
}
