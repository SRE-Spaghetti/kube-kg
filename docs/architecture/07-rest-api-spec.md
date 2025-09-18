## 7. REST API Spec

```yaml
openapi: 3.0.0
info:
  title: Kube-KG Internal API
  version: 1.0.0
  description: An internal API for the Kube-KG service to report health and trigger data refreshes.
servers:
  - url: http://localhost:8080
    description: Local development server

paths:
  /health:
    get:
      summary: Health Check
      description: Reports the health of the service and its dependencies.
      responses:
        '200':
          description: Service is healthy.
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                    example: "OK"
                  dependencies:
                    type: object
                    properties:
                      kubeview:
                        type: string
                        example: "OK"
                      neo4j:
                        type: string
                        example: "OK"
        '503':
          description: Service is unhealthy.
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                    example: "Unavailable"
                  dependencies:
                    type: object
                    properties:
                      kubeview:
                        type: string
                        example: "Unavailable"
                      neo4j:
                        type: string
                        example: "OK"

  /refresh:
    post:
      summary: Trigger Refresh
      description: Triggers a full re-synchronization of the knowledge graph from KubeView.
      responses:
        '222':
          description: Refresh process has been accepted and started in the background.
          content:
            application/json:
              schema:
                type: object
                properties:
                  status:
                    type: string
                    example: "Refresh triggered"

```
