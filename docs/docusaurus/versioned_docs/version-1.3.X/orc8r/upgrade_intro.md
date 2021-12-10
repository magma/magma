---
id: version-1.3.0-upgrade_intro
title: Orchestrator Upgrades: Introduction
hide_title: true
original_id: upgrade_intro
---

# Orchestrator Upgrades: Introduction

Orchestrator upgrades generally follow a standard pattern:

- Upgrade deployment tooling (e.g. Terraform, Helm, etc.)
- Build and publish new application artifacts
- Publish new Helm charts
- Deploy new application artifacts
- Run data migrations and other post-upgrade steps

Every new minor and major release (i.e. non-hotfix) will have a corresponding
documentation section on upgrading from the prior release. Upgrade procedures
are only documented for adjacent releases. Upgrades which skip versions
(e.g. v1.1 to v1.3) are not explicitly supported at this time.

Before every upgrade, we strongly suggest taking a DB snapshot in case you
decide that you need to roll back your application. While the application
components are stateless and can be upgraded or downgraded at any time, the
data migrations that come with each release are not always reversible.
