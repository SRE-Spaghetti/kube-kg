## 4. Epic List

1.  **Epic 1: Foundation & Core Connectivity:** Establish the Go project structure, configuration loading, observability with OpenTelemetry, and a basic, tested client for connecting to the Neo4j database.
2.  **Epic 2: Initial Graph Synchronization:** Develop the KubeView API client, map Kubernetes resources to graph nodes and relationships, and implement the logic to perform a full, one-time synchronization of the entire cluster into Neo4j.
3.  **Epic 3: Real-Time Graph Updates:** Implement the client for the KubeView Server-Sent Events (SSE) stream and develop the event processor to handle `add`, `update`, and `delete` events, keeping the Neo4j graph continuously up-to-date.
4.  **Epic 4: Service API & Orchestration:** Build the internal REST API with `/health` and `/refresh` endpoints, and create the main service logic to orchestrate the initial sync, the real-time event processing, and graceful shutdown.
