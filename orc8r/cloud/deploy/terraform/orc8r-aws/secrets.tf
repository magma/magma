################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

resource "aws_secretsmanager_secret" "orc8r_secrets" {
  name = var.secretsmanager_orc8r_secret
}
