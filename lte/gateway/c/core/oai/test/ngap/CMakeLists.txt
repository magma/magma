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

set(LIB_NGAP $ENV{C_BUILD}/oai/tasks/ngap)
include_directories(${LIB_NGAP})
include_directories("${PROJECT_SOURCE_DIR}")
include_directories("${PROJECT_SOURCE_DIR}/tasks/ngap")

add_library(NGAP_TASK_TEST_LIB
    util_ngap_pkt.h
    util_ngap_setup_failure.cpp
    util_ngap_initiate_ue.cpp
    util_ngap_amf_nas_procedures.cpp
    )
link_directories(${PROJECT_SOURCE_DIR}/tasks/ngap)
target_link_libraries(NGAP_TASK_TEST_LIB
        ${CONFIG}
        COMMON
	LIB_NGAP gtest gtest_main pthread rt yaml-cpp
    )

add_executable(ngap_test test_ngap.cpp)
target_link_libraries(ngap_test NGAP_TASK_TEST_LIB TASK_NGAP)
add_test(test_ngap ngap_test)
