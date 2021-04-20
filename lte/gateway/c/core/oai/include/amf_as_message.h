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
#include "TrackingAreaIdentity.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_38.331.h"
#include "3gpp_38.413.h"
#include "3gpp_36.331.h"
#include "common_types.h"

/*
 * --------------------------------------
 * Network connection establishment type
 * --------------------------------------
 */
#define AMF_ESTABLISH_TYPE_ORIGINATING_SIGNAL 0x10
#define AMF_ESTABLISH_TYPE_EMERGENCY_CALLS 0x20
#define AMF_ESTABLISH_TYPE_ORIGINATING_CALLS 0x30
#define AMF_ESTABLISH_TYPE_TERMINATING_CALLS 0x40

/*
 * --------------------------------------------------------------------------
 *              Access Stratum message types
 * --------------------------------------------------------------------------
 */
#define AS_REQUEST_ 0x0100
#define AS_RESPONSE_ 0x0200
#define AS_INDICATION_ 0x0400
#define AS_CONFIRM_ 0x0800

/* NAS signaling connection establishment */
#define AS_NAS_ESTABLISH_ 0x04
#define AS_NAS_ESTABLISH_REQ_ (AS_NAS_ESTABLISH_ | AS_REQUEST_)
#define AS_NAS_ESTABLISH_IND_ (AS_NAS_ESTABLISH_ | AS_INDICATION_)
#define AS_NAS_ESTABLISH_RSP_ (AS_NAS_ESTABLISH_ | AS_RESPONSE_)
#define AS_NAS_ESTABLISH_CNF_ (AS_NAS_ESTABLISH_ | AS_CONFIRM_)

/* NAS signaling connection release */
#define AS_NAS_RELEASE_ 0x05
#define AS_NAS_RELEASE_REQ_ (AS_NAS_RELEASE_ | AS_REQUEST_)
#define AS_NAS_RELEASE_IND_ (AS_NAS_RELEASE_ | AS_INDICATION_)

/* Uplink information transfer */
#define AS_UL_INFO_TRANSFER_ 0x06
#define AS_UL_INFO_TRANSFER_REQ_ (AS_UL_INFO_TRANSFER_ | AS_REQUEST_)
#define AS_UL_INFO_TRANSFER_CNF_ (AS_UL_INFO_TRANSFER_ | AS_CONFIRM_)
#define AS_UL_INFO_TRANSFER_IND_ (AS_UL_INFO_TRANSFER_ | AS_INDICATION_)

/* Downlink information transfer */
#define AS_DL_INFO_TRANSFER_ 0x07
#define AS_DL_INFO_TRANSFER_REQ_ (AS_DL_INFO_TRANSFER_ | AS_REQUEST_)
#define AS_DL_INFO_TRANSFER_CNF_ (AS_DL_INFO_TRANSFER_ | AS_CONFIRM_)
#define AS_DL_INFO_TRANSFER_IND_ (AS_DL_INFO_TRANSFER_ | AS_INDICATION_)

typedef struct m5g_broadcast_info_ind_s {
#define PLMN_LIST_MAX_SIZE 6
  PLMN_LIST_T(PLMN_LIST_MAX_SIZE) plmn_ids; /* List of PLMN identifiers */
  eci_t cell_id; /* Identity of the cell serving the listed PLMNs */
  tac_t tac;     /* Code of the tracking area the cell belongs to */
} m5g_broadcast_info_ind_t;

typedef struct m5g_cell_info_req_s {
  plmn_t plmn_id; /* Selected PLMN identity           */
  uint8_t rat;    /* Bitmap - set of radio access technologies    */
} m5g_cell_info_req_t;

/*
 * AS->NAS - Cell Information confirm
 * AS search for a suitable cell and respond to NAS. If found, the cell
 * is selected to camp on.
 */
typedef struct m5g_cell_info_cnf_s {
  uint8_t err_code; /* Error code                     */
  eci_t cell_id;    /* Identity of the cell serving the selected PLMN */
  tac_t tac;        /* Code of the tracking area the cell belongs to  */
  AcT_t rat;        /* Radio access technology supported by the cell  */
  uint8_t rsrq;     /* Reference signal received quality         */
  uint8_t rsrp;     /* Reference signal received power       */
} m5g_cell_info_cnf_t;

