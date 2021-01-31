/**
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
/*****************************************************************************

  Source      amf_app_ue_context.h

  Version     0.1

  Date        2020/07/28

  Product     NAS stack

  Subsystem   Access and Mobility Management Function

  Author      Sandeep Kumar Mall

  Description Defines Access and Mobility Management Messages

*****************************************************************************/
#ifndef AMF_APP_UE_CONTEXT_AND_PROC_SEEN
#define AMF_APP_UE_CONTEXT_AND_PROC_SEEN

#include <sstream>
#include <thread>
#ifdef __cplusplus
extern "C" {
#endif
#include "bstrlib.h"
#include "3gpp_23.003.h"
#include "3gpp_24.301.h"
#include "3gpp_24.008.h"
#include "3gpp_38.331.h"
#include "3gpp_38.413.h"
#include "3gpp_24.501.h"
#include "TrackingAreaIdentity.h"
#include "hashtable.h"
#include "ngap_cause.h"
#include "obj_hashtable.h"
#include "queue.h"
#ifdef __cplusplus
};
#endif
#include "amf_fsm.h"
#include "amf_data.h"
#include "amf_common_defs.h"
#include "AmfMessage.h"
#include "M5GRegistrationRequest.h"
#include "ngap_messages_types.h"
//#include "M5GAuthenticationResponse.h"
//#include "amf_securityDef.h"
//=== amf_message.h related merged ==============

//#include "amf_app_ue_context_and_proc.h"
#include "M5GDLNASTransport.h"
#include "M5GRegistrationRequest.h"
#include "M5GRegistrationAccept.h"
#include "M5GIdentityRequest.h"
#include "M5GIdentityResponse.h"
#include "M5GAuthenticationRequest.h"
#include "M5GAuthenticationResponse.h"
#include "M5GSecurityModeCommand.h"
#include "M5GSecurityModeComplete.h"
#include "M5GSecurityModeReject.h"
#include "M5GDeRegistrationRequestUEInit.h"
#include "M5GDeRegistrationAcceptUEInit.h"
#include "M5GDeRegistrationRequestUEInit.h"
#include "M5GDeRegistrationAcceptUEInit.h"
#include "M5GULNASTransport.h"
#include "M5GDLNASTransport.h"

using namespace std;

