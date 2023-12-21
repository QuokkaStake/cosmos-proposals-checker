package types

import (
	"cosmossdk.io/math"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestTallyGetVoted(t *testing.T) {
	t.Parallel()

	tally := Tally{
		{Option: "Yes", Voted: math.LegacyNewDec(2)},
		{Option: "No", Voted: math.LegacyNewDec(3)},
		{Option: "Abstain", Voted: math.LegacyNewDec(5)},
	}

	assert.Equal(t, "20.00%", tally.GetVoted(tally[0]), "Wrong value!")
}

func TestTallyGetQuorum(t *testing.T) {
	t.Parallel()

	tallyInfo := TallyInfo{
		Tally: Tally{
			{Option: "idk", Voted: math.LegacyNewDec(3)},
		},
		TotalVotingPower: math.LegacyNewDec(10),
	}

	assert.Equal(t, "30.00%", tallyInfo.GetQuorum(), "Wrong value!")
}

func TestTallyGetNotVoted(t *testing.T) {
	t.Parallel()

	tallyInfo := TallyInfo{
		Tally: Tally{
			{Option: "idk", Voted: math.LegacyNewDec(3)},
		},
		TotalVotingPower: math.LegacyNewDec(10),
	}

	assert.Equal(t, "70.00%", tallyInfo.GetNotVoted(), "Wrong value!")
}
