---
id: version-1.0.0-deploy_build
title: Building Federation Gateway
hide_title: true
original_id: deploy_build
---

# Building Federation Gateway Components

Start up your Docker daemon, `cd` to where you've cloned Magma, then:

```bash
cd magma/feg/gateway/docker
docker-compose build --parallel
```

If this is your first time building the FeG, this may take a while. When this
job finishes, upload these images to your image registry:

```bash
../../../orc8r/tools/docker/publish.sh -r <REGISTRY> -i gateway_python
../../../orc8r/tools/docker/publish.sh -r <REGISTRY> -i gateway_go
../../../orc8r/tools/docker/publish.sh -r <REGISTRY> -i gateway_radius
```
