# Design

## Dex

`Dex` is a custom resource definition declaratively describing a [Dex](https://dexidp.io/) instance.


```yaml
apiVersion: dex.rudeigerc.dev/v1alpha1
kind: DexCluster
metadata:
  labels:
    app.kubernetes.io/name: dexcluster
    app.kubernetes.io/instance: dexcluster-sample
    app.kubernetes.io/part-of: dex-operator
    app.kubernetes.io/managed-by: kustomize
    app.kubernetes.io/created-by: dex-operator
  name: dexcluster-sample
spec:
  image: dexidp/dex:v2.37.0
  imagePullPolicy: IfNotPresent
  config: |
    issuer: http://127.0.0.1:5556/dex
    storage:
      type: sqlite3
      config:
        file: /tmp/dex.db
    web:
      http: 0.0.0.0:5556
    telemetry:
      http: 0.0.0.0:5558
    staticClients:
    - id: example-app
      redirectURIs:
      - 'http://127.0.0.1:5555/callback'
      name: 'Example App'
      secret: ZXhhbXBsZS1hcHAtc2VjcmV0
    connectors:
    - type: mockCallback
      id: mock
      name: Example
```
