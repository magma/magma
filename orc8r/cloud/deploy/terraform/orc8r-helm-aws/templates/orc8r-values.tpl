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
    extraAnnotations:
      bootstrapLagacy:
        external-dns.alpha.kubernetes.io/hostname: bootstrapper-${controller_hostname}
      clientcertLegacy:
        external-dns.alpha.kubernetes.io/hostname: ${controller_hostname},${api_hostname}
    name: orc8r-bootstrap-legacy
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
  alertmanager:
    create: true
  prometheusConfigurer:
    create: true
    image:
      repository: ${docker_registry}/prometheus-configurer
      tag: "${docker_tag}"
  alertmanagerConfigurer:
    create: true
    image:
      repository: ${docker_registry}/alertmanager-configurer
      tag: "${docker_tag}"
  prometheusCache:
    create: true
    image:
      repository: ${docker_registry}/prometheus-cache
      tag: "${docker_tag}"
    limit: 500000
  grafana:
    create: true
    image:
      repository: ${docker_registry}/grafana
      tag: "${docker_tag}"

  userGrafana:
    image:
      repository: docker.io/grafana/grafana
      tag: 6.6.2
    create: ${create_usergrafana}
    volumes:
      datasources:
        persistentVolumeClaim:
          claimName: ${grafana_pvc_grafanaData}
      dashboardproviders:
        persistentVolumeClaim:
          claimName: ${grafana_pvc_grafanaData}
      dashboards:
        persistentVolumeClaim:
          claimName: ${grafana_pvc_grafanaData}
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
    create: true

    image:
      repository: ${docker_registry}/magmalte
      tag: "${docker_tag}"

    env:
      api_host: ${api_hostname}
      mysql_host: ${nms_db_host}
      mysql_user: ${nms_db_user}
  nginx:
    create: true

    service:
      type: LoadBalancer
      annotations:
        external-dns.alpha.kubernetes.io/hostname: ${nms_hostname}

    deployment:
      spec:
        ssl_cert_name: controller.crt
        ssl_cert_key_name: controller.key

logging:
  enabled: false
