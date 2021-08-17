---
id: hil_tests
title: Hardware In Loop Testing
hide_title: true
---
# Hardware In Loop Tests

Current testing workflow for HIL testing is using Spirent test center to emulate eNodeB, UE and Network host to run scale and performance tests. We cover
gateway-only tests.

HIL tests can be run with different gateway. However for now the automated runs are using physical Protectli box. (https://protectli.com/vault-4-port)

### Lab Setup


Spirent test emulation hardware is hosted in FB lab emulating eNODEB, UE and traffic host elements. gateway are also hosted in same lab. All tests are
executed in worker node in FB lab. Reports and logs are pushed out to aws S3 for debug and analysis.

### Run tests

To setup HIL worker follow instruction on (https://github.com/fbcinternal/ens_magma/tree/master/spirent_automation)

Current Test categories supported are:
1. Sanity (nightly run time - 30 minutes) update badge with latest results on magma main README
1. Performace (nightly run time - 12hrs)
1. Feature tests at scale - (nightly run time - 6 hrs)
1. Availability - Every day for 12hrs period

### HIL SANITY TEST CASES
1. Verify 12 eNodeBs can connect to a Magma Access Gateway
1. Verify 200 UEs at 5 UE/sec can connect to a Magma Access Gateway
1. Verify 400 UEs at 5 UE/sec can connect to a Magma Access Gateway
1. Verify 600 UEs at 5 UE/sec can connect to a Magma Access Gateway
1. Verify 200 UEs across 12 eNodeBs with 2M data per UE
1. Verify 400 UEs across 12 eNodeBs with 500k data per UE
1. Verify 600 UEs across 12 eNodeBs with 500k data per UE
1. Verify 30 UEs across 12 eNodeBs with 500K data changing state from active-idle-active-idle

