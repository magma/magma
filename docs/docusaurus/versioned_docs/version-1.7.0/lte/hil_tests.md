---
id: version-1.7.0-HIL_AGW_tests
title: Hardware In Loop Testing
hide_title: true
original_id: HIL_AGW_tests
---

# Hardware In Loop Tests

Current testing workflow for HIL testing is using Spirent test center to emulate eNodeB, UE and Network host to run scale and performance tests. HIL is focusing on Magma AGW.
HIL tests can be run with virtual or physical gateway. However, for now the automated runs are using physical [SUT-HW](https://protectli.com/vault-4-port) box.

![HIL AGW Last Pass](http://ens-spirent-test-summary.com.s3-us-west-1.amazonaws.com/sanity/hilsanitypass.svg)

## Dashboard

All test results are available to compare on [dashboard](http://automation.fbmagma.ninja/). Please use `username:magma` and `password:magma`.
We can retrieve log and Grafana metrics for each run by clicking on test run result table.

## Lab Setup

Spirent test emulation hardware is hosted in FB lab emulating eNODEB, UE and traffic host elements. Magma AGW is also hosted in same lab. All tests are
executed in worker node in FB lab. Reports and logs are pushed out to AWS S3 for debug and analysis.

### Run tests

To setup HIL worker follow [instruction](https://github.com/fbcinternal/ens_magma/tree/master/spirent_automation)
Current Test categories supported are:

1. Sanity (every new build, run time - 30 minutes) updates badge with latest result on magma main README
1. Performance (nightly, run time - 12hrs)
1. Feature tests at scale - (nightly, run time - 6hrs)
1. Availability - Every day, for 12hrs

### HIL SANITY TEST CASES

1. Verify 12 eNodeBs can connect to a Magma Access Gateway
1. Verify 200 UEs at 5 UE/sec can connect to a Magma Access Gateway
1. Verify 400 UEs at 5 UE/sec can connect to a Magma Access Gateway
1. Verify 600 UEs at 5 UE/sec can connect to a Magma Access Gateway
1. Verify 200 UEs across 12 eNodeBs with 2M data per UE
1. Verify 400 UEs across 12 eNodeBs with 500k data per UE
1. Verify 600 UEs across 12 eNodeBs with 500k data per UE
1. Verify 30 UEs across 12 eNodeBs with 500K data changing state from active-idle-active-idle

### Notification

All test suites run send notification to slack channel which is used as alerting mechanism.
please join slack [channel](https://magmacore.slack.com/archives/C02164DSGPM) for regular update
