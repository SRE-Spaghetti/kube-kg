# Kubernetes Knowledge Graph (Kube-KG) Builder Service - Product Brief

**Project Goal:** Develop a Go service that continuously builds and maintains a Neo44j knowledge graph of Kubernetes resources by consuming the KubeView API and its SSE updates, and exposes a minimal REST interface for management.

**Core Components:**

1.  **Project Structure:** Standard Go module layout.
2.  **Configuration:** Load Neo44j and KubeView API endpoints from environment variables or a configuration file.
3.  **KubeView Client:** A Go client to interact with the KubeView REST API (`/api/namespaces`, `/api/fetch/{namespace}`) and the Server-Sent Events (`/updates`) endpoint.
4.  **Neo44j Client:** A Go client to interact with the Neo44j database, responsible for creating/updating/deleting nodes and relationships.
5.  **Graph Builder/Processor:** Core logic to:
    *   Convert Kubernetes resources into Neo44j nodes with properties.
    *   Identify and create relationships (edges) between resources (e.g., Deployment -> ReplicaSet -> Pod, Service -> Pod, ConfigMap -> Pod, Secret -> Pod).
    *   Handle initial full synchronization.
    *   Process real-time SSE events (`add`, `update`, `delete`) to incrementally update the graph.
6.  **Internal REST API:** A small HTTP server within the Kube-KG service to expose:
    *   `/health`: A health check endpoint.
    *   `/refresh`: An endpoint to trigger a full re-synchronization of the graph.
7.  **Telemetry:** OpenTelemetry instrumentation for tracing and metrics.

**Detailed Plan by Phase:**

**Phase 1: Setup and Core Utilities**

*   **1.1 Project Initialization:**
    *   Create a new Go module: `go mod init github.com/yourorg/kube-kg-builder` (or similar).
    *   Define the basic directory structure:
        ```
        kube-kg-builder/
        ├── main.go
        ├── config/
        │   └── config.go
        ├── pkg/
        │   ├── kubeview/
        │   ├── neo4j/
        │   └── telemetry/
        └── internal/
            ├── graph/
            ├── api/
            └── processor/
        ```
*   **1.2 Configuration Management (`config/config.go`):**
    *   Define a `Config` struct to hold `KubeviewURL`, `Neo44jURI`, `Neo44jUser`, `Neo44jPassword`, `ClientID`.
    *   Implement a function to load these values from environment variables (e.g., using `os.LookupEnv` or a library like `spf13/viper` or `kelseyhightower/envconfig`).
*   **1.3 OpenTelemetry Setup (`pkg/telemetry/telemetry.go`):**
    *   Initialize OpenTelemetry SDK (TracerProvider, MeterProvider).
    *   Configure an exporter (e.g., OTLP to a local collector).
    *   Provide functions to get a `Tracer` and `Meter` instance.
*   **1.4 Neo44j Go Driver Integration (`pkg/neo4j/client.go`):**
    *   Import `github.com/neo4j/neo4j-go-driver/v4/neo4j`.
    *   Create a `Neo44jClient` struct with a `neo4j.Driver` instance.
    *   Implement `NewNeo44jClient(uri, user, password string) (*Neo44jClient, error)`.
    *   Implement `Close()`.
    *   Implement generic `RunCypher(ctx context.Context, query string, params map[string]interface{}) error` for executing Cypher queries.

**Phase 2: KubeView API Client and Initial Graph Ingestion**

*   **2.1 KubeView API Client (`pkg/kubeview/client.go`):**
    *   Define structs for `NamespaceListResult`, `NamespaceResources`, `KubernetesResource` (and its specific types like `Pod`, `Deployment`, `Service`, `ConfigMap`, `Secret`, `ReplicaSet`, `ObjectMeta`) based on `openapi.yaml`.
    *   Create a `KubeviewClient` struct with an `http.Client` and `baseURL`.
    *   Implement `NewKubeviewClient(baseURL, clientID string) (*KubeviewClient, error)`.
    *   Implement `ListNamespaces(ctx context.Context) (*NamespaceListResult, error)` to call `/api/namespaces`.
    *   Implement `FetchNamespaceResources(ctx context.Context, namespace string) (*NamespaceResources, error)` to call `/api/fetch/{namespace}`.
*   **2.2 Graph Node & Relationship Mapping (`internal/graph/mapper.go`):**
    *   `KubernetesResourceToNode(resource KubernetesResource) (nodeLabel string, properties map[string]interface{})`: Converts a K8s resource to a Neo44j node label and properties.
    *   `ExtractRelationships(resource KubernetesResource) ([]Relationship, error)`: Identifies and returns a list of potential relationships (source node UID, target node UID, relationship type, properties). This will involve:
        *   **Owner References:** Parsing `metadata.ownerReferences` to link child resources (e.g., Pods) to their owners (e.g., ReplicaSets, Deployments).
        *   **Label Selectors:** Matching `Service.spec.selector` to `Pod.metadata.labels`.
        *   **Volume References:** Linking Pods to ConfigMaps and Secrets via `spec.volumes`.
        *   **Service Account:** Linking Pods to ServiceAccounts (if applicable and exposed by KubeView).
