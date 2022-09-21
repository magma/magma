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

include_directories("/usr/src/googletest/googlemock/include/")

link_directories(/usr/src/googletest/googlemock/lib/)

add_executable(bstr_test test_bstr.cpp)
target_link_libraries(bstr_test LIB_BSTR gmock_main gtest gtest_main gmock pthread)
add_test(test_bstr bstr_test)

add_executable(3gpp_test test_3gpp.cpp)
target_link_libraries(3gpp_test LIB_3GPP gmock_main gtest gtest_main gmock)
add_test(test_3gpp 3gpp_test)

add_executable(ula_sub_data_test test_ula_subData.cpp)
target_link_libraries(ula_sub_data_test LIB_STORE LIB_S6A_PROXY COMMON TASK_MME_APP gmock_main gtest gtest_main gmock)
add_test(test_ula_subdata ula_sub_data_test)
