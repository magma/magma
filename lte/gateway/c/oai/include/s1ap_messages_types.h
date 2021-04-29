/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice,
 * this list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE
 * ARE DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE
 * LIABLE FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR
 * CONSEQUENTIAL DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF
 * SUBSTITUTE GOODS OR SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS
 * INTERRUPTION) HOWEVER CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN
 * CONTRACT, STRICT LIABILITY, OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE)
 * ARISING IN ANY WAY OUT OF THE USE OF THIS SOFTWARE, EVEN IF ADVISED OF THE
 * POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are
 * those of the authors and should not be interpreted as representing official
 * policies, either expressed or implied, of the FreeBSD Project.
 */
/*! \file s1ap_messages_types.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_S1AP_MESSAGES_TYPES_SEEN
#define FILE_S1AP_MESSAGES_TYPES_SEEN

#include "3gpp_24.008.h"
#include "3gpp_36.401.h"
#include "3gpp_36.413.h"
#include "3gpp_36.331.h"
#include "3gpp_23.003.h"
#include "TrackingAreaIdentity.h"
#include "nas/securityDef.h"

#include "S1ap_Source-ToTarget-TransparentContainer.h"
#include "S1ap_HandoverType.h"

#define S1AP_ENB_DEREGISTERED_IND(mSGpTR)                                      \
  (mSGpTR)->ittiMsg.s1ap_eNB_deregistered_ind
#define S1AP_ENB_INITIATED_RESET_REQ(mSGpTR)                                   \
  (mSGpTR)->ittiMsg.s1ap_enb_initiated_reset_req
#define S1AP_ENB_INITIATED_RESET_ACK(mSGpTR)                                   \
  (mSGpTR)->ittiMsg.s1ap_enb_initiated_reset_ack
#define S1AP_UE_CONTEXT_RELEASE_REQ(mSGpTR)                                    \
  (mSGpTR)->ittiMsg.s1ap_ue_context_release_req
#define S1AP_UE_CONTEXT_RELEASE_COMMAND(mSGpTR)                                \
  (mSGpTR)->ittiMsg.s1ap_ue_context_release_command
#define S1AP_UE_CONTEXT_RELEASE_COMPLETE(mSGpTR)                               \
  (mSGpTR)->ittiMsg.s1ap_ue_context_release_complete
#define S1AP_UE_CONTEXT_MODIFICATION_REQUEST(mSGpTR)                           \
  (mSGpTR)->ittiMsg.s1ap_ue_context_mod_request
#define S1AP_UE_CONTEXT_MODIFICATION_RESPONSE(mSGpTR)                          \
  (mSGpTR)->ittiMsg.s1ap_ue_context_mod_response
#define S1AP_UE_CONTEXT_MODIFICATION_FAILURE(mSGpTR)                           \
  (mSGpTR)->ittiMsg.s1ap_ue_context_mod_failure
#define S1AP_E_RAB_SETUP_REQ(mSGpTR) (mSGpTR)->ittiMsg.s1ap_e_rab_setup_req
#define S1AP_E_RAB_SETUP_RSP(mSGpTR) (mSGpTR)->ittiMsg.s1ap_e_rab_setup_rsp
#define S1AP_E_RAB_MODIFICATION_IND(mSGpTR)                                    \
  (mSGpTR)->ittiMsg.s1ap_e_rab_modification_ind
#define S1AP_E_RAB_MODIFICATION_CNF(mSGpTR)                                    \
  (mSGpTR)->ittiMsg.s1ap_e_rab_modification_cnf
#define S1AP_INITIAL_UE_MESSAGE(mSGpTR)                                        \
  (mSGpTR)->ittiMsg.s1ap_initial_ue_message
#define S1AP_NAS_DL_DATA_REQ(mSGpTR) (mSGpTR)->ittiMsg.s1ap_nas_dl_data_req
#define S1AP_PAGING_REQUEST(mSGpTR) (mSGpTR)->ittiMsg.s1ap_paging_request
#define S1AP_E_RAB_REL_CMD(mSGpTR) (mSGpTR)->ittiMsg.s1ap_e_rab_rel_cmd
#define S1AP_E_RAB_REL_RSP(mSGpTR) (mSGpTR)->ittiMsg.s1ap_e_rab_rel_rsp
#define S1AP_PATH_SWITCH_REQUEST(mSGpTR)                                       \
  (mSGpTR)->ittiMsg.s1ap_path_switch_request
#define S1AP_PATH_SWITCH_REQUEST_ACK(mSGpTR)                                   \
  (mSGpTR)->ittiMsg.s1ap_path_switch_request_ack
#define S1AP_PATH_SWITCH_REQUEST_FAILURE(mSGpTR)                               \
  (mSGpTR)->ittiMsg.s1ap_path_switch_request_failure
#define S1AP_REMOVE_STALE_UE_CONTEXT(mSGpTR)                                   \
  (mSGpTR)->ittiMsg.s1ap_remove_stale_ue_context
#define S1AP_HANDOVER_REQUIRED(mSGpTR) (mSGpTR)->ittiMsg.s1ap_handover_required
#define S1AP_HANDOVER_REQUEST_ACK(mSGpTR)                                      \
  (mSGpTR)->ittiMsg.s1ap_handover_request_ack
#define S1AP_HANDOVER_NOTIFY(mSGpTR) (mSGpTR)->ittiMsg.s1ap_handover_notify

// NOT a ITTI message
typedef struct s1ap_initial_ue_message_s {
  enb_ue_s1ap_id_t enb_ue_s1ap_id : 24;
  ecgi_t e_utran_cgi;
} s1ap_initial_ue_message_t;

typedef struct itti_s1ap_ue_context_mod_req_s {
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id : 24;
/* Use presence mask to identify presence of optional fields */
#define S1AP_UE_CONTEXT_MOD_LAI_PRESENT (1 << 0)
#define S1AP_UE_CONTEXT_MOD_CSFB_INDICATOR_PRESENT (1 << 1)
#define S1AP_UE_CONTEXT_MOD_UE_AMBR_INDICATOR_PRESENT (1 << 2)
  uint8_t presencemask;
  lai_t lai;
  int cs_fallback_indicator;
  ambr_t ue_ambr;
} itti_s1ap_ue_context_mod_req_t;

