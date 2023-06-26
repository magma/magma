---
id: version-1.8.0-deploy_intro
title: Introduction
hide_title: true
original_id: deploy_intro
---

# Introduction

This section walks through installing a production Orchestrator deployment.

We assume you will use the versioned artifacts provided by the project's official artifactory at [linuxfoundation.jfrog.io](https://linuxfoundation.jfrog.io/). If you would like to build and host your own artifacts, see the [Build Orchestrator](./dev_build.md) page.

To deploy orc8r, see the [Manual installation](./deploy_install.md) page.

## Prerequisites

Throughout this guide we assume the `MAGMA_ROOT` environment variable is set to the local directory where you cloned the Magma repository

```sh
export MAGMA_ROOT=PATH_TO_YOUR_MAGMA_CLONE
```

Before deployment, it may be useful to read through the [Magma prerequisites](../basics/prerequisites.md) and [Magma quick start guide](../basics/quick_start_guide.md) sections.

Familiarity with the following is assumed

- AWS
- Kubernetes
- Terraform

The instructions in this section have been tested on macOS and Linux. If you are deploying from a Windows host, some shell commands will likely require adjustments.

## Deployment types

Orc8r deployment type specifies the Orc8r modules which will be included to manage Magma gateways. It supports following deployment types

- `fwa` for fixed wireless deployment, enables management of *AGWs*
- `federated_fwa` for federated fixed wireless deployment, enables management
  of *AGWs and FEGs*
- `all` for all-encompassing deployments, enables management of *AGWs, FEGs,
  and CWAGs*

## Release versioning

Orc8r follows the standard [semantic versioning scheme](https://semver.org/) of `MAJOR.MINOR.PATCH`. Generally speaking, a bump in each version type involves the following

- `MAJOR` considerable change
    - Orc8r: major manual intervention
    - General: large-scale changes, e.g. to the conceptual function of Magma as a whole
- `MINOR` non-trivial change
    - Orc8r: minor manual intervention may be required
    - Orc8r-gateway interface: gateways may need to be updated to new minimum version
- `PATCH` small, backward-compatible changes
    - Security or functionality-critical
    - Updating to newer patch should be seamless, with no manual intervention required

Major and minor releases are tagged off the master branch, then a patch branch is opened starting at that tag. Patch releases are tagged on the respective patch branch.

The current release schedule tags a new minor version 3-4 times per year. Patch releases are tagged on an on-demand basis.

## Deploying specific release

To target a specific release, checkout the Magma repository's relevant release branch when building artifacts. This is also a great place to find relevant Terraform values.

Values for recent Orchestrator releases are summarized below

### v1.8.0

Verified with Terraform version `1.0.11`.

- `v1.8` [patch branch](https://github.com/magma/magma/tree/v1.8)
- `github.com/magma/magma//orc8r/cloud/deploy/terraform/orc8r-aws?ref=v1.8`
Terraform module source
- `1.8.0` Helm chart version

### v1.6.0

Verified with Terraform version `0.15.0`.

- `v1.6` [patch branch](https://github.com/magma/magma/tree/v1.6)
- `github.com/magma/magma//orc8r/cloud/deploy/terraform/orc8r-aws?ref=v1.6`
Terraform module source
- `1.5.23` Helm chart version
- Additional notes
    - `9.6` PostgreSQL target release. Prefer `12.6` for new deployments.

### v1.5.0

Verified with Terraform version `0.14.5`. Terraform `0.15.x` is *not* compatible.

- `v1.5` [patch branch](https://github.com/magma/magma/tree/v1.5)
- `github.com/magma/magma//orc8r/cloud/deploy/terraform/orc8r-aws?ref=v1.5`
Terraform module source
- `1.5.21` Helm chart version
- Additional notes
    - `9.6` PostgreSQL target release, newer versions will likely work as well

### v1.4.0

Verified with Terraform version `0.14.0`.

- `v1.4` [patch branch](https://github.com/magma/magma/tree/v1.4)
- `github.com/magma/magma//orc8r/cloud/deploy/terraform/orc8r-aws?ref=v1.4`
Terraform module source
- `1.5.16` Helm chart version
- Additional notes
    - `9.6` PostgreSQL target release, newer versions will likely work as well

### v1.3.0

Verified with Terraform version `0.13.1`.

- `v1.3` [patch branch](https://github.com/magma/magma/tree/v1.3)
- `github.com/magma/magma//orc8r/cloud/deploy/terraform/orc8r-aws?ref=v1.3`
Terraform module source
- `1.4.36` Helm chart version
- Additional notes
    - `9.6` PostgreSQL target release, newer versions will likely work as well

### v1.2.0

Verified with Terraform version `0.13.1`.

- `v1.2` [patch branch](https://github.com/magma/magma/tree/v1.2)
- `github.com/magma/magma//orc8r/cloud/deploy/terraform/orc8r-aws?ref=v1.2`
Terraform module source
- `1.4.35` Helm chart version
- Additional notes
    - `9.6` PostgreSQL target release, newer versions will likely work as well

### v1.1.0

Verified with Terraform version `0.12.29`.

- `v1.1` [patch branch](https://github.com/magma/magma/tree/v1.1)
- `github.com/magma/magma//orc8r/cloud/deploy/terraform/orc8r-aws?ref=v1.1`
Terraform module source
- `1.4.21` Helm chart version
- Additional notes
    - `9.6` PostgreSQL target release, newer versions will likely work as well
