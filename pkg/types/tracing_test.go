package types

import (
	"testing"

	"github.com/guregu/null/v5"

	"github.com/stretchr/testify/assert"
)

func TestTracingInvalid(t *testing.T) {
	t.Parallel()

	tracing := TracingConfig{Enabled: null.BoolFrom(true)}
	err := tracing.Validate()
	assert.Error(t, err)
}

func TestTracingValid(t *testing.T) {
	t.Parallel()

	tracing := TracingConfig{Enabled: null.BoolFrom(true), OpenTelemetryHTTPHost: "test"}
	err := tracing.Validate()
	assert.NoError(t, err)
}
