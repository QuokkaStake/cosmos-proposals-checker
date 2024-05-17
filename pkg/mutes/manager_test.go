package mutesmanager

import (
	"main/pkg/events"
	"main/pkg/fs"
	"main/pkg/logger"
	"main/pkg/types"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
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

func TestMuteManagerAddAndDeleteMuteIsMuted(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	filesystem := &fs.TestFS{}

	manager := NewMutesManager("out.json", filesystem, log)
	manager.Load()

	manager.AddMute(&Mute{
		Chain:   "chain",
		Expires: time.Now().Add(time.Hour),
	})

	assert.True(t, manager.IsEntryMuted(events.VotedEvent{
		Chain:    &types.Chain{Name: "chain"},
		Proposal: types.Proposal{ID: "proposal"},
	}))
	assert.False(t, manager.IsEntryMuted(events.VotedEvent{
		Chain:    &types.Chain{Name: "chain2"},
		Proposal: types.Proposal{ID: "proposal"},
	}))

	deleted := manager.DeleteMute(&Mute{
		Chain: "chain",
	})
	assert.True(t, deleted)

	assert.False(t, manager.IsEntryMuted(events.VotedEvent{
		Chain:    &types.Chain{Name: "chain"},
		Proposal: types.Proposal{ID: "proposal"},
	}))
	assert.False(t, manager.IsEntryMuted(events.VotedEvent{
		Chain:    &types.Chain{Name: "chain2"},
		Proposal: types.Proposal{ID: "proposal"},
	}))
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

	assert.False(t, manager.IsEntryMuted(events.VotedEvent{
		Chain:    &types.Chain{Name: "chain"},
		Proposal: types.Proposal{ID: "proposal"},
	}))
}

func TestMuteManagerIsNotAlert(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	filesystem := &fs.TestFS{}

	manager := NewMutesManager("", filesystem, log)
	manager.Load()

	assert.False(t, manager.IsEntryMuted(events.ProposalsQueryErrorEvent{
		Chain: &types.Chain{Name: "chain"},
	}))
}
