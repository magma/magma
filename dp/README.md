# Magma Domain Proxy

## Prerequsites
This instructions assumes that mentioned below utilities are installed and available on PATH:
  - [minikube](https://minikube.sigs.k8s.io/docs/start/)
  - [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)
  - [skaffold](https://skaffold.dev/docs/install/)
  - [helm](https://helm.sh/docs/intro/install/)
  - [kubens/kubectx](https://github.com/ahmetb/kubectx#installation)
  - [GNU Make version >= 4.2](https://www.gnu.org/software/make/)

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
make integration_tests
```

## Running integration tests with orc8r deployed.

```
make orc8r_integration_tests
```

## Enabling metrics

Make targets `orc8r` and `orc8r_integration_tests` can be run with metrics.
In this case prometheus and grafana PODS will be included in deployment.
To enable it evironment variable `DP_METRICS` have to be set e.g:

```bash
DP_METRICS=true make orc8r
```
