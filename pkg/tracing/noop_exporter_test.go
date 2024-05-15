package tracing

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/stretchr/testify/require"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

func TestNoopExporterShutdown(t *testing.T) {
	t.Parallel()

	exporter := NewNoopExporter()
	err := exporter.Shutdown(context.Background())

	require.NoError(t, err)
}

func TestNoopExporterExportSpans(t *testing.T) {
	t.Parallel()

	exporter := NewNoopExporter()
	tp := NewTraceProvider(exporter, "1.2.3")

	_, span := tp.Tracer("test").Start(context.Background(), "test")
	defer span.End()

	readOnlySpan, ok := span.(tracesdk.ReadOnlySpan)
	assert.True(t, ok)

	err := exporter.ExportSpans(context.Background(), []tracesdk.ReadOnlySpan{readOnlySpan})

	require.NoError(t, err)
}
