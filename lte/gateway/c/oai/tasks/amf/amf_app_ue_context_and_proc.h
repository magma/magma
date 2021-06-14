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

#pragma once
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
#include <vector>
#include "amf_fsm.h"
#include "amf_data.h"
#include "amf_smfDefs.h"
#include "AmfMessage.h"
#include "amf_app_messages_types.h"
#include "M5GRegistrationRequest.h"
#include "ngap_messages_types.h"

// NAS messages
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

namespace magma5g {
#define NAS5G_TIMER_INACTIVE_ID (-1)
#define SECURITY_MODE_TIMER_EXPIRY_MSECS 6000
#define AUTHENTICATION_TIMER_EXPIRY_MSECS 6000
#define REGISTRATION_ACCEPT_TIMER_EXPIRY_MSECS 6000
#define IDENTITY_TIMER_EXPIRY_MSECS 6000
struct amf_procedures_t;
/*
 * Timer identifier returned when in inactive state (timer is stopped or has
 * failed to be started)
 */
#define AMF_APP_TIMER_INACTIVE_ID (-1)
#define AMF_APP_DELTA_T3512_REACHABILITY_TIMER 4            // in minutes
#define AMF_APP_DELTA_REACHABILITY_IMPLICIT_DETACH_TIMER 0  // in minutes
#define AMF_APP_INITIAL_CONTEXT_SETUP_RSP_TIMER_VALUE 2     // In seconds
#define AMF_APP_UE_CONTEXT_MODIFICATION_TIMER_VALUE 2       // In seconds
#define AMF_APP_PAGING_RESPONSE_TIMER_VALUE 4               // In seconds
#define AMF_APP_ULR_RESPONSE_TIMER_VALUE 3                  // In seconds
#define NAS5G_SECURITY_ALGORITHMS_MINIMUM_LENGTH 1
#define NAS5G_SECURITY_ALGORITHMS_MAXIMUM_LENGTH 2
#define NAS5G_MESSAGE_CONTAINER_MAXIMUM_LENGTH 253

#define IDENTITY_TIMER_EXPIRY_MSECS 6000
#define AUTHENTICATION_TIMER_EXPIRY_MSECS 6000
#define SECURITY_MODE_TIMER_EXPIRY_MSECS 6000
#define REGISTRATION_ACCEPT_TIMER_EXPIRY_MSECS 6000
#define PAGING_TIMER_EXPIRY_MSECS 4000

#define MAX_PAGING_RETRY_COUNT 1
// Header length boundaries of 5GS Mobility Management messages
#define AMF_HEADER_LENGTH sizeof(amf_msg_header)

#define N1_SM_INFO 0x1
#define AMBR_LEN 6
#define PDU_ESTAB_ACCPET_PAYLOAD_CONTAINER_LEN 30
#define PDU_ESTAB_ACCPET_NAS_PDU_LEN 39
#define SSC_MODE_ONE 0x1
#define PDU_ADDR_IPV4_LEN 0x4
#define PDU_ADDR_TYPE 0X1
#define GNB_IPV4_ADDR_LEN 4
#define GNB_TEID_LEN 4

#define NON_AMF_3GPP_ACCESS 2
#define AMF_3GPP_ACCESS_AND_NON_AMF_3GPP_ACCESS 3

// Global variable and needed to increment based on nas procedures
static uint64_t nas_puid = 1;

// Timer structure
typedef struct amf_app_timer_s {
  long id;  /* The timer identifier                 */
  long sec; /* The timer interval value in seconds  */
} amf_app_timer_t;

/* PDU session resource request and release NGAP messages
 * Request and response
 */
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
  up_transport_layer_information_t up_transport_layer_info;
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

// Structure to handle Resource Release Response from gNB
typedef enum {
  NR,
  E_UTRA,
} rat_type_e;

typedef struct volume_timed_report_item_s {
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

typedef struct pdu_session_resource_setup_req_s {
  /*
   * values: 1-1Kbps, 2- 4Kbps, 3- 16Kbps, 4- 64Kbps 24.501 9.11.4.14
   */
  uint8_t units_for_session;
  ngap_ue_aggregate_maximum_bit_rate_t pdu_aggregate_maximum_bit_rate;
  Ngap_PDUSessionID_t Pdu_Session_ID;  // from NGAP
  Ngap_SNSSAI_t Ngap_s_nssai;          // S-NSSAI as defined in TS 23.003 [23]
  pdu_session_resource_setup_request_transfer_t
      pdu_session_resource_setup_request_transfer;
} pdu_session_resource_setup_req_t;

// GTP tunnel id for UPF and gNB exchange infomration
typedef struct teid_upf_gnb_s {
  uint8_t upf_gtp_teid_ip_addr[16];
  uint8_t upf_gtp_teid[4];
  uint8_t gnb_gtp_teid_ip_addr[16];
  uint8_t gnb_gtp_teid[4];
} teid_upf_gnb_t;

// Data get communicated with SMF and stored for reference
typedef struct smf_proc_data_s {
  PDUSessionIdentityMsg pdu_session_identity;
  PTIMsg pti;
  MessageTypeMsg message_type;
  IntegrityProtMaxDataRateMsg integrity_prot_max_data_rate;
  PDUSessionTypeMsg pdu_session_type;
  SSCModeMsg ssc_mode;
} smf_proc_data_t;

// PDU session context part of AMFContext
typedef struct smf_context_s {
  SMSessionFSMState pdu_session_state;
  uint32_t pdu_session_version;
  uint32_t n_active_pdus;
  bool is_emergency;
  uint8_t dl_ambr_unit;
  uint16_t dl_session_ambr;
  uint8_t ul_ambr_unit;
  uint16_t ul_session_ambr;
  QOSRule qos_rules[1];
  teid_upf_gnb_t gtp_tunnel_id;
  smf_proc_data_t smf_proc_data;

  // Request to gnb on PDU establisment request
  pdu_session_resource_setup_req_t pdu_resource_setup_req;
  pdu_session_resource_setup_rsp_t pdu_resource_setup_rsp;
  pdu_session_resource_to_release_list pdu_resource_release_req;
  pdu_session_resource_to_release_list pdu_resource_release_rsp;
} smf_context_t;

typedef struct paging_context_s {
#define MAX_PAGING_RETRY_COUNT 1
  amf_app_timer_t m5_paging_response_timer;
  uint8_t paging_retx_count;
} paging_context_t;

/*
 * Structure of the AMF context established by core for a particular UE
 * --------------------------------------------------------------------
 */
typedef struct amf_context_s {
  bool is_dynamic;    /* Dynamically allocated context indicator         */
  bool is_registered; /* Registration indicator                            */
  bool is_initial_identity_imsi;  // If the IMSI was used for identification in
                                  // the initial NAS message
  bool is_guti_based_registered;  // For future use
  uint32_t member_present_mask;   /* bitmask, see significance of bits below */
  uint32_t member_valid_mask;     /* bitmask, see significance of bits below */
  uint8_t m5gsregistrationtype;
  std::vector<smf_context_t> smf_ctxt_vector;  // smf contents
  // smf_context_t smf_ctxt_vector;  // smf contents
  amf_procedures_t* amf_procedures;
  imsi_t imsi;     /* The IMSI provided by the UE or the AMF, set valid when
                       identification returns IMSI */
  imsi64_t imsi64; /* The IMSI provided by the UE or the AMF, set valid when
                       identification returns IMSI */
  imsi64_t saved_imsi64; /* Useful for 5.4.2.7.c */
  imei_t imei;           /* The IMEI provided by the UE                     */
  imeisv_t imeisv;       /* The IMEISV provided by the UE                   */
  guti_m5_t m5_guti;     /* The GUTI assigned to the UE                     */
  guti_m5_t m5_old_guti; /* The GUTI assigned to the UE                     */
  ksi_t ksi;             /*key set identifier  */
  drx_parameter_t drx_parameter;
  UESecurityCapabilityMsg ue_sec_capability;
  m5g_auth_vector_t
      _vector[MAX_EPS_AUTH_VECTORS]; /* 5GMM authentication vector */
  amf_security_context_t
      _security; /* Current 5GMM security context: The security context which
                    has been activated most recently. Note that a current 5GMM
                    security context originating from either a mapped
                    or native 5GMM security context may exist simultaneously
                    with a native non-current 5GMM security context.*/
  int amf_cause; /* AMF failure cause code                          */
  amf_fsm_state_t amf_fsm_state;
  smf_context_t smf_context;  // Keeps PDU session related info
  void* t3422_arg;
  drx_parameter_t current_drx_parameter; /* stored TAU Request IE Requirement
                                             AMF24.501R15_5.5.3.2.4_4*/
  std::string smf_msg; /* SMF message contained within the initial request*/
  bool is_imsi_only_detach;
  uint8_t reg_id_type;
} amf_context_t;

typedef struct amf_ue_context_s {
  hash_table_uint64_ts_t* imsi_amf_ue_id_htbl;    // data is amf_ue_ngap_id_t
  hash_table_uint64_ts_t* tun11_ue_context_htbl;  // data is amf_ue_ngap_id_t
  hash_table_uint64_ts_t*
      gnb_ue_ngap_id_ue_context_htbl;             // data is amf_ue_ngap_id_t
  obj_hash_table_uint64_t* guti_ue_context_htbl;  // data is amf_ue_ngap_id_t
} amf_ue_context_t;

enum m5gcm_state_t {
  M5GCM_IDLE = 0,
  M5GCM_CONNECTED,
};

/* @ue_m5gmm_context_s
 * @brief Useful parameters to know in AMF application layer.
 */
typedef struct ue_m5gmm_context_s {
  ngap_Cause_t
      ue_context_rel_cause;  // define require for Ngcause in NGAP module
  m5gmm_state_t mm_state;
  m5gcm_state_t ecm_state;
  amf_context_t amf_context;
  teid_t amf_teid_n11;
  // SCTP assoc id
  sctp_assoc_id_t sctp_assoc_id_key;
  // gNB UE NGAP ID,  Unique identity the UE within gNodeB
  gnb_ue_ngap_id_t gnb_ue_ngap_id;
  gnb_ngap_id_key_t gnb_ngap_id_key;
  // AMF UE NGAP ID, Unique identity the UE within AMF
  amf_ue_ngap_id_t amf_ue_ngap_id;
  /* mobile_reachability_timer: Start when UE moves to idle state.
   *             Stop when UE moves to connected state
   */
  amf_app_timer_t m5_mobile_reachability_timer;

  /* implicit_detach_timer: Start at the expiry of Mobile Reachability timer.
   * Stop when UE moves to connected state
   */
  amf_app_timer_t m5_implicit_detach_timer;

  // Initial Context Setup Procedure Guard timer
  amf_app_timer_t m5_initial_context_setup_rsp_timer;

  // UE Context Modification Procedure Guard timer
  amf_app_timer_t m5_ue_context_modification_timer;

  /* Paging Structure */
  paging_context_t paging_context;
  amf_app_timer_t m5_ulr_response_timer;

  // UEContextRequest in  INITIAL UE MESSAGE
  m5g_uecontextrequest_t ue_context_request;
} ue_m5gmm_context_t;

/* Operation on UE context structure
 */
int amf_insert_ue_context(
    amf_ue_ngap_id_t ue_id, amf_ue_context_t* amf_ue_context_p,
    ue_m5gmm_context_s* ue_context_p);
amf_ue_ngap_id_t amf_app_ctx_get_new_ue_id(
    amf_ue_ngap_id_t* amf_app_ue_ngap_id_generator_p);
/* Notify NGAP about the mapping between amf_ue_ngap_id and
 * sctp assoc id + gnb_ue_ngap_id */
void notify_ngap_new_ue_amf_ngap_id_association(
    const ue_m5gmm_context_s* ue_context_p);
// void amf_remove_ue_context(
//    amf_ue_context_t* const amf_ue_context,
//    ue_m5gmm_context_s* const ue_context_p);

ue_m5gmm_context_s* amf_create_new_ue_context(void);
/*Multi PDU Session*/
smf_context_t* amf_insert_smf_context(
    ue_m5gmm_context_s* ue_context, uint8_t pdu_session_id);
smf_context_t* amf_smf_context_exists_pdu_session_id(
    ue_m5gmm_context_s* ue_context, uint8_t id);

// Retrieve required UE context from the respective hash table
amf_context_t* amf_context_get(const amf_ue_ngap_id_t ue_id);
ue_m5gmm_context_s* amf_ue_context_exists_amf_ue_ngap_id(
    const amf_ue_ngap_id_t amf_ue_ngap_id);
ue_m5gmm_context_s* lookup_ue_ctxt_by_imsi(imsi64_t imsi64);
int amf_context_upsert_imsi(amf_context_t* elm) __attribute__((nonnull));

// Set valid imsi
void amf_ctx_set_valid_imsi(
    amf_context_t* ctxt, imsi_t* imsi, const imsi64_t imsi64)
    __attribute__((nonnull)) __attribute__((flatten));

// Set valid attribute
void amf_ctx_set_attribute_valid(
    amf_context_t* ctxt, const uint32_t attribute_bit_pos)
    __attribute__((nonnull)) __attribute__((flatten));

// set attribute present
void amf_ctx_set_attribute_present(
    amf_context_t* ctxt, const int attribute_bit_pos) __attribute__((nonnull))
__attribute__((flatten));

void amf_ctx_clear_attribute_present(
    amf_context_t* const ctxt, const int attribute_bit_pos)
    __attribute__((nonnull)) __attribute__((flatten));

// NAS encode header
typedef struct amf_msg_header_t {
  uint8_t extended_protocol_discriminator;
  uint8_t security_header_type;
  uint8_t message_type;
  uint32_t message_authentication_code;
  uint8_t sequence_number;
} amf_msg_header;

// Release Request routine.
void amf_app_ue_context_release(
    ue_m5gmm_context_s* ue_context_p, ngap_Cause_t cause);

// NAS decode and validaion of IE
typedef struct amf_nas_message_decode_status_s {
  uint8_t integrity_protected_message : 1;
  uint8_t ciphered_message : 1;
  uint8_t mac_matched : 1;
  uint8_t security_context_available : 1;
  int amf_cause;
} amf_nas_message_decode_status_t;

// 5G Mobility Management Messages
union mobility_msg_u {
  RegistrationRequestMsg registrationrequestmsg;
  RegistrationAcceptMsg registrationacceptmsg;
  RegistrationCompleteMsg registrationcompletemsg;
  RegistrationRejectMsg registrationrejectmsg;
  ServiceRequestMsg service_request;
  ServiceAcceptMsg service_accept;
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
  mobility_msg_u() {}
  ~mobility_msg_u() {}
};

// Procedure for NAS5G encoding and decoding
class AMFMsg {
 public:
  amf_msg_header header;
  mobility_msg_u msg;

  AMFMsg() {}
  ~AMFMsg() {}

  int amf_msg_decode_header(
      amf_msg_header* header, const uint8_t* buffer, uint32_t len);
  int amf_msg_encode_header(
      const amf_msg_header* header, uint8_t* buffer, uint32_t len);
  int amf_msg_decode(AMFMsg* msg, uint8_t* buffer, uint32_t len);
  int amf_msg_encode(AMFMsg* msg, uint8_t* buffer, uint32_t len);
};

// union of plain NAS message
typedef struct nas_message_plain_s {
  AMFMsg amf; /* 5GMM Mobility Management messages */
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

typedef enum {
  CN5G_PROC_NONE = 0,
  CN5G_PROC_AUTH_INFO,
} cn5g_proc_type_t;

typedef enum amf_common_proc_type_s {
  AMF_COMM_PROC_NONE = 0,
  AMF_COMM_PROC_AUTH,
  AMF_COMM_PROC_SMC,
  AMF_COMM_PROC_IDENT,
} amf_common_proc_type_t;

enum nas_base_proc_type_t {
  NAS_PROC_TYPE_NONE = 0,
  NAS_PROC_TYPE_AMF,
  NAS_PROC_TYPE_CN,
};

// forward declaration
struct nas5g_base_proc_t;
struct nas_amf_proc_t;
struct nas_amf_registration_proc_t;

// call back routines during procedure handling
typedef int (*success_cb_t)(amf_context_t* amf_ctx);
typedef int (*failure_cb_t)(amf_context_t* amf_ctx);
typedef int (*proc_abort_t)(
    amf_context_t* amf_ctx, nas5g_base_proc_t* nas_proc);
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

// NAS related procedure
struct nas5g_base_proc_t {
  success_cb_t success_notif;
  failure_cb_t failure_notif;
  proc_abort_t abort;
  pdu_in_rej_t fail_in;
  pdu_out_rej_t fail_out;
  time_out_t time_out;
  nas_base_proc_type_t type;  // AMF, SMF, CN
  uint64_t nas_puid;
  nas5g_base_proc_t* parent;
  nas5g_base_proc_t* child;
};

enum nas_amf_proc_type_t {
  NAS_AMF_PROC_TYPE_NONE = 0,
  NAS_AMF_PROC_TYPE_SPECIFIC,
  NAS_AMF_PROC_TYPE_COMMON,
  NAS_AMF_PROC_TYPE_CONN_MNGT,
};

// AMF Specific procedures
struct nas_amf_proc_t {
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

// Deregistration specific elements
typedef enum deregistration_switchoff_e {
  // In the UE to network direction, octate 1, 4th bit
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
  ksi_t ksi;
} amf_deregistration_request_ies_t;

// AMF Specific procedures
typedef struct nas_amf_specific_proc_s {
  nas_amf_proc_t amf_proc;
  amf_specific_proc_type_t type;
} nas_amf_specific_proc_t;

// UL identification routines.
int amf_proc_identification(
    amf_context_t* const amf_context, nas_amf_proc_t* const amf_proc,
    const identity_type2_t type, success_cb_t success, failure_cb_t failure);
int amf_proc_identification_complete(
    const amf_ue_ngap_id_t ue_id, imsi_t* const imsi, imei_t* const imei,
    imeisv_t* const imeisv, uint32_t* const tmsi);

typedef struct nas_amf_auth_proc_s {
  nas_amf_common_proc_t amf_com_proc;
  nas5g_timer_t T3560; /* Authentication timer         */
#define AUTHENTICATION_COUNTER_MAX 5
  unsigned int retransmission_count;
#define EMM_AUTHENTICATION_SYNC_FAILURE_MAX 2
  unsigned int sync_fail_count; /* counter of successive AUTHENTICATION FAILURE
                                   messages 1133                       from the
                                   UE with AMF cause #21 "synch failure" */
  unsigned int mac_fail_count;
  amf_ue_ngap_id_t ue_id;
  bool is_cause_is_registered;  //  could also be done by seeking parent
                                //  procedure
  ksi_t ksi;
  uint8_t rand[AUTH_RAND_SIZE]; /* Random challenge number  */
  uint8_t autn[AUTH_AUTN_SIZE]; /* Authentication token     */
  imsi_t* unchecked_imsi;
  int amf_cause;
} nas_amf_auth_proc_t;

typedef struct nas5g_cn_proc_s {
  nas5g_base_proc_t base_proc;
  cn5g_proc_type_t type;
} nas5g_cn_proc_t;

typedef struct nas5g_cn_procedure_s {
  nas5g_cn_proc_t* proc;
  LIST_ENTRY(nas5g_cn_procedure_s) entries;
} nas5g_cn_procedure_t;

// Clasify all UL NAS messages based on message type
int nas_proc_establish_ind(
    const amf_ue_ngap_id_t ue_id, const bool is_mm_ctx_new,
    const tai_t originating_tai, const ecgi_t ecgi,
    const m5g_rrc_establishment_cause_t as_cause, const s_tmsi_m5_t s_tmsi,
    bstring msg);
// Registration procedure routine
nas_amf_registration_proc_t* get_nas_specific_procedure_registration(
    const amf_context_t* ctxt);
bool is_nas_specific_procedure_registration_running(const amf_context_t* ctxt);

nas_amf_common_proc_t* get_nas5g_common_procedure(
    const amf_context_t* const ctxt, amf_common_proc_type_t proc_type);

// 5G CN Specific procedures
typedef struct nas_amf_common_procedure_s {
  nas_amf_common_proc_t* proc;
  LIST_ENTRY(nas_amf_common_procedure_s) entries;
} nas_amf_common_procedure_t;

// Recheck and change to nas5g, comment
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

/*
0 0 1 initial registration
0 1 0 mobility registration updating
0 1 1 periodic registration updating
1 0 0 emergency registration
*/
enum amf_proc_registration_type_t {
  AMF_REGISTRATION_TYPE_INITIAL = 1,
  AMF_REGISTRATION_TYPE_MOBILITY_UPDATING,
  AMF_REGISTRATION_TYPE_PERODIC_UPDATING,
  AMF_REGISTRATION_TYPE_EMERGENCY,
  AMF_REGISTRATION_TYPE_RESERVED = 7,
};

typedef struct amf_registration_request_ies_s {
  amf_proc_registration_type_t m5gsregistrationtype;
  guti_m5_t* guti;
  imsi_t* imsi;
  imei_t* imei;
  tai_t* last_visited_registered_tai;             // Last visited registered TAI
  ue_network_capability_t ue_network_capability;  // UE security capability
  drx_parameter_t* drx_parameter;  // Requested DRX parameters during paging
  amf_nas_message_decode_status_t decode_status;
} amf_registration_request_ies_t;

struct amf_procedures_t {
  nas_amf_specific_proc_t* amf_specific_proc;
  LIST_HEAD(nas_amf_common_procedures_head_s, nas_amf_common_procedure_s)
  amf_common_procs;
  LIST_HEAD(nas5g_cn_procedures_head_s, nas5g_cn_procedure_s)
  cn_procs;  // triggered by AMF
};

struct nas_amf_registration_proc_t {
  nas_amf_specific_proc_t amf_spec_proc;
#define REGISTRATION_COUNTER_MAX 5
  unsigned int retransmission_count;
  struct nas5g_timer_s T3550;  // AMF message retransmission timer
  bstring amf_msg_out;  // SMF message to be sent within the Registration Accept
                        // message
  amf_registration_request_ies_t* ies;
  amf_ue_ngap_id_t ue_id;
  ksi_t ksi;
  int amf_cause;
  int registration_accept_sent;
};
// NAS security related IEs
class nas_amf_smc_proc_t {
 public:
  nas_amf_common_proc_t amf_com_proc;
  nas5g_timer_t T3560; /* Authentication timer         */
#define SECURITY_COUNTER_MAX 5
  amf_ue_ngap_id_t ue_id;
  unsigned int retransmission_count;  // Retransmission counter
  int ksi;                            // NAS key set identifier
  int eea;                            // Replayed 5G encryption algorithms
  int eia;                            // Replayed 5G integrity algorithms
  int ucs2;                           // Replayed Alphabet
  int selected_eea;                   // Selected 5G encryption algorithms
  int selected_eia;                   // Selected 5G integrity algorithms
  int saved_selected_eea;   // Previous selected 5G encryption algorithms
  int saved_selected_eia;   // Previous selected 5G integrity algorithms
  int saved_eksi;           // Previous ksi
  uint16_t saved_overflow;  // Previous dl_count overflow
  uint8_t saved_seq_num;    // Previous dl_count seq_num
  amf_sc_type_t saved_sc_type;
  bool is_new;  // new security context for SMC header type
  bool imeisv_request;
  void amf_ctx_clear_security(amf_context_t* ctxt) __attribute__((nonnull))
  __attribute__((flatten));
  void amf_ctx_set_security_eksi(amf_context_t* ctxt, ksi_t eksi);
  void amf_ctx_set_security_type(amf_context_t* ctxt, amf_sc_type_t sc_type);
};

nas_amf_smc_proc_t* get_nas5g_common_procedure_smc(const amf_context_t* ctxt);

void amf_app_state_free_ue_context(void** ue_context_node);
int amf_proc_security_mode_control(
    amf_context_t* amf_ctx, nas_amf_specific_proc_t* amf_specific_proc,
    ksi_t ksi, success_cb_t success, failure_cb_t failure);
void amf_proc_create_procedure_registration_request(
    ue_m5gmm_context_s* ue_ctx, amf_registration_request_ies_t* ies);
amf_procedures_t* nas_new_amf_procedures(amf_context_t* amf_context);
int amf_proc_amf_information(ue_m5gmm_context_s* ue_amf_ctx);
int amf_send_registration_accept(amf_context_t* amf_context);

int amf_send_registration_accept(amf_context_t* amf_context);
// UE originated deregistration procedures
int amf_proc_deregistration_request(
    amf_ue_ngap_id_t ue_id, amf_deregistration_request_ies_t* params);
int amf_app_handle_deregistration_req(amf_ue_ngap_id_t ue_id);
void amf_remove_ue_context(
    amf_ue_context_t* amf_ue_context_p, ue_m5gmm_context_s* ue_context_p);
void amf_smf_context_cleanup_pdu_session(ue_m5gmm_context_s* ue_context);

// PDU session related communication to gNB
int pdu_session_resource_setup_request(
    ue_m5gmm_context_s* ue_context, amf_ue_ngap_id_t amf_ue_ngap_id,
    smf_context_t*);
void amf_app_handle_resource_setup_response(
    itti_ngap_pdusessionresource_setup_rsp_t session_seup_resp);
int pdu_session_resource_release_request(
    ue_m5gmm_context_s* ue_context, amf_ue_ngap_id_t amf_ue_ngap_id,
    smf_context_t* smf_ctx);
void amf_app_handle_resource_release_response(
    itti_ngap_pdusessionresource_rel_rsp_t session_rel_resp);
void amf_app_handle_cm_idle_on_ue_context_release(
    itti_ngap_ue_context_release_req_t cm_idle_req);
// Handle UE CONTEXT RELEASE COMMAND in DL to NGAP
void ue_context_release_command(
    amf_ue_ngap_id_t amf_ue_ngap_id, gnb_ue_ngap_id_t gnb_ue_ngap_id,
    Ngcause ng_cause);

// NAS5G encode and decode routines with security header support
int nas5g_message_decode(
    const unsigned char* const buffer, amf_nas_message_t* msg, uint32_t length,
    void* security, amf_nas_message_decode_status_t* status);

int nas5g_message_encode(
    unsigned char* buffer, const amf_nas_message_t* const msg, uint32_t length,
    void* security);

int amf_registration_run_procedure(amf_context_t* amf_context);
int amf_proc_registration_complete(
    amf_ue_ngap_id_t ue_id, bstring smf_msg_p, int amf_cause,
    const amf_nas_message_decode_status_t status);

// Finite state machine handlers
int ue_state_handle_message_initial(
    m5gmm_state_t cur_state, int event, SMSessionFSMState session_state,
    ue_m5gmm_context_s* ue_m5gmm_context, amf_context_t* amf_context);
int ue_state_handle_message_reg_conn(
    m5gmm_state_t, int, SMSessionFSMState, ue_m5gmm_context_s*,
    amf_ue_ngap_id_t, bstring, int, amf_nas_message_decode_status_t);
int ue_state_handle_message_dereg(
    m5gmm_state_t, int event, SMSessionFSMState, ue_m5gmm_context_s*,
    amf_ue_ngap_id_t);
int pdu_state_handle_message(
    m5gmm_state_t, int event, SMSessionFSMState session_state,
    ue_m5gmm_context_s*, amf_smf_t, char*,
    itti_n11_create_pdu_session_response_t*, uint32_t);
nas_amf_ident_proc_t* get_5g_nas_common_procedure_identification(
    const amf_context_t* ctxt);
void amf_delete_registration_proc(amf_context_t* amf_txt);
void amf_delete_registration_ies(amf_registration_request_ies_t** ies);
void amf_delete_child_procedures(
    amf_context_t* amf_txt, struct nas5g_base_proc_t* const parent_proc);
void amf_delete_common_procedure(
    amf_context_t* amf_ctx, nas_amf_common_proc_t** proc);
void delete_wrapper(void** ptr);
}  // namespace magma5g
