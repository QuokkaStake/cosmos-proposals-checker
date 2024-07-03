package http

import (
	"errors"
	"main/assets"
	"main/pkg/constants"
	loggerPkg "main/pkg/logger"
	"main/pkg/tracing"
	"main/pkg/types"
	"net/http"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/require"
)

func TestHttpClientErrorCreating(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	client := NewClient("chain", []string{"https://example.com"}, logger, tracer)
	_, err := client.GetFull("://", nil, types.HTTPPredicateCheckHeightAfter(100), nil)
	require.Error(t, err)
	require.ErrorContains(t, err, "missing protocol scheme")
}

//nolint:paralleltest // disabled due to httpmock usage
func TestHttpClientQueryFail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://example.com",
		httpmock.NewErrorResponder(errors.New("custom error")),
	)
	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	client := NewClient("chain", []string{"https://example.com"}, logger, tracer)

	var response interface{}
	_, err := client.GetFull("https://example.com", &response, types.HTTPPredicateCheckHeightAfter(100), nil)
	require.Error(t, err)
	require.ErrorContains(t, err, "custom error")
}

//nolint:paralleltest // disabled due to httpmock usage
func TestHttpClientPredicateFail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://example.com",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("lcd-error.json")).HeaderAdd(http.Header{
			constants.HeaderBlockHeight: []string{"1"},
		}),
	)
	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	client := NewClient("chain", []string{"https://example.com"}, logger, tracer)

	var response interface{}
	_, err := client.GetFull("https://example.com", &response, types.HTTPPredicateCheckHeightAfter(100), nil)
	require.Error(t, err)
	require.ErrorContains(t, err, "is bigger than the current height")
}

//nolint:paralleltest // disabled due to httpmock usage
func TestHttpClientJsonParseFail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://example.com",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("invalid-json.json")),
	)
	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	client := NewClient("chain", []string{"https://example.com"}, logger, tracer)

	var response interface{}
	_, err := client.GetFull("https://example.com", &response, types.HTTPPredicateAlwaysPass(), nil)
	require.Error(t, err)
	require.ErrorContains(t, err, "invalid character")
}

//nolint:paralleltest // disabled due to httpmock usage
func TestHttpClientOk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://example.com",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("tally.json")),
	)
	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	client := NewClient("chain", []string{"https://example.com"}, logger, tracer)

	var response interface{}
	_, err := client.GetFull("https://example.com", &response, types.HTTPPredicateAlwaysPass(), nil)
	require.NoError(t, err)
}

//nolint:paralleltest // disabled due to httpmock usage
func TestHttpClientGetMultipleFail(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/",
		httpmock.NewErrorResponder(errors.New("custom error")),
	)
	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	client := NewClient("chain", []string{"https://example.com"}, logger, tracer)

	var response interface{}
	errs := client.Get("/", &response, nil)
	require.Len(t, errs, 1)

	firstErr := errs[0].Error
	require.ErrorContains(t, &firstErr, "custom error")
}

//nolint:paralleltest // disabled due to httpmock usage
func TestHttpClientGetMultipleOk(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	httpmock.RegisterResponder(
		"GET",
		"https://example.com/",
		httpmock.NewBytesResponder(200, assets.GetBytesOrPanic("tally.json")),
	)
	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()
	client := NewClient("chain", []string{"https://example.com"}, logger, tracer)

	var response interface{}
	errs := client.Get("/", &response, nil)
	require.Empty(t, errs)
}
