Fully Remote Secrets
===

This is an example of a root Terraform module for an Orchestrator 1.1.0
cluster that stores all sensitive secrets in the same AWS account that the
cluster is deployed in.

This uses S3 to store Terraform state remotely, locked by a DynamoDB table. The
S3 bucket should be created manually before the first Terraform.

All secret values passed into the Orchestrator modules are read from an AWS
secretsmanager secret that should be created and populated manually. With this
setup, it is completely safe to check your root module in to source control,
even in a public repository.

This example demonstrates how to use Secretsmanager to store the sensitive
variables required by the Orchestrator mdoules, but you can use any service
that has a Terraform provider, such as Hashicorp Vault or Lastpass.
