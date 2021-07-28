---
id: contribute_ci_checks
title: Continuous Integration Checks
hide_title: true
---

# Continuous Integration Checks

When a PR is submitted to the Magma repo, a suite of Continuous Integration tests and analysis are executed. Checks marked as required must pass for a PR to merge.

This document attempts to give high level explanations of the purpose of each check, contacts for questions regarding each test, and expectations around CI flow success.

## How to reach out to a check owner

Below are some resources for finding who to contact:
* [GitHub-to-Slack ID mappings](contribute_id_mappings)
* [List of Magma maintainers](https://github.com/orgs/magma/teams/repo-magma-maintain/members)
* [List of `approvers-*` teams and their members](https://github.com/orgs/magma/teams/?query=approvers)

## Blocking Checks

Merge blocking CI checks are listed below.

| Check Name             | Purpose                                       | Owner                  | Remediation Steps                                                            |
| ---------------------- | --------------------------------------------- | ---------------------- | ---------------------------------------------------------------------------- |
| DCO Check              | Check PR is signed off                        | any maintainer         | [PR guidelines](contribute_workflow#guidelines)                              |
| Semantic PR            | PR format checker                             | any maintainer         | [PR guidelines](contribute_workflow#guidelines)                              |
| mergefreeze            | Stop merges during a code freeze              | themarwhal             | N/A                                                                          |
| orc8r-build            | Validate Orc8r builds                         | approvers-orc8r        | [Orc8r build](../basics/quick_start_guide#terminal-tab-2-build-orchestrator) |
| cloud-test             | Run Orc8r unit tests                          | approvers-orc8r        | [Orc8r tests](../orc8r/dev_testing)                                          |
| cloud_lint             | Check Orc8r changes satisfy golint            | approvers-orc8r        | [Orc8r tests](../orc8r/dev_testing)                                          |
| insync-checkin         | Ensure generated files are committed          | approvers-orc8r        | [Orc8r tests](../orc8r/dev_testing)                                          |
| feg-build              | Validate FeG builds                           | approvers-feg          | [FeG build](../feg/deploy_build)                                             |
| feg-precommit          | Run FeG unit tests                            | approvers-feg          | [FeG tests](../feg/dev_testing)                                              |
| nms-build              | Validate NMS builds                           | approvers-nms          | [NMS build](../basics/quick_start_guide#using-the-nms-ui)                    |
| eslint                 | Ensure NMS changes satisfy eslint             | approvers-nms          | [NMS tests](../nms/dev_testing)                                              |
| nms-flow-test          | Validate NMS changes satisfy flow             | approvers-nms          | [NMS tests](../nms/dev_testing)                                              |
| nms-yarn-test          | Run NMS unit tests                            | approvers-nms          | [NMS tests](../nms/dev_testing)                                              |
| nms-e2e-test           | Run NMS end to end tests                      | approvers-nms          | [NMS tests](../nms/dev_testing)                                              |
| lte-test               | Run AGW Python unit tests                     | approvers-agw          | [AGW tests](../lte/dev_unit_testing)                                         |
| mme_test               | Run MME and sctpd unit tests                  | approvers-agw-mme      | [AGW tests](../lte/dev_unit_testing)                                         |
| session_manager_test   | Run SessionD unit tests                       | approvers-agw-sessiond | [AGW tests](../lte/dev_unit_testing)                                         |
| orc8r-gateway-test     | Run Golang unit tests for orc8r/gateway       | approvers-agw          | [AGW tests](../lte/dev_unit_testing)                                         |
| cwag-precommit         | Run CWAG unit tests                           | approvers-cwf          | [CWAG tests](../cwf/dev_testing)                                             |
| Magma-OAI-Jenkins      | OAI's MME integration test run on OAI Jenkins | rdefosse               | N/A                                                                          |
| cwf-operator-build     | Validate CWF deployer builds                  | approvers-cwf          | TODO                                                                         |
| cwf-operator-precommit | Run CWF deployer unit tests                   | approvers-cwf          | TODO                                                                         |

## Non-blocking Checks

The CI checks listed below do not block merging on failure.

| CCheck Name                                   | Purpose                                                          | Point of Contact   | Remediation Steps                                    |
| --------------------------------------------- | ---------------------------------------------------------------- | ------------------ | ---------------------------------------------------- |
| Shellcheck by Reviewdog                       | Annotate PRs with shell script static analysis                   | electronjoe        | N/A                                                  |
| GCC Warnings & Errors / gen_build_container   | Passive Check generates container for GCC builds                 | electronjoe        | N/A                                                  |
| GCC Warnings & Errors / build_oai             | Annotate PRs with any GCC Warnings or Errors for MME             | electronjoe        | N/A                                                  |
| GCC Warnings & Errors / build_session_manager | Annotate PRs with any GCC Warnings or Errors for session_manager | electronjoe        | N/A                                                  |
| Jenkins CWAG Libvirt                          | Run CWF integration tests                                        | themarwhal mattymo | TODO                                                 |
| ci/circleci: c-cpp-codecov                    | Upload AGW C/C++ code coverage                                   | electronjoe        | N/A                                                  |
| ci/circleci: cwag-build                       | Validate CWAG builds                                             | approvers-cwf      | TODO                                                 |
| ci/circleci: feg-lint                         | Check FeG changes satisfies Go linter                            | approvers-feg      | TODO                                                 |
| Python Format Check                           | Ensure Python changes are formatted                              | themarwhal         | [AGW formatting](../lte/dev_unit_testing#format-agw) |
