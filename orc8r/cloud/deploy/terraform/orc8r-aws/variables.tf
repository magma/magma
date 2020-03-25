################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
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

variable "global_tags" {
  default = {}
}
##############################################################################
# K8s configuration
##############################################################################

variable "orc8r_domain_name" {
  description = "Base domain name for AWS Route 53 hosted zone"
  type        = string
}

##############################################################################
# K8s configuration
##############################################################################

variable "cluster_name" {
  description = "Name for the orc8r EKS cluster."
  type        = string
  default     = "orc8r"
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

variable "eks_worker_groups" {
  description = "Worker group configuration for EKS. Default value is 1 worker group consisting of 3 t3.small instances."
  type = list(
    object({
      name                 = string,
      instance_type        = string,
      asg_desired_capacity = number,
      asg_min_size         = number,
      asg_max_size         = number,
      autoscaling_enabled  = bool,
    })
  )
  default = [
    {
      name                 = "wg-1"
      instance_type        = "t3.large"
      asg_desired_capacity = 3
      asg_min_size         = 1
      asg_max_size         = 3
      autoscaling_enabled  = false
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
  description = "Additional IAM users to add to the aws-auth configmap."
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
  description = "Password for the Orchestrator DB."
  type        = string
}

variable "orc8r_db_engine_version" {
  description = "Postgres engine version for Orchestrator DB."
  type        = string
  default     = "9.6.15"
}

##############################################################################
# NMS DB Specs
##############################################################################

variable "nms_db_identifier" {
  description = "Identifier for the RDS instance for NMS."
  type        = string
  default     = "nmsdb"
}

variable "nms_db_storage_gb" {
  description = "Capacity in GB to allocate for NMS RDS instance."
  type        = number
  default     = 16
}

variable "nms_db_instance_class" {
  description = "RDS instance type for NMS DB."
  type        = string
  default     = "db.m4.large"
}

variable "nms_db_name" {
  description = "DB name for NMS RDS instance."
  type        = string
  default     = "magma"
}

variable "nms_db_username" {
  description = "Username for default DB user for NMS DB."
  type        = string
  default     = "magma"
}

variable "nms_db_password" {
  description = "Password for the NMS DB."
  type        = string
}

variable "nms_db_engine_version" {
  description = "MySQL engine version for NMS DB."
  type        = string
  default     = "5.7"
}

##############################################################################
# Secretmanager configuration
##############################################################################

variable "secretsmanager_orc8r_secret" {
  description = "AWS Secretmanager secret to store Orchestrator secrets."
  type        = string
}

##############################################################################
# Elasticsearch configuration
##############################################################################

variable "elasticsearch_domain_name" {
  description = "Name for the AWS Elasticsearch domain."
  type        = string
  default     = null
}

variable "elasticsearch_version" {
  description = "ES version for Elasticsearch domain."
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
