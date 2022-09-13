# Copyright 2022 The Magma Authors.
# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7.2)
set(CMAKE_CXX_FLAGS "-Wno-write-strings -Wno-literal-suffix")

#compile the required protos
set(SUB_LTE_PROTOS subscriberdb apn)
set(SUB_ORC8R_PROTOS digest)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

create_proto_dir("lte" LTE_OUT_DIR)
create_proto_dir("orc8r" ORC8R_OUT_DIR)

generate_cpp_protos("${SUB_LTE_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_OUT_DIR})

generate_cpp_protos("${SUB_ORC8R_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${ORC8R_OUT_DIR})

include_directories(${CMAKE_CURRENT_BINARY_DIR})
include_directories(${LTE_OUT_DIR})
include_directories(${ORC8R_OUT_DIR})

add_library(LIB_STORE sqlite.cpp ${PROTO_SRCS} ${PROTO_HDRS})
target_link_libraries(LIB_STORE MAGMA_CONFIG COMMON  protobuf sqlite3)

target_include_directories(LIB_STORE PUBLIC
    $ENV{MAGMA_ROOT}
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${MAGMA_COMMON_DIR}
    )
