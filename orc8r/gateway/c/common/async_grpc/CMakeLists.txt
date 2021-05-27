# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

PROJECT(MagmaAsyncGrpc)
set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

find_package(MAGMA_LOGGING REQUIRED)

include("${PROJECT_SOURCE_DIR}/../CMakeProtoMacros.txt")

#compile the relevant protos

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(ASYNC_ORC8R_CPP_PROTOS common redis eventd)
set(ASYNC_LTE_CPP_PROTOS session_manager abort_session policydb subscriberdb
    mobilityd pipelined)
set(ASYNC_LTE_GRPC_PROTOS session_manager abort_session mobilityd pipelined)
set(ASYNC_ORC8R_GRPC_PROTOS eventd)
set(MAGMA_INCLUDE_DIR $ENV{C_BUILD}/common/includes)

generate_all_protos("${ASYNC_LTE_CPP_PROTOS}" "${ASYNC_ORC8R_CPP_PROTOS}"
    "${ASYNC_LTE_GRPC_PROTOS}"
    "${ASYNC_ORC8R_GRPC_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}")

message("Async Proto_srcs are ${PROTO_SRCS}")

add_library(ASYNC_GRPC
    GRPCReceiver.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_link_libraries(ASYNC_GRPC PRIVATE MAGMA_LOGGING)

# copy headers to build directory so they can be shared with OAI,
# session_manager, etc.
add_custom_command(TARGET ASYNC_GRPC POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E make_directory ${MAGMA_INCLUDE_DIR})
add_custom_command(TARGET ASYNC_GRPC POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E copy
    ${PROJECT_SOURCE_DIR}/includes/*.h ${MAGMA_INCLUDE_DIR})

target_include_directories(ASYNC_GRPC PUBLIC
    "$<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}>"
    $<TARGET_FILE_DIR:ASYNC_GRPC>
    )

install(TARGETS ASYNC_GRPC EXPORT ASYNC_GRPC_TARGETS
    INCLUDES DESTINATION ""
    ARCHIVE DESTINATION lib)

set(CMAKE_EXPORT_PACKAGE_REGISTRY ON)
export(TARGETS ASYNC_GRPC FILE ASYNC_GRPCConfig.cmake)
export(PACKAGE ASYNC_GRPC)
