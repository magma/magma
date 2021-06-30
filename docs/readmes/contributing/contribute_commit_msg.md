---
id: contribute_commit_msg
title: Commit message style
hide_title: true
---

# Improve git commit messages
This document is about recommended git commit style.  Every project generally
has conventions for writing commit message subject line.  Rather than defining
our own convention, magma can make use of commit subject structure
according to conventionalcommits.org.

##tldr
Commit message can follow the template defined in magma github.
This document covers commit subject line structure and
validation APP.
```
<type>[optional scope]: <description>
```
for example:
```
fix(AGW): fix pyroute2 dependency
```

## Advantage of semantic commits.
1. This state intention of the commit msg in concise manner in PR title.
2. It can be used for automatically compiling change log for magma release
   documentation.
3. Analytics related to number of fixes vs feature commits in each release.
4. This would force developer to rethink when they mix different type of
   changes in same PR.

## type for magma repo
Every commit would need to use `type` for each commit. Following would be the
type available for use.
```
  - feat
  - fix
  - docs
  - style
  - refactor
  - perf
  - test
  - build
  - chore
  - revert
  - proposal
```

## scope for magma repo
After `type` commit subject should have `scope` of the change. `scope` is
optional  since it  can be automatically derived from the change. there is
already a bot to set component for every git PR.

```
  - orc8r
  - nms
  - feg
  - agw
  - mme
  - pipelined
  - sessiond
  - mobilityd
  - subscriberdb
  - policydb
  - enodebd
  - health
  - monitord
  - redirectd
  - smsd
  - envoy_controller
  - ctraced
  - directoryd
  - eventd
  - magmad
  - state
  - ci
  - cwg
  - xwf
  - agw
  - cloud
  - upf
  - smf
  - amf
```

For an up-to-date list of types we support see checkout Semantic Github APP config:
https://github.com/magma/magma/blob/master/.github/semantic.yml

For details on commit specification refer: https://www.conventionalcommits.org/en/v1.0.0/

## commit message body
magma github PR commit msg template covers required structure for commit message.

## Validation of commit mesages
semantic-pull-requests github app validates the commit subject according to
type and scopes defined for magma. This would be helpful while reviewing PR.
Configuration for this APP is in defined in `.github/semantic.yml`, As we add
more components to magma source code this file needs to be updated.

https://github.com/apps/semantic-pull-requests
