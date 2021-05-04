# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

include_directories("${PROJECT_SOURCE_DIR}/../common/config")

add_library(SERVICE_REGISTRY
    ServiceRegistrySingleton.h
    ServiceRegistrySingleton.cpp
    )

target_link_libraries(SERVICE_REGISTRY CONFIG
   grpc++ dl yaml-cpp
  )

# copy headers to build directory so they can be shared with OAI,
# session_manager, etc.
add_custom_command(TARGET SERVICE_REGISTRY POST_BUILD
                   COMMAND ${CMAKE_COMMAND} -E copy
                   ${PROJECT_SOURCE_DIR}/service_registry/*.h
                   $<TARGET_FILE_DIR:SERVICE_REGISTRY>)

target_include_directories(SERVICE_REGISTRY PUBLIC
                   $<TARGET_FILE_DIR:SERVICE_REGISTRY>
)
