package database

import (
	"context"
	"main/pkg/types"
	"testing"
)

func TestStubVarious(t *testing.T) {
	// this one is purely to get more code coverage
	t.Parallel()

	db := &StubDatabase{}
	db.Migrate()
	db.Rollback()
	_, _ = db.IsMuted("chain", "proposal")
	_, _ = db.DeleteMute(&types.Mute{})
	_ = db.UpsertVote(
		&types.Chain{Name: "chain"},
		types.Proposal{ID: "proposal1"},
		&types.Wallet{Address: "address"},
		&types.Vote{},
		context.Background(),
	)
	_, _ = db.GetVote(
		&types.Chain{Name: "chain"},
		types.Proposal{ID: "proposal2"},
		&types.Wallet{Address: "address"},
	)
}
