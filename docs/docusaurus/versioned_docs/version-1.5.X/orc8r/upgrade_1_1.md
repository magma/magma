---
id: version-1.5.0-upgrade_1_1
title: Upgrade to v1.1
hide_title: true
original_id: upgrade_1_1
---

# Upgrade to v1.1

This guide covers upgrading Orchestrator deployments from v1.0 to v1.1.

First, read through [Installing Orchestrator](deploy_install.md) to familiarize
yourself with the installation steps. If you want to perform an online upgrade
(i.e. no downtime on NMS), this guide will walk you through the process of
concurrently deploying the 1.1.0 version of Orchestrator and NMS to your EKS
cluster. You can flip your DNS records to the new application whenever you feel
comfortable to complete the migration.

This guide will assume that you've already set up all the prerequisites,
including developer tooling, a Helm chart repository, and a container registry.

## Create a New Root Module

First, create a new directory somewhere to store your new root Terraform module
for the 1.1.x deployment. We have an example root module at https://github.com/facebookincubator/magma/tree/v1.1/orc8r/cloud/deploy/terraform/orc8r-helm-aws/examples/online-upgrade
that we recommend you use for the upgrade. Copy all the files to your new
directory and change the `source` attribute of both modules in `main.tf` to
`github.com/facebookincubator/magma//orc8r/cloud/deploy/terraform/orc8r-aws` and
`github.com/facebookincubator/magma//orc8r/cloud/deploy/terraform/orc8r-helm-aws`,
respectively.

Once you've got this root module set up and your variables defined, run
`terraform init` in this directory.

## Migrate Old Terraform State

In the following instructions, `OLDTF` refers to the directory where you put
your existing root Terraform module and its state files, and `NEWTF` refers to
the directory where you put your new root Terraform module. If you are storing
your Terraform state remotely, remember to `terraform pull` before running
`terraform state mv`.

Your new Terraform root module needs to know about the infrastructure
components that you created for v1.0 so it doesn't create new copies. Terraform
has a useful utility `terraform state mv` that can help accomplish this:

```bash
cd OLDTF
terraform state mv -state-out=NEWTF/terraform.tfstate 'module.vpc' 'module.orc8r.module.vpc'
terraform state mv -state-out=NEWTF/terraform.tfstate 'aws_security_group.default' 'module.orc8r.aws_security_group.default'
terraform state mv -state-out=NEWTF/terraform.tfstate 'aws_ebs_volume.prometheus-ebs-eks' 'aws_ebs_volume.prometheus-ebs-eks'
terraform state mv -state-out=NEWTF/terraform.tfstate 'aws_ebs_volume.prometheus-configs-ebs-eks' 'aws_ebs_volume.prometheus-configs-ebs-eks'
terraform state mv -state-out=NEWTF/terraform.tfstate 'aws_iam_policy.worker_node_policy' 'aws_iam_policy.worker_node_policy'
terraform state mv -state-out=NEWTF/terraform.tfstate 'aws_db_instance.default' 'module.orc8r.aws_db_instance.default'
terraform state mv -state-out=NEWTF/terraform.tfstate 'aws_db_instance.nms' 'module.orc8r.aws_db_instance.nms'
terraform state mv -state-out=NEWTF/terraform.tfstate 'module.eks' 'module.orc8r.module.eks'
terraform state mv -state-out=NEWTF/terraform.tfstate 'data.template_file.metrics_userdata' 'data.template_file.metrics_userdata'
```

If you added any custom components to your v1.0 root Terraform module, you
should copy the resource blocks over to your new root module and use
`terraform state mv` to let Terraform know that they already exist.

## Migrate Secrets

One change we've made in the v1.1 installation procedure is to store all the
application certificates in AWS Secretsmanager instead of on-disk in the
`secrets` sub-chart. Simply copy your old application secrets (which probably
live under the `secrets` orc8r subchart as a `.secrets` directory) to a new
temporary location on disk.

