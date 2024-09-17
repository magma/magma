# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
#
PYTHON_SRCS=$(MAGMA_ROOT)/lte/gateway/python $(MAGMA_ROOT)/orc8r/gateway/python
PROTO_LIST:=orc8r_protos lte_protos feg_protos dp_protos
SWAGGER_LIST:=lte_swagger_specs orc8r_swagger_specs

# Path to the test files
TESTS=magma/tests \
      magma/policydb/tests \
      magma/enodebd/tests \
      magma/mobilityd/tests \
      magma/pipelined/openflow/tests \
      magma/redirectd/tests \
      magma/subscriberdb/tests \
      magma/monitord/tests \
      magma/kernsnoopd/tests

SUDO_TESTS=magma/mobilityd/tests/test_ip_allocator_dhcp_e2e.py \
      magma/mobilityd/tests/test_ip_allocator_dhcp_with_vlan.py \
      magma/pipelined/tests
