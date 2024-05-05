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

	assert.Equal(t, "🗳️Voting", Proposal{Status: ProposalStatusVoting}.SerializeStatus())
	assert.Equal(t, "💸Deposit", Proposal{Status: ProposalStatusDeposit}.SerializeStatus())
	assert.Equal(t, "🙌 Passed", Proposal{Status: ProposalStatusPassed}.SerializeStatus())
	assert.Equal(t, "🙅‍Rejected", Proposal{Status: ProposalStatusRejected}.SerializeStatus())
	assert.Equal(t, "🤦‍Failed", Proposal{Status: ProposalStatusFailed}.SerializeStatus())
	assert.Equal(t, "test", Proposal{Status: ProposalStatus("test")}.SerializeStatus())

}
