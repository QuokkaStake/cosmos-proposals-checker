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
	assert.Equal(t, muted, true, "Mute should match!")
}

func TestMuteMatchesWithChainSpecified(t *testing.T) {
	t.Parallel()

	mute := &Mute{Chain: "chain"}
	muted := mute.Matches("chain", "proposal")
	assert.Equal(t, muted, true, "Mute should match!")
	muted2 := mute.Matches("chain2", "proposal")
	assert.Equal(t, muted2, false, "Mute should not match!")
}

func TestMuteMatchesWithProposalSpecified(t *testing.T) {
	t.Parallel()

	mute := &Mute{ProposalID: "proposal"}
	muted := mute.Matches("chain", "proposal")
	assert.Equal(t, muted, true, "Mute should match!")
	muted2 := mute.Matches("chain", "proposal2")
	assert.Equal(t, muted2, false, "Mute should not match!")
}

func TestMuteMatchesWithAllSpecified(t *testing.T) {
	t.Parallel()

	mute := &Mute{Chain: "chain", ProposalID: "proposal"}
	muted := mute.Matches("chain", "proposal")
	assert.Equal(t, muted, true, "Mute should match!")
	muted2 := mute.Matches("chain", "proposal2")
	assert.Equal(t, muted2, false, "Mute should not match!")
	muted3 := mute.Matches("chain2", "proposal")
	assert.Equal(t, muted3, false, "Mute should not match!")
}

func TestMutesMatchesIgnoreExpired(t *testing.T) {
	t.Parallel()

	mutes := Mutes{
		Mutes: []*Mute{
			{Expires: time.Now().Add(-time.Hour)},
		},
	}
	muted := mutes.IsMuted("chain", "proposal")
	assert.Equal(t, muted, false, "Mute should not be muted!")
}

func TestMutesMatchesNotIgnoreActual(t *testing.T) {
	t.Parallel()

	mutes := Mutes{
		Mutes: []*Mute{
			{Expires: time.Now().Add(time.Hour)},
		},
	}
	muted := mutes.IsMuted("chain", "proposal")
	assert.Equal(t, muted, true, "Mute should be muted!")
}

func TestMutesAddsMute(t *testing.T) {
	t.Parallel()

	mutes := Mutes{
		Mutes: []*Mute{
			{Chain: "chain1", Expires: time.Now().Add(-time.Hour)},
		},
	}

	mutes.AddMute(&Mute{Chain: "chain2", Expires: time.Now().Add(time.Hour)})
	assert.Equal(t, len(mutes.Mutes), 1, "There should be 1 mute!")
	assert.Equal(t, mutes.Mutes[0].Chain, "chain2", "Chain name should match!")
}

func TestMutesDeletesMute(t *testing.T) {
	t.Parallel()

	mutes := Mutes{
		Mutes: []*Mute{
			{Chain: "chain1", Expires: time.Now().Add(-time.Hour)},
		},
	}

	mutes.AddMute(&Mute{Chain: "chain2", Expires: time.Now().Add(time.Hour)})
	assert.Equal(t, len(mutes.Mutes), 1, "There should be 1 mute!")
	assert.Equal(t, mutes.Mutes[0].Chain, "chain2", "Chain name should match!")
}
