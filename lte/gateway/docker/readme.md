# Containerized AGW

This folder contains container image definitions for AGW services.

## Preparing the host

The containerized AGW requires a specific setup on the host machine.
The script ../agw_install_docker.sh can be used to configure a host
accordingly.

## Building the images

The images can be built with `docker-compose build`.

## Deploying the containerized AGW on AWS

* Build the docker images and push to a registry
* Look at ../../../experimental/cloudstrapper/README.md for an AWS setup
    * If you are only interested in running the AGW, you only need steps 1 and 5.1
    * After running step 5, you should be able to SSH to the created EC2 instance
* On the EC2 instance
    * optional: Copy custom `rootCA.pem` to  `/var/opt/magma/certs` with permissions 400
    * Copy `agw_install_docker.sh` and run script to prepare the host
    * Make changes to `/var/opt/magma/docker/.env` to include the right docker information
      * Example for the registry setting: `DOCKER_REGISTRY=registry.hub.docker.com/arunuke/`
    * Make changes to files in `/var/opt/magma/configs` if needed
        * pipelined, dnsd, enodebd, spgw, mme all will have their eth0/eth1 changed to newer values based on local interface names (eth0 and eth2 references to use the first interface for SGi and eth1 references to use the second interface for S1)
        * pipelined will also set dp_router_enabled to false
    * Make changes to config files and restart services by running `/var/opt/magma/docker/agw_upgrade.sh` or by running the `agw_install_docker.sh` script