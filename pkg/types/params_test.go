package types

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestParamsDescription(t *testing.T) {
	t.Parallel()

	assert.Equal(
		t,
		"Description",
		PercentParam{Description: "Description"}.GetDescription(),
		"Wrong value!",
	)

	assert.Equal(
		t,
		"Description",
		AmountsParam{Description: "Description"}.GetDescription(),
		"Wrong value!",
	)

	assert.Equal(
		t,
		"Description",
		DurationParam{Description: "Description"}.GetDescription(),
		"Wrong value!",
	)
}

func TestPercentParamSerialize(t *testing.T) {
	t.Parallel()

	params := PercentParam{
		Value: 0.4,
	}

	assert.Equal(t, "40.00%", params.Serialize(), "Wrong value!")
}

func TestAmountParamSerialize(t *testing.T) {
	t.Parallel()

	params := AmountsParam{
		Value: []Amount{
			{Denom: "stake", Amount: "100"},
			{Denom: "test", Amount: "100"},
		},
	}

	assert.Equal(t, "100 stake,100 test", params.Serialize(), "Wrong value!")
}

func TestDurationParamSerialize(t *testing.T) {
	t.Parallel()

	params := DurationParam{
		Value: time.Hour + time.Minute,
	}

	assert.Equal(t, "1 hour 1 minute", params.Serialize(), "Wrong value!")
}
