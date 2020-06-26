################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

output "helm_vals" {
  description = "Helm values for the orc8r deployment"
  value       = helm_release.orc8r.values
  sensitive   = true
}
