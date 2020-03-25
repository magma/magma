################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

resource "aws_db_instance" "default" {
  identifier        = var.orc8r_db_identifier
  allocated_storage = var.orc8r_db_storage_gb
  engine            = "postgres"
  engine_version    = var.orc8r_db_engine_version
  instance_class    = var.orc8r_db_instance_class

  name     = var.orc8r_db_name
  username = var.orc8r_db_username
  password = var.orc8r_db_password

  vpc_security_group_ids = [aws_security_group.default.id]

  db_subnet_group_name = module.vpc.database_subnet_group

  skip_final_snapshot = true
  # we only need this as a placeholder value for `terraform destroy` to work,
  # this won't actually create a final snapshot on destroy
  final_snapshot_identifier = "foo"
}

resource "aws_db_instance" "nms" {
  identifier        = var.nms_db_identifier
  allocated_storage = var.nms_db_storage_gb
  engine            = "mysql"
  engine_version    = var.nms_db_engine_version
  instance_class    = var.nms_db_instance_class

  name     = var.nms_db_name
  username = var.nms_db_username
  password = var.nms_db_password

  vpc_security_group_ids = [aws_security_group.default.id]

  db_subnet_group_name = module.vpc.database_subnet_group

  skip_final_snapshot = true
  # we only need this as a placeholder value for `terraform destroy` to work,
  # this won't actually create a final snapshot on destroy
  final_snapshot_identifier = "nms-foo"
}
