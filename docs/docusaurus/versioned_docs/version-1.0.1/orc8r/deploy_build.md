---
id: version-1.0.1-deploy_build
title: Building Orchestrator
hide_title: true
original_id: deploy_build
---

# Building Orchestrator Components

Start up your Docker daemon, `cd` to where you've cloned Magma, then:

```bash
cd magma/orc8r/cloud/docker
./build.py -a
```

If this is your first time building Orchestrator, this may take a while. We
recommend continuing onto the next step (Terraforming cloud resources) in a
different terminal instance or tab and coming back to this section when the
builds are completed.

When this job finishes, upload these images to your image registry:

```bash
export MAGMA_TAG=v1.0.1
for image in proxy controller prometheus-cache config-manager grafana
do
    ../../../orc8r/tools/docker/publish.sh -r REGISTRY -i ${image} -v ${MAGMA_TAG}
done
```

While we're here, you can build and publish the NMS containers as well:

```bash
cd magma/nms/fbcnms-projects/magmalte
docker-compose build magmalte
COMPOSE_PROJECT_NAME=magmalte ../../../orc8r/tools/docker/publish.sh -r REGISTRY -i magmalte -v ${MAGMA_TAG}
```
