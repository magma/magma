---
id: version-1.4.0-contribute_codeowners
title: Codeowners Workflow
hide_title: true
original_id: contribute_codeowners
---

# Codeowners Workflow

This document covers the development workflow for codeowners (maintainers) of the Magma project.

## Guidelines on force-merging a PR

The `master` branch of the `magma/magma` repo is a [protected branch](https://docs.github.com/en/github/administering-a-repository/about-protected-branches). This means

- [PR must be reviewed and accepted](https://docs.github.com/en/github/administering-a-repository/about-protected-branches#require-pull-request-reviews-before-merging) by relevant [codeowners](https://docs.github.com/en/github/creating-cloning-and-archiving-repositories/about-code-owners)
- [PR must pass required CI checks](https://docs.github.com/en/github/administering-a-repository/about-protected-branches#require-status-checks-before-merging)

However, users with `admin` permissions to the `magma/magma` repo can bypass these requirements, performing a "force-merge".

This functionality should be used *only when in dire need*, as it can easily break the build in unexpected ways.

Thus, when an admin force-merges a PR, they *must* comment on the PR describing the reason why it's being force-merged, referencing one of the following guidelines.

### Unexpected Emergency

True emergency, where the existing set of guidelines doesn't fit the need.

There should be a postmortem to decide whether the existing set of guidelines needs to be updated to encompass the experienced emergency.

### Chicken-and-egg

CI is broken, and fixing CI requires overriding a required check.

### UBN during broken CI

CI is broken and blocking a UBN ("unbreak now") fix.

This is still high-risk and should be used with extreme care. One potential example: an existential threat to the project caused by severe security vulnerability.

### Codeowner OOO

For high-priority work: if the PR author is a codeowner for the code they changed, and all other codeowners for that part of the codebase are OOO ("out of office") for an extended period, the codeowner can request a force-merge once all CI checks have passed.

This guideline circumvents the quirk that a PR author can't approve their own PR, even if they are a codeowner for part of the PR's changes. The author *must* be a codeowner, and *must* share permissions with the OOO codeowner (e.g. they're both on the `orc8r-approvers` codeowner team).