namespace magma5g {
#define NAS5G_TIMER_INACTIVE_ID (-1)
class amf_procedures_t;
/*
 * Timer identifier returned when in inactive state (timer is stopped or has
 * failed to be started)
 */
#define AMF_APP_TIMER_INACTIVE_ID (-1)

#define AMF_APP_DELTA_T3512_REACHABILITY_TIMER 4            // in minutes
#define AMF_APP_DELTA_REACHABILITY_IMPLICIT_DETACH_TIMER 0  // in minutes

#define AMF_APP_INITIAL_CONTEXT_SETUP_RSP_TIMER_VALUE 2  // In seconds
#define AMF_APP_UE_CONTEXT_MODIFICATION_TIMER_VALUE 2    // In seconds
#define AMF_APP_PAGING_RESPONSE_TIMER_VALUE 4            // In seconds
#define AMF_APP_ULR_RESPONSE_TIMER_VALUE 3               // In seconds

#define NAS_SECURITY_ALGORITHMS_MINIMUM_LENGTH 1
#define NAS5G_SECURITY_ALGORITHMS_MAXIMUM_LENGTH 2

#define NAS5G_MESSAGE_CONTAINER_MAXIMUM_LENGTH 253

/* Timer structure */
typedef struct amf_app_timer_s {
  long id;  /* The timer identifier                 */
  long sec; /* The timer interval value in seconds  */
} amf_app_timer_t;

//=============================
//#if 1

typedef struct ie_header_t {
  uint16_t protocol_id;
  uint8_t criticality;
  uint8_t pad;
} ie_header;

#if 0
typedef struct gtp_tunnel_s {
  bstring endpoint_ip_address;  // Transport_Layer_Information
  uint8_t gtp_tied[4];
} gtp_tunnel;

typedef struct up_transport_layer_information_s {
  ie_header hdr; 
  gtp_tunnel gtp_tnl;
} up_transport_layer_information_t;

// typedef enum pre_emption_capability_e {
typedef enum {
  SHALL_NOT_TRIGGER_PRE_EMPTION,
  MAY_TRIGGER_PRE_EMPTION,
} pre_emption_capability;

// typedef enum pre_emption_vulnerability_e {
typedef enum {
  NOT_PREEMPTABLE,
  PRE_EMPTABLE,
} pre_emption_vulnerability;

typedef struct allocation_and_retention_priority_s {
  int priority_level;
  pre_emption_capability pre_emption_cap;
  pre_emption_vulnerability pre_emption_vul;
  // pre_emption_capability_t pre_emption_cap;
  // pre_emption_vulnerability_t pre_emption_vul;
} allocation_and_retention_priority;

typedef struct non_dynamic_5QI_descriptor_s {
  int fiveQI;
} non_dynamic_5QI_descriptor;
// Dynamic_5QI not cosidered

typedef struct qos_characteristics_s {
  non_dynamic_5QI_descriptor non_dynamic_5QI_desc;
} qos_characteristics;

typedef struct qos_flow_level_qos_parameters_s {
  qos_characteristics qos_characteristic;
  allocation_and_retention_priority alloc_reten_priority;

} qos_flow_level_qos_parameters;

typedef struct qos_flow_setup_request_item_s {
  uint32_t qos_flow_identifier;
  qos_flow_level_qos_parameters qos_flow_level_qos_param;
  // E-RAB ID is optional spec-38413 - 9.3.4.1
} qos_flow_setup_request_item;

typedef struct qos_flow_request_list_s {
  ie_header hdr; 
  qos_flow_setup_request_item qos_flow_req_item;

} qos_flow_request_list;
#endif  // pdu_res_set_change
// Resource relese request and response
// typedef enum radio_network_layer_cause_e
#if 0
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

typedef struct cause_group_s {
  radio_network_layer network_layer;
  transport_layer_t trasport_layer;
  NAS_t nas;
  protocol_t protocal;
  miscellaneous_t miscellaneous;
} cause_group_t;

typedef struct cause_s {
  cause_group_t cause_group;
} cause_t;

typedef struct pdu_session_resource_release_command_transfer_s {
  cause_t cause;
} pdu_session_resource_release_command_transfer;
#endif

typedef struct pdu_session_resource_to_release_item_s {
  PDUSessionIdentityMsg pdu_session_id;
  pdu_session_resource_release_command_transfer release_command_transfer;
} pdu_session_resource_to_release_item;

typedef struct pdu_session_resource_to_release_list_s {
  pdu_session_resource_to_release_item release_item;
} pdu_session_resource_to_release_list;

// Response failure message
typedef struct pdu_session_resource_setup_unsuccessful_transfer_s {
  cause_t cause;
} pdu_session_resource_setup_unsuccessful_transfer;

typedef struct pdu_session_resource_failed_to_setup_item_s {
  PDUSessionIdentityMsg pdu_session_id;
  pdu_session_resource_setup_unsuccessful_transfer unsuccessful_transfer;
} pdu_session_resource_failed_to_setup_item;

typedef struct amf_pdu_session_resource_setup_res_fail_list_s {
  pdu_session_resource_failed_to_setup_item setup_item;
} amf_pdu_session_resource_setup_res_fail_list;

// Response success message
typedef enum qos_flow_mapping_indication_e {
  UL,
  DL,
} qos_flow_mapping_indication;

typedef struct associated_qos_flow_item_s {
  uint32_t qos_flow_identifier;
  qos_flow_mapping_indication mapping_indication;
} associated_qos_flow_item;

typedef struct dl_qos_flow_per_tnl_info_s {
  up_transport_layer_information_t
      up_transport_layer_info;  // TODO uncomment after pdu_res_set_change
  associated_qos_flow_item flow_item;
} dl_qos_flow_per_tnl_info_t;

typedef struct pdu_session_resource_setup_response_transfer {
  dl_qos_flow_per_tnl_info_t dl_qos_flow_per_tnl_info;
} pdu_session_resource_setup_response_transfer;

typedef struct pdu_session_setup_response_success_item_s {
  PDUSessionIdentityMsg pdu_session_id;
  pdu_session_resource_setup_response_transfer response_transfer;
} pdu_session_setup_response_success_item;

typedef struct amf_pdu_session_resource_setup_res_success_list_s {
  uint16_t no_of_items;
  pdu_session_setup_response_success_item item_rsp_success;
} amf_pdu_session_resource_setup_res_success_list;

typedef struct pdu_session_resource_setup_rsp_s {
  amf_pdu_session_resource_setup_res_success_list
      pdu_ses_resource_rsp_success_list;
  amf_pdu_session_resource_setup_res_fail_list pdu_ses_resource_rsp_fail_list;

} pdu_session_resource_setup_rsp_t;

//=============================
// Structure to handle Resource Release Response from gNB
typedef enum {
  NR,
  E_UTRA,
} rat_type_e;

typedef struct volume_timed_report_item_s {
  // octetString start_timestamp;
  // octetString end_timestamp;
  uint32_t usage_count_ul;
  uint32_t usage_count_dl;
} volume_timed_report_item;

typedef struct pdu_session_usage_report_s {
  rat_type_e rat_type;
  volume_timed_report_item pdu_session_timed_report_list;
} pdu_session_usage_report;

typedef struct qos_flow_usage_report_item_s {
  uint32_t qos_flow_indicator;
  rat_type_e rat_type;
  volume_timed_report_item qos_flows_timed_report_list;
} qos_flow_usage_report_item;

typedef struct secondary_rat_usage_information_s {
  pdu_session_usage_report usage_report;
  qos_flow_usage_report_item report_item;
} secondary_rat_usage_information;

typedef struct pdu_session_resource_release_response_transfer_s {
  secondary_rat_usage_information rat_usage_Information;
} pdu_session_resource_release_response_transfer;

typedef struct pdu_session_resource_released_item_s {
  PDUSessionIdentityMsg pdu_session_id;
  pdu_session_resource_release_response_transfer release_response_transfer;
} pdu_session_resource_released_item;

typedef struct pdu_session_resource_released_list_s {
  pdu_session_resource_released_item released_item;
} pdu_session_resource_released_list;

#if 0
typedef struct amf_ue_aggregate_maximum_bit_rate_s{
   ie_header hdr; 
   uint64_t dl;
   uint64_t ul;
}amf_ue_aggregate_maximum_bit_rate_t;
typedef struct amf_pdn_type_value_s {
   ie_header hdr;
   pdn_type_value_t pdn_type;
}amf_pdn_type_value_t;
//=============================
typedef struct pdu_session_resource_setup_request_transfer_s {
  //ngap_ue_aggregate_maximum_bit_rate_t pdu_aggregate_max_bit_rate;
  amf_ue_aggregate_maximum_bit_rate_t pdu_aggregate_max_bit_rate;
  up_transport_layer_information_t up_transport_layer_info;
  amf_pdn_type_value_t pdu_ip_type;
  qos_flow_request_list qos_flow_setup_request_list;
} pdu_session_resource_setup_request_transfer_t;
#endif  // pdu_res_set_change
//=============================

// Reusing structures
typedef struct pdu_session_resource_setup_req_s {
  /*
   * values: 1-1Kbps, 2- 4Kbps, 3- 16Kbps, 4- 64Kbps 24.501 9.11.4.14
   */
  uint8_t units_for_session;
  ngap_ue_aggregate_maximum_bit_rate_t pdu_aggregate_maximum_bit_rate;
  Ngap_PDUSessionID_t Pdu_Session_ID;  // from NGAP
  Ngap_SNSSAI_t Ngap_s_nssai;          // S-NSSAI as defined in TS 23.003 [23]
  // heavily nested structure
  pdu_session_resource_setup_request_transfer_t
      pdu_session_resource_setup_request_transfer;  // pdu_res_set_change
} pdu_session_resource_setup_req_t;

/*GTP tunnel id for UPF and gNB exchange infomration*/
typedef struct teid_upf_gnb_s {
  uint8_t upf_gtp_teid_ip_addr[16];
  // bstring upf_gtp_teid_ip_addr;
  uint8_t upf_gtp_teid[4];
  // bstring gnb_gtp_teid_ip_addr;
  uint8_t gnb_gtp_teid_ip_addr[16];
  uint8_t gnb_gtp_teid[4];
} teid_upf_gnb_t;

typedef struct smf_proc_data_s {
  PDUSessionIdentityMsg pdu_session_identity;
  PTIMsg pti;
  MessageTypeMsg message_type;
  IntegrityProtMaxDataRateMsg integrity_prot_max_data_rate;
  PDUSessionTypeMsg pdu_session_type;
  SSCModeMsg ssc_mode;
  bstring apn;  // apn is dnn
  // smf_proc_pdn_type_t pdn_type;
  bstring pdn_addr;
} smf_proc_data_t;

typedef struct smf_context_s {
  uint32_t n_active_pdus;
  bool is_emergency;
  uint8_t dl_ambr_unit;
  uint16_t dl_session_ambr;  // Session amber downlink
  uint8_t ul_ambr_unit;
  uint16_t ul_session_ambr;  // Session amber uplink
  QOSRule qos_rules[32];     // QOS rules
  teid_upf_gnb_t gtp_tunnel_id;
  // qos_rule_t qos_rules[255];  // QOS rules
  smf_proc_data_t smf_proc_data;
  // Request to gnb on PDU establisment request
  pdu_session_resource_setup_req_t pdu_resource_setup_req;
  pdu_session_resource_setup_rsp_t pdu_resource_setup_rsp;
  pdu_session_resource_to_release_list pdu_resource_release_req;
  pdu_session_resource_to_release_list pdu_resource_release_rsp;
  // struct nas_timer_s T3489; //TODO implement timer
} smf_context_t;

/*
 * Structure of the AMF context established by the network for a particular UE
 * ---------------------------------------------------------------------------
 */
class amf_context_t {
 public:
  amf_context_t() {}
  ~amf_context_t() {}
  bool is_dynamic;    /* Dynamically allocated context indicator         */
  bool is_registered; /* Registration indicator                            */
  // bool is_emergency; /* Emergency bearer services indicator             */
  bool is_initial_identity_imsi;  // If the IMSI was used for identification in
                                  // the initial NAS message
  bool is_guti_based_registered;
  /*
   * registration_type has type amf_proc_registration_type_t.
   *
   * Here, it is un-typedef'ed as uint8_t to avoid circular dependency issues.
   */
  uint32_t member_present_mask; /* bitmask, see significance of bits below */
  uint32_t member_valid_mask;   /* bitmask, see significance of bits below */

