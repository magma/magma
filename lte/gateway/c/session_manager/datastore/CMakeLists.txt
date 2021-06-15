# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

PROJECT(datastore)
cmake_minimum_required(VERSION 3.7.2)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(DS_ORC8R_CPP_PROTOS common redis)
set(DS_LTE_CPP_PROTOS "")
set(DS_LTE_GRPC_PROTOS "")
set(DS_ORC8R_GRPC_PROTOS "")

generate_all_protos("${DS_LTE_CPP_PROTOS}" "${DS_ORC8R_CPP_PROTOS}"
    "${DS_LTE_GRPC_PROTOS}"
    "${DS_ORC8R_GRPC_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}")

message("Datastore Proto_srcs are ${PROTO_SRCS}")

add_library(DATASTORE
    ObjectMap.h
    Serializers.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_link_libraries(DATASTORE PRIVATE
    tacopie cpp_redis pthread
    )

target_include_directories(DATASTORE PUBLIC
    $<TARGET_FILE_DIR:DATASTORE> # protos
    "$<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}>")
