package config

import (
	"os"
	"testing"
)

func TestLoadConfig(t *testing.T) {
	// Arrange
	os.Setenv("KUBEVIEW_URL", "http://localhost:8080")
	os.Setenv("NEO4J_URI", "neo4j://localhost:7687")
	os.Setenv("NEO4J_USER", "neo4j")
	os.Setenv("NEO4J_PASSWORD", "password")
	os.Setenv("CLIENT_ID", "test-client")
	os.Setenv("OTEL_EXPORTER_OTLP_ENDPOINT", "otel.example.com:4317")

	// Act
	cfg := LoadConfig()

	// Assert
	if cfg.KubeviewURL != "http://localhost:8080" {
		t.Errorf("expected KubeviewURL to be 'http://localhost:8080', got '%s'", cfg.KubeviewURL)
	}
	if cfg.Neo4jURI != "neo4j://localhost:7687" {
		t.Errorf("expected Neo4jURI to be 'neo4j://localhost:7687', got '%s'", cfg.Neo4jURI)
	}
	if cfg.Neo4jUser != "neo4j" {
		t.Errorf("expected Neo4jUser to be 'neo4j', got '%s'", cfg.Neo4jUser)
	}
	if cfg.Neo4jPassword != "password" {
		t.Errorf("expected Neo4jPassword to be 'password', got '%s'", cfg.Neo4jPassword)
	}
	if cfg.ClientID != "test-client" {
		t.Errorf("expected ClientID to be 'test-client', got '%s'", cfg.ClientID)
	}
	if cfg.OtelExporterEndpoint != "otel.example.com:4317" {
		t.Errorf("expected OtelExporterEndpoint to be 'otel.example.com:4317', got '%s'", cfg.OtelExporterEndpoint)
	}
}
