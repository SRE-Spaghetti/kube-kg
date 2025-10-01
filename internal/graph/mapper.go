package graph

import (
	"encoding/json"
	"kube-kg/internal/kubeview"
)

// Node represents a node in the graph.
type Node struct {
	ID         string                 `json:"id"`
	Label      string                 `json:"label"`
	Properties map[string]interface{} `json:"properties"`
}

// Relationship represents a relationship between two nodes in the graph.
type Relationship struct {
	SourceID string `json:"sourceId"`
	TargetID string `json:"targetId"`
	Type     string `json:"type"`
}

// KubernetesResourceToNode converts a Kubernetes resource to a graph node.
func KubernetesResourceToNode(resource kubeview.KubernetesResource) Node {
	properties := make(map[string]interface{})
	properties["name"] = resource.Metadata.Name
	properties["namespace"] = resource.Metadata.Namespace
	properties["creationTimestamp"] = resource.Metadata.CreationTimestamp
	properties["uid"] = resource.Metadata.UID

	// Add all labels and annotations as properties
	for k, v := range resource.Metadata.Labels {
		properties["label."+k] = v
	}
	for k, v := range resource.Metadata.Annotations {
		properties["annotation."+k] = v
	}

	return Node{
		ID:         resource.Metadata.UID,
		Label:      resource.Kind,
		Properties: properties,
	}
}

// ExtractRelationships extracts relationships from a Kubernetes resource.
func ExtractRelationships(resource kubeview.KubernetesResource, resources []kubeview.KubernetesResource) []Relationship {
	var relationships []Relationship

	// Owner references
	type OwnerReference struct {
		UID string `json:"uid"`
	}
	for _, rawOwner := range resource.Metadata.OwnerReferences {
		var owner OwnerReference
		if err := json.Unmarshal(rawOwner, &owner); err == nil {
			relationships = append(relationships, Relationship{
				SourceID: resource.Metadata.UID,
				TargetID: owner.UID,
				Type:     "OWNS",
			})
		}
	}

	// Service selectors
	if resource.Kind == "Service" {
		var serviceSpec struct {
			Selector map[string]string `json:"selector"`
		}
		if err := json.Unmarshal(resource.Spec, &serviceSpec); err == nil {
			for _, other := range resources {
				if other.Kind == "Pod" {
					match := true
					for k, v := range serviceSpec.Selector {
						if other.Metadata.Labels[k] != v {
							match = false
							break
						}
					}
					if match {
						relationships = append(relationships, Relationship{
							SourceID: resource.Metadata.UID,
							TargetID: other.Metadata.UID,
							Type:     "SELECTS",
						})
					}
				}
			}
		}
	}

	// Pod volumes
	if resource.Kind == "Pod" {
		var podSpec struct {
			Volumes []struct {
				Name      string `json:"name"`
				ConfigMap struct {
					Name string `json:"name"`
				} `json:"configMap"`
				Secret struct {
					SecretName string `json:"secretName"`
				} `json:"secret"`
			} `json:"volumes"`
		}
		if err := json.Unmarshal(resource.Spec, &podSpec); err == nil {
			for _, volume := range podSpec.Volumes {
				if volume.ConfigMap.Name != "" {
					for _, other := range resources {
						if other.Kind == "ConfigMap" && other.Metadata.Name == volume.ConfigMap.Name {
							relationships = append(relationships, Relationship{
								SourceID: resource.Metadata.UID,
								TargetID: other.Metadata.UID,
								Type:     "MOUNTS",
							})
						}
					}
				}
				if volume.Secret.SecretName != "" {
					for _, other := range resources {
						if other.Kind == "Secret" && other.Metadata.Name == volume.Secret.SecretName {
							relationships = append(relationships, Relationship{
								SourceID: resource.Metadata.UID,
								TargetID: other.Metadata.UID,
								Type:     "MOUNTS",
							})
						}
					}
				}
			}
		}
	}

	return relationships
}
