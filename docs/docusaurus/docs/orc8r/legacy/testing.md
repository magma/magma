---
id: version-1.0.0-testing
sidebar_label: Testing
title: Testing in Orchestrator
hide_title: true
original_id: testing
---
# Testing in Orchestrator
### Unit Tests
One easy way to test is to run unit tests. This can be done by running..
```
HOST [magma/orc8r/cloud/docker]$ ./build.py --tests
```

### Run the services and check their health
Unit tests are great for checking small logic chunks, 
but another way to test is to run the services and check their status.
The services can be built and started by running
```
docker compose up -d
```

The state of the containers can be checked by running
```
docker-compose ps
docker-compose logs -f
```