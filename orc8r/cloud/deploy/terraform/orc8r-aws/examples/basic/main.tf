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

module orc8r {
  source = "../.."

  region = "us-west-2"

  orc8r_db_password           = "Faceb00k12345"
  secretsmanager_orc8r_secret = "magma-orc8r-test"
  deployment_secrets_bucket   = "magma.orc8r.test"
  orc8r_domain_name           = "orc8r.magma.test"

  vpc_name        = "magma-orc8r-test"
  cluster_name    = "orc8r-test"
  cluster_version = "1.17"
}
