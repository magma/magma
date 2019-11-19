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

# This section is unnecessary if using host networking
S6A_LOCAL_PORT={{ .Values.feg.bind.S6A_LOCAL_PORT | default "3868" }}
S6A_HOST_PORT={{ .Values.feg.bind.S6A_HOST_PORT | default "3869" }}
S6A_NETWORK={{ .Values.feg.bind.S6A_NETWORK | default "sctp" }}

SWX_LOCAL_PORT={{ .Values.feg.bind.SWX_LOCAL_PORT | default "3868" }}
SWX_HOST_PORT={{ .Values.feg.bind.SWX_HOST_PORT | default "3868" }}
SWX_NETWORK={{ .Values.feg.bind.SWX_NETWORK | default "sctp" }}

GX_LOCAL_PORT={{ .Values.feg.bind.GX_LOCAL_PORT | default "3907" }}
GX_HOST_PORT={{ .Values.feg.bind.GX_HOST_PORT | default "0" }}
GX_NETWORK={{ .Values.feg.bind.GX_NETWORK | default "tcp" }}

GY_LOCAL_PORT={{ .Values.feg.bind.GY_LOCAL_PORT | default "3906" }}
GY_HOST_PORT={{ .Values.feg.bind.GY_HOST_PORT | default "0" }}
GY_NETWORK={{ .Values.feg.bind.GY_NETWORK | default "tcp" }}
