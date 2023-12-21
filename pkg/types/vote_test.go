package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVoteResolveVote(t *testing.T) {
	t.Parallel()

	vote := Vote{
		Options: []VoteOption{
			{Option: "Yes", Weight: 1},
			{Option: "No", Weight: 2},
		},
	}

	assert.Equal(t, "Yes, No", vote.ResolveVote(), "Wrong value!")
}

func TestVoteEqualsDifferentLength(t *testing.T) {
	t.Parallel()

	vote1 := &Vote{
		Options: []VoteOption{
			{Option: "Yes", Weight: 1},
			{Option: "No", Weight: 2},
		},
	}

	vote2 := &Vote{
		Options: []VoteOption{
			{Option: "Yes", Weight: 1},
		},
	}

	assert.False(t, vote1.VotesEquals(vote2), "Wrong value!")
}

func TestVoteEqualsDifferentVotesOptions(t *testing.T) {
	t.Parallel()

	vote1 := &Vote{
		Options: []VoteOption{
			{Option: "Yes", Weight: 1},
			{Option: "No", Weight: 2},
		},
	}

	vote2 := &Vote{
		Options: []VoteOption{
			{Option: "Yes", Weight: 1},
			{Option: "Abstain", Weight: 2},
		},
	}

	assert.False(t, vote1.VotesEquals(vote2), "Wrong value!")
}

func TestVoteEqualsDifferentVotesWeight(t *testing.T) {
	t.Parallel()

	vote1 := &Vote{
		Options: []VoteOption{
			{Option: "Yes", Weight: 1},
			{Option: "No", Weight: 2},
		},
	}

	vote2 := &Vote{
		Options: []VoteOption{
			{Option: "Yes", Weight: 1},
			{Option: "No", Weight: 3},
		},
	}

	assert.False(t, vote1.VotesEquals(vote2), "Wrong value!")
}

func TestVoteEqualsSame(t *testing.T) {
	t.Parallel()

	vote1 := &Vote{
		Options: []VoteOption{
			{Option: "Yes", Weight: 1},
			{Option: "No", Weight: 2},
		},
	}

	vote2 := &Vote{
		Options: []VoteOption{
			{Option: "Yes", Weight: 1},
			{Option: "No", Weight: 2},
		},
	}

	assert.True(t, vote1.VotesEquals(vote2), "Wrong value!")
}
