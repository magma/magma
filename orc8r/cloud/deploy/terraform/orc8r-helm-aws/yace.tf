################################################################################
# Copyright 2021 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
################################################################################

resource "aws_iam_user" "yace_user" {
  count = var.cloudwatch_exporter_enabled ? 1 : 0

  name = "yace_user"
}

resource "aws_iam_access_key" "yace_access_key" {
  count = var.cloudwatch_exporter_enabled ? 1 : 0

  user = aws_iam_user.yace_user[0].name
}

data "aws_iam_policy_document" "yace_policy_doc" {
  count = var.cloudwatch_exporter_enabled ? 1 : 0

  statement {
    effect = "Allow"

    actions = [
      "tag:GetResources",
      "cloudwatch:ListTagsForResource",
      "cloudwatch:GetMetricStatistics",
      "cloudwatch:GetMetricData",
      "cloudwatch:ListMetrics"
    ]

    resources = [ "*" ]
  }
}

resource "aws_iam_user_policy" "yace_policy" {
  count = var.cloudwatch_exporter_enabled ? 1 : 0

  name = "yace_policy"
  user = aws_iam_user.yace_user[0].name

  policy = data.aws_iam_policy_document.yace_policy_doc[0].json
}

resource "helm_release" "yace_exporter" {
  count = var.cloudwatch_exporter_enabled ? 1 : 0

  name       = "yace-exporter"
  namespace  = kubernetes_namespace.monitoring.metadata[0].name
  repository = "https://mogaal.github.io/helm-charts/"
  chart      = "prometheus-yace-exporter"
  timeout    = 600
  values = [<<EOT
# Default values for prometheus-yace-exporter.
# This is a YAML-formatted file.
# Declare variables to be passed into your templates.

replicaCount: 1

image:
  repository: quay.io/invisionag/yet-another-cloudwatch-exporter
  tag: v0.25.0-alpha
  pullPolicy: IfNotPresent

nameOverride: ""
fullnameOverride: ""

service:
  type: ClusterIP
  port: 80
  annotations:
    prometheus.io/scrape: "true"
    prometheus.io/port: "5000"
  labels: {}

ingress:
  enabled: false
  annotations: {}
    # kubernetes.io/ingress.class: nginx
    # kubernetes.io/tls-acme: "true"
  hosts:
    - host: chart-example.local
      paths: []

  tls: []
  #  - secretName: chart-example-tls
  #    hosts:
  #      - chart-example.local

resources: {}
  # We usually recommend not to specify default resources and to leave this as a conscious
  # choice for the user. This also increases chances charts run on environments with little
  # resources, such as Minikube. If you do want to specify resources, uncomment the following
  # lines, adjust them as necessary, and remove the curly braces after 'resources:'.
  # limits:
  #   cpu: 100m
  #   memory: 128Mi
  # requests:
  #   cpu: 100m
  #   memory: 128Mi

nodeSelector: {}

tolerations: []

affinity: {}

podAnnoatations: {}

podLabels: {}

aws:
  role:

  # The name of a pre-created secret in which AWS credentials are stored. When
  # set, aws_access_key_id is assumed to be in a field called access_key,
  # aws_secret_access_key is assumed to be in a field called secret_key, and the
  # session token, if it exists, is assumed to be in a field called
  # security_token
  secret:
    name:
    includesSessionToken: false

  # Note: Do not specify the aws_access_key_id and aws_secret_access_key if you specified role or secret.name before
  aws_access_key_id:
  aws_secret_access_key:

serviceAccount:
  # Specifies whether a ServiceAccount should be created
  create: true
  annotations: {}
  labels: {}
  # The name of the ServiceAccount to use.
  # If not set and create is true, a name is generated using the fullname template
  name:

rbac:
  # Specifies whether RBAC resources should be created
  create: true

serviceMonitor:
  # When set true then use a ServiceMonitor to configure scraping
  enabled: true
  # Set the namespace the ServiceMonitor should be deployed
  namespace: ${var.orc8r_kubernetes_namespace}
  # Set how frequently Prometheus should scrape
  interval: 30s
  # Set targetPort for serviceMonitor
  port: http
  # Set path to cloudwatch-exporter telemtery-path
  telemetryPath: /metrics
  # Set labels for the ServiceMonitor, use this to define your scrape label for Prometheus Operator
  # labels:
  # Set timeout for scrape
  timeout: 10s


config: |-
  discovery:
    exportedTagsOnMetrics:
      ec2:
        - Name
      ebs:
        - VolumeId
    jobs:
    - regions:
        - ${var.region}
      type: "es"
      searchTags:
        - Key: magma-uuid
          Value: ${var.magma_uuid}
      metrics:
        - name: FreeStorageSpace
          statistics:
          - 'Sum'
          period: 600
          length: 60
        - name: ClusterStatus.green
          statistics:
          - 'Minimum'
          period: 600
          length: 60
        - name: ClusterStatus.yellow
          statistics:
          - 'Maximum'
          period: 600
          length: 60
        - name: ClusterStatus.red
          statistics:
          - 'Maximum'
          period: 600
          length: 60
    - type: "elb"
      regions:
        - ${var.region}
      searchTags:
        - Key: magma-uuid
          Value: ${var.magma_uuid}
      metrics:
        - name: HealthyHostCount
          statistics:
          - 'Minimum'
          period: 600
          length: 600
        - name: HTTPCode_Backend_4XX
          statistics:
          - 'Sum'
          period: 60
          length: 900
          delay: 300
          nilToZero: true
    - type: "ec2"
      regions:
        - ${var.region}
      searchTags:
        - Key: magma-uuid
          Value: ${var.magma_uuid}
      metrics:
        - name: NetworkIn
          statistics:
          - 'Sum'
          period: 600
          length: 600
        - name: NetworkOut
          statistics:
          - 'Sum'
          period: 60
          length: 900
          delay: 300
          nilToZero: true
  EOT
  ]
  set_sensitive {
    name  = "aws.aws_access_key_id"
    value = var.cloudwatch_exporter_enabled ? aws_iam_access_key.yace_access_key[0].id : ""
  }
  set_sensitive {
    name  = "aws.aws_secret_access_key"
    value = var.cloudwatch_exporter_enabled ? aws_iam_access_key.yace_access_key[0].secret : ""
  }
}