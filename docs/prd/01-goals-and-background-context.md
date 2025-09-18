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
