# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)
PROJECT(MagmaCommonTests)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

include_directories("/usr/src/googletest/googlemock/include/")
link_directories("/usr/src/googletest/googlemock/lib/")

add_executable(yaml_utils_test test_yaml_utils.cpp)
target_link_libraries(yaml_utils_test
    MAGMA_CONFIG yaml-cpp rt
    gmock_main gtest gtest_main gmock
    pthread
    ${GCOV_LIB})
add_test(test_yaml_utils yaml_utils_test)