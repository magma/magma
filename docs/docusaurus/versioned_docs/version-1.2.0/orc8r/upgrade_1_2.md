---
id: version-1.2.0-upgrade_1_2
title: Upgrade to v1.2
hide_title: true
original_id: upgrade_1_2
---

# Upgrade to v1.2

This guide covers upgrading Orchestrator deployments from v1.1 to v1.2.

First, read through [Installing Orchestrator](deploy_install.md) to familiarize
yourself with the installation steps.

We assume you've set up all the prerequisites, including developer tooling,
a Helm chart repository, and a container registry.

We'll upgrade to v1.2 in the following steps

- Upgrade to Helm 3
- Apply updated Terraform modules
- Run Orchestrator data migrations
- Validate and clean up

## Prerequisites

Please strongly consider taking a DB snapshot before attempting the upgrade.
The v1.1 to v1.2 upgrade does not support Helm-only downgrades. That is,
once you perform the Orchestrator data migrations, v1.1 deployments will lose
access to a subset of the migrated data.

If you are using local Terraform state (the default), ensure all Terraform state files (i.e. [`terraform.tfstate`](https://www.terraform.io/docs/state/index.html)) are located in your working directory before proceeding. This means `terraform show` should list existing state (rather than outputting `No state`).
 
## Helm 3 Upgrade

Orchestrator v1.2 requires an upgrade from Helm 2 to Helm 3. Helm provides a
[migration guide](https://helm.sh/blog/migrate-from-helm-v2-to-helm-v3/), which
we'll follow below.

NOTE: for each Helm `2to3` command, consider adding the `--dry-run` switch to
double-check the output meets your expectations.

### Install Helm 3

We'll assume you currently have Helm 2 installed as `helm`. To start, install
Helm 3 as `helm3`

```bash
$ brew install helm
$ alias helm3=/usr/local/Cellar/helm/3.3.1/bin/helm  # or wherever v3 is

$ helm version
Client: &version.Version{SemVer:"v2.16.1", ... }

$ helm3 version
version.BuildInfo{Version:"v3.3.1", ... }
```

### Update Helm config

Next let's update the local Helm configuration

```bash
$ helm3 repo list  # initially empty

Error: no repositories to show

$ helm3 2to3 move config

2020/08/31 21:49:31 WARNING: Helm v3 configuration may be overwritten during this operation.
2020/08/31 21:49:31 [Move Config/confirm] Are you sure you want to move the v2 configuration? [y/N]: y
2020/08/31 21:49:33 Helm v2 configuration will be moved to Helm v3 configuration.
...
2020/08/31 21:49:33 Helm v2 configuration was moved successfully to Helm v3 configuration.


$ helm3 repo list  # now correctly populated

NAME            URL
stable          https://kubernetes-charts.storage.googleapis.com
incubator       http://storage.googleapis.com/kubernetes-charts-incubator
YOUR_REPO       YOUR_REPO_URL
```

### Upgrade Helm deployments

Finally, upgrade existing Helm releases

```bash
$ for release in $(helm list --short) ; do helm3 2to3 convert ${release} ; done

2020/09/01 00:51:35 Release "efs-provisioner" will be converted from Helm v2 to Helm v3.
...
2020/09/01 00:51:36 Release "efs-provisioner" was converted successfully from Helm v2 to Helm v3.
2020/09/01 00:51:36 Release "external-dns" will be converted from Helm v2 to Helm v3.
...
2020/09/01 00:51:36 Release "external-dns" was converted successfully from Helm v2 to Helm v3.
2020/09/01 00:51:36 Release "orc8r" will be converted from Helm v2 to Helm v3.
...
2020/09/01 00:51:37 Release "orc8r" was converted successfully from Helm v2 to Helm v3.

$ helm list  # Helm 2

NAME            REVISION    UPDATED                                           STATUS    CHART                   APP VERSION NAMESPACE
efs-provisioner 1           Mon Aug 31 20:43:01 2020                          DEPLOYED  efs-provisioner-0.11.0  v2.4.0      kube-system
external-dns    1           Mon Aug 31 20:43:01 2020                          DEPLOYED  external-dns-2.19.1     0.6.0       kube-system
orc8r           2           Mon Aug 31 21:26:29 2020                          DEPLOYED  orc8r-1.4.21            1.0         orc8r

$ helm3 list -n orc8r && helm3 list -n kube-system  # Helm 3

NAME            NAMESPACE   REVISION  UPDATED                                 STATUS    CHART                   APP VERSION
orc8r           orc8r       2         2020-09-01 07:48:19.803722138 +0000 UTC deployed  orc8r-1.4.21            1.0

NAME            NAMESPACE   REVISION  UPDATED                                 STATUS    CHART                   APP VERSION
efs-provisioner kube-system 1         2020-09-01 07:39:37.87374189 +0000 UTC  deployed  efs-provisioner-0.11.0  v2.4.0
external-dns    kube-system 1         2020-09-01 07:39:38.45646776 +0000 UTC  deployed  external-dns-2.19.1     0.6.0
```

Note [one of the Helm 3 changes](https://v3.helm.sh/docs/faq/#changes-since-helm-2)
is that Helm 3 release names are now scoped to their K8s namespace.

## Terraform Apply

NOTE: this step involves a `terraform apply`. Please triple check there is
*nothing* related to RDS in the plan (`aws_db_instance` resources). All the
application components are stateless so any mistakes while updating the EKS
cluster are recoverable, but if you drop your RDS database instances you could
end up with unrecoverable data loss.

First, update your `main.tf` file to pull in the v1.2 changes

```hcl-terraform
module orc8r {
  source = "github.com/magma/magma//orc8r/cloud/deploy/terraform/orc8r-aws?ref=v1.2"
  # ...
}

module orc8r-app {
  source = "github.com/magma/magma//orc8r/cloud/deploy/terraform/orc8r-helm-aws?ref=v1.2"
  # ...
  orc8r_chart_version = "1.4.35"
  orc8r_tag           = "MAGMA_TAG"  # from build step
}
```

Refresh Terraform state, then apply the new changes. The `apply` step should
only contain `changes`, no `additions`

```bash
terraform init --upgrade
terraform refresh
terraform apply  # DOUBLE CHECK this output
```

## Data Migrations

We updated the DB schemas for a few services since v1.1. We'll run three
manual migrations to migrate the data.

NOTE: running these scripts is the point of commitment. The v1.1 deployment
will continue to function, and 95% of its functionality will be retained, but
after the non-reversible schema change you will need to upgrade to v1.2 or fall
back to a DB checkpoint.

```bash
$ export CNTLR_POD=$(kubectl get pod -l app.kubernetes.io/component=controller -o jsonpath='{.items[0].metadata.name}')
$ kubectl exec -it ${CNTLR_POD} -- bash

(pod)$ cd /var/opt/magma/bin

(pod)$ ./m010_default_apns -verify

I0729 04:56:41.499084     462 main.go:130] BEGIN MIGRATION
I0729 04:56:41.516114     462 main.go:231] [RUN] INSERT INTO cfg_entities (pk,network_id,type,"key",graph_id,config) VALUES ($1,$2,$3,$4,$5,$6) [ ... ]
I0729 04:56:41.518006     462 main.go:369] [RUN] INSERT INTO cfg_assocs (from_pk,to_pk) VALUES ($1,$2),($3,$4),($5,$6),($7,$8),($9,$10),($11,$12) [ ... ]
I0729 04:56:41.519320     462 main.go:400] [RUN] UPDATE cfg_entities SET graph_id = $1 WHERE (network_id = $2 AND (graph_id = $3 OR graph_id = $4 OR graph_id = $5 OR graph_id = $6 OR graph_id = $7 OR graph_id = $8)) [ ... ]
...
I0729 04:56:42.024077     462 main.go:437] Subscriber {some_network subscriber IMSI0123456789 <nil> 00665aeb-968e-4319-8dc9-260647a4105b [apn-oai.ipv4] [] 1} has APN assocs [oai.ipv4]
I0729 04:56:42.024110     462 main.go:437] Subscriber {some_network subscriber IMSI0123456789 <nil> 00665aeb-968e-4319-8dc9-260647a4105b [apn-oai.ipv4] [] 1} has APN assocs [oai.ipv4]
I0729 04:56:42.031600     462 main.go:175] SUCCESS
I0729 04:56:42.031632     462 main.go:176] END MIGRATION

(pod)$ ./m011_subscriber_config

I0901 10:41:35.613014     267 main.go:56] Subscriber migration successfully completed

(pod)$ ./m012_policy_qos

I0901 10:42:27.457720     274 main.go:97] BEGIN MIGRATION
I0901 10:42:27.487736     274 configurator.go:116] Flipping 2 assocs
...
I0901 10:42:27.489151     274 main.go:111] SUCCESS
I0901 10:42:27.489168     274 main.go:112] END MIGRATION
```

## Validate and Clean Up

To validate, visit the NMS and check out some subscribers. Ensure everything
passes the eye test.

Once you're satisfied with your v1.2 deployment, clean up the outdated Helm 2
state

```bash
$ helm3 2to3 cleanup

WARNING: "Helm v2 Configuration" "Release Data" "Tiller" will be removed.
This will clean up all releases managed by Helm v2. It will not be possible to restore them if you haven't made a backup of the releases.
Helm v2 may not be usable afterwards.

[Cleanup/confirm] Are you sure you want to cleanup Helm v2 data? [y/N]: y
2020/09/01 02:26:04 Helm v2 data will be cleaned up.
...
2020/09/01 02:26:07 Helm v2 data was cleaned up successfully.
```

Finally, make Helm 3 the default Helm executable

```bash
brew unlink helm@2 && brew link helm  # for example
```
