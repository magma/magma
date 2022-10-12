---
id: version-1.8.0-dev_testing
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

where `<TEST_TO_RUN>` is to be replaced by the desired test, e.g. `TestGyReAuth`. Run
`fab --display integ_test` to see more available options.

Note that running a test can be further expedited by running it directly on the CWAG test VM using
`gotestsum`, which echoes how it is run by the `fab` command. In particular, the above command can
be run with `run_tests=False` to do the required set-up

```bash
[HOST] fab integ_test:run_tests=False,provision_vm=False,no_build=True,skip_unit_tests=True,test_re=<TEST_TO_RUN>
```

before logging into the CWAG test VM, navigating to the gateway folder and running the test directly

```bash
[HOST] vagrant ssh cwag_test
[VM] cd magma/cwf/gateway
[VM] gotestsum --format=standard-verbose --packages='./...' -- -test.short -timeout 50m -count 1 -tags=all -run=<TEST_TO_RUN>
```

This command can be used to execute the test multiple times in a way that is faster than can be done
with the `fab` command, which may be helpful during development.
