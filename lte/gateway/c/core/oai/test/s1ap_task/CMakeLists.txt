# Copyright 2021 The Magma Authors.
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

include_directories(${PROJECT_SOURCE_DIR})

add_executable(s1ap_test
        s1ap_mme_test_utils.cpp
        mock_s1ap_op.cpp
        s1ap_test.cpp
        test_s1ap_mme_handlers.cpp
        test_s1ap_mme_handlers_with_injected_state.cpp
        test_s1ap_handle_new_association.cpp
        test_s1ap_state_manager.cpp)

target_link_libraries(s1ap_test
        TASK_S1AP MOCK_TASKS TASK_AMF_APP
        gtest gtest_main
        )

target_include_directories(s1ap_test PUBLIC
        ${CMAKE_CURRENT_SOURCE_DIR}
        ${CHECK_INCLUDE_DIRS}
        )

add_test(test_s1ap s1ap_test)
