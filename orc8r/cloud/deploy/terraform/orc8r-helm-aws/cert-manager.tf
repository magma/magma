resource "helm_release" "cert-manager" {
  version          = "1.6.1"
  name             = "cert-manager"
  chart            = "cert-manager"
  namespace        = "cert-manager"
  repository       = "https://charts.jetstack.io"
  create_namespace = true

  set {
    name  = "installCRDs"
    value = "true"
  }
}

resource "random_id" "cert_manager_random_id" {
  byte_length = 8
}

resource "aws_iam_user" "cert_manager_iam_user" {
  name = "cert_manager_iam_user-${resource.random_id.cert_manager_random_id.id}"
}

resource "aws_iam_access_key" "cert_manager_iam_access_key" {
  user = aws_iam_user.cert_manager_iam_user.name
}

resource "aws_iam_user_policy" "cert_manager_iam_user_policy" {
  name = "cert_manager_user_policy-${resource.random_id.cert_manager_random_id.id}"
  user = aws_iam_user.cert_manager_iam_user.name

  policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Action": "route53:GetChange",
      "Resource": "arn:aws:route53:::change/*"
    },
    {
      "Effect": "Allow",
      "Action": [
        "route53:ChangeResourceRecordSets",
        "route53:ListResourceRecordSets"
      ],
      "Resource": "arn:aws:route53:::hostedzone/*"
    },
    {
      "Effect": "Allow",
      "Action": "route53:ListHostedZonesByName",
      "Resource": "*"
    }
  ]
}
EOF
}
