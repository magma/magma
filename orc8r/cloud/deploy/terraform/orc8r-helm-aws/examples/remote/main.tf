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

locals {
  region = "us-west-2"
  domain = "my.domain.com"
}

# Setup remote state and root secrets
# You will have to manually create the S3 bucket for remote state to work
terraform {
  backend "s3" {
    bucket = "orc8r.release.test.deployment"
    key    = "terraform/terraform.tfstate"
    # Unfortunately terraform doesn't support using locals in backend blocks
    dynamodb_table = "my-dynamodb-table"
    region         = "us-west-2"
  }
}

provider "aws" {
  version = ">= 2.6.0"
  region  = local.region
}

# Lock table for remote terraform state
resource "aws_dynamodb_table" "terraform_locks" {
  name         = "my-dynamodb-table"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "LockID"

  attribute {
    name = "LockID"
    type = "S"
  }
}

# This secretsmanager secret needs to be manually created and populated in the
# AWS console. For this example, you would set the following key-value pairs:
#   orc8r_db_pass
#   docker_registry
#   docker_user
#   docker_pass
#   helm_repo
#   helm_user
#   helm_pass
data "aws_secretsmanager_secret" "root_secrets" {
  name = "orc8r_root_secrets"
}

data "aws_secretsmanager_secret_version" "root_secrets" {
  secret_id = data.aws_secretsmanager_secret.root_secrets.id
}

module "orc8r" {
  # Change this to pull from github with a specified ref
  source = "../../../orc8r-aws"

  region = local.region

  orc8r_db_password           = jsondecode(data.aws_secretsmanager_secret_version.root_secrets.secret_string)["orc8r_db_pass"]
  secretsmanager_orc8r_secret = "orc8r-secrets"
  orc8r_domain_name           = "orc8r.example.com"

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
  # Change this to pull from github with a specified ref
  source = "../.."

  region = local.region

  # This has to match the backend block declared at the top. Unfortunately we
  # have to duplicate this because Terraform evaluates backend blocks before
  # the rest of the module.
  state_backend = "s3"
  state_config = {
    bucket         = "orc8r.release.test.deployment"
    key            = "terraform/terraform.tfstate"
    dynamodb_table = "orc8r-release-test-tf-lock"
    region         = "us-west-2"
  }

  orc8r_domain_name     = module.orc8r.orc8r_domain_name
  orc8r_route53_zone_id = module.orc8r.route53_zone_id
  external_dns_role_arn = module.orc8r.external_dns_role_arn

  secretsmanager_orc8r_name = module.orc8r.secretsmanager_secret_name
  seed_certs_dir            = "~/orc8r.test.secrets/certs"

  orc8r_db_host    = module.orc8r.orc8r_db_host
  orc8r_db_port    = module.orc8r.orc8r_db_port
  orc8r_db_dialect = module.orc8r.orc8r_db_dialect
  orc8r_db_name    = module.orc8r.orc8r_db_name
  orc8r_db_user    = module.orc8r.orc8r_db_user
  orc8r_db_pass    = module.orc8r.orc8r_db_pass

  docker_registry = jsondecode(data.aws_secretsmanager_secret_version.root_secrets.secret_string)["docker_registry"]
  docker_user     = jsondecode(data.aws_secretsmanager_secret_version.root_secrets.secret_string)["docker_user"]
  docker_pass     = jsondecode(data.aws_secretsmanager_secret_version.root_secrets.secret_string)["docker_pass"]

  helm_repo = jsondecode(data.aws_secretsmanager_secret_version.root_secrets.secret_string)["helm_repo"]
  helm_user = jsondecode(data.aws_secretsmanager_secret_version.root_secrets.secret_string)["helm_user"]
  helm_pass = jsondecode(data.aws_secretsmanager_secret_version.root_secrets.secret_string)["helm_pass"]

  eks_cluster_id = module.orc8r.eks_cluster_id

  efs_file_system_id       = module.orc8r.efs_file_system_id
  efs_provisioner_role_arn = module.orc8r.efs_provisioner_role_arn

  elasticsearch_endpoint = module.orc8r.es_endpoint

  orc8r_deployment_type = "fwa"
  orc8r_tag             = "1.4.0"
  deploy_nms            = true
}

output "nameservers" {
  value = module.orc8r.route53_nameservers
}
