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

  skip_final_snapshot = true
  # we only need this as a placeholder value for `terraform destroy` to work,
  # this won't actually create a final snapshot on destroy
  final_snapshot_identifier = "foo"
}
