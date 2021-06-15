# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)
PROJECT(PolicyDb)

set(CMAKE_CXX_STANDARD 14)
set(CMAKE_CXX_STANDARD_REQUIRED ON)
set(CMAKE_CXX_EXTENSIONS OFF)

#compile the relevant protos
include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)
create_proto_dir("lte" LTE_CPP_OUT_DIR)
create_proto_dir("orc8r" ORC8R_CPP_OUT_DIR)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(POLICYDB_ORC8R_CPP_PROTOS common)
generate_cpp_protos("${POLICYDB_ORC8R_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_CPP_OUT_DIR})

set(POLICYDB_LTE_CPP_PROTOS policydb apn mobilityd subscriberdb)
generate_cpp_protos("${POLICYDB_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_CPP_OUT_DIR})

message("Policydb proto_srcs are ${PROTO_SRCS}")

add_library(POLICYDB
    PolicyLoader.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS})

target_link_libraries(POLICYDB PRIVATE
    DATASTORE MAGMA_CONFIG glog)


target_compile_definitions(POLICYDB PUBLIC LOG_WITH_GLOG)
target_include_directories(POLICYDB PUBLIC
    "$<BUILD_INTERFACE:${CMAKE_CURRENT_SOURCE_DIR}>"
    $<TARGET_FILE_DIR:POLICYDB>)
