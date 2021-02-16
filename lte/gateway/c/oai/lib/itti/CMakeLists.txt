# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

# TODO: ITTI_Messages shouldn't depend on S1_Types
set(ASN1RELDIR r15)
set(asn1_generated_dir ${PROJECT_BINARY_DIR}/s1ap)
set(asn1c_call ${CMAKE_CURRENT_SOURCE_DIR}/generate_asn1)
set(S1AP_C_DIR ${asn1_generated_dir}/${ASN1RELDIR})
include_directories("${S1AP_C_DIR}")

set(ITTI_FILES
    intertask_interface.c
    signals.c
    timer.c
)
add_library(LIB_ITTI ${ITTI_FILES})
target_link_libraries(LIB_ITTI
    rt pthread czmq
    COMMON
    LIB_BSTR LIB_HASHTABLE LIB_MESSAGE_UTILS
)
target_include_directories(LIB_ITTI PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
)
