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
    <match eventd>
      @id eventd_elasticsearch
      @type elasticsearch
      @log_level info
      include_tag_key true
      host "#{ENV['OUTPUT_HOST']}"
      port "#{ENV['OUTPUT_PORT']}"
      scheme "#{ENV['OUTPUT_SCHEME']}"
      ssl_version "#{ENV['OUTPUT_SSL_VERSION']}"
      logstash_format true
      logstash_prefix "eventd"
      reconnect_on_error true
      reload_on_failure true
      reload_connections false
      log_es_400_reason true
      <buffer>
        @type file
        path /var/log/fluentd-buffers/eventd.kubernetes.system.buffer
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
extraVolumeMounts:
- mountPath: /certs
  name: certs
  readOnly: true
extraVolumes:
- name: certs
  secret:
    defaultMode: 420
    secretName: fluentd-certs
output:
  host: elasticsearch-master
  port: 9200
  scheme: http
rbac:
  create: false
replicaCount: 2
service:
  annotations:
    external-dns.alpha.kubernetes.io/hostname: fluentd.${dns_domain}
  ports:
  - containerPort: 24224
    name: forward
    protocol: TCP
  type: LoadBalancer
