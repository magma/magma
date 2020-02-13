# Example Terraform for Orchestrator EKS

The contents of this README have been moved to the "Deploying Orchestrator"
section of the docs: https://facebookincubator.github.io/magma.

## Variables

| Name | Description | Type | Default | Required |
|------|-------------|:----:|:-----:|:-----:|
| db_password | The password for the RDS instance | string | "" | **yes** |
| nms_db_password | The password for the nms RDS instance | string | "" | **yes** |
| key_name | The name of the EC2 keypair for SSH access to nodes | string | "" | **yes** |
| region | The AWS region to provision the resources in | string | "eu-west-1" | no |
| vpc_name | The name of the provisioned VPC | string | "orc8r-vpc" | no |
| cluster_name | The name of the provisioned EKS cluster | string | "orc8r" | no |
| map_users | Additional IAM users to add to the aws-auth configmap | list(map(string)) | [] | no 
