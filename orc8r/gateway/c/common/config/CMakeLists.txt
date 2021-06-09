# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PROJECT(MagmaConfig)
cmake_minimum_required(VERSION 3.7.2)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

find_package(MAGMA_LOGGING REQUIRED)
set(MAGMA_INCLUDE_DIR $ENV{C_BUILD}/common/includes)

add_library(MAGMA_CONFIG
    MConfigLoader.cpp
    ServiceConfigLoader.cpp
    YAMLUtils.cpp
    )
target_link_libraries(MAGMA_CONFIG PRIVATE MAGMA_LOGGING glog)

# copy headers to build directory so they can be shared with OAI,
# session_manager, etc.
add_custom_command(TARGET MAGMA_CONFIG POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E make_directory ${MAGMA_INCLUDE_DIR})
add_custom_command(TARGET MAGMA_CONFIG POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E copy
    ${PROJECT_SOURCE_DIR}/includes/*.h ${MAGMA_INCLUDE_DIR})

target_include_directories(MAGMA_CONFIG PUBLIC
    "$<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}>")

install(TARGETS MAGMA_CONFIG EXPORT MAGMA_CONFIG_TARGETS
    INCLUDES DESTINATION includes
    ARCHIVE DESTINATION lib)

set(CMAKE_EXPORT_PACKAGE_REGISTRY ON)
export(TARGETS MAGMA_CONFIG FILE MAGMA_CONFIGConfig.cmake)
export(PACKAGE MAGMA_CONFIG)
