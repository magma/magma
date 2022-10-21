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

# k8s requires provisioner to treat efs as a persistent volume
resource "helm_release" "kubernetes_efs_csi_driver" {
  count = var.orc8r_is_staging_deployment == true ? 0 : 1

  name       = "aws-efs-csi-driver"
  repository = "https://kubernetes-sigs.github.io/aws-efs-csi-driver"
  chart      = "aws-efs-csi-driver"
  version    = "2.2.9"
  namespace  = "kube-system"

  values = [<<VALUES
  controller:
    serviceAccount:
      name: "aws-efs-csi-driver"
      annotations:
        eks.amazonaws.com/role-arn: ${var.efs_csi_driver_arn}
  node:
    serviceAccount:
      name: "aws-efs-csi-driver"
      create: false
  storageClasses:
    - name: ${var.efs_storage_class_name}
      parameters:
        provisioningMode: efs-ap
        fileSystemId: ${var.efs_file_system_id}
        directoryPerms: "700"
  VALUES
  ]
}
