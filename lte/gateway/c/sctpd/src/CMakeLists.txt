# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

cmake_minimum_required(VERSION 3.7)
set(CMAKE_CXX_STANDARD 14)

project(MagmaSctpd)

message("Adding code coverage build and linker flags for BUILD_TESTS")
set(CMAKE_CXX_FLAGS_DEBUG
    "${CMAKE_CXX_FLAGS_DEBUG} -g -O0 -fprofile-arcs -ftest-coverage")
set(CMAKE_LINKER_FLAGS_DEBUG
    "${CMAKE_LINKER_FLAGS_DEBUG} -fprofile-arcs -ftest-coverage")

# generate sctpd protos
include($ENV{MAGMA_ROOT}/orc8r/gateway/c/common/CMakeProtoMacros.txt)

create_proto_dir("lte" LTE_CPP_OUT_DIR)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(SCTPD_LTE_CPP_PROTOS sctpd)
generate_cpp_protos("${SCTPD_LTE_CPP_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_CPP_OUT_DIR})

set(SCTPD_GRPC_PROTOS sctpd)
generate_grpc_protos("${SCTPD_GRPC_PROTOS}" "${PROTO_SRCS}"
    "${PROTO_HDRS}" ${LTE_PROTO_DIR} ${LTE_CPP_OUT_DIR})

find_package(MAGMA_LOGGING REQUIRED)

add_library(SCTPD_LIB
    sctp_assoc.cpp
    sctp_connection.cpp
    sctp_desc.cpp
    sctpd_downlink_impl.cpp
    sctpd_event_handler.cpp
    sctpd_uplink_client.cpp
    util.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
    )

target_compile_definitions(SCTPD_LIB PUBLIC LOG_WITH_GLOG)

target_link_libraries(SCTPD_LIB
    sctp pthread grpc++ grpc protobuf glog MAGMA_LOGGING
    )

target_include_directories(SCTPD_LIB PUBLIC
    ${CMAKE_CURRENT_SOURCE_DIR}
    $<TARGET_FILE_DIR:SCTPD_LIB>
    )

# add sctpd executable
add_executable(sctpd
    sctpd.cpp
    )

target_link_libraries(sctpd
    SCTPD_LIB
    )

if (BUILD_TESTS)
  ENABLE_TESTING()
  ADD_SUBDIRECTORY(test)
endif (BUILD_TESTS)