#!/bin/bash
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#

set -e

# shellcheck disable=SC1091
source /etc/environment

SWAGGER_TMP_DIR=/var/tmp/swagger

# Cleanup swagger tmp directory
rm -rf "$SWAGGER_TMP_DIR"
mkdir -p "$SWAGGER_TMP_DIR"

# Find all the other swagger files, skip hidden paths (e.g. .cache) and copy them
find "$MAGMA_ROOT" -not -path '*/\.*' -regex ".*swagger/.*.yml" -print0 | xargs -I% -0 cp % "$SWAGGER_TMP_DIR"

# Copy all the required swagger files to the tmp directory
cp "$SWAGGER_SPEC" "$SWAGGER_TMP_DIR"

/usr/bin/java -jar "$SWAGGER_CODEGEN_JAR" generate -i "$SWAGGER_TMP_DIR"/swagger.yml -o "$SWAGGER_CODEGEN_OUTPUT" -l python
