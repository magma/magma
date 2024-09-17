esJavaOpts: -Xmx4G -Xms4G
minimumMasterNodes: 1
nodeGroup: data
replicas: 1
resources:
  limits:
    cpu: "1"
    memory: 8Gi
  requests:
    cpu: "1"
    memory: 8Gi
roles:
  data: "true"
  ingest: "true"
  master: "false"
  ml: "false"
  remote_cluster_client: "false"
volumeClaimTemplate:
  accessModes:
  - ReadWriteOnce
  resources:
    requests:
      storage: 50Gi
  storageClassName: nfs
