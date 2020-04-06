################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

provider "aws" {
  version = ">= 2.6.0"
  region  = var.region
}

provider "random" {
  version = "~> 2.1"
}

provider "tls" {
  version = "~> 2.1"
}

provider "template" {
  version = "~> 2.0"
}

data "aws_eks_cluster" "cluster" {
  name = var.eks_cluster_id
}

data "aws_eks_cluster_auth" "cluster" {
  name = var.eks_cluster_id
}

provider "kubernetes" {
  host                   = data.aws_eks_cluster.cluster.endpoint
  cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
  token                  = data.aws_eks_cluster_auth.cluster.token
  load_config_file       = false
  # See https://github.com/terraform-providers/terraform-provider-kubernetes/issues/759
  version = "~> 1.10.0"
}

provider "helm" {
  kubernetes {
    host                   = data.aws_eks_cluster.cluster.endpoint
    cluster_ca_certificate = base64decode(data.aws_eks_cluster.cluster.certificate_authority.0.data)
    token                  = data.aws_eks_cluster_auth.cluster.token
    load_config_file       = false
  }

  service_account = kubernetes_service_account.tiller.metadata.0.name
  namespace       = kubernetes_service_account.tiller.metadata.0.namespace
  tiller_image    = "gcr.io/kubernetes-helm/tiller:v2.16.3"
  install_tiller  = var.install_tiller
  max_history     = 100

  version = "~> 0.10"
}

