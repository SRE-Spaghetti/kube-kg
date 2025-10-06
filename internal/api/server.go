package api

import (
	"encoding/json"
	"net/http"
)

// Server is the HTTP server.
type Server struct {
	router         *http.ServeMux
	kubeviewClient KubeviewClient
	neo4jClient    Neo4jClient
	processor      Processor
}

// NewServer creates a new HTTP server.
func NewServer(kc KubeviewClient, nc Neo4jClient, p Processor) *Server {
	s := &Server{
		router:         http.NewServeMux(),
		kubeviewClient: kc,
		neo4jClient:    nc,
		processor:      p,
	}
	s.routes()
	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) routes() {
	s.router.HandleFunc("/health", s.handleHealth())
	s.router.HandleFunc("/refresh", s.handleRefresh())
}

func (s *Server) handleHealth() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		_, kubeviewHealth := s.kubeviewClient.ListNamespaces(r.Context())
		neo4jHealth := s.neo4jClient.VerifyConnectivity(r.Context())

		status := http.StatusOK
		response := make(map[string]string)

		if kubeviewHealth != nil {
			status = http.StatusServiceUnavailable
			response["kubeview"] = "unavailable"
		} else {
			response["kubeview"] = "ok"
		}

		if neo4jHealth != nil {
			status = http.StatusServiceUnavailable
			response["neo4j"] = "unavailable"
		} else {
			response["neo4j"] = "ok"
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(status)
		err := json.NewEncoder(w).Encode(response)
		if err != nil {
			return
		}
	}
}

func (s *Server) handleRefresh() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		go func() {
			if err := s.processor.InitialSync(r.Context()); err != nil {
				// Log the error, but don't block the response
				// In a real application, you'd use a structured logger
				println("error during initial sync:", err.Error())
			}
		}()
		w.WriteHeader(http.StatusAccepted)
	}
}
