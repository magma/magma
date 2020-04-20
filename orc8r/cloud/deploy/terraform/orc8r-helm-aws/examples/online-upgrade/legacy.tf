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

# The following resources come from the v1.0 terraform file. You can safely
# delete them once you're fully established on v1.1+ and purge the 1.0 Helm
# release from the old namespace.
data "template_file" "metrics_userdata" {
  template = file("${path.module}/scripts/prepare_metrics_instance.sh.tpl")
}

data "aws_iam_policy_document" "worker_node_policy_doc" {
  statement {
    effect = "Allow"

    actions = [
      "ec2:DescribeVolumes",
      "ec2:AttachVolume",
      "ec2:DetachVolume",
    ]

    resources = [
      "arn:aws:ec2:*:*:volume/*",
      "arn:aws:ec2:*:*:instance/*",
    ]
  }
}

# Grab the full IAM policy name from the IAM console under "Policies"
resource "aws_iam_policy" "worker_node_policy" {
  name   = format("magma_eks_worker_node_policy-%s", var.worker_node_policy_suffix)
  policy = data.aws_iam_policy_document.worker_node_policy_doc.json
}

# EBS volume for prometheus metrics.
# Copy the AZ name from the EC2 console for this EBS volume
resource "aws_ebs_volume" "prometheus-ebs-eks" {
  availability_zone = var.prometheus_ebs_az
  size              = var.prometheus_ebs_size

  tags = {
    Name = "orc8r-prometheus-data"
  }
}

# EBS volume for prometheus configs.
# Copy the AZ name from the EC2 console for this EBS volume
resource "aws_ebs_volume" "prometheus-configs-ebs-eks" {
  availability_zone = var.prometheus_ebs_az
  size              = 1

  tags = {
    Name = "orc8r-prometheus-configs"
  }
}
