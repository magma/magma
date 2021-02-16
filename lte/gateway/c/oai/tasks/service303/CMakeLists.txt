# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

find_package(Protobuf REQUIRED)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

add_library(TASK_SERVICE303
  service303.cpp
)

target_link_libraries(TASK_SERVICE303
  COMMON
  ${SERVICE303_LIB}
  ${ASYNC_GRPC}
  ${SERVICE_REGISTRY}
  ${CONFIG}
  LIB_BSTR
  ${PROTOBUF_LIBRARIES}
  )

target_include_directories(TASK_SERVICE303 PUBLIC
  ${MAGMA_LIB_DIR}/service303
  ${MAGMA_LIB_DIR}/async_grpc
  ${MAGMA_LIB_DIR}/config
  ${MAGMA_LIB_DIR}/service_registry
  ${CMAKE_CURRENT_SOURCE_DIR}
)
