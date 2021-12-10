---
id: version-1.5.0-hold_terraform_state_on_s3
title: Hold Terraform state outside my local machine
hide_title: true
original_id: hold_terraform_state_on_s3
---
# Holding the terraform state outside of my local machine

**Description:** It is best to not hold the terraform state file locally on any disk that is susceptible to local failures and loss of state file. Note that the state file is the blueprint of the assets deployed on the public cloud environment and if this gets corrupt, it will be very difficult to manage cloud assets in a programmatic manner.

To aid this vision, terraform supports holding the state file remotely. One example is to use the AWS S3 bucket to store the state file. However, in such a scenario, it is plausible that more than one user may try to access the state file at the same time. For this reason, we can use a dynamo database to lock the database when it is being read and written to

**Environment:** Orc8r in AWS/EKS

**Components:** Orchestrator

**Steps:**

High level steps:

1. [Create an S3 bucket](https://docs.aws.amazon.com/AmazonS3/latest/user-guide/create-bucket.html)

2. [Create a dynamo dB state locking](https://www.terraform.io/docs/language/settings/backends/s3.html#dynamodb-state-locking)

3. Add configuration to the main.tf file under `module orc8r-app` to change the backend(variables `state_backend` and `state_config`)

```
state_backend = "s3"
state_config = {
    bucket = "s3_bucket_name"
    region = "us-west-2"
    dynamodb_table = "dynamo_table_name"
    key = "terraform/terraform.tfstate"
}
```

In the same main.tf file, add the below configuration:

```
output "nameservers" {
    value = module.orc8r.route53_nameservers
}
```

Similarly, at the top of the main.tf file, you can add the following per instructions in the [hashicorp instructions](https://www.terraform.io/docs/backends/types/s3.html#example-configuration)

```
terraform {
    required_version = ">= 0.12.0"
    backend "s3" {
    bucket = "s3_bucket_name"
    dynamodb_table = "dynamo_table_name"
    encrypt = true
    key="terraform/terraform.tfstate"
    region="us-west-2"
    }
}
```

Next time you run terraform init, you will be asked if you would like to migrate. Follow the instructions (starting at
step 6) [here](https://www.terraform.io/docs/cloud/migrate/index.html#step-6-run-terraform-init-to-migrate-the-workspace) to finish the migration.

At this point, you can go check your s3 bucket and validate that you have content written there and it should be
the state file for your deployment!
