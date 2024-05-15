package tracing

import (
	"main/pkg/types"
	"testing"

	"github.com/guregu/null/v5"
	"github.com/stretchr/testify/require"
)

func TestTracerGetExporterNoop(t *testing.T) {
	t.Parallel()

	config := types.TracingConfig{}
	exporter, err := getExporter(config)

	require.NoError(t, err)
	require.NotNil(t, exporter)
}

func TestTracerGetExporterHttpBasic(t *testing.T) {
	t.Parallel()

	config := types.TracingConfig{Enabled: null.BoolFrom(true)}
	exporter, err := getExporter(config)

	require.NoError(t, err)
	require.NotNil(t, exporter)
}

func TestTracerGetExporterHttpComplex(t *testing.T) {
	t.Parallel()

	config := types.TracingConfig{
		Enabled:                   null.BoolFrom(true),
		OpenTelemetryHTTPHost:     "test",
		OpenTelemetryHTTPUser:     "test",
		OpenTelemetryHTTPPassword: "test",
		OpenTelemetryHTTPInsecure: null.BoolFrom(true),
	}
	exporter, err := getExporter(config)

	require.NoError(t, err)
	require.NotNil(t, exporter)
}

func TestTracerGetTracerOk(t *testing.T) {
	t.Parallel()

	config := types.TracingConfig{
		Enabled:                   null.BoolFrom(true),
		OpenTelemetryHTTPHost:     "test",
		OpenTelemetryHTTPUser:     "test",
		OpenTelemetryHTTPPassword: "test",
		OpenTelemetryHTTPInsecure: null.BoolFrom(true),
	}
	tracer, err := InitTracer(config, "v1.2.3")

	require.NoError(t, err)
	require.NotNil(t, tracer)
}

func TestTracerGetNoopTracerOk(t *testing.T) {
	t.Parallel()

	tracer := InitNoopTracer()
	require.NotNil(t, tracer)
}
