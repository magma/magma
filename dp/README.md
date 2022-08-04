# Magma Domain Proxy

## Prerequisites

This instructions assumes that mentioned below utilities are installed and available on PATH:

- [minikube](https://minikube.sigs.k8s.io/docs/start/)
- [kubectl](https://kubernetes.io/docs/tasks/tools/install-kubectl-linux/)
- [skaffold](https://skaffold.dev/docs/install/)
- [helm](https://helm.sh/docs/intro/install/)
- [kubens/kubectx](https://github.com/ahmetb/kubectx#installation)
- [GNU Make version >= 4.2](https://www.gnu.org/software/make/)

## Local deployment without orc8r

```bash
make
```

## Local deployment with orc8r

```bash
make orc8r
```

## Running integration tests

```bash
make integration_tests
```

## Running integration tests with orc8r deployed

```bash
make orc8r_integration_tests
```

## Enabling features

### Metrics

Make targets `orc8r` and `orc8r_integration_tests` can be run with metrics.
In this case prometheus and grafana PODS will be included in deployment.
To enable it evironment variable `DP_METRICS` have to be set e.g:

```bash
DP_METRICS=true make orc8r
```

### Services needed by Certification process

Device certification process require extra services to be deployed.
Crl server and Harness certification service.
To enable them `DP_CERTIFICATION` environment variable have to be set

```bash
DP_CERTIFICATION=true make orc8r
```
