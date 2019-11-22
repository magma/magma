/*
 * Copyright (c) 2015, EURECOM (www.eurecom.fr)
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 * 1. Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 * 2. Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS" AND
 * ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE IMPLIED
 * WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT OWNER OR CONTRIBUTORS BE LIABLE FOR
 * ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL DAMAGES
 * (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR SERVICES;
 * LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER CAUSED AND
 * ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY, OR TORT
 * (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE OF THIS
 * SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 *
 * The views and conclusions contained in the software and documentation are those
 * of the authors and should not be interpreted as representing official policies,
 * either expressed or implied, of the FreeBSD Project.
 */

/*! \file nas_messages_types.h
  \brief
  \author Sebastien ROUX, Lionel Gauthier
  \company Eurecom
  \email: lionel.gauthier@eurecom.fr
*/

#ifndef FILE_NAS_MESSAGES_TYPES_SEEN
#define FILE_NAS_MESSAGES_TYPES_SEEN

#include <stdint.h>

#include "3gpp_23.003.h"
#include "3gpp_29.274.h"
#include "nas/as_message.h"
#include "common_ies.h"
#include "nas/networkDef.h"

#define NAS_BEARER_PARAM(mSGpTR) (mSGpTR)->ittiMsg.nas_bearer_param
#define NAS_AUTHENTICATION_REQ(mSGpTR) (mSGpTR)->ittiMsg.nas_auth_req
#define NAS_AUTHENTICATION_PARAM_REQ(mSGpTR)                                   \
  (mSGpTR)->ittiMsg.nas_auth_param_req
#define NAS_SGS_DETACH_REQ(mSGpTR) (mSGpTR)->ittiMsg.nas_sgs_detach_req
#define NAS_DETACH_REQ(mSGpTR) (mSGpTR)->ittiMsg.nas_detach_req
#define NAS_ERAB_SETUP_REQ(mSGpTR) (mSGpTR)->ittiMsg.itti_erab_setup_req
#define NAS_IMPLICIT_DETACH_UE_IND(mSGpTR)                                     \
  (mSGpTR)->ittiMsg.nas_implicit_detach_ue_ind
#define NAS_NW_INITIATED_DETACH_UE_REQ(mSGpTR)                                 \
  (mSGpTR)->ittiMsg.nas_nw_initiated_detach_ue_req
#define NAS_CS_DOMAIN_LOCATION_UPDATE_REQ(mSGpTR)                              \
  (mSGpTR)->ittiMsg.nas_cs_domain_location_update_req
#define NAS_CS_DOMAIN_LOCATION_UPDATE_ACC(mSGpTR)                              \
  (mSGpTR)->ittiMsg.nas_cs_domain_location_update_acc
#define NAS_CS_DOMAIN_LOCATION_UPDATE_FAIL(mSGpTR)                             \
  (mSGpTR)->ittiMsg.nas_cs_domain_location_update_fail
#define NAS_CS_SERVICE_NOTIFICATION(mSGpTR)                                    \
  (mSGpTR)->ittiMsg.nas_cs_service_notification
#define NAS_DATA_LENGHT_MAX 256
#define NAS_EXTENDED_SERVICE_REQ(mSGpTR)                                       \
  (mSGpTR)->ittiMsg.nas_extended_service_req
#define NAS_TAU_COMPLETE(mSGpTR) (mSGpTR)->ittiMsg.nas_tau_complete
#define NAS_NOTIFY_SERVICE_REJECT(mSGpTR)                                      \
  (mSGpTR)->ittiMsg.nas_notify_service_reject
#define NAS_ERAB_REL_CMD(mSGpTR) (mSGpTR)->ittiMsg.itti_erab_rel_cmd

typedef struct itti_nas_cs_service_notification_s {
  mme_ue_s1ap_id_t ue_id; /* UE lower layer identifier        */
#define NAS_PAGING_ID_IMSI 0X00
#define NAS_PAGING_ID_TMSI 0X01
  uint8_t paging_id; /* Paging UE ID, to be sent in CS Service Notification */
  bstring
    cli; /* If CLI received in Sgsap-Paging_Req,shall sent in CS Service Notification */
} itti_nas_cs_service_notification_t;

typedef struct itti_nas_info_transfer_s {
  mme_ue_s1ap_id_t ue_id; /* UE lower layer identifier        */
  //nas_error_code_t err_code;     /* Transaction status               */
  bstring nas_msg; /* Uplink NAS message           */
} itti_nas_info_transfer_t;


typedef struct itti_erab_setup_req_s {
  mme_ue_s1ap_id_t ue_id; /* UE lower layer identifier   */
  ebi_t ebi;              /* EPS bearer id        */
  bstring nas_msg; /* NAS erab bearer context activation message           */
  bitrate_t mbr_dl;
  bitrate_t mbr_ul;
  bitrate_t gbr_dl;
  bitrate_t gbr_ul;
} itti_erab_setup_req_t;

