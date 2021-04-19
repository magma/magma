#!/usr/bin/env bash

# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# consolidate_protos.bash consolidates magma and imported proto files
# to /tmp/magma_protos for easy consumption by e.g. Wireshark.

set -e

outdir=/tmp/magma_protos/

magma=${MAGMA_ROOT-~/magma}
echo $magma
include=/usr/local/include

ignore=( -not -path '*/migrations/*' )

pushd ${magma}
find -L . -name '*.proto' "${ignore[@]}" | cpio -pdm --insecure ${outdir}
popd

pushd ${include}
find -L . -name '*.proto' "${ignore[@]}" | cpio -pdm --insecure ${outdir}
popd

echo
echo copied protos to ${outdir}