  uint8_t m5gsregistrationtype;
  amf_procedures_t* amf_procedures;

  uint num_registration_request; /* Num registration request received */

  // imsi present mean we know it but was not checked with identity proc, or was
  // not provided in initial message
  imsi_t _imsi;     /* The IMSI provided by the UE or the AMF, set valid when
                       identification returns IMSI */
  imsi64_t _imsi64; /* The IMSI provided by the UE or the AMF, set valid when
                       identification returns IMSI */
  imsi64_t saved_imsi64; /* Useful for 5.4.2.7.c */
  imei_t _imei;          /* The IMEI provided by the UE                     */
  imeisv_t _imeisv;      /* The IMEISV provided by the UE                   */
  // bool                   _guti_is_new; /* The GUTI assigned to the UE is new
  // */
  guti_m5_t _m5_guti;    /* The GUTI assigned to the UE                     */
  guti_m5_t m5_old_guti; /* The GUTI assigned to the UE                     */
  // tai_list_t _tai_list;   /* TACs the the UE is registered to */ tai_t
  // _lvr_tai; tai_t originating_tai;

  ksi_t ksi; /*key set identifier  */
  // ue_network_capability_t _ue_network_capability; // will be use in perodic
  // registration ms_network_capability_t _ms_network_capability;
  drx_parameter_t _drx_parameter;

  int remaining_vectors;  // remaining unused vectors
  m5g_auth_vector_t
      _vector[MAX_EPS_AUTH_VECTORS]; /* 5GMM authentication vector */
  amf_security_context_t
      _security; /* Current 5GMM security context: The security context which
                    has been activated most recently. Note that a current 5GMM
                                       //security context originating from
                    either a mapped or native 5GMM security context may exist
                    simultaneously with a native
                                      // non-current 5GMM security context.*/

  // Requirement MME24.301R10_4.4.2.1_2
  // amf_security_context_t  _non_current_security; /* Non-current 5GMM security
  // context: A native 5GMM security context that is not the current one. A
  // non-current 5GMM
  /* security context may be stored along with a current 5GMM security context
   in the UE and the MME. A non-current 5GMM security context does not contain
   an 5GMM AS security context. A non-current 5GMM security context is either of
   type 'full native' or of type 'partial native'.     */

  int amf_cause; /* EMM failure cause code                          */

  amf_fsm_state_t amf_fsm_state;

  // nas_timer_t T3422; /* Deregister timer         */
  void* t3422_arg;

  smf_context_t smf_context;  // smf contents

  drx_parameter_t _current_drx_parameter; /* stored TAU Request IE Requirement
                                             AMF24.501R15_5.5.3.2.4_4*/

  // TODO: DO BETTER  WITH BELOW
  std::string smf_msg; /* SMF message contained within the initial request*/
  bool is_imsi_only_detach;
};

class amf_ue_context_t {
 public:
  amf_ue_context_t() {}
  ~amf_ue_context_t() {}

  /* hash_table_uint64_ts_t is defined in lib/hastable*/
  hash_table_uint64_ts_t* imsi_amf_ue_id_htbl;    // data is amf_ue_ngap_id_t
  hash_table_uint64_ts_t* tun11_ue_context_htbl;  // data is amf_ue_ngap_id_t
  hash_table_uint64_ts_t*
      gnb_ue_ngap_id_ue_context_htbl;             // data is amf_ue_ngap_id_t
  obj_hash_table_uint64_t* guti_ue_context_htbl;  // data is amf_ue_ngap_id_t
};

enum m5gmm_state_t {
  UE_UNREGISTERED = 0,
  UE_REGISTERED,
};
enum m5gcm_state_t {
  M5GCM_IDLE = 0,
  M5GCM_CONNECTED,
};
class ue_m5gmm_context_s  //:public amf_context_t
{
 public:
  ue_m5gmm_context_s() {}
  ~ue_m5gmm_context_s() {}
  /* msisdn: The basic MSISDN of the UE. The presence is dictated by its storage
   *         in the HSS, set by S6A UPDATE LOCATION ANSWER
   */
  bstring msisdn;

