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

find_package(Check REQUIRED)
find_package(Threads REQUIRED)
find_package(libprotobuf-mutator CONFIG REQUIRED)

set(MME_APP_UE_CONTEXT_IMSI_SRC
    test_mme_app_ue_context.c
    )
set(MME_APP_EMM_DECODE_SRC
    test_mme_app_emm_decode.cpp
    )

add_executable(test_mme_app_ue_context_imsi ${MME_APP_UE_CONTEXT_IMSI_SRC})
add_executable(test_mme_app_emm_decode ${MME_APP_EMM_DECODE_SRC})

target_link_libraries(test_mme_app_ue_context_imsi
    TASK_MME_APP ${CHECK_LIBRARIES} ${CMAKE_THREAD_LIBS_INIT}
    LIB_BSTR LIB_HASHTABLE
    )
target_link_libraries(test_mme_app_emm_decode
    TASK_MME_APP TASK_NAS ${CHECK_LIBRARIES} ${CMAKE_THREAD_LIBS_INIT}
    LIB_BSTR gtest gtest_main
    )

target_include_directories(test_mme_app_ue_context_imsi PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CHECK_INCLUDE_DIRS}
    )
target_include_directories(test_mme_app_emm_decode PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CHECK_INCLUDE_DIRS}
    )

add_test(NAME test_mme_app_ue_context COMMAND test_mme_app_ue_context_imsi)
add_test(NAME test_mme_app_emm_decode COMMAND test_mme_app_emm_decode)

add_executable(mme_app_state_converter_test mme_app_state_converter_test.cpp)

pkg_search_module(OPENSSL openssl REQUIRED)
pkg_search_module(CRYPTO libcrypto REQUIRED)
target_link_libraries(mme_app_state_converter_test
    -Wl,--start-group
        TASK_MME_APP COMMON
    -Wl,--end-group
        ${CRYPTO_LIBRARIES} ${OPENSSL_LIBRARIES}
        gtest gtest_main
        libprotobuf-mutator::protobuf-mutator libprotobuf-mutator::protobuf-mutator-libfuzzer
    )
target_include_directories(mme_app_state_converter_test PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    $ENV{C_BUILD}/oai/tasks/mme_app
    )
add_test(NAME mme_app_state_converter_test COMMAND mme_app_state_converter_test) 
