#!/usr/bin/env bash

################################################################################
# Copyright 2022 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
###############################################################################

set -euo pipefail

###############################################################################
# VARIABLES SECTION
###############################################################################

# The CACHE_KEY is the mandatory first argument.
CACHE_KEY=${1:-}

if [[ -z "$CACHE_KEY" ]]
then
  echo "Required argument CACHE_KEY not set!" >&2
  exit 1
fi

# The BAZEL_REMOTE_PASSWORD is an optional second argument.
BAZEL_REMOTE_PASSWORD=${2:-}

###############################################################################
# FUNCTIONS SECTION
###############################################################################

create_config () {
  local cache_key=$1
  local bazel_remote_password=$2
  if [[ -n "$bazel_remote_password" ]]
  then
    create_config_for_rw_remote_cache "$cache_key" "$bazel_remote_password"
  else
    create_config_for_ro_remote_cache "$cache_key"
  fi
}

create_config_for_rw_remote_cache () {
  local cache_key=$1
  local bazel_remote_password=$2
  sed \
    -e s/~~CACHE_KEY~~/"$cache_key"/ \
    -e s/~~BAZEL_REMOTE_PASSWORD~~/"$bazel_remote_password"/ \
    bazel/bazelrcs/remote_caching_rw.bazelrc
  echo "Configured bazel for read and write access to the remote cache with cache key: $cache_key" 1>&2
}

create_config_for_ro_remote_cache () {
  local cache_key=$1
  sed \
    -e s/~~CACHE_KEY~~/"$cache_key"/ \
    bazel/bazelrcs/remote_caching_ro.bazelrc
  echo "Configured bazel for read-only access to the remote cache with cache key: $cache_key" 1>&2
}

###############################################################################
# SCRIPT SECTION
###############################################################################

create_config "$CACHE_KEY" "$BAZEL_REMOTE_PASSWORD" > bazel/bazelrcs/cache.bazelrc
