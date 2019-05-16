# Orchestrator Secrets

Orchestrator Secrets is used to apply a set of secrets required by magma orchestrator.

## TL;DR;

```bash
# Copy secrets into subchart root
$ mkdir charts/secrets/.secrets && \
    cp -r <secrets>/* charts/secrets/.secrets/

# Apply secrets
helm template charts/secrets \
    --name orc8r-secrets \
    --namespace magma \
    --set=docker.registry=docker.io \
    --set=docker.username=username \
    --set=docker.password=password |
kubectl apply -f -
```

## Overview

This chart installs a set to secrets required by magma orchestrator.
The secrets are expected to be provided as files and placed under
secrets subchart root.
```bash
$ ls charts/secrets/.secrets
certs  configs
$ pwd
magma/orc8r/cloud/helm/orc8r
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
| `secret.certs` | Root relative certs directory. | `.secrets/certs` |
| `secret.configs` | Root relative configs directory. | `.secrets/configs` |
| `docker.registry` | Docker registry. | `""` |
| `docker.username` | Docker registry username. | `""` |
| `docker.password` | Docker registry password. | `""` |