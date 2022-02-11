# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

PROJECT(MagmaCommon)

include(ExternalProject)
set(CMAKE_EXPORT_PACKAGE_REGISTRY ON)
set(CMAKE_EXPORT_COMPILE_COMMANDS ON)
set(CMAKE_VERBOSE_MAKEFILE ON)

# MUST be done before call to 'project'
get_cmake_property(vars CACHE_VARIABLES)
foreach (var ${vars})
  get_property(currentHelpString CACHE "${var}" PROPERTY HELPSTRING)
  # message("${var} = [${${var}}]  --  ${currentHelpString}")
  if ("${currentHelpString}" MATCHES
      "No help, variable specified on the command line." OR
      "${currentHelpString}" STREQUAL "")
    list(APPEND CL_ARGS "-D${var}=${${var}}")
  endif ()
endforeach ()

list(APPEND CL_ARGS "-DCMAKE_BUILD_TYPE=${CMAKE_BUILD_TYPE}")

ExternalProject_Add(MagmaLogging
    SOURCE_DIR ${CMAKE_SOURCE_DIR}/logging
    BINARY_DIR ${CMAKE_BINARY_DIR}/logging
    INSTALL_COMMAND ""
    CMAKE_ARGS ${CL_ARGS})

ExternalProject_Add(AsyncGrpc
    SOURCE_DIR ${CMAKE_SOURCE_DIR}/async_grpc
    BINARY_DIR ${CMAKE_BINARY_DIR}/async_grpc
    INSTALL_COMMAND ""
    DEPENDS MagmaLogging
    CMAKE_ARGS ${CL_ARGS})

ExternalProject_Add(MagmaConfig
    SOURCE_DIR ${CMAKE_SOURCE_DIR}/config
    BINARY_DIR ${CMAKE_BINARY_DIR}/config
    BUILD_ALWAYS 1
    INSTALL_COMMAND ""
    DEPENDS MagmaLogging
    CMAKE_ARGS ${CL_ARGS})

ExternalProject_Add(ServiceRegistry
    SOURCE_DIR ${CMAKE_SOURCE_DIR}/service_registry
    BINARY_DIR ${CMAKE_BINARY_DIR}/service_registry
    BUILD_ALWAYS 1
    INSTALL_COMMAND ""
    DEPENDS MagmaConfig
    CMAKE_ARGS ${CL_ARGS})

ExternalProject_Add(Service303
    SOURCE_DIR ${CMAKE_SOURCE_DIR}/service303
    BINARY_DIR ${CMAKE_BINARY_DIR}/service303
    BUILD_ALWAYS 1
    INSTALL_COMMAND ""
    DEPENDS ServiceRegistry
    DEPENDS MagmaConfig
    DEPENDS MagmaLogging
    CMAKE_ARGS ${CL_ARGS})

ExternalProject_Add(Eventd
    SOURCE_DIR ${CMAKE_SOURCE_DIR}/eventd
    BINARY_DIR ${CMAKE_BINARY_DIR}/eventd
    BUILD_ALWAYS 1
    INSTALL_COMMAND ""
    DEPENDS AsyncGrpc
    DEPENDS ServiceRegistry
    CMAKE_ARGS ${CL_ARGS})

ExternalProject_Add(MagmaSentry
    SOURCE_DIR ${CMAKE_SOURCE_DIR}/sentry
    BINARY_DIR ${CMAKE_BINARY_DIR}/sentry
    BUILD_ALWAYS 1
    INSTALL_COMMAND ""
    DEPENDS MagmaConfig
    CMAKE_ARGS ${CL_ARGS})

ExternalProject_Add(MagmaEbpf
    SOURCE_DIR ${CMAKE_SOURCE_DIR}/ebpf
    BINARY_DIR ${CMAKE_BINARY_DIR}/ebpf
    BUILD_ALWAYS 1
    INSTALL_COMMAND ""
    CMAKE_ARGS ${CL_ARGS})
