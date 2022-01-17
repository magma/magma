---
id: version-1.6.X-contribute_ci_checks
title: Continuous Integration Checks
hide_title: true
original_id: contribute_ci_checks
---

# Continuous Integration Checks

When submitting PRs to the Magma repo, a suite of Continuous Integration tests and analysis are executed (some dependent upon the files modified in the PR).  This document attempts to give high level explanations of the purpose of each check, contacts for questions regarding each test, and expectations around CI flow success.

| Check Name | Purpose | Point of Contact |
|-|-|-|
| CodeQL / Analyze (go) | Security analysis for all Golang |  |
| DCO check / DCO Check | Ensure all Pull Requests are Signed | N/A |
| Shellcheck by Reviewdog | Annotate Pull Requests with shell script static analysis | electronjoe (GH), Scott Moeller (slack) |
| CodeQL / Analyze (javascript) | Security analysis for all Javascript |  |
| CodeQL / Analyze (python) | Security analysis for all Python |  |
| GCC Warnings & Errors / gen_build_container | Passive Check generates container for GCC builds |  |
| GCC Warnings & Errors / build_oai | Annotate Pull Requests with any GCC Warnings or Errors for MME target | electronjoe (GH), Scott Moeller (slack) |
| GCC Warnings & Errors / build_session_manager | Annotate Pull Requests with any GCC Warnings or Errors for session_manager target | electronjoe (GH), Scott Moeller (slack) |
| Jenkins | LTE integration test suite in Jenkins | mattymo (GH), Matthew Mosesohn |
| Jenkins CWAG Libvirt | CWF integration test suite in Jenkins |  themarwhal (GH), Marie Bremner (slack) |
| Jenkins LTE | A lighter weight and low flake subset of LTE tests including...? |  |
| Code scanning results / CodeQL | How does this differ from the other CodeQL / Analyzer stages? |  |
| Magma CI autolabel | ? |  |
| Magma-OAI-Jenkins | OAI's MME integration test run on OAI Jenkins | rdefosse (GH), Raphael Defosseux (slack) |
| ci/circleci: c-cpp-codecov | AGW c/c++ unit tests and push coverage to codecov.io | electronjoe (GH), Scott Moeller (slack) |
| ci/circleci: cloud-test | Orc8r unit tests | hcgatewood (GH), Hunter Gatewood (slack) |
| ci/circleci: cloud_lint | ? |  |
| ci/circleci: cwag-precommit | CWF Unit Tests | themarwhal (GH), Marie Bremner (slack); uri200 (GH), Oriol Batalla (slack) |
| ci/circleci: cwag-build | Validate that CWAG builds | themarwhal (GH), Marie Bremner (slack); uri200 (GH), Oriol Batalla (slack) |
| ci/circleci: cwf-integ-test | Run CWF integration test suite in CircleCI | themarwhal (GH), Marie Bremner (slack); uri200 (GH), Oriol Batalla (slack) |
| ci/circleci: cwf-operator-build | ? |  |
| ci/circleci: cwf-operator-precommit | ? |  |
| ci/circleci: eslint | ? |  |
| ci/circleci: feg-build | 	Validate that FeG builds | themarwhal (GH), Marie Bremner (slack); uri200 (GH), Oriol Batalla (slack); emakeev (GH), emak (slack) |
| ci/circleci: feg-lint | FEG Go linter and code coverage | themarwhal (GH), Marie Bremner (slack);  uri200 (GH), Oriol Batalla (slack); emakeev (GH), emak (slack) |
| ci/circleci: feg-precommit | Run FEG unit tests | themarwhal (GH), Marie Bremner (slack);  uri200 (GH), Oriol Batalla (slack); emakeev (GH), emak (slack) |
| ci/circleci: insync-checkin | ? |  |
| ci/circleci: lte-test | Run AGW Python unit tests | pshelar (GH), Pravin Shelar (slack); ardzoht (GH), Alex Rodriguez (slack) |
| ci/circleci: lte-integ-test | Run LTE integration test suite in CircleCI | ardzoht (GH), Alex Rodriguez (slack); ulaskozat (GH), Ulas Kozat (slack) |
| ci/circleci: mme_test | Validate that MME related unit tests build and pass | themarwhal (GH), Marie Bremner (slack) |
| ci/circleci: nms-build | Validate that the NMS builds | karthiksubravet (GH), karthik subraveti (slack); andreilee (GH), Andre Lee (slack) |
| ci/circleci: nms-e2e-test | ? | karthiksubravet (GH), karthik subraveti (slack); andreilee (GH), Andre Lee (slack) |
| ci/circleci: nms-flow-test | ? | karthiksubravet (GH), karthik subraveti (slack); andreilee (GH), Andre Lee (slack) |
| ci/circleci: nms-yarn-test | ? | karthiksubravet (GH), karthik subraveti (slack); andreilee (GH), Andre Lee (slack) |
| ci/circleci: orc8r-build | Validate that the orc8r builds | |
| ci/circleci: orc8r-gateway-test | ? | |
| ci/circleci: session_manager_test | Validate that the session_manager related unit tests build and pass | themarwhal (GH), Marie Bremner (slack) |
| netlify/* | ? | |
