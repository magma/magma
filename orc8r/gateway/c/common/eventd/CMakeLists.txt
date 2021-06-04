# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)
PROJECT(Eventd)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

set(MAGMA_INCLUDE_DIR $ENV{C_BUILD}/common/includes)
include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)

set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

create_proto_dir("orc8r" CPP_OUT_DIR)
set(EVENTD_PROTOS common eventd)
generate_cpp_protos("${EVENTD_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}"
    ${ORC8R_PROTO_DIR} ${CPP_OUT_DIR})
generate_grpc_protos("${EVENTD_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}"
    ${ORC8R_PROTO_DIR} ${CPP_OUT_DIR})

add_library(EVENTD
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    EventdClient.cpp
    )

find_package(SERVICE_REGISTRY REQUIRED)
find_package(ASYNC_GRPC REQUIRED)
find_package(MAGMA_CONFIG REQUIRED)

target_link_libraries(
    EVENTD
    SERVICE_REGISTRY ASYNC_GRPC
    grpc++ grpc
)

# copy headers to build directory so they can be shared with OAI,
# session_manager, etc.
add_custom_command(TARGET EVENTD POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E make_directory ${MAGMA_INCLUDE_DIR})
add_custom_command(TARGET EVENTD POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E copy
    ${PROJECT_SOURCE_DIR}/includes/*.h ${MAGMA_INCLUDE_DIR})

target_include_directories(EVENTD PUBLIC
    "$<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}>"
    $<TARGET_FILE_DIR:EVENTD>)

install(TARGETS EVENTD EXPORT EVENTD_TARGETS
    INCLUDES DESTINATION includes
    ARCHIVE DESTINATION lib)

set(CMAKE_EXPORT_PACKAGE_REGISTRY ON)
export(TARGETS EVENTD FILE EVENTDConfig.cmake)
export(PACKAGE EVENTD)
