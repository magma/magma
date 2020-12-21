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

#ifndef N11_MESSAGES_TYPES_SEEN
#define N11_MESSAGES_TYPES_SEEN

//-----------------------------------------------------------------------------
/** @struct itti_n11_create_pdu_session_response_t
 *  @brief Create PDU Session Response */

typedef enum SMSessionFSMState_response_s {
  CREATING_0,
  CREATE_1,
  ACTIVE_2,
  INACTIVE_3,
  RELEASED_4
} SMSessionFSMState_response;

typedef enum PduSessionType_response_s {
  IPV4,
  IPV6,
  IPV4IPV6,
  UNSTRUCTURED
} PduSessionType_response;

typedef enum SscMode_response_s {
  SSC_MODE_1,
  SSC_MODE_2,
  SSC_MODE_3
} SscMode_response;

typedef enum M5GSMCause_response_s {
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
  M5GSM_INFORMATION_ELEMENT_NON_EXISTENT_OR_NOT_IMPLEMENTED = 34,
  M5GSM_CONDITIONAL_IE_ERROR                                = 35,
  M5GSM_MESSAGE_NOT_COMPATIBLE_WITH_THE_PROTOCOL_STATE      = 36,
  M5GSM_PROTOCOL_ERROR_UNSPECIFIED                          = 37,
  M5GSM_PTI_ALREADY_IN_USE                                  = 38,
  M5GSM_OPERATION_SUCCESS                                   = 40
} M5GSMCause_response;

typedef enum RedirectAddressType_response_s {
  IPV4_1,
  IPV6_1,
  URL,
  SIP_URI
} RedirectAddressType_response;

typedef struct RedirectServer_response_s {
  RedirectAddressType_response redirect_address_type;
  uint8_t redirect_server_address[16];
} RedirectServer_response;

typedef struct QosRules_response_s {
  uint32_t qos_rule_identifier;
  bool dqr;
  uint32_t number_of_packet_filters;
  uint32_t packet_filter_identifier[16];
  uint32_t qos_rule_precedence;
  bool segregation;
  uint32_t qos_flow_identifier;
} QosRules_response;

typedef struct itti_n11_create_pdu_session_response_s {
  // common context
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  SMSessionFSMState_response sm_session_fsm_state;
  uint32_t sm_session_version;
  // M5GSMSessionContextAccess
  uint8_t pdu_session_id[2];
  PduSessionType_response pdu_session_type;
  SscMode_response selected_ssc_mode;
  QosRules_response authorized_qos_rules[4];
  M5GSMCause_response M5gsm_cause;
  bool always_on_pdu_session_indication;
  SscMode_response allowed_ssc_mode;
  bool M5gsm_congetion_re_attempt_indicator;
  RedirectServer_response pdu_address;
} itti_n11_create_pdu_session_response_t;

#define N11_CREATE_PDU_SESSION_RESPONSE(mSGpTR)                                \
  (mSGpTR)->ittiMsg.n11_create_pdu_session_response

#endif
