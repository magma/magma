# Copyright 2020 The Magma Authors.
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

set(LIB_GTPV2C $ENV{C_BUILD}/oai/lib/gtpv2-c)
include_directories(${LIB_GTPV2C})
include_directories("${PROJECT_SOURCE_DIR}")
include_directories("${PROJECT_SOURCE_DIR}/lib/gtpv2-c/nwgtpv2c-0.11/include")
include_directories("${PROJECT_SOURCE_DIR}/lib/gtpv2-c/nwgtpv2c-0.11/shared")
include_directories(NWGTPV2C_IE_FORMATTER_DIR)

add_executable(gtpv2c_test test_fteid.cpp)
target_link_libraries(gtpv2c_test ${CONFIG} LIB_GTPV2C
        COMMON gtest gtest_main pthread rt yaml-cpp)
add_test(test_gtpv2c gtpv2c_test)
