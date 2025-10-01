package graph

import (
	"encoding/json"
	"kube-kg/internal/kubeview"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func loadTestResource(t *testing.T, path string) kubeview.KubernetesResource {
	data, err := os.ReadFile(path)
	require.NoError(t, err)

	var resource kubeview.KubernetesResource
	err = json.Unmarshal(data, &resource)
	require.NoError(t, err)

	return resource
}

func TestKubernetesResourceToNode(t *testing.T) {
	resource := loadTestResource(t, "testdata/pod.json")

	node := KubernetesResourceToNode(resource)

	assert.Equal(t, "Pod", node.Label)
	assert.Equal(t, "test-pod-uid", node.ID)
	assert.Equal(t, "test-pod", node.Properties["name"])
}

func TestExtractRelationships(t *testing.T) {
	pod := loadTestResource(t, "testdata/pod.json")
	replicaSet := loadTestResource(t, "testdata/replicaset.json")
	service := loadTestResource(t, "testdata/service.json")
	configMap := loadTestResource(t, "testdata/configmap.json")

	resources := []kubeview.KubernetesResource{pod, replicaSet, service, configMap}

	relationships := ExtractRelationships(pod, resources)
	assert.Len(t, relationships, 2)

	relationships = ExtractRelationships(service, resources)
	assert.Len(t, relationships, 1)
	assert.Equal(t, "SELECTS", relationships[0].Type)
	assert.Equal(t, service.Metadata.UID, relationships[0].SourceID)
	assert.Equal(t, pod.Metadata.UID, relationships[0].TargetID)
}
