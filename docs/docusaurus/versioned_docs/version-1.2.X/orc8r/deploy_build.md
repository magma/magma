---
id: version-1.2.0-deploy_build
title: Build Orchestrator
hide_title: true
original_id: deploy_build
---

# Build Orchestrator

We'll build and publish Magma container images and Helm charts from source.

If this is your first time building Orchestrator, each step may take a while.
We recommend continuing to the next step (Terraforming cloud resources) in a
different shell and coming back to this section as each command completes.

## Prerequisites

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

## Build and Publish Container Images

We provide scripts to build and publish images. The publish script is provided
as a starting point, as individual needs may vary.

First define some necessary variables

```bash
export PUBLISH=MAGMA_ROOT/orc8r/tools/docker/publish.sh  # or add to path
export REGISTRY=registry.hub.docker.com/REGISTRY  # or desired registry
export MAGMA_TAG=1.1.0-master  # or desired tag
```

Build and publish Orchestrator images

```bash
cd MAGMA_ROOT/orc8r/cloud/docker
./build.py -a
for image in controller nginx ; do ${PUBLISH} -r ${REGISTRY} -i ${image} -v ${MAGMA_TAG} ; done
```

Build and publish NMS images

```bash
cd MAGMA_ROOT/nms/app/packages/magmalte
docker-compose build magmalte
COMPOSE_PROJECT_NAME=magmalte ${PUBLISH} -r ${REGISTRY} -i magmalte -v ${MAGMA_TAG}
```

## Build and Publish Helm Charts

We'll build the Orchestrator Helm charts, as well as publish them to a
[GitHub repo acting as a Helm chart repo](https://blog.softwaremill.com/hosting-helm-private-repository-from-github-ff3fa940d0b7).

To start, create a private GitHub repo to use as your Helm chart repo. We'll
refer to this as `GITHUB_REPO`.

Next, package the Magma Helm charts and publish them to the GitHub
repo

```bash
mkdir ~/magma-charts && cd ~/magma-charts
git init
helm package MAGMA_ROOT/orc8r/cloud/helm/orc8r/ && helm repo index .
git add . && git commit -m 'Initial chart commit'
git remote add origin GITHUB_REPO_URL && git push -u origin master
```

To confirm, reference the published charts locally

```bash
helm repo add GITHUB_REPO --username GITHUB_USERNAME --password GITHUB_ACCESS_TOKEN \
      'https://raw.githubusercontent.com/GITHUB_USERNAME/GITHUB_REPO/master/'
helm repo update && helm repo list  # should list the GITHUB_REPO repository
helm search repo GITHUB_REPO  # should list the GITHUB_REPO chart
```

Finally, check out the newly-created `index.yaml`. It should contain the
Orchestrator chart's version, which we'll use as a Terraform value in the next
section.
