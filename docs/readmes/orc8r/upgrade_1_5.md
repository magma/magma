---
id: upgrade_1_5
title: Upgrade to v1.5
hide_title: true
---

# Upgrade to v1.5

This guide covers upgrading Orchestrator deployments from either v1.3 or v1.4
to v1.5.

## Prerequisites

**Terraform Version**

Use Terraform version `0.14.5` for this Orc8r upgrade.
Current Terraform scripts are *not* compatible with Terraform `0.15.0`.

**Helm Charts and Application Containers**

Have all application containers, and Helm charts built and published.
These should be from the head of the release branch.
You can follow the [Build Orchestrator](https://docs.magmacore.org/docs/orc8r/deploy_build)
guide for instructions on this process.

**Terraform State**

If you are using local Terraform state (the default), ensure all Terraform
state files are located in your working directory before proceeding.
As a check of this, `terraform show` should list existing state, rather than outputting `No state`.

**Kubernetes Cluster Access**

You should be able to access your Kubernetes cluster from the machine you are
using to initiate the upgrade process through the `kubectl` CLI.

If you cannot access your Kubernetes cluster through `kubectl`, set your
`KUBECONFIG` environment variable to point to your Orc8r cluster's kubeconfig
file.

**Kubernetes Cluster Version**

Have the Kubernetes version running on your cluster ready.
This version number will be used to set `cluster_version` later on during
the upgrade process.

If you are using AWS, and the AWS CLI is configured locally, run

```bash
aws eks describe-cluster --name orc8r
```

and note the `version` listed under `cluster`.

Otherwise, navigate to the AWS UI. Enter service `EKS` and select clusters on
the left hand side of the screen. Note the `Kubernetes version` for the
`orc8r` cluster.

## Upgrade Process

#### 1. NMS DB Data Migration

This process prepares your cluster for Orc8r-NMS DB unification. (In release 1.5 and going forward, the default DB type is Postgres for both NMS and Orc8r.)

*Paste in shell prompt:*

```bash
wget https://raw.githubusercontent.com/magma/magma/master/nms/app/packages/magmalte/scripts/fuji-upgrade/pre-upgrade-migration.sh && chmod +x pre-upgrade-migration.sh && ./pre-upgrade-migration.sh
```

The defaults options should likely work, unless the script cannot find the
right pods, or you have your cluster namespaced differently.

#### 2. Terraform main.tf Sanitization

The following variables should be removed from the main.tf Terraform file:
```
- nms_db_host
- nms_db_name
- nms_db_user
- nms_db_pass
```

If you are using a non-Postgres DB instance for your Orc8r setup, you will
need to modify `orc8r_db_dialect`.
This variable was added in 1.5, and defaults to `postgres`.
Common dialects will be `mysql`, `postgres`, and `mariadb`.
See the [sequelize docs](https://sequelize.org/v5/manual/dialects.html) for a list.

#### 3. Upgrade Terraform Modules

Set the `source` values for each of the Orchestrator modules in your root
Terraform module to point to this release's modules and bump chart and
container versions

```hcl-terraform
# This will likely be found in main.tf

module orc8r {
  source = "github.com/magma/magma//orc8r/cloud/deploy/terraform/orc8r-aws?ref=v1.5"
  # ...
  cluster_version = "1.17"   # set to Kubernetes version found above. 1.17 is used as an example here
}

module orc8r-app {
  source = "github.com/magma/magma//orc8r/cloud/deploy/terraform/orc8r-helm-aws?ref=v1.5"
  # ...
  orc8r_chart_version   = "1.5.20"
  orc8r_tag             = "MAGMA_TAG"  # from build step, e.g. v1.5.0
  orc8r_deployment_type = "fwa"        # valid options: ["fwa", "federated_fwa", "all"]
}
```

Set `cluster_version` to the Kubernetes version found during the
`Prerequisites` section. Bump your chart version to `1.5.20` and `orc8r_tag` to
the semver tag you published your new Orchestrator container images as.
You also need to set the `orc8r_deployment_type` variable to the deployment
type that you intend to deploy. This type sets which Orc8r modules will run.

Then, prepare Terraform for the upgrade

```bash
terraform init --upgrade
terraform refresh
```

#### 4. Terraform Apply

Terraform the upgrade.

Apply the upgrade

```bash
terraform apply     # Check this output VERY carefully!
```

After the Terraform command completes, all application components should be
successfully upgraded.

#### 5. (Optional) Post-Upgrade: Sync Pre-defined Alerts

A new `Duplicate Request Alert` has been added, which triggers when more than
100 alerts within the last 5 minutes have been raised.

To enable this alert, you must "Sync Pre-Defined Alerts" for _*every*_ network
defined in the Orc8r.
See the [NMS Alerts Guide](../nms/alerts#predefined-alerts) for guidance
on this process.

## Additional Details

#### NMS DB Data Migration

Before upgrading Orc8r to 1.5, the NMS DB migration needs to be completed.
This migration copies all NMS data to the same database housing Orc8r data.

These NMS tables are named
```
- Users
- Organizations
- FeatureFlags
- AuditLogEntry
- SequelizeMeta
```

A single script is provided to perform this migration for you, and should
tell you whether or not it was successful.
After migration, the tables listed above should be present in the Orc8r
database.

The script will ask for various inputs, and the defaults can be chosen unless
the script fails to find certain information, such as pod names.
It uses `kubectl` internally to run commands inside the pods.

#### New Pre-defined Alerts

"Duplicate request alert" is raised as a critical alert, when we receive more
than 100 alerts within last 5 minutes.
If you would like to customize this alert, refer to
[NMS alerts guide - custom alert rules](../nms/alerts#custom-alert-rules).

In order to start receiving these alarms, it is necessary to click on the
“Sync Pre-Defined alerts” in the NMS window for _*every*_ network defined in
the Orc8r. For more details on this process, see the [NMS alerts guide](../nms/alerts#predefined-alerts)
Syncing the alerts once on the NMS does not replicate the behavior to other
networks managed by the same NMS.
