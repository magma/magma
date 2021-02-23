/*
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#pragma once

//-----------------------------------------------------------------------------
/** @struct itti_n11_create_pdu_session_response_t
 *  @brief carries the PDU Session Establishment Response from SMF to AMF task
 */

typedef enum sm_session_fsm_state_e {
  CREATING,
  CREATE,
  ACTIVE,
  INACTIVE,
  RELEASED
} sm_session_fsm_state_t;

typedef enum pdu_session_type_e {
  IPV4,
  IPV6,
  IPV4IPV6,
  UNSTRUCTURED
} pdu_session_type_t;

typedef enum ssc_mode_e { SSC_MODE_1, SSC_MODE_2, SSC_MODE_3 } ssc_mode_t;

typedef enum m5g_sm_cause_e {
  M5GSM_OPERATOR_DETERMINED_BARRING                       = 0,
  M5GSM_INSUFFICIENT_RESOURCES                            = 1,
  M5GSM_MISSING_OR_UNKNOWN_DNN                            = 2,
  M5GSM_UNKNOWN_PDU_SESSION_TYPE                          = 3,
  M5GSM_USER_AUTHENTICATION_OR_AUTHORIZATION_FAILED       = 4,
  M5GSM_REQUEST_REJECTED_UNSPECIFIED                      = 5,
  M5GSM_SERVICE_OPTION_NOT_SUPPORTED                      = 6,
  M5GSM_REQUESTED_SERVICE_OPTION_NOT_SUBSCRIBED           = 7,
  M5GSM_SERVICE_OPTION_TEMPORARILY_OUT_OF_ORDER           = 8,
  M5GSM_REGULAR_DEACTIVATION                              = 10,
  M5GSM_NETWORK_FAILURE                                   = 11,
  M5GSM_REACTIVATION_REQUESTED                            = 12,
  M5GSM_INVALID_PDU_SESSION_IDENTITY                      = 13,
  M5GSM_SEMANTIC_ERRORS_IN_PACKET_FILTER                  = 14,
  M5GSM_SYNTACTICAL_ERROR_IN_PACKET_FILTER                = 15,
  M5GSM_OUT_OF_LADN_SERVICE_AREA                          = 16,
  M5GSM_PTI_MISMATCH                                      = 17,
  M5GSM_PDU_SESSION_TYPE_IPV4_ONLY_ALLOWED                = 18,
  M5GSM_PDU_SESSION_TYPE_IPV6_ONLY_ALLOWED                = 19,
  M5GSM_PDU_SESSION_DOES_NOT_EXIST                        = 20,
  M5GSM_INSUFFICIENT_RESOURCES_FOR_SPECIFIC_SLICE_AND_DNN = 21,
  M5GSM_NOT_SUPPORTED_SSC_MODE                            = 22,
  M5GSM_INSUFFICIENT_RESOURCES_FOR_SPECIFIC_SLICE         = 23,
  M5GSM_MISSING_OR_UNKNOWN_DNN_IN_A_SLICE                 = 24,
  M5GSM_INVALID_PTI_VALUE                                 = 25,
  M5GSM_MAXIMUM_DATA_RATE_PER_UE_FOR_USER_PLANE_INTEGRITY_PROTECTION_IS_TOO_LOW =
      26,
  M5GSM_SEMANTIC_ERROR_IN_THE_QOS_OPERATION                 = 27,
  M5GSM_SYNTACTICAL_ERROR_IN_THE_QOS_OPERATION              = 28,
  M5GSM_INVALID_MAPPED_EPS_BEARER_IDENTITY                  = 29,
  M5GSM_SEMANTICALLY_INCORRECT_MESSAGE                      = 30,
  M5GSM_INVALID_MANDATORY_INFORMATION                       = 31,
  M5GSM_MESSAGE_TYPE_NON_EXISTENT_OR_NOT_IMPLEMENTED        = 32,
  M5GSM_MESSAGE_TYPE_NOT_COMPATIBLE_WITH_THE_PROTOCOL_STATE = 33,
  M5GSM_IE_NON_EXISTENT_OR_NOT_IMPLEMENTED                  = 34,
  M5GSM_CONDITIONAL_IE_ERROR                                = 35,
  M5GSM_MESSAGE_NOT_COMPATIBLE_WITH_THE_PROTOCOL_STATE      = 36,
  M5GSM_PROTOCOL_ERROR_UNSPECIFIED                          = 37,
  M5GSM_PTI_ALREADY_IN_USE                                  = 38,
  M5GSM_OPERATION_SUCCESS                                   = 40
} m5g_sm_cause_t;

typedef enum redirect_address_type_e {
  IPV4_1,
  IPV6_1,
  URL,
  SIP_URI
} redirect_address_type_t;

typedef struct redirect_server_s {
  redirect_address_type_t redirect_address_type;
  uint8_t redirect_server_address[16];
} redirect_server_t;

typedef struct itti_n11_create_pdu_session_response_s {
  // common context
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  sm_session_fsm_state_t sm_session_fsm_state;
  uint32_t sm_session_version;
  // M5GSMSessionContextAccess
  uint8_t pdu_session_id;
  pdu_session_type_t pdu_session_type;
  ssc_mode_t selected_ssc_mode;
  m5g_sm_cause_t m5gsm_cause;
  bool always_on_pdu_session_indication;
  ssc_mode_t allowed_ssc_mode;
  bool m5gsm_congetion_re_attempt_indicator;
  redirect_server_t pdu_address;
} itti_n11_create_pdu_session_response_t;

#define N11_CREATE_PDU_SESSION_RESPONSE(mSGpTR)                                \
  (mSGpTR)->ittiMsg.n11_create_pdu_session_response
