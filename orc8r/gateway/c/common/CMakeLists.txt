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

ExternalProject_Add(MagmaLogging
        SOURCE_DIR ${CMAKE_SOURCE_DIR}/logging
        BINARY_DIR ${CMAKE_BINARY_DIR}/logging
        INSTALL_COMMAND "")

ExternalProject_Add(AsyncGrpc
        SOURCE_DIR ${CMAKE_SOURCE_DIR}/async_grpc
        BINARY_DIR ${CMAKE_BINARY_DIR}/async_grpc
        INSTALL_COMMAND ""
        DEPENDS MagmaLogging)

ExternalProject_Add(MagmaConfig
        SOURCE_DIR ${CMAKE_SOURCE_DIR}/config
        BINARY_DIR ${CMAKE_BINARY_DIR}/config
        INSTALL_COMMAND ""
        DEPENDS MagmaLogging)

ExternalProject_Add(ServiceRegistry
        SOURCE_DIR ${CMAKE_SOURCE_DIR}/service_registry
        BINARY_DIR ${CMAKE_BINARY_DIR}/service_registry
        INSTALL_COMMAND ""
        DEPENDS MagmaConfig)

ExternalProject_Add(Service303
        SOURCE_DIR ${CMAKE_SOURCE_DIR}/service303
        BINARY_DIR ${CMAKE_BINARY_DIR}/service303
        INSTALL_COMMAND ""
        DEPENDS ServiceRegistry
        DEPENDS MagmaConfig
        DEPENDS MagmaLogging)

ExternalProject_Add(Eventd
        SOURCE_DIR ${CMAKE_SOURCE_DIR}/eventd
        BINARY_DIR ${CMAKE_BINARY_DIR}/eventd
        INSTALL_COMMAND ""
        DEPENDS AsyncGrpc
        DEPENDS ServiceRegistry)

if (BUILD_TESTS)
  ENABLE_TESTING()
  ADD_SUBDIRECTORY(test)
endif (BUILD_TESTS)
