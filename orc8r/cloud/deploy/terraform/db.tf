################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

resource "aws_db_instance" "default" {
  identifier        = "orc8rdb"
  allocated_storage = 32
  engine            = "postgres"
  engine_version    = "9.6.11"
  instance_class    = "db.m4.large"

  name     = "orc8r"
  username = "orc8r"
  password = var.db_password

  vpc_security_group_ids = [aws_security_group.default.id]

  db_subnet_group_name = module.vpc.database_subnet_group

  skip_final_snapshot = true
  # we only need this as a placeholder value for `terraform destroy` to work,
  # this won't actually create a final snapshot on destroy
  final_snapshot_identifier = "foo"
}

resource "aws_db_instance" "nms" {
  identifier        = "nmsdb"
  allocated_storage = 16
  engine            = "mysql"
  engine_version    = "5.7"
  instance_class    = "db.m4.large"

  name     = "magma"
  username = "magma"
  password = var.nms_db_password

  vpc_security_group_ids = [aws_security_group.default.id]

  db_subnet_group_name = module.vpc.database_subnet_group

  skip_final_snapshot = true
  # we only need this as a placeholder value for `terraform destroy` to work,
  # this won't actually create a final snapshot on destroy
  final_snapshot_identifier = "nms-foo"
}
