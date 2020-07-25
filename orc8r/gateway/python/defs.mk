# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
ifndef PROTO_LIST
PROTO_LIST:=orc8r_protos
endif

TESTS=magma/common/redis/tests \
      magma/common/tests \
      magma/configuration/tests \
      magma/directoryd/tests \
      magma/eventd/tests \
      magma/magmad/check/tests \
      magma/magmad/check/kernel_check/tests \
      magma/magmad/check/machine_check/tests \
      magma/magmad/check/network_check/tests \
      magma/magmad/tests \
      magma/magmad/upgrade/tests \
      magma/magmad/generic_command/tests \
      magma/state/tests \
      magma/tests
