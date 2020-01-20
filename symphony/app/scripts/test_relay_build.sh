#!/bin/bash

# Copyright 2004-present Facebook. All Rights Reserved.
#
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

set -e

ROOT_DIR=$(cd -P "$(dirname "${BASH_SOURCE[0]}")" && hg root)

# shellcheck source=xplat/js/env-utils/setup_env_vars.sh
source "$ROOT_DIR/xplat/js/env-utils/setup_env_vars.sh"

pushd "$ROOT_DIR/xplat/fbc/fbcnms-projects/inventory" >/dev/null
  "$YARN_BINARY" relay "$@"
  repo_status=$(hg status | wc -l)
  if [ "$repo_status" -ne 0 ]; then
    echo "'yarn relay' modified changes. Please run 'yarn relay' from xplat/fbc/fbcnms-projects/inventory" >&2
    exit 1
  fi
popd >/dev/null
