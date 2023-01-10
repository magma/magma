---
id: deploy_build
title: Building Carrier Wifi Gateway
hide_title: true
---

# Building Carrier Wifi Gateway Components

Start up your Docker daemon, `cd` to where you've cloned Magma, then:

```bash
cd magma/cwf/gateway/docker
docker-compose build --parallel
```

If this is your first time building the CWAG, this may take a while. When this
job finishes, upload these images to your image registry:

```bash
../../../orc8r/tools/docker/publish.sh -r <REGISTRY> -i gateway_python
../../../orc8r/tools/docker/publish.sh -r <REGISTRY> -i gateway_go
../../../orc8r/tools/docker/publish.sh -r <REGISTRY> -i gateway_radius
../../../orc8r/tools/docker/publish.sh -r <REGISTRY> -i gateway_pipelined
../../../orc8r/tools/docker/publish.sh -r <REGISTRY> -i gateway_sessiond
```
