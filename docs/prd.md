# Kubernetes Knowledge Graph (Kube-KG) Builder Service Product Requirements Document (PRD)

## 1. Goals and Background Context

### Goals

*   Develop a Go service to build and maintain a real-time knowledge graph of Kubernetes cluster resources.
*   Consume data from the KubeView API, including initial state and continuous updates via Server-Sent Events (SSE).
*   Store and model the Kubernetes resources and their relationships in a Neo4j graph database.
*   Provide an internal REST API for service health checks and to trigger manual data refreshes.
*   Ensure the service is observable with OpenTelemetry for tracing and metrics.

### Background Context

This project aims to solve the challenge of understanding the complex and dynamic relationships between resources in a Kubernetes cluster. By creating a knowledge graph in Neo4j, we can represent the cluster's state in a way that is easily queryable and explorable. The service will act as a bridge, consuming data from the KubeView API, which provides a comprehensive view of cluster resources, and translating it into a graph model. This will enable advanced analysis, visualization, and operational insights that are difficult to achieve by inspecting Kubernetes resources directly.

### Change Log

| Date       | Version | Description                | Author |
| :--------- | :------ | :------------------------- | :----- |
| 2025-09-15 | 1.0     | Initial draft of the PRD.  | John   |

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

## 3. Technical Assumptions

### Repository Structure: Monorepo

*   The entire Go service, including all its packages and internal logic, will be contained within a single Git repository. This simplifies dependency management and the build process for a standalone service.

### Service Architecture: Standalone Service (Monolith)

*   The project will be built as a single, self-contained Go application. This is the most straightforward architecture for the defined scope, avoiding the complexity of a distributed microservices setup which is not required by the brief.

### Testing Requirements: Unit + Integration

*   The project must include both unit tests for individual functions and packages, and integration tests that verify the service's interaction with its external dependencies (KubeView and Neo4j). This ensures both logic correctness and reliable connectivity.

### Additional Technical Assumptions and Requests

*   **Language:** The service will be written in Go.
*   **Configuration:** All configuration will be managed via environment variables.
*   **Key Libraries:** The official Neo4j Go Driver (`github.com/neo4j/neo4j-go-driver`) will be used for database interaction.
*   **Observability:** OpenTelemetry will be integrated for tracing and metrics.
*   **API Client:** A custom Go client will be built to interact with the KubeView API and SSE stream.

## 4. Epic List

1.  **Epic 1: Foundation & Core Connectivity:** Establish the Go project structure, configuration loading, observability with OpenTelemetry, and a basic, tested client for connecting to the Neo4j database.
2.  **Epic 2: Initial Graph Synchronization:** Develop the KubeView API client, map Kubernetes resources to graph nodes and relationships, and implement the logic to perform a full, one-time synchronization of the entire cluster into Neo4j.
3.  **Epic 3: Real-Time Graph Updates:** Implement the client for the KubeView Server-Sent Events (SSE) stream and develop the event processor to handle `add`, `update`, and `delete` events, keeping the Neo4j graph continuously up-to-date.
4.  **Epic 4: Service API & Orchestration:** Build the internal REST API with `/health` and `/refresh` endpoints, and create the main service logic to orchestrate the initial sync, the real-time event processing, and graceful shutdown.

---
### Epic 1: Foundation & Core Connectivity

**Goal:** Establish the Go project structure, configuration loading, observability with OpenTelemetry, and a basic, tested client for connecting to the Neo4j database.

---

#### Story 1.1: Project Initialization & Configuration

*As a developer,*
*I want a standard Go module structure and a robust configuration loading mechanism,*
*so that I can easily manage dependencies and configure the service for different environments.*

**Acceptance Criteria:**

1.  A new Go module is initialized and a directory structure matching the project brief is created.
2.  A `Config` struct is defined to hold all required configuration values (`KubeviewURL`, `Neo4jURI`, `Neo4jUser`, `Neo4jPassword`, `ClientID`).
3.  A function is implemented to load these configuration values from environment variables.
4.  Unit tests exist to verify that configuration is loaded correctly from environment variables.
5.  The application can be run and it will successfully load configuration without errors.

---

#### Story 1.2: Observability Setup

*As an SRE,*
*I want the service to be instrumented with OpenTelemetry,*
*so that I can trace requests and monitor its performance in a standardized way.*

**Acceptance Criteria:**

1.  The OpenTelemetry SDK is initialized when the service starts.
2.  A TracerProvider and MeterProvider are configured.
3.  An OTLP exporter is configured to send telemetry data to a collector.
4.  Helper functions are available to easily get a `Tracer` and `Meter` instance for use throughout the application.
5.  A simple "application.startup" span is successfully created and exported when the application initializes, verifiable in a telemetry backend.

---

#### Story 1.3: Neo4j Client Implementation

*As a developer,*
*I want a dedicated Neo4j client wrapper,*
*so that I can interact with the database in a consistent and testable manner.*

**Acceptance Criteria:**

1.  A `Neo4jClient` struct is created that encapsulates the official Neo4j Go driver.
2.  A `NewNeo4jClient` function correctly initializes a connection to the database using the loaded configuration.
3.  A `Close` method is implemented to gracefully shut down the database connection.
4.  A generic `RunCypher` method is implemented that can execute a given Cypher query with parameters against the database.
5.  Integration tests are written to verify that the client can successfully connect to a Neo4j instance and execute a basic query (e.g., `RETURN 1`).
6.  The application, upon startup, successfully connects to the Neo4j database using the client.

