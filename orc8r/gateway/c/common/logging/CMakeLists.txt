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

set(MAGMA_INCLUDE_DIR $ENV{C_BUILD}/common/includes)

add_library(MAGMA_LOGGING magma_logging_init.cpp magma_logging_init.h magma_logging.h)
# Specify include path for chained dependencies
target_include_directories(MAGMA_LOGGING PUBLIC
    $ENV{MAGMA_ROOT}
    "$<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}>")
target_link_libraries(MAGMA_LOGGING PUBLIC glog)

# copy headers to build directory so they can be shared with OAI,
# session_manager, etc.
add_custom_command(TARGET MAGMA_LOGGING POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E make_directory ${MAGMA_INCLUDE_DIR})
add_custom_command(TARGET MAGMA_LOGGING POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E copy
    ${PROJECT_SOURCE_DIR}/*.h ${MAGMA_INCLUDE_DIR})

install(TARGETS MAGMA_LOGGING EXPORT MAGMA_LOGGING
    INCLUDES DESTINATION ""
    ARCHIVE DESTINATION lib)

# Export package
set(CMAKE_EXPORT_PACKAGE_REGISTRY ON)
export(TARGETS MAGMA_LOGGING FILE MAGMA_LOGGINGConfig.cmake)
export(PACKAGE MAGMA_LOGGING)
