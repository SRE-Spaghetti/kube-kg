package neo4j

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/neo4j"

	"kube-kg/internal/config"
)

func TestNewClient(t *testing.T) {
	ctx := context.Background()

	neo4jContainer, err := neo4j.Run(ctx, "neo4j:5")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, neo4jContainer.Terminate(ctx))
	}()

	uri, err := neo4jContainer.BoltUrl(ctx)
	require.NoError(t, err)

	cfg := &config.Config{
		Neo4jURI:      uri,
		Neo4jUser:     "neo4j",
		Neo4jPassword: "letmein",
	}

	client, err := NewClient(ctx, cfg)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, client.Close(ctx))
	}()

	assert.NotNil(t, client)
}

func TestClient_RunCypher(t *testing.T) {
	ctx := context.Background()

	neo4jContainer, err := neo4j.Run(ctx, "neo4j:5")
	require.NoError(t, err)
	defer func() {
		require.NoError(t, neo4jContainer.Terminate(ctx))
	}()

	uri, err := neo4jContainer.BoltUrl(ctx)
	require.NoError(t, err)

	cfg := &config.Config{
		Neo4jURI:      uri,
		Neo4jUser:     "neo4j",
		Neo4jPassword: "letmein",
	}

	client, err := NewClient(ctx, cfg)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, client.Close(ctx))
	}()

	result, err := client.RunCypher(ctx, "RETURN 1", nil)
	require.NoError(t, err)

	assert.NotNil(t, result)
}
