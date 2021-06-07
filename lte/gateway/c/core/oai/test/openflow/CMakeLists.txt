# Copyright 2020 The Magma Authors.
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

add_compile_options(-std=c++11)

set(OUTPUT_DIR "${CMAKE_CURRENT_BINARY_DIR}")

include_directories("${PROJECT_SOURCE_DIR}/openflow/controller")
include_directories("/usr/src/googletest/googlemock/include/")
link_directories(/usr/src/googletest/googlemock/lib/)
add_executable(openflow_controller_test test_openflow_controller.cpp)
add_executable(imsi_encoder_test test_imsi_encoder.cpp)
add_executable(gtp_app_test test_gtp_app.cpp)

add_library(OPENFLOW_TEST openflow_mocks.h)
target_link_libraries(OPENFLOW_TEST
    COMMON
    lfds710
    LIB_OPENFLOW_CONTROLLER LIB_BSTR LIB_HASHTABLE LIB_ITTI LIB_S1AP TASK_S1AP
    gmock_main gtest gtest_main gmock pthread rt)

set_target_properties(OPENFLOW_TEST PROPERTIES LINKER_LANGUAGE CXX)

target_link_libraries(openflow_controller_test OPENFLOW_TEST)
target_link_libraries(imsi_encoder_test OPENFLOW_TEST)
target_link_libraries(gtp_app_test OPENFLOW_TEST)

add_test(test_openflow_controller openflow_controller_test)
add_test(test_imsi_encoder imsi_encoder_test)
add_test(test_gtp_app gtp_app_test)
