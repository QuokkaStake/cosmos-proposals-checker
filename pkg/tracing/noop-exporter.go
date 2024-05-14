package tracing

import (
	"context"

	tracesdk "go.opentelemetry.io/otel/sdk/trace"
)

type NoopExporter struct {
}

func NewNoopExporter() tracesdk.SpanExporter {
	return &NoopExporter{}
}

func (e *NoopExporter) ExportSpans(ctx context.Context, spans []tracesdk.ReadOnlySpan) error {
	return nil
}

func (e *NoopExporter) Shutdown(ctx context.Context) error {
	return nil
}
