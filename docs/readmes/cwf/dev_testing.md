---
id: dev_testing
title: Test CWAG
hide_title: true
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

CWF integration tests use 3 separate VMs listed below.
`cwf/gateway/fabfile.py` can be used to automate all setup work.

### `cwag-dev`

Runs CWAG docker containers and mock core services needed to run the test.
See `cwf/gateway/docker-compose.integ-test.yml` for the complete list of services.

### `cwag-test`

Runs a UE simulator service and all tests.

### `magma-trfserver`

Runs an iperf3 server to drive traffic through CWAG.

#### Entire test suite

To run all setup work and the entire CWF integration test suite, run

```bash
[HOST] fab integ_test
```

Once the above command has been executed, which means that the set-up of the VMs etc. has been
performed, command-line options can be utilized to rerun the tests without redoing the set-up

```bash
[HOST] fab integ_test:provision_vm=False,no_build=True
```

#### Individual tests

The command above can be further modified to run one integration test at a time

```bash
[HOST] fab integ_test:provision_vm=False,no_build=True,skip_unit_tests=True,test_re=<TEST_TO_RUN>
```

where `<TEST_TO_RUN>` is to be replaced by the desired test, e.g. `TestGyReAuth`.

See `fab --display integ_test` for more information.
