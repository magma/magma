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

module "orc8r" {
  # Change this to pull from GitHub with a specified ref, e.g.
  # source = "github.com/magma/magma//orc8r/cloud/deploy/terraform/orc8r-aws?ref=v1.6"
  source = "../../../orc8r-aws"

  region = "us-west-2"

  # If you performing a fresh Orc8r install, choose a recent Postgres version
  # orc8r_db_engine_version     = "12.6"
  orc8r_db_password = "mypassword" # must be at least 8 characters

  secretsmanager_orc8r_secret = "orc8r-secrets"
  orc8r_domain_name           = "orc8r.example.com"

  orc8r_sns_email             = "admin@example.com"
  enable_aws_db_notifications = true

  vpc_name        = "orc8r"
  cluster_name    = "orc8r"
  cluster_version = "1.17"

  deploy_elasticsearch          = true
  elasticsearch_domain_name     = "orc8r-es"
  elasticsearch_version         = "7.7"
  elasticsearch_instance_type   = "t2.medium.elasticsearch"
  elasticsearch_instance_count  = 2
  elasticsearch_az_count        = 2
  elasticsearch_ebs_enabled     = true
  elasticsearch_ebs_volume_size = 32
  elasticsearch_ebs_volume_type = "gp2"
}

module "orc8r-app" {
  # Change this to pull from GitHub with a specified ref, e.g.
  # source = "github.com/magma/magma//orc8r/cloud/deploy/terraform/orc8r-helm-aws?ref=v1.6"
  source = "../.."

  region = "us-west-2"

  orc8r_domain_name     = module.orc8r.orc8r_domain_name
  orc8r_route53_zone_id = module.orc8r.route53_zone_id
  external_dns_role_arn = module.orc8r.external_dns_role_arn

  secretsmanager_orc8r_name = module.orc8r.secretsmanager_secret_name
  seed_certs_dir            = "~/secrets/certs"

  orc8r_db_host    = module.orc8r.orc8r_db_host
  orc8r_db_port    = module.orc8r.orc8r_db_port
  orc8r_db_dialect = module.orc8r.orc8r_db_dialect
  orc8r_db_name    = module.orc8r.orc8r_db_name
  orc8r_db_user    = module.orc8r.orc8r_db_user
  orc8r_db_pass    = module.orc8r.orc8r_db_pass

  # Note that this can be any container registry provider
  docker_registry = "https://docker.artifactory.magmacore.org/artifactory/docker"
  docker_user     = ""
  docker_pass     = ""

  # Note that this can be any Helm chart repo provider
  helm_repo       = "https://docker.artifactory.magmacore.org/artifactory/helm"
  helm_user       = ""
  helm_pass       = ""
  eks_cluster_id = module.orc8r.eks_cluster_id

  efs_file_system_id       = module.orc8r.efs_file_system_id
  efs_provisioner_role_arn = module.orc8r.efs_provisioner_role_arn

  elasticsearch_endpoint       = module.orc8r.es_endpoint
  elasticsearch_disk_threshold = tonumber(module.orc8r.es_volume_size * 75 / 100)

  orc8r_deployment_type = "fwa"
  orc8r_tag             = "1.6.0"
}

output "nameservers" {
  value = module.orc8r.route53_nameservers
}
