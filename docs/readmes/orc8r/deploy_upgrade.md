---
id: deploy_upgrade
title: Upgrading from 1.0
hide_title: true
---
# Upgrading from 1.0

First, read through [Installing Orchestrator](deploy_install.md) to familiarize
yourself with the installation steps.

Create a new root Terraform module somewhere. Just like a fresh installation,
start by following the basic example root module in the `orc8r-helm-aws`
module. If you've already deployed the NMS, you can leave `deploy_nms` set to
`true`.

Pay careful attention to match the variables in your new root module to the
ones in your old one unless you specifically want to change some configuration.
If you change the instance type for EKS worker nodes, you will have to
terminate the existing instances one-by-one in the EC2 console after applying
Terraform and wait for autoscaling to replace the node.

## Migrating Old State

In the following instructions, `OLDTF` refers to the directory where you put
your existing root Terraform module and its state files, and `NEWTF` refers to
the directory where you put your new root Terraform module.

```bash
cd OLDTF
terraform state mv -state-out=NEWTF/terraform.tfstate 'module.vpc' 'module.orc8r.module.vpc'
terraform state mv -state-out=NEWTF/terraform.tfstate 'eks' 'module.orc8r.module.eks'
terraform state mv -state-out=NEWTF/terraform.tfstate 'aws_security_group.default' 'module.orc8r.aws_security_group.default'
terraform state mv -state-out=NEWTF/terraform.tfstate 'aws_db_instance.default' 'module.orc8r.aws_db_instance.default'
terraform state mv -state-out=NEWTF/terraform.tfstate 'aws_db_instance.nms' 'module.orc8r.aws_db_instance.nms'
```

Next, copy your old cert secrets onto your local disk if you've stored them
remotely. Then set the `seed_certs_dir` variable in the configuration of the
`orc8r-helm-aws` module to point to the location on disk where these certs
are.

Because the `helm_release` resource in Terraform doesn't yet support importing,
you can either continue to maintain the release outside of Terraform (not
recommended) or `helm delete --purge orc8r` before running `terraform apply` if
you can tolerate a few minutes of downtime.

At this point, you are ready to follow
[Installing Orchestrator](deploy_install.md). If you would like to migrate your
existing timeseries data, come back to the next section after you've finished
deploying the application again. Otherwise, you are good to go.

## Migrating Timeseries Data

Since 1.0, we've updated the storage solution for timeseries data from an EBS
volume mounted on a specific worker node to an EFS-based PersistentVolume. If
you want to keep your old timeseries data with this migration, you can do the
following:

1. Find the worker node that the new prometheus pod is running on
2. Attach the old metrics EBS volume to that worker node
3. Temporarily stop the prometheus pod (scale the deployment to 0)
4. Attach the new metricsdata EFS volume to the worker node
5. Manually `cp` the old data directory over to the new one by SSH'ing into the
worker node.
6. Scale the prometheus deployment back up to 1 replica

In most cases, we suggest throwing away the old data if you can because the
data migration procedure can take some time.
