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
