#!/bin/bash
# Copyright (c) Facebook, Inc. and its affiliates.
# All rights reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

set -e

cd "$(dirname "$0")/.."
docker build -f docusaurus/Dockerfile -t docusaurus-doc .
docker stop docs_container || true
docker run --rm -p 3000:3000 -d --name docs_container docusaurus-doc

rm -rf ./web
docker cp docs_container:/app/website/build/Magma ./web

echo "Navigate to http://127.0.0.1:3000/magma/web/ to see the docs."
