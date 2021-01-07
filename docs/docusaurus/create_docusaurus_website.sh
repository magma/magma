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

cd "$(dirname "$0")/.."
docker build -f docusaurus/Dockerfile -t docusaurus-doc .
docker stop docs_container || true
docker run --rm -p 3000:3000 -d --name docs_container docusaurus-doc

echo ""
echo "Navigate to http://localhost:3000/ to see the docs."
