---
id: s1ap_tests
title: S1AP Integration Tests
hide_title: true
---
# S1AP Integration Tests

Current testing workflow for VM-only S1AP integration tests. We cover
gateway-only tests and some general notes.

CRITICAL NOTE: S1AP integration tests are supposed to be run in a headless mode,
i.e., AGW should not be connected to Orc8r. This is quite critical as S1AP tester
makes local configurations to Magma AGW in accordance with the testing scenario (e.g.,
subscriber, APN, policy rules, etc.). If an Orc8r is connected, these configurations
would be overwritten periodically and also lead to restart of services, both of which will
interfere with the test scenario.

<!-- TODO: Update this document once integration tests with cloud are also supported -->

Our VM-only tests use 3 Vagrant-managed VMs hosted on the local device (laptop):

- *magma* (a.k.a. magma-dev) or *magma_deb*, both of which act as a gateway, with the difference
being the way magma is installed on them (see [below](#gateway-vm-setup) for more details)
- *magma_test*, i.e. s1ap_tester
- *magma_trfserver*, i.e. an Iperf server to generate uplink/downlink traffic

## Gateway-only tests

These tests use all 3 VMs listed above. The *magma_test* VM abstracts away the
UE and eNodeB, the *magma_trfserver* emulates the Internet, while the *magma*/*magma_deb* VM
acts as the gateway between *magma_test* and *magma_trfserver*.

### Gateway VM setup

There are two options for setting up the gateway VM, with the difference being the magma installation
method: it can either be installed by directly compiling the services or from debian packages, the
latter being the method by which magma is usually deployed. For everyday development, directly
compiling the services is recommended, while the debian installation is useful for testing packages
before release.

> **Warning**: These two VMs use the same network configuration, so one must only run one of them
at a time.

#### Directly compiling services

Spin up and provision the gateway VM. From `magma/lte/gateway` on the host machine run
`vagrant up magma && vagrant ssh magma`.

Now build and run the services with bazel.

`cd $MAGMA_ROOT && bazel/scripts/build_and_run_bazelified_agw.sh`

#### Debian installation

Spin up the *magma_deb* VM. The services start automatically:

1. From `magma/lte/gateway` on the host machine: `vagrant up magma_deb && vagrant ssh magma_deb`
1. To check the services are running, run `systemctl list-units --type=service magma@*`

> **Warning**: During provisioning, the latest magma gateway debian build from the magma artifactory
is installed. That is, the deployed gateway might not match your local repository state.

### Test VM setup

Spin up and provision the s1ap tester's VM. From `magma/lte/gateway` on the host machine: `vagrant up magma_test`

### Run tests

From `$MAGMA_ROOT/lte/gateway/python/integ_tests` on the *magma_test* VM, run
either individual tests or the full suite of tests. A safe, non-flaky test to
run is `s1aptests/test_attach_detach.py`.

- Individual test(s): `make selected_tests TESTS="<test(s)_to_run_space_separated>"`
- All Sanity tests: `make integ_test`
- All Non-Sanity tests: `make nonsanity`
- Minimal set of tests to be executed before committing changes to magma repository: `make precommit`
- Run with *-i* flag to enable continuous test runs (ignoring the failing test(s), if any):\
 `make -i precommit` or `make -i integ_test`
- Set *enable-flaky-retry=true* to re-run the failing test(s) to identify flaky behavior:\
`make precommit enable-flaky-retry=true` or `make integ_test enable-flaky-retry=true`

> **Note**: The traffic tests will fail as the traffic server is not running in this
setup. Look at the section below on running traffic tests.

#### Run tests with Bazel

Alternatively, the integration tests can be run with Bazel.

> **Note**: Running the integration tests with the `run_integ_tests.sh` script results in slightly longer test runtimes. Many features of the script may nonetheless be convenient for local development.

From inside the repository on the *magma_test* VM, execute `$MAGMA_ROOT/bazel/scripts/run_integ_tests.sh` to run all Sanity tests.

It is also possible to run individual tests, run the Non-Sanity tests, list tests by classification, re-run failed tests and to run the set-up and tear-down individually for all test classifications.

To see the full list of features available, please use the help flag: `$MAGMA_ROOT/bazel/scripts/run_integ_tests.sh --help`.

### Running uplink/downlink traffic tests

1. On the *magma* or *magma_deb* VM, run, `disable-tcp-checksumming`

1. On the *magma_test* VM, `disable-tcp-checksumming`

1. Start the traffic server VM from the host, `vagrant up magma_trfserver && vagrant ssh magma_trfserver`

1. From *magma_trfserver* VM, run `disable-tcp-checksumming && trfgen-server`

Running `make integ_test` in *magma_test* VM should succeed now.

## Testing stateless Access Gateway

The Access Gateway by default runs a set of stateful services, which means that
whenever the services are restarted, all previous state of UEs and eNodeBs, and
they need to reconnect and re-register. Alternatively, we can switch the Access
Gateway to be stateless, as shown below, so that all UE state is preserved
across service restarts.

Note that this is a feature in development, so some tests from the integ_test
suite may not pass.

All the tests below assume you have completed the Gateway Setup and Test VM
Setup described above.

### Testing one stateless service

#### Stateless MME

This section describes how to test whether MME service is persisting state to Redis.

On gateway VM:

1. Disable Pipelined, Mobilityd, Sctpd and Sessiond from restarting when MME
restarts.
    1. `cd /etc/systemd/system`
    1. comment out the line `PartOf=magma@mme.service` from the following files
 (you will need sudo privileges):\
 magma@mobilityd.service, magma@pipelined.service, magma@sessiond.service and
sctpd.service
    1. `sudo systemctl daemon-reload`

1. In `/etc/magma/mme.yml`, set `use_stateless` to true
1. Clean up all the state in redis: `redis-cli -p 6380 FLUSHALL`. This might
throw a "Could not connect" error if magma@redis service is not running. Start
the redis service with `sudo service magma@redis start` and then try again.
1. `magma-restart`

On test VM:

1. Basic attach/detach test where MME is restarted mid-way:\
  `make selected_tests TESTS=s1aptests/test_attach_detach_with_mme_restart.py`

1. Attach with uplink UDP traffic, where MME is restarted while UDP traffic is
flowing:\
 `make selected_tests TESTS=s1aptests/test_attach_ul_udp_data_with_mme_restart.py`\
 , make sure traffic server VM is running (as described in traffic tests above) and
TCP checksum is disabled on all VMs.

#### Stateless Mobilityd

This section describes how to test whether Mobilityd service is persisting state to Redis.

On gateway VM:

1. Disable MME from restarting when Mobilityd restarts.
    1. comment out the line `PartOf=magma@mobilityd.service` from the MME system
service file `/etc/systemd/system/magma@mme.service` (you will need sudo privileges)
    1. `sudo systemctl daemon-reload`

1. Clean up all the state in redis: `redis-cli -p 6380 FLUSHALL`. This might
throw a "Could not connect" error if magma@redis service is not running. Start
the redis service with `sudo service magma@redis start` and then try again.

1. `magma-restart`

On test VM:

1. `cd $MAGMA_ROOT/lte/gateway/python && make`
1. `cd integ_tests && make`
1. Basic attach/detach test where Mobilityd is restarted mid-way:\
 `make selected_tests TESTS=s1aptests/test_attach_detach_with_mobilityd_restart.py`

1. Test IP blocks are maintained across service restart\
 `make selected_tests TESTS=s1aptests/test_attach_detach_multiple_ip_blocks_mobilityd_restart.py`

#### Stateless Pipelined

This section describes how to test whether Pipelined service is persisting state to Redis.

On gateway VM:

1. Disable MME from restarting when Pipelined restarts.
    1. comment out the line `PartOf=magma@pipelined.service` from the MME system
service file `/etc/systemd/system/magma@mme.service` (you will need sudo privileges)
    1. `sudo systemctl daemon-reload`

1. In `/etc/magma/pipelined.yml`, set `clean_restart` to `false`

1. Clean up all the state in redis: `redis-cli -p 6380 FLUSHALL`. This might
throw a "Could not connect" error if magma@redis service is not running. Start
the redis service with `sudo service magma@redis start` and then try again.

1. `magma-restart`

On test VM:

1. `cd $MAGMA_ROOT/lte/gateway/python && make`
1. `cd integ_tests && make`
1. UDP traffic test where Pipelined is restarted mid-way:
 `make selected_tests TESTS=s1aptests/test_attach_ul_udp_data_with_pipelined_restart.py`

### Testing stateless gateway with all services

To test the gateway with all services being stateless,

1. On gateway VM, follow steps 1 and 2 for each of Stateless MME, Mobilityd and Pipelined as
listed above

1. On test VM, you can run any of the test cases for individual service restarts
listed above. Further, you can test attach with uplink UDP traffic, where
multiple services are restarted while UDP traffic is flowing:\
 `make selected_tests TESTS=s1aptests/test_attach_ul_udp_data_with_multiple_service_restart.py`\
 , make sure traffic server VM is running (as described in traffic tests above) and
TCP checksum is disabled on all VMs.

## Notes

- Restart the *magma* VM (`vagrant reload magma`) on an assertion error involving `ENB_S1_SETUP_RESP.` This is a known issue.
- See *[Bindings for Magma's REST API](https://fb.quip.com/4tmUAtlox4Oy)* for notes on the Python bindings for our REST API generated by [swagger-codegen](https://github.com/swagger-api/swagger-codegen).
- It may be cleaner to set the host using the [configuration class](https://github.com/swagger-api/swagger-codegen/blob/master/samples/client/petstore/python/petstore_api/configuration.py). This is also where we can set SSL options.
