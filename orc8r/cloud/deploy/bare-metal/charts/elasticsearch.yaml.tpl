clusterName: "es-magma"
replicas: 1
minimumMasterNodes: 1
rbac:
  create: true
antiAffinity: "soft"
service:
  annotations:
    external-dns.alpha.kubernetes.io/hostname: elasticsearch.${dns_domain}
data:
  storageClass: ${storage_class}
master:
  storageClass: ${storage_class}
