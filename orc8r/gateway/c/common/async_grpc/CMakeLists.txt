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

add_library(ASYNC_GRPC
    GRPCReceiver.cpp
    )

target_link_libraries(ASYNC_GRPC PRIVATE MAGMA_LOGGING)

target_include_directories(ASYNC_GRPC PUBLIC
    "$<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}>"
    )

install(TARGETS ASYNC_GRPC EXPORT ASYNC_GRPC_TARGETS
    INCLUDES DESTINATION ""
    ARCHIVE DESTINATION lib)

set(CMAKE_EXPORT_PACKAGE_REGISTRY ON)
export(TARGETS ASYNC_GRPC FILE ASYNC_GRPCConfig.cmake)
export(PACKAGE ASYNC_GRPC)

# Copy headers to build directory so they can be shared with session_manager
# Can be removed when session manager moves to a super build model.
set(MAGMA_INCLUDE_DIR $ENV{C_BUILD}/common/includes)
add_custom_command(TARGET ASYNC_GRPC POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E make_directory ${MAGMA_INCLUDE_DIR})
add_custom_command(TARGET ASYNC_GRPC POST_BUILD
    COMMAND ${CMAKE_COMMAND} -E copy
    ${PROJECT_SOURCE_DIR}/includes/*.h ${MAGMA_INCLUDE_DIR})