typedef struct m5g_cell_info_ind_s {
  eci_t cell_id; /* Identity of the new serving cell      */
  tac_t tac;     /* Code of the tracking area the cell belongs to */
} m5g_cell_info_ind_t;

/*
 * NAS->AS - Paging Information request
 * NAS requests the AS that NAS signaling messages or user data is pending
 * to be sent.
 */
typedef struct m5g_paging_req_s {
  s_tmsi_m5_t s_tmsi; /* UE identity                  */
  uint8_t cn_domain;  /* Core network domain              */
} m5g_paging_req_t;

/* Type of the call associated to the RRC connection establishment */
typedef enum amf_as_call_type_s {
  AMF_AS_TYPE_ORIGINATING_SIGNAL = AMF_ESTABLISH_TYPE_ORIGINATING_SIGNAL,
  AMF_AS_TYPE_EMERGENCY_CALLS    = AMF_ESTABLISH_TYPE_EMERGENCY_CALLS,
  AMF_AS_TYPE_ORIGINATING_CALLS  = AMF_ESTABLISH_TYPE_ORIGINATING_CALLS,
  AMF_AS_TYPE_TERMINATING_CALLS  = AMF_ESTABLISH_TYPE_TERMINATING_CALLS,

} amf_as_call_type_t;

/* from 6.1.2 ETSI TS 138 331 V15.8.0 (2020-01)*/
/* m5G RRC Cause */
/*{emergency, highPriorityAccess, mt-Access, mo-signaling,
  mo-Data, mo-VoiceCall, mo-VideoCall, mo-SMS, mps-PriorityAccess,
  mcs-PriorityAccess, spare6, spare5, spare4, spare3, spare2, spare1}*/
/* --------------------------------------------------------------------------
 *          NAS signaling connection establishment
 * --------------------------------------------------------------------------
 */

/* Cause of RRC connection establishment */
typedef enum m5g_as_cause_s {
  M5G_AS_CAUSE_UNKNOWN       = 0,
  M5G_AS_CAUSE_EMERGENCY     = M5G_EMERGENCY,
  M5G_AS_CAUSE_HIGH_PRIO     = M5G_HIGH_PRIORITY_ACCESS,
  M5G_AS_CAUSE_MT_ACCESS     = M5G_MT_ACCESS,
  M5G_AS_CAUSE_MO_SIGNAL     = M5G_MO_SIGNALLING,
  M5G_AS_CAUSE_MO_DATA       = M5G_MO_DATA,
  M5G_AS_CAUSE_MO_VOICE_CALL = M5G_MO_VOICE_CALL,
  M5G_AS_CAUSE_MO_VIDEO_CALL = M5G_MO_VIDEOCALL,
  M5G_AS_CAUSE_MO_SMS        = M5G_MO_SMS,
  M5G_AS_CAUSE_MPS_PRIO      = M5G_MPS_PRIORITYACCESS,
  M5G_AS_CAUSE_MCS_PRIO      = M5G_MCS_PRIORITYACCESS,
  M5G_AS_CAUSE_V1020         = DELAY_TOLERANT_ACCESS_V1020
} m5g_as_cause_t;

/*
 * NAS->AS - NAS signaling connection establishment request
 * NAS requests the AS to perform the RRC connection establishment procedure
 * to transfer initial NAS message to the network while UE is in IDLE mode.
 */

typedef struct nas5g_establish_req_s {
  m5g_as_cause_t cause;    /* RRC connection establishment cause   */
  amf_as_call_type_t type; /* RRC associated call type             */
  s_tmsi_m5_t s_tmsi;      /* UE identity                          */
  plmn_t plmn_id;          /* Selected PLMN identity               */
  bstring initial_nas_msg; /* Initial NAS message to transfer      */
} nas5g_establish_req_t;

/*
 * AS->NAS - NAS signaling connection establishment indication
 * AS transfers the initial NAS message to the NAS.
 */
