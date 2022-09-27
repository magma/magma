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
method: it can either be installed via `Bazel` or from debian packages, the latter being the
method by which magma is usually deployed. For everyday development, the `Bazel` installation is
recommended, while the debian installation is useful for testing packages before release.

> **Warning**: These two VMs use the same network configuration, so one must only run one of them
at a time.

#### Make installation

Spin up and provision the gateway VM, then build and start its services:

1. From `magma/lte/gateway` on the host machine: `vagrant up magma && vagrant ssh magma`
1. Now in the gateway VM (TODO: Make this more comfortable):
   1. 
   ```bash
    vagrant up magma
    vagrant ssh magma
    sudo apt-get update
    sudo DEBIAN_FRONTEND=noninteractive apt-get -y dist-upgrade
    cd ~/magma; bazel/scripts/remote_cache_bazelrc_setup.sh magma-dev-vm false
    sudo sed -i "s@#precedence ::ffff:0:0/96  100@precedence ::ffff:0:0/96  100@" /etc/gai.conf;
    cd ~/magma; bazel build --profile=bazel_profile_lte_integ_tests `bazel query "kind(.*_binary, //orc8r/... union //lte/... union //feg/...)"`;
    sudo sed -i "s@precedence ::ffff:0:0/96  100@#precedence ::ffff:0:0/96  100@" /etc/gai.conf;
    cd ~/magma; bazel/scripts/link_scripts_for_bazel_integ_tests.sh;
    sudo cp $MAGMA_ROOT/lte/gateway/deploy/roles/magma/files/systemd_bazel/* /etc/systemd/system/
    sudo systemctl daemon-reload
    sudo service magma@magmad start
   ```

#### Debian installation

Spin up the *magma_deb* VM. The services start automatically:

1. From `magma/lte/gateway` on the host machine: `vagrant up magma_deb && vagrant ssh magma_deb`
2. To check the services are running, run `systemctl list-units --type=service magma@*`

> **Warning**: During provisioning, the latest magma gateway debian build from the magma artifactory
is installed. That is, the deployed gateway might not match your local repository state.

### Test VM setup

Spin up and provision the s1ap tester's VM.

1. From `magma/lte/gateway` on the host machine: `vagrant up magma_test && vagrant ssh magma_test`

### Run tests

From `$MAGMA_ROOT/bazel/scripts/` on the *magma_test* VM, run
either individual tests or the full suite of tests. A safe, non-flaky test to
run is `s1aptests/test_attach_detach.py`.

- Individual test(s): `run_integ_tests.sh //lte/gateway/python/integ_tests/s1aptests:test_attach_detach`
- All Sanity tests: `run_integ_tests.sh`
- All Non-Sanity tests: `run_integ_tests.sh --nonsanity`
- Minimal set of tests to be executed before committing changes to magma repository: `run_integ_tests.sh --precommit`
- All tests are run even if some of fail during the process. At the end a summary will be provided.
- You can enable flaky-retry behavior to re-run the failing test(s) to identify flaky behavior:\
`run_integ_tests.sh --retry-on-failure` or `run_integ_tests.sh --retry-on-failure --precommit --retry-attempts 5`
- For more information about the usage and features of this script execute the following command:\
`run_integ_tests.sh --help`

**Note**: The traffic tests will fail as traffic server is not running in this
setup. Look at the section below on running traffic tests.

### Running uplink/downlink traffic tests

1. On the *magma* or *magma_deb* VM, run, `disable-tcp-checksumming`

1. On the *magma_test* VM, `disable-tcp-checksumming`

1. Start the traffic server VM from the host, `vagrant up magma_trfserver && vagrant ssh magma_trfserver`

1. From *magma_trfserver* VM, run `disable-tcp-checksumming && trfgen-server`

Running `$MAGMA_ROOT/bazel/scripts/run_integ_tests.sh` in *magma_test* VM should succeed now.

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
  `$MAGMA_ROOT/bazel/scripts/run_integ_tests.sh //lte/gateway/python/integ_tests/s1aptests:test_attach_detach_with_mme_restart`

1. Attach with uplink UDP traffic, where MME is restarted while UDP traffic is
flowing:\
 `$MAGMA_ROOT/bazel/scripts/run_integ_tests.sh //lte/gateway/python/integ_tests/s1aptests:test_attach_ul_udp_data_with_mme_restart`\
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

1. Basic attach/detach test where Mobilityd is restarted mid-way:\
 `$MAGMA_ROOT/bazel/scripts/run_integ_tests.sh //lte/gateway/python/integ_tests/s1aptests:test_attach_detach_with_mobilityd_restart`

1. Test IP blocks are maintained across service restart\
 `$MAGMA_ROOT/bazel/scripts/run_integ_tests.sh //lte/gateway/python/integ_tests/s1aptests:test_attach_detach_multiple_ip_blocks_mobilityd_restart`

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

1. UDP traffic test where Pipelined is restarted mid-way:
 `$MAGMA_ROOT/bazel/scripts/run_integ_tests.sh //lte/gateway/python/integ_tests/s1aptests:test_attach_ul_udp_data_with_pipelined_restart`

### Testing stateless gateway with all services

To test the gateway with all services being stateless,

1. On gateway VM, follow steps 1 and 2 for each of Stateless MME, Mobilityd and Pipelined as
listed above

1. On test VM, you can run any of the test cases for individual service restarts
listed above. Further, you can test attach with uplink UDP traffic, where
multiple services are restarted while UDP traffic is flowing:\
 `$MAGMA_ROOT/bazel/scripts/run_integ_tests.sh //lte/gateway/python/integ_tests/s1aptests:test_attach_ul_udp_data_with_multiple_service_restart`\
 , make sure traffic server VM is running (as described in traffic tests above) and
TCP checksum is disabled on all VMs.

## Notes

- Restart the *magma* VM (`vagrant reload magma`) on an assertion error involving `ENB_S1_SETUP_RESP.` This is a known issue.
- See *[Bindings for Magma's REST API](https://fb.quip.com/4tmUAtlox4Oy)* for notes on the Python bindings for our REST API generated by [swagger-codegen](https://github.com/swagger-api/swagger-codegen).
- It may be cleaner to set the host using the [configuration class](https://github.com/swagger-api/swagger-codegen/blob/master/samples/client/petstore/python/petstore_api/configuration.py). This is also where we can set SSL options.
