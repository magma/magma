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
resource "helm_release" "efs_provisioner" {
  count = var.orc8r_is_staging_deployment == true ? 0 : 1

  name       = var.efs_provisioner_name
  repository = local.stable_helm_repo
  chart      = "efs-provisioner"
  version    = "0.11.0"
  namespace  = "kube-system"
  keyring    = ""

  values = [<<VALUES
  efsProvisioner:
    efsFileSystemId: ${var.efs_file_system_id}
    awsRegion: ${var.region}
    path: /pv-volume
    provisionerName: aws-efs
    storageClass:
      name: ${var.efs_storage_class_name}
  podAnnotations:
    iam-assumable-role: ${var.efs_provisioner_role_arn}
  VALUES
  ]
}
