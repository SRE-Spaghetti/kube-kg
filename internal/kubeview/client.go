package kubeview

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// NamespaceListResult represents the response from the /api/namespaces endpoint.
type NamespaceListResult struct {
	Namespaces     []string `json:"namespaces"`
	ClusterHost    string   `json:"clusterHost"`
	Version        string   `json:"version"`
	BuildInfo      string   `json:"buildInfo"`
	Mode           string   `json:"mode"`
	PodLogsEnabled bool     `json:"podLogsEnabled"`
}

// NamespaceResources represents the response from the /api/fetch/{namespace} endpoint.
type NamespaceResources map[string]json.RawMessage

// KubernetesResource represents a generic Kubernetes resource.
type KubernetesResource struct {
	APIVersion string          `json:"apiVersion"`
	Kind       string          `json:"kind"`
	Metadata   ObjectMeta      `json:"metadata"`
	Spec       json.RawMessage `json:"spec"`
	Status     json.RawMessage `json:"status,omitempty"`
}

// ObjectMeta represents the metadata of a Kubernetes object.
type ObjectMeta struct {
	Name              string            `json:"name"`
	Namespace         string            `json:"namespace"`
	UID               string            `json:"uid"`
	ResourceVersion   string            `json:"resourceVersion"`
	CreationTimestamp time.Time         `json:"creationTimestamp"`
	Labels            map[string]string `json:"labels"`
	Annotations       map[string]string `json:"annotations"`
	OwnerReferences   []json.RawMessage `json:"ownerReferences"`
	Finalizers        []string          `json:"finalizers"`
	ManagedFields     []json.RawMessage `json:"managedFields"`
}

// Client is a client for the KubeView API.
type Client struct {
	httpClient *http.Client
	baseURL    string
	tracer     trace.Tracer
}

// NewClient creates a new KubeView API client.
func NewClient(baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		baseURL: baseURL,
		tracer:  otel.Tracer("kube-kg/internal/kubeview"),
	}
}

// ListNamespaces fetches the list of namespaces from the KubeView API.
func (c *Client) ListNamespaces(ctx context.Context) (*NamespaceListResult, error) {
	ctx, _ = c.tracer.Start(ctx, "ListNamespaces")
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/namespaces", c.baseURL), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			slog.Error("failed to close response body", "err", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result NamespaceListResult
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}

// FetchNamespaceResources fetches the resources for a given namespace from the KubeView API.
func (c *Client) FetchNamespaceResources(ctx context.Context, namespace string, clientID string) (NamespaceResources, error) {
	ctx, span := c.tracer.Start(ctx, "FetchNamespaceResources", trace.WithAttributes(
		attribute.String("namespace", namespace),
		attribute.String("clientID", clientID),
	))
	defer span.End()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf("%s/api/fetch/%s?clientID=%s", c.baseURL, namespace, clientID), nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}
	defer func(body io.ReadCloser) {
		if err := body.Close(); err != nil {
			slog.Error("failed to close response body", "err", err)
		}
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result NamespaceResources
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return result, nil
}
