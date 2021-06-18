################################################################################
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
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

  backup_retention_period = var.orc8r_db_backup_retention
  backup_window           = var.orc8r_db_backup_window

  allow_major_version_upgrade = true
  skip_final_snapshot = true
  # we only need this as a placeholder value for `terraform destroy` to work,
  # this won't actually create a final snapshot on destroy
  final_snapshot_identifier = "foo"
}

resource "aws_sns_topic" "sns_orc8r_topic" {
  name = var.orc8r_sns_name
}

resource "aws_sns_topic_subscription" "sns_orc8r_db_subscription_email" {
  count     = var.enable_aws_db_notifications ? 1: 0
  topic_arn = aws_sns_topic.sns_orc8r_topic.arn
  protocol  = "email"
  endpoint  = var.orc8r_sns_email
}

resource "aws_db_event_subscription" "default" {
  name      = var.orc8r_db_event_subscription
  sns_topic = aws_sns_topic.sns_orc8r_topic.arn
  source_type = "db-instance"
  source_ids = [aws_db_instance.default.id]
  event_categories = ["failure", "maintenance", "notification", "restoration"]
}
