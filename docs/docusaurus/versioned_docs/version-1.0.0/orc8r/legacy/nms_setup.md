---
id: version-1.0.0-nms_setup
title: Setting up the NMS
hide_title: true
original_id: nms_setup
---
# Setting up the NMS
The NMS is the UI for managing, configuring, and monitoring networks. To set up the NMS, we will first need the orc8r to be setup (the NMS needs magma certs for the API).

## Setup
By default, the NMS looks for API certs in `magma/.cache/test_certs`, and uses `host.docker.internal:9443` as the API host. If you wish to use different API certs and/or a different API host, you can create a `.env` file within `magma/nms/fbcnms-projects/magmalte` and specify them there.
```bash
API_HOST=example.com
API_CERT_FILENAME=/path/to/api_operator_cert.pem
API_PRIVATE_KEY_FILENAME=/path/to/operator_pk.pem
``` 

## Running the NMS
In the `magmalte` directory, start docker containers and create a test user:
```bash
HOST [magma]$ cd nms/fbcnms-projects/magmalte
HOST [magma/nms/fbcnms-projects/magmalte]$ docker-compose up -d
HOST [magma/nms/fbcnms-projects/magmalte]$ ./scripts/dev_setup.sh
```
You may get an error if you run `dev_setup.sh` immediately after `docker-compose up -d`. To resolve this, wait a bit before running the script to let migrations run.

Once you have started the docker containers and created a test user using the `dev_setup.sh` script, go to https://localhost and login with test credentials `admin@magma.test` and `password1234`.

Note: if you want to name a user other than `admin@magma.test`, you can run `setAdminPassword`, like so:
```bash
HOST [magma/nms/fbcnms-projects/magmalte]$ docker-compose run magmalte yarn run setAdminPassword admin@magma.test password1234
```