  ngap_Cause_t
      ue_context_rel_cause;  // define require for Ngcause in NGAP module
  m5gmm_state_t mm_state;
  m5gcm_state_t ecm_state;

  /* Last known 5G cell, set by nas_registration_req_t */
  ecgi_t e_utran_cgi;  // from 3gpp 23.003

  /* cell_age: Time elapsed since the last 5G Cell Global Identity was
   *           acquired. Set by nas_auth_param_req_t
   */
  time_t cell_age;
  /* TODO: add csg_id */
  /* TODO: add csg_membership */
  /* TODO Access mode: Access mode of last known ECGI when the UE was active */

  /* apn_config_profile: set by UPDATE LOCATION ANSWER */
  apn_config_profile_t apn_config_profile;

  amf_context_t amf_context;

  /* access_restriction_data: The access restriction subscription information.
   *           set by UPDATE LOCATION ANSWER
   */
  ard_t access_restriction_data;

  bstring apn_oi_replacement;
  teid_t amf_teid_n11;
  /* SCTP assoc id */
  sctp_assoc_id_t sctp_assoc_id_key;

  /* gNB UE NGAP ID,  Unique identity the UE within gNodeB */
  gnb_ue_ngap_id_t gnb_ue_ngap_id;
  // gnb_ue_ngap_id_t gnb_ue_ngap_id : 24;

  gnb_ngap_id_key_t gnb_ngap_id_key;

  /* AMF UE NGAP ID, Unique identity the UE within AMF */
  amf_ue_ngap_id_t amf_ue_ngap_id;

  /* Subscribed UE-AMBR: The Maximum Aggregated uplink and downlink MBR values
   *           to be shared across all Non-GBR bearers according to the
   *           subscription of the user. Set by SMF UPDATE LOCATION ANSWER
   */
  ambr_t subscribed_ue_ambr;
  /* used_ue_ambr: The currently used Maximum Aggregated uplink and downlink
   *           MBR values to be shared across all Non-GBR bearers.
   *           Set by S6A UPDATE LOCATION ANSWER
   */
  ambr_t used_ue_ambr;
  /* rau_tau_timer: Indicates a subscribed Periodic RAU/TAU Timer value
   *           Set by S6A UPDATE LOCATION ANSWER
   */
  rau_tau_timer_t rau_tau_timer;

  int nb_active_pdn_contexts;
  // pdn_context_t* pdn_contexts[MAX_APN_PER_UE];

  // TODO Required during dual connectivity communication
  // bearer_context_t* bearer_contexts[BEARERS_PER_UE];

  /* ue_radio_capability: Store the radio capabilities as received in
   *           S1AP UE capability indication message
   */
  bstring ue_radio_capability;

  /* mobile_reachability_timer: Start when UE moves to idle state.
   *             Stop when UE moves to connected state
   */
  amf_app_timer_t m5_mobile_reachability_timer;
  /* implicit_detach_timer: Start at the expiry of Mobile Reachability timer.
   * Stop when UE moves to connected state
   */
  amf_app_timer_t m5_implicit_detach_timer;
  /* Initial Context Setup Procedure Guard timer */
  amf_app_timer_t m5_initial_context_setup_rsp_timer;
  /* UE Context Modification Procedure Guard timer */
  amf_app_timer_t m5_ue_context_modification_timer;
  /* Timer for retrying paging messages */
  amf_app_timer_t m5_paging_response_timer;
  /* send_ue_purge_request: If true AMF shall send - Purge Req to
   * delete contexts at HSS
   */
  bool send_ue_purge_request;

  bool hss_initiated_detach;
  bool location_info_confirmed_in_hss;
  /* S6a- update location request guard timer */
  amf_app_timer_t m5_ulr_response_timer;

  uint8_t registration_type;
  lai_t lai;
  /* granted_service_t: informs the granted service to UE */
  // TODO required during dual connectivity (DC)
  // granted_service_t m5_granted_service;
  /*  paging_proceeding_flag (PPF) shall set to true, when UE moves to
   * connected. Indicates that paging procedure can be prooceeded, Is set to
   * false, due to "Inactivity of UE including lack of periodic TAU"
   */
  bool ppf;

#define SUBSCRIPTION_UNKNOWN false
#define SUBSCRIPTION_KNOWN true
  bool subscription_known;
  ambr_t used_ambr;
  subscriber_status_t subscriber_status;
  network_access_mode_t network_access_mode;

  bool path_switch_req;
  // LIST_HEAD(n11_procedures_s, amf_app_n11_proc_s) * n11_procedures;
};
/** @class ue_m5gmm_context_s
 *  @brief Useful parameters to know in AMF application layer. They are set
 * according to 3GPP TS.23.518 #6.1.6.2.25
 */
class ue_mm_context {
 public:
  ue_mm_context() {}
  ~ue_mm_context() {}

  /* msisdn: The basic MSISDN of the UE. The presence is dictated by its storage
   *         in the UDM, set by N8 UPDATE LOCATION ANSWER
   */
  std::string imsi;
  bool supi_UnauthInd;
  // std::string gpsiList[] array(Gpsi);
  std::string pei;
  uint64_t udmGroupId;   // NfGroupId
  uint64_t ausfGroupId;  // NfGroupId;
  std::string routingIndicator;
  // std::auto groupList[] array(GroupId);
  std::string drxParameter;
  std::string subRfsp;
  uint32_t usedRfsp;  // RfspIndex type;
  ambr_t subUeAmbr;
  bool smsSupport;
  std::string smsfId;  // NfInstanceId type
  std::string
      seafData;  // SeafData which will come while AUSF communication for AUTN.
  // M5GMM_Capability_msg m5gMmCapability; //5GMmCapability
  std::string pcfId;           // NfInstanceId
  std::string pcfAmPolicyUri;  // Uri
  // std::auto amPolicyReqTriggerList;//array(PolicyReqTrigger)
  std::string pcfUePolicyUri;  // Uri
  // std::auto uePolicyReqTriggerList; //array(PolicyReqTrigger)
  std::string hpcfId;                  // NfInstanceId
  std::string restrictedRatList;       // array(RatType)
  std::string forbiddenAreaList;       // array(Area)
  std::string serviceAreaRestriction;  // ServiceAreaRestriction
  std::string restrictedCnList;        // array(CoreNetworkType)
  std::string eventSubscriptionList;   // array(AmfEventSubscription)
  std::string mmContextList;           // array(MmContext)
  std::string sessionContextList;      // array(PduSessionContext)
  std::string traceData;               // TraceData
  /* SCTP assoc id */
  sctp_assoc_id_t sctp_assoc_id_key;

