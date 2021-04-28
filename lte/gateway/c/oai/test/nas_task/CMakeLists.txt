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

set(MME_APP_LIB_DIR $ENV{C_BUILD}/oai/tasks/mme_app)
set(NAS_LIB_DIR $ENV{C_BUILD}/oai/tasks/nas)

include_directories(${MME_APP_LIB_DIR})
include_directories(${NAS_LIB_DIR})
include_directories("${PROJECT_SOURCE_DIR}")
include_directories("/usr/src/googletest/googlemock/include/")

add_library(NAS_TASK_TEST_LIB
    state_creators.h
    state_creators.cpp
    )

link_directories(/usr/src/googletest/googlemock/lib/)

pkg_search_module(OPENSSL openssl REQUIRED)
include_directories(${OPENSSL_INCLUDE_DIRS})

pkg_search_module(CRYPTO libcrypto REQUIRED)
include_directories(${CRYPTO_INCLUDE_DIRS})

target_link_libraries(NAS_TASK_TEST_LIB
    TASK_MME_APP TASK_NAS 
    ${CRYPTO_LIBRARIES} ${OPENSSL_LIBRARIES}
    gmock_main gtest gtest_main gmock
    pthread rt yaml-cpp prometheus-cpp folly
    )

set_target_properties(NAS_TASK_TEST_LIB PROPERTIES LINKER_LANGUAGE CXX)

foreach (nas_test nas_state_converter)
  add_executable(${nas_test}_test test_${nas_test}.cpp)
  target_link_libraries(${nas_test}_test NAS_TASK_TEST_LIB)
  add_test(test_${nas_test} ${nas_test}_test)
endforeach (nas_test)
