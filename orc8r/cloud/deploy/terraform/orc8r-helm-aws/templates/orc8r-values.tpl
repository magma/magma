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

imagePullSecrets:
  - name: ${image_pull_secret}

secrets:
  create: false
secret:
  certs: ${certs_secret}
  configs:
    orc8r: ${configs_secret}
  envdir: ${envdir_secret}

# certs sub-chart configuration.
certs:
  create: ${managed_certs_create}
  enabled: ${managed_certs_enabled}
  domainName: ${managed_certs_domain_name}
  nms:
    customIssuer: ${nms_custom_issuer}
  route53:
    enabled: ${managed_certs_route53_enabled}
    region: ${region}

nginx:
  create: true

  podDisruptionBudget:
    enabled: true
  image:
    repository: ${docker_registry}/nginx
    tag: "${docker_tag}"
  replicas: ${nginx_replicas}
  service:
    enabled: true
    legacyEnabled: true
    annotations:
      service.beta.kubernetes.io/aws-load-balancer-additional-resource-tags: "magma-uuid=${magma_uuid}"
    extraAnnotations:
      proxy:
        external-dns.alpha.kubernetes.io/hostname: ${api_hostname}
      bootstrapLagacy:
        external-dns.alpha.kubernetes.io/hostname: bootstrapper-${controller_hostname}
      clientcertLegacy:
        external-dns.alpha.kubernetes.io/hostname: ${controller_hostname}
    name: orc8r-bootstrap-nginx
    type: LoadBalancer
  spec:
    hostname: ${controller_hostname}

controller:
  podDisruptionBudget:
    enabled: true
  image:
    repository: ${docker_registry}/controller
    tag: "${docker_tag}"
    env:
      orc8r_domain_name: "${orc8r_domain_name}"
      version_tag: "${docker_tag}"
      helm_version_tag: "${orc8r_chart_version}"
  replicas: ${controller_replicas}
  spec:
    database:
      db: ${orc8r_db_name}
      host: ${orc8r_db_host}
      port: ${orc8r_db_port}
      user: ${orc8r_db_user}
    service_registry:
      mode: "k8s"

metrics:
  imagePullSecrets:
    - name: ${image_pull_secret}
  metrics:
    volumes:
      prometheusData:
        volumeSpec:
          persistentVolumeClaim:
            claimName: ${metrics_pvc_promdata}
      prometheusConfig:
        volumeSpec:
          persistentVolumeClaim:
            claimName: ${metrics_pvc_promcfg}

  prometheus:
    create: ${enable_metrics}
    includeOrc8rAlerts: true
    prometheusCacheHostname: ${prometheus_cache_hostname}
    alertmanagerHostname: ${alertmanager_hostname}

  alertmanager:
    create: ${enable_metrics}

  prometheusConfigurer:
    create: ${enable_metrics}
    image:
      repository: docker.io/facebookincubator/prometheus-configurer
      tag: ${prometheus_configurer_version}
    prometheusURL: ${prometheus_url}

  alertmanagerConfigurer:
    create: ${enable_metrics}
    image:
      repository: docker.io/facebookincubator/alertmanager-configurer
      tag: ${alertmanager_configurer_version}
    alertmanagerURL: ${alertmanager_url}

  prometheusCache:
    create: ${enable_metrics}
    image:
      repository: docker.io/facebookincubator/prometheus-edge-hub
      tag: 1.1.0
    limit: 500000
  grafana:
    create: ${enable_metrics}

  userGrafana:
    image:
      repository: docker.io/grafana/grafana
      tag: 6.6.2
    create: ${create_usergrafana}
    volumes:
      datasources:
        volumeSpec:
          persistentVolumeClaim:
            claimName: ${grafana_pvc_grafanaDatasources}
      dashboardproviders:
        volumeSpec:
          persistentVolumeClaim:
            claimName: ${grafana_pvc_grafanaProviders}
      dashboards:
        volumeSpec:
          persistentVolumeClaim:
            claimName: ${grafana_pvc_grafanaDashboards}
      grafanaData:
        volumeSpec:
          persistentVolumeClaim:
            claimName: ${grafana_pvc_grafanaData}

  thanos:
    enabled: ${thanos_enabled}

    compact:
      nodeSelector:
        ${thanos_compact_selector}

    store:
      nodeSelector:
        ${thanos_store_selector}

    query:
      nodeSelector:
        ${thanos_query_selector}

    objstore:
      type: S3
      config:
        bucket: ${thanos_bucket}
        endpoint: s3.${region}.amazonaws.com
        region: ${region}
        access_key: ${thanos_aws_access_key}
        secret_key: ${thanos_aws_secret_key}
        insecure: false
        signature_version2: false
        put_user_metadata: {}
        http_config:
          idle_conn_timeout: 0s
          response_header_timeout: 0s
          insecure_skip_verify: false
        trace:
          enable: false
        part_size: 0

nms:
  enabled: ${deploy_nms}

  imagePullSecrets:
    - name: ${image_pull_secret}

  secret:
    certs: ${nms_certs_secret}

  certs:
    enabled: ${nms_managed_certs_enabled}

  magmalte:
    create: true

    image:
      repository: ${docker_registry}/magmalte
      tag: "${docker_tag}"

    env:
      api_host: ${api_hostname}
      mysql_db: ${orc8r_db_name}
      mysql_dialect: ${orc8r_db_dialect}
      mysql_host: ${orc8r_db_host}
      mysql_port: ${orc8r_db_port}
      mysql_user: ${orc8r_db_user}
      mysql_pass: ${orc8r_db_pass}
      grafana_address: ${user_grafana_hostname}
      version_tag: "${docker_tag}"

  nginx:
    create: true

    service:
      type: LoadBalancer
      annotations:
        external-dns.alpha.kubernetes.io/hostname: "${nms_hostname}"
        service.beta.kubernetes.io/aws-load-balancer-additional-resource-tags: "magma-uuid=${magma_uuid}"

    deployment:
      spec:
        ssl_cert_name: controller.crt
        ssl_cert_key_name: controller.key

logging:
  enabled: ${enable_logging}

dp:
  create: ${dp_enabled}

  configuration_controller:
    sasEndpointUrl: "${dp_sas_endpoint_url}"
    image:
      repository: "${docker_registry}/configuration-controller"
      tag: "${docker_tag}"

    database:
      driver: postgres
      db: ${orc8r_db_name}
      host: ${orc8r_db_host}
      port: ${orc8r_db_port}
      user: ${orc8r_db_user}
      pass: ${orc8r_db_pass}

  protocol_controller:
    enabled: false
    image:
      repository: "${docker_registry}/protocol-controller"
      tag: "${docker_tag}"

  radio_controller:
    image:
      repository: "${docker_registry}/radio-controller"
      tag: "${docker_tag}"

    database:
      driver: postgres
      db: ${orc8r_db_name}
      host: ${orc8r_db_host}
      port: ${orc8r_db_port}
      user: ${orc8r_db_user}
      pass: ${orc8r_db_pass}

  active_mode_controller:
    image:
      repository: "${docker_registry}/active-mode-controller"
      tag: "${docker_tag}"

  db_service:
    image:
      repository: "${docker_registry}/db-service"
      tag: "${docker_tag}"

    database:
      driver: postgres
      db: ${orc8r_db_name}
      host: ${orc8r_db_host}
      port: ${orc8r_db_port}
      user: ${orc8r_db_user}
      pass: ${orc8r_db_pass}
