#!/bin/bash

# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

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

pushd "$ROOT_DIR/xplat/fbc-mobile-app" >/dev/null
  "$YARN_BINARY" "$@"
  "$YARN_BINARY" relay "$@"
  repo_status=$(hg status | wc -l)
  if [ "$repo_status" -ne 0 ]; then
    echo "'yarn relay' modified changes. Please run 'yarn relay' from xplat/fbc-mobile-app" >&2
    exit 1
  fi
popd >/dev/null
