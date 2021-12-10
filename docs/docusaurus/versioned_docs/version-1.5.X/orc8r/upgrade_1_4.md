---
id: version-1.5.0-upgrade_1_4
title: Upgrade to v1.4
hide_title: true
original_id: upgrade_1_4
---

# Upgrade to v1.4

This guide covers upgrading Orchestrator deployments from v1.3 to v1.4.
Upgrades which skip versions (e.g. v1.2 to v1.4) are not explicitly documented
or supported at this time.

The v1.4 upgrade follows a standard procedure:

- Upgrade Terraform modules
- Bump version numbers

## Prerequisites

Build and publish the application containers and Helm charts on the head of the
release branch by following the [Build Orchestrator](https://docs.magmacore.org/docs/orc8r/deploy_build)
documentation.

If you are using local Terraform state (the default), ensure all Terraform
state files are located in your working directory before proceeding. This means
`terraform show` should list existing state, rather than outputting `No state`.

Now fetch the Kubernetes version that your cluster is running. If the AWS CLI
is configured locally, run

```
aws eks describe-cluster --name orc8r
```

and note the `version` listed under `cluster`.

Otherwise, navigate to the AWS UI. Enter service `EKS` and select clusters on
the left hand side of the screen. Note the `Kubernetes version` for the
`orc8r` cluster. Use this version to set `cluster_version` below.

## Upgrade Terraform modules

Set the `source` values for each of the Orchestrator modules in your root
Terraform module to point to this release's modules and bump chart and
container versions

```hcl-terraform
# This will likely be found in main.tf

module orc8r {
  source = "github.com/magma/magma//orc8r/cloud/deploy/terraform/orc8r-aws?ref=v1.4"
  # ...
  cluster_version = "1.17"   # set to Kubernetes version found above. 1.17 is used as an example here
}

module orc8r-app {
  source = "github.com/magma/magma//orc8r/cloud/deploy/terraform/orc8r-helm-aws?ref=v1.4"
  # ...
  orc8r_chart_version   = "1.5.16"
  orc8r_tag             = "MAGMA_TAG"  # from build step, e.g. v1.4.0
  orc8r_deployment_type = "fwa"        # valid options: ["fwa", "federated_fwa", "all"]
}
```

Set `cluster_version` to the Kubernetes version found during the
`Prerequisites` section. Bump your chart version to `1.5.16` and `orc8r_tag` to
the semver tag you published your new Orchestrator container images as.
You also need to set the `orc8r_deployment_type` variable to the deployment
type that you intend to deploy. This type sets which orc8r modules will run.

Then, prepare Terraform for the upgrade

```bash
terraform init --upgrade
terraform refresh
```

## Terraform apply

Terraform the upgrade. We'll also need to manually rescale the Prometheus
deployment, as a workaround for quirks in how its deployment handles upgrading
(for context, see [issue #4580](https://github.com/magma/magma/issues/4580)).

Apply the upgrade

```bash
terraform apply     # Check this output VERY carefully!
```

Once the above command indicates it's applying the Helm upgrade, switch to
a separate terminal tab and manually scale the Prometheus deployment. You can
re-run the Terraform command if it times out.

```bash
kubectl --namespace orc8r scale deployment orc8r-prometheus --replicas=0 && kubectl --namespace orc8r scale deployment orc8r-prometheus --replicas=1
```

After the Terraform command completes, all application components should be
successfully upgraded.
