package neutron

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
func TestParamsFail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config := &types.Chain{
		Name:                 "chain",
		LCDEndpoints:         []string{"https://example.com"},
		NeutronSmartContract: "neutron1436kxs0w2es6xlqpp9rd35e3d0cjnw4sv8j3a7483sgks29jqwgshlt6zh",
	}
	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()

	fetcher := NewFetcher(config, logger, tracer)

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmwasm/wasm/v1/contract/neutron1436kxs0w2es6xlqpp9rd35e3d0cjnw4sv8j3a7483sgks29jqwgshlt6zh/smart/eyJjb25maWciOnt9fQ==",
		httpmock.NewErrorResponder(errors.New("custom error")),
	)

	params, errs := fetcher.GetChainParams(context.Background())

	require.Len(t, errs, 1)
	require.ErrorContains(t, errs[0], "custom error")
	require.Nil(t, params)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestParamsOk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	config := &types.Chain{
		Name:                 "chain",
		LCDEndpoints:         []string{"https://example.com"},
		NeutronSmartContract: "neutron1436kxs0w2es6xlqpp9rd35e3d0cjnw4sv8j3a7483sgks29jqwgshlt6zh",
	}
	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()

	fetcher := NewFetcher(config, logger, tracer)

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmwasm/wasm/v1/contract/neutron1436kxs0w2es6xlqpp9rd35e3d0cjnw4sv8j3a7483sgks29jqwgshlt6zh/smart/eyJjb25maWciOnt9fQ==",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("neutron-params.json")),
	)

	params, err := fetcher.GetChainParams(context.Background())

	require.Empty(t, err)
	require.NotNil(t, params)
}
