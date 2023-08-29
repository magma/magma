---
id: load_tests
title: Magma AGW services load testing
hide_title: true
---

# Magma VM Load testing

Current testing workflow for VM-only load tests for AGW stateless services
(directoryd, mobilityd, pipelined, policydb, sessiond, subscriberdb).

This process makes use of a golang tool `gHZ` that will call a specific gRPC
interface / function and retrieve the results of response time, success or
failure, and so on. The options for configuration and such can be looked
[here](https://ghz.sh/docs/options).

### Gateway VM setup

Spin up and provision the gateway VM. From `magma/lte/gateway` on the host machine run `vagrant up magma && vagrant ssh magma`.

Next build and start the magma services. See the [integ test documentation](https://magma.github.io/magma/docs/next/lte/s1ap_tests#s1ap-integration-tests) for details.

### Run tests

On the gateway VM, run `$MAGMA_ROOT/bazel/scripts/run_load_tests.sh` from anywhere inside the magma folder.

This will run every load test that is defined inside the `LOAD_TEST_LIST` of the script.

For example, for mobilityd's `AllocateIPRequest`, there is a
`load_test_mobilityd.py` which contains an `allocate` command
(run as `load_test_mobilityd.py allocate`).

### Results

After running the load tests, the results for each of them will be saved under
`/var/tmp` on the gateway VM, named as `result_<name_of_grpc_function>` as a JSON file.
These can be uploaded to different monitoring tools, an example being
[gHZ-web](https://ghz.sh/docs/web/intro),
which acts as a web server and contains an API to show results over time.

### Notes

- [gHZ reference](https://ghz.sh/)
