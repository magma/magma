# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)
PROJECT(MagmaSentry)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

# compile the needed macros
include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)
create_proto_dir("orc8r" ORC8R_CPP_OUT_DIR)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(SMGR_ORC8R_CPP_PROTOS common mconfig/mconfigs)
generate_cpp_protos("${SMGR_ORC8R_CPP_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_CPP_OUT_DIR})

add_library(MAGMA_SENTRY
  SentryWrapper.cpp
  ${PROTO_SRCS}
  ${PROTO_HDRS})

find_package(MAGMA_CONFIG REQUIRED)

## Find + link sentry if it exists
find_library(SENTRY_LIB NAMES sentry sentry-native)
if (SENTRY_LIB)
  # Include dir
  message("Building with sentry!")
  find_path(SENTRY_INCLUDE_DIR NAMES sentry.h )
  target_link_libraries(MAGMA_SENTRY ${SENTRY_LIB} MAGMA_CONFIG)
  set (CMAKE_CXX_FLAGS_DEBUG "${CMAKE_CXX_FLAGS_DEBUG} -DSENTRY_ENABLED=1")
  set (CMAKE_CXX_FLAGS_RELWITHDEBINFO "${CMAKE_CXX_FLAGS_RELWITHDEBINFO} -DSENTRY_ENABLED=1")
endif()

target_include_directories(MAGMA_SENTRY PUBLIC
    $ENV{MAGMA_ROOT}
    ${CMAKE_CURRENT_SOURCE_DIR}
    $<TARGET_FILE_DIR:MAGMA_SENTRY>)

install(TARGETS MAGMA_SENTRY EXPORT MAGMA_SENTRY_TARGETS
    INCLUDES DESTINATION includes
    ARCHIVE DESTINATION lib)

set(CMAKE_EXPORT_PACKAGE_REGISTRY ON)
export(TARGETS MAGMA_SENTRY FILE MAGMA_SENTRYConfig.cmake)
export(PACKAGE MAGMA_SENTRY)
