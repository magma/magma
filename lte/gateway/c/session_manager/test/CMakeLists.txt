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

set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

include_directories("${PROJECT_SOURCE_DIR}")
include_directories("/usr/src/googletest/googlemock/include/")

add_library(SESSIOND_TEST_LIB
    Consts.h
    Matchers.h
    ProtobufCreators.cpp
    ProtobufCreators.h
    SessiondMocks.h
    SessionStateTester.h
    )

link_directories(/usr/src/googletest/googlemock/lib/)
target_link_libraries(SESSIOND_TEST_LIB SESSION_MANAGER gmock_main gtest gtest_main gmock pthread rt)

foreach (session_test session_credit local_enforcer cloud_reporter
    session_manager_handler sessiond_integ session_state
    session_store store_client stored_state proxy_responder_handler
    metering_reporter local_enforcer_wallet_exhaust charging_grant
    usage_monitor upf_node_state)
  add_executable(${session_test}_test test_${session_test}.cpp)
  target_link_libraries(${session_test}_test SESSIOND_TEST_LIB)
  add_test(test_${session_test} ${session_test}_test)
endforeach (session_test)
