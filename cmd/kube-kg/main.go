package main

import (
	"context"
	"log/slog"
	"os"

	"log"
	"net/http"

	"kube-kg/internal/api"
	"kube-kg/internal/config"
	"kube-kg/internal/kubeview"
	"kube-kg/internal/neo4j"
	"kube-kg/internal/observability"
	"kube-kg/internal/processor"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	ctx := context.Background()
	cfg := config.LoadConfig()

	// Initialize OpenTelemetry
	tracerProvider, err := observability.InitTracerProvider(ctx, cfg.OtelExporterEndpoint)
	if err != nil {
		logger.Error("failed to initialize tracer provider", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		if err := tracerProvider.Shutdown(ctx); err != nil {
			logger.Error("failed to shutdown tracer provider", slog.Any("error", err))
		}
	}()

	_, err = observability.InitMeterProvider(ctx, cfg.OtelExporterEndpoint)
	if err != nil {
		logger.Error("failed to initialize meter provider", slog.Any("error", err))
		os.Exit(1)
	}

	// Initialize Neo4j client
	neo4jClient, err := neo4j.NewClient(ctx, cfg)
	if err != nil {
		logger.Error("failed to initialize Neo4j client", slog.Any("error", err))
		os.Exit(1)
	}
	defer func() {
		if err := neo4jClient.Close(ctx); err != nil {
			logger.Error("failed to close Neo4j client", slog.Any("error", err))
		}
	}()

	kubeviewClient := kubeview.NewClient(cfg.KubeviewURL)

	processor := processor.NewProcessor(kubeviewClient, neo4jClient)

	logger.Info("Configuration loaded",
		slog.String("KubeviewURL", cfg.KubeviewURL),
		slog.String("Neo4jURI", cfg.Neo4jURI),
		slog.String("Neo4jUser", cfg.Neo4jUser),
		slog.String("Neo4jPassword", cfg.Neo4jPassword),
		slog.String("ClientID", cfg.ClientID),
		slog.String("OtelExporterEndpoint", cfg.OtelExporterEndpoint),
	)

	// Create a startup span
	tracer := observability.GetTracer()
	_, span := tracer.Start(ctx, "application.startup")
	span.End()

	logger.Info("Application started")

	// Start the HTTP server
	server := api.NewServer(kubeviewClient, neo4jClient, processor)
	logger.Info("Starting HTTP server on port 8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatalf("could not start server: %v", err)
	}
}
