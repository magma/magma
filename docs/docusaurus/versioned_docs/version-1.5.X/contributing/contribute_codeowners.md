---
id: version-1.5.0-contribute_codeowners
title: Codeowners Workflow
hide_title: true
original_id: contribute_codeowners
---

# Codeowners Workflow

This document covers the development workflow for codeowners (maintainers) of the Magma project.

## Shepherds not gatekeepers

The intent is for Magma Codeowerns to behave as shepherds of contributions, not gatekeepers.  This outlook asks more of the codeowner, but aims to maximally leverage community involvement and goodwill.  The following are some guidelines we hope Magma Codeowners will integrate into their reviews.

### Show instead of tell

For nuanced or unclear asks by the code reviewer, it is drastically more efficient for the contributor (and often for all parties involved) if the reviewer can show the contributor what is being asked of them in the PR, rather than tell them.  This could be done in a number of ways.

- Direct code suggestion in GH Review comment
- Create a quick GitHub Gist that is shared in the review comment
- Suggesting a direct edit to the PR via Gihub UI
- Cherry-picking the PR in a reviewer branch and prototyping the suggestion, sharing branch

The benefits of this include.

- Fewer round trips between reviewer and contributor
- Faster PR revisions by the contributor
- Can reveal problems with the reviewer ask - to the reviewer

In net these result in fewer abandoned PRs, fewer merge conflicts due to delays, more efficient use of everybody's time, cleaner code commits (more likely to employ best practices) and can highlight to the code owner work that needs to be done to make the code base more aligned with their intent.

### Avoid moving goalposts

It is a common phenomenon to see new changes that are precipitated by an in-progress or in-review PR.  Magma Codeowners should resist the urge to move goalposts to include new refactoring or substantial changes in asks of the contributor.  Instead, follow-up work should be immediately described in a GH Issue and assigned either to the codeowner, the PR contributor (if they are game), or some other Magma contributor.  The rationale is that moving goalposts results in many stalled PRs - and our objective is not to stall PRs. Further, the best PR is a small PR - and movement of goalposts almost always results in larger PRs (one could use this metric to determine whether it's reaonsable to ask for a change in direction - does it enlarge the PR?).


### PRs with limited unit / integration tests

Where Magma has made it reasonably low-friction, is is good practice for Magma Codeowners to request tests be included with community PRs.  Unfortunately, there exist plenty of Magma code regions which do not readily admit to test today due to architecture, libraries in play, or build configuration issues.  Magma is committed to improving this situation but codeowners should not hold contributors to a test bar that has not yet been made accessible to them.

### Un-Testable PRs

Some community PRs may contain changes that are untestable given Magma resources or existing test infrastructure. This covers a wide range of possibilities - but as an example perhaps a new radio is being added that is a mild API change compared to some existing supported model.  If the change is deemed unlikely to break existing supported and tested implementations, and is easily reverted, then the risk to allowing a merge of a community-tested PR seems worth acceptance without confirmation of health.

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
