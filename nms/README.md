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
Some dev workflows and commands live directly in `magma/nms`.
We plan to merge this directory with the main application directory.
Here, `package.json` specifies dependencies necessary for development workflows and testing.
Specifically, triggering unit tests, e2e tests, eslint, and flow tests are done from this directory.

---
## Running Local Dev NMS
> NOTE: This guide is written for development directly using Docker.
> You may need to make adjustments if using MiniKube.

### Pre-requisites
**Docker Settings**
- Ensure you have Docker installed
- Configure Docker to at least 4 GiB allocated memory for Orc8r + NMS

**Magma Orc8r**
- The NMS requires a working connection to the Magma Orc8r to function.
- For local development, bring up the Orc8r docker containers.
- See `magma/orc8r/cloud/docker/docker-compose.yml`

**NMS Metrics and Graphing**
- To ensure that NMS has access to metrics and graphing functionality, bring up the Prometheus and Grafana containers.
- See `magma/orc8r/cloud/docker/docker-compose.metrics.yml`
- This is optional during development

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

---
## NMS Troubleshooting

### First Time Running NMS
We recommend that the first time you bring up the Magma NMS,
you do so without any local changes,
allowing you to troubleshoot your setup.
This should correspond to a recent pull of the `main`/`master` branch from Github.

Here are some example logs from the Magma NMS `magmalte` container from a successful startup:
```
➜  magmalte git:(1-9-nms-orgs-fix) docker-compose logs -f magmalte
Attaching to magmalte_magmalte_1
magmalte_1     | wait-for-it.sh: waiting 30 seconds for postgres:5432
magmalte_1     | wait-for-it.sh: postgres:5432 is available after 0 seconds
magmalte_1     | yarn run v1.22.4
magmalte_1     | $ nodemon scripts/server
magmalte_1     | [nodemon] 2.0.6
magmalte_1     | [nodemon] reading config ./nodemon.json
magmalte_1     | [nodemon] to restart at any time, enter `rs`
magmalte_1     | [nodemon] or send SIGHUP to 48 to restart
magmalte_1     | [nodemon] watching path(s): config/**/* scripts/**/* server/**/* grafana/**/* alerts/**/*
magmalte_1     | [nodemon] watching extensions: js,mjs,json
magmalte_1     | [nodemon] starting `node -r '@fbcnms/babel-register' scripts/server.js`
magmalte_1     | [nodemon] spawning
magmalte_1     | [nodemon] child pid: 60
magmalte_1     | [nodemon] watching 39 files
...
magmalte_1     | 2021-09-11T22:15:40.797Z [scripts/server.js] info: Development server started on port 8081
....
magmalte_1     | webpack built fc3715696396ffc97329 in 36938ms
magmalte_1     | 2021-09-11T22:16:12.510Z [express-middleware/webpackSmartMiddleware.js] warn: Version: webpack 4.46.0
magmalte_1     | Time: 36938ms
magmalte_1     | Built at: 09/11/2021 10:16:12 PM
magmalte_1     |     Asset      Size  Chunks                    Chunk Names
magmalte_1     |  login.js  43.5 KiB   login  [emitted]         login
magmalte_1     |   main.js  5.37 MiB    main  [emitted]  [big]  main
magmalte_1     | master.js  1.03 MiB  master  [emitted]  [big]  master
magmalte_1     | vendor.js  41.7 MiB  vendor  [emitted]  [big]  vendor
magmalte_1     | Entrypoint main [big] = vendor.js main.js
magmalte_1     | Entrypoint login [big] = vendor.js login.js
magmalte_1     | Entrypoint master [big] = vendor.js master.js
...
magmalte_1     | 2021-09-11T22:16:12.511Z [express-middleware/webpackSmartMiddleware.js] info: Compiled with warnings.
```

### Checking NMS Logs
Run the following from `magma/nms/packages/magmalte`
```
docker-compose logs -f magmalte
```

Depending on the issues you run into, you may need to check the logs of other containers, both NMS and Orc8r related.

If the NMS application server crashes, you will likely find the following error log:
```
[nodemon] app crashed - waiting for file changes before starting...
```

### Cannot Connect to NMS Through Browser
If your NMS containers are up and running, it is likely that the NMS app server has crashed.
Check the `magmalte` container logs.

### NMS Loads a Blank Page
This likely corresponds to React errors.
If this is the case, you should be able to see the relevant error logs through your web browser's developer tools.

---
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

---
## Updating NMS for API Changes

### How to re-generate MagmaAPIBindings.js
Run `./build.py --generate` in `${MAGMA_ROOT}/orc8r/cloud/docker`

This re-generates various files, including `MagmaAPIBindings.js` for NMS.

### How to re-generate MagmaAPIBindings.js (old)
1. Place an up-to-date `swagger.yml` in `magma/nms/packages/magmalte`
   You can get the full `swagger.yml` at `{orc8r domain}/swagger/v1/spec`
   (e.g. `https://localhost:9443/swagger/v1/spec`)
2. Run `scripts/generateAPIFromSwagger.sh` from `magma/nms/packages/magmalte`
3. Delete `swagger.yml` afterwards

---
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
