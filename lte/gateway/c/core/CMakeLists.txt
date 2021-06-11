# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

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

set(CMAKE_EXPORT_COMPILE_COMMANDS ON)
set(CMAKE_VERBOSE_MAKEFILE ON)

project(MagmaConvergedCore)

include(ExternalProject)

ExternalProject_Add(MagmaCommon
    SOURCE_DIR "$ENV{MAGMA_ROOT}/orc8r/gateway/c/common"
    BUILD_ALWAYS 1
    BINARY_DIR "${CMAKE_BINARY_DIR}/common"
    INSTALL_COMMAND "")

ExternalProject_Add(MagmaCore
    SOURCE_DIR "${CMAKE_SOURCE_DIR}/oai"
    BUILD_ALWAYS 1
    BINARY_DIR "${CMAKE_BINARY_DIR}/oai"
    INSTALL_COMMAND ""
    DEPENDS MagmaCommon
    CMAKE_ARGS ${CL_ARGS})

