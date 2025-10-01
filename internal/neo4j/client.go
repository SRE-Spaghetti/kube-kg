package neo4j

import (
	"context"
	"fmt"
	"kube-kg/internal/config"
	"kube-kg/internal/graph"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"
)

// Client wraps the Neo4j driver.
type Client struct {
	driver neo4j.DriverWithContext
}

// NewClient creates a new Neo4j client and connects to the database.
func NewClient(ctx context.Context, cfg *config.Config) (*Client, error) {
	driver, err := neo4j.NewDriverWithContext(
		cfg.Neo4jURI,
		neo4j.BasicAuth(cfg.Neo4jUser, cfg.Neo4jPassword, ""),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create Neo4j driver: %w", err)
	}

	if err := driver.VerifyConnectivity(ctx); err != nil {
		return nil, fmt.Errorf("failed to verify Neo4j connectivity: %w", err)
	}

	return &Client{driver: driver}, nil
}

// Close closes the Neo4j driver.
func (c *Client) Close(ctx context.Context) error {
	return c.driver.Close(ctx)
}

// Begin starts a new transaction.
func (c *Client) Begin(ctx context.Context) (neo4j.ExplicitTransaction, error) {
	session := c.driver.NewSession(ctx, neo4j.SessionConfig{})
	tx, err := session.BeginTransaction(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to begin transaction: %w", err)
	}
	return tx, nil
}

// MergeNode merges a node in the graph.
func (c *Client) MergeNode(ctx context.Context, tx neo4j.ExplicitTransaction, node graph.Node) error {
	query := `
	MERGE (n:%s {uid: $uid})
	SET n += $props
	`
	query = fmt.Sprintf(query, node.Label)

	params := map[string]interface{}{
		"uid":   node.ID,
		"props": node.Properties,
	}

	_, err := tx.Run(ctx, query, params)
	return err
}

// MergeRelationship merges a relationship in the graph.
func (c *Client) MergeRelationship(ctx context.Context, tx neo4j.ExplicitTransaction, rel graph.Relationship) error {
	query := `
	MATCH (source {uid: $sourceId})
	MATCH (target {uid: $targetId})
	MERGE (source)-[:%s]->(target)
	`
	query = fmt.Sprintf(query, rel.Type)

	params := map[string]interface{}{
		"sourceId": rel.SourceID,
		"targetId": rel.TargetID,
	}

	_, err := tx.Run(ctx, query, params)
	return err
}
