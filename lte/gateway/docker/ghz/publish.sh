#!/bin/bash
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

REGISTRY=$1
VERSION=latest

publish () {
  docker tag ghz_gateway_"$1":latest "${REGISTRY}"/ghz_gateway_"$1":"${VERSION}"
  docker push "${REGISTRY}/"ghz_gateway_"$1":"${VERSION}"
}

publish python
publish c
