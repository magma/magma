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

resource "aws_route53_zone" "orc8r" {
  name = format("%s.", var.orc8r_domain_name)
}

# policy required by external dns
data "aws_iam_policy_document" "external_dns" {
  statement {
    actions = [
      "route53:ChangeResourceRecordSets",
    ]

    resources = [
      "arn:aws:route53:::hostedzone/${aws_route53_zone.orc8r.id}",
    ]
  }

  statement {
    actions = [
      "route53:ListHostedZones",
      "route53:ListResourceRecordSets",
    ]

    resources = ["*"]
  }
}

# create external dns policy from above document
resource "aws_iam_role_policy" "external_dns" {
  policy = data.aws_iam_policy_document.external_dns.json
  role   = aws_iam_role.external_dns.id
}

# allow eks workers to assume external dns role
resource "aws_iam_role" "external_dns" {
  name_prefix        = "ExternalDNSRole"
  assume_role_policy = data.aws_iam_policy_document.eks_worker_assumable.json
  tags               = var.global_tags
}
