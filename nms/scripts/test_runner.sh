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
  "$YARN_BINARY" run jest "$@"
popd >/dev/null

