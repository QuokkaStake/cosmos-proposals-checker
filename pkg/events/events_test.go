package events

import (
	"main/pkg/types"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFinishedVotingEvent(t *testing.T) {
	t.Parallel()

	event := FinishedVotingEvent{}
	assert.Equal(t, "finished_voting", event.Name())
	assert.False(t, event.IsAlert())
}

func TestProposalsQueryErrorEvent(t *testing.T) {
	t.Parallel()

	event := ProposalsQueryErrorEvent{}
	assert.Equal(t, "proposals_query_error", event.Name())
	assert.False(t, event.IsAlert())
}

func TestGenericErrorEvent(t *testing.T) {
	t.Parallel()

	event := GenericError{}
	assert.Equal(t, "generic_error", event.Name())
	assert.False(t, event.IsAlert())
}

func TestVoteQueryErrorEvent(t *testing.T) {
	t.Parallel()

	event := VoteQueryError{}
	assert.Equal(t, "vote_query_error", event.Name())
	assert.False(t, event.IsAlert())
}

func TestNotVotedEvent(t *testing.T) {
	t.Parallel()

	event := NotVotedEvent{
		Chain:    &types.Chain{Name: "chain"},
		Proposal: types.Proposal{ID: "proposal"},
		Wallet:   &types.Wallet{Address: "wallet"},
	}
	assert.Equal(t, "not_voted", event.Name())
	assert.True(t, event.IsAlert())
	assert.Equal(t, "chain", event.GetChain().Name)
	assert.Equal(t, "proposal", event.GetProposal().ID)
	assert.Equal(t, "wallet", event.GetWallet().Address)
}
func TestVotedEvent(t *testing.T) {
	t.Parallel()

	event := VotedEvent{
		Chain:    &types.Chain{Name: "chain"},
		Proposal: types.Proposal{ID: "proposal"},
		Wallet:   &types.Wallet{Address: "wallet"},
	}
	assert.Equal(t, "voted", event.Name())
	assert.True(t, event.IsAlert())
	assert.Equal(t, "chain", event.GetChain().Name)
	assert.Equal(t, "proposal", event.GetProposal().ID)
	assert.Equal(t, "wallet", event.GetWallet().Address)
}

func TestRevotedEvent(t *testing.T) {
	t.Parallel()

	event := RevotedEvent{
		Chain:    &types.Chain{Name: "chain"},
		Proposal: types.Proposal{ID: "proposal"},
		Wallet:   &types.Wallet{Address: "wallet"},
	}
	assert.Equal(t, "revoted", event.Name())
	assert.True(t, event.IsAlert())
	assert.Equal(t, "chain", event.GetChain().Name)
	assert.Equal(t, "proposal", event.GetProposal().ID)
	assert.Equal(t, "wallet", event.GetWallet().Address)
}
