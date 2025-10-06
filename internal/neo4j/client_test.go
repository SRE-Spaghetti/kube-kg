package neo4j

import (
	"context"
	"kube-kg/internal/config"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go/modules/neo4j"
)

func TestNewClient(t *testing.T) {
	ctx := context.Background()

	neo4jContainer, err := neo4j.Run(ctx, "neo4j:5", neo4j.WithAdminPassword("password"))
	require.NoError(t, err)
	defer func() {
		require.NoError(t, neo4jContainer.Terminate(ctx))
	}()

	uri, err := neo4jContainer.BoltUrl(ctx)
	require.NoError(t, err)

	cfg := &config.Config{
		Neo4jURI:      uri,
		Neo4jUser:     "neo4j",
		Neo4jPassword: "password",
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

	neo4jContainer, err := neo4j.Run(ctx, "neo4j:5", neo4j.WithAdminPassword("password"))
	require.NoError(t, err)
	defer func() {
		require.NoError(t, neo4jContainer.Terminate(ctx))
	}()

	uri, err := neo4jContainer.BoltUrl(ctx)
	require.NoError(t, err)

	cfg := &config.Config{
		Neo4jURI:      uri,
		Neo4jUser:     "neo4j",
		Neo4jPassword: "password",
	}

	client, err := NewClient(ctx, cfg)
	require.NoError(t, err)
	defer func() {
		require.NoError(t, client.Close(ctx))
	}()

	tx, err := client.Begin(ctx)
	require.NoError(t, err)
	defer func() { require.NoError(t, tx.Close(ctx)) }()

	_, err = tx.Run(ctx, "CREATE (n:Test {name: 'test'})", nil)
	require.NoError(t, err)

	result, err := tx.Run(ctx, "MATCH (n:Test) RETURN n.name", nil)
	require.NoError(t, err)

	assert.True(t, result.Next(ctx))
	assert.Equal(t, "test", result.Record().Values[0])

	require.NoError(t, tx.Commit(ctx))
}
