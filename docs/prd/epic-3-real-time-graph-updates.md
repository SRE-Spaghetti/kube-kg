### Epic 3: Real-Time Graph Updates

**Goal:** Implement the client for the KubeView Server-Sent Events (SSE) stream and develop the event processor to handle `add`, `update`, and `delete` events, keeping the Neo4j graph continuously up-to-date.

---

#### Story 3.1: KubeView SSE Client

*As a developer,*
*I want a client that can connect to the KubeView SSE stream and parse incoming events,*
*so that I can receive real-time updates about Kubernetes resource changes.*

**Acceptance Criteria:**

1.  A function `StreamUpdates` is implemented that connects to the `/updates?clientID={id}` endpoint.
2.  The function correctly reads the `text/event-stream` format, parsing `event:` and `data:` lines.
3.  The `data` JSON is successfully unmarshaled into a `kubeview.Event` struct.
4.  Parsed events are sent to a Go channel for consumption by the event processor.
5.  `ping` events are handled correctly to maintain the connection without being sent to the event channel.
6.  The client includes a retry mechanism with backoff to automatically reconnect if the connection is lost.
7.  An integration test verifies that the client can connect to a mock SSE stream, parse events, and send them to the channel.

---

#### Story 3.2: Real-Time Event Processor

*As a developer,*
*I want a processor that consumes resource events and applies the corresponding changes to the Neo4j graph,*
*so that the graph remains a consistent, real-time reflection of the cluster state.*

**Acceptance Criteria:**

1.  A function `StartEventProcessor` is created that listens for events on the channel provided by the SSE client.
2.  For an `add` event, the processor uses the mapping logic (from Story 2.2) and the `Neo4jClient` to `MERGE` the new node and its relationships into the graph.
3.  For an `update` event, the processor uses the `Neo4jClient` to `MERGE` the node's new properties and re-evaluates its relationships, adding or removing them as needed.
4.  For a `delete` event, the processor uses the `Neo4jClient` to execute a `DETACH DELETE` query to remove the node and all its relationships.
5.  Each event processed is tracked within its own OpenTelemetry trace.
6.  Unit tests are written for the processor logic for each event type (`add`, `update`, `delete`), mocking the `Neo4jClient`.
7.  An integration test verifies that when an event is sent to the input channel, the correct change is persisted in the test Neo4j database.
