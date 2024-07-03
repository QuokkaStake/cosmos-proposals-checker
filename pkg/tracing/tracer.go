package tracing

import (
	"context"
	"encoding/base64"
	"main/pkg/types"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

func getExporter(config types.TracingConfig) tracesdk.SpanExporter {
	if config.Enabled.Bool {
		opts := []otlptracehttp.Option{
			otlptracehttp.WithEndpoint(config.OpenTelemetryHTTPHost),
		}

		if config.OpenTelemetryHTTPInsecure.Bool {
			opts = append(opts, otlptracehttp.WithInsecure())
		}

		if config.OpenTelemetryHTTPUser != "" && config.OpenTelemetryHTTPPassword != "" {
			auth := config.OpenTelemetryHTTPUser + ":" + config.OpenTelemetryHTTPPassword
			token := base64.StdEncoding.EncodeToString([]byte(auth))
			opts = append(opts, otlptracehttp.WithHeaders(map[string]string{
				"Authorization": "Basic " + token,
			}))
		}

		exporter, _ := otlptracehttp.New(
			context.Background(),
			opts...,
		)
		return exporter
	}

	return NewNoopExporter()
}

func InitTracer(config types.TracingConfig, version string) trace.Tracer {
	exporter := getExporter(config)
	tp := NewTraceProvider(exporter, version)
	otel.SetTracerProvider(tp)

	return tp.Tracer("main")
}

func InitNoopTracer() trace.Tracer {
	tp := NewTraceProvider(NewNoopExporter(), "1.2.3")
	otel.SetTracerProvider(tp)

	return tp.Tracer("main")
}
