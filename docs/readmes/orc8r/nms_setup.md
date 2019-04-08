---
id: nms_setup
title: Setting up the NMS
---
The NMS is the UI for managing, configuring, and monitoring networks. To set up the NMS, we will first need to set up a couple of prerequisites.

## Prerequisites
- Magma must be set up (the NMS needs magma certs for the API)
- Docker must be installed

## Setup
By default, the NMS looks for API certs in `magma/.cache/test_certs`, and uses `192.168.80.10:9443` as the API host. If you wish to use different API certs and/or a different API host, you can create a `.env` file within `magma/nms/fbcnms-projects/magmalte` and specify them there.
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
HOST [magma/nms/fbcnms-projects/magmalte]$ ./scripts/create_test_user.sh
```
You may get an error if you run `create_test_user.sh` immediately after `docker-compose up -d`. To resolve this, wait a bit before running `create_test_user.sh` to let migrations run.

Once you have started the docker containers and created a test user, go to https://localhost and login with test credentials `admin@magma.test` and `password1234`.

Note: if you want to name a user other than `admin@magma.test`, you can run `setAdminPassword`, like so:
```bash
HOST [magma/nms/fbcnms-projects/magmalte]$ docker-compose run magmalte yarn run setAdminPassword admin@magma.test password1234
```
