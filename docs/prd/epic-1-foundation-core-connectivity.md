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
