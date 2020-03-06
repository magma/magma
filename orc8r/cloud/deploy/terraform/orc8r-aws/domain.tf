################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

resource "aws_route53_zone" "orc8r" {
  name = format("%s.", var.orc8r_domain_name)
}
