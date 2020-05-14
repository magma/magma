################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

# stable helm repository
data "helm_repository" "stable" {
  name = "stable"
  url  = "https://kubernetes-charts.storage.googleapis.com"
}

# incubator helm repository
data "helm_repository" "incubator" {
  name = "incubator"
  url  = "http://storage.googleapis.com/kubernetes-charts-incubator"
}

# helm tiller service account
resource "kubernetes_service_account" "tiller" {
  count = var.existing_tiller_service_account_name == null ? 1 : 0

  metadata {
    name      = "tiller"
    namespace = var.tiller_namespace
  }

  automount_service_account_token = true
}

# helm tiller cluster role
resource "kubernetes_cluster_role_binding" "tiller" {
  count = var.existing_tiller_service_account_name == null ? 1 : 0

  metadata {
    name = "tiller"
  }

  role_ref {
    kind      = "ClusterRole"
    name      = "cluster-admin"
    api_group = "rbac.authorization.k8s.io"
  }

  subject {
    kind      = "ServiceAccount"
    name      = "tiller"
    api_group = ""
    namespace = var.tiller_namespace
  }
}