---
### Epic 2: Initial Graph Synchronization

**Goal:** Develop the KubeView API client, map Kubernetes resources to graph nodes and relationships, and implement the logic to perform a full, one-time synchronization of the entire cluster into Neo4j.

---

#### Story 2.1: KubeView API Client

*As a developer,*
*I want a client to interact with the KubeView REST API,*
*so that I can fetch the list of namespaces and all resources within them.*

**Acceptance Criteria:**

1.  A `KubeviewClient` struct is created with an `http.Client`.
2.  Go structs are defined to accurately represent the JSON responses from `/api/namespaces` and `/api/fetch/{namespace}`, based on the `openapi.yaml`.
3.  A `ListNamespaces` method is implemented that correctly calls the `/api/namespaces` endpoint and unmarshals the response.
4.  A `FetchNamespaceResources` method is implemented that correctly calls the `/api/fetch/{namespace}` endpoint and unmarshals the response.
5.  The client is instrumented with OpenTelemetry to trace the API calls.
6.  Integration tests are written to verify the client can successfully connect to a KubeView instance and parse the responses.

---

#### Story 2.2: Kubernetes Resource to Graph Mapping Logic

*As a developer,*
*I want a set of pure functions that can transform Kubernetes resource objects into a graph representation,*
*so that the data transformation logic is decoupled and easily testable.*

**Acceptance Criteria:**

1.  A function `KubernetesResourceToNode` is created that takes a KubeView resource object and returns a generic representation of a Neo4j node (label and properties).
2.  A function `ExtractRelationships` is created that takes a KubeView resource object and returns a list of potential relationships (e.g., based on `ownerReferences`, label selectors, volume mounts).
3.  The mapping logic correctly identifies the `uid` of a resource as the primary identifier for nodes.
4.  The relationship extraction logic correctly identifies `ownerReferences` to link child resources to their parents (e.g., Pod to ReplicaSet).
5.  The relationship extraction logic correctly identifies `Service` selectors and links them to `Pods` with matching labels.
6.  Unit tests are written to verify the mapping and relationship extraction for all supported resource types (Pod, Deployment, Service, ConfigMap, Secret, ReplicaSet).

---

#### Story 2.3: Initial Synchronization Processor

*As a developer,*
*I want a processor that orchestrates the full synchronization of Kubernetes resources into Neo4j,*
*so that I can populate the graph with the initial state of the cluster.*

**Acceptance Criteria:**

1.  An `InitialSync` function is created that takes the `KubeviewClient` and `Neo4jClient` as dependencies.
2.  The function first calls the `KubeviewClient` to get all resources from all namespaces.
3.  For each resource, it uses the mapping logic from Story 2.2 to generate nodes and relationships.
4.  It then uses the `Neo4jClient` to execute idempotent `MERGE` queries to create/update the nodes and relationships in the database.
5.  The entire synchronization process is wrapped in a single OpenTelemetry trace.
6.  All Neo4j write operations for a given namespace are executed within a single transaction.
7.  An integration test verifies that running `InitialSync` against a mock KubeView API correctly populates a test Neo4j database with the expected nodes and relationships.

---
### Epic 3: Real-Time Graph Updates

**Goal:** Implement the client for the KubeView Server-Sent Events (SSE) stream and develop the event processor to handle `add`, `update`, and `delete` events, keeping the Neo4j graph continuously up-to-date.

---

#### Story 3.1: KubeView SSE Client

*As a developer,*
*I want a client that can connect to the KubeView SSE stream and parse incoming events,*
*so that I can receive real-time updates about Kubernetes resource changes.*

**Acceptance Criteria:**

1.  A function `StreamUpdates` is implemented that connects to the `/updates?clientID={id}` endpoint.
2.  The function correctly reads the `text/event-stream` format, parsing `event:` and `data:` lines.
3.  The `data` JSON is successfully unmarshaled into a `kubeview.Event` struct.
4.  Parsed events are sent to a Go channel for consumption by the event processor.
5.  `ping` events are handled correctly to maintain the connection without being sent to the event channel.
6.  The client includes a retry mechanism with backoff to automatically reconnect if the connection is lost.
7.  An integration test verifies that the client can connect to a mock SSE stream, parse events, and send them to the channel.

---

#### Story 3.2: Real-Time Event Processor

*As a developer,*
*I want a processor that consumes resource events and applies the corresponding changes to the Neo4j graph,*
*so that the graph remains a consistent, real-time reflection of the cluster state.*

**Acceptance Criteria:**

1.  A function `StartEventProcessor` is created that listens for events on the channel provided by the SSE client.
2.  For an `add` event, the processor uses the mapping logic (from Story 2.2) and the `Neo4jClient` to `MERGE` the new node and its relationships into the graph.
3.  For an `update` event, the processor uses the `Neo4jClient` to `MERGE` the node's new properties and re-evaluates its relationships, adding or removing them as needed.
4.  For a `delete` event, the processor uses the `Neo4jClient` to execute a `DETACH DELETE` query to remove the node and all its relationships.
5.  Each event processed is tracked within its own OpenTelemetry trace.
6.  Unit tests are written for the processor logic for each event type (`add`, `update`, `delete`), mocking the `Neo4jClient`.
7.  An integration test verifies that when an event is sent to the input channel, the correct change is persisted in the test Neo4j database.

---
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
