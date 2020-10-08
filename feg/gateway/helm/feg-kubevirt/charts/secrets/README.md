# FeG Secrets

Feg Secrets is used to apply a set of secrets for a federated gateway. These 
secrets provide permanent gateway identification so that gateway pods will be 
recreated with the same hwid and challenge key.

## TL;DR;

```bash
# Copy gw_challenge and snowflake into temp dir
    cp -r <secrets>/* ../temp/
$ ls temp/
snowflake gw_challenge.key

# Apply secrets
helm template charts/secrets \
    --name <feg release name> \
    --set-file secrets.gwinfo.gw_challenge.key=../temp/gw_challenge.key \
    --set-file secrets.gwinfo.snowflake=../temp/snowflake \
    --namespace magma | kubectl -n magma apply -f -
```

## Overview

This chart installs a secret that serves as identifiers for the gateway. 
The secrets are expected to be provided as files and placed in temp dir.
```bash
$ ls temp/
snowflake  gw_challenge.key
```

## Creating Gateway Info
If creating a gateway for the first time, you'll need to create a snowflake
and challenge key before installing the gateway. To do so:

```
$ docker login <DOCKER REGISTRY>
$ docker pull <DOCKER REGISTRY>/gateway_python:<IMAGE_VERSION>
$ docker run -d <DOCKER_REGISTRY>/gateway_python:<IMAGE_VERSION> python3.5 -m magma.magmad.main

This will output a container ID such as
f3bc383a95db16f2e448fdf67cac133a5f9019375720b59477aebc96bacd05a9

Run the following, substituting your container ID here
$ docker cp <container ID>:/etc/snowflake charts/secrets/.secrets
$ docker cp <container ID>:/var/opt/magma/certs/gw_challenge.key /charts/secrets/.secrets
```

Otherwise if redeploying a gateway with permanent gwinfo, copy the existing 
snowflake from `etc/snowflake` and challenge key at 
`/var/opt/magma/certs/gw_challenge.key`

## Configuration

The following table lists the configurable secret locations and 
their default values.

| Parameter        | Description     | Default   |
| ---              | ---             | ---       |
| `create` | Set to ``true`` to create feg secrets. | `false` |
| `secret.enabled` | Enable gwinfo secrets. | `false` |
| `secret.gwinfo.gw_challenge.key` | gw_challenge.key file | `""` |
| `secret.gwinfo.snowflake` | snowflake file. | `""` |
