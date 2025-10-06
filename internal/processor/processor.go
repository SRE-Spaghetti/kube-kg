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

// InitialSync performs an initial synchronization of the Kubernetes cluster state to Neo4j.
func InitialSync(ctx context.Context, kubeClient *kubeview.Client, neo4jClient *neo4j.Client) error {
	ctx, span := otel.Tracer("kube-kg/internal/processor").Start(ctx, "InitialSync")
	defer span.End()

	namespaceResult, err := kubeClient.ListNamespaces(ctx)
	if err != nil {
		return fmt.Errorf("failed to list namespaces: %w", err)
	}

	for _, namespace := range namespaceResult.Namespaces {
		span.AddEvent(fmt.Sprintf("processing namespace: %s", namespace))
		rawResources, err := kubeClient.FetchNamespaceResources(ctx, namespace, "initial-sync")
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

		tx, err := neo4jClient.Begin(ctx)
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
							if err := neo4jClient.MergeNode(ctx, tx, node); err != nil {
								if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
									slog.Error("failed to rollback transaction", "err", rollbackErr)
								}
								return fmt.Errorf("failed to merge node: %w", err)
							}
			relationships := graph.ExtractRelationships(resource, resources)
			for _, rel := range relationships {
									if err := neo4jClient.MergeRelationship(ctx, tx, rel); err != nil {
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
