/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

/*****************************************************************************

Source      as_message.h

Version     0.1

Date        2012/10/18

Product     NAS stack

Subsystem   Application Programming Interface

Author      Frederic Maurel, Lionel GAUTHIER

Description Defines the messages supported by the Access Stratum sublayer
        protocol (usually RRC and S1AP for E-UTRAN) and functions used
        to encode and decode

*****************************************************************************/
#ifndef FILE_AS_MESSAGE_H_SEEN
#define FILE_AS_MESSAGE_H_SEEN

#include "nas/commonDef.h"
#include "nas/networkDef.h"
#include "3gpp_24.007.h"
#include "3gpp_24.301.h"
#include "3gpp_23.003.h"
#include "3gpp_36.331.h"
#include "3gpp_36.401.h"
#include "TrackingAreaIdentity.h"
#include "common_types.h"

/****************************************************************************/
/*********************  G L O B A L    C O N S T A N T S  *******************/
/****************************************************************************/

/*
 * --------------------------------------------------------------------------
 *              Access Stratum message types
 * --------------------------------------------------------------------------
 */
#define AS_REQUEST 0x0100
#define AS_RESPONSE 0x0200
#define AS_INDICATION 0x0400
#define AS_CONFIRM 0x0800

/*
 * --------------------------------------------------------------------------
 *          Access Stratum message identifiers
 * --------------------------------------------------------------------------
 */

/* Broadcast information */
#define AS_BROADCAST_INFO 0x01
#define AS_BROADCAST_INFO_IND (AS_BROADCAST_INFO | AS_INDICATION)

/* Paging information */
#define AS_PAGING 0x03
#define AS_PAGING_REQ (AS_PAGING | AS_REQUEST)

/* NAS signalling connection establishment */
#define AS_NAS_ESTABLISH 0x04
#define AS_NAS_ESTABLISH_REQ (AS_NAS_ESTABLISH | AS_REQUEST)
#define AS_NAS_ESTABLISH_IND (AS_NAS_ESTABLISH | AS_INDICATION)
#define AS_NAS_ESTABLISH_RSP (AS_NAS_ESTABLISH | AS_RESPONSE)
#define AS_NAS_ESTABLISH_CNF (AS_NAS_ESTABLISH | AS_CONFIRM)

/* NAS signalling connection release */
#define AS_NAS_RELEASE 0x05
#define AS_NAS_RELEASE_REQ (AS_NAS_RELEASE | AS_REQUEST)
#define AS_NAS_RELEASE_IND (AS_NAS_RELEASE | AS_INDICATION)

/* Uplink information transfer */
#define AS_UL_INFO_TRANSFER 0x06
#define AS_UL_INFO_TRANSFER_REQ (AS_UL_INFO_TRANSFER | AS_REQUEST)
#define AS_UL_INFO_TRANSFER_CNF (AS_UL_INFO_TRANSFER | AS_CONFIRM)
#define AS_UL_INFO_TRANSFER_IND (AS_UL_INFO_TRANSFER | AS_INDICATION)

/* Downlink information transfer */
#define AS_DL_INFO_TRANSFER 0x07
#define AS_DL_INFO_TRANSFER_REQ (AS_DL_INFO_TRANSFER | AS_REQUEST)
#define AS_DL_INFO_TRANSFER_CNF (AS_DL_INFO_TRANSFER | AS_CONFIRM)
#define AS_DL_INFO_TRANSFER_IND (AS_DL_INFO_TRANSFER | AS_INDICATION)

/* Radio Access Bearer establishment */
#define AS_ACTIVATE_BEARER_CONTEXT 0x08
#define AS_ACTIVATE_BEARER_CONTEXT_REQ (AS_ACTIVATE_BEARER_CONTEXT | AS_REQUEST)

/* Radio Access Bearer release */
#define AS_RAB_RELEASE 0x09
#define AS_RAB_RELEASE_REQ (AS_RAB_RELEASE | AS_REQUEST)
#define AS_RAB_RELEASE_IND (AS_RAB_RELEASE | AS_INDICATION)

/* Deactivate Bearer */
#define AS_DEACTIVATE_BEARER_CONTEXT 0xa
#define AS_DEACTIVATE_BEARER_CONTEXT_REQ                                       \
  (AS_DEACTIVATE_BEARER_CONTEXT | AS_REQUEST)

