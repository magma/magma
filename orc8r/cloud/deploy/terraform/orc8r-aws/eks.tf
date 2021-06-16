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

resource "tls_private_key" "eks_workers" {
  count = var.eks_worker_group_key == null ? 1 : 0

  algorithm = "RSA"
}

resource "aws_key_pair" "eks_workers" {
  count = var.eks_worker_group_key == null ? 1 : 0

  key_name_prefix = var.cluster_name
  public_key      = tls_private_key.eks_workers[0].public_key_openssh
}

module "eks" {
  source  = "terraform-aws-modules/eks/aws"
  version = "~> 17.0.3"

  cluster_name    = var.cluster_name
  cluster_version = var.cluster_version

  vpc_id       = module.vpc.vpc_id
  subnets      = length(module.vpc.private_subnets) > 0 ? module.vpc.private_subnets : module.vpc.public_subnets

  cluster_enabled_log_types = [
    "api",
    "audit",
    "authenticator",
    "controllerManager",
    "scheduler",
  ]

  workers_group_defaults = {
    key_name = var.eks_worker_group_key == null ? aws_key_pair.eks_workers[0].key_name : var.eks_worker_group_key
  }
  worker_additional_security_group_ids = concat([aws_security_group.default.id], var.eks_worker_additional_sg_ids)
  workers_additional_policies          = var.eks_worker_additional_policy_arns
  worker_groups                        = var.thanos_enabled ? concat(var.eks_worker_groups, var.thanos_worker_groups) : var.eks_worker_groups

  map_roles = var.eks_map_roles
  map_users = var.eks_map_users

  tags = var.global_tags
}

# role assume policy for eks workers
data "aws_iam_policy_document" "eks_worker_assumable" {
  statement {
    principals {
      identifiers = ["ec2.amazonaws.com"]
      type        = "Service"
    }
    actions = ["sts:AssumeRole"]
  }

  statement {
    principals {
      identifiers = [module.eks.worker_iam_role_arn]
      type        = "AWS"
    }
    actions = ["sts:AssumeRole"]
  }
}
