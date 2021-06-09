# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PROJECT(datastore)
cmake_minimum_required(VERSION 3.7.2)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

include("${PROJECT_SOURCE_DIR}/../CMakeProtoMacros.txt")
set(MAGMA_INCLUDE_DIR $ENV{C_BUILD}/common/includes)
find_package(MAGMA_LOGGING REQUIRED)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(ASYNC_ORC8R_CPP_PROTOS common redis)
set(ASYNC_LTE_CPP_PROTOS session_manager policydb mobilityd subscriberdb)
set(ASYNC_LTE_GRPC_PROTOS session_manager)
set(ASYNC_ORC8R_GRPC_PROTOS "")

generate_all_protos("${ASYNC_LTE_CPP_PROTOS}" "${ASYNC_ORC8R_CPP_PROTOS}"
    "${ASYNC_LTE_GRPC_PROTOS}"
    "${ASYNC_ORC8R_GRPC_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}")

message("Async Proto_srcs are ${PROTO_SRCS}")

add_library(DATASTORE
    ObjectMap.h
    Serializers.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_link_libraries(DATASTORE PRIVATE
    tacopie cpp_redis pthread
    )

target_link_libraries(DATASTORE PRIVATE MAGMA_LOGGING)

# copy headers to build directory so they can be shared with OAI,
# session_manager, etc.
add_custom_command(TARGET DATASTORE POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E make_directory ${MAGMA_INCLUDE_DIR})
add_custom_command(TARGET DATASTORE POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E copy
    ${PROJECT_SOURCE_DIR}/includes/*.hpp ${MAGMA_INCLUDE_DIR})

target_include_directories(DATASTORE PUBLIC
    $<TARGET_FILE_DIR:DATASTORE> # protos
    "$<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}>")

install(TARGETS DATASTORE EXPORT DATASTORE_TARGETS
    INCLUDES DESTINATION includes
    ARCHIVE DESTINATION lib)

set(CMAKE_EXPORT_PACKAGE_REGISTRY ON)
export(TARGETS DATASTORE FILE DATASTOREConfig.cmake)
export(PACKAGE DATASTORE)
