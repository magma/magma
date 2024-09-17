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
add_executable(service_registry_test test_service_registry.cpp)

target_link_libraries(service_registry_test
    gmock_main gtest gtest_main gmock
    ${GCOV_LIB} SERVICE_REGISTRY pthread)
add_test(test_service_registry service_registry_test)
