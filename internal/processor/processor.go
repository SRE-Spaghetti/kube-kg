package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"

	"kube-kg/internal/graph"
	"kube-kg/internal/kubeview"
	"kube-kg/internal/neo4j"

	"go.opentelemetry.io/otel"
)

// Processor handles the synchronization of Kubernetes data to Neo4j.

type Processor struct {
	kubeClient  *kubeview.Client
	neo4jClient *neo4j.Client
}

// NewProcessor creates a new Processor.
func NewProcessor(kubeClient *kubeview.Client, neo4jClient *neo4j.Client) *Processor {
	return &Processor{
		kubeClient:  kubeClient,
		neo4jClient: neo4jClient,
	}
}

// InitialSync performs an initial synchronization of the Kubernetes cluster state to Neo4j.
func (p *Processor) InitialSync(ctx context.Context) error {
	ctx, span := otel.Tracer("kube-kg/internal/processor").Start(ctx, "InitialSync")
	defer span.End()

	namespaceResult, err := p.kubeClient.ListNamespaces(ctx)
	if err != nil {
		return fmt.Errorf("failed to list namespaces: %w", err)
	}

	for _, namespace := range namespaceResult.Namespaces {
		span.AddEvent(fmt.Sprintf("processing namespace: %s", namespace))
		rawResources, err := p.kubeClient.FetchNamespaceResources(ctx, namespace, "initial-sync")
		if err != nil {
			return fmt.Errorf("failed to fetch resources for namespace %s: %w", namespace, err)
		}

		var resources []kubeview.KubernetesResource
		for _, rawResourceArray := range rawResources {
			var resourceSlice []kubeview.KubernetesResource
			if err := json.Unmarshal(rawResourceArray, &resourceSlice); err != nil {
				return fmt.Errorf("failed to unmarshal resource slice: %w", err)
			}
			resources = append(resources, resourceSlice...)
		}

		tx, err := p.neo4jClient.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction: %w", err)
		}
		defer func() {
			if err := tx.Close(ctx); err != nil {
				slog.Error("failed to close transaction", "err", err)
			}
		}()

		for _, resource := range resources {
			node := graph.KubernetesResourceToNode(resource)
							if err := p.neo4jClient.MergeNode(ctx, tx, node); err != nil {
								if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
									slog.Error("failed to rollback transaction", "err", rollbackErr)
								}
								return fmt.Errorf("failed to merge node: %w", err)
							}
			relationships := graph.ExtractRelationships(resource, resources)
			for _, rel := range relationships {
									if err := p.neo4jClient.MergeRelationship(ctx, tx, rel); err != nil {
										if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
											slog.Error("failed to rollback transaction", "err", rollbackErr)
										}
										return fmt.Errorf("failed to merge relationship: %w", err)
									}			}
		}

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed to commit transaction: %w", err)
		}
	}

	return nil
}

// StartEventProcessor starts a goroutine to process events from the KubeView SSE stream.
func StartEventProcessor(ctx context.Context, eventChan <-chan kubeview.Event, neo4jClient *neo4j.Client) {
	go func() {
		for {
			select {
			case <-ctx.Done():
				slog.Info("stopping event processor")
				return
			case event := <-eventChan:
				processEvent(ctx, event, neo4jClient)
			}
		}
	}()
}

func processEvent(ctx context.Context, event kubeview.Event, neo4jClient *neo4j.Client) {
	ctx, span := otel.Tracer("kube-kg/internal/processor").Start(ctx, "processEvent")
	defer span.End()

	switch event.Type {
	case "add", "update":
		handleAddOrUpdate(ctx, event, neo4jClient)
	case "delete":
		handleDelete(ctx, event, neo4jClient)
	default:
		slog.Warn("unknown event type", "type", event.Type)
	}
}

func handleAddOrUpdate(ctx context.Context, event kubeview.Event, neo4jClient *neo4j.Client) {
	node := graph.KubernetesResourceToNode(event.Object)
	relationships := graph.ExtractRelationships(event.Object, nil) // In a real-time scenario, we might need to fetch related resources

	tx, err := neo4jClient.Begin(ctx)
	if err != nil {
		slog.Error("failed to begin transaction", "err", err)
		return
	}
	defer func() {
		if err := tx.Close(ctx); err != nil {
			slog.Error("failed to close transaction", "err", err)
		}
	}()

	if err := neo4jClient.MergeNode(ctx, tx, node); err != nil {
		slog.Error("failed to merge node", "err", err)
		if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
			slog.Error("failed to rollback transaction", "err", rollbackErr)
		}
		return
	}

	for _, rel := range relationships {
		if err := neo4jClient.MergeRelationship(ctx, tx, rel); err != nil {
			slog.Error("failed to merge relationship", "err", err)
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				slog.Error("failed to rollback transaction", "err", rollbackErr)
			}
			return
		}
	}

	if err := tx.Commit(ctx); err != nil {
		slog.Error("failed to commit transaction", "err", err)
	}
}

func handleDelete(ctx context.Context, event kubeview.Event, neo4jClient *neo4j.Client) {
	if err := neo4jClient.DeleteNode(ctx, event.Object.Metadata.UID); err != nil {
		slog.Error("failed to delete node", "err", err)
	}
}
