# Magma Domain Proxy


This instructions assumes that both minikube, kubectl and skaffold utilities are installed and available on PATH.
Detailed instruction on how to install utilities can be found here:
  - [minikube](https://minikube.sigs.k8s.io/docs/start/)
  - [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)
  - [skaffold](https://skaffold.dev/docs/install/)
  - [helm](https://helm.sh/docs/intro/install/)
  - [kubens/kubectx](https://github.com/ahmetb/kubectx#installation)

## Local deployment without orc8r.

```
make
```

## Local deployment with orc8r.

```
make orc8r
```

## Running integration tests.

```
CI=true make
```

## Running integration tests with orc8r deployed.

```
make orc8r_integration_tests
```
