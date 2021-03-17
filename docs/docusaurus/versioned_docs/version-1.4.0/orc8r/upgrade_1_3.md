---
id: version-1.4.0-upgrade_1_3
title: Upgrade to v1.3
hide_title: true
original_id: upgrade_1_3
---

# Upgrade to v1.3

This guide covers upgrading Orchestrator deployments from v1.2 to v1.3.
Upgrades which skip versions (e.g. v1.1 to v1.3) are not explicitly documented
or supported at this time.

The v1.3 upgrade follows a standard procedure:

- Upgrade the Terraform modules
- Bump version numbers
- Run data migrations

## Prerequisites

There are quite a few data migrations to run after this upgrade. We recommend
taking a DB snapshot of your RDS instances before continuing so that there is
a path to downgrade.

Build and publish the application containers on the head of the release branch
by following the documentation, and package and upload version `1.4.36` of the
orc8r Helm chart as well.

If you are using local Terraform state (the default), ensure all Terraform state files (i.e. [`terraform.tfstate`](https://www.terraform.io/docs/state/index.html)) are located in your working directory before proceeding. This means `terraform show` should list existing state (rather than outputting `No state`).

## Upgrade Terraform Modules

Set the `source` values for each of the Orchestrator modules in your root
Terraform module to point to this release's modules and bump chart and
container versions:

```hcl-terraform
module orc8r {
  source = "github.com/magma/magma//orc8r/cloud/deploy/terraform/orc8r-aws?ref=v1.3"
  # ...
}

module orc8r-app {
  source = "github.com/magma/magma//orc8r/cloud/deploy/terraform/orc8r-helm-aws?ref=v1.3"
  # ...
  orc8r_chart_version = "1.4.36"
  orc8r_tag           = "MAGMA_TAG"  # from build step, e.g. v1.3.0
}
```

You should bump your chart version from `1.4.35` to `1.4.36` and `orc8r_tag` to
whatever semver tag you published your new Orchestrator container images as.

IMPORTANT: check the "Build Orchestrator" section for the Git hash to check out
before packaging the orc8r Helm chart. You must deploy `1.4.36`, not `1.4.37`
on the head of the release branch.

Then,

```bash
terraform init --upgrade
terraform refresh
```

## Terraform Apply

Before upgrading the deployment, manually scale the Prometheus pods to 0
replicas, otherwise the Prometheus version upgrade in this release will fail:

```bash
kubectl scale deployment orc8r-prometheus --replicas=0

terraform apply     # Check this output VERY carefully!
```

After applying Terraform changes, all your application components should be
upgraded and Prometheus should be back up on v2.20.1.

## Data Migrations

> **_NOTE:_** If you're upgrading to release tag v1.3.0 specifically, `m014_enodeb_config` is not required.

```
$ export CNTLR_POD=$(kubectl get pod -l app.kubernetes.io/component=controller -o jsonpath='{.items[0].metadata.name}')
$ kubectl exec -it ${CNTLR_POD} -- bash

(pod)$ cd /var/opt/magma/bin

(pod)$ ./m013_relay_split

...

(pod)$ ./m013_policy_ipv6

...

(pod)$ ./m014_enodeb_config

...
```
