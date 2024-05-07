package state

import (
	"github.com/stretchr/testify/assert"
	"main/pkg/fs"
	"main/pkg/logger"
	"main/pkg/types"
	"testing"
)

func TestStateManagerLoadNotExisting(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	filesystem := &fs.TestFS{}

	manager := NewStateManager("non-existing.json", filesystem, log)
	manager.Load()

	assert.Empty(t, manager.State.ChainInfos)
}

func TestStateManagerLoadFailedUnmarshaling(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	filesystem := &fs.TestFS{}

	manager := NewStateManager("invalid-json.json", filesystem, log)
	manager.Load()

	assert.Empty(t, manager.State.ChainInfos)
}

func TestStateManagerLoadSuccess(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	filesystem := &fs.TestFS{}

	manager := NewStateManager("valid-state.json", filesystem, log)
	manager.Load()

	assert.NotEmpty(t, manager.State.ChainInfos)
}

func TestManagerSaveWithError(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	filesystem := &fs.TestFS{WithWriteError: true}

	manager := NewStateManager("out.json", filesystem, log)
	manager.Load()
	manager.Save()

	assert.Empty(t, manager.State.ChainInfos)
}

func TestManagerSaveWithoutError(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	filesystem := &fs.TestFS{}

	manager := NewStateManager("out.json", filesystem, log)
	manager.Load()
	manager.Save()

	assert.Empty(t, manager.State.ChainInfos)
}

func TestManagerCommitState(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	filesystem := &fs.TestFS{}

	state := NewState()
	state.SetProposal(&types.Chain{Name: "chain"}, types.Proposal{ID: "id"})

	manager := NewStateManager("out.json", filesystem, log)
	manager.Load()
	assert.Empty(t, manager.State.ChainInfos)

	manager.CommitState(state)
	assert.NotEmpty(t, manager.State.ChainInfos)
}
