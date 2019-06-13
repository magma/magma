#!/bin/bash

# Copyright 2004-present Facebook. All Rights Reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

set -e

ROOT_DIR=$(cd -P "$(dirname "${BASH_SOURCE[0]}")" && hg root)

# shellcheck source=xplat/js/env-utils/setup_env_vars.sh
source "$ROOT_DIR/xplat/js/env-utils/setup_env_vars.sh"

pushd "$ROOT_DIR/xplat/fbc" >/dev/null
  "$INSTALL_NODE_MODULES"

  # manually install GRPC binaries
  # TODO remove this when node-pre-gyp@0.13.1 is released (see D15591571)
  mkdir -p node_modules/grpc/src/node/extension_binary/
  tar -xzf "$ROOT_DIR/xplat/third-party/grpc/v1.20.3/node-v57-linux-x64-glibc.tar.gz" -C node_modules/grpc/src/node/extension_binary/
  tar -xzf "$ROOT_DIR/xplat/third-party/grpc/v1.20.3/node-v64-linux-x64-glibc.tar.gz" -C node_modules/grpc/src/node/extension_binary/

  "$YARN_BINARY" run test "$@"
popd >/dev/null

