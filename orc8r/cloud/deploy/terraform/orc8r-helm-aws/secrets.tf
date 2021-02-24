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

# null resource with local-exec provisioner to seed Secretsmanager with all
# relevant orc8r secrets
# We use a null resource with a local-exec provisioner instead of an external
# data source because this script should only be run on creation and on-demand
# via terraform taint. This is the behavior of the local-exec provisioner - see
# https://www.terraform.io/docs/provisioners/#creation-time-provisioners
# This does introduce an unfortunate side effect that users have to first
# target this resource with terraform apply before doing a full tf apply.
# This is because the data property of the k8s secret resources below cannot
# be conditionally evaluated based on the existence of this null resource and
# so will error out during the plan if secretsmanager hasn't been seeded.
resource "null_resource" orc8r_seed_secrets {
  provisioner "local-exec" {
    command = <<EOT
      ${path.module}/scripts/create_orc8r_secrets.py \
        '${var.secretsmanager_orc8r_name}' '${var.region}' \
        "${var.seed_certs_dir}"
    EOT
  }
}

locals {
  orc8r_cert_names = [
    "rootCA.key",
    "rootCA.pem",
    "controller.key",
    "controller.crt",
    "certifier.key",
    "certifier.pem",
    "bootstrapper.key",
    "admin_operator.pem",
  ]

  fluentd_cert_names = [
    "fluentd.key",
    "fluentd.pem",
    "certifier.pem",
  ]

  nms_cert_names = [
    "admin_operator.pem",
    "admin_operator.key.pem",
    "controller.crt",
    "controller.key",
  ]
}

resource "kubernetes_secret" "orc8r_certs" {
  metadata {
    name      = "orc8r-certs"
    namespace = kubernetes_namespace.orc8r.metadata[0].name
  }

  data = {
    for name in local.orc8r_cert_names : name => jsondecode(data.aws_secretsmanager_secret_version.orc8r_secrets.secret_string)[name]
  }

  depends_on = [null_resource.orc8r_seed_secrets]
}

resource "kubernetes_secret" "nms_certs" {
  count = var.deploy_nms ? 1 : 0

  metadata {
    name      = "nms-certs"
    namespace = kubernetes_namespace.orc8r.metadata[0].name
  }

  data = {
    for name in local.nms_cert_names : name => jsondecode(data.aws_secretsmanager_secret_version.orc8r_secrets.secret_string)[name]
  }
}

resource "kubernetes_secret" "orc8r_configs" {
  metadata {
    name      = "orc8r-configs"
    namespace = kubernetes_namespace.orc8r.metadata[0].name
  }

  data = {
    "metricsd.yml" = yamlencode({
      "profile" : "prometheus",
      "prometheusQueryAddress" : var.thanos_enabled ? format("http://%s-thanos-query-http:10902", var.helm_deployment_name) : format("http://%s-prometheus:9090", var.helm_deployment_name),

      "alertmanagerApiURL" : format("http://%s-alertmanager:9093/api/v2", var.helm_deployment_name),
      "prometheusConfigServiceURL" : format("http://%s-prometheus-configurer:9100/v1", var.helm_deployment_name),
      "alertmanagerConfigServiceURL" : format("http://%s-alertmanager-configurer:9101/v1", var.helm_deployment_name),
    })

    "orchestrator.yml" = yamlencode({
      "useGRPCExporter": true,
      "prometheusGRPCPushAddress" : format("%s-prometheus-cache:9092", var.helm_deployment_name),
      "prometheusPushAddresses" : [
        format("http://%s-prometheus-cache:9091/metrics", var.helm_deployment_name),
      ],
    })

    "elastic.yml" = yamlencode({
      "elasticHost" : var.elasticsearch_endpoint == null ? "elastic" : var.elasticsearch_endpoint
      "elasticPort" : 80,
    })

    "analytics.yml" = yamlencode({
      "exportMetrics": var.analytics_export_enabled == null ? false : var.analytics_export_enabled,
      "metricsPrefix": var.analytics_metrics_prefix == null ? "" : var.analytics_metrics_prefix,
      "appSecret": var.analytics_app_secret == null ? "" : var.analytics_app_secret,
      "appID": var.analytics_app_id == null ? "" : var.analytics_app_id,
      "metricExportURL": var.analytics_metric_export_url == null ? "" : var.analytics_metric_export_url,
      "categoryName": var.analytics_category_name == null ? "" : var.analytics_category_name,
    })
  }
}

resource "kubernetes_secret" "orc8r_envdir" {
  metadata {
    name      = "orc8r-envdir"
    namespace = kubernetes_namespace.orc8r.metadata[0].name
  }
}

resource "kubernetes_secret" "fluentd_certs" {
  metadata {
    name      = "fluentd-certs"
    namespace = kubernetes_namespace.orc8r.metadata[0].name
  }

  data = {
    for name in local.fluentd_cert_names : name => jsondecode(data.aws_secretsmanager_secret_version.orc8r_secrets.secret_string)[name]
  }

  depends_on = [null_resource.orc8r_seed_secrets]
}

data "aws_secretsmanager_secret" "orc8r_secrets" {
  name = var.secretsmanager_orc8r_name
}

data "aws_secretsmanager_secret_version" "orc8r_secrets" {
  secret_id = data.aws_secretsmanager_secret.orc8r_secrets.id
}
