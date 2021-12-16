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
#include <string.h>
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/lib/3gpp/3gpp_38.401.h"
#include "lte/gateway/c/core/oai/include/s6a_messages_types.h"
#include "lte/gateway/c/core/oai/include/amf_config.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_session_manager_pco.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_recv.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_sap.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"
#include "lte/gateway/c/core/oai/lib/n11/SmfServiceClient.h"
#include "lte/gateway/c/core/oai/lib/n11/M5GMobilityServiceClient.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_timer_management.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_common.h"
#include "lte/gateway/c/core/oai/tasks/nas/api/mme/mme_api.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_defs.h"
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_smf_packet_handler.h"
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_client_servicer.h"

using magma5g::AsyncM5GMobilityServiceClient;
using magma5g::AsyncSmfServiceClient;

extern amf_config_t amf_config;
namespace magma5g {
#define IMSI_LEN 15

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
    PDUSessionEstablishmentRequestMsg* message,
    std::shared_ptr<smf_context_t> smf_ctx) {
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
  smf_ctx->gtp_tunnel_id.gnb_gtp_teid = 0x0;
}

/***************************************************************************
**                                                                        **
** Name:    clear_amf_smf_context()                                       **
**                                                                        **
** Description: clear smf_context on session release                      **
**                                                                        **
**                                                                        **
***************************************************************************/
void clear_amf_smf_context(amf_context_t& amf_context, uint8_t pdu_session_id) {
  OAILOG_DEBUG(
      LOG_AMF_APP, "clearing saved context associated with the pdu session\n");
  auto it = amf_context.smf_ctxt_map.find(pdu_session_id);
  if (it != amf_context.smf_ctxt_map.end()) {
    amf_context.smf_ctxt_map.erase(it);
  } else {
    OAILOG_WARNING(LOG_AMF_APP, "PDU Session is not found");
  }
}

int pdu_session_release_request_process(
    ue_m5gmm_context_s* ue_context, std::shared_ptr<smf_context_t> smf_ctx,
    amf_ue_ngap_id_t amf_ue_ngap_id, bool retransmit) {
  int rc = RETURNerror;
  OAILOG_DEBUG(
      LOG_AMF_APP, "sending PDU session resource release request to gNB \n");

  rc = pdu_session_resource_release_request(
      ue_context, amf_ue_ngap_id, smf_ctx, retransmit);

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
    std::shared_ptr<smf_context_t> smf_ctx) {
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
    AMFClientServicer::getInstance().release_ipv4_address(
        imsi, smf_ctx->dnn.c_str(), &(smf_ctx->pdu_address.ipv4_address));
  }

  OAILOG_DEBUG(
      LOG_AMF_APP, "clear saved context associated with the PDU session\n");
  clear_amf_smf_context(ue_context->amf_context, amf_smf_msg.pdu_session_id);
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
}

static int pdu_session_resource_release_t3592_handler(
    zloop_t* loop, int timer_id, void* arg) {
  OAILOG_INFO(
      LOG_AMF_APP, "T3592: pdu_session_resource_release_t3592_handler\n");

  amf_ue_ngap_id_t amf_ue_ngap_id = 0;
  uint8_t pdu_session_id          = 0;
  ue_pdu_id_t uepdu_id;
  std::shared_ptr<smf_context_t> smf_ctx;
  char imsi[IMSI_BCD_DIGITS_MAX + 1];
  int rc = 0;

  if (!amf_pop_pdu_timer_arg(timer_id, &uepdu_id)) {
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
    smf_ctx = amf_get_smf_context_by_pdu_session_id(ue_context, pdu_session_id);

    if (smf_ctx == NULL) {
      OAILOG_ERROR(
          LOG_AMF_APP, "T3592:pdu session  not found for session_id = %u\n",
          pdu_session_id);
      OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
    }
  } else {
    OAILOG_ERROR(
        LOG_AMF_APP,
        "T3592: ue context not found for UE ID = " AMF_UE_NGAP_ID_FMT,
        amf_ue_ngap_id);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
  }

  OAILOG_WARNING(
      LOG_AMF_APP, "T3592: timer id: %ld expired for pdu_session_id: %d\n",
      smf_ctx->T3592.id, pdu_session_id);

  smf_ctx->retransmission_count += 1;

  OAILOG_ERROR(
      LOG_AMF_APP, "T3592: Incrementing retransmission_count to %d\n",
      smf_ctx->retransmission_count);

  if (smf_ctx->retransmission_count < REGISTRATION_COUNTER_MAX) {
    /* Send entity Registration accept message to the UE */

    ue_pdu_id_t id = {amf_ue_ngap_id, pdu_session_id};

    pdu_session_release_request_process(
        ue_context, smf_ctx, amf_ue_ngap_id, true);

    smf_ctx->T3592.id = amf_pdu_start_timer(
        PDUE_SESSION_RELEASE_TIMER_MSECS, TIMER_REPEAT_ONCE,
        pdu_session_resource_release_t3592_handler, id);

  } else {
    /* Abort the registration procedure */
    OAILOG_ERROR(
        LOG_AMF_APP,
        "T3592: Maximum retires:%d, for PDU_SESSION_RELEASE_COMPELETE done "
        "hence Abort "
        "the pdu sesssion release "
        "procedure\n",
        smf_ctx->retransmission_count);
  }
  OAILOG_FUNC_RETURN(LOG_NAS_AMF, RETURNok);
}

