variable "region" {
  description = "AWS region to deploy to. This should match your v1.0 Terraform."
  type        = string
}

variable "vpc_name" {
  description = "Name of the VPC that the v1.0 infra is deployed in. This should match your v1.0 Terraform."
  type        = string
}

# The v1.0 terraform didn't specify any private subnets and put the EKS nodes
# into the public subnets. If you specify private subnets here after
# migrating, the worker nodes will be reorganized into the new private
# subnets when you terminate them and autoscaling brings them back up.
# This may result in some downtime (a rolling restart with 2+ proxy and
# controller replicas will mostly avoid this).
variable "vpc_configuration" {
  description = "Configuration of the VPC that the v1.0 chart is deployed in. This should match your v1.0 Terraform."
  type = object({
    cidr            = string
    public_subnets  = list(string)
    private_subnets = list(string)
    db_subnets      = list(string)
  })
  default = {
    cidr            = "10.10.0.0/16"
    public_subnets  = ["10.10.1.0/24", "10.10.2.0/24", "10.10.3.0/24"]
    private_subnets = []
    db_subnets      = ["10.10.11.0/24", "10.10.12.0/24", "10.10.13.0/24"]
  }
}

variable "orc8r_db_configuration" {
  description = "Configuration of the Orchestrator Postgres instance. This should match the v1.0 Terraform."
  type = object({
    identifier     = string
    storage_gb     = number
    engine_version = string
    instance_class = string
  })
  default = {
    identifier     = "orc8rdb"
    storage_gb     = 32
    engine_version = "9.6.11"
    instance_class = "db.m4.large"
  }
}

variable "orc8r_db_password" {
  description = "Password for the Orchestrator Postgres instance. This should match the v1.0 Terraform."
  type        = string
}

variable "secretsmanager_secret_name" {
  description = "Name for the Secretsmanager secret that the orc8r-aws module will create."
  type        = string
}

variable "orc8r_domain" {
  description = "Root domain or subdomain for your Orchestrator deployment (e.g. orc8r.mydomain.com)."
  type        = string
}

variable "ssh_key_name" {
  description = "Name of the SSH key you created for the v1.0 infra."
  type        = string
}

variable "eks_cluster_name" {
  description = "Name of the EKS cluster that the v1.0 application is deployed on. This should match your v1.0 Terraform."
  type        = string
}

variable "additional_eks_worker_groups" {
  description = "Additional EKS worker nodes to spin up while the v1.1.0 application is deployed concurrently with the v1.0 application."
  type        = any
  default = [
    {
      name                 = "wg-1"
      instance_type        = "t3.large"
      asg_min_size         = 3
      asg_desired_capacity = 3

      tags = [
        {
          key                 = "orc8r-node-type"
          value               = "orc8r-worker-node"
          propagate_at_launch = true
        },
      ]
    }
  ]
}

variable "eks_map_users" {
  description = "Additional users you want to grant access to EKS to. This should match your v1.0 Terraform or those users will lose k8s access."
  type        = any
  default     = []
}

variable "deploy_elasticsearch" {
  description = "Deploy elasticsearch cluster for log aggregation (default false)."
  type        = bool
  default     = false
}

# Note: the elasticsearch service linked role only needs to be created once
# per account, so if you've already created this set this to false
variable "deploy_elasticsearch_linked_role" {
  description = "Deploy ES linked role if ES is deployed."
  type        = bool
  default     = true
}

variable "elasticsearch_domain_name" {
  description = "Name for ES domain"
  type        = string
  default     = "orc8r-es-domain"
}

variable "elasticsearch_domain_configuration" {
  description = "Configuration for the ES domain"
  type = object({
    version         = string
    instance_type   = string
    instance_count  = number
    az_count        = number
    ebs_enabled     = bool
    ebs_volume_size = number
    ebs_volume_type = string
  })
  default = {
    version         = "7.4"
    instance_type   = "t2.medium.elasticsearch"
    instance_count  = 3
    az_count        = 3
    ebs_enabled     = true
    ebs_volume_size = 32
    ebs_volume_type = "gp2"
  }
}

variable "docker_registry" {
  description = "URL to your Docker registry"
  type        = string
}

variable "docker_user" {
  description = "Username for your Docker registry"
  type        = string
}

variable "docker_pass" {
  description = "Password for your Docker user"
  type        = string
}

variable "helm_repo" {
  description = "URL to your Helm repo. Don't forget the protocol prefix (e.g. https://)"
  type        = string
}

variable "helm_user" {
  description = "Username for your Helm repo"
  type        = string
}

variable "helm_pass" {
  description = "Password for your Helm user"
  type        = string
}

variable "seed_certs_dir" {
  description = "Directory with your Orchestrator certificates."
  type        = string
}

variable "new_deployment_name" {
  description = "New name for the v1.1.0 Helm deployment. This must be different than your old v1.0 deployment (which was probably 'orc8r')"
  type        = string
}

variable "orc8r_chart_version" {
  description = "Chart version for the Helm deployment"
  type        = string
  default     = "1.4.21"
}

variable "orc8r_container_tag" {
  description = "Container tag to deploy"
  type        = string
}

variable "orc8r_controller_replicas" {
  description = "How many controller pod replicas to deploy"
  type        = number
  default     = 2
}

variable "orc8r_proxy_replicas" {
  description = "How many proxy pod replicas to deploy"
  type        = number
  default     = 2
}

variable "deploy_nms" {
  description = "Whether to deploy NMS. You can leave this set to true for the online upgrade, unlike the from-scratch v1.1.0 installation."
  type        = bool
  default     = true
}

variable "worker_node_policy_suffix" {
  description = "The name suffix of the custom IAM node policy from the v1.0 Terraform root module. This policy name will begin with magma_eks_worker_node_policy."
  type        = string
}

variable "prometheus_ebs_az" {
  description = "Availability zone that the Prometheus worker node and EBS volume are located in. Find this in the EC2 console."
  type        = string
}

variable "prometheus_ebs_size" {
  description = "Size of the EBS volume for the Prometheus data EBS volume. This should match your v1.0 Terraform."
  type        = number
  default     = 64
}

variable "metrics_worker_subnet_id" {
  description = "Subnet ID of the metrics worker instance. Find this in the EC2 console (the instance will have the tag orc8r-node-type: orc8r-prometheus-node)."
  type        = string
}
