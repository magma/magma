---
id: deployment
title: Deployment
hide_title: true
---
# Deployment

The FeG installation process assumes that the necessary FeG docker images
are available in a docker registry. If this isn't the case, refer to
the [FeG Docker Setup](docker_setup)
for how to build and publish the images.

## Installation

The installation is done using `install_gateway.sh` located at `magma/orc8r/tools/docker`.
There are 3 files needed in addition to the install script:

* rootCA.pem
* control_proxy.yml
* .env

To install:

```console
INSTALL_HOST [~/]$ sudo ./install_gateway.sh
```

After this completes, you should see: `Installed successfully!!`

## Registration

After installation, the next step is to register the gateway with the Orchestrator.
To do so:

```console
INSTALL_HOST [~/]$ cd /var/opt/magma/docker
INSTALL_HOST [/var/opt/magma/docker]$ docker-compose exec magmad /usr/local/bin/show_gateway_info.py
```

This will output a hardware ID and a challenge key. This information must be registered with
the NMS. To do this, go to NMS and select Network Management in the lower left hand corner.
Then select the appropriate network that the gateway should be added to. Navigate to the
"Configure Gateways" tab and click "Add Gateway". Here you will be prompted for the gateway
information that was shown, as well as a gateway name and gateway ID.
To verify that the gateway was correctly registered, run:

```console
INSTALL_HOST [~/]$ cd /var/opt/magma/docker
INSTALL_HOST [/var/opt/magma/docker]$ docker-compose exec magmad /usr/local/bin/checkin_cli.py
```

## Configuration

The final step is to configure the gateway. To do this, go to NMS and click on the edit icon
at the right side of the newly listed gateway. Select the `Magma` tab and enter the appropriate
configs.

At this point more specific configuration can be added (LTE, WiFi, etc.). This is dependent
upon the deployment.

## Upgrades

The installation process places the `upgrade_gateway.sh` script at `/var/opt/magma/docker`.
To upgrade, run:

```console
INSTALL_HOST [~/]$ cd /var/opt/magma/docker
INSTALL_HOST [/var/opt/magma/docker]$ sudo ./upgrade_gateway.sh <github_tag>
```

This will upgrade the gateway with the github version supplied.
