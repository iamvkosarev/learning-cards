package tracing

import (
	"context"
	"github.com/iamvkosarev/learning-cards/internal/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.34.0"
)

func SetupTracingProvider(ctx context.Context, cfg config.OTelTracing) (*trace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(
		ctx, otlptracegrpc.WithEndpoint(cfg.TraceGRPCExporterEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, err
	}
	ops := []trace.TracerProviderOption{
		trace.WithBatcher(exporter),
		trace.WithResource(
			resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceName(cfg.ServiceName),
			),
		),
	}
	if cfg.AlwaysSample {
		ops = append(ops, trace.WithSampler(trace.AlwaysSample()))
	}
	tp := trace.NewTracerProvider(
		ops...,
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(
		propagation.NewCompositeTextMapPropagator(
			propagation.TraceContext{},
			propagation.Baggage{},
		),
	)
	return tp, nil
}
