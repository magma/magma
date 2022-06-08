# Copyright 2022 The Magma Authors.
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

include_directories("/usr/src/googletest/googlemock/include/")

link_directories(/usr/src/googletest/googlemock/lib/)

set(MME_APP_MAP_TEST_SRC
    test_mme_app_map.cpp
    )
add_executable(mme_app_map_test ${MME_APP_MAP_TEST_SRC})
target_link_libraries(mme_app_map_test
    TASK_MME_APP ${CHECK_LIBRARIES} ${CMAKE_THREAD_LIBS_INIT}
    gtest gtest_main
    )

target_include_directories(mme_app_map_test PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CHECK_INCLUDE_DIRS}
    )

add_test(NAME test_mme_app_map COMMAND mme_app_map_test)
