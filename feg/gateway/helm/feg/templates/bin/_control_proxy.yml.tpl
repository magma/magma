---
#
# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

# nghttpx config will be generated here and used
nghttpx_config_location: /var/tmp/nghttpx.conf

# Location for certs
rootca_cert: /var/opt/magma/certs/rootCA.pem
gateway_cert: /var/opt/magma/certs/gateway.crt
gateway_key: /var/opt/magma/certs/gateway.key

# Listening port of the proxy for local services. The port would be closed
# for the rest of the world.
local_port: {{ .Values.feg.proxy.local_port }}

# Cloud address for reaching out to the cloud.
cloud_address: {{ .Values.feg.proxy.cloud_address }}
cloud_port: {{ .Values.feg.proxy.cloud_port }}

bootstrap_address: {{ .Values.feg.proxy.bootstrap_address }}
bootstrap_port: {{ .Values.feg.proxy.bootstrap_port }}

# Option to use nghttpx for proxying. If disabled, the individual
# services would establish the TLS connections themselves.
proxy_cloud_connections: True

# Allows http_proxy usage if the environment variable is present
allow_http_proxy: True
