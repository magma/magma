---
id: deploy_terraform_options
title: Terraform Options
hide_title: true
---

# Terraform Options

This page describes a set of useful Terraform options, and how to apply them to your Orc8r deployment.

## [Store Terraform state in AWS](https://www.terraform.io/docs/language/settings/backends/s3.html)

Storing [Terraform state](https://www.terraform.io/docs/language/state/index.html) locally can cause data loss and team-wise synchronization issues. Terraform recommends storing this state remotely when possible.

This walkthrough covers storing Terraform state in an AWS S3 bucket, with synchronization provided by AWS DynamoDB.

### 1. [Create an S3 bucket](https://docs.aws.amazon.com/AmazonS3/latest/user-guide/create-bucket.html)

### 2. [Create a DynamoDB instance](https://docs.aws.amazon.com/amazondynamodb/latest/developerguide/Introduction.html) (optional)

### 3. [Update Terraform configuration](https://www.terraform.io/docs/backends/types/s3.html#example-configuration)

Add configuration to the main Terraform module file (likely `main.tf` file)

```hcl
terraform {
  required_version = ">= 0.15.0"
  backend "s3" {
    bucket = "YOUR_S3_BUCKET_NAME"
    dynamodb_table = "YOUR_DYNAMODB_TABLE_NAME"
    encrypt = true
    key = "terraform/terraform.tfstate"
    region = "us-west-2"  # or chosen region
  }
}

# ...

module "orc8r-app" {
  state_backend = "s3"
  state_config = {
    # Contents should match block above
    bucket = "YOUR_S3_BUCKET_NAME"
    dynamodb_table = "YOUR_DYNAMODB_TABLE_NAME"
    encrypt = true
    key = "terraform/terraform.tfstate"
    region = "us-west-2"
  }
  # ...
}

# ...
```

### 4. Migrate

Next time you run `terraform init`, you will be asked if you'd like to migrate. Follow the [migration instructions](https://www.terraform.io/docs/cloud/migrate/index.html#step-6-run-terraform-init-to-migrate-the-workspace) to complete the migration.

Finally, check your S3 bucket is non-empty, meaning the Terraform state has been successfully written.

## Override Terraform module values

The Orc8r Terraform modules are configurable.

This walkthrough provides an example of how to override the Orc8r Terraform module variables -- specifically, to override the deployed Postgres version.

### 1. Introduction

The root Terraform module (i.e. your `main.tf` file) is used to override variables from imported modules (i.e. Orc8r modules).

### 2. Locate variable

Let's assume we want to install Postgres 12.6. This infra-related value is managed by the `orc8r` module.

Checking the contents of [orc8r-aws/variables.tf](https://github.com/magma/magma/blob/master/orc8r/cloud/deploy/terraform/orc8r-aws/variables.tf), we'll find the following block

```
variable "orc8r_db_engine_version" {
  description = "Postgres engine version for Orchestrator DB."
  type        = string
  default     = "9.6.15"
}
```

### 3. Override variable

To override this value, go to your `main.tf` and define the variable in the `orc8r` module

```
module "orc8r" {
  # ...
  orc8r_db_engine_version     = "12.6"
```

At this point, you can run `terraform init`, `terraform plan`, and `terraform apply` to proceed with the deployment.

## Managing Scaled Orc8r Deployments

In our default deployments, the current recommendation is to bring up orc8r instances(eks_worker_groups) as “t3.large” and DB instance(orc8r_db_instance_class) as “db.m4.large” instance. These default recommendations works well for small to medium deployments (< 50 gateways) and (< 10k subscribers). 

For medium large to high scale deployments i.e > 50 gateways and > 10k subscribers we recommend the following.

### RDS

Setting ‘orc8r_db_instance_class’ to db.m4.xlarge

### Prometheus

Prometheus also requires an instance with larger memory footprint when number of metrics scraped increases. We would recommend using t3.xlarge instance and pinning prometheus service with a larger node instance.

```
--- a/orc8r/cloud/deploy/terraform/orc8r-aws/variables.tf
+++ b/orc8r/cloud/deploy/terraform/orc8r-aws/variables.tf
@@ -102,12 +102,21 @@ variable "eks_worker_groups" {
 {
 name = "wg-1"
 instance_type = "t3.large"
- asg_desired_capacity = 3
+ asg_desired_capacity = 2
 asg_min_size = 1
- asg_max_size = 3
+ asg_max_size = 2
 autoscaling_enabled = false
 kubelet_extra_args = "" // object types must be identical (see thanos_worker_groups)
 },
+ {
+ name = "wg-2"
+ instance_type = "t3.xlarge"
+ asg_desired_capacity = 1
+ asg_min_size = 1
+ asg_max_size = 1
+ autoscaling_enabled = false
+ kubelet_extra_args = "" // object types must be identical (see thanos_worker_groups)
+ },
 ]
 }
 diff --git a/orc8r/cloud/deploy/terraform/orc8r-helm-aws/templates/orc8r-values.tpl b/orc8r/cloud/deploy/terraform/orc8r-helm-aws/templates/orc8r-values.tpl
index d3b2e3837..f6a07c9e8 100644
--- a/orc8r/cloud/deploy/terraform/orc8r-helm-aws/templates/orc8r-values.tpl
+++ b/orc8r/cloud/deploy/terraform/orc8r-helm-aws/templates/orc8r-values.tpl
@@ -83,6 +83,8 @@ metrics:
     includeOrc8rAlerts: true
     prometheusCacheHostname: ${prometheus_cache_hostname}
     alertmanagerHostname: ${alertmanager_hostname}
+    nodeSelector:
+      node.kubernetes.io/instance-type: t3.xlarge

   alertmanager:
     create: true
 
```

**Snapshotting**
curl -XPOST http://localhost:9090/api/v1/admin/tsdb/snapshot

**Delete Series**
https://prometheus.io/docs/prometheus/latest/querying/api/

### Elasticsearch Instance

Recommend tuning “elasticsearch_ebs_volume_size” and “elasticsearch_disk_threshold” to ensure that elasticsearch can deal with volume of events and logs. Currently elasticsearch curator cleans up indices based on age and space consumed. We would highly recommend customizing elasticsearch_curator for your own deployment and ensuring that important logs are backed up. 

**Snapshotting**
https://www.elastic.co/guide/en/elasticsearch/reference/current/snapshot-restore.html#snapshot-restore


