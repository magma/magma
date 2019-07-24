# Example Terraform for Orchestrator EKS

This directory contains a bare-bones Terraform script to set up an EKS cluster
which the cloud component of Orchestrator can be deployed to. Copy this
directory to a location of your choice and `terraform apply` to provision the
resources.

You should check in the generated .tfstate files after a `terraform apply`, so
consider copying this directory into your project's source control repository.
The terraform script will also output a kubeconfig file for `kubectl` and an
aws auth config map to manage access to the EKS cluster. These files should
also be checked in.

## Prerequisites

This script depends on Terraform >= 0.12.0, and you'll need a few other
dependencies to manage and deploy orc8r to the cluster. If you're on a Mac:

```bash
brew install terraform aws-iam-authenticator kubernetes-cli kubernetes-helm awscli
```

Next `cd` to where you've copied the contents of this directory and

```bash
terraform init
```

In the AWS console, create or import a new keypair to enable SSH access to the
worker nodes of the EKS cluster. This can be found in the EC2 dashboard under
"Key Pairs". If you're creating a new key pair, make sure not to lose the
private key :)

![](keypair.png)

Next, create a `vars.tfvars` file in this directory, *add it to your source
control's .ignore*, and specify your desired RDS password and the name of the
keypair that you imported or created in the above step:

```
$ cat vars.tfvars
db_password = "foobar"
key_name = "my_key"
```

Now you're ready to move on:

## Terraform Workflow

You should always run a `terraform plan` before `terraform apply`. This will
let you preview what Terraform is going to do. Always sanity-check the output
of `terraform plan` before you run `terraform apply`.

1. `terraform plan`: check the output to make sure you're not deleting any
resources which shouldn't be deleted
2. `terraform apply`
3. Check the generated .tfstate files into source control

In a team setting, refrain from running `terraform apply` before tfstate files
have been checked in after another apply. You may end up with nasty conflicts
in the state files. If that happens, just `terraform refresh` and commit those
changes.

## Setup Steps After Terraforming

First find the public IP address of the metrics instance using
```bash
export METRICS_IP=$(aws ec2 describe-instances --filters Name=tag:orc8r-node-type,Values=orc8r-prometheus-node --query 'Reservations[*].Instances[0].PublicIpAddress' --output text)
echo $METRICS_IP
```
The Prometheus config manager application expects some configuration files to
be seeded in the EBS config volume (don't forget to use the correct private
key in `scp` with the `-i` flag:

```bash
scp -r config_defaults ec2-user@$METRICS_IP:~
ssh ec2-user@$METRICS_IP
[ec2-user@<metrics-ip> ~]$ sudo cp -r config_defaults/. /configs/prometheus
```

Now you've got your infra set up, we can move on to configuring the EKS cluster.

Assuming you don't have an existing Kubeconfig file in `~/.kube/config`:

```bash
cp ./kubeconfig_orc8r ~/.kube/config
```

Now we can set up access to the EKS cluster:

```bash
kubectl apply -f config-map-aws-auth_orc8r.yaml
kubectl create namespace magma
```

Label the EKS worker nodes appropriately, so we can schedule the metrics pod on
the metrics worker and the Orchestrator pods on the Orchestrator worker nodes:

```bash
export AWS_DEFAULT_REGION=...
aws ec2 describe-instances --filters Name=tag:orc8r-node-type,Values=orc8r-worker-node \
  --query 'Reservations[].Instances[].[PrivateDnsName]' --output text \
  | xargs -I % kubectl -n magma label nodes % worker-type=controller --overwrite

aws ec2 describe-instances --filters Name=tag:orc8r-node-type,Values=orc8r-prometheus-node \
  --query 'Reservations[].Instances[].[PrivateDnsName]' --output text \
  | xargs -I % kubectl -n magma label nodes % worker-type=metrics --overwrite
```

At this point, if everything succeeded, you're ready to move on to the initial
Helm deployment. Follow the README in `magma/orc8r/cloud/helm`.

## Variables

| Name | Description | Type | Default | Required |
|------|-------------|:----:|:-----:|:-----:|
| db_password | The password for the RDS instance | string | "" | **yes** |
| nms_db_password | The password for the nms RDS instance | string | "" | **yes** |
| key_name | The name of the EC2 keypair for SSH access to nodes | string | "" | **yes** |
| region | The AWS region to provision the resources in | string | "eu-west-1" | no |
| vpc_name | The name of the provisioned VPC | string | "orc8r-vpc" | no |
| cluster_name | The name of the provisioned EKS cluster | string | "orc8r" | no |
| map_users | Additional IAM users to add to the aws-auth configmap | list(map(string)) | [] | no |
| map_users_count | How many users are in map_users | string | 0 | no
