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

data "aws_iam_policy_document" "es-management" {
  count = var.deploy_elasticsearch ? 1 : 0

  statement {
    effect = "Allow"

    actions = [
      "es:*",
    ]

    principals {
      identifiers = ["*"]
      type        = "AWS"
    }

    resources = [
      "${aws_elasticsearch_domain.es[0].arn}/*",
    ]
  }
}

resource "aws_elasticsearch_domain_policy" "es_management_access" {
  count = var.deploy_elasticsearch ? 1 : 0

  domain_name     = aws_elasticsearch_domain.es[0].domain_name
  access_policies = data.aws_iam_policy_document.es-management[0].json
}

resource "aws_iam_service_linked_role" "es" {
  count = var.deploy_elasticsearch && var.deploy_elasticsearch_service_linked_role ? 1 : 0

  aws_service_name = "es.amazonaws.com"
}

locals {
  elasticsearch_available_subnets = length(module.vpc.private_subnets) > 0 ? module.vpc.private_subnets : module.vpc.public_subnets
}

resource "aws_elasticsearch_domain" "es" {
  count = var.deploy_elasticsearch ? 1 : 0

  domain_name           = var.elasticsearch_domain_name
  elasticsearch_version = var.elasticsearch_version

  cluster_config {
    instance_type  = var.elasticsearch_instance_type
    instance_count = var.elasticsearch_instance_count

    dedicated_master_enabled = var.elasticsearch_dedicated_master_enabled
    dedicated_master_type    = var.elasticsearch_dedicated_master_enabled ? var.elasticsearch_dedicated_master_type : null
    dedicated_master_count   = var.elasticsearch_dedicated_master_enabled ? var.elasticsearch_dedicated_master_count : null

    zone_awareness_enabled = true
    zone_awareness_config {
      availability_zone_count = var.elasticsearch_az_count
    }
  }

  advanced_options = {
    "rest.action.multi.allow_explicit_index" = "true"
  }

  vpc_options {
    subnet_ids         = slice(local.elasticsearch_available_subnets, 0, min(var.elasticsearch_az_count, 3))
    security_group_ids = [aws_security_group.default.id]
  }

  dynamic "ebs_options" {
    for_each = var.elasticsearch_ebs_enabled ? [true] : []
    content {
      ebs_enabled = true
      volume_size = var.elasticsearch_ebs_volume_size
      volume_type = var.elasticsearch_ebs_volume_type
      iops        = var.elasticsearch_ebs_iops
    }
  }

  snapshot_options {
    automated_snapshot_start_hour = 0
  }

  tags = merge(
    var.global_tags,
    var.elasticsearch_domain_tags,
  )

  depends_on = [aws_iam_service_linked_role.es]
}
