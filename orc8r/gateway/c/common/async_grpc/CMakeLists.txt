# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.


include_directories("${PROJECT_SOURCE_DIR}/../common/logging")

#compile the relevant protos

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(ASYNC_ORC8R_CPP_PROTOS common redis eventd)
set(ASYNC_LTE_CPP_PROTOS session_manager abort_session policydb subscriberdb mobilityd subscriberauth)
set(ASYNC_LTE_GRPC_PROTOS session_manager abort_session mobilityd subscriberauth)
set(ASYNC_ORC8R_GRPC_PROTOS eventd)

generate_all_protos("${ASYNC_LTE_CPP_PROTOS}" "${ASYNC_ORC8R_CPP_PROTOS}"
  "${ASYNC_LTE_GRPC_PROTOS}"
  "${ASYNC_ORC8R_GRPC_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}")

message("Async Proto_srcs are ${PROTO_SRCS}")

add_library(ASYNC_GRPC
    GRPCReceiver.cpp
    ${PROTO_SRCS}
    ${PROTO_HDRS}
)

# copy headers to build directory so they can be shared with OAI,
# session_manager, etc.
add_custom_command(TARGET ASYNC_GRPC POST_BUILD
                   COMMAND ${CMAKE_COMMAND} -E copy
                   ${PROJECT_SOURCE_DIR}/async_grpc/*.h
                   $<TARGET_FILE_DIR:ASYNC_GRPC>)

target_include_directories(ASYNC_GRPC PUBLIC
                   $<TARGET_FILE_DIR:ASYNC_GRPC>
)
