################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

# efs file system for eks persistent volumes
resource "aws_efs_file_system" "eks_pv" {
  tags = {
    Name = "${var.efs_project_name}.k8s.pv.local"
  }
}

# efs mount target for eks persistent volumes
resource "aws_efs_mount_target" "eks_pv_mnt" {
  file_system_id  = aws_efs_file_system.eks_pv.id
  security_groups = [aws_security_group.default.id]
  subnet_id       = length(var.vpc_private_subnets) > 0 ? module.vpc.private_subnets[count.index] : module.vpc.public_subnets[count.index]
  count           = length(var.vpc_private_subnets) > 0 ? length(var.vpc_private_subnets) : length(var.vpc_public_subnets)
}

# allow eks workers to assume efs provisioner role
resource "aws_iam_role" "efs_provisioner" {
  name_prefix        = "EFSProvisionerRole"
  assume_role_policy = data.aws_iam_policy_document.eks_worker_assumable.json
  tags               = var.global_tags
}

# grant efs read only policy to efs provisioner
resource "aws_iam_role_policy_attachment" "efs_provisioner" {
  policy_arn = "arn:aws:iam::aws:policy/AmazonElasticFileSystemReadOnlyAccess"
  role       = aws_iam_role.efs_provisioner.id
}
