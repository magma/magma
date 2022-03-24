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

set(TEST_SHARED_SRCS 
        s1ap_mme_test_utils.cpp
        s1ap_test.cpp
        s1ap_mme_handlers_test_fixture.cpp
)
set(TEST_SHARED_LINKS
        TASK_S1AP MOCK_TASKS TASK_AMF_APP
        gtest gtest_main
)

add_executable(s1ap_mme_handler_ue_test 
        test_s1ap_mme_handlers_ue.cpp
        ${TEST_SHARED_SRCS})
target_link_libraries(s1ap_mme_handler_ue_test ${TEST_SHARED_LINKS})
add_test(test_s1ap_mme_handler_ue s1ap_mme_handler_ue_test)

add_executable(s1ap_mme_handlers_handover_test 
        test_s1ap_mme_handlers_handover.cpp
        ${TEST_SHARED_SRCS})
target_link_libraries(s1ap_mme_handlers_handover_test ${TEST_SHARED_LINKS})
add_test(test_s1ap_mme_handlers_handover s1ap_mme_handlers_handover_test)

add_executable(s1ap_mme_handlers_erab_test 
        test_s1ap_mme_handlers_erab.cpp
        ${TEST_SHARED_SRCS})
target_link_libraries(s1ap_mme_handlers_erab_test ${TEST_SHARED_LINKS})
add_test(test_s1ap_mme_handlers_erab s1ap_mme_handlers_erab_test)

add_executable(s1ap_mme_handlers_s1_enb_test 
        test_s1ap_mme_handlers_s1_enb.cpp
        ${TEST_SHARED_SRCS})
target_link_libraries(s1ap_mme_handlers_s1_enb_test ${TEST_SHARED_LINKS})
add_test(test_s1ap_mme_handlers_s1_enb s1ap_mme_handlers_s1_enb_test)

add_executable(s1ap_handle_new_association_test
        test_s1ap_handle_new_association.cpp
        ${TEST_SHARED_SRCS})
target_link_libraries(s1ap_handle_new_association_test ${TEST_SHARED_LINKS})
add_test(test_s1ap_handle_new_association s1ap_handle_new_association_test)

add_executable(s1ap_state_manager_test
        test_s1ap_state_manager.cpp
        ${TEST_SHARED_SRCS})
target_link_libraries(s1ap_state_manager_test ${TEST_SHARED_LINKS})
add_test(test_s1ap_state_manager s1ap_state_manager_test)

add_executable(s1ap_state_converter_test
        test_s1ap_state_converter.cpp
        ${TEST_SHARED_SRCS})
target_link_libraries(s1ap_state_converter_test ${TEST_SHARED_LINKS})
add_test(test_s1ap_state_converter s1ap_state_converter_test)


