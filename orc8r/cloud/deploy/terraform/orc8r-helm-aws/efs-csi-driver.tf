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
resource "helm_release" "aws_efs_csi_driver" {
  name             = "aws-efs-csi-driver"
  chart            = "aws-efs-csi-driver"
  repository       = "https://kubernetes-sigs.github.io/aws-efs-csi-driver"
  version          = "2.2.9"
  namespace        = var.efs_namespace
  create_namespace = true

  values = [<<VALUES
  controller:
    serviceAccount:
      create: true
      name: ${var.efs_service_account}
      annotations:
        eks.amazonaws.com/role-arn: ${aws_iam_role.efs_csi_driver.arn}
  node:
    serviceAccount:
      create: false
      name: ${var.efs_service_account}
      annotations:
        eks.amazonaws.com/role-arn: ${aws_iam_role.efs_csi_driver.arn}
  storageClasses:
    - name: efs-sc
      parameters:
        provisioningMode: efs-ap
        fileSystemId: ${aws_efs_file_system.eks.id}
        directoryPerms: "700"
  VALUES
  ]
}