Then, update the application certs to include 2 new components (replace
YOURDOMAIN.COM with the domain you've reserved for Orchestrator):

```bash
cd MYSECRETS/certs
openssl genrsa -out fluentd.key 2048
openssl req -x509 -new -nodes -key fluentd.key -sha256 -days 3650 \
    -out fluentd.pem -subj "/C=US/CN=fluentd.YOURDOMAIN.COM"
```

## Define Terraform Variables

The variables that you can define in your `vars.tfvars` are documented in the
README.md in your new root Terraform module. In addition to `vars.tfvars`,
if you changed the default worker node configuration for your EKS cluster when
deploying v1.0, update that accordingly in `main.tf`. Most of this
configuration should match your v1.0 Terraform - we are aiming for a smooth
import of existing components.

Importantly, `seed_certs_dir` needs to point the the `certs` subdirectory
inside the temporary secrets directory you created above. At this point it is
actually safe to trash the other subdirectories (`envdir` and `configs`).

If you changed the EKS worker group configuration in your v1.0 deployment,
also update `eks_worker_groups` in `main.tf` to match.

## Initial Terraform

```bash
terraform plan -target=module.orc8r -var-file=vars.tfvars
```

Pay VERY close attention to the output of the plan to make sure that nothing
unexpected is getting deleted. If you have any questions about Terraform's
proposed plan, please drop us a question on the mailing list and we can take a
look. A misconfiguration could result in downtime or some complicated recovery
procedures.

For reference, our Terraform plan for this step ended up with:

```bash
Plan: 21 to add, 5 to change, 5 to destroy.
```

with Elasticsearch enabled. The 5 to change/destroy were mostly launch
configurations for the EKS worker groups. If you find yourself with something
dramatically different, check over your `vars.tfvars` to make sure everything
lines up with your v1.0 Terraform configuration.

Most importantly, you should triple check that there is *nothing* related to
RDS in the plan (`aws_db_instance` resources). All the application components
are stateless so any mistakes while updating the EKS cluster are recoverable,
but if you drop your RDS database instances you could end up with unrecoverable
data loss.

When you are convinced that your new Terraform module won't break anything:

```bash
terraform apply -target=module.orc8r -var-file=vars.tfvars
```

## Application Terraform

We will deploy the v1.1 application concurrently with the v1.0 application,
just in a different namespace and under a different Helm deployment name. At
this point, the deployment procedure is pretty much like the from-scratch
installation:

```bash
$ terraform apply -target=module.orc8r-app.null_resource.orc8r_seed_secrets -var-file=vars.tfvars

Apply complete! Resources: 1 added, 0 changed, 0 destroyed.

$ terraform apply -var-file=vars.tfvars

Apply complete! Resources: 16 added, 0 changed, 0 destroyed.
```

At this point, you should see all the v1.1 application pods in the namespace
that you chose for the upgraded deployment.

## Data Migrations

We updated the DB schemas for a few services since v1.0. You'll have to run a
pair of manual migrations to migrate the data. These scripts are idempotent
and will not affect the v1.0 deployment.

```bash
# Replace orc8r with your v1.1 k8s namespace if you changed the name
$ export CNTLR_POD=$(kubectl --namespace orc8r get pod -l app.kubernetes.io/component=controller -o jsonpath='{.items[0].metadata.name}')
$ kubectl --namespace orc8r exec -it ${CNTLR_POD} bash

(pod)$ cd /var/opt/magma/bin
(pod)$ ./m005_certifier_to_blobstore -verify

... 149 main.go:49] BEGIN MIGRATION
... 149 datastore_to_blobstore.go:136] [RUN] INSERT INTO certificate_info_blobstore (network_id,type,"key",value,version) SELECT ('placeholder_network') AS network_id, ('certificate_info') AS type, "key", value, generation_number FROM certificate_info_db
... 149 datastore_to_blobstore.go:84] SUCCESS
... 149 main.go:53] END MIGRATION
... 149 main.go:85] [manually verify] serial number count: 42
... 149 main.go:97] [manually verify] key-value pair: {key: REDACTED, value: id:<gateway:<hardware_id:"redacted-1234-1234-1234-redacted" > > not_before:<seconds:42 nanos:42 > not_after:<seconds:42 nanos:42 > }
...

(pod)$ ./m008_accessd_to_blobstore -verify

... 156 main.go:49] BEGIN MIGRATION
... 156 datastore_to_blobstore.go:136] [RUN] INSERT INTO access_control_blobstore (network_id,type,"key",value,version) SELECT ('placeholder_network') AS network_id, ('access_control') AS type, "key", value, generation_number FROM access_control
... 156 datastore_to_blobstore.go:84] SUCCESS
... 156 main.go:52] END MIGRATION
... 156 main.go:85] [manually verify] number of operators: 1
... 156 main.go:97] [manually verify] operator-acl pair: {operator: operator:"admin_operator" , acl: operator:<operator:"admin_operator" > entities:<key:"Id_Wildcard_Gateway" value:<id:<wildcard:<> > permissions:3 > > entities:<key:"Id_Wildcard_Network" value:<id:<wildcard:<type:Network > > permissions:3 > > entities:<key:"Id_Wildcard_Operator" value:<id:<wildcard:<type:Operator > > permissions:3 > > }
```

These 2 CLIs will spit out a few records from the migrated tables as output.
Verify that the Go structs in the output don't have all zeroed values as a
sanity check. The exact count of migrated records will differ depending on
your specific deployment but it should be nonzero for both commands.

### First-Time NMS Setup

With the new multi-tenancy support in the NMS introduced in v1.1 you have to
create a new admin user in the `master` organization to set up access for
other tenants:

```bash
kubectl --namespace orc8r exec -it \
    $(kubectl --namespace orc8r get pod -l app.kubernetes.io/component=magmalte -o jsonpath='{.items[0].metadata.name}') -- \
    yarn setAdminPassword master ADMIN_USER_EMAIL ADMIN_USER_PASSWORD
```

When you flip DNS over to the services in the v1.1 namespace, you'll be able to
set up access for your NMS tenants at `https://master.nms.yourdomain.com`.

You will have to reprovision accounts for your existing users using
organizations as the old `nms.yourdomain.com` URL will no longer resolve to the
NMS. If you have a lot of users to migrate, you can do this using SQL directly
on the backing NMS MySQL database but you may find it simpler to create new
accounts for all your existing users in the frontend instead.

See the NMS user guides to understand how to set up tenants and users in the
new NMS.

### Migrating Timeseries Data and Configuration

We used an EBS volume mounted to a specific worker node to persist Prometheus
configuration and data in 1.0. In v1.1 we are now using PersistentVolume and
EFS to store this data so we can align with the "no pets" philosophy of k8s.

This does mean that if you want to keep your timeseries data from v1.0, you'll
have to migrate it by hand to the v1.1 EFS volume. We don't have an exact guide
for you to follow, but if you really want to keep this data, we successfully
migrated our old timeseries data with this procedure:

1. `kubectl scale` the Prometheus pods in both namespaces to 0
2. Attach and mount the EFS volume for Prometheus data in v1.1 to the `metrics`
EC2 worker node in the EC2 console or using the `aws` CLI on the node
3. SSH into the `metrics` EC2 worker node then `rsync` the Prometheus data
from the EBS volume (mounted at `/prometheusData`) to the EFS volume
4. `kubectl scale` the Prometheus pods back up to 1 in both namespaces

We observed a transfer rate around 3MB/sec with this procedure. You will lose
timeseries data for the duration of the data transfer.

If you've created a nontrivial number of alerts using the NMS, you can also
migrate those manually over to the new namespace. Use `kubectl cp` to move
all the files under `/etc/alertmanager` of the `orc8r-alertmanager` pod in 1.0
to your local host then up to the same directory of the same pod in the 1.1
namespace. You may find it more straightforward to recreate all the alerts in
the frontend instead.

## Flipping the Switch

The v1.0.x LTE AGWs are compatible with the v1.1 Orchestrator, so you can
simply swap your DNS records over to the new application. If things don't look
good, just change the DNS records back.

An important change we made in 1.1 was using AWS Route53 to automatically
resolve DNS to appropriate k8s services. So instead of CNAME records in your
registrar for your domain/subdomain, use NS records instead and set the list
of nameservers to the `nameservers` output from the root Terraform module.
Because NS records don't support wildcards, you'll have to add new records
for the following subdomains:

- `master.nms.` to access the admin UI to create and manage NMS tenants
- `<tenant>.nms.` for each of your tenants (replace `<tenant>` with the
organization ID you set up in the NMS)

You can also remove the old CNAME record for the `nms.` subdomain, that will
no longer resolve to the application (all access must be through an
organization).

See the NMS user guides to understand how to set up tenants in the new
application.

## Cleaning Up

When you are satisfied with your v1.1 deployment and you've upgraded all your
LTE AGW's in the field, you can safely purge and clean up the old v1.0
components.

To delete the v1.0 application components, `helm delete --purge <name>`. You
can safely `kubectl delete namespace <name>` after this (make sure you pick
the namespace that you used for v1.0, not v1.1).

For the online upgrade we spun up new worker nodes to handle the extra load of
deploying 2 application versions at the same time. You can scale your EKS
worker groups back down once you purge the legacy deployment. Since we are
using only cattle in v1.1, you can delete the extra kubelet args and the
custom metrics userdata as well. For a deployment handling up to 100 LTE AGWs,
3x `t3.large` instances will do the job just fine.

You can safely delete all the Terraform resources in `legacy.tf` at this point
as well. After this, your root Terraform module should look a lot like
`examples/basic` in the `orc8r-helm-aws` Terraform module in Magma.
