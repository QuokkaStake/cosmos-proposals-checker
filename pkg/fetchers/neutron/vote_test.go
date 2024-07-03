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
func TestVoteFail(t *testing.T) {
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
		"https://example.com/cosmwasm/wasm/v1/contract/neutron1436kxs0w2es6xlqpp9rd35e3d0cjnw4sv8j3a7483sgks29jqwgshlt6zh/smart/eyJnZXRfdm90ZSI6eyJwcm9wb3NhbF9pZCI6OTM2LCJ2b3RlciI6Im5ldXRyb24xeHF6OXBlbXo1ZTV6eWNhYTg5a3lzNWF3Nm04cmhnc3YwN3ZhM2QifX0=",
		httpmock.NewErrorResponder(errors.New("custom error")),
	)

	vote, height, err := fetcher.GetVote(
		"936",
		"neutron1xqz9pemz5e5zycaa89kys5aw6m8rhgsv07va3d",
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
		Name:                 "chain",
		LCDEndpoints:         []string{"https://example.com"},
		NeutronSmartContract: "neutron1436kxs0w2es6xlqpp9rd35e3d0cjnw4sv8j3a7483sgks29jqwgshlt6zh",
	}
	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()

	fetcher := NewFetcher(config, logger, tracer)

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/cosmwasm/wasm/v1/contract/neutron1436kxs0w2es6xlqpp9rd35e3d0cjnw4sv8j3a7483sgks29jqwgshlt6zh/smart/eyJnZXRfdm90ZSI6eyJwcm9wb3NhbF9pZCI6OTM2LCJ2b3RlciI6Im5ldXRyb24xeHF6OXBlbXo1ZTV6eWNhYTg5a3lzNWF3Nm04cmhnc3YwN3ZhM2QifX0=",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("neutron-vote-not-found.json")),
	)

	vote, height, err := fetcher.GetVote(
		"936",
		"neutron1xqz9pemz5e5zycaa89kys5aw6m8rhgsv07va3d",
		0,
		context.Background(),
	)

	require.Error(t, err)
	require.Zero(t, height)
	require.Nil(t, vote)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestVoteOk(t *testing.T) {
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
		"https://example.com/cosmwasm/wasm/v1/contract/neutron1436kxs0w2es6xlqpp9rd35e3d0cjnw4sv8j3a7483sgks29jqwgshlt6zh/smart/eyJnZXRfdm90ZSI6eyJwcm9wb3NhbF9pZCI6NDIsInZvdGVyIjoibmV1dHJvbjEwM2wwMjVmdzJ4N2Q4azhweDdkMDd2dTY2N3gwaDVwNWdrZTR5diJ9fQ==",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("neutron-vote-ok.json")),
	)

	vote, height, err := fetcher.GetVote(
		"42",
		"neutron103l025fw2x7d8k8px7d07vu667x0h5p5gke4yv",
		0,
		context.Background(),
	)

	require.Nil(t, err)
	require.Zero(t, height)
	require.NotNil(t, vote)
	require.Equal(t, types.VoteOptions{
		{Option: "👌Yes", Weight: 1},
	}, vote.Options)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestVoteUnknown(t *testing.T) {
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
		"https://example.com/cosmwasm/wasm/v1/contract/neutron1436kxs0w2es6xlqpp9rd35e3d0cjnw4sv8j3a7483sgks29jqwgshlt6zh/smart/eyJnZXRfdm90ZSI6eyJwcm9wb3NhbF9pZCI6NDIsInZvdGVyIjoibmV1dHJvbjEwM2wwMjVmdzJ4N2Q4azhweDdkMDd2dTY2N3gwaDVwNWdrZTR5diJ9fQ==",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("neutron-vote-unknown.json")),
	)

	vote, height, err := fetcher.GetVote(
		"42",
		"neutron103l025fw2x7d8k8px7d07vu667x0h5p5gke4yv",
		0,
		context.Background(),
	)

	require.Nil(t, err)
	require.Zero(t, height)
	require.NotNil(t, vote)
	require.Equal(t, types.VoteOptions{
		{Option: "unknown", Weight: 1},
	}, vote.Options)
}
