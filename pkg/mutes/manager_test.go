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

	muted1, err1 := manager.IsEntryMuted(events.VotedEvent{
		Chain: &types.Chain{Name: "chain"},
	})

	assert.True(t, muted1)
	require.NoError(t, err1)

	muted2, err2 := manager.IsEntryMuted(events.VotedEvent{
		Chain: &types.Chain{Name: "chain2"},
	})
	assert.False(t, muted2)
	require.NoError(t, err2)

	deleted, err := manager.DeleteMute(&types.Mute{
		Chain: null.StringFrom("chain"),
	})
	assert.True(t, deleted)
	require.NoError(t, err)

	muted3, err3 := manager.IsEntryMuted(events.VotedEvent{
		Chain:    &types.Chain{Name: "chain"},
		Proposal: types.Proposal{ID: "proposal"},
	})
	assert.False(t, muted3)
	require.NoError(t, err3)

	muted4, err4 := manager.IsEntryMuted(events.VotedEvent{
		Chain:    &types.Chain{Name: "chain2"},
		Proposal: types.Proposal{ID: "proposal"},
	})
	assert.False(t, muted4)
	require.NoError(t, err4)

	deleted2, err2 := manager.DeleteMute(&types.Mute{
		Chain: null.StringFrom("chain"),
	})
	assert.False(t, deleted2)
	require.NoError(t, err2)
}

func TestMuteManagerIsNotAlert(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	db := &databasePkg.StubDatabase{}
	manager := NewMutesManager(log, db)

	muted, err := manager.IsEntryMuted(events.ProposalsQueryErrorEvent{
		Chain: &types.Chain{Name: "chain"},
	})

	require.NoError(t, err)
	assert.False(t, muted)
}

func TestMuteManagerGetAllMutes(t *testing.T) {
	t.Parallel()

	log := logger.GetNopLogger()
	db := &databasePkg.StubDatabase{}
	manager := NewMutesManager(log, db)

	allMutes1, err1 := manager.GetAllMutes()
	assert.Empty(t, allMutes1)
	require.NoError(t, err1)

	err := manager.AddMute(&types.Mute{
		Chain:   null.StringFrom("chain"),
		Expires: time.Now().Add(time.Hour),
	})
	require.NoError(t, err)

	allMutes2, err2 := manager.GetAllMutes()
	assert.NotEmpty(t, allMutes2)
	require.NoError(t, err2)
}
