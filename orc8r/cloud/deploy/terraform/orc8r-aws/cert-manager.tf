resource "random_id" "cert_manager_route53_random_id" {
  count = var.setup_cert_manager ? 1 : 0

  byte_length = 8
}

data "aws_iam_policy_document" "cert_manager_route53_iam_policy_document" {
  count = var.setup_cert_manager ? 1 : 0

  statement {
    actions = [
      "route53:GetChange"
    ]
    resources = [
      "arn:aws:route53:::change/*"
    ]
  }

  statement {
    actions = [
      "route53:ChangeResourceRecordSets",
      "route53:ListResourceRecordSets"
    ]
    resources = [
      "arn:aws:route53:::hostedzone/*"
    ]
  }

  statement {
    actions = [
      "route53:ListHostedZonesByName"
    ]
    resources = [
      "*"
    ]
  }
}

data "aws_iam_policy_document" "cert_manager_route53_iam_policy_document_assume" {
  count = var.setup_cert_manager ? 1 : 0

  statement {
    actions = ["sts:AssumeRoleWithWebIdentity"]

    principals {
      type        = "Federated"
      identifiers = [module.eks.oidc_provider_arn]
    }

    condition {
      test     = "StringEquals"
      variable = "${replace(module.eks.cluster_oidc_issuer_url, "https://", "")}:sub"

      values = [
        "system:serviceaccount:cert-manager:cert-manager",
      ]
    }
  }
}

resource "aws_iam_policy" "cert_manager_route53_iam_policy" {
  count = var.setup_cert_manager ? 1 : 0

  name   = "cert_manager_route53_iam_policy-${resource.random_id.cert_manager_route53_random_id.0.id}"
  path   = "/"
  policy = data.aws_iam_policy_document.cert_manager_route53_iam_policy_document.0.json
}

resource "aws_iam_role" "cert_manager_route53_iam_role" {
  count = var.setup_cert_manager ? 1 : 0

  name                = "cert_manager_route53_iam_role-${resource.random_id.cert_manager_route53_random_id.0.id}"
  assume_role_policy  = data.aws_iam_policy_document.cert_manager_route53_iam_policy_document_assume.0.json
  managed_policy_arns = [aws_iam_policy.cert_manager_route53_iam_policy.0.arn]
}
