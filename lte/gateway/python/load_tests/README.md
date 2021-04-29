---
id: load_tests
title: Magma AGW services load testing
hide_title: true
---

# Magma VM Load testing
Current testing workflow for VM-only load tests for AGW stateless services
(mobilityd, pipelined, sessiond, MME). Note: At the moment, only mobilityd
and pipelined is supported, sessiond and MME load tets will be introduced soon.

This process makes use of a golang tool `gHZ` that will call a specific gRPC
interface / function and retrieve the results of response time, success or
failure, and so on. The options for configuration and such can be looked
[here](https://ghz.sh/docs/options).

### Gateway VM setup

Spin up and provision the gateway VM, then make and start its services:

1. From `magma/lte/gateway` on the host machine:
   `vagrant up magma && vagrant ssh magma`
2. Now in the gateway VM: `cd $MAGMA_ROOT/lte/gateway && make run`

### Run tests

From `$MAGMA_ROOT/lte/gateway/python/load_tests` on the *magma* VM, run:

* All tests: `make load_test`

This will run every load test script that is defined on
`$MAGMA_ROOT/lte/gateway/python/load_tests/defs.mk`, the pattern to define a
new load test follows `<name_of_load_test_script>.py:<name_of_grpc_function>`.

For example, for mobilityd's `AllocateIPRequest`, there is a
`load_test_mobilityd.py` which contains an `allocate` command
(run as `load_test_mobilityd.py allocate`).

These can also be triggered from the host (laptop), under
`$MAGMA_ROOT/lte/gateway/python/load_tests`, calling `fab load_test`.

### Results

After running the load tests, the results for each of them will be saved under
`/tmp` on *magma* VM, named as `result_<name_of_grpc_function>` as a JSON file,
these can be uploaded to different monitoring tools, an example being
[gHZ-web](https://ghz.sh/docs/web/intro),
which acts as a web server and contains an API to show results over time.

### Notes

- [gHZ reference](https://ghz.sh/)
