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

add_library(LIAGENTD_TEST_LIB
    Consts.h
    LIAgentdMocks.h
    )

link_directories(/usr/src/googletest/googlemock/lib/)
target_link_libraries(LIAGENTD_TEST_LIB LI_AGENT gmock_main gtest gtest_main
    gmock rt)
target_include_directories(LIAGENTD_TEST_LIB PUBLIC
    ${PROJECT_SOURCE_DIR})

foreach (li_agent_test pdu_generator)
  add_executable(${li_agent_test}_test test_${li_agent_test}.cpp)
  target_link_libraries(${li_agent_test}_test LIAGENTD_TEST_LIB)
  add_test(test_${li_agent_test} ${li_agent_test}_test)
endforeach (li_agent_test)
