---
id: p006_mandatory_integration_tests_for_each_PR.md
title: Enable integration tests as a mandatory check required to merge any PR
hide_title: true
---# Enable integration tests as a mandatory check required to merge any PR
- Feature owners: `@119vik`, `@tdzik`, `@mattymo`
- Feedback requested from: `@apad`, `@zer0tweets`, `@ttx`
## Problem descriptions
At the moment, buggy codes gets merged into master branch time after time, which affects development of new features, slowing developers down. This happens because some PRs get merged without performing proper QA by the code maintainers who are not willing to wait until integration tests have executed. Here are some numbers which can explain what does "wait too long" means:
 - magma_integration_tests pipeline takes ~1 hour to execute and Jenkins can process only 1 run/hour (does not support concurrent builds)
 - existing test env is a VirtualBox VMs managed by Vagrant on top of KVM VMs which are running on baremetall node is very slow. 
 - average magma_integration_tests queue length was 30-40 jobs (before opt-in strategy was implemented)
 - CWAG-integration-test and LTE-integration-test jobs take 1 hour each (before optimisation.) Jenkins can process approximtely 13 executions per hour. 
 Recently we switched to "Opt-In strategy" for PRs affecting the gateway code base. These PRs must be marked with one of the following labels: "component: agw" OR "component: cwf" OR "component: feg". Implementation of this strategy decreased load on jenkins (went from 30+ to 2-10 per day). At the same time, this means that PRs which were not marked with the correct label are not tested and may be introducing new bugs.
 
 We can draw a conclusion from above and outline the following pain points:
 - Some PRs are still merged into master without integration tests
 - Integration tests are slow (actual execution is slow, regardless of VM setup time) so we can't realistically run them for each commit of each PR.
 - Reviewers are not applying tags to PRs and waiting for them to complete (or skipping gatekeeping tests)
 - Magma_integration_tests pipeline is the largest bottleneck in our CI process
 

## Proposed solutions

### Integration tests are slow(which is more or less expected) so we can't run them for each commit of each PR.
We should execute integration tests for each and every PR which is going to be merged in to master. But there's no sense to do it for each and every commit.
Solution 
Stage 1:
 - Ensure that for each commit of each PR all unit test jobs are executed by CircleCI
 - Ensure that each PR was reviewed and accepted by component maintainer
 - Ensure that each PR which succesfully passed unit testing and maintainer review gets "ready-to-merge" label assigned.
 - Ensure that each PR marked with label "ready-to-merge" was tested by all integration tests.
 Stage 2:
 - Ensure that each PR marked with label "ready-to-merge", which passed integration tests is automatically merged by CI.
 Stage 3:
 - Ensure that no side effects appear when PRs are automatically merged by CI (so speculative CI implementation in fact).

### Some PRs are still merged into master without integration tests
Make Jenkins as a mandatory gatekeeper for merging each PR. Proposed workflow:
 - Reviewer marks appropriate labels to run necessary relevant tests (such as CWAG / NMS / etc.) before indicating "ready-to-merge" (tag or some keyword in a comment)
 - Jenkins triggers all tests
 - When all jobs are green, Jenkins merges the PR (or alternatively set a test result "Jenkins CI Gate: Passed")

### Reviewers are not applying tags to PRs and waiting for them to complete (or skipping gatekeeping tests)
We run twice daily(EU time / US night time) master branch regression test builds (all tests) and look for regressions that could have gotten merged into master. Regressions can happen for the following reasons:
- Reviewer did not wait for tests
- Two conflicting PRs were merged that had logic that breaks something
- External factors not covered by PR process (such as a library package update)
This regression pipeline runs twice a day(EU time / US night time) on a schedule. If the build fails, a GitHub issue is created tagging Magma CI team for urgent attention. 

### Magma_integration_tests pipeline is the largest bottleneck in our CI process
At the moment magma_integration_tests pipeline is structured in such a way that it's using a Jenkins slave as a jump host to a real test environment. We only have one real test env for the magma_integration_tests for now. We need to do some simple refactoring to ensure that magma_integration_tests utilizes a Jenkins slave as a real test env (like CWAG and LTE integration tests do). This means replacing all the 'ssh <user>@<host> - '<command>'' with corresponding 'sh <command> ' statements in the pipeline definition. This simple refactoring should allow Jenkins to handle up to 16 builds per hour (instead of just 1 at the present moment).

### Existing test env is a VirtualBox VMs managed by Vagrant on top of KVM VMs which are running on baremetall node is very slow.
As we found as a result of performance tests current setup gives us low performace because of Virtualbox running in KVM. Best option we can suggest so far is start using libvirt based vagrant boxes which 2-3 times faster then Virtualbox.

## Execution plan
Steps:
- Refactor magma_integration_tests to enable concurrent execution for multiple PRs
- Prepare Jenkins pipeline which can build libvirt boxes for vagrant and publish using community account.
- Refactor all existing pipelines to use libvirt vagrant boxes.
- Configure Jenkins to run Nightly Builds against master with All integration tests
- Configure Jenkins to run All Integration tests against any PR marked with "ready-to-merge" label
- Configure Jenkins to raise an Github Issue in case of nightly build failure
- Configure Jenkins to merge PR when all integration tests passed for any PR marked with "ready-to-merge" label

Optional Steps:
- Setup Speculative CI. CI which can merge several PRs together in temporary branch - confirm that code is still stable with unit and integration tests and merge temporary branch into master.
