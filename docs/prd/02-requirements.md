## 2. Requirements

### Functional

1.  **FR1:** The service must connect to and consume the KubeView REST API endpoints (`/api/namespaces`, `/api/fetch/{namespace}`).
2.  **FR2:** The service must connect to and consume the KubeView Server-Sent Events (SSE) endpoint (`/updates`) for real-time updates.
3.  **FR3:** The service must perform an initial full synchronization of all Kubernetes resources from KubeView upon startup.
4.  **FR4:** The service must process real-time SSE events (`add`, `update`, `delete`) to incrementally update the knowledge graph.
5.  **FR5:** The service must convert Kubernetes resources (e.g., Deployment, Pod, Service) into corresponding nodes in the Neo4j graph.
6.  **FR6:** The service must identify and create relationships (edges) between resource nodes in Neo4j based on Kubernetes concepts like owner references, label selectors, and volume mounts.
7.  **FR7:** The service must expose an internal REST API with a `/health` endpoint to report the service's operational status and its connectivity to downstream dependencies (KubeView, Neo4j).
8.  **FR8:** The service must expose an internal REST API with a `/refresh` endpoint to trigger a full re-synchronization of the graph on demand.
9.  **FR9:** The service must handle the deletion of a Kubernetes resource by removing the corresponding node and its relationships from the graph.

### Non-Functional

1.  **NFR1:** The service must be written in Go.
2.  **NFR2:** All configuration, including KubeView URL and Neo4j credentials, must be loadable from environment variables.
3.  **NFR3:** All Neo4j write operations must be idempotent to handle retries and potential duplicate events.
4.  **NFR4:** The service must include OpenTelemetry instrumentation for distributed tracing and metrics.
5.  **NFR5:** The service must handle connection drops to the KubeView SSE endpoint with an automatic retry mechanism.
6.  **NFR6:** The service must manage Neo4j writes within transactions for atomicity and performance.
7.  **NFR7:** The service must support graceful shutdown, cleanly closing client connections and stopping the HTTP server.
8.  **NFR8:** The `clientID` used for the KubeView SSE stream must be unique per service instance to avoid conflicts.