  /* gNB UE NGAP ID,  Unique identity the UE within gNodeB */
  gnb_ue_ngap_id_t gnb_ue_ngap_id : 24;

  gnb_ngap_id_key_t gnb_ngap_id_key;

  /* MME UE S1AP ID, Unique identity the UE within MME */
  amf_ue_ngap_id_t amf_ue_ngap_id;
};

class amf_app_ue_context : public amf_ue_context_t, public ue_m5gmm_context_s {
 public:
  amf_app_ue_context() {}
  ~amf_app_ue_context() {}
  // check & create state information.
  int amf_insert_ue_context(
      const amf_ue_context_t* amf_ue_context,
      const ue_m5gmm_context_s* ue_context_p);
  amf_ue_ngap_id_t amf_app_ctx_get_new_ue_id(
      amf_ue_ngap_id_t* amf_app_ue_ngap_id_generator_p);
  // Notify NGAP about the mapping between amf_ue_ngap_id and
  // sctp assoc id + gnb_ue_ngap_id
  void notify_ngap_new_ue_amf_ngap_id_association(
      const ue_m5gmm_context_s* ue_context_p);
  void amf_remove_ue_context(
      amf_ue_context_t* const amf_ue_context,
      ue_m5gmm_context_s* const ue_context_p);
};

ue_m5gmm_context_s* amf_create_new_ue_context(void);

/* Retrive required UE context from respective hash table*/
amf_context_t* amf_context_get(const amf_ue_ngap_id_t ue_id);
ue_m5gmm_context_s* amf_ue_context_exists_amf_ue_ngap_id(
    const amf_ue_ngap_id_t amf_ue_ngap_id);
int amf_context_upsert_imsi(amf_context_t* elm) __attribute__((nonnull));

void amf_ctx_set_valid_imsi(
    amf_context_t* ctxt, imsi_t* imsi, const imsi64_t imsi64)
    __attribute__((nonnull)) __attribute__((flatten));
void amf_ctx_set_attribute_valid(
    amf_context_t* ctxt, const uint32_t attribute_bit_pos)
    __attribute__((nonnull)) __attribute__((flatten));
void amf_ctx_set_attribute_present(
    amf_context_t* ctxt, const int attribute_bit_pos) __attribute__((nonnull))
__attribute__((flatten));

//========== merged from amf_message.h =======

class amf_msg_header;
// moved from amf_msfdefs.com
/* Header length boundaries of 5GS Mobility Management messages  */
#define AMF_HEADER_LENGTH sizeof(amf_msg_header)
#define AMF_HEADER_MINIMUM_LENGTH AMF_HEADER_LENGTH
#define AMF_HEADER_MAXIMUM_LENGTH AMF_HEADER_LENGTH
class amf_msg_header {
 public:
  uint8_t extended_protocol_discriminator;
  uint8_t security_header_type;
  uint8_t message_type;
  uint8_t sequence_number;
};

// moved from amf_app_msg.h
class amf_app_msg {
 public:
  void amf_app_ue_context_release(
      ue_m5gmm_context_s* ue_context_p, ngap_Cause_t cause);
};

// move from amf_nas_message.h
class amf_nas_message_decode_status_t {
 public:
  amf_nas_message_decode_status_t() {}
  ~amf_nas_message_decode_status_t() {}
  uint8_t integrity_protected_message : 1;
  uint8_t ciphered_message : 1;
  uint8_t mac_matched : 1;
  uint8_t security_context_available : 1;
  int amf_cause;
};

class AMFMsg {
 public:
  AMFMsg() {}

  ~AMFMsg() {}

  amf_msg_header header;

  RegistrationRequestMsg registrationrequestmsg;

  RegistrationAcceptMsg registrationacceptmsg;

  RegistrationCompleteMsg registrationcompletemsg;

  RegistrationRejectMsg registrationrejectmsg;

  IdentityRequestMsg identityrequestmsg;

  IdentityResponseMsg identityresponsemsg;

  AuthenticationRequestMsg authenticationrequestmsg;

  AuthenticationResponseMsg authenticationresponsemsg;

  AuthenticationRejectMsg authenticationrejectmsg;

  AuthenticationFailureMsg authenticationfailuremsg;

  SecurityModeCommandMsg securitymodecommandmsg;

  SecurityModeCompleteMsg securitymodecompletemsg;

  SecurityModeRejectMsg securitymodereject;

  DeRegistrationRequestUEInitMsg deregistrationequesmsg;

  DeRegistrationAcceptUEInitMsg deregistrationacceptmsg;

  ULNASTransportMsg uplinknas5gtransport;

  DLNASTransportMsg downlinknas5gtransport;

  // SERVICE REQUEST
  int amf_msg_decode_header(
      amf_msg_header* header, const uint8_t* buffer, uint32_t len);

  int amf_msg_encode_header(
      const amf_msg_header* header, uint8_t* buffer, uint32_t len);

  int amf_msg_decode(AMFMsg* msg, uint8_t* buffer, uint32_t len);

