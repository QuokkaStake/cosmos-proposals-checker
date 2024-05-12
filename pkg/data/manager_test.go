package data

import (
	fetchersPkg "main/pkg/fetchers"
	"main/pkg/logger"
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
	})

	assert.NotNil(t, dataManager)
}

func TestDataManagerGetTallyWithError(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	dataManager := &Manager{
		Logger:   *log,
		Chains:   types.Chains{{Name: "chain"}},
		Fetchers: []fetchersPkg.Fetcher{&fetchersPkg.TestFetcher{WithTallyError: true}},
	}

	tallies, err := dataManager.GetTallies()
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
	}

	tallies, err := dataManager.GetTallies()
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
	}

	tallies, err := dataManager.GetTallies()
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
	}

	params, err := dataManager.GetParams()
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
	}

	params, err := dataManager.GetParams()
	require.NoError(t, err)
	assert.NotEmpty(t, params)
}
