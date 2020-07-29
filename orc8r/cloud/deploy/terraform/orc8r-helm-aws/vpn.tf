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

resource "helm_release" "ovpn" {
  count = var.deploy_openvpn ? 1 : 0

  chart      = "openvpn"
  name       = "openvpn"
  namespace  = kubernetes_namespace.orc8r.metadata[0].name
  repository = local.stable_helm_repo

  # TCP ovpn because ELB does not support UDP
  values = [<<EOT
  openvpn:
    OVPN_K8S_POD_NETWORK: null
    OVPN_K8S_POD_SUBNET: null
    OVPN_PROTO: tcp
    redirectGateway: false
    serverConf: |-
      client-to-client
      duplicate-cn
  service:
    annotations:
      external-dns.alpha.kubernetes.io/hostname: vpn.${var.orc8r_domain_name}
  persistence:
    existingClaim: ${kubernetes_persistent_volume_claim.storage["openvpn"].metadata.0.name}
  EOT
  ]

  # Cert creation can take some time
  timeout = 900
}
