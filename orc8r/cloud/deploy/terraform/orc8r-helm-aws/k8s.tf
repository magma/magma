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

resource "kubernetes_namespace" "orc8r" {
  metadata {
    name = var.orc8r_kubernetes_namespace
  }
}

resource "kubernetes_namespace" "monitoring" {
  count = var.orc8r_is_staging_deployment == true ? 0 : 1

  metadata {
    name = var.monitoring_kubernetes_namespace
  }
}

# external dns maps route53 to ingress resources
resource "helm_release" "external_dns" {
  count = var.orc8r_is_staging_deployment == true ? 0 : 1

  name       = var.external_dns_deployment_name
  repository = local.stable_helm_repo
  chart      = "external-dns"
  version    = "2.19.1"
  namespace  = "kube-system"
  keyring    = ""

  values = [<<VALUES
  rbac:
    create: true
  aws:
    assumeRoleArn: ${var.external_dns_role_arn}
  zoneIdFilters:
    - ${var.orc8r_route53_zone_id}
  VALUES
  ]
}

resource "helm_release" "cert-manager" {
  count = var.deploy_cert_manager_helm_chart ? 1 : 0

  version          = "1.6.1"
  name             = "cert-manager"
  chart            = "cert-manager"
  namespace        = "cert-manager"
  repository       = "https://charts.jetstack.io"
  create_namespace = true

  set {
    name  = "installCRDs"
    value = true
  }
  set {
    name  = "serviceAccount.annotations.eks\\.amazonaws\\.com/role-arn"
    value = var.cert_manager_route53_iam_role_arn
  }
  set {
    name  = "securityContext.fsGroup"
    value = 1001
  }
  set {
    name  = "extraArgs"
    value = "{--issuer-ambient-credentials}"
  }
}
