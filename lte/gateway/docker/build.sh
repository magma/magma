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
# build container for python services
docker build . -f services/build/Dockerfile.python -t pythonbuilder:latest

build (){
  docker build . -f services/$1/Dockerfile -t $1
}
#
build mobilityd
build enodebd
build health
build policydb
build smsd
build subscriberdb
build ctraced
build magmad
build state
build directoryd
build pipelined

cd ../../../

docker build . -f lte/gateway/docker/mme/Dockerfile.ubuntu20.04 -t mme_builder:latest

cd lte/gateway/docker

docker build . -f services/build/Dockerfile.c -t cbuilder:latest

build mme
build sctpd
