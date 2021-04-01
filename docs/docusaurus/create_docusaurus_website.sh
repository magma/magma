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

# spin until localhost:3000 returns HTTP code 200.
function spin() {
  nsec=${1-10}
  spin='-\|/'
  i=0
  while [[ "$(curl -s -o /dev/null -w '%{http_code}' localhost:3000)" != "200" ]]; do
    i=$(( i + 1 ))
    j=$(( i % 4 ))
    printf "\r${spin:$j:1}"
    sleep 1
  done
  printf "\r \n"
}

docker-compose down
docker build -t magma_docusaurus .
docker-compose up -d

echo ''
echo 'NOTE: README changes will live-reload. Sidebar changes require re-running this script.'
echo ''
echo 'Waiting for Docusaurus site to come up...'
echo 'If you want to follow the build logs, run docker-compose logs -f docusaurus'
spin
echo 'Navigate to http://localhost:3000/ to see the docs.'

open 'http://localhost:3000/docs/next/basics/introduction.html' || true

