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

func TestProposalInVoting(t *testing.T) {
	t.Parallel()

	assert.True(t, Proposal{Status: ProposalStatusVoting}.IsInVoting())
	assert.False(t, Proposal{Status: ProposalStatusPassed}.IsInVoting())
}

func TestProposalStatusSerialize(t *testing.T) {
	t.Parallel()

	assert.Equal(t, "🗳️Voting", ProposalStatusVoting.String())
	assert.Equal(t, "💸Deposit", ProposalStatusDeposit.String())
	assert.Equal(t, "🙌 Passed", ProposalStatusPassed.String())
	assert.Equal(t, "🙅‍Rejected", ProposalStatusRejected.String())
	assert.Equal(t, "🤦‍Failed", ProposalStatusFailed.String())
	assert.Equal(t, "test", ProposalStatus("test").String())
}
