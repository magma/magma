esJavaOpts: -Xmx2G -Xms2G
imageTag: 7.12.0
minimumMasterNodes: 1
replicas: 1
resources:
  limits:
    cpu: "2"
    memory: 4Gi
  requests:
    cpu: "2"
    memory: 4Gi
roles:
  data: "false"
  ingest: "false"
  master: "true"
  ml: "false"
  remote_cluster_client: "false"
volumeClaimTemplate:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 1Gi
  storageClassName: nfs
replicas: 1
minimumMasterNodes: 1
rbac:
  create: true
antiAffinity: "soft"
service:
  annotations:
    external-dns.alpha.kubernetes.io/hostname: elasticsearch.${dns_domain}
