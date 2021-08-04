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

build (){
  docker build . -f services/$1/Dockerfile -t $1
}

# build container for python services
docker build . -f services/build/Dockerfile.python -t pythonbuilder:latest

# python services
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
build eventd
build control_proxy


docker-compose build

