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
set(MME_APP_TEST_SRC
    mme_app_test.cpp
    test_mme_app_emm_decode.cpp
    test_mme_app_config.cpp
    test_mme_procedures.cpp
    )

add_executable(test_mme_app_ue_context_imsi ${MME_APP_UE_CONTEXT_IMSI_SRC})
add_executable(mme_app_test ${MME_APP_TEST_SRC})

target_link_libraries(test_mme_app_ue_context_imsi
    TASK_MME_APP ${CHECK_LIBRARIES} ${CMAKE_THREAD_LIBS_INIT}
    LIB_BSTR LIB_HASHTABLE
    )
target_link_libraries(mme_app_test
    TASK_MME_APP TASK_NAS ${CHECK_LIBRARIES} ${CMAKE_THREAD_LIBS_INIT}
    LIB_BSTR LIB_ITTI MOCK_TASKS gtest gtest_main ${CRYPTO_LIBRARIES} ${OPENSSL_LIBRARIES}
    ${NETTLE_LIBRARIES}
    )

add_test(NAME test_mme_app_ue_context COMMAND test_mme_app_ue_context_imsi)
add_test(NAME test_mme_app COMMAND mme_app_test)
