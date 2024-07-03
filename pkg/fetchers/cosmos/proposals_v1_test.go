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
func TestProposalsV1Fail(t *testing.T) {
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

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/gov/v1/proposals?pagination.limit=1000&pagination.offset=0&pagination.count_total=1",
		httpmock.NewErrorResponder(errors.New("custom error")),
	)

	proposals, height, err := fetcher.GetAllProposals(
		0,
		context.Background(),
	)

	require.Error(t, err)
	require.NotEmpty(t, err.NodeErrors)
	require.Len(t, err.NodeErrors, 1)
	firstError := err.NodeErrors[0].Error
	require.ErrorContains(t, &firstError, "custom error")
	require.Zero(t, height)
	require.Empty(t, proposals)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestProposalsV1LcdError(t *testing.T) {
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

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/gov/v1/proposals?pagination.limit=1000&pagination.offset=0&pagination.count_total=1",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("lcd-error.json")),
	)

	proposals, height, err := fetcher.GetAllProposals(
		0,
		context.Background(),
	)

	require.Error(t, err)
	require.Error(t, err.QueryError)
	require.Zero(t, height)
	require.Empty(t, proposals)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestProposalsV1Success(t *testing.T) {
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

	proposals, height, err := fetcher.GetAllProposals(
		0,
		context.Background(),
	)

	require.Error(t, err)
	require.Zero(t, height)
	require.Len(t, proposals, 218)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestProposalsV1SuccessKujira(t *testing.T) {
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

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/gov/v1/proposals?pagination.limit=1000&pagination.offset=0&pagination.count_total=1",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("proposals_v1_kujira.json")),
	)

	proposals, height, err := fetcher.GetAllProposals(
		0,
		context.Background(),
	)

	require.Error(t, err)
	require.Zero(t, height)
	require.Len(t, proposals, 100)
}
