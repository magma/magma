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

#include "common_types.h"
//-----------------------------------------------------------------------------
/** @struct itti_n11_create_pdu_session_response_t
 *  @brief Create PDU Session Response */

/***********************pdu_res_set_change starts*************************/
typedef enum {
  SHALL_NOT_TRIGGER_PRE_EMPTION,
  MAY_TRIGGER_PRE_EMPTION,
} pre_emption_capability;

typedef enum sm_session_fsm_state_e {
  CREATING,
  CREATE,
  ACTIVE,
  INACTIVE,
  RELEASED
} sm_session_fsm_state_t;

typedef enum {
  NOT_PREEMPTABLE,
  PRE_EMPTABLE,
} pre_emption_vulnerability;

typedef struct m5g_allocation_and_retention_priority_s {
  int priority_level;
  pre_emption_capability pre_emption_cap;
  pre_emption_vulnerability pre_emption_vul;
} m5g_allocation_and_retention_priority;

typedef struct non_dynamic_5QI_descriptor_s {
  int fiveQI;
} non_dynamic_5QI_descriptor;

typedef struct qos_characteristics_s {
  non_dynamic_5QI_descriptor non_dynamic_5QI_desc;
} qos_characteristics_t;

typedef struct qos_flow_level_qos_parameters_s {
  qos_characteristics_t qos_characteristic;
  m5g_allocation_and_retention_priority alloc_reten_priority;
} qos_flow_level_qos_parameters;

typedef struct qos_flow_setup_request_item_s {
  uint32_t qos_flow_identifier;
  qos_flow_level_qos_parameters qos_flow_level_qos_param;
  // E-RAB ID is optional spec-38413 - 9.3.4.1
} qos_flow_setup_request_item;

typedef struct qos_flow_request_list_s {
  qos_flow_setup_request_item qos_flow_req_item;
} qos_flow_request_list_t;

typedef struct amf_pdn_type_value_s {
  pdn_type_value_t pdn_type;
} amf_pdn_type_value_t;

typedef struct gtp_tunnel_s {
  bstring endpoint_ip_address;  // Transport_Layer_Information
  uint8_t gtp_tied[4];
} gtp_tunnel;

typedef struct up_transport_layer_information_s {
  gtp_tunnel gtp_tnl;
} up_transport_layer_information_t;

typedef struct amf_ue_aggregate_maximum_bit_rate_s {
  uint64_t dl;
  uint64_t ul;
} amf_ue_aggregate_maximum_bit_rate_t;

typedef struct pdu_session_resource_setup_request_transfer_s {
  amf_ue_aggregate_maximum_bit_rate_t pdu_aggregate_max_bit_rate;
  up_transport_layer_information_t up_transport_layer_info;
  amf_pdn_type_value_t pdu_ip_type;
  qos_flow_request_list_t qos_flow_setup_request_list;
} pdu_session_resource_setup_request_transfer_t;

/***********************pdu_res_set_change ends*************************/

typedef enum SMSessionFSMState_response_s {
  CREATING_0,
  CREATE_1,
  ACTIVE_2,
  INACTIVE_3,
  RELEASED_4
} SMSessionFSMState_response;

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

typedef struct QosRules_response_s {
  uint32_t qos_rule_identifier;
  bool dqr;
  uint32_t number_of_packet_filters;
  uint32_t packet_filter_identifier[16];
  uint32_t qos_rule_precedence;
  bool segregation;
  uint32_t qos_flow_identifier;
} QosRules_response;

typedef struct AggregatedMaximumBitrate_respose_t {
  uint32_t max_bandwidth_ul;
  uint32_t max_bandwidth_dl;
} AggregatedMaximumBitrate_response;
typedef enum AmbrUnit_response_e {
  Kbps_0  = 0,
  Kbps_1  = 1,
  Kbps_4  = 2,
  Kbps_16 = 3,
  Kbps_64 = 4
} AmbrUnit_response;

typedef struct SessionAmbr_reponse_s {
  AmbrUnit_response downlink_unit_type;
  uint32_t downlink_units;  // Only to use lower 2 bytes (16 bit values)
  AmbrUnit_response uplink_unit_type;
  uint32_t uplink_units;  // Only to use lower 2 bytes (16 bit values)
} SessionAmbr_response;

