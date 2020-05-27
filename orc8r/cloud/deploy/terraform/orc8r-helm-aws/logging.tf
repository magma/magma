################################################################################
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
################################################################################

locals {
  fluentd_hostname = format("fluentd.%s", var.orc8r_domain_name)
}

resource "helm_release" "fluentd" {
  count = var.elasticsearch_endpoint == null ? 0 : 1

  name       = "fluentd"
  namespace  = kubernetes_namespace.orc8r.metadata[0].name
  repository = data.helm_repository.stable.id
  chart      = "fluentd"
  version    = "2.3.2"
  keyring    = ""

  values = [<<EOT
  replicaCount: 2
  output:
    host: ${var.elasticsearch_endpoint}
    port: 443
    scheme: https
  rbac:
    create: false
  service:
    annotations:
      external-dns.alpha.kubernetes.io/hostname: ${local.fluentd_hostname}
    type: LoadBalancer
    ports:
      - name: "forward"
        protocol: TCP
        containerPort: 24224
  configMaps:
    forward-input.conf: |-
      <source>
        @type forward
        port 24224
        bind 0.0.0.0
        <transport tls>
          ca_path /certs/certifier.pem
          cert_path /certs/fluentd.pem
          private_key_path /certs/fluentd.key
          client_cert_auth true
        </transport>
      </source>
    output.conf: |-
      <match **>
        @id elasticsearch
        @type elasticsearch
        @log_level info
        include_tag_key true
        host "#{ENV['OUTPUT_HOST']}"
        port "#{ENV['OUTPUT_PORT']}"
        scheme "#{ENV['OUTPUT_SCHEME']}"
        ssl_version "#{ENV['OUTPUT_SSL_VERSION']}"
        logstash_format true
        logstash_prefix "magma"
        reconnect_on_error true
        reload_on_failure true
        reload_connections false
        <buffer>
          @type file
          path /var/log/fluentd-buffers/kubernetes.system.buffer
          flush_mode interval
          retry_type exponential_backoff
          flush_thread_count 2
          flush_interval 5s
          retry_forever
          retry_max_interval 30
          chunk_limit_size "#{ENV['OUTPUT_BUFFER_CHUNK_LIMIT']}"
          queue_limit_length "#{ENV['OUTPUT_BUFFER_QUEUE_LIMIT']}"
          overflow_action block
        </buffer>
      </match>
  extraVolumes:
    - name: certs
      secret:
        defaultMode: 420
        secretName: ${kubernetes_secret.fluentd_certs.metadata.0.name}
  extraVolumeMounts:
    - name: certs
      mountPath: /certs
      readOnly: true
  EOT
  ]
}

# helm chart for cleanning old indices.
resource "helm_release" "elasticsearch_curator" {
  count = var.elasticsearch_endpoint == null ? 0 : 1

  name       = "elasticsearch-curator"
  repository = data.helm_repository.stable.id
  chart      = "elasticsearch-curator"
  namespace  = "monitoring"
  version    = "2.1.3"
  keyring    = ""

  values = [<<EOT
  configMaps:
    config_yml: |-
      ---
      client:
        hosts:
          - "${var.elasticsearch_endpoint}"

    action_file_yml: |-
      ---
      actions:
        1:
          action: delete_indices
          description: "Clean up ES by deleting old indices"
          options:
            timeout_override:
            continue_if_exception: False
            disable_action: False
            ignore_empty_list: True
          filters:
          - filtertype: age
            source: name
            direction: older
            timestring: '%Y.%m.%d'
            unit: days
            unit_count: ${var.elasticsearch_retention_days}
            field:
            stats_result:
            epoch:
            exclude: False
  EOT
  ]
}

# TODO: add helm chart for k8s cluster logging as optional component
