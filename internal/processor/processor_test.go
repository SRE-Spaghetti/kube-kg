package processor

import (
	"context"
	"kube-kg/internal/config"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"kube-kg/internal/kubeview"
	"kube-kg/internal/neo4j"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	neo4jcontainer "github.com/testcontainers/testcontainers-go/modules/neo4j"
)

func TestInitialSync(t *testing.T) {
	ctx := context.Background()

	// Set up mock KubeView server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/namespaces":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{"namespaces":["default"]}`))
		case "/api/fetch/default":
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte(`{
				"pods": [
					{
						"apiVersion": "v1",
						"kind": "Pod",
						"metadata": {
							"name": "test-pod",
							"namespace": "default",
							"uid": "test-pod-uid"
						}
					}
				]
			}`))
		}
	}))
	defer server.Close()

	kubeClient := kubeview.NewClient(server.URL)

	// Set up Neo4j container
	neo4jContainer, err := neo4jcontainer.Run(ctx, "neo4j:5", neo4jcontainer.WithAdminPassword("password"))
	require.NoError(t, err)
	defer func() {
		if err := neo4jContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	uri, err := neo4jContainer.BoltUrl(ctx)
	require.NoError(t, err)

	cfg := &config.Config{
		Neo4jURI:      uri,
		Neo4jUser:     "neo4j",
		Neo4jPassword: "password",
	}

	neo4jClient, err := neo4j.NewClient(ctx, cfg)
	require.NoError(t, err)
	defer func() { require.NoError(t, neo4jClient.Close(ctx)) }()

	processor := NewProcessor(kubeClient, neo4jClient)

	// Run initial sync
	err = processor.InitialSync(ctx)
	require.NoError(t, err)

	// Verify data in Neo4j
	tx, err := neo4jClient.Begin(ctx)
	require.NoError(t, err)
	defer func() { require.NoError(t, tx.Close(ctx)) }()

	result, err := tx.Run(ctx, "MATCH (n:Pod) RETURN count(n) AS count", nil)
	require.NoError(t, err)

	require.True(t, result.Next(ctx))
	assert.Equal(t, int64(1), result.Record().Values[0])
}

func TestEventProcessor(t *testing.T) {
	ctx := context.Background()

	// Set up Neo4j container
	neo4jContainer, err := neo4jcontainer.Run(ctx, "neo4j:5", neo4jcontainer.WithAdminPassword("password"))
	require.NoError(t, err)
	defer func() {
		if err := neo4jContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	uri, err := neo4jContainer.BoltUrl(ctx)
	require.NoError(t, err)

	cfg := &config.Config{
		Neo4jURI:      uri,
		Neo4jUser:     "neo4j",
		Neo4jPassword: "password",
	}

	neo4jClient, err := neo4j.NewClient(ctx, cfg)
	require.NoError(t, err)
	defer func() { require.NoError(t, neo4jClient.Close(ctx)) }()

	eventChan := make(chan kubeview.Event)
	go StartEventProcessor(ctx, eventChan, neo4jClient)

	// Test ADD event
	eventChan <- kubeview.Event{
		Type: "add",
		Object: kubeview.KubernetesResource{
			APIVersion: "v1",
			Kind:       "Pod",
			Metadata: kubeview.ObjectMeta{
				Name: "test-pod",
				UID:  "test-pod-uid",
			},
		},
	}

	// Allow time for processing
	// In a real test, you might use channels or other synchronization mechanisms
	time.Sleep(2 * time.Second)

	// Verify data in Neo4j
	tx, err := neo4jClient.Begin(ctx)
	require.NoError(t, err)

	result, err := tx.Run(ctx, "MATCH (n:Pod) RETURN count(n) AS count", nil)
	require.NoError(t, err)
	require.True(t, result.Next(ctx))
	assert.Equal(t, int64(1), result.Record().Values[0])

	require.NoError(t, tx.Close(ctx))

	// Test DELETE event
	eventChan <- kubeview.Event{
		Type: "delete",
		Object: kubeview.KubernetesResource{
			Metadata: kubeview.ObjectMeta{
				UID: "test-pod-uid",
			},
		},
	}

	time.Sleep(2 * time.Second)

	tx, err = neo4jClient.Begin(ctx)
	require.NoError(t, err)

	result, err = tx.Run(ctx, "MATCH (n:Pod) RETURN count(n) AS count", nil)
	require.NoError(t, err)
	require.True(t, result.Next(ctx))
	assert.Equal(t, int64(0), result.Record().Values[0])

	require.NoError(t, tx.Close(ctx))
}