/* NAS Cause */
typedef enum nas_cause_s {
  NAS_CAUSE_IMSI_UNKNOWN_IN_HSS = EMM_CAUSE_IMSI_UNKNOWN_IN_HSS,
  NAS_CAUSE_ILLEGAL_UE          = EMM_CAUSE_ILLEGAL_UE,
  NAS_CAUSE_ILLEGAL_ME          = EMM_CAUSE_ILLEGAL_ME,
  NAS_CAUSE_UE_IDENTITY_CANT_BE_DERIVED_BY_NW =
      EMM_CAUSE_UE_IDENTITY_CANT_BE_DERIVED_BY_NW,
  NAS_CAUSE_IMPLICITLY_DETACHED      = EMM_CAUSE_IMPLICITLY_DETACHED,
  NAS_CAUSE_IMEI_NOT_ACCEPTED        = EMM_CAUSE_IMEI_NOT_ACCEPTED,
  NAS_CAUSE_EPS_SERVICES_NOT_ALLOWED = EMM_CAUSE_EPS_NOT_ALLOWED,
  NAS_CAUSE_EPS_SERVICES_AND_NON_EPS_SERVICES_NOT_ALLOWED =
      EMM_CAUSE_BOTH_NOT_ALLOWED,
  NAS_CAUSE_PLMN_NOT_ALLOWED          = EMM_CAUSE_PLMN_NOT_ALLOWED,
  NAS_CAUSE_TRACKING_AREA_NOT_ALLOWED = EMM_CAUSE_TA_NOT_ALLOWED,
  NAS_CAUSE_ROAMING_NOT_ALLOWED_IN_THIS_TRACKING_AREA =
      EMM_CAUSE_ROAMING_NOT_ALLOWED,
  NAS_CAUSE_EPS_NOT_ALLOWED_IN_PLMN = EMM_CAUSE_EPS_NOT_ALLOWED_IN_PLMN,
  NAS_CAUSE_NO_SUITABLE_CELLS_IN_TRACKING_AREA = EMM_CAUSE_NO_SUITABLE_CELLS,
  NAS_CAUSE_CSG_NOT_AUTHORIZED                 = EMM_CAUSE_CSG_NOT_AUTHORIZED,
  NAS_CAUSE_NOT_AUTHORIZED_IN_PLMN    = EMM_CAUSE_NOT_AUTHORIZED_IN_PLMN,
  NAS_CAUSE_NO_EPS_BEARER_CTX_ACTIVE  = EMM_CAUSE_NO_EPS_BEARER_CTX_ACTIVE,
  NAS_CAUSE_MSC_NOT_REACHABLE         = EMM_CAUSE_MSC_NOT_REACHABLE,
  NAS_CAUSE_NETWORK_FAILURE           = EMM_CAUSE_NETWORK_FAILURE,
  NAS_CAUSE_CS_DOMAIN_NOT_AVAILABLE   = EMM_CAUSE_CS_DOMAIN_NOT_AVAILABLE,
  NAS_CAUSE_ESM_FAILURE               = EMM_CAUSE_ESM_FAILURE,
  NAS_CAUSE__MAC_FAILURE              = EMM_CAUSE_MAC_FAILURE,
  NAS_CAUSE_SYNCH_FAILURE             = EMM_CAUSE_SYNCH_FAILURE,
  NAS_CAUSE_CONGESTION                = EMM_CAUSE_CONGESTION,
  NAS_CAUSE_SECURITY_MISMATCH         = EMM_CAUSE_UE_SECURITY_MISMATCH,
  NAS_CAUSE_SECURITY_MODE_REJECTED    = EMM_CAUSE_SECURITY_MODE_REJECTED,
  NAS_CAUSE_NON_EPS_AUTH_UNACCEPTABLE = EMM_CAUSE_NON_EPS_AUTH_UNACCEPTABLE,
  NAS_CAUSE_CS_SERVICE_NOT_AVAILABLE  = EMM_CAUSE_CS_SERVICE_NOT_AVAILABLE
} nas_cause_t;

/*
 * --------------------------------------------------------------------------
 *          Access Stratum message global parameters
 * --------------------------------------------------------------------------
 */