  int AmfMsgEncode(AMFMsg* msg, uint8_t* buffer, uint32_t len);
};

/* union of plain NAS message */
typedef struct nas_message_plain_s {
  AMFMsg amf; /* 5GMM Mobility Management messages */
  // SMFMsg smf; /*TODO 5GS Session Management messages  */
} nas_message_plain_t;

typedef struct nas_message_security_protected_s {
  amf_msg_header header;
  nas_message_plain_t plain;
} nas_message_security_protected_t;

typedef struct amf_nas_message_s {
  amf_msg_header header;
  nas_message_security_protected_t security_protected;
  nas_message_plain_t plain;
} amf_nas_message_t;

//=========== merged from amf_nas5g_proc.h======

typedef enum {
  CN5G_PROC_NONE = 0,
  CN5G_PROC_AUTH_INFO,
} cn5g_proc_type_t;

typedef enum amf_common_proc_type_s {
  AMF_COMM_PROC_NONE = 0,
  AMF_COMM_PROC_GUTI,
  AMF_COMM_PROC_AUTH,
  AMF_COMM_PROC_SMC,
  AMF_COMM_PROC_IDENT,
  AMF_COMM_PROC_INFO,
} amf_common_proc_type_t;

enum nas_base_proc_type_t {
  NAS_PROC_TYPE_NONE = 0,
  NAS_PROC_TYPE_AMF,
  NAS_PROC_TYPE_SMF,
  NAS_PROC_TYPE_CN,
};

class nas5g_base_proc_t;
class nas_amf_proc_t;
class nas_amf_registration_proc_t;

typedef int (*success_cb_t)(amf_context_t* amf_ctx);
typedef int (*failure_cb_t)(amf_context_t* amf_ctx);
typedef int (*proc_abort_t)(
    amf_context_t* amf_ctx, nas5g_base_proc_t* nas_proc);

typedef int (*pdu_in_resp_t)(
    amf_context_t* amf_ctx, void* arg);  // can be RESPONSE, COMPLETE, ACCEPT
typedef int (*pdu_in_rej_t)(amf_context_t* amf_ctx, void* arg);  // REJECT.
typedef int (*pdu_out_rej_t)(
    amf_context_t* amf_ctx, nas5g_base_proc_t* nas_proc);  // REJECT.
typedef void (*time_out_t)(void* arg);

typedef int (*sdu_out_delivered_t)(
    amf_context_t* amf_ctx, nas_amf_proc_t* nas_proc);
typedef int (*sdu_out_not_delivered_t)(
    amf_context_t* amf_ctx, nas_amf_proc_t* nas_proc);
typedef int (*sdu_out_not_delivered_ho_t)(
    amf_context_t* amf_ctx, nas_amf_proc_t* nas_proc);
/*Global variable and needed to increment based on nas procedures*/
static uint64_t nas_puid = 1;
#if 1  // reverted from amf_nas_common.h
class nas5g_base_proc_t {
 public:
  success_cb_t success_notif;
  failure_cb_t failure_notif;
  proc_abort_t abort;

  // PDU interface
  // pdu_in_resp_t           resp_in;
  pdu_in_rej_t fail_in;
  pdu_out_rej_t fail_out;
  time_out_t time_out;
  nas_base_proc_type_t type;  // AMF, SMF, CN
  uint64_t nas_puid;          // procedure unique identifier for internal use

  class nas5g_base_proc_t* parent;
  class nas5g_base_proc_t* child;
};

enum nas_amf_proc_type_t {
  NAS_AMF_PROC_TYPE_NONE = 0,
  NAS_AMF_PROC_TYPE_SPECIFIC,
  NAS_AMF_PROC_TYPE_COMMON,
  NAS_AMF_PROC_TYPE_CONN_MNGT,
};
// AMF Specific procedures
class nas_amf_proc_t {
 public:
  nas5g_base_proc_t base_proc;
  nas_amf_proc_type_t type;  // specific, common, connection management
  // SDU interface
  sdu_out_delivered_t delivered;
  sdu_out_not_delivered_t not_delivered;
  sdu_out_not_delivered_ho_t not_delivered_ho;

  amf_fsm_state_t previous_amf_fsm_state;
};

typedef struct nas_amf_common_proc_s {
  nas_amf_proc_t amf_proc;
  amf_common_proc_type_t type;
} nas_amf_common_proc_t;

enum amf_specific_proc_type_t {
  AMF_SPEC_PROC_TYPE_NONE = 0,
  AMF_SPEC_PROC_TYPE_REGISTRATION,
  AMF_SPEC_PROC_TYPE_DEREGISTRATION,
  AMF_SPEC_PROC_TYPE_TAU,
};

/*Deregistration specific elements*/

typedef enum deregistration_switchoff_e {
  /*In the UE to network direction, octate 1, 4th bit*/
  AMF_NORMAL_DEREGISTRATION = 0,
  AMF_SWITCHOFF_DEREGISTRATION,
} deregistration_switchoff;

typedef enum deregistration_access_type_e {
  AMF_3GPP_ACCESS = 1,
  AMF_NONE_3GPP_ACCESS,
  AMF_3GPP_ACCESS_AND_NONE_3GPP_ACCESS,
} deregistration_access_type;

typedef struct amf_deregistration_request_ies_s {
  deregistration_switchoff de_reg_type;
  deregistration_access_type de_reg_access_type;
  // bool switch_off;
  ksi_t ksi;
  TmsiM5GSMobileIdentity* tmsi;
  SuciM5GSMobileIdentity* suci;
  ImsiM5GSMobileIdentity* imsi;
  ImeiM5GSMobileIdentity* imei;
  GutiM5GSMobileIdentity* guti;
  amf_nas_message_decode_status_t decode_status;
} amf_deregistration_request_ies_t;

// AMF Specific procedures
class nas_amf_specific_proc_t {
 public:
  nas_amf_proc_t amf_proc;
  amf_specific_proc_type_t type;
};

#endif

class identification : public amf_context_t {
 public:
  char amf_identity_type_str[5][15] = {"NOT AVAILABLE", "IMSI", "IMEI",
                                       "IMEISV", "TMSI"};
  // static char *amf_identity_type_str[] = {"NOT AVAILABLE", "IMSI", "IMEI",
  // "IMEISV", "TMSI"}; static const char* amf_identity_type_str[] = {"NOT
  // AVAILABLE", "IMSI", "IMEI", "IMEISV", "TMSI"};
  int amf_proc_identification(
      amf_context_t* const amf_context, nas_amf_proc_t* const amf_proc,
      const identity_type2_t type, success_cb_t success, failure_cb_t failure);
  int amf_proc_identification_complete(
      const amf_ue_ngap_id_t ue_id, imsi_t* const imsi, imei_t* const imei,
      imeisv_t* const imeisv, uint32_t* const tmsi);
};

// typedef struct nas_amf_common_proc_s {
//  nas_amf_proc_t amf_proc;
//  amf_common_proc_type_t type;
//} nas_amf_common_proc_t;

class nas_amf_auth_proc_t {
 public:
  nas_amf_common_proc_t amf_com_proc;
  nas5g_timer_t T3560; /* Authentication timer         */
#define AUTHENTICATION_COUNTER_MAX 5
  unsigned int retransmission_count;
#define EMM_AUTHENTICATION_SYNC_FAILURE_MAX 2
  unsigned int
      sync_fail_count; /* counter of successive AUTHENTICATION FAILURE messages
                      from the UE with AMF cause #21 "synch failure" */
  unsigned int mac_fail_count;
  amf_ue_ngap_id_t ue_id;
  bool is_cause_is_registered;  //  could also be done by seeking parent
                                //  procedure
  ksi_t ksi;
  uint8_t rand[AUTH_RAND_SIZE]; /* Random challenge number  */
  uint8_t autn[AUTH_AUTN_SIZE]; /* Authentication token     */
  imsi_t* unchecked_imsi;
  int amf_cause;
};

typedef struct nas5g_cn_proc_s {
  nas5g_base_proc_t base_proc;
  cn5g_proc_type_t type;
} nas5g_cn_proc_t;

typedef struct nas5g_cn_procedure_s {
  nas5g_cn_proc_t* proc;
  LIST_ENTRY(nas5g_cn_procedure_s) entries;
} nas5g_cn_procedure_t;

class nas_5g_auth_info_proc_t {
 public:
  nas5g_cn_proc_t cn_proc;
  success_cb_t success_notif;
  failure_cb_t failure_notif;
  bool request_sent;
  uint8_t nb_vectors;
  // m5g_vector_t* vector[MAX_5G_AUTH_VECTORS];//TODO Check with Sandeep
  int nas_cause;
  amf_ue_ngap_id_t ue_id;
  bool resync;  // Indicates whether the authentication information is requested
                // due to sync failure
};

nas_5g_auth_info_proc_t* nas_new_5gcn_auth_info_procedure(
    amf_context_t* const amf_context);

class m5g_authentication : public amf_context_t {
 public:
  int amf_proc_authentication_ksi(
      amf_context_t* amf_context,
      nas_amf_specific_proc_t* const amf_specific_proc, ksi_t ksi,
      const uint8_t* const rand, const uint8_t* const autn,
      success_cb_t success, failure_cb_t failure);

