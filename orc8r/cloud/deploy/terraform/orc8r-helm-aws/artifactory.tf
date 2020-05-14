################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

locals {
  dockercfg = {
    (var.docker_registry) = {
      username = var.docker_user
      password = var.docker_pass
    }
  }
}

resource "kubernetes_secret" "artifactory" {
  metadata {
    name      = "artifactory"
    namespace = kubernetes_namespace.orc8r.metadata[0].name
  }

  data = { ".dockercfg" = jsonencode(local.dockercfg) }
  type = "kubernetes.io/dockercfg"
}

data "helm_repository" "artifactory" {
  name     = "orc8r_artifactory"
  url      = var.helm_repo
  username = var.helm_user
  password = var.helm_pass
}
