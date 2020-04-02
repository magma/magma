# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

# Copyright (c) 2016-present, Facebook, Inc.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree. An additional grant
# of patent rights can be found in the PATENTS file in the same directory.

COMPOSE_PROJECT_NAME=cwf
DOCKER_REGISTRY={{ .Values.cwf.image.docker_registry }}
DOCKER_USERNAME={{ .Values.cwf.image.username }}
DOCKER_PASSWORD={{ .Values.cwf.image.password }}
IMAGE_VERSION={{ .Values.cwf.image.tag }}
BUILD_CONTEXT={{ .Values.cwf.repo.url }}#{{ .Values.cwf.repo.branch }}

ROOTCA_PATH=/var/opt/magma/certs/rootCA.pem
CONTROL_PROXY_PATH=/etc/magma/control_proxy.yml
CONFIGS_TEMPLATES_PATH=/etc/magma/templates

CERTS_VOLUME=/var/opt/magma/certs
CONFIGS_OVERRIDE_VOLUME=/var/opt/magma/configs
CONFIGS_DEFAULT_VOLUME=/etc/magma

{{ if .Values.cwf.env }}
{{ .Values.cwf.env }}
{{- end }}
