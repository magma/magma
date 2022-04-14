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

create_proto_dir("orc8r" ORC8R_OUT_DIR)
# Common protos
set(ORC8R_CPP_PROTOS common )
generate_cpp_protos("${ORC8R_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_OUT_DIR})
include_directories("${PROJECT_SOURCE_DIR}")
include_directories("${PROJECT_SOURCE_DIR}/tasks/s6a")

include_directories("/usr/src/googletest/googlemock/include/")
link_directories(/usr/src/googletest/googlemock/lib/)

set(S6A_TEST_CONFIG_SRC
    s6a_recv_async_grpc_messages.cpp
    )

add_executable(s6a_test ${S6A_TEST_CONFIG_SRC})

target_link_libraries(s6a_test
    ${CHECK_LIBRARIES} ${CMAKE_THREAD_LIBS_INIT}
    TASK_ASYNC_GRPC_SERVICE LIB_BSTR LIB_ITTI MOCK_TASKS
    gtest gtest_main
    )

target_include_directories(s6a_test PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CHECK_INCLUDE_DIRS}
    )

add_test(NAME test_s6a COMMAND s6a_test)
