## 14. Security

### 14.1. Input Validation

-   **Validation Library:** Go standard library (`encoding/json` for struct unmarshalling).
-   **Validation Location:** In the `kubeview` client, immediately after receiving data from the KubeView API or SSE stream.
-   **Required Rules:**
    -   All incoming data from KubeView must be successfully unmarshalled into the defined Go structs.
    -   Any unmarshalling errors must be logged, and the invalid data must be discarded.

### 14.2. Authentication & Authorization

-   **Auth Method:** The internal API (`/health`, `/refresh`) does not require authentication as it is not exposed outside the cluster. If this changes, token-based authentication (e.g., JWT) should be implemented.
-   **Session Management:** Not applicable.
-   **Required Patterns:**
    -   The service must use credentials (e.g., API key, username/password) to connect to the KubeView and Neo4j services. These credentials must be managed as secrets.

### 14.3. Secrets Management

-   **Development:** Use a `.env` file (which must be in `.gitignore`) to load secrets into environment variables.
-   **Production:** Use Kubernetes Secrets to mount secrets as environment variables into the running container.
-   **Code Requirements:**
    -   NEVER hardcode secrets in the source code.
    -   Access secrets only via the `config` component.
    -   NEVER log secrets or any personally identifiable information (PII).

### 14.4. API Security

-   **Rate Limiting:** Not required for the internal API.
-   **CORS Policy:** Not applicable as the API is not a public web API.
-   **Security Headers:** Not required for the internal API.
-   **HTTPS Enforcement:** The internal API will run on plain HTTP. TLS termination should be handled by the Kubernetes ingress controller or a service mesh if the API is ever exposed.

### 14.5. Data Protection

-   **Encryption at Rest:** The Neo4j database must be configured to encrypt data at rest.
-   **Encryption in Transit:** The connections to KubeView and Neo4j must use TLS.
-   **PII Handling:** The service handles Kubernetes metadata, which can be sensitive. Access to the Neo4j database must be strictly controlled.
-   **Logging Restrictions:** Do not log the full body of Kubernetes resource objects. Only log metadata like the resource `uid`, `kind`, `name`, and `namespace`.

### 14.6. Dependency Security

-   **Scanning Tool:** `govulncheck` and `trivy fs`.  `trivy image` should be used to scan the docker image.
-   **Update Policy:** Dependencies will be reviewed and updated on a regular basis.
-   **Approval Process:** New dependencies must be approved by the project lead.

### 14.7. Security Testing

-   **SAST Tool:** `gosec` will be integrated into the CI/CD pipeline to perform static analysis security testing.
-   **DAST Tool:** Not in scope for this internal service.
-   **Penetration Testing:** Not in scope for this internal service.
