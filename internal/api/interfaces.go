package api

import (
	"context"

	"kube-kg/internal/kubeview"
)

// KubeviewClient is the interface for the Kubeview client.
type KubeviewClient interface {
	ListNamespaces(ctx context.Context) (*kubeview.NamespaceListResult, error)
}

// Neo4jClient is the interface for the Neo4j client.
type Neo4jClient interface {
	VerifyConnectivity(ctx context.Context) error
}

// Processor is the interface for the processor.
type Processor interface {
	InitialSync(ctx context.Context) error
}