typedef struct itti_s1ap_ue_context_mod_resp_s {
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id : 24;
} itti_s1ap_ue_context_mod_resp_t;

typedef struct itti_s1ap_ue_context_mod_resp_fail_s {
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id : 24;
  int64_t cause;
} itti_s1ap_ue_context_mod_resp_fail_t;

typedef struct itti_s1ap_initial_ctxt_setup_req_s {
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id : 24;

  /* Key eNB */
  uint8_t kenb[32];

  ambr_t ambr;
  ambr_t apn_ambr;

  /* EPS bearer ID */
  unsigned ebi : 4;

  /* QoS */
  qci_t qci;
  priority_level_t prio_level;
  pre_emption_vulnerability_t pre_emp_vulnerability;
  pre_emption_capability_t pre_emp_capability;

  /* S-GW TEID for user-plane */
  teid_t teid;
  /* S-GW IP address for User-Plane */
  ip_address_t s_gw_address;
} itti_s1ap_initial_ctxt_setup_req_t;

typedef struct itti_s1ap_ue_cap_ind_s {
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id : 24;
  uint8_t* radio_capabilities;
  size_t radio_capabilities_length;
} itti_s1ap_ue_cap_ind_t;

#define S1AP_ITTI_UE_PER_DEREGISTER_MESSAGE 128
typedef struct itti_s1ap_eNB_deregistered_ind_s {
  uint16_t nb_ue_to_deregister;
  enb_ue_s1ap_id_t enb_ue_s1ap_id[S1AP_ITTI_UE_PER_DEREGISTER_MESSAGE];
  mme_ue_s1ap_id_t mme_ue_s1ap_id[S1AP_ITTI_UE_PER_DEREGISTER_MESSAGE];
  uint32_t enb_id;
} itti_s1ap_eNB_deregistered_ind_t;

typedef enum s1ap_reset_type_e {
  RESET_ALL = 0,
  RESET_PARTIAL
} s1ap_reset_type_t;

typedef struct s1_sig_conn_id_s {
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id;
} s1_sig_conn_id_t;

typedef struct itti_s1ap_enb_initiated_reset_req_s {
  uint32_t sctp_assoc_id;
  uint16_t sctp_stream_id;
  uint32_t enb_id;
  s1ap_reset_type_t s1ap_reset_type;
  uint32_t num_ue;
  s1_sig_conn_id_t* ue_to_reset_list;
} itti_s1ap_enb_initiated_reset_req_t;

