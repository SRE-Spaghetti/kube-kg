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