/***************************************************************************
**                                                                        **
** Name:    amf_smf_process_pdu_session_packet()                          **
**                                                                        **
** Description: handler to send session request to SMF                    **
**                                                                        **
**                                                                        **
***************************************************************************/
int amf_smf_process_pdu_session_packet(
    amf_ue_ngap_id_t ue_id, ULNASTransportMsg* msg, int amf_cause) {
  int rc                = RETURNok;
  amf_smf_t amf_smf_msg = {};
  std::shared_ptr<smf_context_t> smf_ctx;
  char imsi[IMSI_BCD_DIGITS_MAX + 1]        = {0};
  protocol_configuration_options_t* msg_pco = nullptr;

  if (!msg) {
    return RETURNerror;
  }

  if (amf_cause != AMF_CAUSE_SUCCESS) {
    rc = amf_pdu_session_establishment_reject(
        ue_id, msg->payload_container.smf_msg.header.pdu_session_id,
        msg->payload_container.smf_msg.header.procedure_transaction_id,
        amf_cause);
    return rc;
  }

  ue_m5gmm_context_s* ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (!ue_context) {
    OAILOG_ERROR(
        LOG_AMF_APP, "ue context not found for the ue_id :" AMF_UE_NGAP_ID_FMT,
        ue_id);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }

  if (msg->payload_container.smf_msg.header.message_type ==
      PDU_SESSION_ESTABLISHMENT_REQUEST) {
    M5GSmCause cause = amf_smf_get_smcause(ue_id, msg);

    if (cause != M5GSmCause::INVALID_CAUSE) {
      OAILOG_DEBUG(
          LOG_AMF_APP,
          "PDU Session establishment request rejecting with cause %u",
          static_cast<uint8_t>(cause));
      rc = amf_pdu_session_establishment_reject(
          ue_id, msg->payload_container.smf_msg.header.pdu_session_id,
          msg->payload_container.smf_msg.header.procedure_transaction_id,
          static_cast<uint8_t>(cause));
      return rc;
    }

    M5GMmCause mm_cause = amf_smf_validate_context(ue_id, msg);
    if (mm_cause == M5GMmCause::MAX_PDU_SESSIONS_REACHED) {
      OAILOG_ERROR(
          LOG_AMF_APP,
          "Max pdu session limit reached, Rejecting new session for the "
          "ue_id :" AMF_UE_NGAP_ID_FMT,
          ue_id);
      rc = handle_sm_message_routing_failure(ue_id, msg, mm_cause);
      return rc;
    }
    smf_ctx = amf_get_smf_context_by_pdu_session_id(
        ue_context, msg->payload_container.smf_msg.header.pdu_session_id);

    if (smf_ctx && smf_ctx->duplicate_pdu_session_est_req_count > 0) {
      OAILOG_DEBUG(
          LOG_AMF_APP, "Duplicate PDU Session Establishment Request, Dropped");
      return rc;
    }
  }
  IMSI64_TO_STRING(ue_context->amf_context.imsi64, imsi, 15);
  if (msg->payload_container.smf_msg.header.message_type ==
      PDU_SESSION_ESTABLISHMENT_REQUEST) {
    smf_ctx = amf_insert_smf_context(
        ue_context, msg->payload_container.smf_msg.header.pdu_session_id);
  } else {
    smf_ctx = amf_get_smf_context_by_pdu_session_id(
        ue_context, msg->payload_container.smf_msg.header.pdu_session_id);
  }

  if (smf_ctx == NULL) {
    OAILOG_ERROR(
        LOG_AMF_APP, "pdu session  not found for session_id = %u\n",
        msg->payload_container.smf_msg.header.pdu_session_id);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, RETURNerror);
  }

  amf_smf_msg.pdu_session_id =
      msg->payload_container.smf_msg.header.pdu_session_id;
  // Process the decoded NAS message
  switch (msg->payload_container.smf_msg.header.message_type) {
    case PDU_SESSION_ESTABLISHMENT_REQUEST: {
      amf_cause = amf_smf_handle_pdu_establishment_request(
          &(msg->payload_container.smf_msg), &amf_smf_msg);

      OAILOG_INFO(LOG_AMF_APP, "Copy the contents from message to context \n");
      msg_pco = &(msg->payload_container.smf_msg.msg.pdu_session_estab_request
                      .protocolconfigurationoptions.pco);

      /* Copy the pco contents from Message to smf_context */
      sm_copy_protocol_configuration_options(&(smf_ctx->pco), msg_pco);

      /* Free the memory from Message Structure */
      sm_free_protocol_configuration_options(&msg_pco);

      if (amf_cause != SMF_CAUSE_SUCCESS) {
        rc = amf_pdu_session_establishment_reject(
            ue_id, msg->payload_container.smf_msg.header.pdu_session_id,
            msg->payload_container.smf_msg.header.procedure_transaction_id,
            static_cast<uint8_t>(amf_cause));

        return rc;
      }

      smf_ctx->sst = msg->nssai.sst;
      if (msg->nssai.sd[0]) {
        memcpy(smf_ctx->sd, msg->nssai.sd, SD_LENGTH);
      }
      set_amf_smf_context(
          &(msg->payload_container.smf_msg.msg.pdu_session_estab_request),
          smf_ctx);
      memset(
          amf_smf_msg.u.establish.gnb_gtp_teid_ip_addr, '\0',
          sizeof(amf_smf_msg.u.establish.gnb_gtp_teid_ip_addr));

      amf_smf_msg.u.establish.gnb_gtp_teid = 0x0;
      memcpy(
          amf_smf_msg.u.establish.gnb_gtp_teid_ip_addr,
          smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr, GNB_IPV4_ADDR_LEN);

      amf_smf_msg.u.establish.gnb_gtp_teid =
          smf_ctx->gtp_tunnel_id.gnb_gtp_teid;

      // Initialize DNN
      char* default_dnn = bstr2cstr(amf_config.default_dnn, '?');

      int index_dnn    = 0;
      bool ue_sent_dnn = true;
      std::string dnn_string;

      if (msg->dnn.len <= 1) {
        ue_sent_dnn = false;
        dnn_string  = default_dnn;
      } else {
        dnn_string.assign(
            reinterpret_cast<char*>(msg->dnn.dnn), msg->dnn.len - 1);
      }

      int validate = amf_validate_dnn(
          &ue_context->amf_context, dnn_string, &index_dnn, ue_sent_dnn);
      free(default_dnn);

      if (validate == RETURNok) {
        smf_dnn_ambr_select(smf_ctx, ue_context, index_dnn);
      } else {
        OAILOG_INFO(
            LOG_AMF_APP,
            "DNN is not Supported or not Subscribed, reject with a cause: 91 "
            "\n");
        M5GMmCause cause_dnn_reject =
            M5GMmCause::DNN_NOT_SUPPORTED_OR_NOT_SUBSCRIBED;
        rc = handle_sm_message_routing_failure(ue_id, msg, cause_dnn_reject);
        ue_context->amf_context.smf_ctxt_map.erase(
            msg->payload_container.smf_msg.header.pdu_session_id);
        return rc;
      }

      smf_ctx->smf_proc_data.pti.pti =
          msg->payload_container.smf_msg.msg.pdu_session_estab_request.pti.pti;
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
      smf_ctx->smf_proc_data.pti.pti = msg->payload_container.smf_msg.msg
                                           .pdu_session_release_request.pti.pti;
      smf_ctx->retransmission_count = 0;
      if (RETURNok == pdu_session_release_request_process(
                          ue_context, smf_ctx, ue_id, false)) {
        OAILOG_INFO(
            LOG_AMF_APP,
            "T3592: PDU_SESSION_RELEASE_REQUEST timer T3592 with id  %ld "
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
            "= %ld\n",
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
** Name:    smf_dnn_ambr_select()                                         **
**                                                                        **
** Description: Copy dnn and ambr info in smf context                     **
**                                                                        **
**                                                                        **
***************************************************************************/
void smf_dnn_ambr_select(
    const std::shared_ptr<smf_context_t>& smf_ctx,
    ue_m5gmm_context_s* ue_context, int index_dnn) {
  smf_ctx->dnn.assign(
      reinterpret_cast<char*>(ue_context->amf_context.apn_config_profile
                                  .apn_configuration[index_dnn]
                                  .service_selection),
      strlen(ue_context->amf_context.apn_config_profile
                 .apn_configuration[index_dnn]
                 .service_selection));
  OAILOG_INFO(LOG_AMF_APP, "dnn selected %s\n", smf_ctx->dnn.c_str());

  memcpy(
      &smf_ctx->smf_ctx_ambr,
      &ue_context->amf_context.apn_config_profile.apn_configuration[index_dnn]
           .ambr,
      sizeof(ambr_t));
}
/***************************************************************************
**                                                                        **
** Name:    amf_smf_get_smcause()                                         **
**                                                                        **
** Description: function to handle PDU Session Establishment Failures     **
**                                                                        **
**                                                                        **
***************************************************************************/
M5GSmCause amf_smf_get_smcause(amf_ue_ngap_id_t ue_id, ULNASTransportMsg* msg) {
  std::shared_ptr<smf_context_t> smf_ctx;
  M5GSmCause cause = M5GSmCause::INVALID_CAUSE;

  ue_m5gmm_context_s* ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  if (!ue_context) {
    return cause;
  }
  /*
  Cause #27 – Missing or unknown DNN
  the external DNN because the DNN was not included
  although required or if the DNN could not be resolved.
  */
  if (msg->dnn.len <= 1 &&
      (ue_context->amf_context.apn_config_profile.nb_apns == 0)) {
    cause = M5GSmCause::MISSING_OR_UNKNOWN_DNN;
    return cause;
  }

  /*
  Cause #28 – Unknown PDU session type
  the requested PDU session type could not be recognized or is not allowed.
  */
  M5GPduSessionType session_type = static_cast<M5GPduSessionType>(
      msg->payload_container.smf_msg.msg.pdu_session_estab_request
          .pdu_session_type.type_val);
  if (session_type == M5GPduSessionType::UNSTRUCTURED ||
      session_type == M5GPduSessionType::ETHERNET) {
    cause = M5GSmCause::UNKNOWN_PDU_SESSION_TYPE;
    return cause;
  }

  /*
  Cause #43 – Invalid PDU session identity
  Usecase: If AMF receives a new PDU Session Establishment Request and
  AMF has already another active PDU Session in progress with the Same PDU
  Session IDentity, AMF will not process PDUSession Establishment requests for 5
  times, after that AMF rejects with reason #43.
  */
  M5GRequestType requestType =
      static_cast<M5GRequestType>(msg->request_type.type_val);
  smf_ctx = amf_get_smf_context_by_pdu_session_id(
      ue_context, msg->payload_container.smf_msg.header.pdu_session_id);

  if ((msg->payload_container.smf_msg.header.message_type ==
       PDU_SESSION_ESTABLISHMENT_REQUEST) &&
      (requestType == M5GRequestType::INITIAL_REQUEST) && smf_ctx &&
      (smf_ctx->pdu_session_state == ACTIVE)) {
    if (smf_ctx->duplicate_pdu_session_est_req_count >=
        MAX_UE_INITIAL_PDU_SESSION_ESTABLISHMENT_REQ_ALLOWED - 1) {
      cause = M5GSmCause::INVALID_PDU_SESSION_IDENTITY;
    } else {
      smf_ctx->duplicate_pdu_session_est_req_count += 1;
      OAILOG_INFO(
          LOG_AMF_APP,
          "Duplicate Initial PDU Session establishment request received");
    }
  }
  return cause;
}
/***************************************************************************
**                                                                        **
** Name:    amf_smf_validate_context()                                    **
**                                                                        **
** Description: function to check if max PDU sessions reached             **
**                                                                        **
**                                                                        **
***************************************************************************/
M5GMmCause amf_smf_validate_context(
    amf_ue_ngap_id_t ue_id, ULNASTransportMsg* msg) {
  M5GMmCause mm_cause = M5GMmCause::UNKNOWN_CAUSE;

  ue_m5gmm_context_s* ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  M5GRequestType requestType =
      static_cast<M5GRequestType>(msg->request_type.type_val);
  /*
   * 1) the Payload container type IE is set to "N1 SM information" and
   * 2) the Request type IE is set to "initial request" or "existing PDU
   * session" the AMF determines that the PLMN's maximum number of PDU sessions
   * has already been reached for the UE, the AMF shall send back to the UE the
   * 5GSM message which was not forwarded and 5GMM cause #65
   */

  if ((N1_SM_INFO == msg->payload_container_type.type_val) &&
      ((M5GRequestType::INITIAL_REQUEST == requestType) ||
       (M5GRequestType::EXISTING_PDU_SESSION == requestType))) {
    if (ue_context->amf_context.smf_ctxt_map.size() >=
        MAX_UE_PDU_SESSION_LIMIT) {
      mm_cause = M5GMmCause::MAX_PDU_SESSIONS_REACHED;
    }
  }
  return mm_cause;
}
/***************************************************************************
**                                                                        **
** Name:    amf_validate_dnn()                                            **
**                                                                        **
** Description:                                                           **
** This function validates the DNN string received from UE or from        **
** mme.yml againt apn stored in amf_context for a particular imsi.        **
***************************************************************************/
int amf_validate_dnn(
    const amf_context_s* amf_ctxt_p, std::string dnn_string, int* index,
    bool ue_sent_dnn) {
  // Validating apn_configuration_s
  if (dnn_string.empty()) {
    return RETURNok;
  }
  for (uint8_t i = 0; i < amf_ctxt_p->apn_config_profile.nb_apns; i++) {
    if (strcmp(
            amf_ctxt_p->apn_config_profile.apn_configuration[i]
                .service_selection,
            dnn_string.c_str()) == 0) {
      *index = i;
      return RETURNok;
    }
  }
  *index = 0;
  return ue_sent_dnn ? RETURNerror : RETURNok;
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

  for (const auto& it : ue_context->amf_context.smf_ctxt_map) {
    std::shared_ptr<smf_context_t> smf_context = it.second;

    if (smf_context->pdu_address.pdn_type == IPv4) {
      char ip_str[INET_ADDRSTRLEN];

      inet_ntop(
          AF_INET, &(smf_context->pdu_address.ipv4_address.s_addr), ip_str,
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
    const std::shared_ptr<smf_context_t>& smf_ctx, paa_t* address_info) {
  OAILOG_INFO(LOG_AMF_APP, "SMF context PDU address updated\n");
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
  ue_m5gmm_context_s* ue_context;
  std::shared_ptr<smf_context_t> smf_ctx;
  imsi64_t imsi64;
  int rc = RETURNerror;

  IMSI_STRING_TO_IMSI64(response_p->imsi, &imsi64);
  ue_context = lookup_ue_ctxt_by_imsi(imsi64);

  if (ue_context == NULL) {
    OAILOG_ERROR(
        LOG_AMF_APP, "UE Context for [%s] not found \n",
        reinterpret_cast<char*>(response_p->imsi));
    return rc;
  }

  smf_ctx = amf_get_smf_context_by_pdu_session_id(
      ue_context, response_p->pdu_session_id);
  if (NULL == smf_ctx) {
    OAILOG_ERROR(
        LOG_AMF_APP, "Smf Context not found for pdu session id: [%s] \n",
        reinterpret_cast<char*>(response_p->pdu_session_id));
    return rc;
  }

  rc = amf_update_smf_context_pdu_ip(smf_ctx, &(response_p->paa));

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
        response_p->gnb_gtp_teid_ip_addr, ip_str, smf_ctx->smf_ctx_ambr);

    if (rc < 0) {
      OAILOG_ERROR(LOG_AMF_APP, "Create IPV4 Session \n");
    }
  }

  return rc;
}

int amf_send_n11_update_location_req(amf_ue_ngap_id_t ue_id) {
  OAILOG_FUNC_IN(LOG_AMF_APP);
  ue_m5gmm_context_s* ue_context_p = NULL;
  int rc                           = RETURNok;

  OAILOG_INFO(
      LOG_AMF_APP,
      "Sending UPDATE LOCATION REQ to subscriberd, ue_id = " AMF_UE_NGAP_ID_FMT,
      ue_id);

  ue_context_p = amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (ue_context_p) {
    OAILOG_INFO(
        LOG_AMF_APP, "IMSI HANDLED =%lu\n", ue_context_p->amf_context.imsi64);
  } else {
    OAILOG_ERROR(
        LOG_AMF_APP, "ue context not found for the ue_id= " AMF_UE_NGAP_ID_FMT,
        ue_id);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
  }

  s6a_update_location_req_t* s6a_ulr_p = new s6a_update_location_req_t();

  IMSI64_TO_STRING(
      ue_context_p->amf_context.imsi64, s6a_ulr_p->imsi, IMSI_LENGTH);

  s6a_ulr_p->imsi_length    = strlen(s6a_ulr_p->imsi);
  s6a_ulr_p->initial_attach = INITIAL_ATTACH;
  plmn_t visited_plmn       = {0};
  COPY_PLMN(visited_plmn, ue_context_p->amf_context.originating_tai.plmn);
  memcpy(&s6a_ulr_p->visited_plmn, &visited_plmn, sizeof(plmn_t));
  s6a_ulr_p->rat_type = RAT_NG_RAN;

  // Set regional_subscription flag
  s6a_ulr_p->supportedfeatures.regional_subscription = true;

  rc = AsyncSmfServiceClient::getInstance().n11_update_location_req(s6a_ulr_p);

  delete s6a_ulr_p;

  OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
}

/****************************************************************************
 **                                                                        **
 ** Name        :  handle_sm_message_routing_failure()                     **
 **                                                                        **
 ** Description :  Send the Downlink Transport with 5GMM Cause to gnb      **
 **                                                                        **
 ** Inputs      :  amf_ue_ngap_id_t :   pdusession response message        **
 **                ULNASTransportMsg:   received uplinktransport msg       **
 **                                                                        **
 **  Return     :  RETURNok, RETURNerror                                   **
 **                                                                        **
 ***************************************************************************/
int handle_sm_message_routing_failure(
    amf_ue_ngap_id_t ue_id, ULNASTransportMsg* ulmsg, M5GMmCause m5gmmcause) {
  int rc                   = RETURNok;
  DLNASTransportMsg* dlmsg = nullptr;
  uint32_t bytes           = 0;
  uint32_t len             = 0;
  bstring buffer;
  ue_m5gmm_context_s* ue_context = nullptr;
  amf_nas_message_t msg          = {};

  /*
        AMF shall perform if Max PDU Session limit exceeds
        a) include the PDU session ID in the PDU session ID IE;
        b) set the Payload container type IE to "N1 SM information";
        c) set the Payload container IE to the 5GSM message which was not
     forwarded; and d) set the specific 5GMM cause IE.
  */
  ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (!ue_context) {
    OAILOG_ERROR(
        LOG_AMF_APP, "UE Context not found for UE ID: " AMF_UE_NGAP_ID_FMT,
        ue_id);
    return RETURNerror;
  }

  // Message construction for PDU Establishment Reject
  // NAS-5GS (NAS) PDU
  msg.security_protected.plain.amf.header.extended_protocol_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  msg.security_protected.plain.amf.header.message_type = DLNASTRANSPORT;
  msg.header.security_header_type =
      SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED;
  msg.header.extended_protocol_discriminator = M5G_MOBILITY_MANAGEMENT_MESSAGES;
  msg.header.sequence_number =
      ue_context->amf_context._security.dl_count.seq_num;

  dlmsg = &msg.security_protected.plain.amf.msg.downlinknas5gtransport;

  // AmfHeader
  dlmsg->extended_protocol_discriminator.extended_proto_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  len++;
  dlmsg->spare_half_octet.spare  = 0x00;
  dlmsg->sec_header_type.sec_hdr = SECURITY_HEADER_TYPE_NOT_PROTECTED;
  len++;
  dlmsg->message_type.msg_type = DLNASTRANSPORT;
  len++;
  dlmsg->payload_container.iei = PAYLOAD_CONTAINER;

  // SmfMsg
  dlmsg->payload_container_type.iei      = 0;
  dlmsg->payload_container_type.type_val = N1_SM_INFO;
  len++;
  dlmsg->pdu_session_identity.iei =
      static_cast<uint8_t>(M5GIei::PDU_SESSION_IDENTITY_2);
  len++;
  dlmsg->pdu_session_identity.pdu_session_id =
      ulmsg->payload_container.smf_msg.header.pdu_session_id;
  len++;

  dlmsg->m5gmm_cause.iei         = static_cast<uint8_t>(M5GIei::M5GMM_CAUSE);
  dlmsg->m5gmm_cause.m5gmm_cause = static_cast<uint8_t>(m5gmmcause);
  len += 2;

  // Payload container IE from ulmsg
  dlmsg->payload_container.copy(ulmsg->payload_container);

  len += 2;  // 2 bytes for container.len
  len += dlmsg->payload_container.len;

  /* Ciphering algorithms, EEA1 and EEA2 expects length to be mode of 4,
   * so length is modified such that it will be mode of 4
   */
  AMF_GET_BYTE_ALIGNED_LENGTH(len);
  if (msg.header.security_header_type != SECURITY_HEADER_TYPE_NOT_PROTECTED) {
    amf_msg_header* header = &msg.security_protected.plain.amf.header;
    /*
     * Expand size of protected NAS message
     */
    len += NAS_MESSAGE_SECURITY_HEADER_SIZE;
    /*
     * Set header of plain NAS message
     */
    header->extended_protocol_discriminator = M5G_MOBILITY_MANAGEMENT_MESSAGES;
    header->security_header_type = SECURITY_HEADER_TYPE_NOT_PROTECTED;
  }

  buffer = bfromcstralloc(len, "\0");
  bytes  = nas5g_message_encode(
      buffer->data, &msg, len, &ue_context->amf_context._security);
  if (bytes > 0) {
    buffer->slen = bytes;
    rc           = amf_app_handle_nas_dl_req(ue_id, buffer, M5G_AS_SUCCESS);

  } else {
    OAILOG_WARNING(LOG_AMF_APP, "NAS encode failed \n");
    bdestroy_wrapper(&buffer);
    return RETURNerror;
  }
  return rc;
}

/****************************************************************************
 **                                                                        **
 ** Name        :  construct_pdu_session_reject_dl_req()                   **
 **                                                                        **
 ** Description :  Construct Session Establishment Reject Struct           **
 **                                                                        **
 ** Inputs      :  sequence_number     : seq to construct Secure msgs      **
 **                session_id          : PDU Session Identity              **
 **                pti                 : Procedure transaction identity    **
 **                is_security_enabled : indcates to construct plain msg   **
 **                                      secure msg                        **
 **                msg                 : out parameter Session             **
 **                                      Establishment Reject              **
 **                                                                        **
 **  Return     :  len                 : buffer required for               **
 **                                      DLNASTransportMsg                 **
 **                                                                        **
 ***************************************************************************/

int construct_pdu_session_reject_dl_req(
    uint8_t sequence_number, uint8_t session_id, uint8_t pti, uint8_t cause,
    bool is_security_enabled, amf_nas_message_t* msg) {
  uint32_t len             = 0;
  uint32_t container_len   = 0;
  DLNASTransportMsg* dlmsg = nullptr;

  if (nullptr == msg) {
    return 0;
  }
  // Message construction for PDU Establishment Reject
  // NAS-5GS (NAS) PDU
  if (is_security_enabled) {
    msg->security_protected.plain.amf.header.extended_protocol_discriminator =
        M5G_MOBILITY_MANAGEMENT_MESSAGES;
    msg->security_protected.plain.amf.header.message_type = DLNASTRANSPORT;
    msg->header.security_header_type =
        SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED_CYPHERED;
    dlmsg = &msg->security_protected.plain.amf.msg.downlinknas5gtransport;
  } else {
    msg->plain.amf.header.extended_protocol_discriminator =
        M5G_MOBILITY_MANAGEMENT_MESSAGES;
    msg->plain.amf.header.message_type = DLNASTRANSPORT;
    msg->header.security_header_type   = SECURITY_HEADER_TYPE_NOT_PROTECTED;
    dlmsg = &msg->plain.amf.msg.downlinknas5gtransport;
  }

  msg->header.extended_protocol_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  msg->header.sequence_number = sequence_number;

  // AmfHeader
  dlmsg->extended_protocol_discriminator.extended_proto_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  len++;
  dlmsg->spare_half_octet.spare  = 0x00;
  dlmsg->sec_header_type.sec_hdr = SECURITY_HEADER_TYPE_NOT_PROTECTED;
  len++;
  dlmsg->message_type.msg_type = DLNASTRANSPORT;
  len++;
  dlmsg->payload_container.iei = PAYLOAD_CONTAINER;

  // SmfMsg
  dlmsg->payload_container_type.iei      = 0;
  dlmsg->payload_container_type.type_val = N1_SM_INFO;
  len++;
  dlmsg->pdu_session_identity.iei =
      static_cast<uint8_t>(M5GIei::PDU_SESSION_IDENTITY_2);
  len++;
  dlmsg->pdu_session_identity.pdu_session_id = session_id;
  len++;

  SmfMsg& pdu_sess_est_reject = dlmsg->payload_container.smf_msg;
  // header
  pdu_sess_est_reject.header.extended_protocol_discriminator =
      M5G_SESSION_MANAGEMENT_MESSAGES;
  pdu_sess_est_reject.header.pdu_session_id = session_id;
  pdu_sess_est_reject.header.message_type   = PDU_SESSION_ESTABLISHMENT_REJECT;
  pdu_sess_est_reject.header.procedure_transaction_id = pti;

  // Smf NAS message
  pdu_sess_est_reject.msg.pdu_session_estab_reject
      .extended_protocol_discriminator.extended_proto_discriminator =
      M5G_SESSION_MANAGEMENT_MESSAGES;
  container_len++;
  pdu_sess_est_reject.msg.pdu_session_estab_reject.pdu_session_identity
      .pdu_session_id = session_id;
  container_len++;
  pdu_sess_est_reject.msg.pdu_session_estab_reject.pti.pti = pti;
  container_len++;
  pdu_sess_est_reject.msg.pdu_session_estab_reject.message_type.msg_type =
      PDU_SESSION_ESTABLISHMENT_REJECT;
  container_len++;
  pdu_sess_est_reject.msg.pdu_session_estab_reject.m5gsm_cause.cause_value =
      cause;
  container_len++;

  dlmsg->payload_container.len = container_len;
  len += PAYLOAD_CONTAINER_TAG_LENGTH;
  len += dlmsg->payload_container.len;

  /* Ciphering algorithms, EEA1 and EEA2 expects length to be mode of 4,
   * so length is modified such that it will be mode of 4
   */
  AMF_GET_BYTE_ALIGNED_LENGTH(len);
  if (SECURITY_HEADER_TYPE_NOT_PROTECTED != msg->header.security_header_type) {
    amf_msg_header* header = &msg->security_protected.plain.amf.header;
    /*
     * Expand size of protected NAS message
     */
    len += NAS_MESSAGE_SECURITY_HEADER_SIZE;
    /*
     * Set header of plain NAS message
     */
    header->extended_protocol_discriminator = M5G_MOBILITY_MANAGEMENT_MESSAGES;
    header->security_header_type = SECURITY_HEADER_TYPE_NOT_PROTECTED;
  }

  return len;
}

/****************************************************************************
 **                                                                        **
 ** Name        :  amf_pdu_session_establishment_reject()                  **
 **                                                                        **
 ** Description :  Send the Downlink Transport with 5GMM Cause to gnb      **
 **                                                                        **
 ** Inputs      :  amf_ue_ngap_id_t : pdusession response message          **
 **                session_id       : PDU Session Inputs                   **
 **                pti              : Procedure transaction identity       **
 **                5GSM cause       : 5GSM cause                           **
 **                                                                        **
 **  Return     :  RETURNok, RETURNerror                                   **
 **                                                                        **
 ***************************************************************************/
int amf_pdu_session_establishment_reject(
    amf_ue_ngap_id_t ue_id, uint8_t session_id, uint8_t pti, uint8_t cause) {
  int rc         = RETURNok;
  uint32_t bytes = 0;
  uint32_t len   = 0;
  bstring buffer;
  ue_m5gmm_context_s* ue_context = nullptr;
  amf_nas_message_t msg          = {};

  ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);

  if (!ue_context) {
    OAILOG_ERROR(
        LOG_AMF_APP, "UE Context not found for UE ID: " AMF_UE_NGAP_ID_FMT,
        ue_id);
    return RETURNerror;
  }

  len = construct_pdu_session_reject_dl_req(
      ue_context->amf_context._security.dl_count.seq_num, session_id, pti,
      cause, true, &msg);

  if (len <= 0) {
    OAILOG_WARNING(LOG_AMF_APP, "PDU Construction is failed \n");
    return RETURNerror;
  }
  buffer = bfromcstralloc(len, "\0");
  bytes  = nas5g_message_encode(
      buffer->data, &msg, len, &ue_context->amf_context._security);

  if (bytes > 0) {
    buffer->slen = bytes;
    rc           = amf_app_handle_nas_dl_req(ue_id, buffer, M5G_AS_SUCCESS);

  } else {
    OAILOG_WARNING(LOG_AMF_APP, "NAS encode failed \n");
    bdestroy_wrapper(&buffer);
    rc = RETURNerror;
  }
  return rc;
}

}  // namespace magma5g
