Online v1.1.0 Orchestrator Upgrade
===

IMPORTANT: Read the "Upgrading from 1.0" Orchestrator documentation for details
on how to use this module to upgrade your deployment. If you just
`terraform apply` here things will probably go very poorly.

The files in this directory can be copied and used as-is to perform an online
upgrade of the Orchestrator application from v1.0.x to v1.1.x. Please read the
descriptions of all the variables very carefully, and double-check the output
of `terraform apply` before applying changes. A misconfiguration may result in
significant downtime.

A lot of variables in this module have been set with defaults equal to the
legacy Terraform root module's defaults. If you copied that as-is to deploy
your v1.0.x application, you can leave these default values alone. If something
does not match up with your old Terraform configuration, you will see an
unexpected planned step in the output of `terraform apply` before it asks
for confirmation.

Detailed upgrade steps can be found in the Orchestrator documentation at
https://magma.github.io/magma.

## Providers

| Name | Version |
|------|---------|
| aws | n/a |
| template | n/a |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:-----:|
| additional\_eks\_worker\_groups | Additional EKS worker nodes to spin up while the v1.1.0 application is deployed concurrently with the v1.0 application. | `any` | <pre>[<br>  {<br>    "asg_desired_capacity": 3,<br>    "asg_min_size": 3,<br>    "instance_type": "t3.large",<br>    "name": "wg-1",<br>    "tags": [<br>      {<br>        "key": "orc8r-node-type",<br>        "propagate_at_launch": true,<br>        "value": "orc8r-worker-node"<br>      }<br>    ]<br>  }<br>]</pre> | no |
| deploy\_elasticsearch | Deploy elasticsearch cluster for log aggregation (default false). | `bool` | `false` | no |
| deploy\_elasticsearch\_linked\_role | Deploy ES linked role if ES is deployed. | `bool` | `true` | no |
| deploy\_nms | Whether to deploy NMS. You can leave this set to true for the online upgrade, unlike the from-scratch v1.1.0 installation. | `bool` | `true` | no |
| docker\_pass | Password for your Docker user | `string` | n/a | yes |
| docker\_registry | URL to your Docker registry | `string` | n/a | yes |
| docker\_user | Username for your Docker registry | `string` | n/a | yes |
| eks\_cluster\_name | Name of the EKS cluster that the v1.0 application is deployed on. This should match your v1.0 Terraform. | `string` | n/a | yes |
| eks\_map\_users | Additional users you want to grant access to EKS to. This should match your v1.0 Terraform or those users will lose k8s access. | `any` | `[]` | no |
| elasticsearch\_domain\_configuration | Configuration for the ES domain | <pre>object({<br>    version         = string<br>    instance_type   = string<br>    instance_count  = number<br>    az_count        = number<br>    ebs_enabled     = bool<br>    ebs_volume_size = number<br>    ebs_volume_type = string<br>  })</pre> | <pre>{<br>  "az_count": 3,<br>  "ebs_enabled": true,<br>  "ebs_volume_size": 32,<br>  "ebs_volume_type": "gp2",<br>  "instance_count": 3,<br>  "instance_type": "t2.medium.elasticsearch",<br>  "version": "7.4"<br>}</pre> | no |
| elasticsearch\_domain\_name | Name for ES domain | `string` | `"orc8r-es-domain"` | no |
| helm\_pass | Password for your Helm user | `string` | n/a | yes |
| helm\_repo | URL to your Helm repo. Don't forget the protocol prefix (e.g. https://) | `string` | n/a | yes |
| helm\_user | Username for your Helm repo | `string` | n/a | yes |
| metrics\_worker\_subnet\_id | Subnet ID of the metrics worker instance. Find this in the EC2 console (the instance will have the tag orc8r-node-type: orc8r-prometheus-node). | `string` | n/a | yes |
| new\_deployment\_name | New name for the v1.1.0 Helm deployment. This must be different than your old v1.0 deployment (which was probably 'orc8r') | `string` | n/a | yes |
| orc8r\_chart\_version | Chart version for the Helm deployment | `string` | `"1.4.21"` | no |
| orc8r\_container\_tag | Container tag to deploy | `string` | n/a | yes |
| orc8r\_controller\_replicas | How many controller pod replicas to deploy | `number` | `2` | no |
| orc8r\_db\_configuration | Configuration of the Orchestrator Postgres instance. This should match the v1.0 Terraform. | <pre>object({<br>    identifier     = string<br>    storage_gb     = number<br>    engine_version = string<br>    instance_class = string<br>  })</pre> | <pre>{<br>  "engine_version": "9.6.11",<br>  "identifier": "orc8rdb",<br>  "instance_class": "db.m4.large",<br>  "storage_gb": 32<br>}</pre> | no |
| orc8r\_db\_password | Password for the Orchestrator Postgres instance. This should match the v1.0 Terraform. | `string` | n/a | yes |
| orc8r\_domain | Root domain or subdomain for your Orchestrator deployment (e.g. orc8r.mydomain.com). | `string` | n/a | yes |
| orc8r\_proxy\_replicas | How many proxy pod replicas to deploy | `number` | `2` | no |
| prometheus\_ebs\_az | Availability zone that the Prometheus worker node and EBS volume are located in. Find this in the EC2 console. | `string` | n/a | yes |
| prometheus\_ebs\_size | Size of the EBS volume for the Prometheus data EBS volume. This should match your v1.0 Terraform. | `number` | `64` | no |
| region | AWS region to deploy to. This should match your v1.0 Terraform. | `string` | n/a | yes |
| secretsmanager\_secret\_name | Name for the Secretsmanager secret that the orc8r-aws module will create. | `string` | n/a | yes |
| seed\_certs\_dir | Directory with your Orchestrator certificates. | `string` | n/a | yes |
| ssh\_key\_name | Name of the SSH key you created for the v1.0 infra. | `string` | n/a | yes |
| vpc\_configuration | Configuration of the VPC that the v1.0 chart is deployed in. This should match your v1.0 Terraform. | <pre>object({<br>    cidr            = string<br>    public_subnets  = list(string)<br>    private_subnets = list(string)<br>    db_subnets      = list(string)<br>  })</pre> | <pre>{<br>  "cidr": "10.10.0.0/16",<br>  "db_subnets": [<br>    "10.10.11.0/24",<br>    "10.10.12.0/24",<br>    "10.10.13.0/24"<br>  ],<br>  "private_subnets": [],<br>  "public_subnets": [<br>    "10.10.1.0/24",<br>    "10.10.2.0/24",<br>    "10.10.3.0/24"<br>  ]<br>}</pre> | no |
| vpc\_name | Name of the VPC that the v1.0 infra is deployed in. This should match your v1.0 Terraform. | `string` | n/a | yes |
| worker\_node\_policy\_suffix | The name suffix of the custom IAM node policy from the v1.0 Terraform root module. This policy name will begin with magma\_eks\_worker\_node\_policy. | `string` | n/a | yes |

## Outputs

| Name | Description |
|------|-------------|
| nameservers | n/a |
| vals | n/a |
