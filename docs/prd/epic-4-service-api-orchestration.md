### Epic 4: Service API & Orchestration

**Goal:** Build the internal REST API with `/health` and `/refresh` endpoints, and create the main service logic to orchestrate the initial sync, the real-time event processing, and graceful shutdown.

---

#### Story 4.1: Internal REST API

*As an operator,*
*I want an HTTP API to check the health of the service and trigger a manual refresh,*
*so that I can monitor the service and manage its state.*

**Acceptance Criteria:**

1.  An HTTP server is started using the standard `net/http` package.
2.  A `HealthHandler` is implemented for the `/health` endpoint.
3.  The health handler checks connectivity to both KubeView and Neo4j and returns `200 OK` only if both are reachable, otherwise it returns `503 Service Unavailable`.
4.  A `RefreshHandler` is implemented for the `/refresh` endpoint.
5.  The refresh handler triggers the `InitialSync` function (from Epic 2) in a new goroutine to avoid blocking the request.
6.  The refresh handler immediately returns a `202 Accepted` status code.
7.  Integration tests are written to verify the behavior of both the `/health` and `/refresh` endpoints.

---

#### Story 4.2: Main Service Orchestration

*As a developer,*
*I want a main entrypoint that correctly initializes and manages the lifecycle of all service components,*
*so that the service starts, runs, and stops reliably.*

**Acceptance Criteria:**

1.  The `main.go` file contains the main application entrypoint.
2.  The `main` function orchestrates the startup sequence:
    *   Load configuration (from Story 1.1).
    *   Initialize OpenTelemetry (from Story 1.2).
    *   Create the `Neo4jClient` and `KubeviewClient` (from Stories 1.3 and 2.1).
    *   Start the `InitialSync` process in a goroutine (from Story 2.3).
    *   Start the KubeView SSE stream and the event processor in separate goroutines (from Epic 3).
    *   Start the internal HTTP server (from Story 4.1).
3.  The application listens for OS signals (`SIGINT`, `SIGTERM`) to trigger a graceful shutdown.
4.  Upon receiving a shutdown signal, the application cleanly closes the KubeView and Neo4j client connections and stops the HTTP server.
5.  An end-to-end test plan is documented in the README, describing how to run the service and verify its core functionality (e.g., check for nodes in Neo4j, call the health endpoint).