typedef struct TeidSet_response_s {
  uint8_t teid[4];
  uint8_t end_ipv4_addr[16];
} TeidSet_response;

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
  QosRules_response
      authorized_qos_rules[4];  // TODO 32 in NAS5g,3 in pcap, revisit later
  SessionAmbr_response session_ambr;
  bool always_on_pdu_session_indication;
  ssc_mode_t allowed_ssc_mode;
  bool m5gsm_congetion_re_attempt_indicator;
  redirect_server_t pdu_address;

  qos_flow_request_list_t qos_list;
  TeidSet_response upf_endpoint;
  uint8_t procedure_trans_identity[2];
} itti_n11_create_pdu_session_response_t;

#define N11_CREATE_PDU_SESSION_RESPONSE(mSGpTR)                                \
  (mSGpTR)->ittiMsg.n11_create_pdu_session_response

#define N11_NOTIFICATION_RECEIVED(mSGpTR)                                      \
  (mSGpTR)->ittiMsg.n11_notification_received

// Resource relese request
// typedef enum radio_network_layer_cause_e
typedef enum {
  UNSPECIFIED,
  TXNRELOCOVERALL_EXPIRY,
  SUCCESSFUL_HANDOVER,
  RELEASE_DUE_TO_NG_RAN_GENERATED_REASON,
  RELEASE_DUE_TO_5GC_GENERATED,
  REASON,
  HANDOVER_CANCELLED,
  PARTIAL_HANDOVER,
  HANDOVER_FAILURE_IN_TARGET_5GC_NGRAN_NODE_OR_TARGET_SYSTEM,
  HANDOVER_TARGET_NOT_ALLOWED,
  TNGRELOCOVERALL_EXPIRY,
  TNGRELOCPREP_EXPIRY,
  CELL_NOT_AVAILABLE,
  UNKNOWN_TARGET_ID,
  NO_RADIO_RESOURCES_AVAILABLE_IN_TARGET_CELL,
  UNKNOWN_LOCAL_UE_NGAP_ID,
  INCONSISTENT_REMOTE_UE_NGAP_ID,
  HANDOVER_DESIRABLE_FOR_RADIO_REASONS,
  TIME_CRITICAL_HANDOVER,
  RESOURCE_OPTIMISATION_HANDOVER,
  REDUCE_LOAD_IN_SERVING_CELL,
  USER_INACTIVITY,
  RADIO_CONNECTION_WITH_UE_LOST,
  RADIO_RESOURCES_NOT_AVAILABLE,
  INVALID_QOS_COMBINATION,
  FAILURE_IN_THE_RADIO_INTERFACE_PROCEDURE,
  INTERACTION_WITH_OTHER_PROCEDURE,
  UNKNOWN_PDU_SESSION_ID,
  UNKNOWN_QOS_FLOW_ID,
  MULTIPLE_PDU_SESSION_ID_INSTANCES,
  MULTIPLE_QOS_FLOW_ID_INSTANCES,
  ENCRYPTION_AND_OR_INTEGRITY_PROTECTION_ALGORITHMS_NOT_SUPPORTED,
  NG_INTRA_SYSTEM_HANDOVER_TRIGGERED,
  XN_HANDOVER_TRIGGERED,
  NOT_SUPPORTED_5QI_VALUE,
  UE_CONTEXT_TRANSFER,
  IMS_VOICE_EPS_FALLBACK_OR_RAT_FALLBACK_TRIGGERED,
  UP_INTEGRITY_PROTECTION_NOT_POSSIBLE,
  UP_CONFIDENTIALITY_PROTECTION_NOT_POSSIBLE,
  SLICE_NOT_SUPPORTED,
  UE_IN_RRC_INACTIVE_STATE_NOT_REACHABLE,
  REDIRECTION,
  RESOURCES_NOT_AVAILABLE_FOR_THE_SLICE,
  UE_MAXIMUM_INTEGRITY_PROTECTED_DATA_RATE_REASON,
  RELEASE_DUE_TO_CN_DETECTED_MOBILITY,
  N26_INTERFACE_NOT_AVAILABLE,
  RELEASE_DUE_TO_PRE_EMPTION,
} radio_network_layer_cause;

typedef struct radio_network_layer_s {
  radio_network_layer_cause nw_layer_cause;
} radio_network_layer;

typedef enum {
  TRANSPORT_RESOURCE_UNAVAILABLE,
  UNSPECIFIED_TL,
} transport_layer_cause;

