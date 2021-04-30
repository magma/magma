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

set(ITTI_LIB_DIR $ENV{C_BUILD}/oai/lib/itti)
# include_directories(${ITTI_LIB_DIR})

# include_directories("${PROJECT_SOURCE_DIR}")
include_directories("/usr/src/googletest/googlemock/include/")

add_library(ITTI_TEST_LIB

    )

link_directories(/usr/src/googletest/googlemock/lib/)

target_link_libraries(ITTI_TEST_LIB
    LIB_ITTI
    gmock_main gtest gtest_main gmock
    )

set_target_properties(ITTI_TEST_LIB PROPERTIES LINKER_LANGUAGE CXX)

add_executable(itti_test test_itti.cpp)
target_link_libraries(itti_test ITTI_TEST_LIB)
add_test(test_itti itti_test)
