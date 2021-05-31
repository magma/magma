# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)
PROJECT(ServiceRegistry)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

set(MAGMA_INCLUDE_DIR $ENV{C_BUILD}/common/includes)

add_library(SERVICE_REGISTRY
    ServiceRegistrySingleton.cpp
    )

find_package(MAGMA_CONFIG REQUIRED)

target_link_libraries(SERVICE_REGISTRY MAGMA_CONFIG
    grpc++ dl yaml-cpp
    )

# copy headers to build directory so they can be shared with OAI,
# session_manager, etc.
add_custom_command(TARGET SERVICE_REGISTRY POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E make_directory ${MAGMA_INCLUDE_DIR})
add_custom_command(TARGET SERVICE_REGISTRY POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E copy
    ${PROJECT_SOURCE_DIR}/includes/*.h ${MAGMA_INCLUDE_DIR})

target_include_directories(SERVICE_REGISTRY PUBLIC
    "$<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}>")

install(TARGETS SERVICE_REGISTRY EXPORT SERVICE_REGISTRY_TARGETS
    INCLUDES DESTINATION includes
    ARCHIVE DESTINATION lib)

set(CMAKE_EXPORT_PACKAGE_REGISTRY ON)
export(TARGETS SERVICE_REGISTRY FILE SERVICE_REGISTRYConfig.cmake)
export(PACKAGE SERVICE_REGISTRY)
