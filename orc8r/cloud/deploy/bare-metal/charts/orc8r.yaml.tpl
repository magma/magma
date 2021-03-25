controller:
  image:
    repository: ${img_repo}/controller
    tag: ${controller_tag}
  podDisruptionBudget:
    enabled: true
  replicas: 2
  spec:
    database:
      driver: mysql
      sql_dialect: maria
      db: ${orc8r_db_user}
      host: ${orc8r_db_host}
      pass: ${orc8r_db_pass}
      port: 3306
      user: ${orc8r_db_user}
  service:
    portEnd: 9121

imagePullSecrets:
- name: artifactory
logging:
  enabled: false
metrics:
  alertmanager:
    create: true
  alertmanagerConfigurer:
    alertmanagerURL: orc8r-alertmanager:9093
    create: true
    image:
      repository: docker.io/facebookincubator/alertmanager-configurer
      tag: 1.0.0
  grafana:
    create: false
  imagePullSecrets:
  - name: artifactory
  metrics:
    volumes:
      prometheusConfig:
        volumeSpec:
          persistentVolumeClaim:
            claimName: promcfg
      prometheusData:
        volumeSpec:
          persistentVolumeClaim:
            claimName: promdata
  prometheus:
    create: true
    includeOrc8rAlerts: true
  prometheusCache:
    create: true
    image:
      repository: docker.io/facebookincubator/prometheus-edge-hub
      tag: 1.0.0
    limit: 500000
  prometheusConfigurer:
    create: true
    image:
      repository: docker.io/facebookincubator/prometheus-configurer
      tag: 1.0.0
    prometheusURL: orc8r-prometheus:9090
  userGrafana:
    create: true
    image:
      repository: docker.io/grafana/grafana
      tag: 6.6.2
    volumes:
      dashboardproviders:
        persistentVolumeClaim:
          claimName: grafanaproviders
      dashboards:
        persistentVolumeClaim:
          claimName: grafanadashboards
      datasources:
        persistentVolumeClaim:
          claimName: grafanadatasources
      grafanaData:
        persistentVolumeClaim:
          claimName: grafanadata

nms:
  enabled: true
  imagePullSecrets:
  - name: artifactory
  magmalte:
    env:
      api_host: api.${dns_domain}
      grafana_address: orc8r-user-grafana:3000
      mysql_db: ${nms_db_user}
      mysql_host: ${nms_db_host}
      mysql_pass: ${nms_db_pass}
      mysql_port: 3306
      mysql_user: ${nms_db_user}
      mysql_dialect: mariadb
    image:
      pullPolicy: IfNotPresent
      repository: ${img_repo}/magmalte
      tag: ${nms_tag}
    manifests:
      deployment: true
      rbac: false
      secrets: true
      service: true
  nginx:
    deployment:
      spec:
        ssl_cert_key_name: controller.key
        ssl_cert_name: controller.crt
    manifests:
      configmap: true
      deployment: true
      rbac: false
      secrets: true
      service: true
    service:
      annotations:
        external-dns.alpha.kubernetes.io/hostname: '*.nms.${dns_domain}'
      type: LoadBalancer
  secret:
    certs: nms-certs
nginx:
  image:
    repository: ${img_repo}/nginx
    tag: ${nginx_tag}
  podDisruptionBudget:
    enabled: true
  replicas: 2
  service:
    enabled: true
    extraAnnotations:
      bootstrapLagacy:
        external-dns.alpha.kubernetes.io/hostname: bootstrapper-controller.${dns_domain}
      clientcertLegacy:
        external-dns.alpha.kubernetes.io/hostname: controller.${dns_domain}
      proxy:
        external-dns.alpha.kubernetes.io/hostname: api.${dns_domain}
    legacyEnabled: true
    name: orc8r-bootstrap-legacy
    type: LoadBalancer
  spec:
    hostname: controller.${dns_domain}

secret:
  certs: orc8r-certs
  configs:
    orc8r: orc8r-configs
  envdir: orc8r-envdir
secrets:
  create: false
