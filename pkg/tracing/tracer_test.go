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
	exporter := getExporter(config)
	require.NotNil(t, exporter)
}

func TestTracerGetExporterHttpBasic(t *testing.T) {
	t.Parallel()

	config := types.TracingConfig{Enabled: null.BoolFrom(true)}
	exporter := getExporter(config)
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
	exporter := getExporter(config)
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
	tracer := InitTracer(config, "v1.2.3")
	require.NotNil(t, tracer)
}

func TestTracerGetNoopTracerOk(t *testing.T) {
	t.Parallel()

	tracer := InitNoopTracer()
	require.NotNil(t, tracer)
}
