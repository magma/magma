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

include_directories("/usr/src/googletest/googlemock/include/")
link_directories(/usr/src/googletest/googlemock/lib/)

set(MME_APP_UE_CONTEXT_IMSI_SRC
    test_mme_app_ue_context.c
    )

set(MME_APP_NAS_ENCODE_DECODE_SRC
test_mme_app_nas_encode_decode.cpp
    )

set(NAS_EPS_QUALITY_OF_SERVICE
    test_EPSQualityOfService.cpp
    )

set(NAS_APN_AMBR
    test_ApnAggregateMaximumBitRate.cpp
    )

set(MME_APP_TEST_SRC
    mme_app_test.cpp
    mme_app_test_util.cpp
    test_mme_app_emm_encode_decode.cpp
    test_mme_app_esm_encode_decode.cpp
    test_mme_app_config.cpp
    test_mme_procedures.cpp
    test_sgw_config.cpp
    )

add_executable(test_mme_app_ue_context_imsi ${MME_APP_UE_CONTEXT_IMSI_SRC})
add_executable(test_nas_encode_decode ${MME_APP_NAS_ENCODE_DECODE_SRC})
add_executable(test_nas_eps_quality_of_service ${NAS_EPS_QUALITY_OF_SERVICE})
add_executable(test_nas_apn_ambr ${NAS_APN_AMBR})
add_executable(mme_app_test ${MME_APP_TEST_SRC})

target_link_libraries(test_mme_app_ue_context_imsi
    TASK_MME_APP ${CHECK_LIBRARIES} ${CMAKE_THREAD_LIBS_INIT}
    LIB_BSTR LIB_HASHTABLE
    )

target_link_libraries(test_nas_encode_decode
    TASK_MME_APP TASK_NAS ${CHECK_LIBRARIES} ${CMAKE_THREAD_LIBS_INIT}
    LIB_BSTR gtest gtest_main
    )

target_link_libraries(test_nas_apn_ambr
    TASK_MME_APP ${CHECK_LIBRARIES} ${CMAKE_THREAD_LIBS_INIT}
    LIB_BSTR LIB_HASHTABLE gtest gtest_main
    )

target_link_libraries(test_nas_eps_quality_of_service
    TASK_MME_APP ${CHECK_LIBRARIES} ${CMAKE_THREAD_LIBS_INIT}
    LIB_BSTR LIB_HASHTABLE gtest gtest_main
    )

target_link_libraries(mme_app_test
    TASK_MME_APP TASK_NAS TASK_AMF_APP ${CHECK_LIBRARIES} ${CMAKE_THREAD_LIBS_INIT}
    LIB_BSTR LIB_ITTI MOCK_TASKS gtest gtest_main ${CRYPTO_LIBRARIES} ${OPENSSL_LIBRARIES}
    ${NETTLE_LIBRARIES}
    )

target_include_directories(test_mme_app_ue_context_imsi PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CHECK_INCLUDE_DIRS}
    )
target_include_directories(test_nas_encode_decode PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CHECK_INCLUDE_DIRS}
    )

target_include_directories(test_nas_eps_quality_of_service PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CHECK_INCLUDE_DIRS}
    )

target_include_directories(test_nas_apn_ambr PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CHECK_INCLUDE_DIRS}
    )

add_test(NAME test_mme_app_ue_context COMMAND test_mme_app_ue_context_imsi)
add_test(NAME test_mme_app_nas_encode_decode COMMAND test_nas_encode_decode)
add_test(NAME test_nas_eps_quality_of_service COMMAND test_nas_eps_quality_of_service)
add_test(NAME test_nas_apn_ambr COMMAND test_nas_apn_ambr)
add_test(NAME test_mme_app COMMAND mme_app_test)
