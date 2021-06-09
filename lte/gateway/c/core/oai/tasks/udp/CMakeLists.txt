# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

set(S1AP_C_DIR ${PROJECT_BINARY_DIR}/tasks/s1ap/r15)
set(NGAP_C_DIR ${PROJECT_BINARY_DIR}/tasks/ngap/r16)
include_directories(${S1AP_C_DIR})
include_directories(${NGAP_C_DIR})

add_library(TASK_UDP udp_primitives_server.c)

target_link_libraries(TASK_UDP
    LIB_BSTR LIB_HASHTABLE
    )
