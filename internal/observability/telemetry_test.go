package observability

import (
	"context"
	"testing"
)

func TestGetTracer(t *testing.T) {
	// Arrange
	_, err := InitTracerProvider(context.Background(), "localhost:4317")
	if err != nil {
		t.Fatalf("failed to initialize tracer provider: %v", err)
	}

	// Act
	tracer := GetTracer()

	// Assert
	if tracer == nil {
		t.Error("expected a tracer instance, got nil")
	}
}

func TestGetMeter(t *testing.T) {
	// Arrange
	_, err := InitMeterProvider(context.Background(), "localhost:4317")
	if err != nil {
		t.Fatalf("failed to initialize meter provider: %v", err)
	}

	// Act
	meter := GetMeter()

	// Assert
	if meter == nil {
		t.Error("expected a meter instance, got nil")
	}
}
