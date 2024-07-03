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
	"github.com/stretchr/testify/assert"
)

//nolint:paralleltest // disabled due to httpmock usage
func TestParamsFailAll(t *testing.T) {
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
		"https://example.com/cosmos/gov/v1beta1/params/deposit",
		httpmock.NewErrorResponder(errors.New("custom error")),
	)

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/gov/v1beta1/params/voting",
		httpmock.NewErrorResponder(errors.New("custom error")),
	)

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/gov/v1beta1/params/tallying",
		httpmock.NewErrorResponder(errors.New("custom error")),
	)

	params, errs := fetcher.GetChainParams(
		context.Background(),
	)

	assert.Len(t, errs, 3)
	assert.Nil(t, params)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestParamsFailOk(t *testing.T) {
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
		"https://example.com/cosmos/gov/v1beta1/params/deposit",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("params_deposit.json")),
	)

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/gov/v1beta1/params/voting",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("params_voting.json")),
	)

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmos/gov/v1beta1/params/tallying",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("params_tallying.json")),
	)

	params, errs := fetcher.GetChainParams(
		context.Background(),
	)

	assert.Empty(t, errs)
	assert.NotNil(t, params)
	assert.Len(t, params.Params, 6)
	assert.Equal(t, "14 days", params.Params[0].Serialize())
	assert.Equal(t, "14 days", params.Params[1].Serialize())
	assert.Equal(t, "250000000 uatom", params.Params[2].Serialize())
	assert.Equal(t, "40.00%", params.Params[3].Serialize())
	assert.Equal(t, "50.00%", params.Params[4].Serialize())
	assert.Equal(t, "33.40%", params.Params[5].Serialize())
}
