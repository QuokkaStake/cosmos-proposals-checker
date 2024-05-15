package data

import (
	"context"
	fetchersPkg "main/pkg/fetchers"
	"main/pkg/logger"
	"main/pkg/tracing"
	"main/pkg/types"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/stretchr/testify/assert"
)

func TestDataManagerNew(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	dataManager := NewManager(log, types.Chains{
		{Name: "chain"},
	}, tracing.InitNoopTracer())

	assert.NotNil(t, dataManager)
}

func TestDataManagerGetTallyWithError(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	dataManager := &Manager{
		Logger:   *log,
		Chains:   types.Chains{{Name: "chain"}},
		Fetchers: []fetchersPkg.Fetcher{&fetchersPkg.TestFetcher{WithTallyError: true}},
		Tracer:   tracing.InitNoopTracer(),
	}

	tallies, err := dataManager.GetTallies(context.Background())
	require.Error(t, err)
	assert.Empty(t, tallies)
}

func TestDataManagerGetTallyEmpty(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	dataManager := &Manager{
		Logger:   *log,
		Chains:   types.Chains{{Name: "chain"}},
		Fetchers: []fetchersPkg.Fetcher{&fetchersPkg.TestFetcher{}},
		Tracer:   tracing.InitNoopTracer(),
	}

	tallies, err := dataManager.GetTallies(context.Background())
	require.NoError(t, err)
	assert.Empty(t, tallies)
}

func TestDataManagerGetTallyOk(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	dataManager := &Manager{
		Logger:   *log,
		Chains:   types.Chains{{Name: "chain"}},
		Fetchers: []fetchersPkg.Fetcher{&fetchersPkg.TestFetcher{WithTallyNotEmpty: true}},
		Tracer:   tracing.InitNoopTracer(),
	}

	tallies, err := dataManager.GetTallies(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, tallies)
}

func TestDataManagerGetParamsWithError(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	dataManager := &Manager{
		Logger:   *log,
		Chains:   types.Chains{{Name: "chain"}},
		Fetchers: []fetchersPkg.Fetcher{&fetchersPkg.TestFetcher{WithParamsError: true}},
		Tracer:   tracing.InitNoopTracer(),
	}

	params, err := dataManager.GetParams(context.Background())
	require.Error(t, err)
	assert.Empty(t, params)
}

func TestDataManagerGetParamsOk(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	dataManager := &Manager{
		Logger:   *log,
		Chains:   types.Chains{{Name: "chain"}},
		Fetchers: []fetchersPkg.Fetcher{&fetchersPkg.TestFetcher{}},
		Tracer:   tracing.InitNoopTracer(),
	}

	params, err := dataManager.GetParams(context.Background())
	require.NoError(t, err)
	assert.NotEmpty(t, params)
}
