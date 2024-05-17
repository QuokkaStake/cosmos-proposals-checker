package mutesmanager

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMuteMatchesNoParams(t *testing.T) {
	t.Parallel()

	mute := &Mute{}
	muted := mute.Matches("chain", "proposal")
	assert.True(t, muted, "Mute should match!")
}

func TestMuteMatchesWithChainSpecified(t *testing.T) {
	t.Parallel()

	mute := &Mute{Chain: "chain"}
	muted := mute.Matches("chain", "proposal")
	assert.True(t, muted, "Mute should match!")
	muted2 := mute.Matches("chain2", "proposal")
	assert.False(t, muted2, "Mute should not match!")
}

func TestMuteMatchesWithProposalSpecified(t *testing.T) {
	t.Parallel()

	mute := &Mute{ProposalID: "proposal"}
	muted := mute.Matches("chain", "proposal")
	assert.True(t, muted, "Mute should match!")
	muted2 := mute.Matches("chain", "proposal2")
	assert.False(t, muted2, "Mute should not match!")
}

func TestMuteMatchesWithAllSpecified(t *testing.T) {
	t.Parallel()

	mute := &Mute{Chain: "chain", ProposalID: "proposal"}
	muted := mute.Matches("chain", "proposal")
	assert.True(t, muted, "Mute should match!")
	muted2 := mute.Matches("chain", "proposal2")
	assert.False(t, muted2, "Mute should not match!")
	muted3 := mute.Matches("chain2", "proposal")
	assert.False(t, muted3, "Mute should not match!")
}

func TestMutesMatchesIgnoreExpired(t *testing.T) {
	t.Parallel()

	mutes := Mutes{
		Mutes: []*Mute{
			{Expires: time.Now().Add(-time.Hour)},
		},
	}
	muted := mutes.IsMuted("chain", "proposal")
	assert.False(t, muted, "Mute should not be muted!")
}

func TestMutesMatchesNotIgnoreActual(t *testing.T) {
	t.Parallel()

	mutes := Mutes{
		Mutes: []*Mute{
			{Expires: time.Now().Add(time.Hour)},
		},
	}
	muted := mutes.IsMuted("chain", "proposal")
	assert.True(t, muted, "Mute should be muted!")
}

func TestMutesAddsMuteNew(t *testing.T) {
	t.Parallel()

	mutes := Mutes{
		Mutes: []*Mute{
			{Chain: "chain1", Expires: time.Now().Add(-time.Hour)},
		},
	}

	mutes.AddMute(&Mute{Chain: "chain2", Expires: time.Now().Add(time.Hour)})
	assert.Len(t, mutes.Mutes, 1, "There should be 1 mute!")
	assert.Equal(t, "chain2", mutes.Mutes[0].Chain, "Chain name should match!")
}

func TestMutesAddsMuteOverride(t *testing.T) {
	t.Parallel()

	expireTime := time.Now().Add(time.Hour)

	mutes := Mutes{
		Mutes: []*Mute{
			{Chain: "chain1", ProposalID: "proposal1", Expires: expireTime, Comment: "comment1"},
			{Chain: "chain1", ProposalID: "proposal2", Expires: expireTime, Comment: "comment2"},
		},
	}

	newExpireTime := time.Now().Add(3 * time.Hour)

	mutes.AddMute(&Mute{
		Chain:      "chain1",
		ProposalID: "proposal1",
		Expires:    newExpireTime,
		Comment:    "newcomment",
	})
	assert.Len(t, mutes.Mutes, 2)
	assert.Equal(t, "chain1", mutes.Mutes[0].Chain)
	assert.Equal(t, "proposal1", mutes.Mutes[0].ProposalID)
	assert.Equal(t, "newcomment", mutes.Mutes[0].Comment)
	assert.Equal(t, newExpireTime, mutes.Mutes[0].Expires)
}

func TestMutesDeleteMuteNotExisting(t *testing.T) {
	t.Parallel()

	mutes := Mutes{
		Mutes: []*Mute{
			{Chain: "chain1", Expires: time.Now().Add(time.Hour)},
		},
	}

	deleted := mutes.DeleteMute(&Mute{Chain: "chain2"})
	assert.False(t, deleted)
	assert.Len(t, mutes.Mutes, 1)
}

func TestMutesDeleteMuteExisting(t *testing.T) {
	t.Parallel()

	mutes := Mutes{
		Mutes: []*Mute{
			{Chain: "chain2", ProposalID: "proposal1", Expires: time.Now().Add(time.Hour)},
			{Chain: "chain1", ProposalID: "proposal2", Expires: time.Now().Add(time.Hour)},
			{Chain: "chain1", ProposalID: "proposal1", Expires: time.Now().Add(time.Hour)},
			{Chain: "chain2", ProposalID: "proposal2", Expires: time.Now().Add(time.Hour)},
		},
	}

	deleted := mutes.DeleteMute(&Mute{Chain: "chain1", ProposalID: "proposal1"})
	assert.True(t, deleted)
	assert.Len(t, mutes.Mutes, 3)
}
