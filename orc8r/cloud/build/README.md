# Orc8r Docker setup

*NOTE: Docker support for orc8r is not yet official*

## Container setup

Orc8r consists of 2 containers: one for the proxy, and one for all the
controller services. We use supervisord to spin multiple services within
these containers.

NOTE: The multiple services per container model was adopted to model the
current Vagrant setup and for easier migration, and we will soon migrate to
one container per microservice model which is more appropriate.

For development, we use a postgresql container as the datastore. For
production, it is advisable to use a hosted solution like AWS RDS.

## How to build the images

Since orc8r can include modules outside the magma codebase, we use a wrapper
python script which creates a temporary folder as the docker build context.
The temporary folder contains all the modules necessary, and is created based
on the module.yml config.

Build the docker images using:
```
./build.py
```
To use a different module.yml file, run something similar to:
```
MAGMA_MODULES_FILE=../../../modules.yml ./build.py
```

## How to run

To run and manage the containers, use the following commands:
```
docker-compose up -d
docker-compose ps
docker-compose down
```
To tail the logs from the containers, use one of these commands:
```
docker-compose logs -f
docker-compose logs -f controller
```
To create a shell inside a container, run:
```
docker-compose exec controller bash
```
