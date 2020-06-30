---
id: deploy_build
title: Building Orchestrator
hide_title: true
---

# Building Orchestrator Components

Start up your Docker daemon, `cd` to where you've cloned Magma, then

```bash
cd magma/orc8r/cloud/docker
./build.py -a
```

If this is your first time building Orchestrator, this may take a while. We
recommend continuing onto the next step (Terraforming cloud resources) in a
different terminal instance or tab and coming back to this section when the
builds are completed.

When this job finishes, upload these images to your image registry as follows.
We provide a script to publish an individual image as a starting point, though
individual needs may vary.

To start, define some necessary variables

```bash
export PUBLISH=../../../orc8r/tools/docker/publish.sh  # or add to path
export REGISTRY=registry.hub.docker.com/YOURREGISTRY  # replace with desired registry
export MAGMA_TAG=v1.1.0-master  # or alternative desired tag
```

Publish orc8r images to the registry

```bash
for image in proxy controller prometheus-cache alertmanager-configurer prometheus-configurer grafana; do
    ${PUBLISH} -r ${REGISTRY} -i ${image} -v ${MAGMA_TAG}
done
```

While we're here, build and publish NMS images as well

```bash
cd magma/symphony/app/fbcnms-projects/magmalte
docker-compose build magmalte
COMPOSE_PROJECT_NAME=magmalte ${PUBLISH} -r ${REGISTRY} -i magmalte -v ${MAGMA_TAG}
```
