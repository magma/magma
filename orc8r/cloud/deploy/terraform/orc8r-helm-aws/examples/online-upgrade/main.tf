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

module orc8r {
  source = "../../../orc8r-aws"

  region = var.region

  # VPC
  vpc_name             = var.vpc_name
  vpc_cidr             = var.vpc_configuration.cidr
  vpc_public_subnets   = var.vpc_configuration.public_subnets
  vpc_private_subnets  = var.vpc_configuration.private_subnets
  vpc_database_subnets = var.vpc_configuration.db_subnets

  # RDS
  orc8r_db_identifier     = var.orc8r_db_configuration.identifier
  orc8r_db_storage_gb     = var.orc8r_db_configuration.storage_gb
  orc8r_db_engine_version = var.orc8r_db_configuration.engine_version
  orc8r_db_instance_class = var.orc8r_db_configuration.instance_class
  orc8r_db_password       = var.orc8r_db_password

  secretsmanager_orc8r_secret = var.secretsmanager_secret_name

  orc8r_domain_name = var.orc8r_domain
  cluster_name      = var.eks_cluster_name

  eks_worker_group_key              = var.ssh_key_name
  eks_worker_additional_policy_arns = [aws_iam_policy.worker_node_policy.arn]

  # Once you are fully migrated to v1.1+ on the application side, you can
  # safely delete this custom worker groups configuration or set it to
  # something else. On v1.1+ we don't reserve any worker nodes for specific
  # payloads since we migrated to EFS for prometheus configuration and storage.
  # IMPORTANT: Make sure the first list matches your v1.0 Terraform module!
  eks_worker_groups = concat(
    [
      {
        name                 = "wg-1"
        instance_type        = "t3.small"
        asg_desired_capacity = 2
        kubelet_extra_args   = "--node-labels=worker-type=controller"

        tags = [
          {
            key                 = "orc8r-node-type"
            value               = "orc8r-worker-node"
            propagate_at_launch = true
          },
        ]
      },
      {
        name                 = "wg-metrics"
        instance_type        = "t3.medium"
        asg_desired_capacity = 1
        kubelet_extra_args   = "--node-labels=worker-type=metrics"

        subnets = [var.metrics_worker_subnet_id]

        additional_userdata = data.template_file.metrics_userdata.rendered

        tags = [
          {
            key                 = "orc8r-node-type"
            value               = "orc8r-prometheus-node"
            propagate_at_launch = true
          },
        ]
      },
    ],
    var.additional_eks_worker_groups,
  )

  eks_map_users = var.eks_map_users

  # ES
  deploy_elasticsearch                     = var.deploy_elasticsearch
  elasticsearch_domain_name                = var.elasticsearch_domain_name
  elasticsearch_version                    = var.elasticsearch_domain_configuration.version
  elasticsearch_instance_type              = var.elasticsearch_domain_configuration.instance_type
  elasticsearch_instance_count             = var.elasticsearch_domain_configuration.instance_count
  elasticsearch_az_count                   = var.elasticsearch_domain_configuration.az_count
  elasticsearch_ebs_enabled                = var.elasticsearch_domain_configuration.ebs_enabled
  elasticsearch_ebs_volume_size            = var.elasticsearch_domain_configuration.ebs_volume_size
  elasticsearch_ebs_volume_type            = var.elasticsearch_domain_configuration.ebs_volume_type
  deploy_elasticsearch_service_linked_role = var.deploy_elasticsearch_linked_role
}

module orc8r-app {
  source = "../.."

  region = var.region

  orc8r_domain_name     = module.orc8r.orc8r_domain_name
  orc8r_route53_zone_id = module.orc8r.route53_zone_id
  external_dns_role_arn = module.orc8r.external_dns_role_arn

  secretsmanager_orc8r_name = module.orc8r.secretsmanager_secret_name
  seed_certs_dir            = var.seed_certs_dir

  # Tiller would already have been installed with v1.0
  install_tiller                       = false
  existing_tiller_service_account_name = "tiller"
  helm_deployment_name                 = var.new_deployment_name

  orc8r_db_host    = module.orc8r.orc8r_db_host
  orc8r_db_dialect = module.orc8r.orc8r_db_dialect
  orc8r_db_name    = module.orc8r.orc8r_db_name
  orc8r_db_user    = module.orc8r.orc8r_db_user
  orc8r_db_pass    = module.orc8r.orc8r_db_pass

  docker_registry = var.docker_registry
  docker_user     = var.docker_user
  docker_pass     = var.docker_pass

  helm_repo = var.helm_repo
  helm_user = var.helm_user
  helm_pass = var.helm_pass

  eks_cluster_id = module.orc8r.eks_cluster_id

  efs_file_system_id       = module.orc8r.efs_file_system_id
  efs_provisioner_role_arn = module.orc8r.efs_provisioner_role_arn

  elasticsearch_endpoint = module.orc8r.es_endpoint

  orc8r_chart_version       = var.orc8r_chart_version
  orc8r_tag                 = var.orc8r_container_tag
  orc8r_controller_replicas = var.orc8r_controller_replicas
  orc8r_proxy_replicas      = var.orc8r_proxy_replicas

  deploy_nms = var.deploy_nms
}

output "nameservers" {
  value = module.orc8r.route53_nameservers
}

output "vals" {
  value     = module.orc8r-app.helm_vals
  sensitive = true
}
