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

# helm tiller service account
resource "kubernetes_service_account" "tiller" {
  count = var.install_tiller == true ? 1 : 0

  metadata {
    name      = "tiller"
    namespace = var.tiller_namespace
  }

  automount_service_account_token = true
}

# helm tiller cluster role
resource "kubernetes_cluster_role_binding" "tiller" {
  count = var.install_tiller == true ? 1 : 0

  metadata {
    name = "tiller"
  }

  role_ref {
    kind      = "ClusterRole"
    name      = "cluster-admin"
    api_group = "rbac.authorization.k8s.io"
  }

  subject {
    kind      = "ServiceAccount"
    name      = "tiller"
    api_group = ""
    namespace = var.tiller_namespace
  }
}
