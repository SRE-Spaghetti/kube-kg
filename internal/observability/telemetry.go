package observability

import (
	"context"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/metric"
	sdkMetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.20.0"
	"go.opentelemetry.io/otel/trace"
)

// InitTracerProvider initializes and returns a new OpenTelemetry TracerProvider.
func InitTracerProvider(ctx context.Context, endpoint string) (*sdkTrace.TracerProvider, error) {
	exporter, err := otlptracegrpc.New(ctx, otlptracegrpc.WithInsecure(), otlptracegrpc.WithEndpoint(endpoint))
	if err != nil {
		return nil, err
	}

	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceNameKey.String("kube-kg"),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdkTrace.NewTracerProvider(
		sdkTrace.WithBatcher(exporter),
		sdkTrace.WithResource(res),
	)

	otel.SetTracerProvider(tp)

	return tp, nil
}

// InitMeterProvider initializes and returns a new OpenTelemetry MeterProvider.
func InitMeterProvider(ctx context.Context, endpoint string) (metric.MeterProvider, error) {
	exporter, err := otlpmetricgrpc.New(ctx, otlpmetricgrpc.WithInsecure(), otlpmetricgrpc.WithEndpoint(endpoint))
	if err != nil {
		return nil, err
	}

	mp := sdkMetric.NewMeterProvider(sdkMetric.WithReader(sdkMetric.NewPeriodicReader(exporter)))
	otel.SetMeterProvider(mp)

	return mp, nil
}

// GetTracer returns a new tracer instance.
func GetTracer() trace.Tracer {
	return otel.Tracer("kube-kg")
}

// GetMeter returns a new meter instance.
func GetMeter() metric.Meter {
	return otel.Meter("kube-kg")
}
