################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

module "vpc" {
  source  = "terraform-aws-modules/vpc/aws"
  version = "~> 2.17.0"

  name = var.vpc_name
  azs  = data.aws_availability_zones.available.names

  cidr             = var.vpc_cidr
  public_subnets   = var.vpc_public_subnets
  private_subnets  = var.vpc_private_subnets
  database_subnets = var.vpc_database_subnets

  create_database_subnet_group = true
  enable_dns_hostnames         = true
  enable_nat_gateway           = length(var.vpc_private_subnets) > 0 ? true : false
  single_nat_gateway           = length(var.vpc_private_subnets) > 0 ? true : false

  tags = merge(
    var.vpc_extra_tags,
    var.global_tags,
    {
      "kubernetes.io/cluster/${var.cluster_name}" = "shared"
    },
  )

  public_subnet_tags = {
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
    "kubernetes.io/role/elb"                    = 1
  }

  private_subnet_tags = {
    "kubernetes.io/cluster/${var.cluster_name}" = "shared"
    "kubernetes.io/role/elb"                    = 1
  }
}
