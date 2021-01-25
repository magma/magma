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

include_directories("${PROJECT_SOURCE_DIR}/../common/logging")

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(ASYNC_ORC8R_CPP_PROTOS common redis)
set(ASYNC_LTE_CPP_PROTOS session_manager policydb mobilityd subscriberdb)
set(ASYNC_LTE_GRPC_PROTOS session_manager)
set(ASYNC_ORC8R_GRPC_PROTOS "")

generate_all_protos("${ASYNC_LTE_CPP_PROTOS}" "${ASYNC_ORC8R_CPP_PROTOS}"
        "${ASYNC_LTE_GRPC_PROTOS}"
        "${ASYNC_ORC8R_GRPC_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}")

message("Async Proto_srcs are ${PROTO_SRCS}")

add_library(DATASTORE
    RedisMap.hpp
    ObjectMap.h
    Serializers.cpp
    Serializers.h
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_link_libraries(DATASTORE
   tacopie cpp_redis pthread
   )

# copy headers to build directory so they can be shared with OAI,
# session_manager, etc.
add_custom_command(TARGET DATASTORE POST_BUILD
                   COMMAND ${CMAKE_COMMAND} -E copy
                   ${PROJECT_SOURCE_DIR}/datastore/*.h*
                   $<TARGET_FILE_DIR:DATASTORE>)

target_include_directories(DATASTORE PUBLIC
                  $<TARGET_FILE_DIR:DATASTORE>
)
