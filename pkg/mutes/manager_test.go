package mutesmanager

import (
	databasePkg "main/pkg/database"
	"main/pkg/events"
	"main/pkg/logger"
	"main/pkg/types"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/guregu/null/v5"

	"github.com/stretchr/testify/assert"
)

func TestMuteManagerAddAndDeleteMuteIsMuted(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	db := &databasePkg.StubDatabase{}
	manager := NewMutesManager(log, db)

	err := manager.AddMute(&types.Mute{
		Chain:   null.StringFrom("chain"),
		Expires: time.Now().Add(time.Hour),
	})
	require.NoError(t, err)

	assert.True(t, manager.IsEntryMuted(events.VotedEvent{
		Chain:    &types.Chain{Name: "chain"},
		Proposal: types.Proposal{ID: "proposal"},
	}))
	assert.False(t, manager.IsEntryMuted(events.VotedEvent{
		Chain:    &types.Chain{Name: "chain2"},
		Proposal: types.Proposal{ID: "proposal"},
	}))

	deleted, err := manager.DeleteMute(&types.Mute{
		Chain: null.StringFrom("chain"),
	})
	assert.True(t, deleted)
	require.NoError(t, err)

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
	db := &databasePkg.StubDatabase{}
	manager := NewMutesManager(log, db)

	err := manager.AddMute(&types.Mute{
		Chain:   null.StringFrom("chain"),
		Expires: time.Now().Add(time.Hour),
	})
	require.NoError(t, err)

	assert.False(t, manager.IsEntryMuted(events.VotedEvent{
		Chain:    &types.Chain{Name: "chain"},
		Proposal: types.Proposal{ID: "proposal"},
	}))
}

func TestMuteManagerIsNotAlert(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	db := &databasePkg.StubDatabase{}
	manager := NewMutesManager(log, db)

	assert.False(t, manager.IsEntryMuted(events.ProposalsQueryErrorEvent{
		Chain: &types.Chain{Name: "chain"},
	}))
}
