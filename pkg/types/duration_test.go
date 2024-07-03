package types

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestDurationMarshal(t *testing.T) {
	t.Parallel()

	duration := Duration{Duration: 30 * time.Second}
	bytes, err := duration.MarshalJSON()
	require.NoError(t, err)
	require.Equal(t, []byte("\"30s\""), bytes)
}

func TestDurationInvalidInput(t *testing.T) {
	t.Parallel()

	duration := Duration{}
	err := duration.UnmarshalJSON([]byte{})
	require.Error(t, err)
}

func TestDurationNotString(t *testing.T) {
	t.Parallel()

	duration := Duration{}
	err := json.Unmarshal([]byte("3"), &duration)
	require.Error(t, err)
}

func TestDurationInvalidString(t *testing.T) {
	t.Parallel()

	duration := Duration{}
	err := json.Unmarshal([]byte("\"asd\""), &duration)
	require.Error(t, err)
}

func TestDurationValid(t *testing.T) {
	t.Parallel()

	duration := Duration{}
	err := json.Unmarshal([]byte("\"30s\""), &duration)
	require.NoError(t, err)
	require.Equal(t, 30*time.Second, duration.Duration)
}
