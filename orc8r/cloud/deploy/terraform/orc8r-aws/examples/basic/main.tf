################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

module orc8r {
  source = "../.."

  region = "us-west-2"

  nms_db_password                   = "Faceb00k12345"
  orc8r_db_password                 = "Faceb00k12345"
  secretsmanager_artifactory_secret = "magma-orc8r-test"
  deployment_secrets_bucket         = "magma.orc8r.test"
  orc8r_domain_name                 = "orc8r.magma.test"

  vpc_name     = "magma-orc8r-test"
  cluster_name = "orc8r-test"
}
