package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestProposalGetVotingTime(t *testing.T) {
	t.Parallel()

	proposal := Proposal{
		EndTime: time.Now().Add(time.Hour).Add(time.Minute),
	}

	assert.Equal(t, "1 hour 1 minute", proposal.GetTimeLeft(), "Wrong value!")
}
