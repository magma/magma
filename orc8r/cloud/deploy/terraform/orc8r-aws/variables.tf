################################################################################
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
################################################################################

variable "region" {
  description = "AWS region to deploy Orchestrator components into. The chosen region must provide EKS."
  type        = string
}

data "aws_availability_zones" "available" {}

##############################################################################
# Module flags
##############################################################################

variable "deploy_elasticsearch" {
  description = "Flag to deploy AWS Elasticsearch service as the datasink for aggregated logs."
  type        = bool
  default     = false
}

variable "deploy_elasticsearch_service_linked_role" {
  description = "Flag to deploy AWS Elasticsearch service linked role with cluster. If you've already created an ES service linked role for another cluster, you should set this to false."
  type        = bool
  default     = true
}


variable "enable_aws_db_notifications" {
  description = "Flag to enable AWS RDS notifications"
  type        = bool
  default     = false
}

variable "magma_uuid" {
  description = "UUID to identify Orc8r deployment"
  type        = string
  default     = "default"
}

variable "global_tags" {
  default = {}
}
##############################################################################
# K8s configuration
##############################################################################

variable "orc8r_domain_name" {
  description = "Base domain name for AWS Route 53 hosted zone."
  type        = string
}

##############################################################################
# K8s configuration
##############################################################################

variable "cluster_name" {
  description = "Name for the Orchestrator EKS cluster."
  type        = string
  default     = "orc8r"
}

variable "cluster_version" {
  description = "Kubernetes version for the EKS cluster."
  type        = string
  default     = "1.17"
}

variable "eks_worker_group_key" {
  description = "If specified, the worker nodes for EKS will use this EC2 keypair."
  type        = string
  default     = null
}

variable "eks_worker_additional_sg_ids" {
  description = "Additional security group IDs to attach to EKS worker nodes."
  type        = list(string)
  default     = []
}

variable "eks_worker_additional_policy_arns" {
  description = "Additional IAM policy ARNs to attach to EKS worker nodes."
  type        = list(string)
  default     = []
}

variable "eks_worker_groups" {
  # Check the docs at https://github.com/terraform-aws-modules/terraform-aws-eks
  # for the complete set of valid properties for these objects.
  description = "Worker group configuration for EKS. Default value is 1 worker group consisting of 3 t3.large instances."
  type        = any
  default = [
    {
      name                 = "wg-1"
      instance_type        = "t3.large"
      asg_desired_capacity = 3
      asg_min_size         = 1
      asg_max_size         = 3
      autoscaling_enabled  = false
      kubelet_extra_args = "" // object types must be identical (see thanos_worker_groups)
    },
  ]
}

variable "thanos_worker_groups" {
  # Check the docs at https://github.com/terraform-aws-modules/terraform-aws-eks
  # for the complete set of valid properties for these objects. This worker group
  # exists because some thanos components (compact) require significant instance
  # storage to operate.
  # Use label key 'compute-type' to specify the node used by nodeSelector
  # in the helm release
  description = "Worker group configuration for Thanos. Default consists of 1 group consisting of 1 m5d.xlarge for thanos."
  type        = any
  default = [
    {
      name                 = "thanos-1"
      instance_type        = "m5d.xlarge"
      asg_desired_capacity = 1
      asg_min_size         = 1
      asg_max_size         = 1
      autoscaling_enabled  = false
      kubelet_extra_args = "--node-labels=compute-type=thanos"
    },
  ]

}

variable "eks_map_roles" {
  description = "EKS IAM role mapping. Note that by default, the creator of the cluster will be in the system:master group."
  type = list(
    object({
      rolearn  = string,
      username = string,
      groups   = list(string),
    })
  )
  default = []
}

variable "eks_map_users" {
  description = "Additional IAM users to add to the aws-auth ConfigMap."
  type = list(object({
    userarn  = string
    username = string
    groups   = list(string)
  }))
  default = []
}

##############################################################################
# EFS configuration
##############################################################################

variable "efs_project_name" {
  description = "Project name for EFS file system"
  type        = string
  default     = "orc8r"
}

##############################################################################
# VPC configuration
##############################################################################

# TODO: support an existing VPC

variable "vpc_name" {
  description = "Name for the VPC that will contain all the Orchestrator components."
  type        = string
  default     = "orc8r_vpc"
}

variable "vpc_cidr" {
  description = "CIDR block for the VPC."
  type        = string
  default     = "10.10.0.0/16"
}

variable "vpc_public_subnets" {
  description = "CIDR blocks for the VPC's public subnets."
  type        = list(string)
  default     = ["10.10.0.0/24", "10.10.1.0/24", "10.10.2.0/24"]
}

