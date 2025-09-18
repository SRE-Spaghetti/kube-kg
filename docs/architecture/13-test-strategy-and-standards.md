## 13. Test Strategy and Standards

### 13.1. Testing Philosophy

-   **Approach:** Test-After. Tests will be written after the implementation of a feature.
-   **Coverage Goals:** A target of 80% unit test coverage will be enforced.
-   **Test Pyramid:** The testing strategy will follow the testing pyramid: a large number of unit tests, a smaller number of integration tests, and a minimal set of end-to-end tests.

### 13.2. Test Types and Organization

#### Unit Tests

-   **Framework:** Go `testing` package
-   **File Convention:** `_test.go`
-   **Location:** In the same package as the code being tested.
-   **Mocking Library:** `gomock`
-   **Coverage Requirement:** 80%

**AI Agent Requirements:**
-   Generate tests for all public methods.
-   Cover edge cases and error conditions.
-   Follow AAA pattern (Arrange, Act, Assert).
-   Mock all external dependencies.

#### Integration Tests

-   **Scope:** Test the interaction between the service and its external dependencies (KubeView, Neo4j).
-   **Location:** In a separate package with a `_test` suffix (e.g., `processor_test`).
-   **Test Infrastructure:**
    -   **KubeView:** A mock KubeView server will be used to serve static JSON and SSE data.
    -   **Neo4j:** Testcontainers will be used to spin up a real Neo4j database in a Docker container for each test run.

#### End-to-End Tests

-   **Framework:** Manual test plan documented in the `README.md`.
-   **Scope:** Test the entire service from end to end, from the KubeView API to the Neo4j database.
-   **Environment:** A dedicated `staging` environment.
-   **Test Data:** A pre-defined set of Kubernetes manifests will be applied to the `staging` cluster to create a known state for testing.

### 13.3. Test Data Management

-   **Strategy:** Test data will be stored in a `/testdata` directory.
-   **Fixtures:** JSON files will be used as fixtures for mock KubeView API responses.
-   **Factories:** Not required for this project.
-   **Cleanup:** The Testcontainers library will automatically clean up the Neo4j database after each integration test run.

### 13.4. Continuous Testing

-   **CI Integration:** All unit and integration tests will be run on every push to the repository using GitHub Actions.
-   **Performance Tests:** Not in scope for the initial version.
-   **Security Tests:** Not in scope for the initial version.
