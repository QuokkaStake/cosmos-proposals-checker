package types

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestJsonErrorMarshalJson(t *testing.T) {
	t.Parallel()

	jsonErr := JSONError{error: "error"}
	value, err := jsonErr.MarshalJSON()

	assert.NoError(t, err)
	assert.Equal(t, []byte("\"error\""), value)
}

func TestJsonErrorUnmarshalJson(t *testing.T) {
	t.Parallel()

	jsonErr := JSONError{}
	err := jsonErr.UnmarshalJSON([]byte("error"))

	assert.NoError(t, err)
	assert.Equal(t, "error", jsonErr.error)
}

func TestJsonErrorToString(t *testing.T) {
	t.Parallel()

	jsonErr := JSONError{error: "error"}
	assert.Equal(t, "error", jsonErr.Error())
}
