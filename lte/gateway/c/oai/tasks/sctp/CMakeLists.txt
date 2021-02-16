# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

set(S1AP_C_DIR ${PROJECT_BINARY_DIR}/s1ap/r15)
include_directories(${S1AP_C_DIR})

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(SCTPD_LTE_CPP_PROTOS sctpd)
set(SCTPD_GRPC_PROTOS sctpd)

generate_all_protos(
  "${SCTPD_LTE_CPP_PROTOS}" "" "${SCTPD_GRPC_PROTOS}" ""
  "${PROTO_SRCS}" "${PROTO_HDRS}"
)

add_library(TASK_SCTP_SERVER
    sctpd_downlink_client.cpp
    sctpd_uplink_server.cpp
    sctp_itti_messaging.c
    sctp_primitives_server.c
    ${PROTO_SRCS}
    ${PROTO_HDRS}
)
target_link_libraries(TASK_SCTP_SERVER
    COMMON
    LIB_BSTR LIB_HASHTABLE
    grpc++ grpc protobuf
)
target_include_directories(TASK_SCTP_SERVER PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CMAKE_CURRENT_BINARY_DIR}
    ${PROJECT_BINARY_DIR}
)
