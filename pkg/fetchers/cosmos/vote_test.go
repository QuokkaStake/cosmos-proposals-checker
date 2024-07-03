package cosmos

import (
	"context"
	"errors"
	"main/assets"
	loggerPkg "main/pkg/logger"
	"main/pkg/tracing"
	"main/pkg/types"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/jarcoal/httpmock"
)

//nolint:paralleltest // disabled due to httpmock usage
func TestVoteFail(t *testing.T) {
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
		"https://example.com/cosmos/gov/v1beta1/proposals/936/votes/cosmos1xqz9pemz5e5zycaa89kys5aw6m8rhgsvtp9lt2",
		httpmock.NewErrorResponder(errors.New("custom error")),
	)

	vote, height, err := fetcher.GetVote(
		"936",
		"cosmos1xqz9pemz5e5zycaa89kys5aw6m8rhgsvtp9lt2",
		0,
		context.Background(),
	)

	require.Error(t, err)
	require.NotEmpty(t, err.NodeErrors)
	require.Zero(t, height)
	require.Nil(t, vote)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestVoteNotFound(t *testing.T) {
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
		"https://example.com/cosmos/gov/v1beta1/proposals/936/votes/cosmos1xqz9pemz5e5zycaa89kys5aw6m8rhgsvtp9lt2",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("vote-not-found.json")),
	)

	vote, height, err := fetcher.GetVote(
		"936",
		"cosmos1xqz9pemz5e5zycaa89kys5aw6m8rhgsvtp9lt2",
		0,
		context.Background(),
	)

	require.Error(t, err)
	require.Zero(t, height)
	require.Nil(t, vote)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestVoteLcdError(t *testing.T) {
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
		"https://example.com/cosmos/gov/v1beta1/proposals/936/votes/cosmos1xqz9pemz5e5zycaa89kys5aw6m8rhgsvtp9lt2",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("lcd-error.json")),
	)

	vote, height, err := fetcher.GetVote(
		"936",
		"cosmos1xqz9pemz5e5zycaa89kys5aw6m8rhgsvtp9lt2",
		0,
		context.Background(),
	)

	require.Error(t, err)
	require.Error(t, err.QueryError)
	require.Zero(t, height)
	require.Nil(t, vote)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestVoteOk(t *testing.T) {
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
		"https://example.com/cosmos/gov/v1beta1/proposals/936/votes/cosmos1xqz9pemz5e5zycaa89kys5aw6m8rhgsvtp9lt2",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("vote.json")),
	)

	vote, height, err := fetcher.GetVote(
		"936",
		"cosmos1xqz9pemz5e5zycaa89kys5aw6m8rhgsvtp9lt2",
		0,
		context.Background(),
	)

	require.Nil(t, err)
	require.Zero(t, height)
	require.NotNil(t, vote)
	require.Equal(t, types.VoteOptions{
		{Option: "🤬No with veto", Weight: 1},
	}, vote.Options)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestVoteOkUnknown(t *testing.T) {
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
		"https://example.com/cosmos/gov/v1beta1/proposals/936/votes/cosmos1xqz9pemz5e5zycaa89kys5aw6m8rhgsvtp9lt2",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("vote-unknown.json")),
	)

	vote, height, err := fetcher.GetVote(
		"936",
		"cosmos1xqz9pemz5e5zycaa89kys5aw6m8rhgsvtp9lt2",
		0,
		context.Background(),
	)

	require.Nil(t, err)
	require.Zero(t, height)
	require.NotNil(t, vote)
	require.Equal(t, types.VoteOptions{
		{Option: "test", Weight: 1},
	}, vote.Options)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestVoteOkOld(t *testing.T) {
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
		"https://example.com/cosmos/gov/v1beta1/proposals/936/votes/cosmos1xqz9pemz5e5zycaa89kys5aw6m8rhgsvtp9lt2",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("vote-old.json")),
	)

	vote, height, err := fetcher.GetVote(
		"936",
		"cosmos1xqz9pemz5e5zycaa89kys5aw6m8rhgsvtp9lt2",
		0,
		context.Background(),
	)

	require.Nil(t, err)
	require.Zero(t, height)
	require.NotNil(t, vote)
	require.Equal(t, types.VoteOptions{
		{Option: "🤬No with veto", Weight: 1},
	}, vote.Options)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestVoteOkOldUnknown(t *testing.T) {
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
		"https://example.com/cosmos/gov/v1beta1/proposals/936/votes/cosmos1xqz9pemz5e5zycaa89kys5aw6m8rhgsvtp9lt2",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("vote-old-unknown.json")),
	)

	vote, height, err := fetcher.GetVote(
		"936",
		"cosmos1xqz9pemz5e5zycaa89kys5aw6m8rhgsvtp9lt2",
		0,
		context.Background(),
	)

	require.Nil(t, err)
	require.Zero(t, height)
	require.NotNil(t, vote)
	require.Equal(t, types.VoteOptions{
		{Option: "test", Weight: 1},
	}, vote.Options)
}
