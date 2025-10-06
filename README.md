# Kube Knowledge Graph service (kube-kg)

This is a standalone service that will be part of the SRE Orchestrator and will build up the knowledge graph in Neo4j
from Kubernetes resources.

This is the initial prompt used with the BMAD Method

```
*plan I want to create a new service that will continually run and convert Kubernetes resources in to cypher statements
and insert them in to a running Neo4j database to create a knowledge graph. The service will connect to the REST service
defined in @openapi.yaml and here it will query the kubeview API for Kubernetes resource details. It should first list
the namespaces, and then for each namespace get the resources. Each resource should become a Node in Neo4j and the
resources attributes should become properties of the node. Where resources reference each other they should be linked by
an Edge in the Neo4j database. For example deployments might be linked to replicasets which might be linked to pods.
Services, Configmaps and Secrets might also be linked to pods. The service should also connect to the /updates endpoint
and listen out for any Server Sent Events that indicate changes in the Kubernetes deployment. It should hadle these
events asynchronously and add or delete Nodes on the Neo4j knowledge graph appropriate to the event data receieved.
I would like the service to be developed in the Go programming language and have the address of the Neo4j server and the
kubeview server be configurable. The service itself should expose a minimal REST interface itself that allows a refresh
of the configuration to be made (and creating, updating and deleting Nodes in Neo4j as appropriate) and to check the
health and status of the service. The Go code should be instrumented with open telemetry.
```

## End-to-End Testing

This plan describes how to manually run and verify the core functionality of the `kube-kg` service.

### 1. Prerequisites

- A running Kubernetes cluster.
- A running `kubeview` service that can connect to the cluster.
- A running Neo4j database.

### 2. Running the Service

To run the service, you must provide connection details for Kubeview and Neo4j via environment variables.

Open a terminal and execute the following command, replacing the placeholder values with your actual configuration:

```sh
KUBEVIEW_URL="http://your-kubeview-host:8000" \
NEO4J_URI="neo4j://your-neo4j-host:7687" \
NEO4J_USER="neo4j" \
NEO4J_PASSWORD="your-password" \
go run ./cmd/kube-kg
```

### 3. Verification Steps

#### a. Verify Service Startup

After running the command above, you should see log messages indicating that the service has started, including:
- `Configuration loaded successfully`
- `Starting initial cluster synchronization`
- `Started real-time event processor`
- `Starting HTTP server`

#### b. Check Health Endpoint

In a separate terminal, use `curl` to check the service's health endpoint. (Note: The health endpoint is implemented in a future story, but the server should be running).

```sh
curl http://localhost:8080/health
```

You should receive a response from the server.

#### c. Verify Data in Neo4j

1.  Wait for the `Initial cluster synchronization completed successfully` log message.
2.  Open your Neo4j Browser.
3.  Run the following Cypher query to see if nodes have been created:
    ```cypher
    MATCH (n) RETURN n LIMIT 25;
    ```
4.  You should see a graph containing nodes that represent your Kubernetes resources.

#### d. Verify Graceful Shutdown

1.  Go back to the terminal where the service is running.
2.  Press `Ctrl+C`.
3.  You should see log messages indicating a graceful shutdown, such as:
    - `Shutting down server...`
    - `Server gracefully stopped`

This completes the manual end-to-end test.
