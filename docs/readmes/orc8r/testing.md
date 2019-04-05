---
id: testing
sidebar_label: Testing
title: Testing in Orchestrator
hide_title: true
---
# Testing in Orchestrator
### Unit Tests
One easy way to test is to run unit tests. This can be done by running
```
magma/orc8r/cloud: make test
```

### Run the services and check their health
Unit tests are great for checking small logic chunks, 
but another way to test is to run the services and check their status.
The services can be built and started by running
```
magma/orc8r/cloud: make run
```
or 
```
magma/orc8r/cloud: make restart
```
to restart the services without rebuilding.

The state of these services can be checked by running
```
sudo service magma@* status
```
for all services, and
```
sudo service magma@SERVICE_NAME status
```
for a specific service.

Run 
```
sudo journalctl -u magma@SERVICE_NAME
```
to see logs relating to a specific service.
