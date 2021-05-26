---
id: deploy_build
title: Build Orchestrator
hide_title: true
---

# Build Orchestrator

We'll build and publish Magma container images and Helm charts from source.

If this is your first time building Orchestrator, each step may take a while.
We recommend continuing to the next step (Terraforming cloud resources) in a
different shell and coming back to this section as each command completes.

## Prerequisites

We assume `MAGMA_ROOT` is set as described in the
[deployment intro](./deploy_intro.md).

We'll go over how to publish images to Docker Hub and charts to a GitHub repo.
For this you'll need

- a [Docker Hub account](https://docs.docker.com/docker-hub/), to host your
`REGISTRY`
- a [GitHub personal access token](https://docs.github.com/github/authenticating-to-github/creating-a-personal-access-token),
which we'll call `GITHUB_ACCESS_TOKEN`

You can of course publish to a container registry or
[chart repository](https://helm.sh/docs/topics/chart_repository/) of your
choice.

First, start up your Docker daemon then
[log in](https://docs.docker.com/engine/reference/commandline/login/) to the
Docker Hub registry

```bash
$ docker login

Login with your Docker ID to push and pull images from Docker Hub. If you don't have a Docker ID, head over to https://hub.docker.com to create one.
Username: REGISTRY_USERNAME
Password:
Login Succeeded
```

## Build and publish container images

We provide scripts to build and publish images. The publish script is provided
as a starting point, as individual needs may vary.

First define some necessary variables

```bash
export PUBLISH=${MAGMA_ROOT}/orc8r/tools/docker/publish.sh  # or add to path
export REGISTRY=registry.hub.docker.com/REGISTRY  # or desired registry
export MAGMA_TAG=1.4.0-master  # or desired tag
```

Checkout the v1.4 release branch
```bash
cd ${MAGMA_ROOT}
git fetch origin
git checkout -b v1.4 origin/v1.4
```

Build and publish Orchestrator images

```bash
cd ${MAGMA_ROOT}/orc8r/cloud/docker
./build.py --all
for image in controller nginx ; do ${PUBLISH} -r ${REGISTRY} -i ${image} -v ${MAGMA_TAG} ; done
```

Build and publish NMS images

```bash
cd ${MAGMA_ROOT}/nms/app/packages/magmalte
docker-compose build magmalte
COMPOSE_PROJECT_NAME=magmalte ${PUBLISH} -r ${REGISTRY} -i magmalte -v ${MAGMA_TAG}
```

## Build and publish Helm charts

NOTE: only choose one of the below options and then run the script

### Option 1: Publish to artifactory

We'll build the Orchestrator Helm charts, as well as publish them to an artifactory.

Define some necessary variables

```bash
export HELM_CHART_ARTIFACTORY_URL=ARTIFACTORY_URL
export HELM_CHART_MUSEUM_REPO=REPO_NAME
export HELM_CHART_MUSEUM_USERNAME=ARTIFACTORY_USERNAME
export HELM_CHART_MUSEUM_TOKEN=ARTIFACTORY_ACCESS_TOKEN
```

Next we'll run the package script. This script will package and publish the
necessary helm charts to the artifactory. The script expects a deployment
type to be provided, which will determine which orc8r modules are deployed.

### Option 2: Publish to a private Github repo

We'll build the Orchestrator Helm charts, as well as publish them to a
[GitHub repo acting as a Helm chart repo](https://blog.softwaremill.com/hosting-helm-private-repository-from-github-ff3fa940d0b7).

To start, create a private GitHub repo to use as your Helm chart repo. We'll
refer to this as `GITHUB_REPO`.

Define some necessary variables

```bash
export GITHUB_REPO=GITHUB_REPO_NAME
export GITHUB_REPO_URL=GITHUB_REPO_URL
export GITHUB_USERNAME=GITHUB_USERNAME
export GITHUB_ACCESS_TOKEN=GITHUB_ACCESS_TOKEN
```

Next we'll run the package script. This script will package and publish the
necessary Helm charts to the `GITHUB_REPO`. The script expects a deployment
type to be provided, which will determine which orc8r modules are deployed.

### Run the package script

The valid deployment type options are

- `fwa`
- `federated_fwa`
- `all`

Run the package script

```bash
$ ${MAGMA_ROOT}/orc8r/tools/helm/package.sh -d fwa # or chosen deployment type

...

Uploaded orc8r charts successfully.
```
You can add `-v` option to overwrite the versions of the chart.
```bash
${MAGMA_ROOT}/orc8r/tools/helm/package.sh -d fwa  -v 1.5
```
