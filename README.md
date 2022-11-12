[![CI](https://github.com/jsirianni/go-statsd-example/actions/workflows/ci.yaml/badge.svg)](https://github.com/jsirianni/go-statsd-example/actions/workflows/ci.yaml)

# StatsD Example

An example application that emits StatsD metrics. Intended to be consumed by
an [Open Telemetry](https://opentelemetry.io/) distribution, such as
[observIQ OTEL Collector](https://github.com/observIQ/observiq-otel-collector).

## Configuration

Configure using environment variables

| Name                               | Default  | Description |
| ---------------------------------- | -------- | ----------- |
| STATSD_HOST                        | Required | IP address or hostname of the StatsD collector. |
| STATSD_APP_NAME                    | Required | Application name, which is used as the metric name prefix. |
| OTEL_INCLUDE_RESOURCES             | `false`  | Whether or not to include `k8s.pod.name` and `k8s.namespace.name` tags. |
| OTEL_RESOURCE_ATTRIBUTES_POD_NAME  |          | The pod name, usually set using the Kubernetes [downward api](https://kubernetes.io/docs/concepts/workloads/pods/downward-api/). |
| OTEL_RESOURCE_ATTRIBUTES_NAMESPACE |          | The namespace name, usually set using the Kubernetes [downward api](https://kubernetes.io/docs/concepts/workloads/pods/downward-api/). |

## Example Usage

This example demonstrates how to send statsd metrics to a remote collector.
- STATSD_APP_NAME: `api`
- OTEL_INCLUDE_RESOURCES: `true`
- OTEL_RESOURCE_ATTRIBUTES_POD_NAME: Set with downward api
- OTEL_RESOURCE_ATTRIBUTES_NAMESPACE: Set with downward api

This example will forward metrics with the prefix `api` to the host `statsd-collector` (clusterIP service).
Metrics will include the tags `k8s.pod.name` and `k8s.pod.namespace`, which are dynamically set
based on the values injected by the downward API.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api
spec:
  selector:
    matchLabels:
      app: api
  template:
    metadata:
      labels:
        app: api
    spec:
      containers:
        - name: api
          image: ghcr.io/jsirianni/go-statsd-example:0.0.3
          env:
            - name: STATSD_APP_NAME
              value: api
            - name: STATSD_HOST
              value: "statsd-collector"
            - name: OTEL_INCLUDE_RESOURCES
              value: "true"
            - name: OTEL_RESOURCE_ATTRIBUTES_POD_NAME
              valueFrom:
                fieldRef:
                  apiVersion: v1
                  fieldPath: metadata.name
            - name: OTEL_RESOURCE_ATTRIBUTES_NAMESPACE
              valueFrom:
                fieldRef:
                  apiVersion: v1
                  fieldPath: metadata.namespace
```