variable "vpc_private_subnets" {
  description = "CIDR blocks for the VPC's private subnets."
  type        = list(string)
  default     = ["10.10.100.0/24", "10.10.101.0/24", "10.10.102.0/24"]
}

variable "vpc_database_subnets" {
  description = "CIDR blocks for the VPC's database subnets."
  type        = list(string)
  default     = ["10.10.200.0/24", "10.10.201.0/24", "10.10.202.0/24"]
}

variable "vpc_extra_tags" {
  description = "Tags to add to the VPC."
  default     = {}
}

##############################################################################
# Orchestrator DB Specs
##############################################################################

variable "orc8r_db_identifier" {
  description = "Identifier for the RDS instance for Orchestrator."
  type        = string
  default     = "orc8rdb"
}

variable "orc8r_db_storage_gb" {
  description = "Capacity in GB to allocate for Orchestrator RDS instance."
  type        = number
  default     = 64
}

variable "orc8r_db_instance_class" {
  description = "RDS instance type for Orchestrator DB."
  type        = string
  default     = "db.m4.large"
}

variable "orc8r_db_name" {
  description = "DB name for Orchestrator RDS instance."
  type        = string
  default     = "orc8r"
}

variable "orc8r_db_username" {
  description = "Username for default DB user for Orchestrator DB."
  type        = string
  default     = "orc8r"
}

variable "orc8r_db_password" {
  description = "Password for the Orchestrator DB. Must be at least 8 characters."
  type        = string
}

variable "orc8r_db_engine_version" {
  description = "Postgres engine version for Orchestrator DB."
  type        = string
  default     = "9.6.15"
}

variable "orc8r_db_dialect" {
  description = "Database dialect for Orchestrator DB."
  type        = string
  default     = "postgres"
}

variable "orc8r_db_backup_retention" {
  description = "Database backup retention period"
  type        = number
  default     = 7
}

variable "orc8r_db_backup_window" {
  description = "Database daily backup window in UTC with a 30-minute minimum"
  type        = string
  default     = "01:00-01:30"
}

variable "orc8r_db_event_subscription" {
  description = "Database event subscription"
  type        = string
  default     = "orc8r-rds-events"
}


##############################################################################
# Secretmanager configuration
##############################################################################

variable "secretsmanager_orc8r_secret" {
  description = "AWS Secret Manager secret to store Orchestrator secrets."
  type        = string
}

##############################################################################
# Elasticsearch configuration
##############################################################################

variable "elasticsearch_domain_name" {
  description = "Name for the ES domain."
  type        = string
  default     = null
}

variable "elasticsearch_version" {
  description = "ES version for ES domain."
  default     = "7.1"
}

variable "elasticsearch_instance_type" {
  description = "AWS instance type for ES domain."
  type        = string
  default     = null
}

variable "elasticsearch_instance_count" {
  description = "Number of instances to allocate for ES domain."
  type        = number
  default     = null
}

variable "elasticsearch_dedicated_master_enabled" {
  description = "Enable/disable dedicated master nodes for ES."
  type        = bool
  default     = false
}

variable "elasticsearch_dedicated_master_type" {
  description = "Instance type for ES dedicated master nodes."
  type        = string
  default     = null
}

variable "elasticsearch_dedicated_master_count" {
  description = "Number of dedicated ES master nodes."
  type        = number
  default     = null
}

variable "elasticsearch_az_count" {
  description = "AZ count for ES."
  type        = number
  default     = 2
}

variable "elasticsearch_ebs_enabled" {
  description = "Use EBS for ES storage. See https://aws.amazon.com/elasticsearch-service/pricing/ to check if your chosen instance types can use non-EBS storage."
  type        = bool
  default     = false
}

variable "elasticsearch_ebs_volume_size" {
  description = "Size in GB to allocate for ES EBS data volumes."
  type        = number
  default     = null
}

variable "elasticsearch_ebs_volume_type" {
  description = "EBS volume type for ES data volumes."
  type        = string
  default     = null
}

variable "elasticsearch_ebs_iops" {
  description = "IOPS for ES EBS volumes."
  type        = number
  default     = null
}

variable "elasticsearch_domain_tags" {
  description = "Extra tags for the ES domain."
  default     = {}
}

variable "thanos_enabled" {
  description = "Enable thanos infrastructure"
  type        = bool
  default     = false
}

##############################################################################
# Simple Notification Service (SNS) configuration
##############################################################################

variable "orc8r_sns_name" {
  description = "SNS for Orc8r to redirect alerts and notifications"
  type        = string
  default     = "orc8r-sns"
  }

variable "orc8r_sns_email" {
  description = "SNS email endpoint to send notifications"
  type        = string
  default     = ""
}
