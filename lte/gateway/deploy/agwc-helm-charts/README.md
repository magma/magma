# AGW Helm Deployment

## Configuration

The following table list the configurable parameters of the agw chart and their default values.

| Parameter        | Description     | Default   |
| ---              | ---             | ---       |
| `imagePullSecrets` | Reference to one or more secrets to be used when pulling images. | `[]` |
| `secrets.create` | Create agwc secrets. See charts/secrets subchart. | `false` |
| `secret.certs` | Secret name containing agwc certs. | `agwc-secrets-certs` |
| `secret.configs` | Secret name containing agwc configs. | `agwc-secrets-configs` |
| `secret.envdir` | Secret name containing agwc envdir. | `agwc-secrets-envdir` |
