package kubeview

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestClient_ListNamespaces(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/namespaces", r.URL.Path)
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"namespaces": ["default", "kube-system"],
			"clusterHost": "https://10.96.0.1:443",
			"version": "2.1.1",
			"buildInfo": "stable 89f388f 2025-07-14",
			"mode": "in-cluster",
			"podLogsEnabled": true
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	result, err := client.ListNamespaces(context.Background())

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.Namespaces, 2)
	assert.Equal(t, "default", result.Namespaces[0])
}

func TestClient_FetchNamespaceResources(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/api/fetch/default", r.URL.Path)
		assert.Equal(t, "test-client", r.URL.Query().Get("clientID"))
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{
			"pods": [
				{
					"apiVersion": "v1",
					"kind": "Pod",
					"metadata": {
						"name": "test-pod",
						"namespace": "default"
					}
				}
			]
		}`))
	}))
	defer server.Close()

	client := NewClient(server.URL)
	result, err := client.FetchNamespaceResources(context.Background(), "default", "test-client")

	require.NoError(t, err)
	assert.NotNil(t, result)
	assert.Contains(t, result, "pods")

	var pods []KubernetesResource
	err = json.Unmarshal((result["pods"]), &pods)
	require.NoError(t, err)

	assert.Len(t, pods, 1)
	assert.Equal(t, "Pod", pods[0].Kind)
	assert.Equal(t, "test-pod", pods[0].Metadata.Name)
}
