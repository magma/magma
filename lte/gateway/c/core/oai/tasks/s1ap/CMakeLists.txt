# Copyright 2020 The Magma Authors.

# This source code is licensed under the BSD-style license found in the
# LICENSE file in the root directory of this source tree.

# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

add_list1_option(S1AP_VERSION R15 "S1AP Asn.1 grammar version" R15)
add_definitions(-DASN1_MINIMUM_VERSION=923)
set(ASN1RELDIR r15)
set(S1AP_DIR ${CMAKE_CURRENT_SOURCE_DIR})
set(asn1_generated_dir ${CMAKE_CURRENT_BINARY_DIR})
set(S1AP_ASN_DIR ${CMAKE_CURRENT_SOURCE_DIR}/messages/asn1/${ASN1RELDIR})
set(S1AP_ASN_FILES "s1ap-15.6.0.asn1")
set(S1AP_C_DIR ${CMAKE_CURRENT_BINARY_DIR}/${ASN1RELDIR})
include_directories(${S1AP_C_DIR})
include_directories(${S1AP_DIR})
set(GENERATED_FULL_DIR ${S1AP_C_DIR})
set(ASN1_SOURCE_DIR ${S1AP_ASN_DIR}/${S1AP_ASN_FILES})
set(ENV{ASN1C_PREFIX} "S1ap_")

file(MAKE_DIRECTORY ${GENERATED_FULL_DIR})
set(s1ap_generate_code_done_flag ${GENERATED_FULL_DIR}/S1AP_GENERATE_CODE_DONE)

if (${ASN1_SOURCE_DIR} IS_NEWER_THAN ${s1ap_generate_code_done_flag})
  file(REMOVE ${GENERATED_FULL_DIR}/${ASN1C_PREFIX}*.c ${GENERATED_FULL_DIR}/${ASN1C_PREFIX}*.h)
  execute_process(COMMAND asn1c -pdu=all -fcompound-names -gen-PER -no-gen-example -fno-include-deps -D ${GENERATED_FULL_DIR} ${ASN1_SOURCE_DIR} RESULT_VARIABLE ret)
  if (NOT ${ret} STREQUAL 0)
    message(FATAL_ERROR "${ret}")
  endif (NOT ${ret} STREQUAL 0)
  execute_process(COMMAND bash "-c" "egrep -lRZ \"18446744073709551615\" ${GENERATED_FULL_DIR} | xargs -0 -l sed -i -e \"s/18446744073709551615/18446744073709551615u/g\"")
endif ()
# TOUCH not in cmake 3.10
file(WRITE ${s1ap_generate_code_done_flag})

file(GLOB S1AP_source ${S1AP_C_DIR}/*.c)
list(REMOVE_ITEM S1AP_source ${S1AP_C_DIR}/converter-sample.c)

file(GLOB s1ap_h ${S1AP_C_DIR}/*.c)
set(s1ap_h ${s1ap_h})

add_library(LIB_S1AP
    ${S1AP_source}
    )
target_link_libraries(LIB_S1AP
    LIB_BSTR LIB_HASHTABLE
    )
target_include_directories(LIB_S1AP PUBLIC
    ${S1AP_C_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}/r15
    )

set(S1AP_STATE_CPP_PROTOS s1ap_state common_types)

list(APPEND PROTO_SRCS "")
list(APPEND PROTO_HDRS "")

set(STATE_OUT_DIR "${CMAKE_CURRENT_BINARY_DIR}/lte/protos/oai")
generate_cpp_protos("${S1AP_STATE_CPP_PROTOS}" "${PROTO_SRCS}" "${PROTO_HDRS}"
    "${STATE_PROTO_DIR}" "${STATE_OUT_DIR}")

add_library(TASK_S1AP
    "${PROTO_SRCS}"
    "${PROTO_HDRS}"
    ${S1AP_DIR}/s1ap_mme_encoder.c
    ${S1AP_DIR}/s1ap_mme_decoder.c
    ${S1AP_DIR}/s1ap_mme_handlers.c
    ${S1AP_DIR}/s1ap_mme_nas_procedures.c
    ${S1AP_DIR}/s1ap_mme.c
    ${S1AP_DIR}/s1ap_mme_itti_messaging.c
    ${S1AP_DIR}/s1ap_mme_ta.c
    ${S1AP_DIR}/s1ap_state.cpp
    ${S1AP_DIR}/s1ap_state_manager.cpp
    ${S1AP_DIR}/s1ap_state_converter.cpp
    )
target_link_libraries(TASK_S1AP
    MAGMA_CONFIG
    COMMON
    LIB_BSTR LIB_HASHTABLE
    TASK_SERVICE303 TASK_MME_APP
    cpp_redis tacopie
    )
target_include_directories(TASK_S1AP PUBLIC
    ${CMAKE_CURRENT_BINARY_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}
    ${CMAKE_CURRENT_SOURCE_DIR}/messages
    ${CMAKE_CURRENT_SOURCE_DIR}/messages/asn1
    )
