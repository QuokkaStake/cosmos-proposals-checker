package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParamsFormatQuorum(t *testing.T) {
	t.Parallel()

	params := ChainWithVotingParams{
		Quorum: 0.4,
	}

	assert.Equal(t, "40.00%", params.FormatQuorum(), "Wrong value!")
}

func TestParamsFormatThreshold(t *testing.T) {
	t.Parallel()

	params := ChainWithVotingParams{
		Threshold: 0.4,
	}

	assert.Equal(t, "40.00%", params.FormatThreshold(), "Wrong value!")
}

func TestParamsFormatVetoThreshold(t *testing.T) {
	t.Parallel()

	params := ChainWithVotingParams{
		VetoThreshold: 0.4,
	}

	assert.Equal(t, "40.00%", params.FormatVetoThreshold(), "Wrong value!")
}

func TestParamsFormatFormatMinDepositAmount(t *testing.T) {
	t.Parallel()

	params := ChainWithVotingParams{
		MinDepositAmount: []Amount{
			{Denom: "stake", Amount: "100"},
			{Denom: "test", Amount: "100"},
		},
	}

	assert.Equal(t, "100 stake,100 test", params.FormatMinDepositAmount(), "Wrong value!")
}