/* Error code */
typedef enum nas_error_code_s {
  AS_SUCCESS = 1,          /* Success code, transaction is going on    */
  AS_TERMINATED_NAS,       /* Transaction terminated by NAS        */
  AS_TERMINATED_AS,        /* Transaction terminated by AS         */
  AS_NON_DELIVERED_DUE_HO, /* Failure code                 */
  AS_FAILURE               /* Failure code, stand also for lower
                            * layer failure AS_LOWER_LAYER_FAILURE */
} nas_error_code_t;

/* Core network domain */
typedef enum core_network_s {
  AS_PS = 1, /* Packet-Switched  */
  AS_CS      /* Circuit-Switched */
} core_network_t;

/* SAE Temporary Mobile Subscriber Identity */
typedef struct as_stmsi_s {
  uint8_t mme_code; /* MME code that allocated the GUTI     */
  uint32_t m_tmsi;  /* M-Temporary Mobile Subscriber Identity   */
} as_stmsi_t;

/* Radio Access Bearer identity */
typedef uint8_t as_rab_id_t;

/****************************************************************************/
/************************  G L O B A L    T Y P E S  ************************/
/****************************************************************************/

/*
 * --------------------------------------------------------------------------
 *              Broadcast information
 * --------------------------------------------------------------------------
 */

/*
 * AS->NAS - Broadcast information indication
 * AS may asynchronously report to NAS available PLMNs within specific
 * location area
 */
typedef struct broadcast_info_ind_s {
#define PLMN_LIST_MAX_SIZE 6
  PLMN_LIST_T(PLMN_LIST_MAX_SIZE) plmn_ids; /* List of PLMN identifiers */
  eci_t cell_id; /* Identity of the cell serving the listed PLMNs */
  tac_t tac;     /* Code of the tracking area the cell belongs to */
} broadcast_info_ind_t;

/*
 * --------------------------------------------------------------------------
 *     Cell information relevant for cell selection processing
 * --------------------------------------------------------------------------
 */

/* Radio access technologies supported by the network */
#define AS_GSM (1 << NET_ACCESS_GSM)
#define AS_COMPACT (1 << NET_ACCESS_COMPACT)
#define AS_UTRAN (1 << NET_ACCESS_UTRAN)
#define AS_EGPRS (1 << NET_ACCESS_EGPRS)
#define AS_HSDPA (1 << NET_ACCESS_HSDPA)
#define AS_HSUPA (1 << NET_ACCESS_HSUPA)
#define AS_HSDUPA (1 << NET_ACCESS_HSDUPA)
#define AS_EUTRAN (1 << NET_ACCESS_EUTRAN)

/*
 * NAS->AS - Cell Information request
 * NAS request AS to search for a suitable cell belonging to the selected
 * PLMN to camp on.
 */
typedef struct cell_info_req_s {
  plmn_t plmn_id; /* Selected PLMN identity           */
  uint8_t rat;    /* Bitmap - set of radio access technologies    */
} cell_info_req_t;

/*
 * AS->NAS - Cell Information confirm
 * AS search for a suitable cell and respond to NAS. If found, the cell
 * is selected to camp on.
 */
typedef struct cell_info_cnf_s {
  uint8_t err_code; /* Error code                     */
  eci_t cell_id;    /* Identity of the cell serving the selected PLMN */
  tac_t tac;        /* Code of the tracking area the cell belongs to  */
  AcT_t rat;        /* Radio access technology supported by the cell  */
  uint8_t rsrq;     /* Reference signal received quality         */
  uint8_t rsrp;     /* Reference signal received power       */
} cell_info_cnf_t;

/*
 * AS->NAS - Cell Information indication
 * AS may change cell selection if a more suitable cell is found.
 */
typedef struct cell_info_ind_s {
  eci_t cell_id; /* Identity of the new serving cell      */
  tac_t tac;     /* Code of the tracking area the cell belongs to */
} cell_info_ind_t;

/*
 * --------------------------------------------------------------------------
 *              Paging information
 * --------------------------------------------------------------------------
 */

/* Paging cause */
typedef enum paging_cause_s {
  AS_CONNECTION_ESTABLISH, /* Establish NAS signalling connection  */
  AS_EPS_ATTACH,           /* Perform local detach and initiate EPS
                            * attach procedure         */
  AS_CS_FALLBACK           /* Inititate CS fallback procedure  */
} paging_cause_t;

/*
 * NAS->AS - Paging Information request
 * NAS requests the AS that NAS signalling messages or user data is pending
 * to be sent.
 */
