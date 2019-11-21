#!/usr/bin/env bash
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

cd ${NMS_ROOT}
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
COMPOSE_PROJECT_NAME=magmalte ../../../../orc8r/tools/docker/publish.sh -r ${DOCKER_REGISTRY} -i magmalte -v ${TRAVIS_COMMIT}
