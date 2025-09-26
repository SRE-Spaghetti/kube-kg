package neo4j

import (
	"context"
	"fmt"

	"github.com/neo4j/neo4j-go-driver/v5/neo4j"

	"kube-kg/internal/config"
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

// RunCypher executes a Cypher query against the database.
func (c *Client) RunCypher(ctx context.Context, query string, params map[string]interface{}) (result neo4j.ResultWithContext, err error) {
	session := c.driver.NewSession(ctx, neo4j.SessionConfig{})
	defer func() {
		if closeErr := session.Close(ctx); closeErr != nil && err == nil {
			err = fmt.Errorf("failed to close session: %w", closeErr)
		}
	}()

	result, err = session.Run(ctx, query, params)
	if err != nil {
		err = fmt.Errorf("failed to run Cypher query: %w", err)
		return
	}

	return result, nil
}
