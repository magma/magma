# Magma Network Management System - Development
The Magma NMS provides an enterprise grade GUI for provisioning and operating magma based networks.

> NOTE:
>
> This document is written to help with development of the Magma NMS
>
> See the [**docs**](https://docs.magmacore.org/docs/next/nms/overview) for a user-focused guide.

## Project Layout

There are two main directories, and `package.json` files.

### Main Application Directory
In `magma/nms/packages/magmalte`, the main NMS application lives.
Here, `package.json` specifies NMS application dependencies.

### Dev Workflows Directory
In `magma/nms`, this is to be merged with the main application directory.
Here, `package.json` specifies dependencies necessary for development workflows and testing.
Specifically, triggering unit tests, e2e tests, eslint, and flow tests are done from this directory.

## Running Local Dev NMS
> NOTE: This guide is written for development directly using Docker.
> You may need to make adjustments if using MiniKube.

### Pre-requisites
The NMS requires a working connection to the Magma Orc8r to function.
For local development, bring up the Orc8r docker containers.
See `magma/orc8r/cloud/docker/docker-compose.yml`
To ensure that NMS has access to metrics and graphing functionality, bring up the Prometheus and Grafana containers.
See `magma/orc8r/cloud/docker/docker-compose.metrics.yml`

### Build NMS Docker Image
Run the following command from `magma/nms/packages/magmalte`
```
docker-compose build magmalte
```

### Run NMS Containers
Run the following command from `magma/nms/packages/magmalte`
```
docker-compose up -d
```

## Checking NMS Logs
Run the following from `magma/nms/packages/magmalte`
```
docker-compose logs -f magmalte
```
Depending on the issues you run into, you may need to check the logs of other containers, both NMS and Orc8r related.

## Testing

### Install Dependencies
Running eslint and flow tests requires installing dependencies.

Install node and npm if you haven't already

Install yarn if you don't already have it:
```
npm install --global yarn
```

Run the following command from `magma/nms` to install these dependencies:
```
yarn
```

### Eslint
Run from `magma/nms`
```
yarn run eslint ./
```

### Flow
Run from `magma/nms`
```
flow
```

### Unit Tests
Run from `magma/nms`
```
yarn run test
```

### Testing Coverage

Run `yarn test --coverage`

## Updating NMS for API Changes

### How to re-generate MagmaAPIBindings.js
Run `./build.py -g` in `magma/orc8r/cloud`

This re-generates various files, including `MagmaAPIBindings.js` for NMS.

### How to re-generate MagmaAPIBindings.js (old)
1. Place an up-to-date `swagger.yml` in `magma/nms/packages/magmalte`
   You can get the full `swagger.yml` at `{orc8r domain}/swagger/v1/spec`
   (e.g. `https://localhost:9443/swagger/v1/spec`)
2. Run `scripts/generateAPIFromSwagger.sh` from `magma/nms/packages/magmalte`
3. Delete `swagger.yml` afterwards


## Accessing Local Dev NMS

### Multitenancy and Organizations

Multitenancy is supported in the Magma NMS. Each tenant is called an "organization".
Each organization owns a subset of the networks provisioned on Orchestrator, and the special `master` organization administrates organizations in the system.

Users in organizations log into the NMS using a subdomain that matches their organization name.
For example, users of a *NewOrg* organization in the NMS would access the NMS using http://<NEWORG>.localhost:8081/nms

Note that any NMS user can only access the organization it was created under.

### First-time Setup
When you deploy the NMS for the first time, you'll need to create a user that has access to the master organization.

Run the following command from `magma/nms/packages/magmalte` and make sure to substitute `ADMIN_USER_EMAIL` and `ADMIN_USER_PASSWORD` with your desired email and password.
```
docker-compose exec magmalte yarn setAdminPassword master ADMIN_USER_EMAIL ADMIN_USER_PASSWORD
```

Access the `master` (http://master.localhost:8081/master) portal to create your first organization.
Create a new super user for that organization, and then you can login as that user for your new organization.

For example, if you created an organization called `magma-test`, you can access the NMS at http://magma-test.localhost:8081/nms

### First-time Setup (Fast)

Run the following from `magma/nms/packages/magmalte` to create two users, one for `master` organization, and another for `magma-test`.
The username and password for both will be `admin@magma.test` and `password1234`
```
./scripts/dev_setup.sh
```

### Master Portal
The master portal allows management of organizations.

http://master.localhost:8081/master

### Admin Portal
The admin portal allows management for users of an organization.

For a `magma-test` organization, you would access the admin portal at the following:
http://magma-test.localhost:8081/admin

### NMS Portal
Regular users will do their regular network management through the NMS portal.

For a `magma-test` organization, you would access the NMS at the following:
http://magma-test.localhost:8081/nms
