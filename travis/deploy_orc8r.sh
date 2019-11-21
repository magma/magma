#!/usr/bin/env bash
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

cd ${TRAVIS_BUILD_DIR}/magma/orc8r/cloud/docker
echo "$DOCKER_PASSWORD" | docker login -u "$DOCKER_USERNAME" --password-stdin
for image in "$@"
do
    ../../../orc8r/tools/docker/publish.sh -r ${DOCKER_REGISTRY} -i ${image} -v ${TRAVIS_COMMIT}
done
