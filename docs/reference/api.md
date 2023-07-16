# API Reference

## Packages
- [dex.rudeigerc.dev/v1alpha1](#dexrudeigercdevv1alpha1)


## dex.rudeigerc.dev/v1alpha1

Package v1alpha1 contains API Schema definitions for the dex v1alpha1 API group

### Resource Types
- [DexCluster](#dexcluster)



#### DexCluster



DexCluster is the Schema for the dexclusters API



| Field | Description |
| --- | --- |
| `apiVersion` _string_ | `dex.rudeigerc.dev/v1alpha1`
| `kind` _string_ | `DexCluster`
| `metadata` _[ObjectMeta](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#objectmeta-v1-meta)_ | Refer to Kubernetes API documentation for fields of `metadata`. |
| `spec` _[DexClusterSpec](#dexclusterspec)_ |  |


#### DexClusterSpec



DexClusterSpec defines the desired state of DexCluster

_Appears in:_
- [DexCluster](#dexcluster)

| Field | Description |
| --- | --- |
| `config` _string_ |  |
| `replicas` _integer_ | Number of desired pods. This is a pointer to distinguish between explicit zero and not specified. Defaults to 1. |
| `image` _string_ |  |
| `imagePullPolicy` _[PullPolicy](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#pullpolicy-v1-core)_ |  |
| `imagePullSecrets` _[LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.27/#localobjectreference-v1-core) array_ |  |




