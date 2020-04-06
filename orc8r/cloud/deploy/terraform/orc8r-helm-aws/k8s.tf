################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

resource "kubernetes_namespace" "orc8r" {
  metadata {
    name = var.orc8r_kubernetes_namespace
  }
}

# external dns maps route53 to ingress rosources
resource "helm_release" "external_dns" {
  name       = "external-dns"
  repository = data.helm_repository.stable.id
  chart      = "external-dns"
  version    = "2.19.1"
  namespace  = "kube-system"
  keyring    = ""

  values = [<<VALUES
  rbac:
    create: true
  aws:
    assumeRoleArn: ${var.external_dns_role_arn}
  zoneIdFilters:
    - ${var.orc8r_route53_zone_id}
  VALUES
  ]
}
