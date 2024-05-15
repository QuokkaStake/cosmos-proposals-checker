package tracing

import (
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.24.0"
)

func NewTraceProvider(exp tracesdk.SpanExporter, version string) *tracesdk.TracerProvider {
	// Ensure default SDK resources and the required service name are set.
	r, _ := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("cosmos-proposals-checker"),
			semconv.ServiceVersionKey.String(version),
		),
	)

	return tracesdk.NewTracerProvider(
		tracesdk.WithBatcher(exp),
		tracesdk.WithResource(r),
	)
}
