---
id: version-1.5.0-override_values_of_terraform_files
title: Override values of Terraform files
hide_title: true
original_id: override_values_of_terraform_files
---
# Override values of Terraform files

**Description:** This guide describes how to override values that are part of the terraform files that are used in the github repository.

**Environment:** Orc8r in AWS/EKS

**Components:** Orchestrator

**Steps:**

The main.tf file is the location where we would override any element that we do not want to inherit from the
source files located on github.

Let’s look at the following example:

Let’s assume that we want to control the version of the PostGREs database that is deployed for the orc8r-app module and lock it to version 9.6.18.

If we look inside `magma/orc8r/cloud/deploy/terraform/orc8r-aws/db.tf`

```
resource "aws_db_instance" "default" {
  identifier        = var.orc8r_db_identifier
  allocated_storage = var.orc8r_db_storage_gb
  engine            = "postgres"
  engine_version    = var.orc8r_db_engine_version
  instance_class    = var.orc8r_db_instance_class
```

We note that the “engine_version” is actually a variable. This variable is defined in the variables.tf file in the same directory as follows.

```
variable "orc8r_db_engine_version" {
  description = "Postgres engine version for Orchestrator DB."
  type        = string
  default     = "9.6.15"
}
```

According to this construct, if this value is not specified, it defaults to “9.6.15”. To override this value, go to `main.tf` and define the variable in the appropriate module. In this example, add `orc8r_db_engine_version = “9.6.18` under the orc8r module.

At this point, you can run terraform init, terraform plan, terraform apply to proceed with the deployment!

```
module "orc8r" {
  # Change this to pull from github with a specified ref
  source = "../../../orc8r-aws"

  region = "us-west-2"

  nms_db_password             = "mypassword" # must be at least 8 characters
  orc8r_db_password           = "mypassword" # must be at least 8 characters
  secretsmanager_orc8r_secret = "orc8r-secrets"
  orc8r_domain_name           = "orc8r.example.com"
  orc8r_db_engine_version     = "9.6.18"
```

At this point, you can run `terraform init`, `terraform plan`, `terraform apply` to proceed with the deployment!
