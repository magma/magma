imagePullSecrets:
  - name: ${image_pull_secret}

secrets:
  create: false
secret:
  certs: ${certs_secret}
  configs:
    orc8r: ${configs_secret}
  envdir: ${envdir_secret}

proxy:
  podDisruptionBudget:
    enabled: true
  image:
    repository: ${docker_registry}/proxy
    tag: "${docker_tag}"
  replicas: ${proxy_replicas}
  service:
    enabled: true
    legacyEnabled: true
    %{~ if create_nginx ~}
    %{~ else ~}
    extraAnnotations:
      bootstrapLagacy:
        external-dns.alpha.kubernetes.io/hostname: bootstrapper-${controller_hostname}
      clientcertLegacy:
        external-dns.alpha.kubernetes.io/hostname: ${controller_hostname},${api_hostname}
    %{~ endif ~}
    name: orc8r-bootstrap-legacy
    type: LoadBalancer
  spec:
    hostname: ${controller_hostname}

nginx:
  create: ${create_nginx}

  podDisruptionBudget:
    enabled: true
  image:
    repository: ${docker_registry}/nginx
    tag: "${docker_tag}"
  replicas: ${nginx_replicas}
  service:
    enabled: true
    legacyEnabled: true
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
  replicas: ${controller_replicas}
  spec:
    database:
      db: ${orc8r_db_name}
      host: ${orc8r_db_host}
      port: ${orc8r_db_port}
      user: ${orc8r_db_user}

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
    create: true
    includeOrc8rAlerts: true
    prometheusCacheHostname: ${prometheus_cache_hostname}
    alertmanagerHostname: ${alertmanager_hostname}

  alertmanager:
    create: true

  prometheusConfigurer:
    create: true
    image:
      repository: ${docker_registry}/prometheus-configurer
      tag: "${docker_tag}"
    prometheusURL: ${prometheus_url}

  alertmanagerConfigurer:
    create: true
    image:
      repository: ${docker_registry}/alertmanager-configurer
      tag: "${docker_tag}"
    alertmanagerURL: ${alertmanager_url}

  prometheusCache:
    create: true
    image:
      repository: docker.io/facebookincubator/prometheus-edge-hub
      tag: 1.0.0
    limit: 500000
  grafana:
    create: false

  userGrafana:
    image:
      repository: docker.io/grafana/grafana
      tag: 6.6.2
    create: ${create_usergrafana}
    volumes:
      datasources:
        persistentVolumeClaim:
          claimName: ${grafana_pvc_grafanaDatasources}
      dashboardproviders:
        persistentVolumeClaim:
          claimName: ${grafana_pvc_grafanaProviders}
      dashboards:
        persistentVolumeClaim:
          claimName: ${grafana_pvc_grafanaDashboards}
      grafanaData:
        persistentVolumeClaim:
          claimName: ${grafana_pvc_grafanaData}

nms:
  enabled: ${deploy_nms}

  imagePullSecrets:
    - name: ${image_pull_secret}

  secret:
    certs: ${nms_certs_secret}

  magmalte:
    manifests:
      secrets: true
      deployment: true
      service: true
      rbac: false

    image:
      repository: ${docker_registry}/magmalte
      tag: "${docker_tag}"

    env:
      api_host: ${controller_hostname}
      mysql_host: ${nms_db_host}
      mysql_user: ${nms_db_user}
      grafana_address: ${user_grafana_hostname}

  nginx:
    manifests:
      configmap: true
      secrets: true
      deployment: true
      service: true
      rbac: false

    service:
      type: LoadBalancer
      annotations:
        external-dns.alpha.kubernetes.io/hostname: "${nms_hostname}"

    deployment:
      spec:
        ssl_cert_name: controller.crt
        ssl_cert_key_name: controller.key

logging:
  enabled: false
