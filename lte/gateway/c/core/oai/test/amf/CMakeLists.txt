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

# Dependencies
find_package(Check REQUIRED)
find_package(Threads REQUIRED)
pkg_search_module(OPENSSL openssl REQUIRED)
include_directories(${OPENSSL_INCLUDE_DIRS})

pkg_search_module(CRYPTO libcrypto REQUIRED)
include_directories(${CRYPTO_INCLUDE_DIRS})

pkg_search_module(NETTLE nettle REQUIRED)
include_directories(${NETTLE_INCLUDE_DIRS})


include_directories(${PROJECT_SOURCE_DIR})
include_directories("${PROJECT_SOURCE_DIR}/include/nas")

include_directories("/usr/src/googletest/googlemock/include/")
link_directories(/usr/src/googletest/googlemock/lib/)

set(AMF_APP_TEST_SRC
    amf_app_test.cpp
    util_nas5g_registration_pkt.cpp
    util_nas5g_auth_fail_pkt.cpp
    util_nas5g_ul_nas_pdu_decode.cpp
    util_nas5g_service_request_pkt.cpp
    util_nas5g_security_mode_reject_pkt.cpp
    test_amf_encode_decode.cpp
    amf_app_test_util.cpp
    test_amf_procedures.cpp
    test_amf_algorithm_selection.cpp
    test_amf_common_utils.cpp
    util_s6a_update_location.cpp
    util_s6a_update_location.h
    test_amf_map.cpp
    test_amf_stateless.cpp
    )

add_executable(amf_app_test ${AMF_APP_TEST_SRC})

target_link_libraries(amf_app_test
    TASK_AMF_APP LIB_N11 LIB_NAS5G LIB_SECU ${CHECK_LIBRARIES} ${CMAKE_THREAD_LIBS_INIT}
    LIB_BSTR LIB_ITTI MOCK_TASKS gtest gtest_main crypt ${CRYPTO_LIBRARIES} ${OPENSSL_LIBRARIES}
    ${NETTLE_LIBRARIES}
    )

add_test(NAME test_amf_app COMMAND amf_app_test)
