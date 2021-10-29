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

variable "state_backend" {
  description = "State backend for terraform (e.g. s3, local)."
  type        = string
  default     = "local"
}

variable "state_config" {
  description = "Optional config for state backend. The object type will depend on your backend."
  default     = null
}

data "aws_availability_zones" "available" {}

##############################################################################
# DNS configuration
##############################################################################

variable "orc8r_domain_name" {
  description = "Base domain name for Orchestrator."
  type        = string
}

variable "orc8r_route53_zone_id" {
  description = "Route53 zone ID of Orchestrator domain name for ExternalDNS."
  type        = string
}

variable "external_dns_role_arn" {
  description = "IAM role ARN for ExternalDNS."
  type        = string
}

##############################################################################
# Kubernetes configuration
##############################################################################

variable "eks_cluster_id" {
  description = "EKS cluster ID for the K8s cluster."
  type        = string
}

variable "orc8r_kubernetes_namespace" {
  description = "K8s namespace to install main Orchestrator components into."
  type        = string
  default     = "orc8r"
}

variable "monitoring_kubernetes_namespace" {
  description = "K8s namespace to install Orchestrator monitoring components into."
  type        = string
  default     = "monitoring"
}

##############################################################################
# General Orchestrator configuration
##############################################################################

variable "deploy_nms" {
  description = "Flag to deploy NMS."
  type        = bool
  default     = true
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
  description = "DB name for Orchestrator database connection."
  type        = string
}

variable "orc8r_db_host" {
  description = "DB hostname for Orchestrator database connection."
  type        = string
}

variable "orc8r_db_port" {
  description = "DB port for Orchestrator database connection."
  type        = number
  default     = 5432
}

variable "orc8r_db_user" {
  description = "DB username for Orchestrator database connection."
  type        = string
}

variable "nms_db_name" {
  description = "DB name for NMS database connection."
  type        = string
}

variable "nms_db_host" {
  description = "DB hostname for NMS database connection."
  type        = string
}

variable "nms_db_user" {
  description = "DB username for NMS database connection."
  type        = string
}

##############################################################################
# Helm configuration
##############################################################################

variable "existing_tiller_service_account_name" {
  description = "Name of existing Tiller service account to use for Helm."
  type        = string
  default     = null
}

variable "tiller_namespace" {
  description = "Namespace where Tiller is installed or should be installed into."
  type        = string
  default     = "kube-system"
}

variable "install_tiller" {
  description = "Install Tiller in the cluster or not."
  type        = bool
  default     = true
}

variable "helm_deployment_name" {
  description = "Name for the Helm release."
  type        = string
  default     = "orc8r"
}

variable "orc8r_deployment_type" {
  description = "Type of orc8r deployment (fixed wireless access, federated fixed wireless access, or all modules)"
  type        = string
  validation {
    condition = (
      var.orc8r_deployment_type == "fwa" ||
      var.orc8r_deployment_type == "federated_fwa" ||
      var.orc8r_deployment_type == "all"
    )
    error_message = "The orc8r_deployment_type value must be one of ['fwa', 'federated_fwa', 'all']."
  }
}

variable "orc8r_chart_version" {
  description = "Version of the core orchestrator Helm chart to install."
  type        = string
  default     = "1.5.8"
}

variable "cwf_orc8r_chart_version" {
  description = "Version of the orchestrator cwf module Helm chart to install."
  type        = string
  default     = "0.2.0"
}

variable "fbinternal_orc8r_chart_version" {
  description = "Version of the orchestrator fbinternal module Helm chart to install."
  type        = string
  default     = "0.2.0"
}

variable "feg_orc8r_chart_version" {
  description = "Version of the orchestrator feg module Helm chart to install."
  type        = string
  default     = "0.2.1"
}

variable "lte_orc8r_chart_version" {
  description = "Version of the orchestrator lte module Helm chart to install."
  type        = string
  default     = "0.2.1"
}

variable "wifi_orc8r_chart_version" {
  description = "Version of the orchestrator wifi module Helm chart to install."
  type        = string
  default     = "0.2.0"
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
  description = "ID of the EFS file system to use for K8s persistent volumes."
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
  description = "Retention period in days of Elasticsearch indices."
  type        = number
  default     = 7
}

##############################################################################
# Secret configuration and values
##############################################################################

variable "secretsmanager_orc8r_name" {
  description = "Name of the AWS Secrets Manager secret where Orchestrator deployment secrets will be stored."
  type        = string
}

variable "orc8r_db_pass" {
  description = "Orchestrator DB password."
  type        = string
}

variable "nms_db_pass" {
  description = "NMS DB password."
  type        = string
}

variable "docker_registry" {
  description = "Docker registry to pull Orchestrator containers from."
  type        = string
}

variable "docker_user" {
  description = "Docker username to login to registry with."
  type        = string
}

variable "docker_pass" {
  description = "Docker registry password."
  type        = string
}

variable "seed_certs_dir" {
  description = "Directory on LOCAL disk where Orchestrator certificates are stored to seed Secrets Manager values. Home directory and env vars will be expanded."
  type        = string
}

variable "helm_repo" {
  description = "Helm repository URL for Orchestrator charts."
  type        = string
}

variable "helm_user" {
  description = "Helm username to login to repository with."
  type        = string
}

variable "helm_pass" {
  description = "Helm repository password."
  type        = string
}

##############################################################################
# Other deployment flags
##############################################################################

variable "deploy_openvpn" {
  description = "Flag to deploy OpenVPN server into cluster. This is useful if you want to remotely access AGWs."
  type        = bool
  default     = false
}

##############################################################################
# Thanos Object Storage
##############################################################################

variable "thanos_enabled" {
  description = "Deploy thanos components and object storage"
  type        = bool
  default     = false
}

variable "thanos_object_store_bucket_name" {
  description = "Bucket name for s3 object storage. Must be globally unique"
  type        = string
  default     = ""
}

variable "thanos_query_node_selector" {
  description = "NodeSelector value to specify which node to run thanos query pod on. Default is 'thanos' to be deployed on the default thanos worker group."
  type        = string
  default     = "thanos"
}


variable "thanos_compact_node_selector" {
  description = "NodeSelector value to specify which node to run thanos compact pod on. Label is 'compute-type:<value>'"
  type        = string
  default     = ""
}

variable "thanos_store_node_selector" {
  description = "NodeSelector value to specify which node to run thanos store pod on. Label is 'compute-type:<value>'"
  type        = string
  default     = ""
}

##############################################################################
# Analytics Service
##############################################################################
variable "analytics_export_enabled" {
  description = "Deploy thanos components and object storage"
  type        = bool
  default     = false
}

variable "analytics_metrics_prefix" {
  description = "Bucket name for s3 object storage. Must be globally unique"
  type        = string
  default     = ""
}

variable "analytics_app_secret" {
  description = "App secret for which the metrics is to be exported to"
  type = string
  default = ""
}


variable "analytics_app_id" {
  description = "App ID for which the metrics is to be exported to"
  type = string
  default = ""
}

variable "analytics_metric_export_url" {
  description = "Metric Export URL"
  type = string
  default = ""
}

variable "analytics_category_name" {
  description = "Category under which the exported metrics will be placed under"
  type = string
  default = "magma"
}
