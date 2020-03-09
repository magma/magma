################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

# k8s requires provisioner to treat efs as a persistent volume
resource "helm_release" "efs_provisioner" {
  name       = "efs-provisioner"
  repository = data.helm_repository.stable.id
  chart      = "efs-provisioner"
  version    = "0.11.0"
  namespace  = "kube-system"
  keyring    = ""

  values = [<<VALUES
  efsProvisioner:
    efsFileSystemId: ${var.efs_file_system_id}
    awsRegion: ${var.region}
    path: /pv-volume
    provisionerName: aws-efs
    storageClass:
      name: efs
  podAnnotations:
    iam-assumable-role: ${var.efs_provisioner_role_arn}
  VALUES
  ]
}
