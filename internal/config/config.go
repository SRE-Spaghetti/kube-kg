package config

import (
	"fmt"
	"os"
)

// Config holds all configuration for the service.
type Config struct {
	KubeviewURL          string
	Neo4jURI             string
	Neo4jUser            string
	Neo4jPassword        string
	ClientID             string
	OtelExporterEndpoint string
}

// LoadConfig loads configuration from environment variables.
func LoadConfig() *Config {
	otelEndpoint := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if otelEndpoint == "" {
		otelEndpoint = "localhost:4317"
	}
	clientId := os.Getenv("CLIENT_ID")
	if clientId == "" {
		clientId = fmt.Sprintf("Client-%d", os.Getpid())
	}
	return &Config{
		KubeviewURL:          os.Getenv("KUBEVIEW_URL"),
		Neo4jURI:             os.Getenv("NEO4J_URI"),
		Neo4jUser:            os.Getenv("NEO4J_USER"),
		Neo4jPassword:        os.Getenv("NEO4J_PASSWORD"),
		ClientID:             clientId,
		OtelExporterEndpoint: otelEndpoint,
	}
}