typedef struct paging_req_s {
  s_tmsi_t s_tmsi;   /* UE identity                  */
  uint8_t cn_domain; /* Core network domain              */
} paging_req_t;

/*
 * AS->NAS - Paging Information indication
 * AS reports to the NAS that appropriate procedure has to be initiated.
 */
typedef struct paging_ind_s {
  paging_cause_t cause; /* Paging cause                 */
} paging_ind_t;

/*
 * --------------------------------------------------------------------------
 *          NAS signalling connection establishment
 * --------------------------------------------------------------------------
 */

/* Cause of RRC connection establishment */
typedef enum as_cause_s {
  AS_CAUSE_UNKNOWN   = 0,
  AS_CAUSE_EMERGENCY = EMERGENCY,
  AS_CAUSE_HIGH_PRIO = HIGH_PRIORITY_ACCESS,
  AS_CAUSE_MT_ACCESS = MT_ACCESS,
  AS_CAUSE_MO_SIGNAL = MO_SIGNALLING,
  AS_CAUSE_MO_DATA   = MO_DATA,
  AS_CAUSE_V1020     = DELAY_TOLERANT_ACCESS_V1020
} as_cause_t;

/* Type of the call associated to the RRC connection establishment */
typedef enum as_call_type_s {
  AS_TYPE_ORIGINATING_SIGNAL = NET_ESTABLISH_TYPE_ORIGINATING_SIGNAL,
  AS_TYPE_EMERGENCY_CALLS    = NET_ESTABLISH_TYPE_EMERGENCY_CALLS,
  AS_TYPE_ORIGINATING_CALLS  = NET_ESTABLISH_TYPE_ORIGINATING_CALLS,
  AS_TYPE_TERMINATING_CALLS  = NET_ESTABLISH_TYPE_TERMINATING_CALLS,
  AS_TYPE_MO_CS_FALLBACK     = NET_ESTABLISH_TYPE_MO_CS_FALLBACK
} as_call_type_t;

/*
 * NAS->AS - NAS signalling connection establishment request
 * NAS requests the AS to perform the RRC connection establishment procedure
 * to transfer initial NAS message to the network while UE is in IDLE mode.
 */
typedef struct nas_establish_req_s {
  as_cause_t cause;        /* RRC connection establishment cause   */
  as_call_type_t type;     /* RRC associated call type             */
  s_tmsi_t s_tmsi;         /* UE identity                          */
  plmn_t plmn_id;          /* Selected PLMN identity               */
  bstring initial_nas_msg; /* Initial NAS message to transfer      */
} nas_establish_req_t;

/*
 * AS->NAS - NAS signalling connection establishment indication
 * AS transfers the initial NAS message to the NAS.
 */
typedef struct nas_establish_ind_s {
  mme_ue_s1ap_id_t ue_id; /* UE lower layer identifier               */
  tai_t tai; /* Indicating the Tracking Area from which the UE has sent the NAS
                message.                         */
  ecgi_t ecgi;         /* Indicating the cell from which the UE has sent the NAS
                          message.                         */
  as_cause_t as_cause; /* Establishment cause                     */
  s_tmsi_t s_tmsi;     /* UE identity optional field, if not present, value is
                          NOT_A_S_TMSI */
  bstring initial_nas_msg; /* Initial NAS message to transfer         */
} nas_establish_ind_t;

/*
 * NAS->AS - NAS signalling connection establishment response
 * NAS responds to the AS that initial answer message has to be provided to
 * the UE.
 */
typedef struct nas_establish_rsp_s {
  mme_ue_s1ap_id_t ue_id;    /* UE lower layer identifier   */
  s_tmsi_t s_tmsi;           /* UE identity                 */
  nas_error_code_t err_code; /* Transaction status          */
  bstring nas_msg;           /* NAS message to transfer     */
  uint32_t nas_ul_count;     /* UL NAS COUNT                */
  uint16_t selected_encryption_algorithm;
  uint16_t selected_integrity_algorithm;
  uint8_t csfb_response;
#define SERVICE_TYPE_PRESENT (1 << 0)
  uint8_t presencemask; /* Indicates the presence of some params like service
                           type */
  uint8_t service_type;
} nas_establish_rsp_t;

/*
 * AS->NAS - NAS signalling connection establishment confirm
 * AS transfers the initial answer message to the NAS.
 */
