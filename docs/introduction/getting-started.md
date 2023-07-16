# Getting Started

Dex operator is a Kubernetes operator for [Dex](https://dexidp.io/), an identity service that uses OpenID Connect to drive authentication for other apps.

## Install Dex Operator

```shell
# Clone the repository
gh repo clone rudeigerc/dex-operator
# Install CRDs
make install
# Deploy Dex Operator
make deploy
```

## Create a Dex instance

```shell
kubectl apply -f config/samples/dex_v1alpha1_dexcluster.yaml
```
