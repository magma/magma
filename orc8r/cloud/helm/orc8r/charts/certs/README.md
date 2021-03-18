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
helm install certs . \
  --set dnsDomain=magma.shubhamtatvamasi.com
```

## Configuration

The following table lists the configurable certs parameters.

| Parameter        | Description     | Default   |
| ---              | ---             | ---       |
| `create` | Set to ``true`` to create orc8r certs. | `false` |
| `dnsDomain` | Domain Name | `localhost` |
| `duration` | Certificate duration | `87600h` |
| `adminOperator.customIssuer` | Custom Issuer | `""` |
| `adminOperator.pkcs12.password` | pkcs12 password | `password` |
| `nms.customIssuer` | Custom Issuer | `""` |
| `controller.customIssuer` | Custom Issuer | `""` |
| `preInstallChecks.enabled` | Pre Install Checks | `true` |