*   **2.3 Initial Graph Synchronization (`internal/processor/sync.go`):**
    *   `InitialSync(ctx context.Context, kubeviewClient *kubeview.Client, neo44jClient *neo4j.Client)`:
        *   Calls `kubeviewClient.ListNamespaces()`.
        *   For each namespace:
            *   Calls `kubeviewClient.FetchNamespaceResources()`.
            *   Iterates through all resource types (Pods, Deployments, etc.).
            *   For each resource:
                *   Calls `KubernetesResourceToNode` and creates/updates the node in Neo44j using `MERGE`.
                *   Calls `ExtractRelationships` and creates/updates relationships in Neo44j using `MERGE`.
        *   **Transaction Management:** Wrap Neo44j writes in transactions for atomicity and performance.
        *   **Error Handling:** Implement retries and logging.

**Phase 3: Real-time Updates (SSE Processing)**

*   **3.1 KubeView SSE Client (`pkg/kubeview/sse.go`):**
    *   Implement `StreamUpdates(ctx context.Context, clientID string) (<-chan *Event, <-chan error)`:
        *   Connects to `/updates?clientID={id}`.
        *   Reads the `text/event-stream` using `bufio.Scanner`.
        *   Parses `event:` and `data:` lines.
        *   Unmarshals `data` into a `kubeview.Event` struct (which contains the `KubernetesResource`).
        *   Sends parsed events to an `Event` channel and errors to an `error` channel.
        *   Handles `ping` events (no-op, just keep connection alive).
        *   Includes retry logic for connection drops.
*   **3.2 Event Processor (`internal/processor/events.go`):**
    *   `StartEventProcessor(ctx context.Context, eventChan <-chan *kubeview.Event, neo44jClient *neo4j.Client)`:
        *   Spawns a goroutine to listen on the `eventChan`.
        *   For each event:
            *   **`add` event:**
                *   Call `KubernetesResourceToNode` and `neo44jClient.RunCypher` with `MERGE` for the node.
                *   Call `ExtractRelationships` and `neo44jClient.RunCypher` with `MERGE` for relationships.
            *   **`update` event:**
                *   Call `KubernetesResourceToNode` and `neo44jClient.RunCypher` with `MERGE` for the node (updating properties).
                *   Re-evaluate relationships: potentially delete old relationships and create new ones if the resource's references have changed. This might involve fetching existing relationships for the node and comparing.
            *   **`delete` event:**
                *   `neo44jClient.RunCypher` to `MATCH (n {uid: $uid}) DETACH DELETE n`. This will remove the node and all its relationships.
        *   **Asynchronous Handling:** Process events concurrently but ensure Neo44j operations are properly synchronized (e.g., using a worker pool or a single goroutine for all Neo44j writes to avoid deadlocks/conflicts).

**Phase 4: Internal REST API and Main Service Logic**

*   **4.1 Internal API Handlers (`internal/api/handlers.go`):**
    *   `HealthHandler(w http.ResponseWriter, r *http.Request)`:
        *   Checks connectivity to Neo44j and KubeView.
        *   Returns `200 OK` if both are reachable, `503 Service Unavailable` otherwise.
    *   `RefreshHandler(w http.ResponseWriter, r *http.Request, syncFunc func(context.Context))`:
        *   Triggers a full `InitialSync` in a new goroutine.
        *   Returns `202 Accepted` immediately.
*   **4.2 Main Service Entry Point (`main.go`):**
    *   Load configuration.
    *   Initialize OpenTelemetry.
    *   Create `Neo44jClient` and `KubeviewClient`.
    *   Start `InitialSync` in a goroutine.
    *   Start `StreamUpdates` and `StartEventProcessor` in separate goroutines.
    *   Set up the HTTP server for the internal REST API using `net/http` and register handlers.
    *   Graceful shutdown: Listen for OS signals (e.g., `SIGINT`, `SIGTERM`) to close clients and stop the HTTP server cleanly.

**Key Considerations & Potential Challenges:**

*   **Relationship Complexity:** The most challenging aspect will be accurately identifying and maintaining relationships between diverse Kubernetes resources. Start with explicit `ownerReferences` and label selectors, then expand to more implicit links.
*   **Data Model in Neo44j:** Carefully design the Neo44j schema (node labels, properties, relationship types) to reflect Kubernetes concepts clearly.
*   **Idempotency:** All Neo44j operations should be idempotent to handle retries and potential duplicate events from SSE. `MERGE` is crucial here.
*   **Error Handling and Observability:** Implement robust error handling, logging, and leverage OpenTelemetry for deep insights into the service's operation.
*   **Concurrency Control:** Ensure that initial sync and real-time event processing don't conflict when writing to Neo44j. A single writer goroutine or careful locking might be necessary.
*   **Resource Deletion:** When a resource is deleted, ensure all its associated nodes and relationships are removed from Neo44j. `DETACH DELETE` is useful for this.
*   **`clientID` Management:** The `clientID` for KubeView needs to be unique per Kube-KG instance. Consider generating a UUID.
