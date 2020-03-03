################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

output "eks_cluster_id" {
  description = "Cluster ID for the EKS cluster"
  value       = module.eks.cluster_id
}

output "kubeconfig" {
  description = "kubectl config file to access the EKS cluster"
  value       = module.eks.kubeconfig
}

output "eks_config_map_aws_auth" {
  description = "A k8s configmap to allow authentication to the EKS cluster."
  value       = module.eks.config_map_aws_auth
  sensitive   = true
}

output "s3_secret_bucket" {
  description = "Name of the S3 bucket for Orchestrator secrets"
  value       = aws_s3_bucket.orc8r_bucket.id
}

output "efs_file_system_id" {
  description = "ID of the EFS file system created for k8s persistent volumes."
  value       = aws_efs_file_system.eks_pv.id
}

output "efs_provisioner_role_arn" {
  description = "ARN of the IAM role for the EFS provisioner."
  value       = aws_iam_role.efs_provisioner.arn
}
