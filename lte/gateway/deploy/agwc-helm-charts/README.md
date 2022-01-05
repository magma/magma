# AGW Helm Deployment

## Configuration

The following table list the configurable parameters of the agw chart and their default values.

| Parameter        | Description     | Default   |
| ---              | ---             | ---       |
| `secret.certs` | Secret name containing agwc certs. | `agwc-secrets-certs` |
| `image.repository` | Repository for agwc images. | `` |
| `image.pullPolicy` | Pull policy for agwc images. | `` |
| `image.tag` | Tag for agwc image. | `latest` |
| `config.domainName` | Orchestrator domain name. | `` |
| `persistant.name` | Secret name containing agwc certs. | `agwc-claim` |

## Installation

# Create magma namespace

~ kubectl create namespace magma

# Create rootca certificate secret needed to communicate with orc8r

~ kubectl create secret generic agwc-secret-certs --from-file=rootCA.pem=rootCA.pem --namespace magma

# Deploy an AGW after updating values.yaml

~ cd lte/gateway/deploy/agwc-helm-charts
~ helm --debug install agwc --namespace magma . --values=values.yaml

# Delete the AGW if needed using,

~ helm uninstall agwc --namespace magma