typedef struct transport_layer_s {
  transport_layer_cause cause;
} transport_layer_t;

typedef enum {
  NORMAL_RELEASE,
  AUTHENTICATION_FAILURE_NAS,  //#defined on AUTHENTICATION_FAILURE
  DEREGISTER,
  UNSPECIFIED_NAS_CAUSE,
} NAS_cause;

typedef struct NAS_s {
  NAS_cause cause;
} NAS_t;

typedef enum {
  TRANSFER_SYNTAX_ERROR,
  ABSTRACT_SYNTAX_ERROR_REJECT,
  ABSTRACT_SYNTAX_ERROR_IGNORE_AND_NOTIFY,
  MESSAGE_NOT_COMPATIBLE_WITH_RECEIVER_STATE,
  SEMANTIC_ERROR,
  ABSTRACT_SYNTAX_ERROR_FALSELY_CONSTRUCTED_MESSAGE,
  UNSPECIFIED_PROTOCOL,
} protocol_cause;

typedef struct Protocol_s {
  protocol_cause cause;
} protocol_t;

typedef enum {
  CONTROL_PROCESSING_OVERLOAD,
  NOT_ENOUGH_USER_PLANE_PROCESSING_RESOURCES,
  HARDWARE_FAILURE,
  O_AND_M_INTERVENTION,
  UNKNOWN_PLMN,
  UNSPECIFIED_MISC,
} miscellaneous_cause;

typedef struct miscellaneous_s {
  miscellaneous_cause cause;
} miscellaneous_t;
typedef enum {
  RADIO_NETWORK_LAYER_GROUP = 1,
  TRANSPORT_LAYER_GROUP,
  NAS_GROUP,
  PROTOCOL_GROUP,
  MISCELLANEOUS_GROUP,
} cause_group_e;

typedef struct cause_group_s {
  cause_group_e cause_group_type;
  union {
    radio_network_layer network_layer;
    transport_layer_t trasport_layer;
    NAS_t nas;
    protocol_t protocal;
    miscellaneous_t miscellaneous;
  } u_group;

} cause_group_t;

typedef struct cause_s {
  cause_group_t cause_group;
} cause_t;

typedef struct pdu_session_resource_release_command_transfer_s {
  cause_t cause;
} pdu_session_resource_release_command_transfer;

#define N11_NOTIFICATION_RECEIVED(mSGpTR)                                      \
  (mSGpTR)->ittiMsg.n11_notification_received

// RequestType
typedef enum RequestType_received_s {
  INITIAL_REQUEST                = 0,
  EXISTING_PDU_SESSION           = 1,
  INITIAL_EMERGENCY_REQUEST      = 2,
  EXISTING_EMERGENCY_PDU_SESSION = 3,
  MODIFICATION_REQUEST           = 4,
} RequestType_received;

// M5GSMCapability
typedef struct M5GSMCapability_received_s {
  bool reflective_qos;
  bool multi_homed_ipv6_pdu_session;
} M5GSMCapability_received;

#define N11_CREATE_PDU_SESSION_RESPONSE(mSGpTR)                                \
  (mSGpTR)->ittiMsg.n11_create_pdu_session_response

typedef enum {
  PDU_SESSION_INACTIVE_NOTIFY,         // AMF <=> SMF
  UE_IDLE_MODE_NOTIFY,                 // AMF  => SMF
  UE_PAGING_NOTIFY,                    // SMF  => AMF
  UE_PERIODIC_REG_ACTIVE_MODE_NOTIFY,  // AMF  => SMF
  PDU_SESSION_STATE_NOTIFY,            // SMF <=> AMF
  UE_SERVICE_REQUEST_ON_PAGING,        // AMF <=> SMF
} notify_ue_event;

typedef struct itti_n11_received_notification_s {
  // common context
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  SMSessionFSMState_response sm_session_fsm_state;
  uint32_t sm_session_version;
  // rat specific
  uint32_t pdu_session_id;
  RequestType_received request_type;
  // PduSessionType_response pdu_session_type;
  M5GSMCapability_received m5g_sm_capability;
  m5g_sm_cause_t m5gsm_cause;
  pdu_session_type_t pdu_session_type;
  // Idle/paging/periodic_reg events and UE state notification
  notify_ue_event notify_ue_evnt;
} itti_n11_received_notification_t;
