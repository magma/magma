# Copyright 2020 The Magma Authors.
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)

set(CMAKE_CXX_STANDARD 11)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

include_directories("${PROJECT_SOURCE_DIR}")
include_directories("${PROJECT_SOURCE_DIR}/tasks/amf")

pkg_search_module(OPENSSL openssl REQUIRED)
include_directories(${OPENSSL_INCLUDE_DIRS})

pkg_search_module(CRYPTO libcrypto REQUIRED)
include_directories(${CRYPTO_INCLUDE_DIRS})

include_directories("${PROJECT_SOURCE_DIR}/include/nas")

add_library(AMF_TASK_TEST_LIB
    util_nas5g_pkt.h
    util_nas5g_registration_pkt.cpp
    util_nas5g_auth_fail_pkt.cpp
    util_nas5g_ul_nas_pdu_decode.cpp
    util_nas5g_service_request_pkt.cpp
    )
link_directories(${PROJECT_SOURCE_DIR}/tasks/amf)
target_link_libraries(AMF_TASK_TEST_LIB
        ${CONFIG}
        COMMON
	LIB_NAS5G TASK_AMF_APP gtest gtest_main pthread rt yaml-cpp
        ${CRYPTO_LIBRARIES} ${OPENSSL_LIBRARIES}
        ${NETTLE_LIBRARIES}
    )

add_executable(amf_test test_amf.cpp)
target_link_libraries(amf_test AMF_TASK_TEST_LIB TASK_AMF_APP)
add_test(test_amf amf_test)
