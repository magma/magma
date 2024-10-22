---
id: version-1.7.0-deploy_terraform_options
title: Terraform Options
hide_title: true
original_id: deploy_terraform_options
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

```hcl
variable "orc8r_db_engine_version" {
  description = "Postgres engine version for Orchestrator DB."
  type        = string
  default     = "9.6.15"
}
```

### 3. Override variable

To override this value, go to your `main.tf` and define the variable in the `orc8r` module

```hcl
module "orc8r" {
  # ...
  orc8r_db_engine_version     = "12.6"
```

At this point, you can run `terraform init`, `terraform plan`, and `terraform apply` to proceed with the deployment.
