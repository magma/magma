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

# compile the needed protos
set(DRD_ORC8R_CPP_PROTOS common directoryd)
set(DRD_LTE_CPP_PROTOS "")
set(DRD_ORC8R_GRPC_PROTOS directoryd)
set(DRD_LTE_GRPC_PROTOS "")

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

generate_all_protos("${DRD_LTE_CPP_PROTOS}" "${DRD_ORC8R_CPP_PROTOS}"
  "${DRD_LTE_GRPC_PROTOS}" "${DRD_ORC8R_GRPC_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}")

include_directories(${CMAKE_CURRENT_BINARY_DIR})

add_library(LIB_DIRECTORYD
  directoryd.cpp
  GatewayDirectorydClient.cpp
  ${PROTO_SRCS}
  ${PROTO_HDRS}
  )

target_link_libraries(LIB_DIRECTORYD
  grpc grpc++
  LIB_MOBILITY_CLIENT ${ASYNC_GRPC} ${SERVICE_REGISTRY}
  ${CONFIG}
)

target_include_directories(LIB_DIRECTORYD PUBLIC
  ${MAGMA_LIB_DIR}/async_grpc
  ${MAGMA_LIB_DIR}/config
  ${MAGMA_LIB_DIR}/service_registry
  ${CMAKE_CURRENT_SOURCE_DIR}
)