typedef struct nas5g_establish_ind_s {
  amf_ue_ngap_id_ty ue_id; /* UE lower layer identifier               */
  tai_t tai; /* Indicating the Tracking Area from which the UE has sent the NAS
                message.                         */
  ecgi_t ecgi; /* Indicating the cell from which the UE has sent the NAS
                  message.                         */
  m5g_as_cause_t as_cause; /* Establishment cause                     */
  s_tmsi_m5_t s_tmsi; /* UE identity optional field, if not present, value is
                      NOT_A_S_TMSI */
  bstring initial_nas_msg; /* Initial NAS message to transfer         */
} nas5g_establish_ind_t;

/*
 * --------------------------------------------------------------------------
 *          Access Stratum message global parameters
 * --------------------------------------------------------------------------
 */

/* Error code */
typedef enum nas5g_error_code_s {
  M5G_AS_SUCCESS = 1,          /* Success code, transaction is going on    */
  M5G_AS_TERMINATED_NAS,       /* Transaction terminated by NAS        */
  M5G_AS_TERMINATED_AS,        /* Transaction terminated by AS         */
  M5G_AS_NON_DELIVERED_DUE_HO, /* Failure code                 */
  M5G_AS_FAILURE               /* Failure code, stand also for lower
                                * layer failure AS_LOWER_LAYER_FAILURE */
} nas5g_error_code_t;

/*
 * NAS->AS - NAS signaling connection establishment response
 * NAS responds to the AS that initial answer message has to be provided to
 * the UE.
 */
typedef struct nas5g_establish_rsp_s {
  amf_ue_ngap_id_ty ue_id;     /* UE lower layer identifier   */
  s_tmsi_m5_t s_tmsi;          /* UE identity                 */
  nas5g_error_code_t err_code; /* Transaction status          */
  bstring nas_msg;             /* NAS message to transfer     */
  uint32_t nas_ul_count;       /* UL NAS COUNT                */
  uint16_t selected_encryption_algorithm;
  uint16_t selected_integrity_algorithm;
#define M5G_SERVICE_TYPE_PRESENT (1 << 0)
  uint8_t presencemask; /* Indicates the presence of some params like service
                           type */
  uint8_t m5g_service_type;
} nas5g_establish_rsp_t;

/*
 * AS->NAS - NAS signaling connection establishment confirm
 * AS transfers the initial answer message to the NAS.
 */
typedef struct nas5g_establish_cnf_s {
  amf_ue_ngap_id_ty ue_id;     /* UE lower layer identifier   */
  nas5g_error_code_t err_code; /* Transaction status          */
  bstring nas_msg;             /* NAS message to transfer     */
  uint32_t ul_nas_count;
  uint16_t selected_encryption_algorithm;
  uint16_t selected_integrity_algorithm;
} nas5g_establish_cnf_t;

/*
 * --------------------------------------------------------------------------
 *          NAS signaling connection release
 * --------------------------------------------------------------------------
 */

/* Release cause */
typedef enum m5g_release_cause_s {
  M5G_AS_AUTHENTICATION_FAILURE = 1, /* Authentication procedure failed   */
  M5G_AS_DEREGISTRATION /* Deregistration requested                  */
} m5g_release_cause_t;

/*
 * NAS->AS - NAS signaling connection release request
 * NAS requests the termination of the connection with the UE.
 */
typedef struct nas5g_release_req_s {
  amf_ue_ngap_id_ty ue_id;   /* UE lower layer identifier    */
  s_tmsi_m5_t s_tmsi;        /* UE identity                  */
  m5g_release_cause_t cause; /* Release cause                */
} nas5g_release_req_t;

/*
 * AS->NAS - NAS signaling connection release indication
 * AS reports that connection has been terminated by the network.
 */
typedef struct nas5g_release_ind_s {
  m5g_release_cause_t cause; /* Release cause            */
} nas5g_release_ind_t;

/*
 * --------------------------------------------------------------------------
 *              NAS information transfer
 * --------------------------------------------------------------------------
 */

/*
 * NAS->AS - Uplink data transfer request
 * NAS requests the AS to transfer uplink information to the NAS that
 * operates at the network side.
 */
