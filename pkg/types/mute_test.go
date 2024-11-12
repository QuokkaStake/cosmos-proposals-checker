package types

import (
	"testing"

	"github.com/guregu/null/v5"
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

	mute := &Mute{Chain: null.StringFrom("chain")}
	muted := mute.Matches("chain", "proposal")
	assert.True(t, muted, "Mute should match!")
	muted2 := mute.Matches("chain2", "proposal")
	assert.False(t, muted2, "Mute should not match!")
}

func TestMuteMatchesWithProposalSpecified(t *testing.T) {
	t.Parallel()

	mute := &Mute{ProposalID: null.StringFrom("proposal")}
	muted := mute.Matches("chain", "proposal")
	assert.True(t, muted, "Mute should match!")
	muted2 := mute.Matches("chain", "proposal2")
	assert.False(t, muted2, "Mute should not match!")
}

func TestMuteMatchesWithAllSpecified(t *testing.T) {
	t.Parallel()

	mute := &Mute{Chain: null.StringFrom("chain"), ProposalID: null.StringFrom("proposal")}
	muted := mute.Matches("chain", "proposal")
	assert.True(t, muted, "Mute should match!")
	muted2 := mute.Matches("chain", "proposal2")
	assert.False(t, muted2, "Mute should not match!")
	muted3 := mute.Matches("chain2", "proposal")
	assert.False(t, muted3, "Mute should not match!")
}