typedef struct itti_s1ap_enb_initiated_reset_ack_s {
  uint32_t sctp_assoc_id;
  uint16_t sctp_stream_id;
  s1ap_reset_type_t s1ap_reset_type;
  uint32_t num_ue;
  s1_sig_conn_id_t* ue_to_reset_list;
} itti_s1ap_enb_initiated_reset_ack_t;

// List of possible causes for MME generated UE context release command towards
// eNB
enum s1cause {
  S1AP_INVALID_CAUSE = 0,
  S1AP_NAS_NORMAL_RELEASE,
  S1AP_NAS_DETACH,
  S1AP_RADIO_EUTRAN_GENERATED_REASON,
  S1AP_RADIO_UNKNOWN_E_RAB_ID,
  S1AP_IMPLICIT_CONTEXT_RELEASE,
  S1AP_INITIAL_CONTEXT_SETUP_FAILED,
  S1AP_SCTP_SHUTDOWN_OR_RESET,
  S1AP_INVALID_ENB_ID,
  S1AP_INVALID_MME_UE_S1AP_ID,
  S1AP_CSFB_TRIGGERED,
  S1AP_NAS_UE_NOT_AVAILABLE_FOR_PS,
  S1AP_SYSTEM_FAILURE,
  S1AP_RADIO_MULTIPLE_E_RAB_ID,
  S1AP_NAS_MME_OFFLOADING,
  S1AP_NAS_MME_PENDING_OFFLOADING
};
typedef struct itti_s1ap_ue_context_release_command_s {
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id : 24;
  enum s1cause cause;
} itti_s1ap_ue_context_release_command_t;

typedef struct itti_s1ap_ue_context_release_req_s {
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id : 24;
  uint32_t enb_id;
  enum s1cause relCause;
  S1ap_Cause_t cause;
} itti_s1ap_ue_context_release_req_t;

typedef struct itti_s1ap_dl_nas_data_req_s {
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id : 24;
  bstring nas_msg; /* Downlink NAS message             */
} itti_s1ap_nas_dl_data_req_t;

typedef struct itti_s1ap_ue_context_release_complete_s {
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id : 24;
} itti_s1ap_ue_context_release_complete_t;

typedef struct itti_s1ap_remove_stale_ue_context_s {
  uint32_t enb_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id;
} itti_s1ap_remove_stale_ue_context_t;

typedef enum s1ap_csfb_indicator_e {
  CSFB_REQUIRED,
  CSFB_HIGH_PRIORITY
} s1ap_csfb_indicator_t;

typedef enum s1ap_cn_domain_e { CN_DOMAIN_PS, CN_DOMAIN_CS } s1ap_cn_domain_t;

typedef struct itti_s1ap_paging_request_s {
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  uint8_t imsi_length;
  tmsi_t m_tmsi;
  mme_code_t mme_code;
  uint32_t sctp_assoc_id;
#define S1AP_PAGING_ID_IMSI 0X0
#define S1AP_PAGING_ID_STMSI 0X1
  uint8_t paging_id;
  s1ap_cn_domain_t domain_indicator;
  uint8_t tai_list_count;
  paging_tai_list_t paging_tai_list[TRACKING_AREA_IDENTITY_MAX_NUM_OF_TAIS];
} itti_s1ap_paging_request_t;

typedef struct itti_s1ap_initial_ue_message_s {
  sctp_assoc_id_t sctp_assoc_id;  // key stored in MME_APP for MME_APP forward
                                  // NAS response to S1AP
  uint32_t enb_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id;
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  bstring nas;
  tai_t tai; /* Indicating the Tracking Area from which the UE has sent the NAS
                message. */
  ecgi_t ecgi; /* Indicating the cell from which the UE has sent the NAS
                  message. */
  rrc_establishment_cause_t rrc_establishment_cause; /* Establishment cause */

  bool is_s_tmsi_valid;
  bool is_csg_id_valid;
  bool is_gummei_valid;
  s_tmsi_t opt_s_tmsi;
  csg_id_t opt_csg_id;
  gummei_t opt_gummei;
  // void                opt_cell_access_mode;
  // void                opt_cell_gw_transport_address;
  // void                opt_relay_node_indicator;
  /* Transparent message from s1ap to be forwarded to MME_APP or
   * to S1AP if connection establishment is rejected by NAS.
   */
  s1ap_initial_ue_message_t transparent;
} itti_s1ap_initial_ue_message_t;

