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
func TestTallyFail(t *testing.T) {
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
		"https://example.com/cosmwasm/wasm/v1/contract/neutron1436kxs0w2es6xlqpp9rd35e3d0cjnw4sv8j3a7483sgks29jqwgshlt6zh/smart/eyJyZXZlcnNlX3Byb3Bvc2FscyI6IHsibGltaXQiOiAxMDAwfX0=",
		httpmock.NewErrorResponder(errors.New("custom error")),
	)

	tallies, err := fetcher.GetTallies(context.Background())

	require.Error(t, err)
	require.ErrorContains(t, err, "custom error")
	require.Empty(t, tallies.TallyInfos)
	require.Nil(t, tallies.Chain)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestTallyOk(t *testing.T) {
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
		"https://example.com/cosmwasm/wasm/v1/contract/neutron1436kxs0w2es6xlqpp9rd35e3d0cjnw4sv8j3a7483sgks29jqwgshlt6zh/smart/eyJyZXZlcnNlX3Byb3Bvc2FscyI6IHsibGltaXQiOiAxMDAwfX0=",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("neutron-proposals.json")),
	)

	tallies, err := fetcher.GetTallies(context.Background())
	require.NoError(t, err)
	require.NotNil(t, tallies.Chain)
	require.Len(t, tallies.TallyInfos, 2)
}
