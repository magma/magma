################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
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
