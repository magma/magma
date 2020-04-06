# orc8r-helm-aws

This is a Terraform module which installs the Orchestrator application and all
supporting components into an EKS cluster.

## Providers

| Name | Version |
|------|---------|
| aws | >= 2.6.0 |
| helm | ~> 0.10 |
| kubernetes | ~> 1.10.0 |
| null | n/a |
| terraform | n/a |

## Inputs

| Name | Description | Type | Default | Required |
|------|-------------|------|---------|:-----:|
| deploy\_nms | Flag to deploy NMS | `bool` | n/a | yes |
| deploy\_openvpn | Flag to deploy openvpn server into cluster. This is useful if you want to remotely access AGW's. | `bool` | `false` | no |
| docker\_pass | Docker registry password | `string` | n/a | yes |
| docker\_registry | Docker registry to pull orc8r containers from | `string` | n/a | yes |
| docker\_user | Docker username to login to registry with | `string` | n/a | yes |
| efs\_file\_system\_id | ID of the EFS file system to use for k8s persistent volumes. | `string` | n/a | yes |
| efs\_provisioner\_role\_arn | ARN of the IAM role for the EFS provisioner. | `string` | n/a | yes |
| eks\_cluster\_id | EKS cluster ID for the kubernetes cluster | `string` | n/a | yes |
| elasticsearch\_endpoint | Endpoint of the Elasticsearch datasink for aggregated logs and events. | `string` | n/a | yes |
| elasticsearch\_retention\_days | Retention period in days of ES indices. | `number` | `7` | no |
| external\_dns\_role\_arn | IAM role ARN for external-dns | `string` | n/a | yes |
| helm\_pass | Helm repository password | `string` | n/a | yes |
| helm\_repo | Helm repository URL for orc8r charts | `string` | n/a | yes |
| helm\_user | Helm username to login to repositoriy with | `string` | n/a | yes |
| install\_tiller | Install tiller in the cluster or not | `bool` | `true` | no |
| nms\_db\_host | DB hostname for NMS database connection | `string` | n/a | yes |
| nms\_db\_name | DB name for NMS database connection | `string` | n/a | yes |
| nms\_db\_pass | NMS DB password | `string` | n/a | yes |
| nms\_db\_user | DB username for NMS database connection | `string` | n/a | yes |
| orc8r\_chart\_version | Version of the Orhcestrator Helm chart to install | `string` | n/a | yes |
| orc8r\_controller\_replicas | Replica count for Orchestrator controller pods. | `number` | `2` | no |
| orc8r\_db\_host | DB hostname for Orchestrator database connection | `string` | n/a | yes |
| orc8r\_db\_name | DB name for Orchestrator database connection | `string` | n/a | yes |
| orc8r\_db\_pass | Orchestrator DB password | `string` | n/a | yes |
| orc8r\_db\_port | DB port for Orchestrator database connection | `number` | `5432` | no |
| orc8r\_db\_user | DB username for Orchestrator database connection | `string` | n/a | yes |
| orc8r\_domain\_name | Base domain name for Orchestrator | `string` | n/a | yes |
| orc8r\_kubernetes\_namespace | K8s namespace to install main Orchestrator components into. | `string` | `"orc8r"` | no |
| orc8r\_proxy\_replicas | Replica count for Orchestrator proxy pods. | `number` | `2` | no |
| orc8r\_route53\_zone\_id | Route53 zone ID of Orchestrator domain name for external-DNS | `string` | n/a | yes |
| orc8r\_tag | Image tag for Orchestrator components. | `string` | `""` | no |
| region | AWS region to deploy Orchestrator components into. The chosen region must provide EKS. | `string` | n/a | yes |
| secretsmanager\_orc8r\_name | Name of the AWS secretsmanager secret where Orchestrator deployment secrets will be stored. | `string` | n/a | yes |
| seed\_certs\_dir | Directory on LOCAL disk where orc8r certificates are stored to seed Secretsmanager values. Home directory and env vars will be expanded. | `string` | n/a | yes |

## Outputs

No output.
