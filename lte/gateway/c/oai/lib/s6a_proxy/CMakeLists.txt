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
set(S6A_FEG_CPP_PROTOS s6a_proxy)
set(S6A_FEG_GRPC_PROTOS s6a_proxy)
set(S6A_LTE_CPP_PROTOS mconfig/mconfigs)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

create_proto_dir("feg" FEG_OUT_DIR)
create_proto_dir("lte" LTE_OUT_DIR)

generate_cpp_protos("${S6A_FEG_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})
generate_grpc_protos("${S6A_FEG_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${FEG_PROTO_DIR} ${FEG_OUT_DIR})

generate_cpp_protos("${S6A_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})

include_directories(${CMAKE_CURRENT_BINARY_DIR})

add_library(LIB_S6A_PROXY
  S6aClient.cpp
  s6a_client_api.cpp
  itti_msg_to_proto_msg.cpp
  proto_msg_to_itti_msg.cpp
  ${PROTO_SRCS}
  ${PROTO_HDRS}
  )

target_link_libraries(LIB_S6A_PROXY
  COMMON
  LIB_MOBILITY_CLIENT ${ASYNC_GRPC} ${SERVICE_REGISTRY}
  ${CONFIG} LIB_BSTR LIB_HASHTABLE
)
target_include_directories(LIB_S6A_PROXY PUBLIC
  ${MAGMA_LIB_DIR}/async_grpc
  ${MAGMA_LIB_DIR}/config
  ${MAGMA_LIB_DIR}/service_registry
  ${CMAKE_CURRENT_SOURCE_DIR}
)
