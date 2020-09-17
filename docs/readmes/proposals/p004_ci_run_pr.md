---
id: p004_ci_run_pr
title: Run CI on Pull-Request
hide_title: true
---

# Run CI on Pull-Requests

- Feature owners: `@tmdzk`
- Feedback requested from: `@apad`, `@xjtian`, `@rckclmbr`

## Summary

Currently Pull Request are triggering a batch of test (docker build, flow test, lint, containerized tests..) However we're currently not testing the integration test (lte-integ-test, cwag-integ-test, xwfm-test) for 2 mains reasons

1. Those tests are running on private hardware protected by a VPN
2. Those tests are using sensitive credentials that cannot be shared without a code check of the pull request.

## Motivation

The current method of testing a pull request is to code review it, make sure all test that can be run on public infra (CircleCI) pass and then merge the pull request that will automatically trigger lte-integ-test, cwag-integ-test, xwfm-test on private hardware.
Sometimes the new code will break one of those test and we have to go back to the beginning of the loop above, new PR, wait for the tests that can run on public infra (CircleCi) then merge and wait for those lte-integ-test, cwag-integ-test, xwfm-test to rerun again.

This process could be simplified if we could run those specific tests on PR, and avoid back and forth on PRs or breaking master.

## Goals

The goals of running all tests on PRs are:

- Catch errors as soon as possible and avoid "Hotfix PRs"
- Keep master cleaner
- Make sure contributors feel comfortable contributing without having the "fear" to break master

## Proposal

```
+---------------+          +----------------+        +-------------+
|               |          |                |        |             |
|    PR Created |          | Maintainer     |        |             |
|               +--------->+ approval       +------->+ Branch      |
|               |          |                |        | creation    |
+---------------+          +----------------+        +------+------+
                                                            |
                                                            |
                                                            |
                                                            v
                         +-------------------+      +-------+--------+
                         |                   |      |                |
                         |                   |      |                |
                         | Report back on the+<-----+Triggers Circle |
                         | PR                |      |job integ-test  |
                         |                   |      +----------------+
                         +-------------------+
```

The flow is pretty straight forward and is composed in 5 steps:

1. PR creation
2. A maintainer review the code and make sure no malicious code was inteded to be injected. When It's done, the contributor comments with
`/ok-to-test`
3. A bot/Ci job will run through (schedule or github webhook) and create a branch directly on magma/magma, this branch will be prefixed with `test-<pr_number>` and tag the PR as `ok-to-test`
4. CircleCI triggers those specific tests for the newly created branches
5. The job/bot reports back to the PR.

## Technical Details

### Maintainer approval

The maintainer approval for testing is a simple comment on the PR that will be read by the bot/job. After reviewing the code just write a comment saying
`/ok-to-test`.

### Banch Creation

When the comment has been drop on the PR we have 2 ways to trigger the branch creation.
1. Through a github webhook which lead to host an api and a custom bot that we'll have to write
2. A simple CircleCI job that goes through all recent Pull requests and make sure branches have been created and test triggered

When the job/bot has been trigger, it will have a github key that allow it to create a branch, so the job will just fetch the code in the PR
and create a branch with it.

### CircleCI triggers

The special formatting of the branch name `test-<pr_number>` will automatically trigger the test we'll just create a new filter by branch with `master` and `test-*`

There is another solution which is to use directly CircleCI api.


## Timeline of Work

- Choose technical details depending on the investment we want to put in that project (simple CI job/a full bot..)
- Create the branch generator from maintainer's comment on PR
- Change CircleCI config.yml to trigger test on all branch preprend with `test-`
- Create the report back on PR callback
- Test the full pieces together using forks of magma

## Futur improvement/ Proposal

1. We can definitely think about a set of command like `/retest`, `/run-<specific-test>`. We can even only manage our PR with this bot (closing/approving etc)

2. This proposal was inspired by some opensource project that use this system of maintainer's comment and a bot. (eg: jenkinsX)
When you dig you can see that they are using [Prow](https://github.com/kubernetes/test-infra/tree/master/prow) which is part of the Kubernetes repo.
This could definitely be something we can explore as a CI/CD solution in the future when we'll want to merge all CIs together especially because
the "approval system" is already in built and it interacts well with github using `/command`.

Example of a [Prow dashboard](https://prow.k8s.io/)


## References and Thanks

Thanks to @rckclmbr for helping with that

Reference to JenkinsX opensource workflow (example here : https://github.com/jenkins-x/jx/pull/7544#issuecomment-682097513)