typedef struct itti_erab_rel_cmd_s {
  mme_ue_s1ap_id_t ue_id;   /* UE lower layer identifier   */
  ebi_t ebi;                /* EPS bearer id        */
  bstring nas_msg;          /* NAS erab bearer context activation message */
  uint8_t n_bearers;        /*Number of bearers to be released*/
  ebi_t* bearers_to_be_rel; /* Bearers to be released*/
} itti_erab_rel_cmd_t;

typedef struct itti_nas_attach_req_s {
  /* TODO: Set the correct size */
  char apn[100];
  char imsi[16];
#define INITIAL_REQUEST (0x1)
  unsigned initial : 1;
  s1ap_initial_ue_message_t transparent;
} itti_nas_attach_req_t;

typedef struct itti_nas_auth_req_s {
  /* UE imsi */
  char imsi[16];

#define NAS_FAILURE_OK 0x0
#define NAS_FAILURE_IND 0x1
  unsigned failure : 1;
  int cause;
} itti_nas_auth_req_t;

typedef struct itti_nas_auth_rsp_s {
  char imsi[16];
} itti_nas_auth_rsp_t;

typedef struct itti_nas_auth_param_req_s {
  /* UE identifier */
  mme_ue_s1ap_id_t ue_id;

  /* Imsi of the UE (In case of initial request) */
  char imsi[16];
  uint8_t imsi_length;

  /* Indicates whether the procedure corresponds to a new connection or not */
  uint8_t initial_req : 1;

  uint8_t re_synchronization : 1;
  uint8_t auts[14];
  uint8_t num_vectors;
} itti_nas_auth_param_req_t;

typedef struct itti_nas_detach_req_s {
  /* UE identifier */
  mme_ue_s1ap_id_t ue_id;
  long cause;
} itti_nas_detach_req_t;

typedef struct itti_nas_implicit_detach_ue_ind_s {
  /* UE identifier */
  mme_ue_s1ap_id_t ue_id;
} itti_nas_implicit_detach_ue_ind_t;

typedef struct itti_nas_nw_initiated_detach_ue_req_s {
  /* UE identifier */
  mme_ue_s1ap_id_t ue_id;
#define HSS_INITIATED_EPS_DETACH 0x00
#define SGS_INITIATED_IMSI_DETACH 0x01
#define MME_INITIATED_EPS_DETACH 0x02
  uint8_t detach_type;
} itti_nas_nw_initiated_detach_ue_req_t;

typedef struct itti_nas_extended_service_req_s {
  /* UE identifier */
  mme_ue_s1ap_id_t ue_id;
  uint8_t servType; /* service type */
  /* csfb_response is valid only if service type Mobile Terminating CSFB */
  uint8_t csfb_response;
} itti_nas_extended_service_req_t;

typedef struct itti_nas_sgs_detach_req_s {
  /* UE identifier */
  mme_ue_s1ap_id_t ue_id;
  /* detach type */
  uint8_t detach_type;
} itti_nas_sgs_detach_req_t;

typedef struct itti_nas_cs_domain_location_update_req_s {
/* UE identifier */
#define ATTACH_REQ (1 << 0)
#define TAU_REQUEST (1 << 1)
  uint8_t msg_type;
  mme_ue_s1ap_id_t ue_id;
  uint8_t attach_type;
  additional_updt_t add_updt_type;
  uint8_t tau_updt_type; /*TAU Update type - Normal Update, Periodic,*/
} itti_nas_cs_domain_location_update_req_t;

typedef struct itti_nas_cs_domain_location_update_acc_s {
/* UE identifier */
#define MOBILE_IDENTITY (1 << 0)
#define ADD_UPDT_TYPE (1 << 1)
  uint8_t presencemask;
  mme_ue_s1ap_id_t ue_id;
  lai_t laicsfb;
  mobile_identity_t mobileid;
  additional_updt_result_t add_updt_res;
  bool is_sgs_assoc_exists;
} itti_nas_cs_domain_location_update_acc_t;

typedef struct itti_nas_cs_domain_location_update_fail_s {
/* UE identifier */
#define LAI (1 << 0)
  uint8_t presencemask;
  mme_ue_s1ap_id_t ue_id;
  int reject_cause;
  lai_t laicsfb;
} itti_nas_cs_domain_location_update_fail_t;

typedef struct itti_nas_tau_complete {
  /* UE identifier */
  mme_ue_s1ap_id_t ue_id;
} itti_nas_tau_complete_t;

/* ITTI message used to intimate service reject for ongoing service request procedure
 * from mme_app to nas
 */
typedef struct itti_nas_notify_service_reject_s {
  mme_ue_s1ap_id_t ue_id;
  uint8_t emm_cause;
#define INTIAL_CONTEXT_SETUP_PROCEDURE_FAILED 0x00
#define UE_CONTEXT_MODIFICATION_PROCEDURE_FAILED 0x01
#define MT_CALL_CANCELLED_BY_NW_IN_IDLE_STATE 0x02
#define MT_CALL_CANCELLED_BY_NW_IN_CONNECTED_STATE 0x03
  uint8_t failed_procedure;
} itti_nas_notify_service_reject_t;

#endif /* FILE_NAS_MESSAGES_TYPES_SEEN */
