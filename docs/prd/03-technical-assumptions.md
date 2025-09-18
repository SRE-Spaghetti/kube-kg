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
