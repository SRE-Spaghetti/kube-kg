# kube-kg Helm Chart

This chart deploys the `kube-kg` service to a Kubernetes cluster.

## Dependencies

This chart has optional dependencies on the following charts:

* `kubeview`
* `jaeger`
* `neo4j`

These dependencies can be enabled or disabled in the `values.yaml` file.

## Installation

To install the chart, run the following command:

```sh
helm install my-release ./charts/kube-kg
```

## Configuration

The following table lists the configurable parameters of the `kube-kg` chart and their default values.

| Parameter | Description | Default |
| :--- | :--- | :--- |
| `replicaCount` | Number of replicas to deploy. | `1` |
| `image.repository` | Image repository. | `kube-kg` |
| `image.pullPolicy` | Image pull policy. | `IfNotPresent` |
| `image.tag` | Image tag. | `latest` |
| `service.type` | Service type. | `ClusterIP` |
| `service.port` | Service port. | `80` |
| `kubeview.enabled` | Enable the kubeview dependency. | `true` |
| `jaeger.enabled` | Enable the jaeger dependency. | `false` |
| `neo4j.enabled` | Enable the neo4j dependency. | `true` |