typedef struct itti_s1ap_e_rab_setup_req_s {
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id;

  // Applicable for non-GBR E-RABs
  bool ue_aggregate_maximum_bit_rate_present;
  ue_aggregate_maximum_bit_rate_t ue_aggregate_maximum_bit_rate;

  // E-RAB to Be Setup List
  e_rab_to_be_setup_list_t e_rab_to_be_setup_list;

} itti_s1ap_e_rab_setup_req_t;

typedef struct itti_s1ap_e_rab_setup_rsp_s {
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id;

  // E-RAB to Be Setup List
  e_rab_setup_list_t e_rab_setup_list;

  // Optional
  e_rab_list_t e_rab_failed_to_setup_list;

} itti_s1ap_e_rab_setup_rsp_t;

typedef struct itti_s1ap_e_rab_rel_cmd_s {
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id;

  // Applicable for non-GBR E-RABs
  bool ue_aggregate_maximum_bit_rate_present;
  ue_aggregate_maximum_bit_rate_t ue_aggregate_maximum_bit_rate;

  // E-RAB to Be Released List
  e_rab_list_t e_rab_to_be_rel_list;
  bstring nas_pdu;
} itti_s1ap_e_rab_rel_cmd_t;

typedef struct itti_s1ap_e_rab_rel_rsp_s {
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id;

  // E-RAB to Be Setup List
  e_rab_list_t e_rab_rel_list;

  // Optional
  e_rab_list_t e_rab_failed_to_rel_list;

} itti_s1ap_e_rab_rel_rsp_t;

typedef struct itti_s1ap_path_switch_request_s {
  uint32_t sctp_assoc_id;
  uint32_t enb_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id : 24;
  e_rab_to_be_switched_in_downlink_list_t e_rab_to_be_switched_dl_list;
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  tai_t tai;
  ecgi_t ecgi;
  uint16_t encryption_algorithm_capabilities;
  uint16_t integrity_algorithm_capabilities;
} itti_s1ap_path_switch_request_t;

typedef struct itti_s1ap_path_switch_request_ack_s {
  uint32_t sctp_assoc_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id : 24;
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  // Security key
  uint8_t nh[AUTH_NEXT_HOP_SIZE]; /* Next Hop security key*/
  uint8_t ncc;                    /* next hop chaining count */
} itti_s1ap_path_switch_request_ack_t;

typedef struct itti_s1ap_path_switch_request_failure_s {
  uint32_t sctp_assoc_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id : 24;
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
} itti_s1ap_path_switch_request_failure_t;

typedef struct itti_s1ap_e_rab_modification_ind_s {
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id;
  // E-RAB to be Modified List
  e_rab_to_be_modified_bearer_mod_ind_list_t e_rab_to_be_modified_list;
  e_rab_not_to_be_modified_bearer_mod_ind_list_t e_rab_not_to_be_modified_list;
  // Optional
} itti_s1ap_e_rab_modification_ind_t;

typedef struct itti_s1ap_e_rab_modification_cnf_s {
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  enb_ue_s1ap_id_t enb_ue_s1ap_id;
  // E-RAB Modify List
  e_rab_modify_bearer_mod_conf_list_t e_rab_modify_list;
  // Optional
  e_rab_list_t e_rab_failed_to_modify_list;
} itti_s1ap_e_rab_modification_cnf_t;

typedef struct itti_s1ap_handover_required_s {
  uint32_t sctp_assoc_id;
  uint32_t enb_id;
  S1ap_Cause_t cause;
  S1ap_HandoverType_t handover_type;
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  bstring src_tgt_container;
} itti_s1ap_handover_required_t;

typedef struct itti_s1ap_handover_request_ack_s {
  uint32_t source_assoc_id;
  uint32_t target_assoc_id;
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  enb_ue_s1ap_id_t src_enb_ue_s1ap_id;
  enb_ue_s1ap_id_t tgt_enb_ue_s1ap_id;
  uint32_t source_enb_id;
  uint32_t target_enb_id;
  S1ap_HandoverType_t handover_type;
  bstring tgt_src_container;
} itti_s1ap_handover_request_ack_t;

typedef struct itti_s1ap_handover_notify_s {
  mme_ue_s1ap_id_t mme_ue_s1ap_id;
  uint32_t target_enb_id;
  uint32_t target_sctp_assoc_id;
  ecgi_t ecgi;
  enb_ue_s1ap_id_t target_enb_ue_s1ap_id;
  e_rab_admitted_list_t e_rab_admitted_list;
} itti_s1ap_handover_notify_t;
#endif /* FILE_S1AP_MESSAGES_TYPES_SEEN */
