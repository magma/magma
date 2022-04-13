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

variable "external_dns_deployment_name" {
  description = "Name of the external dns helm deployment"
  type        = string
  default     = "external-dns"
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

variable "orc8r_is_staging_deployment" {
  description = <<EOT
    Indicates if the orc8r-app being deploy is a staging environment.
    Staging environment does not deploy Logging, Metrics and Alerts
    EOT
  type        = bool
  default     = false
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

variable "orc8r_db_dialect" {
  description = "DB dialect for Orchestrator database connection."
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
  default     = "1.5.27"
}

variable "cwf_orc8r_chart_version" {
  description = "Version of the orchestrator cwf module Helm chart to install."
  type        = string
  default     = "0.2.2"
}

variable "feg_orc8r_chart_version" {
  description = "Version of the orchestrator feg module Helm chart to install."
  type        = string
  default     = "0.2.5"
}

variable "lte_orc8r_chart_version" {
  description = "Version of the orchestrator lte module Helm chart to install."
  type        = string
  default     = "0.2.6"
}

variable "wifi_orc8r_chart_version" {
  description = "Version of the orchestrator wifi module Helm chart to install."
  type        = string
  default     = "0.2.2"
}

variable "dp_orc8r_chart_version" {
  description = "Version of the orchestrator domain proxy module Helm chart to install."
  type        = string
  default     = "0.1.0"
}

variable "orc8r_tag" {
  description = "Image tag for Orchestrator components."
  type        = string
  default     = "1.6.1"
}

variable "magma_uuid" {
  description = "UUID to identify Orc8r deployment"
  type        = string
  default     = "default"
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

variable "efs_provisioner_name" {
  description = "Name of the efs provisioner helm deployment"
  type        = string
  default     = "efs-provisioner"
}

variable "efs_storage_class_name" {
  description = "Name of the Storage class"
  type        = string
  default     = "efs"
}

##############################################################################
# Log aggregation configuration
##############################################################################

variable "elasticsearch_endpoint" {
  description = "Endpoint of the Elasticsearch datasink for aggregated logs and events."
  type        = string
  default     = null
}

variable "elasticsearch_disk_threshold" {
  description = "Size threshold in GB."
  type        = number
  default     = 10
}

variable "elasticsearch_retention_days" {
  description = "Retention period in days of Elasticsearch indices."
  type        = number
  default     = 7
}

variable "elasticsearch_port" {
  description = "Port Elastic search is listening."
  type        = number
  default     = 443
}

variable "elasticsearch_use_ssl" {
  description = "Defines if elasicsearch curator should speak to ELK HTTP or HTTPS."
  type        = string
  default     = "True"
}

variable "elasticsearch_curator_log_level" {
  description = "Defines Elasticsearch curator logging level."
  type        = string
  default     = "INFO"
}

variable "elasticsearch_curator_name" {
  description = "Name of the elasticsearch-curator helm deployment"
  type        = string
  default     = "elasticsearch-curator"
}

variable "fluentd_deployment_name" {
  description = "Name of the fluentd helm deployment"
  type        = string
  default     = "fluentd"
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

variable "docker_registry" {
  description = "Docker registry to pull Orchestrator containers from."
  type        = string
  default     = "docker.artifactory.magmacore.org"
}

variable "docker_user" {
  description = "Docker username to login to registry with."
  type        = string
  default     = ""
}

variable "docker_pass" {
  description = "Docker registry password."
  type        = string
  default     = ""
}

variable "seed_certs_dir" {
  description = "Directory on LOCAL disk where Orchestrator certificates are stored to seed Secrets Manager values. Home directory and env vars will be expanded."
  type        = string
}

variable "helm_repo" {
  description = "Helm repository URL for Orchestrator charts."
  type        = string
  default     = "https://artifactory.magmacore.org/artifactory/helm/"
}

variable "helm_user" {
  description = "Helm username to login to repository with."
  type        = string
  default     = ""
}

variable "helm_pass" {
  description = "Helm repository password."
  type        = string
  default     = ""
}

##############################################################################
# Managed Certificates from cert-manager
##############################################################################

variable "cert_manager_route53_iam_role_arn" {
  description = "IAM role ARN for cert-manager."
  type        = string
  default     = null
}

variable "deploy_cert_manager_helm_chart" {
  description = "Deploy cert-manager helm chart."
  type        = bool
  default     = false
}

variable "managed_certs_create" {
  description = "This will generate certificates that will be stored in kubernetes secrets."
  type        = bool
  default     = false
}

variable "managed_certs_enabled" {
  description = "This will enable controller pods to use managed certificates."
  type        = bool
  default     = false
}

variable "nms_managed_certs_enabled" {
  description = "This will enable NMS nginx pod to use managed certificate."
  type        = bool
  default     = false
}

variable "nms_custom_issuer" {
  description = "Certificate issuer on Route53 for Let's Encrypt."
  type        = string
  default     = "orc8r-route53-issuer"
}

variable "managed_certs_route53_enabled" {
  description = "Use Route53 as DNS Provider."
  type        = bool
  default     = true
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
  type        = string
  default     = ""
}


variable "analytics_app_id" {
  description = "App ID for which the metrics is to be exported to"
  type        = string
  default     = ""
}

variable "analytics_metric_export_url" {
  description = "Metric Export URL"
  type        = string
  default     = ""
}

variable "analytics_category_name" {
  description = "Category under which the exported metrics will be placed under"
  type        = string
  default     = "magma"
}


##############################################################################
# Other dependency variables
##############################################################################

variable "prometheus_configurer_version" {
  description = "Image version for prometheus configurer."
  type        = string
  default     = "1.0.4"
}

variable "alertmanager_configurer_version" {
  description = "Image version for alertmanager configurer."
  type        = string
  default     = "1.0.4"
}

variable "cloudwatch_exporter_enabled" {
  description = "Enable cloudwatch exporter"
  default     = false
  type        = bool
}


##############################################################################
# Domain proxy variables
##############################################################################

variable "dp_enabled" {
  description = "Enable domain proxy"
  type        = bool
  default     = false
}

variable "dp_sas_endpoint_url" {
  description = "Sas endpoint url where to connect DP to."
  type        = string
  default     = ""
}

variable "dp_api_prefix" {
  description = "Protocol controller api prefix."
  type        = string
  default     = "/sas/v1"
}

variable "dp_sas_crt" {
  description = "SAS certificate filename."
  type        = string
  default     = "tls.crt"
}

variable "dp_sas_key" {
  description = "SAS private key filename."
  type        = string
  default     = "tls.key"
}

variable "dp_sas_ca" {
  description = "SAS CA filename."
  type        = string
  default     = "ca.crt"
}
