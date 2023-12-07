package types

import (
	"errors"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestQueryErrorSerializeWithQueryError(t *testing.T) {
	t.Parallel()

	queryError := QueryError{
		QueryError: errors.New("test error"),
	}

	serializedError := queryError.Error()
	assert.Equal(t, "test error", serializedError, "Error mismatch!")
}

func TestQueryErrorSerializeWithoutQueryError(t *testing.T) {
	t.Parallel()

	queryError := QueryError{
		NodeErrors: []NodeError{
			{Node: "test", Error: NewJSONError(errors.New("test error"))},
			{Node: "test2", Error: NewJSONError(errors.New("test error2"))},
		},
	}

	serializedError := queryError.Error()
	assert.Equal(
		t,
		"All LCD requests failed:\n#1: test -> test error\n#2: test2 -> test error2\n",
		serializedError,
		"Error mismatch!",
	)
}
