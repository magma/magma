# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

COMPOSE_PROJECT_NAME=feg
DOCKER_REGISTRY={{ .Values.feg.image.docker_registry }}
DOCKER_USERNAME={{ .Values.feg.image.username }}
DOCKER_PASSWORD={{ .Values.feg.image.password }}
IMAGE_VERSION={{ .Values.feg.image.tag }}
BUILD_CONTEXT={{ .Values.feg.repo.url }}#{{ .Values.feg.repo.branch }}

ROOTCA_PATH=/var/opt/magma/certs/rootCA.pem
CONTROL_PROXY_PATH=/etc/magma/control_proxy.yml
SNOWFLAKE_PATH=/etc/snowflake

CERTS_VOLUME=/var/opt/magma/certs
CONFIGS_VOLUME=/var/opt/magma/configs

{{ if .Values.feg.env }}
{{ .Values.feg.env }}
{{- end }}
