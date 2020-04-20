---
id: version-1.1.0-s1ap_tests
title: S1AP Integration Tests
hide_title: true
original_id: s1ap_tests
---
# S1AP Integration Tests
Current testing workflow for VM-only S1AP integration tests. We cover
gateway-only tests and some general notes.

TODO: Update this document once integration tests with cloud are also supported

Our VM-only tests use 3 Vagrant-managed VMs hosted on the local device (laptop):

- *magma*, i.e. magma-dev or gateway
- *magma_test*, i.e. s1ap_tester
- *magma_trfserver*, i.e. an Iperf server to generate uplink/downlink traffic

## Gateway-only tests

These tests use all 3 VMs listed above. The *magma_test* VM abstracts away the
UE and eNodeB, the *magma_trfserver* emulates the Internet, while the *magma* VM
acts as the gateway between *magma_test* and *magma_trfserver*.

### Gateway VM setup

Spin up and provision the gateway VM, then make and start its services:

1. From `magma/lte/gateway` on the host machine: `vagrant up magma && vagrant ssh magma`
1. Now in the gateway VM: `cd $MAGMA_ROOT/lte/gateway && make run`

### Test VM setup

Spin up and provision the s1ap tester's VM, make, then make in the integ_tests directory.

1. From `magma/lte/gateway` on the host machine: `vagrant up magma_test && vagrant ssh magma_test`
1. Now in the *magma_test* VM:
    1. `cd $MAGMA_ROOT/lte/gateway/python && make`
    1. `cd integ_tests && make`

### Run tests

From `$MAGMA_ROOT/lte/gateway/python/integ_tests` on the *magma_test* VM, run
either individual tests or the full suite of tests. A safe, non-flaky test to
run is `s1aptests/test_attach_detach.py`.

* Individual test(s): `make integ_test TESTS=<test(s)_to_run>`
* All tests: `make integ_test`

**Note**: The traffic tests will fail as traffic server is not running in this
setup. Look at the section below on running traffic tests.

### Running uplink/downlink traffic tests

1. On the *magma* VM, run, `disable-tcp-checksumming`

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

 `cd /etc/systemd/system`

 comment out the line `PartOf=magma@mme.service` from the following files
 (you will need sudo privileges):

 magma@mobilityd.service, magma@pipelined.service, magma@sessiond.service and
sctpd.service

 `sudo systemctl daemon-reload`

1. In `/etc/magma/mme.yml`, set `use_stateless` to true
1. Clean up all the state in redis: `redis-cli -p 6380 FLUSHALL`. This might
throw a "Could not connect" error if magma@redis service is not running. Start
the redis service with `sudo service magma@redis start` and then try again.
1. `cd $MAGMA_ROOT/lte/gateway; make restart`

On test VM:
1. Basic attach/detach test where MME is restarted mid-way:

  `make integ_test TESTS=s1aptests/test_attach_detach_with_mme_restart.py`

1. Attach with uplink UDP traffic, where MME is restarted while UDP traffic is
flowing:

 `make integ_test TESTS=s1aptests/test_attach_ul_udp_data_with_mme_restart.py`

 , make sure traffic server VM is running (as described in traffic tests above) and
TCP checksum is disabled on all VMs.

#### Stateless Mobilityd
This section describes how to test whether Mobilityd service is persisting state to Redis.

On gateway VM:

1. Disable MME from restarting when Mobilityd restarts.

 comment out the line `PartOf=magma@mobilityd.service` from the MME system
service file `/etc/systemd/system/magma@mme.service` (you will need sudo privileges):

 `sudo systemctl daemon-reload`

1. In `/etc/magma/mobilityd.yml`, set `persist_to_redis` to `true`

1. Clean up all the state in redis: `redis-cli -p 6380 FLUSHALL`. This might
throw a "Could not connect" error if magma@redis service is not running. Start
the redis service with `sudo service magma@redis start` and then try again.

1. `cd $MAGMA_ROOT/lte/gateway; make restart`

On test VM:
1. `cd $MAGMA_ROOT/lte/gateway/python && make`
1. `cd integ_tests && make`
1. Basic attach/detach test where Mobilityd is restarted mid-way:

 `make integ_test TESTS=s1aptests/test_attach_detach_with_mobilityd_restart.py`

1. Test IP blocks are maintained across service restart

 `make integ_test TESTS=s1aptests/test_attach_detach_multiple_ip_blocks_mobilityd_restart.py`

#### Stateless Pipelined
This section describes how to test whether Pipelined service is persisting state to Redis.

On gateway VM:

1. Disable MME from restarting when Pipelined restarts.

 comment out the line `PartOf=magma@pipelined.service` from the MME system
service file `/etc/systemd/system/magma@mme.service` (you will need sudo privileges):

 `sudo systemctl daemon-reload`

1. In `/etc/magma/pipelined.yml`, set `clean_restart` to `false`

1. Clean up all the state in redis: `redis-cli -p 6380 FLUSHALL`. This might
throw a "Could not connect" error if magma@redis service is not running. Start
the redis service with `sudo service magma@redis start` and then try again.

1. `cd $MAGMA_ROOT/lte/gateway; make restart`

On test VM:
1. `cd $MAGMA_ROOT/lte/gateway/python && make`
1. `cd integ_tests && make`
1. UDP traffic test where Pipelined is restarted mid-way:

 `make integ_test TESTS=s1aptests/test_attach_ul_udp_data_with_pipelined_restart.py`


### Testing stateless gateway with all services

To test the gateway with all services being stateless,

1. On gateway VM, follow steps 1 and 2 for each of Stateless MME, Mobilityd and Pipelined as
listed above

1. On test VM, you can run any of the test cases for individual service restarts
listed above. Further, you can test attach with uplink UDP traffic, where
multiple services are restarted while UDP traffic is flowing:

 `make integ_test TESTS=s1aptests/test_attach_ul_udp_data_with_multiple_service_restart.py`

 , make sure traffic server VM is running (as described in traffic tests above) and
TCP checksum is disabled on all VMs.

## Notes

- Restart the *magma* VM (`vagrant reload magma`) on an assertion error involving `ENB_S1_SETUP_RESP.` This is a known issue.
- See *[Bindings for Magma's REST API](https://fb.quip.com/4tmUAtlox4Oy)* for notes on the Python bindings for our REST API generated by [swagger-codegen](https://github.com/swagger-api/swagger-codegen).
- It may be cleaner to set the host using the [configuration class](https://github.com/swagger-api/swagger-codegen/blob/master/samples/client/petstore/python/petstore_api/configuration.py). This is also where we can set SSL options.
