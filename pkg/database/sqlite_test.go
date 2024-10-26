package database

import (
	"context"
	"main/pkg/logger"
	"main/pkg/types"
	"testing"
	"time"

	"github.com/guregu/null/v5"
	"github.com/stretchr/testify/require"
)

//nolint:paralleltest
func TestSqliteMigrations(t *testing.T) {
	db := NewSqliteDatabase(logger.GetNopLogger(), types.DatabaseConfig{Path: "db.sqlite"})
	db.Init()

	db.Rollback() // to catch the error when there's nothing to rollback
	db.Migrate()
	db.Rollback()
	db.Migrate()
	db.Migrate() // to catch the error when there's nothing to migrate

	err := db.Destroy()
	require.NoError(t, err)
}

//nolint:paralleltest
func TestSqliteProposal(t *testing.T) {
	db := NewSqliteDatabase(logger.GetNopLogger(), types.DatabaseConfig{Path: "db.sqlite"})
	db.Init()
	db.Migrate()

	chain := &types.Chain{Name: "chain"}
	proposal := types.Proposal{ID: "proposal"}

	proposalFromDB, err := db.GetProposal(chain, proposal.ID)
	require.Nil(t, proposalFromDB)
	require.NoError(t, err)

	err = db.UpsertProposal(chain, proposal)
	require.NoError(t, err)

	proposalFromDB2, err := db.GetProposal(chain, proposal.ID)
	require.NotNil(t, proposalFromDB2)
	require.NoError(t, err)

	err = db.Destroy()
	require.NoError(t, err)
}

//nolint:paralleltest
func TestSqliteVote(t *testing.T) {
	db := NewSqliteDatabase(logger.GetNopLogger(), types.DatabaseConfig{Path: "db.sqlite"})
	db.Init()
	db.Migrate()

	chain := &types.Chain{Name: "chain"}
	proposal := types.Proposal{ID: "proposal"}
	wallet := &types.Wallet{Address: "address"}
	vote := &types.Vote{Options: types.VoteOptions{
		{Option: "YES", Weight: 1},
	}}

	proposalFromDB, err := db.GetProposal(chain, proposal.ID)
	require.Nil(t, proposalFromDB)
	require.NoError(t, err)

	voteFromDB, err := db.GetVote(chain, proposal, wallet)
	require.Nil(t, voteFromDB)
	require.NoError(t, err)

	err = db.UpsertVote(chain, proposal, wallet, vote, context.Background())
	require.NoError(t, err)

	voteFromDB2, err := db.GetVote(chain, proposal, wallet)
	require.NotNil(t, voteFromDB2)
	require.NoError(t, err)

	err = db.Destroy()
	require.NoError(t, err)
}

//nolint:paralleltest
func TestSqliteLastBlockHeight(t *testing.T) {
	db := NewSqliteDatabase(logger.GetNopLogger(), types.DatabaseConfig{Path: "db.sqlite"})
	db.Init()
	db.Migrate()

	chain := &types.Chain{Name: "chain"}

	entryFromDB, err := db.GetLastBlockHeight(chain, "key")
	require.Zero(t, entryFromDB)
	require.NoError(t, err)

	err = db.UpsertLastBlockHeight(chain, "key", 123)
	require.NoError(t, err)

	entryFromDB2, err := db.GetLastBlockHeight(chain, "key")
	require.Equal(t, int64(123), entryFromDB2)
	require.NoError(t, err)

	err = db.Destroy()
	require.NoError(t, err)
}

//nolint:paralleltest
func TestSqliteMutes(t *testing.T) {
	db := NewSqliteDatabase(logger.GetNopLogger(), types.DatabaseConfig{Path: "db.sqlite"})
	db.Init()
	db.Migrate()

	mute := &types.Mute{
		Chain:      null.StringFrom("chain"),
		ProposalID: null.StringFrom("proposal"),
		Expires:    time.Now().Add(time.Hour),
	}

	mutesFromDB, err := db.GetAllMutes()
	require.Empty(t, mutesFromDB)
	require.NoError(t, err)

	isMuted, err := db.IsMuted("chain", "proposal")
	require.False(t, isMuted)
	require.NoError(t, err)

	err = db.UpsertMute(mute)
	require.NoError(t, err)

	mutesFromDB2, err := db.GetAllMutes()
	require.NotEmpty(t, mutesFromDB2)
	require.NoError(t, err)

	isMuted2, err := db.IsMuted("chain", "proposal")
	require.True(t, isMuted2)
	require.NoError(t, err)

	deleted, err := db.DeleteMute(mute)
	require.True(t, deleted)
	require.NoError(t, err)

	deleted2, err := db.DeleteMute(mute)
	require.False(t, deleted2)
	require.NoError(t, err)

	deleted3, err := db.DeleteMute(&types.Mute{
		Chain:      null.NewString("", false),
		ProposalID: null.NewString("", false),
	})
	require.False(t, deleted3)
	require.NoError(t, err)

	deleted4, err := db.DeleteMute(&types.Mute{
		Chain:      null.NewString("chain", true),
		ProposalID: null.NewString("", false),
	})
	require.False(t, deleted4)
	require.NoError(t, err)

	deleted5, err := db.DeleteMute(&types.Mute{
		Chain:      null.NewString("", false),
		ProposalID: null.NewString("proposal", true),
	})
	require.False(t, deleted5)
	require.NoError(t, err)

	err = db.Destroy()
	require.NoError(t, err)
}
