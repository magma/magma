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

set (CMAKE_CXX_FLAGS "-Wno-write-strings -Wno-literal-suffix")
add_definitions(-DDEBUG_IS_ON=1)

# compile the needed protos
set(PCEF_ORC8R_CPP_PROTOS mconfig/mconfigs)
set(PCEF_LTE_CPP_PROTOS policydb session_manager)
set(PCEF_ORC8R_GRPC_PROTOS "")
set(PCEF_LTE_GRPC_PROTOS session_manager)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

generate_all_protos("${PCEF_LTE_CPP_PROTOS}" "${PCEF_ORC8R_CPP_PROTOS}"
  "${PCEF_LTE_GRPC_PROTOS}" "${PCEF_ORC8R_GRPC_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}")

include_directories(${CMAKE_CURRENT_BINARY_DIR})

add_library(LIB_PCEF
  pcef_handlers.cpp
  PCEFClient.cpp
  ${PROTO_SRCS}
  ${PROTO_HDRS}
  )

target_link_libraries(LIB_PCEF
  COMMON
  ${ITTI_LIB} ${ASYNC_GRPC} ${SERVICE_REGISTRY}
  ${CONFIG} LIB_BSTR LIB_HASHTABLE LIB_MOBILITY_CLIENT yaml-cpp
)
target_include_directories(LIB_PCEF PUBLIC
  ${MAGMA_LIB_DIR}/async_grpc
  ${MAGMA_LIB_DIR}/config
  ${MAGMA_LIB_DIR}/service_registry
  ${CMAKE_CURRENT_SOURCE_DIR}
)
