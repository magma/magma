---
id: docker_setup
title: FeG Docker Setup
hide_title: true
---
# FeG Docker Setup

*NOTE: Docker support for FeG is not yet official*

The FeG can be run in two different modes: development or production.
The development mode is setup to make tests easy to run and it runs all of the
services in one container. On the other hand, in the production environment,
each service runs in its own container. In production, containers are optimized
for size and so many utilities are removed.

## Requirements

To run the FeG with docker, both docker and docker compose must be installed.
* Follow [these steps](https://docs.docker.com/install/) to install docker
* Follow [these steps](https://docs.docker.com/compose/install/) to install docker compose

NOTE: If you are running the FeG on Mac, you will need to increase the memory
limit of the docker daemon to at least 4GB to build the images. Otherwise,
when building the Go image, you may see an error message similar to this:
`/usr/local/go/pkg/tool/linux_amd64/link: signal: killed`.

The certificates (eg. rootCA.key) are expected to be in the `.cache/test_certs`
folder. That folder is mounted into the appropriate containers.

## Development

Follow these steps to run the FeG in a docker container for development:
1. `cd magma/feg/gateway/docker`
2. `docker-compose up -d`
3. `docker exec -ti feg /bin/bash`

Now, your shell should be inside a container with the source code
and configs mounted. The code can be changed externally, and commands such as
`make gen build test precommit run` should succeed.
Please run `make precommit` in this container before submitting a patch.

The development environment has a single container, `feg`, which is used to
run all of the magma services. Systemd is used for managing the services,
and so it should be set as the init system in `magmad.yml`.

## Production

Follow these steps to run the FeG in docker containers for production:
1. `cd magma/feg/gateway/docker`
2. `docker-compose -f docker-compose.prod.yml up -d`

Each service should now be running in its own docker container. Docker is used
to manage the services, and so it should be set as the init system in
`magmad.yml`. Also, ensure that `enable_systemd_tailer = False` in `magmad.yml`.

To manage the containers the following commands are useful:
* `docker-compose -f docker-compose.prod.yml ps` (get status of each container)
* `docker-compose -f docker-compose.prod.yml logs -f` (tail logs of all containers)
* `docker logs -f <service name>` (tail logs of a particular service)
* `docker-compose -f docker-compose.prod.yml down` (stop all services)

## Publishing the images

To push the production images to a private docker registry, use the following script:
```
[/magma/feg/gateway/docker]$ ../../../orc8r/tools/docker/publish.sh -r <REGISTRY> -i gateway_python
[/magma/feg/gateway/docker]$ ../../../orc8r/tools/docker/publish.sh -r <REGISTRY> -i gateway_go
```