typedef struct m5g_ul_info_transfer_req_s {
  amf_ue_ngap_id_ty ue_id; /* UE lower layer identifier        */
  s_tmsi_m5_t s_tmsi;      /* UE identity              */
  bstring nas_msg;         /* Uplink NAS message           */
} m5g_ul_info_transfer_req_t;

/*
 * AS->NAS - Uplink data transfer confirm
 * AS immediately notifies the NAS whether uplink information has been
 * successfully sent to the network or not.
 */
typedef struct m5g_ul_info_transfer_cnf_s {
  amf_ue_ngap_id_ty ue_id;     /* UE lower layer identifier        */
  nas5g_error_code_t err_code; /* Transaction status               */
} m5g_ul_info_transfer_cnf_t;
/*
 * AS->NAS - Uplink data transfer indication
 * AS delivers the uplink information message to the NAS that operates
 * at the network side.
 */
typedef struct m5g_ul_info_transfer_ind_s {
  amf_ue_ngap_id_ty ue_id; /* UE lower layer identifier        */
  bstring nas_msg;         /* Uplink NAS message           */
} m5g_ul_info_transfer_ind_t;

/*
 * NAS->AS - Downlink data transfer request
 * NAS requests the AS to transfer downlink information to the NAS that
 * operates at the UE side.
 */
typedef struct m5g_dl_info_transfer_req_s {
  amf_ue_ngap_id_ty ue_id;     /* UE lower layer identifier        */
  s_tmsi_m5_t s_tmsi;          /* UE identity              */
  bstring nas_msg;             /* Uplink NAS message           */
  nas5g_error_code_t err_code; /* Transaction status               */
} m5g_dl_info_transfer_req_t;

/*
 * AS->NAS - Downlink data transfer confirm
 * AS immediately notifies the NAS whether downlink information has been
 * successfully sent to the network or not.
 */
typedef m5g_ul_info_transfer_cnf_t m5g_dl_info_transfer_cnf_t;

/*
 * AS->NAS - Downlink data transfer indication
 * AS delivers the downlink information message to the NAS that operates
 * at the UE side.
 */
typedef m5g_ul_info_transfer_ind_t m5g_dl_info_transfer_ind_t;

/*
 * --------------------------------------------------------------------------
 *          Radio Access Pdu Session establishment
 * --------------------------------------------------------------------------
 */

/*
 * NAS->AS - Radio access Pdu session establishment request
 * NAS requests the AS to allocate transmission resources to radio access
 * Pdu session initialized at the network side.
 */
typedef struct activate_pdusession_context_req_s {
  amf_ue_ngap_id_ty ue_id; /* UE lower layer identifier        */
  psi_t psi;               /* Pdu session id    */
  bitrate_t mbr_dl;
  bitrate_t mbr_ul;
  bitrate_t gbr_dl;
  bitrate_t gbr_ul;
  bstring nas_msg; /* NAS message to transfer     */
} activate_pdusession_context_req_t;

typedef struct amf_as_message_s {
  uint16_t msg_id;
  union {
    m5g_broadcast_info_ind_t broadcast_info_ind;
    m5g_cell_info_req_t cell_info_req;
    m5g_cell_info_cnf_t cell_info_cnf;
    m5g_cell_info_ind_t cell_info_ind;
    m5g_paging_req_t paging_req;
    nas5g_establish_req_t nas_establish_req;
    nas5g_establish_ind_t nas_establish_ind;
    nas5g_establish_rsp_t nas_establish_rsp;
    nas5g_establish_cnf_t nas_establish_cnf;
    nas5g_release_req_t nas_release_req;
    nas5g_release_ind_t nas_release_ind;
    m5g_ul_info_transfer_req_t ul_info_transfer_req;
    m5g_ul_info_transfer_cnf_t ul_info_transfer_cnf;
    m5g_ul_info_transfer_ind_t ul_info_transfer_ind;
    m5g_dl_info_transfer_req_t dl_info_transfer_req;
    m5g_dl_info_transfer_cnf_t dl_info_transfer_cnf;
    m5g_dl_info_transfer_ind_t dl_info_transfer_ind;
    activate_pdusession_context_req_t activate_pdusession_context_req;
  } __attribute__((__packed__)) msg;
} amf_as_message_t;
