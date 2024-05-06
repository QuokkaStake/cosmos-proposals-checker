package logger_test

import (
	loggerPkg "main/pkg/logger"
	"main/pkg/types"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGetDefaultLogger(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetDefaultLogger()
	require.NotNil(t, logger)
}

func TestGetLoggerInvalidLogLevel(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			require.Fail(t, "Expected to have a panic here!")
		}
	}()

	loggerPkg.GetLogger(types.LogConfig{LogLevel: "invalid"})
}

func TestGetLoggerValidPlain(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetLogger(types.LogConfig{LogLevel: "info"})
	require.NotNil(t, logger)
}

func TestGetLoggerValidJSON(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetLogger(types.LogConfig{LogLevel: "info", JSONOutput: true})
	require.NotNil(t, logger)
}

func TestGetLoggerNop(t *testing.T) {
	t.Parallel()

	logger := loggerPkg.GetNopLogger()
	require.NotNil(t, logger)
}