typedef struct nas_establish_cnf_s {
  mme_ue_s1ap_id_t ue_id;    /* UE lower layer identifier   */
  nas_error_code_t err_code; /* Transaction status          */
  bstring nas_msg;           /* NAS message to transfer     */
  uint32_t ul_nas_count;
  uint16_t selected_encryption_algorithm;
  uint16_t selected_integrity_algorithm;
} nas_establish_cnf_t;

/*
 * --------------------------------------------------------------------------
 *          NAS signalling connection release
 * --------------------------------------------------------------------------
 */

/* Release cause */
typedef enum release_cause_s {
  AS_AUTHENTICATION_FAILURE = 1, /* Authentication procedure failed   */
  AS_DETACH                      /* Detach requested                  */
} release_cause_t;

/*
 * NAS->AS - NAS signalling connection release request
 * NAS requests the termination of the connection with the UE.
 */
typedef struct nas_release_req_s {
  mme_ue_s1ap_id_t ue_id; /* UE lower layer identifier    */
  s_tmsi_t s_tmsi;        /* UE identity                  */
  release_cause_t cause;  /* Release cause                */
} nas_release_req_t;

/*
 * AS->NAS - NAS signalling connection release indication
 * AS reports that connection has been terminated by the network.
 */
typedef struct nas_release_ind_s {
  release_cause_t cause; /* Release cause            */
} nas_release_ind_t;

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
typedef struct ul_info_transfer_req_s {
  mme_ue_s1ap_id_t ue_id; /* UE lower layer identifier        */
  s_tmsi_t s_tmsi;        /* UE identity              */
  bstring nas_msg;        /* Uplink NAS message           */
} ul_info_transfer_req_t;

/*
 * AS->NAS - Uplink data transfer confirm
 * AS immediately notifies the NAS whether uplink information has been
 * successfully sent to the network or not.
 */
typedef struct ul_info_transfer_cnf_s {
  mme_ue_s1ap_id_t ue_id;    /* UE lower layer identifier        */
  nas_error_code_t err_code; /* Transaction status               */
} ul_info_transfer_cnf_t;

/*
 * AS->NAS - Uplink data transfer indication
 * AS delivers the uplink information message to the NAS that operates
 * at the network side.
 */
typedef struct ul_info_transfer_ind_s {
  mme_ue_s1ap_id_t ue_id; /* UE lower layer identifier        */
  bstring nas_msg;        /* Uplink NAS message           */
} ul_info_transfer_ind_t;

/*
 * NAS->AS - Downlink data transfer request
 * NAS requests the AS to transfer downlink information to the NAS that
 * operates at the UE side.
 */
typedef struct dl_info_transfer_req_s {
  mme_ue_s1ap_id_t ue_id;    /* UE lower layer identifier        */
  s_tmsi_t s_tmsi;           /* UE identity              */
  bstring nas_msg;           /* Uplink NAS message           */
  nas_error_code_t err_code; /* Transaction status               */
} dl_info_transfer_req_t;

/*
 * AS->NAS - Downlink data transfer confirm
 * AS immediately notifies the NAS whether downlink information has been
 * successfully sent to the network or not.
 */
typedef ul_info_transfer_cnf_t dl_info_transfer_cnf_t;

/*
 * AS->NAS - Downlink data transfer indication
 * AS delivers the downlink information message to the NAS that operates
 * at the UE side.
 */
typedef ul_info_transfer_ind_t dl_info_transfer_ind_t;

/*
 * --------------------------------------------------------------------------
 *          Radio Access Bearer establishment
 * --------------------------------------------------------------------------
 */

/*
 * NAS->AS - Radio access bearer establishment request
 * NAS requests the AS to allocate transmission resources to radio access
 * bearer initialized at the network side.
 */
typedef struct activate_bearer_context_req_s {
  mme_ue_s1ap_id_t ue_id; /* UE lower layer identifier        */
  ebi_t ebi;              /* EPS bearer id    */
  bitrate_t mbr_dl;
  bitrate_t mbr_ul;
  bitrate_t gbr_dl;
  bitrate_t gbr_ul;
  bstring nas_msg; /* NAS message to transfer     */
} activate_bearer_context_req_t;

/*
 * AS->NAS - Radio access bearer establishment indication
 * AS notifies the NAS that specific radio access bearer has to be setup.
 */
typedef struct rab_establish_ind_s {
  mme_ue_s1ap_id_t ue_id; /* UE lower layer identifier        */
  ebi_t ebi;              /* EPS bearer id    */
} rab_establish_ind_t;

