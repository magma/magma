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

find_package(Protobuf REQUIRED)
# compile the needed protos
list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")
# compile MME APP STATE PROTOS
set(MMEAPP_STATE_CPP_PROTOS mme_nas_state nas_state spgw_state common_types
  std_3gpp_types)

set(LTE_OUT_MME_STATE_DIR "${CMAKE_CURRENT_BINARY_DIR}/lte/protos/oai")
generate_cpp_protos("${MMEAPP_STATE_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${STATE_PROTO_DIR} ${LTE_OUT_MME_STATE_DIR})

# compile the needed protos
set(HA_LTE_CPP_PROTOS ha_orc8r)
set(HA_LTE_GRPC_PROTOS ha_orc8r)
set(HA_ORC8R_CPP_PROTOS "")
set(HA_ORC8R_GRPC_PROTOS "")
generate_all_protos("${HA_LTE_CPP_PROTOS}" "${HA_ORC8R_CPP_PROTOS}" "${HA_LTE_GRPC_PROTOS}"
  "${HA_ORC8R_CPP_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}")

include_directories(${CMAKE_CURRENT_BINARY_DIR})
include_directories(${LTE_OUT_MME_STATE_DIR})

add_library(TASK_HA
    ha_task.c
    ha_service_handler.cpp
    HaClient.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_link_libraries(TASK_HA
    COMMON
    LIB_BSTR LIB_HASHTABLE protobuf
    ${PROTOBUF_LIBRARIES}
    TASK_MME_APP
    TASK_S1AP
)

target_include_directories(TASK_HA PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
)
