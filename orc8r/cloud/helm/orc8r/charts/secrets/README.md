# Orchestrator Secrets

Orchestrator Secrets is used to apply a set of secrets required by magma orchestrator.

## TL;DR;

```bash
# Copy secrets into temp directory
$ mkdir ../secretstemp && \
    cp -r <secrets>/* ../secretstemp/

# Apply secrets
helm template charts/secrets \
    --name orc8r-secrets \
    --namespace magma \
    --set-file secret.cert.files.rootCA.pem=../secretstemp/rootCA.pem \
    --set-file secret.cert.files.controller.crt=../secretstemp/controller.crt \
    --set-file secret.cert.files.controller.key=../secretstemp/controller.key \
    --set-file secret.cert.files.admin_operator.pem=../secretstemp/admin_operator.pem \
    --set-file secret.cert.files.admin_operator.key.pem=../secretstemp/admin_operator.key.pem \
    --set-file secret.cert.files.fluentd.pem=../secretstemp/fluentd.pem \
    --set-file secret.cert.files.fluentd.key=../secretstemp/fluentd.key \
    --set=docker.registry=docker.io \
    --set=docker.username=username \
    --set=docker.password=password |
kubectl apply -f -
```

Note: If deploying metrics with Thanos add `--set=thanos_enabled=true` when templating

## Overview

This chart installs a set to secrets required by magma orchestrator.
The secrets are expected to be provided as files and placed under temp Directory

```bash
$ ls ../secretstemp
rootCA.pem controller.crt controller.key admin_operator.pem admin_operator.key.pem fluentd.pem fluentd.key

```
## Image Pull Secret

This chart can also be used to install image pull secret which is required
when working with private docker registries.

## Configuration

The following table lists the configurable secret locations, 
docker credentials and their default values.

| Parameter        | Description     | Default   |
| ---              | ---             | ---       |
| `create` | Set to ``true`` to create orc8r secrets. | `false` |
| `secret.certs.enabled` | Enable certs. | `false` |
| `secret.certs.files.controller.crt` | controller pem file. | `""` |
| `secret.certs.files.controller.key` | controller key file. | `""` |
| `secret.certs.files.rootCA.pem` | rootCA.pem file. | `""` |
| `secret.configs.enabled` | Enable configs. | `false` |
| `docker.registry` | Docker registry. | `""` |
| `docker.username` | Docker registry username. | `""` |
| `docker.password` | Docker registry password. | `""` |
