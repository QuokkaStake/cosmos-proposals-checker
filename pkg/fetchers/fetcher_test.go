package fetchers

import (
	"main/pkg/fetchers/cosmos"
	"main/pkg/fetchers/neutron"
	loggerPkg "main/pkg/logger"
	"main/pkg/tracing"
	"main/pkg/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFetcher(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	tracer := tracing.InitNoopTracer()

	neutronFetcher := GetFetcher(&types.Chain{Type: "neutron"}, logger, tracer)
	assert.IsType(t, &neutron.Fetcher{}, neutronFetcher)

	cosmosFetcher := GetFetcher(&types.Chain{}, logger, tracer)
	assert.IsType(t, &cosmos.RPC{}, cosmosFetcher)
}
