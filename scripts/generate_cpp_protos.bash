#!/usr/bin/env bash
# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

# generate_cpp_protos.bash invokes protoc to produce .pb.cc to a dst dir
# Note: this is just a dev helper script; the generate_cpp_protos CMake macro is defined in
# orc8r/gateway/c/common/CMakeProtoMacros.txt

if [[ $# -eq 0 ]] ; then
    echo 'usage: build_protos.bash <dst dir, e.g. $MAGMA_ROOT/build>'
    exit 0
fi

mr="${MAGMA_ROOT:-/home/vagrant/magma}"

dst=$1
mkdir -p $dst

echo "- LTE pb generaton";for i in $mr/lte/protos/*.proto $mr/lte/protos/oai/*.proto $mr/lte/protos/mconfig/*.proto;do protoc -I $mr/orc8r/protos/prometheus -I $mr --cpp_out="$dst" "$i";  echo "   $i"; done;
echo "- LTE grpc.pb generation"; for i in $mr/lte/protos/*.proto;do protoc -I $mr/orc8r/protos/prometheus -I $mr -I $mr/lte/protos --grpc_out="$dst" --plugin=protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin $i; echo "   $i"; done;
echo "- FEG pb generaton"; for i in $mr/feg/protos/*.proto;do protoc -I $mr/orc8r/protos/prometheus -I $mr --cpp_out="$dst" "$i"; echo "   $i"; done;
echo "- FEG grpc.pb generation"; for i in $mr/lte/protos/*.proto;do protoc -I $mr/orc8r/protos/prometheus -I $mr -I $mr/feg/protos --grpc_out="$dst" --plugin=protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin  $i; echo "   $i"; done;
echo "- Orc8r pb generaton"; for i in $mr/orc8r/protos/*.proto;do protoc -I $mr/orc8r/protos/prometheus -I $mr --cpp_out="$dst" "$i";echo "   $i"; done;
echo "- Orc8r grpc.pb generation"; for i in $mr/lte/protos/*.proto;do protoc -I $mr/orc8r/protos/prometheus -I $mr -I $mr/lte/protos --grpc_out="$dst" --plugin=protoc-gen-grpc=/usr/local/bin/grpc_cpp_plugin  $i; echo "   $i"; done;
