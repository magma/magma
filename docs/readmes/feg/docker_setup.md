---
id: docker_setup
title: FeG Docker Setup
hide_title: true
---
# FeG Docker Setup

The FeG runs each service in its own Docker container.
Production services are defined in `docker-compose.yml`.
Development services are defined in `docker-compose.override.yml`. 
The development `test` service is used to run unit tests and regenerate Swagger/Protobuf code.
The development `test` service can also be used to perform other development-related procedures.

## Requirements

To run the FeG with docker, both docker and docker compose must be installed.
* Follow [these steps](https://docs.docker.com/install/) to install docker
* Follow [these steps](https://docs.docker.com/compose/install/) to install docker compose

NOTE: If you are running the FeG on Mac, you will need to increase the memory
limit of the docker daemon to at least 4GB to build the images. Otherwise,
when building the Go image, you may see an error message similar to this:
`/usr/local/go/pkg/tool/linux_amd64/link: signal: killed`.

The `rootCA.pem` certificate must be located in the `.cache/test_certs` folder,
so that it can be mounted into the appropriate containers from there.

## Development

Follow these steps to run the FeG services:
1. `cd magma/feg/gateway/docker`
2. `docker-compose build`
3. `docker-compose up -d`

Each service should now be running in each of its containers. 
By default, both production and development services should be running. 
To place a shell into the test container, run the command:

`docker-compose exec test /bin/bash`

The test container contains the mounted source code and configuration settings.
The mounted source code and configuration settings can be changed externally 
and the changes will be reflected inside the test container. 
Run the command `make precommit` in the container before submitting a patch.

To make changes to currently running FeG services, the containers must be rebuilt and restarted:
1. `docker-compose down`
2. `docker-compose build`
3. `docker-compose up -d` 

To manage the containers, the following commands are useful:
* `docker-compose ps` (get status of each container)
* `docker-compose logs -f` (tail logs of all containers)
* `docker-compose logs -f <service name>` (tail logs of a particular service)
* `docker-compose down` (stop all services)

## Publishing the images

To push production images to a private docker registry, use the following script:
```
[/magma/feg/gateway/docker]$ ../../../orc8r/tools/docker/publish.sh -r <REGISTRY> -i gateway_python
[/magma/feg/gateway/docker]$ ../../../orc8r/tools/docker/publish.sh -r <REGISTRY> -i gateway_go
```
