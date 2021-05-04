---
id: version-1.5.0-dev_unit_testing
title: Unit Testing
hide_title: true
original_id: dev_unit_testing
---

# Testing Tips

This guide covers tips for quickly validating AGW changes.

## Run all unit tests on the dev VM

In general, all unit testing for AGW is done on the magma dev VM. 
To SSH into the VM, run
```bash
[HOST] cd $MAGMA_ROOT/lte/gateway
[HOST] vagrant up magma
[HOST] vagrant ssh magma
```

To run all existing unit tests, run
```bash
[VM] cd magma/lte/gateway
[VM] make test
```
Note: Running all unit tests can take close to 15 minutes.


### Run Python unit tests on the dev VM

To run only the Python unit tests, run
```bash
[VM] cd magma/lte/gateway
[VM] make test_python
```
The list of services to test are configured in the following files. 
* `orc8r/gateway/python/defs.mk`
* `lte/gateway/python/defs.mk`

### Run C/C++ unit tests on the dev VM

We have several C/C++ services that live in `lte/gateway/c/`. 
To run tests for those services, run

```bash
[VM] cd magma/lte/gateway
[VM] make test_<service_directory_name> # Ex: make test_session_manager
```
