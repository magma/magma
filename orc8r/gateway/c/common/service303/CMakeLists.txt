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

#compile the relevant protos
create_proto_dir("orc8r" CPP_OUT_DIR)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(SERVICE303_PROTOS common metricsd service303)
generate_cpp_protos("${SERVICE303_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}"
  ${ORC8R_PROTO_DIR} ${CPP_OUT_DIR})

set(SERVICE303_GRPC_PROTOS service303)
generate_grpc_protos("${SERVICE303_GRPC_PROTOS}" "${PROTO_SRCS}"
  "${PROTO_HDRS}" ${ORC8R_PROTO_DIR} ${CPP_OUT_DIR})

generate_prometheus_proto("${PROTO_SRCS}" "${PROTO_HDRS}" ${CPP_OUT_DIR})

message("Proto_srcs are ${PROTO_SRCS}")

add_library(SERVICE303_LIB
  MagmaService.cpp
  MetricsSingleton.cpp
  MetricsHelpers.cpp
  ProcFileUtils.cpp
  ${PROTO_SRCS}
  ${PROTO_HDRS}
)

target_link_libraries(SERVICE303_LIB
  prometheus-cpp protobuf grpc grpc++
  SERVICE_REGISTRY LOGGING
)

# copy headers to build directory so they can be shared with OAI,
# session_manager, etc.
add_custom_command(TARGET SERVICE303_LIB POST_BUILD
                   COMMAND ${CMAKE_COMMAND} -E copy
                   ${PROJECT_SOURCE_DIR}/service303/*.h
                   $<TARGET_FILE_DIR:SERVICE303_LIB>)

target_include_directories(SERVICE303_LIB PUBLIC
  $<TARGET_FILE_DIR:SERVICE303_LIB>
  "/usr/local/include/prometheus"
)
