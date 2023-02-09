#!/bin/bash
# Copyright 2023 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# This script should be run from $MAGMA_ROOT

# Relative path to cache directory from $MAGMA_ROOT
RELATIVE_CACHE_DIR="$1"
CUTOFF_MB="$2"

# See https://stackoverflow.com/a/27485157 for reference.
CACHE_SIZE_MB=$(du -smc "$RELATIVE_CACHE_DIR" | grep "$RELATIVE_CACHE_DIR" | cut -f1)
echo "Total size of the Bazel cache (rounded up to MBs): $CACHE_SIZE_MB"

if [[ "$CACHE_SIZE_MB" -gt "$CUTOFF_MB" ]]; then
    echo "Cache exceeds the cut-off size, emptying the cache. The following build will be slower and refill the cache."
    rm -rf "$RELATIVE_CACHE_DIR"
fi
