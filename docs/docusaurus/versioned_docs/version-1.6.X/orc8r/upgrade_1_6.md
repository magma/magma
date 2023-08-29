---
id: version-1.6.X-upgrade_1_6
title: Upgrade to v1.6
hide_title: true
original_id: upgrade_1_6
---

# Upgrade to v1.6

Before proceeding, read the [Orc8r upgrade introduction](./upgrade_intro.md).

This guide covers upgrading Orchestrator deployments from v1.5 to v1.6. The v1.6 upgrade follows a standard procedure

- Upgrade Terraform modules
- Terraform install

## Prerequisites

### Terraform version

This document was tested using Terraform version `0.15.0`. You can use [`tfenv`](https://github.com/tfutils/tfenv) to install a particular Terraform version.

### Orc8r artifacts

As of the v1.6 release, Orc8r artifacts are now available on a publicly-hosted artifactory. This means you no longer need to build from source or manage your own artifactory.

See the [host custom artifacts](#host-custom-artifacts) section if you would rather build your own artifacts.

## Upgrade

### 1. Update Terraform module

Update your root Terraform module

```terraform
# This will likely be found in main.tf

module orc8r {
  source = "github.com/magma/magma//orc8r/cloud/deploy/terraform/orc8r-aws?ref=v1.6"

  # Optional: enable AWS DB failure notifications
  # orc8r_sns_email             = "admin@example.com"
  # enable_aws_db_notifications = true

  # ...
}

module orc8r-app {
  # IMPORTANT: delete all docker_ and helm_ lines

  source = "github.com/magma/magma//orc8r/cloud/deploy/terraform/orc8r-helm-aws?ref=v1.6"
  orc8r_tag = "v1.6.0"

  # Optional: prune ES indices when disk usage exceeds 75% of allocated ES storage
  # elasticsearch_disk_threshold = tonumber(module.orc8r.es_volume_size * 75 / 100)

  # ...
}
```

Pull in the module changes

```bash
terraform init --upgrade
terraform refresh
```

### 2. Terraform apply

Apply the upgrade

```bash
terraform apply  # check this output VERY carefully
```

## Appendix

### Postgres version

[AWS is terminating support for Postgres 9.6 on November 11, 2021](https://aws.amazon.com/blogs/database/upgrading-from-amazon-rds-for-postgresql-version-9-5/), in accordance with the [Postgres versioning policy](https://www.postgresql.org/support/versioning/). Postgres 9.6 was the default Orc8r DB target for the past several Magma releases.

To upgrade your Postgres database to the correct version follow [these instructions](https://magma.github.io/magma/docs/orc8r/rds_upgrade#logs-and-validation).

### Host custom artifacts

If you wish to use your own artifacts, follow the follow the [Build Orchestrator](https://magma.github.io/magma/docs/orc8r/deploy_build) page to build and publish your own artifacts from the [v1.6 branch of the Magma repo](https://github.com/magma/magma/tree/v1.6), then update your `main.tf` file to point to the appropriate artifactories.
