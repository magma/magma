---
id: team-processes
title: Internal Processes
---

### Release Process

Symphony is using the continuous push methodology. This means that new code is pushed to production whenever changes are done, usually every 30 minutes.
This enables the Symphony team to move fast and react to partner requests in real time. Bugs are fixed in a matter of hours, and new requests are developed in a matter of days.

The product is protected by numerous automated tests. Unit tests, integrations tests and UI E2E tests are all in place to block the push in case a major feature was broken.

### SEV Process

The Symphony team is taking any breakage in the product seriously. 
Every week one team member is an “oncall” - responsible for the health of production. He is constantly fixing bugs, monitoring any report from our partners and improving the quality of the tool.
Whenever a serious problem occurs, we are opening a “SEV”. SEV is an incident report, that requires “all hands on deck”. It means that all of the team is focused on solving the issue ASAP.
The severity of the SEV is as follows:

* SEV 3
    * A high priority feature is not working in prod (e.g. connect links, pyinventory)
    * High number of intermittent failures
* SEV2
    * 1-2 prod partners are down
    * Internal partner is down
* SEV1
    * "Production" is down (Inventory\WO is inaccessible for all partners)
    * Data layer is inconsistent and partner's data is lost

Our commitment towards fixing SEVs is as follows:
* SEV 3: Fix  during regular business hours.
* SEV 2: Fix with reasonable after hours. Feel free to ping others or even wake up people until the problem is resolved.
* SEV 1: All hands on deck until fixed! Fix, even with unreasonable after hours.
 
After the SEV is mitigated, the team is having a “postmortem” meeting to review the SEV. SEV timeline, cause, time to mitigate- all numbers are reviewed. The team is leaving this meeting with a set of tasks to do: add automated tests, improve code infrastructure, fix bugs, etc.

After every SEV, our goal is not only to fix the problem, but also to fix the code in such a way that similar SEVs will never occur.

### Deprecating APIs
If changes are done to the schema we will mark old endpoints as deprecated.
Here you will find the list of deprecated endpoints [list of deprecated endpoints](graphql-breaking-changes.md)

Partners are expected to upgrade their code to the new version.
After <TBD> months of Deprecated state we will delete old endpoints.
 
<< TBD more details >>