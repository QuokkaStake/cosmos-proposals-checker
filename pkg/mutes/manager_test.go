package mutesmanager

import (
	"github.com/stretchr/testify/assert"
	"main/pkg/fs"
	"main/pkg/logger"
	"testing"
	"time"
)

func TestMuteManagerLoadWithoutPath(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	filesystem := &fs.TestFS{}

	manager := NewMutesManager("", filesystem, log)
	manager.Load()

	assert.Empty(t, manager.Mutes.Mutes)
}

func TestMuteManagerLoadNotExistingPath(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	filesystem := &fs.TestFS{}

	manager := NewMutesManager("not-existing.json", filesystem, log)
	manager.Load()

	assert.Empty(t, manager.Mutes.Mutes)
}

func TestMuteManagerLoadInvalidJson(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	filesystem := &fs.TestFS{}

	manager := NewMutesManager("invalid-json.json", filesystem, log)
	manager.Load()

	assert.Empty(t, manager.Mutes.Mutes)
}

func TestMuteManagerLoadValidJson(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	filesystem := &fs.TestFS{}

	manager := NewMutesManager("valid-mutes.json", filesystem, log)
	manager.Load()

	assert.NotEmpty(t, manager.Mutes.Mutes)
}

func TestMuteManagerSaveWithoutPath(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	filesystem := &fs.TestFS{}

	manager := NewMutesManager("", filesystem, log)
	manager.Load()
	manager.Save()

	assert.Empty(t, manager.Mutes.Mutes)
}

func TestMuteManagerSaveWithError(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	filesystem := &fs.TestFS{WithWriteError: true}

	manager := NewMutesManager("out.json", filesystem, log)
	manager.Load()
	manager.Save()

	assert.Empty(t, manager.Mutes.Mutes)
}

func TestMuteManagerSaveWithoutError(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	filesystem := &fs.TestFS{}

	manager := NewMutesManager("out.json", filesystem, log)
	manager.Load()
	manager.Save()

	assert.Empty(t, manager.Mutes.Mutes)
}

func TestMuteManagerAddMuteIsMuted(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	filesystem := &fs.TestFS{}

	manager := NewMutesManager("out.json", filesystem, log)
	manager.Load()

	manager.AddMute(&Mute{
		Chain:   "chain",
		Expires: time.Now().Add(time.Hour),
	})

	assert.True(t, manager.IsMuted("chain", "proposal"))
	assert.False(t, manager.IsMuted("chain2", "proposal"))
}

func TestMuteManagerIsMutedNoPath(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	filesystem := &fs.TestFS{}

	manager := NewMutesManager("", filesystem, log)
	manager.Load()

	manager.AddMute(&Mute{
		Chain:   "chain",
		Expires: time.Now().Add(time.Hour),
	})

	assert.False(t, manager.IsMuted("chain", "proposal"))
}
