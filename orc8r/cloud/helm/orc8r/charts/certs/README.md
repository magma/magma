# Orchestrator Certificates

Orchestrator Certificates is used to apply a set of secrets required by magma orchestrator.

## Prerequisite

Install cert-manager:
```bash
helm repo add jetstack https://charts.jetstack.io
helm repo update

helm install cert-manager jetstack/cert-manager \
  --create-namespace \
  --namespace cert-manager \
  --set installCRDs=true
```

## Install certs

Install certs:
```bash
helm install certs .
```

## Configuration

The following table lists the configurable certs parameters.

| Parameter        | Description     | Default   |
| ---              | ---             | ---       |
| `create` | Set to `true` to create orc8r certs. | `false` |
| `domainName` | Domain Name | `localhost` |
| `bootstrapper.duration` | Certificate duration | `8760h` |
| `bootstrapper.renewBefore` | Renew Certificate before expiring | `24h` |
| `certifier.duration` | Certificate Duration | `87600h` |
| `certifier.renewBefore` | Renew Certificate before expiring | `24h` |
| `adminOperator.duration` | Certificate Duration | `87600h` |
| `adminOperator.renewBefore` | Renew Certificate before expiring | `24h` |
| `adminOperator.customIssuer` | Custom Issuer | `""` |
| `adminOperator.pkcs12.password` | pkcs12 password | `password` |
| `nms.duration` | Certificate Duration | `87600h` |
| `nms.renewBefore` | Renew Certificate before expiring | `24h` |
| `nms.customIssuer` | Custom Issuer | `""` |
| `controller.duration` | Certificate Duration | `87600h` |
| `controller.renewBefore` | Renew Certificate before expiring | `24h` |
| `controller.customIssuer` | Custom Issuer | `""` |
| `fluentd.duration` | Certificate Duration | `87600h` |
| `fluentd.renewBefore` | Renew Certificate before expiring | `24h` |
| `root.duration` | Certificate Duration | `87600h` |
| `root.renewBefore` | Renew Certificate before expiring | `24h` |
| `preInstallChecks.enabled` | Pre Install Checks | `true` |
