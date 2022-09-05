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

# Needed to determine if build with CMAKE or Bazel
set(CMAKE_CXX_FLAGS "${CMAKE_C_FLAGS} -DCMAKE_BUILD")

set(SPGW_LIB_DIR $ENV{C_BUILD}/oai/tasks/sgw)

include_directories("/usr/src/googletest/googlemock/include/")

add_library(SPGW_TASK_TEST_LIB
    state_creators.hpp
    state_creators.cpp
    spgw_test_util.h
    spgw_test_util.cpp
    mock_spgw_op.hpp
    mock_spgw_op.cpp
    spgw_test.cpp
    spgw_procedures_test_fixture.hpp
    spgw_procedures_test_fixture.cpp
    )

link_directories(/usr/src/googletest/googlemock/lib/)

target_link_libraries(SPGW_TASK_TEST_LIB
    COMMON LIB_BSTR
    TASK_SGW MOCK_TASKS
    gmock_main gtest gtest_main gmock
    pthread rt yaml-cpp
    TASK_GRPC_SERVICE
    )

set_target_properties(SPGW_TASK_TEST_LIB PROPERTIES LINKER_LANGUAGE CXX)

foreach (sgw_test spgw_service_impl spgw_state_converter 
         spgw_procedures spgw_procedures_dedicated_bearer spgw_procedures_session
         pgw_pco spgw_procedures_with_injected_state)
  add_executable(${sgw_test}_test test_${sgw_test}.cpp)
  target_link_libraries(${sgw_test}_test SPGW_TASK_TEST_LIB TASK_GRPC_SERVICE)
  add_test(test_${sgw_test} ${sgw_test}_test)
endforeach (sgw_test)
