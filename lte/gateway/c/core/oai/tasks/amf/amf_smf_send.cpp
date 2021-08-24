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

#include <sstream>
#ifdef __cplusplus
extern "C" {
#endif
#include "log.h"
#include "conversions.h"
#include "3gpp_38.401.h"
#ifdef __cplusplus
}
#endif
#include "amf_recv.h"
#include "M5gNasMessage.h"
#include "common_defs.h"
#include "amf_app_ue_context_and_proc.h"
#include "SmfServiceClient.h"
#include "M5GMobilityServiceClient.h"
#include "amf_app_timer_management.h"

using magma5g::AsyncM5GMobilityServiceClient;
using magma5g::AsyncSmfServiceClient;

namespace magma5g {
#define IMSI_LEN 15
#define AMF_CAUSE_SUCCESS 1

static int pdu_session_resource_release_t3592_handler(
    zloop_t* loop, int timer_id, void* arg);

/***************************************************************************
**                                                                        **
** Name:    esm_pt_is_reserved()                                          **
**                                                                        **
** Description: Check Validity of Procedure Transaction Identity          **
**                                                                        **
**                                                                        **
***************************************************************************/
int esm_pt_is_reserved(int pti) {
  return (
      (pti != PROCEDURE_TRANSACTION_IDENTITY_UNASSIGNED_t) &&
      (pti > PROCEDURE_TRANSACTION_IDENTITY_LAST_t));
}

/***************************************************************************
**                                                                        **
** Name:    amf_smf_handle_pdu_establishment_request()                    **
**                                                                        **
** Description: Handler for PDU Establishment Requests                    **
**                                                                        **
**                                                                        **
***************************************************************************/
int amf_smf_handle_pdu_establishment_request(
    SmfMsg* msg, amf_smf_t* amf_smf_msg) {
  int smf_cause = SMF_CAUSE_SUCCESS;
  OAILOG_DEBUG(
      LOG_AMF_APP,
      "AMF SMF Handler- Received PDN Connectivity Request message ");

  // Procedure transaction identity checking
  if ((msg->header.procedure_transaction_id ==
       PROCEDURE_TRANSACTION_IDENTITY_UNASSIGNED_t) ||
      esm_pt_is_reserved(msg->header.procedure_transaction_id)) {
    amf_smf_msg->u.establish.cause_value = SMF_CAUSE_INVALID_PTI_VALUE;
    OAILOG_DEBUG(
        LOG_AMF_APP, "smf_cause : %u", amf_smf_msg->u.establish.cause_value);
    return (amf_smf_msg->u.establish.cause_value);
  } else {
    amf_smf_msg->u.establish.pti = msg->header.procedure_transaction_id;
  }

  // Get the value of the PDN type indicator
  if (msg->msg.pdu_session_estab_request.pdu_session_type.type_val ==
      PDN_TYPE_IPV4) {
    amf_smf_msg->u.establish.pdu_session_type = NET_PDN_TYPE_IPV4;
  } else if (
      msg->msg.pdu_session_estab_request.pdu_session_type.type_val ==
      PDN_TYPE_IPV6) {
    amf_smf_msg->u.establish.pdu_session_type = NET_PDN_TYPE_IPV6;
  } else if (
      msg->msg.pdu_session_estab_request.pdu_session_type.type_val ==
      PDN_TYPE_IPV4V6) {
    amf_smf_msg->u.establish.pdu_session_type = NET_PDN_TYPE_IPV4V6;
  } else {
    // Unknown PDN type
    amf_smf_msg->u.establish.cause_value = SMF_CAUSE_UNKNOWN_PDN_TYPE;
    return (amf_smf_msg->u.establish.cause_value);
  }
  amf_smf_msg->u.establish.pdu_session_id = msg->header.pdu_session_id;
  amf_smf_msg->u.establish.cause_value    = smf_cause;
  return (smf_cause);
}

/***************************************************************************
**                                                                        **
** Name:    amf_smf_handle_pdu_release_request                            **
**                                                                        **
** Description: handler for PDU session release                           **
**                                                                        **
**                                                                        **
***************************************************************************/
int amf_smf_handle_pdu_release_request(SmfMsg* msg, amf_smf_t* amf_smf_msg) {
  int smf_cause                         = SMF_CAUSE_SUCCESS;
  amf_smf_msg->u.release.pti            = msg->header.procedure_transaction_id;
  amf_smf_msg->u.release.pdu_session_id = msg->header.pdu_session_id;
  amf_smf_msg->u.release.cause_value    = smf_cause;
  return (smf_cause);  // TODO add error checking and return
                       // appropriate cause value
}

/***************************************************************************
**                                                                        **
** Name:    amf_send_pdusession_reject()                                  **
**                                                                        **
** Description: send PDU session reject to UE                             **
**                                                                        **
**                                                                        **
***************************************************************************/
int amf_send_pdusession_reject(
    SmfMsg* reject_req, uint8_t session_id, uint8_t pti, uint8_t cause) {
  uint8_t buffer[5];
  int rc;
  reject_req->header.extended_protocol_discriminator =
      M5G_SESSION_MANAGEMENT_MESSAGES;
  reject_req->header.pdu_session_id           = session_id;
  reject_req->header.procedure_transaction_id = pti;
  reject_req->header.message_type = PDU_SESSION_ESTABLISHMENT_REJECT;
  reject_req->msg.pdu_session_estab_reject.m5gsm_cause.cause_value = cause;
  rc = reject_req->SmfMsgEncodeMsg(reject_req, buffer, 5);
  if (rc > 0) {
    // TODO: Send the message to AS for nas encode
    // and forward to NGAP. Nagetive scenario.
  }
  return rc;
}

/***************************************************************************
**                                                                        **
** Name:    set_amf_smf_context()                                         **
**                                                                        **
** Description: set the smf_context in amf_context                        **
**                                                                        **
**                                                                        **
***************************************************************************/
void set_amf_smf_context(
    PDUSessionEstablishmentRequestMsg* message, smf_context_t* smf_ctx) {
  smf_ctx->smf_proc_data.pdu_session_identity = message->pdu_session_identity;
  smf_ctx->smf_proc_data.pti                  = message->pti;
  smf_ctx->smf_proc_data.message_type         = message->message_type;
  smf_ctx->smf_proc_data.integrity_prot_max_data_rate =
      message->integrity_prot_max_data_rate;
  smf_ctx->smf_proc_data.pdu_session_type = message->pdu_session_type;
  smf_ctx->smf_proc_data.ssc_mode         = message->ssc_mode;
  smf_ctx->pdu_session_version            = 0;  // Initializing pdu version to 0
  memset(
      smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr, '\0',
      sizeof(smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr));
  memset(
      smf_ctx->gtp_tunnel_id.gnb_gtp_teid, '\0',
      sizeof(smf_ctx->gtp_tunnel_id.gnb_gtp_teid));
}

/***************************************************************************
**                                                                        **
** Name:    clear_amf_smf_context()                                       **
**                                                                        **
** Description: clear smf_context on session release                      **
**                                                                        **
**                                                                        **
***************************************************************************/
void clear_amf_smf_context(smf_context_t* smf_ctx) {
  OAILOG_DEBUG(
      LOG_AMF_APP, "clearing saved context associated with the pdu session\n");
  memset(
      &(smf_ctx->smf_proc_data.pdu_session_identity), 0,
      sizeof(smf_ctx->smf_proc_data.pdu_session_identity));
  memset(&(smf_ctx->smf_proc_data.pti), 0, sizeof(smf_ctx->smf_proc_data.pti));
  memset(
      &(smf_ctx->smf_proc_data.message_type), 0,
      sizeof(smf_ctx->smf_proc_data.message_type));
  memset(
      &(smf_ctx->smf_proc_data.integrity_prot_max_data_rate), 0,
      sizeof(smf_ctx->smf_proc_data.integrity_prot_max_data_rate));
  memset(
      &(smf_ctx->smf_proc_data.pdu_session_type), 0,
      sizeof(smf_ctx->smf_proc_data.pdu_session_type));
  memset(
      &(smf_ctx->smf_proc_data.ssc_mode), 0,
      sizeof(smf_ctx->smf_proc_data.ssc_mode));
}

int pdu_session_release_request_process(
    ue_m5gmm_context_s* ue_context, smf_context_t* smf_ctx,
    amf_ue_ngap_id_t amf_ue_ngap_id) {
  int rc                = 1;
  amf_smf_t amf_smf_msg = {};
  // amf_cause = amf_smf_handle_pdu_release_request(
  //              msg, &amf_smf_msg);

  int smf_cause             = SMF_CAUSE_SUCCESS;
  amf_smf_msg.u.release.pti = smf_ctx->smf_proc_data.pti.pti;
  amf_smf_msg.u.release.pdu_session_id =
      smf_ctx->smf_proc_data.pdu_session_identity.pdu_session_id;
  amf_smf_msg.u.release.cause_value = smf_cause;

  OAILOG_DEBUG(
      LOG_AMF_APP, "sending PDU session resource release request to gNB \n");

  rc =
      pdu_session_resource_release_request(ue_context, amf_ue_ngap_id, smf_ctx);

  if (rc != RETURNok) {
    OAILOG_DEBUG(
        LOG_AMF_APP,
        "PDU session resource release request to gNB failed"
        "\n");
  } else {
    ue_pdu_id_t id = {
        amf_ue_ngap_id,
        smf_ctx->smf_proc_data.pdu_session_identity.pdu_session_id};

    smf_ctx->T3592.id = amf_pdu_start_timer(
        PDUE_SESSION_RELEASE_TIMER_MSECS, TIMER_REPEAT_ONCE,
        pdu_session_resource_release_t3592_handler, id);
  }

  return rc;
}

int pdu_session_resource_release_complete(
    ue_m5gmm_context_s* ue_context, amf_smf_t amf_smf_msg,
    smf_context_t* smf_ctx) {
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  int rc = 1;

  IMSI64_TO_STRING(ue_context->amf_context.imsi64, imsi, 15);

  if (smf_ctx->n_active_pdus) {
    /* Execute PDU Session Release and notify to SMF */
    rc = pdu_state_handle_message(
        ue_context->mm_state, STATE_PDU_SESSION_RELEASE_COMPLETE,
        smf_ctx->pdu_session_state, ue_context, amf_smf_msg, imsi, NULL, 0);
  }

  OAILOG_INFO(
      LOG_AMF_APP, "notifying SMF about PDU session release n_active_pdus=%d\n",
      smf_ctx->n_active_pdus);

  if (smf_ctx->pdu_address.pdn_type == IPv4) {
    // Clean up the Mobility IP Address
    AsyncM5GMobilityServiceClient::getInstance().release_ipv4_address(
        imsi, reinterpret_cast<const char*>(smf_ctx->apn),
        &(smf_ctx->pdu_address.ipv4_address));
  }

  OAILOG_DEBUG(
      LOG_AMF_APP, "clear saved context associated with the PDU session\n");
  clear_amf_smf_context(smf_ctx);
}

static int pdu_session_resource_release_t3592_handler(
    zloop_t* loop, int timer_id, void* arg) {
  OAILOG_INFO(
      LOG_AMF_APP, "T3592: pdu_session_resource_release_t3592_handler\n");

  amf_ue_ngap_id_t amf_ue_ngap_id = 0;
  uint8_t pdu_session_id          = 0;
  ue_pdu_id_t uepdu_id;
  smf_context_t* smf_ctx = NULL;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  int rc = 0;

  if (!amf_pdu_get_timer_arg(timer_id, &uepdu_id)) {
    OAILOG_WARNING(
        LOG_AMF_APP, "T3550: Invalid Timer Id expiration, Timer Id: %u\n",
        timer_id);
    OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
  }

  amf_ue_ngap_id = uepdu_id.ue_id;
  pdu_session_id = uepdu_id.pdu_id;

  ue_m5gmm_context_s* ue_context =
      amf_ue_context_exists_amf_ue_ngap_id(amf_ue_ngap_id);

  if (ue_context) {
    IMSI64_TO_STRING(ue_context->amf_context.imsi64, imsi, 15);
    smf_ctx = amf_smf_context_exists_pdu_session_id(ue_context, pdu_session_id);

    if (smf_ctx == NULL) {
      OAILOG_ERROR(
          LOG_AMF_APP, "T3592:pdu session  not found for session_id = %u\n",
          pdu_session_id);
      OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
    }
  } else {
    OAILOG_ERROR(
        LOG_AMF_APP, "T3592: ue context not found for the ue_id=%u\n",
        amf_ue_ngap_id);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
  }

  OAILOG_WARNING(
      LOG_AMF_APP, "T3592: timer id: %d expired for pdu_session_id: %d\n",
      smf_ctx->T3592.id, pdu_session_id);

  smf_ctx->retransmission_count += 1;

  OAILOG_ERROR(
      LOG_AMF_APP, "T3592: Incrementing retransmission_count to %d\n",
      smf_ctx->retransmission_count);

  if (smf_ctx->retransmission_count < REGISTRATION_COUNTER_MAX) {
    /* Send entity Registration accept message to the UE */

    pdu_session_release_request_process(ue_context, smf_ctx, amf_ue_ngap_id);
  } else {
    /* Abort the registration procedure */
    OAILOG_ERROR(
        LOG_AMF_APP,
        "T3592: Maximum retires:%d, for PDU_SESSION_RELEASE_COMPELETE done "
        "hence Abort "
        "the pdu sesssion release "
        "procedure\n",
        smf_ctx->retransmission_count);
    // To abort the registration procedure
    // amf_proc_registration_abort(amf_ctx, ue_amf_context);
    // pdu_session_resource_release_abort()
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
}

/***************************************************************************
**                                                                        **
** Name:    amf_smf_send()                                                **
**                                                                        **
** Description: handler to send session request to SMF                    **
**                                                                        **
**                                                                        **
***************************************************************************/
int amf_smf_send(
    amf_ue_ngap_id_t ue_id, ULNASTransportMsg* msg, int amf_cause) {
  int rc = 1;
  SmfMsg reject_req;
  amf_smf_t amf_smf_msg  = {};
  smf_context_t* smf_ctx = NULL;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  if (amf_cause != AMF_CAUSE_SUCCESS) {
    rc = amf_send_pdusession_reject(
        &reject_req, msg->payload_container.smf_msg.header.pdu_session_id,
        msg->payload_container.smf_msg.header.procedure_transaction_id,
        amf_cause);
    return rc;
  }

  ue_m5gmm_context_s* ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  if (ue_context) {
    IMSI64_TO_STRING(ue_context->amf_context.imsi64, imsi, 15);
    if (msg->payload_container.smf_msg.header.message_type ==
        PDU_SESSION_ESTABLISHMENT_REQUEST) {
      smf_ctx = amf_insert_smf_context(
          ue_context, msg->payload_container.smf_msg.header.pdu_session_id);
    } else {
      smf_ctx = amf_smf_context_exists_pdu_session_id(
          ue_context, msg->payload_container.smf_msg.header.pdu_session_id);
    }
    if (smf_ctx == NULL) {
      OAILOG_ERROR(
          LOG_AMF_APP, "pdu session  not found for session_id = %u\n",
          msg->payload_container.smf_msg.header.pdu_session_id);
      OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
    }
  } else {
    OAILOG_ERROR(LOG_AMF_APP, "ue context not found for the ue_id=%u\n", ue_id);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
  }

  // Process the decoded NAS message
  switch (msg->payload_container.smf_msg.header.message_type) {
    case PDU_SESSION_ESTABLISHMENT_REQUEST: {
      amf_cause = amf_smf_handle_pdu_establishment_request(
          &(msg->payload_container.smf_msg), &amf_smf_msg);

      if (amf_cause != SMF_CAUSE_SUCCESS) {
        rc = amf_send_pdusession_reject(
            &reject_req, msg->payload_container.smf_msg.header.pdu_session_id,
            msg->payload_container.smf_msg.header.procedure_transaction_id,
            amf_cause);
        return rc;
      }
      set_amf_smf_context(
          &(msg->payload_container.smf_msg.msg.pdu_session_estab_request),
          smf_ctx);
      memset(
          amf_smf_msg.u.establish.gnb_gtp_teid_ip_addr, '\0',
          sizeof(amf_smf_msg.u.establish.gnb_gtp_teid_ip_addr));
      memset(
          amf_smf_msg.u.establish.gnb_gtp_teid, '\0',
          sizeof(amf_smf_msg.u.establish.gnb_gtp_teid));
      memcpy(
          amf_smf_msg.u.establish.gnb_gtp_teid_ip_addr,
          smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr, GNB_IPV4_ADDR_LEN);
      memcpy(
          amf_smf_msg.u.establish.gnb_gtp_teid,
          smf_ctx->gtp_tunnel_id.gnb_gtp_teid, GNB_TEID_LEN);

      // Initialize default APN
      memcpy(smf_ctx->apn, "internet", strlen("internet") + 1);

      // send request to SMF over grpc
      /*
       * Execute the Grpc Send call of PDU establishment Request from AMF to SMF
       */
      rc = pdu_state_handle_message(
          // ue_context->mm_state, STATE_PDU_SESSION_ESTABLISHMENT_REQUEST,
          REGISTERED_CONNECTED, STATE_PDU_SESSION_ESTABLISHMENT_REQUEST,
          // smf_ctx->pdu_session_state, ue_context, amf_smf_msg, imsi, NULL,
          // 0);
          SESSION_NULL, ue_context, amf_smf_msg, imsi, NULL, 0);
    } break;
    case PDU_SESSION_RELEASE_REQUEST: {
      smf_ctx->retransmission_count = 0;
      if (RETURNok ==
          pdu_session_release_request_process(ue_context, smf_ctx, ue_id)) {
        OAILOG_INFO(
            LOG_AMF_APP,
            "T3592: PDU_SESSION_RELEASE_REQUEST timer T3592 with id  %d "
            "Started\n",
            smf_ctx->T3592.id);
      }
    } break;
    case PDU_SESSION_RELEASE_COMPLETE: {
      if (smf_ctx->T3592.id != NAS5G_TIMER_INACTIVE_ID) {
        amf_pdu_stop_timer(smf_ctx->T3592.id);
        OAILOG_INFO(
            LOG_AMF_APP,
            "T3592: after stop PDU_SESSION_RELEASE_REQUEST timer T3592 with id "
            "= %d\n",
            smf_ctx->T3592.id);
        smf_ctx->T3592.id = NAS5G_TIMER_INACTIVE_ID;
      }
      amf_cause = amf_smf_handle_pdu_release_request(
          &(msg->payload_container.smf_msg), &amf_smf_msg);

      pdu_session_resource_release_complete(ue_context, amf_smf_msg, smf_ctx);
    } break;
    default:
      break;
  }
  return rc;
}

/***************************************************************************
**                                                                        **
** Name:    amf_smf_notification_send()                                   **
**                                                                        **
** Description:                                                           **
** This function for UE idle event notification to SMF or single PDU      **
** session state change to Inactive state and notify to SMF.              **
** 4 types of events are used in proto.                                   **
** PDU_SESSION_INACTIVE_NOTIFY => use for single PDU session notify       **
** UE_IDLE_MODE_NOTIFY     => use for idle mode support                   **
** UE_PAGING_NOTIFY                                                       **
** UE_PERIODIC_REG_ACTIVE_MODE_NOTIFY                                     **
**                                                                        **
***************************************************************************/
int amf_smf_notification_send(
    amf_ue_ngap_id_t ue_id, ue_m5gmm_context_s* ue_context,
    notify_ue_event notify_event_type) {
  /* Get gRPC structure of notification to be filled common and
   * rat type elements.
   * Only need  to be filled IMSI and ue_state_idle of UE
   */
  magma::lte::SetSmNotificationContext notify_req;
  auto* req_common       = notify_req.mutable_common_context();
  auto* req_rat_specific = notify_req.mutable_rat_specific_notification();
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  IMSI64_TO_STRING(ue_context->amf_context.imsi64, imsi, 15);

  req_common->mutable_sid()->mutable_id()->assign(imsi);
  if (notify_event_type == UE_IDLE_MODE_NOTIFY) {
    req_rat_specific->set_notify_ue_event(
        magma::lte::NotifyUeEvents::UE_IDLE_MODE_NOTIFY);
  } else if (notify_event_type == UE_SERVICE_REQUEST_ON_PAGING) {
    req_rat_specific->set_notify_ue_event(
        magma::lte::NotifyUeEvents::UE_SERVICE_REQUEST_ON_PAGING);
  }

  auto it = ue_context->amf_context.smf_ctxt_vector.begin();
  if (it != ue_context->amf_context.smf_ctxt_vector.end()) {
    smf_context_t smf_context = *it;

    if (smf_context.pdu_address.pdn_type == IPv4) {
      char ip_str[INET_ADDRSTRLEN];

      inet_ntop(
          AF_INET, &(smf_context.pdu_address.ipv4_address.s_addr), ip_str,
          INET_ADDRSTRLEN);
      req_common->set_ue_ipv4((char*) ip_str);
    }
  }
  // Set the PDU Address

  OAILOG_DEBUG(
      LOG_AMF_APP,
      " Notification gRPC filled with IMSI %s and "
      "ue_state_idle is set to true \n",
      imsi);

  AsyncSmfServiceClient::getInstance().set_smf_notification(notify_req);

  return RETURNok;
}

/***************************************************************************
**                                                                        **
** Name:    amf_smf_context_exists_pdu_session_id()                       **
**                                                                        **
** Description: Update IP Addrss information in SMF Context               **
**                                                                        **
**                                                                        **
***************************************************************************/
int amf_update_smf_context_pdu_ip(
    char* imsi, uint8_t* apn, uint32_t pdu_session_id, paa_t* address_info) {
  ue_m5gmm_context_s* ue_context;
  smf_context_t* smf_ctx;
  imsi64_t imsi64;
  int rc = RETURNerror;

  IMSI_STRING_TO_IMSI64(imsi, &imsi64);

  ue_context = lookup_ue_ctxt_by_imsi(imsi64);
  if (ue_context == NULL) {
    return rc;
  }

  smf_ctx = amf_smf_context_exists_pdu_session_id(ue_context, pdu_session_id);
  if (NULL == smf_ctx) {
    return rc;
  }

  memcpy(&(smf_ctx->pdu_address), address_info, sizeof(paa_t));

  return RETURNok;
}

/***************************************************************************
**                                                                        **
** Name:    amf_smf_context_exists_pdu_session_id()                       **
**                                                                        **
** Description: Update IP Addrss information in SMF Context               **
**                                                                        **
**                                                                        **
***************************************************************************/
int amf_smf_handle_ip_address_response(
    itti_amf_ip_allocation_response_t* response_p) {
  int rc = RETURNerror;

  rc = amf_update_smf_context_pdu_ip(
      response_p->imsi, response_p->apn, response_p->pdu_session_id,
      &(response_p->paa));

  if (rc < 0) {
    OAILOG_ERROR(
        LOG_AMF_APP,
        "SMF Context for PDU not found or Address "
        "type not supported\n");
    return rc;
  }

  if (response_p->paa.pdn_type == IPv4) {
    char ip_str[INET_ADDRSTRLEN];

    inet_ntop(
        AF_INET, &(response_p->paa.ipv4_address.s_addr), ip_str,
        INET_ADDRSTRLEN);

    rc = amf_smf_create_ipv4_session_grpc_req(
        response_p->imsi, response_p->apn, response_p->pdu_session_id,
        response_p->pdu_session_type, response_p->gnb_gtp_teid, response_p->pti,
        response_p->gnb_gtp_teid_ip_addr, ip_str);

    if (rc < 0) {
      OAILOG_ERROR(LOG_AMF_APP, "Create IPV4 Session \n");
    }
  }

  return rc;
}

}  // namespace magma5g