/*
 * NAS->AS - Radio access bearer establishment response
 * NAS responds to AS whether the specified radio access bearer has been
 * successfully setup or not.
 */
typedef struct rab_establish_rsp_s {
  mme_ue_s1ap_id_t ue_id;    /* UE lower layer identifier        */
  ebi_t ebi;                 /* EPS bearer id    */
  nas_error_code_t err_code; /* Transaction status               */
} rab_establish_rsp_t;

/*
 * AS->NAS - Radio access bearer establishment confirm
 * AS notifies NAS whether the specified radio access bearer has been
 * successfully setup at the UE side or not.
 */
typedef struct rab_establish_cnf_s {
  mme_ue_s1ap_id_t ue_id;    /* UE lower layer identifier        */
  ebi_t ebi;                 /* EPS bearer id    */
  nas_error_code_t err_code; /* Transaction status               */
} rab_establish_cnf_t;

/*
 * --------------------------------------------------------------------------
 *              Radio Access Bearer release
 * --------------------------------------------------------------------------
 */

/*
 * NAS->AS - Radio access bearer release request
 * NAS requests the AS to release transmission resources previously allocated
 * to specific radio access bearer at the network side.
 */
typedef struct rab_release_req_s {
  s_tmsi_t s_tmsi;    /* UE identity                      */
  as_rab_id_t rab_id; /* Radio access bearer identity     */
} rab_release_req_t;

/*
 * AS->NAS - Radio access bearer release indication
 * AS notifies NAS that specific radio access bearer has been released.
 */
typedef struct rab_release_ind_s {
  as_rab_id_t rab_id; /* Radio access bearer identity     */
} rab_release_ind_t;

/*
 * NAS->AS - Deactivate EPS Bearer context request
 * NAS requests the AS to deactivate bearer
 */
typedef struct deactivate_bearer_context_req_s {
  mme_ue_s1ap_id_t ue_id; /* UE lower layer identifier        */
  ebi_t ebi;              /* EPS bearer id    */
  bstring nas_msg;        /* NAS message to transfer     */
} deactivate_bearer_context_req_t;

/*
 * --------------------------------------------------------------------------
 *  Structure of the AS messages handled by the network sublayer
 * --------------------------------------------------------------------------
 */
typedef struct as_message_s {
  uint16_t msg_id;
  union {
    broadcast_info_ind_t broadcast_info_ind;
    cell_info_req_t cell_info_req;
    cell_info_cnf_t cell_info_cnf;
    cell_info_ind_t cell_info_ind;
    paging_req_t paging_req;
    paging_ind_t paging_ind;
    nas_establish_req_t nas_establish_req;
    nas_establish_ind_t nas_establish_ind;
    nas_establish_rsp_t nas_establish_rsp;
    nas_establish_cnf_t nas_establish_cnf;
    nas_release_req_t nas_release_req;
    nas_release_ind_t nas_release_ind;
    ul_info_transfer_req_t ul_info_transfer_req;
    ul_info_transfer_cnf_t ul_info_transfer_cnf;
    ul_info_transfer_ind_t ul_info_transfer_ind;
    dl_info_transfer_req_t dl_info_transfer_req;
    dl_info_transfer_cnf_t dl_info_transfer_cnf;
    dl_info_transfer_ind_t dl_info_transfer_ind;
    activate_bearer_context_req_t activate_bearer_context_req;
    rab_establish_ind_t rab_establish_ind;
    rab_establish_rsp_t rab_establish_rsp;
    rab_establish_cnf_t rab_establish_cnf;
    rab_release_req_t rab_release_req;
    rab_release_ind_t rab_release_ind;
    deactivate_bearer_context_req_t deactivate_bearer_context_req;
  } msg;
} as_message_t;

/****************************************************************************/
/********************  G L O B A L    V A R I A B L E S  ********************/
/****************************************************************************/

/****************************************************************************/
/******************  E X P O R T E D    F U N C T I O N S  ******************/
/****************************************************************************/

int as_message_decode(const char* buffer, as_message_t* msg, size_t length);

int as_message_encode(char* buffer, as_message_t* msg, size_t length);

/* Implemented in the network_api.c body file */
int as_message_send(as_message_t* as_msg);

#endif /* FILE_AS_MESSAGE_H_SEEN*/
