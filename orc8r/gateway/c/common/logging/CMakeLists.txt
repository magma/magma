# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PROJECT(MagmaLogging)
cmake_minimum_required(VERSION 3.7.2)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

add_library(MAGMA_LOGGING INTERFACE)
# Specify include path for chained dependencies
target_include_directories(MAGMA_LOGGING INTERFACE ${CMAKE_CURRENT_SOURCE_DIR})

install(TARGETS MAGMA_LOGGING EXPORT MAGMA_LOGGING
    INCLUDES DESTINATION ""
    ARCHIVE DESTINATION lib)

# Export package
set(CMAKE_EXPORT_PACKAGE_REGISTRY ON)
export(TARGETS MAGMA_LOGGING FILE MAGMA_LOGGINGConfig.cmake)
export(PACKAGE MAGMA_LOGGING)