  int amf_proc_authentication(
      amf_context_t* amf_context,
      nas_amf_specific_proc_t* const amf_specific_proc, success_cb_t success,
      failure_cb_t failure);

  int amf_proc_authentication_failure(
      amf_ue_ngap_id_t ue_id, int amf_cause, const_bstring auts);

  int amf_proc_authentication_complete(
      amf_ue_ngap_id_t ue_id, AuthenticationResponseMsg* msg, int amf_cause,
      const unsigned char* res);

  int amf_registration_security(amf_context_t* amf_context);

  void set_notif_callbacks_for_5g_auth_proc(nas_amf_auth_proc_t* auth_proc);
  void set_callbacks_for_5g_auth_proc(nas_amf_auth_proc_t* auth_proc);
  void set_callbacks_for_5g_auth_info_proc(
      nas_5g_auth_info_proc_t* auth_info_proc);
  int amf_send_authentication_request(
      amf_context_t* amf_context, nas_amf_auth_proc_t* auth_proc);
};

class nas_proc {
 public:
  int nas_proc_establish_ind(
      const amf_ue_ngap_id_t ue_id, const bool is_mm_ctx_new,
      const tai_t originating_tai, const ecgi_t ecgi,
      const m5g_rrc_establishment_cause_t as_cause, const s_tmsi_m5_t s_tmsi,
      bstring msg);
  nas_amf_registration_proc_t* get_nas_specific_procedure_registration(
      const amf_context_t* ctxt);
  bool is_nas_specific_procedure_registration_running(
      const amf_context_t* ctxt);
  int nas5g_message_decode(
      unsigned char* buffer, amf_nas_message_t* nas_msg, int length,
      amf_security_context_t* amf_security_context,
      amf_nas_message_decode_status_t* decode_status);
  nas_amf_common_proc_t* get_nas5g_common_procedure(
      const amf_context_t* const ctxt, amf_common_proc_type_t proc_type);
};
// 5G CN Specific procedures
typedef struct nas_amf_common_procedure_s {
  nas_amf_common_proc_t* proc;
  LIST_ENTRY(nas_amf_common_procedure_s) entries;
} nas_amf_common_procedure_t;

typedef struct nas_amf_ident_proc_s {
  nas_amf_common_proc_t amf_com_proc;
  nas5g_timer_t T3570; /* Identification timer         */
#define IDENTIFICATION_COUNTER_MAX 5
  unsigned int retransmission_count;
  amf_ue_ngap_id_t ue_id;
  bool is_cause_is_registered;  //  could also be done by seeking parent
                                //  procedure
  identity_type2_t identity_type;
} nas_amf_ident_proc_t;

//======== merged from amf_proc.h================

enum amf_proc_registration_type_t {
  AMF_REGISTRATION_TYPE_INITIAL = 1,
  AMF_REGISTRATION_TYPE_MOBILITY_UPDATING,
  AMF_REGISTRATION_TYPE_PERODIC_UPDATING,
  AMF_REGISTRATION_TYPE_EMERGENCY,
  AMF_REGISTRATION_TYPE_RESERVED = 7,
};
class amf_registration_request_ies_t : public RegistrationRequestMsg {
 public:
  amf_registration_request_ies_t() {}
  ~amf_registration_request_ies_t() {}
  // need to put registration ies.
  bool is_initial;
  amf_proc_registration_type_t
      m5gsregistrationtype;  // m5gsregistrationtype=AMF_REGISTRATION_TYPE_RESERVED;
  // additional_update_type_t additional_update_type;
  bool is_native_sc;
  // ksi_t ksi; ngKSI_t
  bool is_native_guti;
  guti_m5_t* guti;
  imsi_t* imsi;
  imei_t* imei;
  tai_t* last_visited_registered_tai;  // Last visited registered TAI
  tai_t* originating_tai;
  ecgi_t* originating_ecgi;                        //
  ue_network_capability_t ue_network_capability;   // UE security capability
  ms_network_capability_t* ms_network_capability;  // 5GMM capability
  drx_parameter_t* drx_parameter;  // Requested DRX parameters during paging
  bstring smf_msg;
  amf_nas_message_decode_status_t decode_status;
};

typedef struct nas5g_proc_mess_sign_s {
  uint64_t puid;
#define NAS5G_MSG_DIGEST_SIZE 16
  uint8_t digest[NAS5G_MSG_DIGEST_SIZE];
  size_t digest_length;
  size_t nas_msg_length;
} nas5g_proc_mess_sign_t;

class amf_procedures_t {
 public:
  nas_amf_specific_proc_t* amf_specific_proc;
  LIST_HEAD(nas_amf_common_procedures_head_s, nas_amf_common_procedure_s)
  amf_common_procs;  // TODO -  NEED-RECHECK
  LIST_HEAD(nas5g_cn_procedures_head_s, nas5g_cn_procedure_s)
  cn_procs;  // triggered by AMF
             // nas_amf_con_mngt_proc_t* amf_con_mngt_proc;

