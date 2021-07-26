---
id: dev_unit_testing
title: Test AGW
hide_title: true
---

# Test Access Gateway

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

### Test Python AGW services

To run only the Python unit tests, run

```bash
[VM] cd magma/lte/gateway
[VM] make test_python
```

The list of services to test are configured in the following files. 
* `orc8r/gateway/python/defs.mk`
* `lte/gateway/python/defs.mk`

### Test C/C++ AGW services

We have several C/C++ services that live in `lte/gateway/c/`. 
To run tests for those services, run

```bash
[VM] cd magma/lte/gateway
[VM] make test_<service_directory_name> # Ex: make test_session_manager
```

### Test Go AGW services

We have several Go implementations of AGW services that live in `orc8r/gateway/go`.
To test any changes, run

```bash
[VM] cd magma/orc8r/gateway/go
[VM] go test ./...
```

## Format AGW

### Format Python

Docker is required for running the steps below. 
To use the `--diff` flag, the script will have to be on your host machine where the Magma repository lives.
Refer to the script at `lte/gateway/python/precommit.py` for all available commands, but the main ones are as follows.

```bash
cd $MAGMA/lte/gateway/python
# build the base image
./precommit.py --build

# un the flake8 linter by specifying paths
./precommit.py --lint -p PATH1 PATH2
# run the flake8 linter on all modified files in HEAD vs master
# this command can only be run on your host
./precommit.py --lint --diff

# run all available formatters by specifying paths
./precommit.py --format -p PATH1 PATH2
# run all available formatters on all modified files in HEAD vs master
# this command can only be run on your host
./precommit.py --format --diff
```

### Format C/C++

To run formatting for each C/C++ services, run

```bash
[VM] cd magma/lte/gateway
[VM] make format_<service_directory_name> # Ex: make format_session_manager
```