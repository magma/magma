---
id: version-1.7.0-dev_testing
title: Test CWAG
hide_title: true
original_id: dev_testing
---

# Test Carrier Wifi Access Gateway

This guide covers tips for quickly validating Carrier Wifi Access Gateway changes.

## Dev environment

In general, all unit testing for CWAG is done on the CWAG test VM.

To SSH into the VM, run

```bash
[HOST] cd $MAGMA_ROOT/cwf/gateway
[HOST] vagrant up cwag_test
[HOST] vagrant ssh cwag_test
```

The commands shown below should be run inside the test VM unless specified otherwise.

## Format and verify build

To run all existing unit tests, run

```bash
[VM] make -C ${MAGMA_ROOT}/cwf/gateway precommit
[VM] make -C ${MAGMA_ROOT}/cwf/gateway/integ_tests precommit
```

## Run unit tests

To run all existing unit tests, run

```bash
[VM] make -C ${MAGMA_ROOT}/cwf/gateway test
```

## Run integration tests

### Prerequisite

You need to install a vagrant plugin:

```bash
vagrant plugin install vagrant-disksize
```

### Test setup

CWF integration tests use 3 separate VMs listed below.
`cwf/gateway/fabfile.py` can be used to automate all setup work.

### `cwag-dev`

Runs CWAG docker containers and mock core services needed to run the test.
See `cwf/gateway/docker-compose.integ-test.yml` for the complete list of services.

### `cwag-test`

Runs a UE simulator service and all tests.

### `magma-trfserver`

Runs an iperf3 server to drive traffic through CWAG.

To run all setup work and the entire CWF integration test suite, run

```bash
[HOST] fab integ_test
```

See `fab --display integ_test` for more information.
