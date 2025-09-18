## 3. Tech Stack

### 3.1. Cloud Infrastructure

-   **Provider:** Cloud Agnostic
-   **Key Services:** The service is designed to run in any Kubernetes environment.
-   **Deployment Regions:** Not Applicable

### 3.2. Technology Stack Table

| Category | Technology | Version | Purpose | Rationale |
| :--- | :--- |:--------| :--- | :--- |
| **Language** | Go | 1.24.6  | Primary development language | Required by PRD. Strong performance and concurrency. |
| **Database** | Neo4j | 5.x     | Graph Database | Required by PRD for storing the knowledge graph. |
| **API** | net/http | 1.24.6  | Internal REST API | Standard library, sufficient for the simple internal API. |
| **Observability** | OpenTelemetry | 1.24.6  | Tracing and Metrics | Required by PRD. Vendor-agnostic observability. |
| **DB Driver** | neo4j-go-driver | 5.15.0  | Neo4j Database Driver | Official and recommended driver for Go. |
| **SSE Client** | r3labs/sse/v2 | 2.0.2   | Server-Sent Events Client | Well-regarded library for SSE, simplifies client implementation. |
