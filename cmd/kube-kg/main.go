package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"kube-kg/internal/api"
	"kube-kg/internal/config"
	"kube-kg/internal/kubeview"
	"kube-kg/internal/neo4j"
	"kube-kg/internal/observability"
	"kube-kg/internal/processor"
)

func main() {
	// Setup structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	// Create a main context that can be cancelled to signal shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Load configuration
	cfg := config.LoadConfig()
	slog.Info("Configuration loaded successfully")

	// Initialize OpenTelemetry
	tracerProvider, err := observability.InitTracerProvider(ctx, cfg.OtelExporterEndpoint)
	if err != nil {
		slog.Error("failed to initialize tracer provider", "error", err)
		os.Exit(1)
	}
	defer func() {
		// Use a background context for shutdown to ensure it completes
		if err := tracerProvider.Shutdown(context.Background()); err != nil {
			slog.Error("failed to shutdown tracer provider", "error", err)
		}
	}()

	// Initialize clients
	neo4jClient, err := neo4j.NewClient(ctx, cfg)
	if err != nil {
		slog.Error("failed to initialize Neo4j client", "error", err, "url", cfg.Neo4jURI)
		os.Exit(1)
	}

	if cfg.KubeviewURL == "" {
		slog.Error("failed to initialize Kubeview client", "url", cfg.KubeviewURL)
		os.Exit(2)
	}
	kubeviewClient := kubeview.NewClient(cfg.KubeviewURL)

	// Initialize the processor
	proc := processor.NewProcessor(kubeviewClient, neo4jClient)

	// Start initial synchronization in a background goroutine
	go func() {
		slog.Info("Starting initial cluster synchronization")
		if err := proc.InitialSync(ctx); err != nil {
			slog.Error("initial sync failed", "error", err)
		} else {
			slog.Info("Initial cluster synchronization completed successfully")
		}
	}()

	// Create channel for KubeView events and start the event processor
	eventChan := make(chan kubeview.Event)
	kubeviewClient.StreamUpdates(ctx, cfg.ClientID, eventChan)
	processor.StartEventProcessor(ctx, eventChan, neo4jClient)
	slog.Info("Started real-time event processor")

	// Setup and start HTTP server
	server := api.NewServer(kubeviewClient, neo4jClient, proc)
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: server,
	}

	go func() {
		slog.Info("Starting HTTP server", "address", httpServer.Addr)
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP server failed", "error", err)
			os.Exit(1)
		}
	}()

	// Wait for shutdown signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// Initiate graceful shutdown
	slog.Info("Shutting down server...")

	// Create a context with a timeout for shutdown
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	// Shutdown HTTP server
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		slog.Error("HTTP server shutdown failed", "error", err)
	}

	// Close Neo4j client connection
	if err := neo4jClient.Close(shutdownCtx); err != nil {
		slog.Error("Failed to close Neo4j client", "error", err)
	}

	// Cancel the main context to signal background processes to stop
	cancel()

	slog.Info("Server gracefully stopped")
}
