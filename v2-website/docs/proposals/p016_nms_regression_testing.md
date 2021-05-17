---
id: p016_nms_regression_test
title: NMS Regression Testing
hide_title: true
---

# Proposal: NMS Regression Testing

Author(s): @andreilee

Last updated: 04/22/2021

Discussion at
[https://github.com/magma/magma/issues/4888](https://github.com/magma/magma/issues/4888).

## Abstract


The Magma NMS currently lacks many end-to-end tests for critical worksflows.
Given the current velocity of contributions to the Magma NMS, a tool such as
[Testim](https://www.testim.io/) will be used and integrated with Magma's
CI framework to cover the gap in testing, with limited time investment from
involved developers. 

## Background

As of the time of this document's conception, Magma has high feature velocity.
Without infinite developer resources, this presents a trade-off between
spending time on feature work and validation.
The use of tools such as Testim are intended to free up time for developers to
focus on feature work.

## Proposal

Testim will be used for daily testing of various workflows on NMS.
This testing will be limited to what can be automatically tested without
interaction with eNodeBs, CPEs, or UE devices, and will focus on interactivity
with the Magma NMS.

A representative test set of NMS workflows useful for verification of pure-play
FWA usage of Magma will be created on Testim.
The Testim account will be managed by the Magma Core Foundation, and the
subscription will be paid for by Facebook.

The test suite will be triggered to run daily through Github Actions.
A kubernetes cluster will be set up to run daily images of the Magma
orchestrator and NMS. The tests will run on this test cluster.

**Daily Reports**

A Slack channel will be used for contributors to see the results of tests.
The results of each test should be visible in the Slack channel, and a link
to the test run should also be provided, as well as a reference to which
container image tags were used for the test run.

**Budget**

Facebook will pay for the associated costs with using the service provided by
Testim, for at least one year.

**Tests**

The following tests are proposed to verify pure-play FWA usage of the Magma
NMS:

1. Create/update/delete an LTE network
2. Verifying dialogs for PLMN restriction, dataplans, and IMEI restrictions
2. Create/update/delete an LTE AGW
3. Create/update/delete LTE AGWs with IP allocation schemes: IP_POOL, bridged, STATIC
4. Create and register a managed eNodeB configuration
5. Create and register an unmanaged eNodeB configuration
6. Verify LTE networks dashboard
7. Verify Grafana dashboard
8. Verify charting components
9. Verify presence of logs/events components
10. Create/update/delete a gateway pool and register the LTE AGW
11. Create/update/delete a subscriber through single and bulk edits
12. Create/update/delete an APN
13. Create/update/delete a policy
14. Create/update/delete a traffic profile
15. Create/update/delete a rating group
16. Search through NMS metrics explorer
17. Add/update/delete an alert receiver
18. Add/update/delete an alert rule
19. Create/update/delete an NMS user
20. Create/update/delete an NMS organization


## Implementation

Implementation will proceed in 4 phases.
The entire implementation is proposed to be completed by the end of June 2021,
and should validate the Magma NMS for version 1.6 and onwards.

**Phase 1: Proposal Review and Github Comments**

*April 19 - 26*

The proposal will be visible on Github for at least a week to solicit comments
and feedback. At the end of this phase, there should be some consensus on the
plan to add NMS E2E/regression testing.

**Phase 2: Testim Test Suite Creation**

The test suite will be created on the Testim account.

At the end of this phase, there should be a full test suite capable of
verifying the use of Magma NMS for pure-play FWA usage. This test suite will
need to be run manually.

**Phase 3: Daily Github Actions Triggering Tests**

Github Action integration will be completed during this phase.
At the end of this phase, the full test suite should run on a daily basis.
This test suite should run against the Magma staging setup, currently
maintained by Facebook. Results of the tests should be posted to the Magma
Slack.

**Phase 4: Kubernetes Test Cluster for NMS Testing**

A separate kubernetes test cluster will be set up specifically for NMS testing.
This cluster will run daily builds of orchestrator and NMS.
The test suite will run against the NMS hosted on this test cluster to get
an up-to-date view of the NMS.
