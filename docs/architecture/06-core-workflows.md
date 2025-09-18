## 6. Core Workflows

### 6.1. Initial Synchronization Workflow

This diagram illustrates the process of performing a full synchronization of the Kubernetes cluster state when the service starts or when a manual refresh is triggered.

```mermaid
sequenceDiagram
    participant main as MainApp
    participant proc as Processor
    participant kv as KubeViewClient
    participant n4j as Neo4jClient
    participant graph as GraphMapper

    main->>proc: StartInitialSync()
    proc->>kv: ListNamespaces()
    kv-->>proc: namespaces
    loop for each namespace
        proc->>kv: FetchNamespaceResources(ns)
        kv-->>proc: resources
        proc->>n4j: BeginTransaction()
        loop for each resource
            proc->>graph: KubernetesResourceToNode(res)
            graph-->>proc: node
            proc->>n4j: MERGE (node)
            proc->>graph: ExtractRelationships(res)
            graph-->>proc: relationships
            loop for each relationship
                proc->>n4j: MERGE (relationship)
            end
        end
        proc->>n4j: CommitTransaction()
    end
```

### 6.2. Real-Time Event Processing Workflow

This diagram illustrates how the service processes a single real-time event from the KubeView SSE stream.

```mermaid
sequenceDiagram
    participant sse as KubeViewSSE
    participant kv as KubeViewClient
    participant proc as Processor
    participant n4j as Neo4jClient
    participant graph as GraphMapper

    sse-->>kv: SSE Event (add/update/delete)
    kv-->>proc: event
    alt event is 'add' or 'update'
        proc->>graph: KubernetesResourceToNode(event.resource)
        graph-->>proc: node
        proc->>n4j: MERGE (node)
        proc->>graph: ExtractRelationships(event.resource)
        graph-->>proc: relationships
        proc->>n4j: MERGE (relationships)
    else event is 'delete'
        proc->>n4j: DETACH DELETE (node where uid=event.resource.uid)
    end
```
