#!/usr/bin/env bash
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set -x
# build container

# echo ${DOCKER_PASSWORD} > /tmp/passfile

build (){
  docker build . -f services/$1/Dockerfile -t $1
  #cp built-python > mobilityd
  #../../../orc8r/tools/docker/publish.sh -r ${DOCKER_REGISTRY} -i $1 -u ${DOCKER_USERNAME} -p /tmp/passfile
}

build build
# build mobilityd
# build enodedb
# build health
# build monitord
# build pipelined
# build pkt_tester
# build policydb
# build redirectd
# build smsd
# build subscriberd
# build tests