  int nas_proc_mess_sign_next_location;  // next index in array
#define MAX_NAS_PROC_MESS_SIGN 3
  nas5g_proc_mess_sign_t nas_proc_mess_sign[MAX_NAS_PROC_MESS_SIGN];
};

/*
0 0 1 initial registration
0 1 0 mobility registration updating
0 1 1 periodic registration updating
1 0 0 emergency registration
*/

//============= merged from amf_nas_common_defs.h=========

class nas_amf_registration_proc_t  //: public amf_registration_request_ies_t
{
 public:
  nas_amf_registration_proc_t() {}
  ~nas_amf_registration_proc_t() {}
  nas_amf_specific_proc_t amf_spec_proc;
  // struct nas_timer_s T3450; // AMF message retransmission timer
  //#define REGISTRATION_COUNTER_MAX 5
  int registration_accept_sent;
  bool registration_reject_sent;
  bool registration_complete_received;
  guti_t guti;
  bstring amf_msg_out;  // SMF message to be sent within the Registration Accept
                        // message
  amf_registration_request_ies_t* ies;
  amf_ue_ngap_id_t ue_id;
  ksi_t ksi;
  int amf_cause;
};

// typedef struct nas_amf_common_proc_s {
// nas_amf_proc_t amf_proc;
// amf_common_proc_type_t type;
//} nas_amf_common_proc_t;

//===========/end  merged from amf_nas5g_proc.h======

//=========== moved from data.h ============
class nas_amf_smc_proc_t {
 public:
  nas_amf_common_proc_t amf_com_proc;
  nas5g_timer_t T3560; /* Authentication timer         */
#define SECURITY_COUNTER_MAX 5
  amf_ue_ngap_id_t ue_id;
  unsigned int retransmission_count; /* Retransmission counter    */
  int ksi;                           /* NAS key set identifier                */
  int eea;                           /* Replayed 5G encryption algorithms    */
  int eia;                           /* Replayed 5G integrity algorithms     */
  int ucs2;                          /* Replayed Alphabet                     */

  int selected_eea;        /* Selected 5G encryption algorithms    */
  int selected_eia;        /* Selected 5G integrity algorithms     */
  int saved_selected_eea;  /* Previous selected 5G encryption algorithms    */
  int saved_selected_eia;  /* Previous selected 5G integrity algorithms     */
  int saved_eksi;          /* Previous ksi     */
  uint16_t saved_overflow; /* Previous dl_count overflow     */
  uint8_t saved_seq_num;   /* Previous dl_count seq_num     */
  amf_sc_type_t saved_sc_type;
  bool notify_failure; /* Indicates whether the identification
                        * procedure failure shall be notified
                        * to the ongoing EMM procedure */
  bool is_new;         /* new security context for SMC header type */
  bool imeisv_request;
  void amf_ctx_clear_security(amf_context_t* ctxt) __attribute__((nonnull))
  __attribute__((flatten));
  void amf_ctx_set_security_eksi(amf_context_t* ctxt, ksi_t eksi);
  void amf_ctx_set_security_type(amf_context_t* ctxt, amf_sc_type_t sc_type);
};

typedef struct nas_amf_info_proc_s {
  nas_amf_common_proc_t amf_com_proc;
} nas_amf_info_proc_t;

nas_5g_auth_info_proc_t* nas_new_5gcn_auth_info_procedure(
    amf_context_t* const amf_context);

amf_fsm_state_t amf_fsm_get_state(amf_context_t* amf_context);

void amf_free_send_authentication_request(AuthenticationRequestMsg* amf_msg);

// nas_amf_common_proc_t* get_nas5g_common_procedure(
//    const amf_context_t* ctxt, amf_common_proc_type_t proc_type);
nas_amf_smc_proc_t* get_nas5g_common_procedure_smc(const amf_context_t* ctxt);
nas5g_cn_proc_t* get_nas5g_cn_procedure(
    const amf_context_t* ctxt, cn5g_proc_type_t proc_type);
nas_5g_auth_info_proc_t* get_nas5g_cn_procedure_auth_info(
    const amf_context_t* ctxt);

void amf_app_state_free_ue_context(void** ue_context_node);
int amf_proc_security_mode_control(
    amf_context_t* amf_ctx, nas_amf_specific_proc_t* amf_specific_proc,
    ksi_t ksi, success_cb_t success, failure_cb_t failure);

void amf_proc_create_procedure_registration_request(
    ue_m5gmm_context_s* ue_ctx, amf_registration_request_ies_t* ies);

amf_procedures_t* _nas_new_amf_procedures(amf_context_t* amf_context);
int amf_proc_amf_informtion(ue_m5gmm_context_s* ue_amf_ctx);

// UE originated deregistration procedures
int amf_proc_deregistration_request(
    amf_ue_ngap_id_t ue_id, amf_deregistration_request_ies_t* params);
int amf_app_handle_deregistration_req(amf_ue_ngap_id_t ue_id);
void amf_remove_ue_context(
    amf_ue_context_t* amf_ue_context_p, ue_m5gmm_context_s* ue_context_p);
// TODO Refactor required: resource setup request and release procedure
int pdu_session_resource_setup_request(
    ue_m5gmm_context_s* ue_context, amf_ue_ngap_id_t amf_ue_ngap_id);
void amf_app_handle_resource_setup_response(
    itti_ngap_pdusessionresource_setup_rsp_t session_seup_resp);
int pdu_session_resource_release_request(
    ue_m5gmm_context_s* ue_context, amf_ue_ngap_id_t amf_ue_ngap_id);
void amf_app_handle_resource_release_response(
    itti_ngap_pdusessionresource_rel_rsp_t session_rel_resp);

}  // namespace magma5g
#endif
