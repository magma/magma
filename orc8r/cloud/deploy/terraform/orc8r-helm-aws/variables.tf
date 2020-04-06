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

variable "state_backend" {
  description = "State backend for terraform (e.g. s3, local)"
  type        = string
  default     = "local"
}

data "aws_availability_zones" "available" {}

##############################################################################
# DNS configuration
##############################################################################

variable "orc8r_domain_name" {
  description = "Base domain name for Orchestrator"
  type        = string
}

variable "orc8r_route53_zone_id" {
  description = "Route53 zone ID of Orchestrator domain name for external-DNS"
  type        = string
}

variable "external_dns_role_arn" {
  description = "IAM role ARN for external-dns"
  type        = string
}

##############################################################################
# Kubernetes configuration
##############################################################################

variable "eks_cluster_id" {
  description = "EKS cluster ID for the kubernetes cluster"
  type        = string
}

variable "orc8r_kubernetes_namespace" {
  description = "K8s namespace to install main Orchestrator components into."
  type        = string
  default     = "orc8r"
}

##############################################################################
# General Orchestrator configuration
##############################################################################

variable "deploy_nms" {
  description = "Flag to deploy NMS"
  type        = bool
}

variable "orc8r_controller_replicas" {
  description = "Replica count for Orchestrator controller pods."
  type        = number
  default     = 2
}

variable "orc8r_proxy_replicas" {
  description = "Replica count for Orchestrator proxy pods."
  type        = number
  default     = 2
}

variable "orc8r_db_name" {
  description = "DB name for Orchestrator database connection"
  type        = string
}

variable "orc8r_db_host" {
  description = "DB hostname for Orchestrator database connection"
  type        = string
}

variable "orc8r_db_port" {
  description = "DB port for Orchestrator database connection"
  type        = number
  default     = 5432
}

variable "orc8r_db_user" {
  description = "DB username for Orchestrator database connection"
  type        = string
}

variable "nms_db_name" {
  description = "DB name for NMS database connection"
  type        = string
}

variable "nms_db_host" {
  description = "DB hostname for NMS database connection"
  type        = string
}

variable "nms_db_user" {
  description = "DB username for NMS database connection"
  type        = string
}

##############################################################################
# Helm configuration
##############################################################################

variable "install_tiller" {
  description = "Install tiller in the cluster or not"
  type        = bool
  default     = true
}

variable "orc8r_chart_version" {
  description = "Version of the Orhcestrator Helm chart to install"
  type        = string
}

variable "orc8r_tag" {
  description = "Image tag for Orchestrator components."
  type        = string
  default     = ""
}

##############################################################################
# EFS configuration
##############################################################################

variable "efs_file_system_id" {
  description = "ID of the EFS file system to use for k8s persistent volumes."
  type        = string
}

variable "efs_provisioner_role_arn" {
  description = "ARN of the IAM role for the EFS provisioner."
  type        = string
}

##############################################################################
# Log aggregation configuration
##############################################################################

variable "elasticsearch_endpoint" {
  description = "Endpoint of the Elasticsearch datasink for aggregated logs and events."
  type        = string
  default     = null
}

variable "elasticsearch_retention_days" {
  description = "Retention period in days of ES indices."
  type        = number
  default     = 7
}

##############################################################################
# Secret configuration and values
##############################################################################

variable "secretsmanager_orc8r_name" {
  description = "Name of the AWS secretsmanager secret where Orchestrator deployment secrets will be stored."
  type        = string
}

variable "orc8r_db_pass" {
  description = "Orchestrator DB password"
  type        = string
}

variable "nms_db_pass" {
  description = "NMS DB password"
  type        = string
}

variable "docker_registry" {
  description = "Docker registry to pull orc8r containers from"
  type        = string
}

variable "docker_user" {
  description = "Docker username to login to registry with"
  type        = string
}

variable "docker_pass" {
  description = "Docker registry password"
  type        = string
}

variable "seed_certs_dir" {
  description = "Directory on LOCAL disk where orc8r certificates are stored to seed Secretsmanager values. Home directory and env vars will be expanded."
  type        = string
}

variable "helm_repo" {
  description = "Helm repository URL for orc8r charts"
  type        = string
}

variable "helm_user" {
  description = "Helm username to login to repositoriy with"
  type        = string
}

variable "helm_pass" {
  description = "Helm repository password"
  type        = string
}

##############################################################################
# Other deployment flags
##############################################################################

variable "deploy_openvpn" {
  description = "Flag to deploy openvpn server into cluster. This is useful if you want to remotely access AGW's."
  type        = bool
  default     = false
}
