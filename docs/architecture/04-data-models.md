## 4. Data Models

### 4.1. KubernetesResource Node

**Purpose:** To represent any Kubernetes resource within the graph.

**Key Attributes:**
-   `uid`: `string` - The unique ID of the Kubernetes resource. This will be the primary key for the node.
-   `kind`: `string` - The type of the resource (e.g., 'Pod', 'Service', 'Deployment'). This will be used as the node's label in Neo4j.
-   `name`: `string` - The name of the resource.
-   `namespace`: `string` - The namespace the resource belongs to.
-   `properties`: `map<string, any>` - A map containing all other metadata from the Kubernetes resource object.

**Relationships:**
-   **`OWNS`**: A relationship from a parent resource to a child resource. This is derived from the `ownerReferences` field in a Kubernetes resource. For example, a `ReplicaSet` node would have an `OWNS` relationship to a `Pod` node.
-   **`SELECTS`**: A relationship from a `Service` to a `Pod`. This is derived from the `selector` field in a `Service` and the labels on a `Pod`.
-   **`MOUNTS`**: A relationship from a `Pod` to a `ConfigMap` or `Secret`. This is derived from the `volumes` and `volumeMounts` fields in a `Pod` specification.
