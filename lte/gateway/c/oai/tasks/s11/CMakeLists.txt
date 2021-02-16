# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

include_directories(${S1AP_C_DIR})
include_directories(${PROJECT_SOURCE_DIR}/lib/gtpv2-c/gtpv2c_ie_formatter/shared)
include_directories(converged-mme/lte/gateway/c/oai/lib/gtpv2-c/nwgtpv2c-0.11/include)
include_directories(converged-mme/lte/gateway/c/oai/lib/gtpv2-c/nwgtpv2c-0.11/shared)

add_library(TASK_S11
s11_mme_bearer_manager.c
s11_mme_session_manager.c
s11_tasks.c
s11_common.c
s11_ie_formatter.c
)

target_link_libraries(TASK_S11
    COMMON
    LIB_BSTR LIB_HASHTABLE LIB_SGS_CLIENT LIB_GTPV2C TASK_MME_APP TASK_UDP
    ${PROTOBUF_LIBRARIES}
)
