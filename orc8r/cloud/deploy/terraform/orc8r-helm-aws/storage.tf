################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

resource "kubernetes_persistent_volume_claim" "storage" {
  for_each = {
    promcfg = {
      access_mode = "ReadWriteMany"
      storage     = "1Gi"
    }
    promdata = {
      access_mode = "ReadWriteOnce"
      storage     = "64Gi"
    }
    grafanadata = {
      access_mode = "ReadWriteMany"
      storage     = "2Gi"
    }
    grafanadashboards = {
      access_mode = "ReadWriteMany"
      storage     = "2Gi"
    }
    grafanaproviders = {
      access_mode = "ReadWriteMany"
      storage     = "100M"
    }
    grafanadatasources = {
      access_mode = "ReadWriteMany"
      storage     = "100M"
    }
    openvpn = {
      access_mode = "ReadWriteOnce"
      storage     = "2M"
    }
  }

  metadata {
    name      = each.key
    namespace = kubernetes_namespace.orc8r.metadata[0].name
  }

  spec {
    access_modes = [each.value.access_mode]
    resources {
      requests = {
        storage = each.value.storage
      }
    }
    storage_class_name = "efs"
  }

  depends_on = [helm_release.efs_provisioner]
}
