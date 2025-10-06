package api

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"kube-kg/internal/kubeview"
)

// MockKubeviewClient is a mock implementation of the KubeviewClient interface.
type MockKubeviewClient struct {
	mock.Mock
}

func (m *MockKubeviewClient) ListNamespaces(ctx context.Context) (*kubeview.NamespaceListResult, error) {
	args := m.Called(ctx)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*kubeview.NamespaceListResult), args.Error(1)
}

// MockNeo4jClient is a mock implementation of the Neo4jClient interface.
type MockNeo4jClient struct {
	mock.Mock
}

func (m *MockNeo4jClient) VerifyConnectivity(ctx context.Context) error {
	args := m.Called(ctx)
	return args.Error(0)
}

// MockProcessor is a mock implementation of the Processor interface.
type MockProcessor struct {
	mock.Mock
	initialSyncCalled chan bool
}

func (m *MockProcessor) InitialSync(ctx context.Context) error {
	args := m.Called(ctx)
	m.initialSyncCalled <- true
	return args.Error(0)
}

func TestHealthHandler(t *testing.T) {
	t.Run("should return 200 OK when both services are healthy", func(t *testing.T) {
		kc := new(MockKubeviewClient)
		nc := new(MockNeo4jClient)
		p := new(MockProcessor)

		kc.On("ListNamespaces", mock.Anything).Return(&kubeview.NamespaceListResult{}, nil)
		nc.On("VerifyConnectivity", mock.Anything).Return(nil)

		server := NewServer(kc, nc, p)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rr := httptest.NewRecorder()

		server.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusOK, rr.Code)
		assert.JSONEq(t, `{"kubeview":"ok","neo4j":"ok"}`, rr.Body.String())
	})

	t.Run("should return 503 Service Unavailable when kubeview is down", func(t *testing.T) {
		kc := new(MockKubeviewClient)
		nc := new(MockNeo4jClient)
		p := new(MockProcessor)

		kc.On("ListNamespaces", mock.Anything).Return(nil, assert.AnError)
		nc.On("VerifyConnectivity", mock.Anything).Return(nil)

		server := NewServer(kc, nc, p)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rr := httptest.NewRecorder()

		server.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
		assert.JSONEq(t, `{"kubeview":"unavailable","neo4j":"ok"}`, rr.Body.String())
	})

	t.Run("should return 503 Service Unavailable when neo4j is down", func(t *testing.T) {
		kc := new(MockKubeviewClient)
		nc := new(MockNeo4jClient)
		p := new(MockProcessor)

		kc.On("ListNamespaces", mock.Anything).Return(&kubeview.NamespaceListResult{}, nil)
		nc.On("VerifyConnectivity", mock.Anything).Return(assert.AnError)

		server := NewServer(kc, nc, p)

		req := httptest.NewRequest(http.MethodGet, "/health", nil)
		rr := httptest.NewRecorder()

		server.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusServiceUnavailable, rr.Code)
		assert.JSONEq(t, `{"kubeview":"ok","neo4j":"unavailable"}`, rr.Body.String())
	})
}

func TestRefreshHandler(t *testing.T) {
	t.Run("should return 202 Accepted and trigger initial sync", func(t *testing.T) {
		kc := new(MockKubeviewClient)
		nc := new(MockNeo4jClient)
		p := &MockProcessor{
			initialSyncCalled: make(chan bool, 1),
		}

		p.On("InitialSync", mock.Anything).Return(nil)

		server := NewServer(kc, nc, p)

		req := httptest.NewRequest(http.MethodPost, "/refresh", nil)
		rr := httptest.NewRecorder()

		server.ServeHTTP(rr, req)

		assert.Equal(t, http.StatusAccepted, rr.Code)

		select {
		case <-p.initialSyncCalled:
			// success
		case <-time.After(1 * time.Second):
			t.Fatal("InitialSync was not called within 1 second")
		}
	})
}
