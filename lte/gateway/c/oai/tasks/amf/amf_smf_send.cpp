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
#include "amf_app_ue_context_and_proc.h"
#include "amf_recv.h"
#include "amf_smfDefs.h"
#include "SmfMessage.h"
#include "M5gNasMessage.h"
#include "M5GCommonDefs.h"
#include "amf_common_defs.h"
#include "M5GULNASTransport.h"
#include "SmfServiceClient.h"

using namespace std;
namespace magma5g {
// Check Procedure Transcation Identiy Valid or Invalid
int esm_pt_is_reserved(int pti) {
  return (
      (pti != PROCEDURE_TRANSACTION_IDENTITY_UNASSIGNED_t) &&
      (pti > PROCEDURE_TRANSACTION_IDENTITY_LAST_t));
}

// SMF procedure handler for PDU Establishment Request
int amf_smf_procedure_handler::amf_smf_handle_pdu_establishment_request(
    SmfMsg* msg, amf_smf_t* amf_smf_msg) {
  int rc        = RETURNerror;
  int smf_cause = SMF_CAUSE_SUCCESS;

  OAILOG_INFO(
      LOG_AMF_APP,
      "AMF SMF Handler- Received PDN Connectivity Request message ");

  // Procedure transaction identity checking
  if ((msg->header.procedure_transaction_id ==
       PROCEDURE_TRANSACTION_IDENTITY_UNASSIGNED_t) ||
      esm_pt_is_reserved(msg->header.procedure_transaction_id)) {
    // OAILOG_ERROR(LOG_NAS_AMF,
    // "AMF SMF Handler   - Invalid PTI value pti= %x",
    // msg->header.procedure_transaction_id);//TODO uncomment log
    amf_smf_msg->u.establish.cause_value = SMF_CAUSE_INVALID_PTI_VALUE;
    return (amf_smf_msg->u.establish.cause_value);
  } else {
    amf_smf_msg->u.establish.pti = msg->header.procedure_transaction_id;
  }

  // Get the value of the PDN type indicator
  if (msg->pdu_session_estab_request.pdu_session_type.type_val ==
      PDN_TYPE_IPV4) {
    amf_smf_msg->u.establish.pdu_session_type = NET_PDN_TYPE_IPV4;
  } else if (
      msg->pdu_session_estab_request.pdu_session_type.type_val ==
      PDN_TYPE_IPV6) {
    amf_smf_msg->u.establish.pdu_session_type = NET_PDN_TYPE_IPV6;
  } else if (
      msg->pdu_session_estab_request.pdu_session_type.type_val ==
      PDN_TYPE_IPV4V6) {
    amf_smf_msg->u.establish.pdu_session_type = NET_PDN_TYPE_IPV4V6;
  } else {
    // Unkown PDN type
    amf_smf_msg->u.establish.cause_value = SMF_CAUSE_UNKNOWN_PDN_TYPE;
    return (amf_smf_msg->u.establish.cause_value);
  }
  amf_smf_msg->u.establish.pdu_session_id = msg->header.pdu_session_id;
  amf_smf_msg->u.establish.cause_value    = smf_cause;

  // Return the ESM cause value
  return (smf_cause);
}

// SMF procedure for PDU session release
int amf_smf_procedure_handler::amf_smf_handle_pdu_release_request(
    SmfMsg* msg, amf_smf_t* amf_smf_msg) {
  int smf_cause                         = SMF_CAUSE_SUCCESS;
  amf_smf_msg->u.release.pti            = msg->header.procedure_transaction_id;
  amf_smf_msg->u.release.pdu_session_id = msg->header.pdu_session_id;
  amf_smf_msg->u.release.cause_value    = smf_cause;
  return (smf_cause);  // TODO add error checking as needed and return
                       // appropriate cause value
}

// Send PDU Session Reject Message
int amf_send_pdusession_reject(
    SmfMsg* reject_req, uint8_t session_id, uint8_t pti, uint8_t cause) {
  uint8_t buffer[5];
  int rc;
  reject_req->header.extended_protocol_discriminator           = 0x2e;
  reject_req->header.pdu_session_id                            = session_id;
  reject_req->header.procedure_transaction_id                  = pti;
  reject_req->header.message_type                              = 0xC3;
  reject_req->pdu_session_estab_reject.m5gsm_cause.cause_value = cause;

  rc = reject_req->SmfMsgEncodeMsg(reject_req, buffer, 5);
  if (rc > 0) {
    // TODO
    // Send the message to AS
  }

  return rc;
}

extern ue_m5gmm_context_s
    ue_m5gmm_global_context;  // TODO AMF_TEST global var to temporarily store
                              // context inserted to ht
void set_amf_smf_context(
    PDUSessionEstablishmentRequestMsg* message, smf_context_t* smf_ctx) {
  smf_ctx->smf_proc_data.pdu_session_identity = message->pdu_session_identity;
  smf_ctx->smf_proc_data.pti                  = message->pti;
  smf_ctx->smf_proc_data.message_type         = message->message_type;
  smf_ctx->smf_proc_data.integrity_prot_max_data_rate =
      message->integrity_prot_max_data_rate;
  smf_ctx->smf_proc_data.pdu_session_type = message->pdu_session_type;
  smf_ctx->smf_proc_data.ssc_mode         = message->ssc_mode;
  smf_ctx->pdu_session_version = 0; //Initializing pdu version with 0
  memset(smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr, '\0', 
		  sizeof(smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr));
  memset(smf_ctx->gtp_tunnel_id.gnb_gtp_teid, '\0', 
		  sizeof(smf_ctx->gtp_tunnel_id.gnb_gtp_teid));

  //Removing hard coded values.
  //  strcpy(smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr, "10.32.116.1");
  //uint8_t buff_ip[] = {0x0a, 0x20, 0x74, 0x01};
  //memcpy(smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr, buff_ip, 4);
  //uint8_t buff_teid[] = {0x01, 0x00, 0x00, 0x09};
  //memcpy(smf_ctx->gtp_tunnel_id.gnb_gtp_teid, buff_teid,
  //    4);  // TODO get the gnb_gtp_teid_ip_addr and gnb_gtp_teid from
           // PDUSessionResourceSetupResponse
}

void clear_amf_smf_context(smf_context_t* smf_ctx) {
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
// Send SMF Message through grpc to SMF

int amf_procedure_handler::amf_smf_send(
    amf_ue_ngap_id_t ue_id, ULNASTransportMsg* msg, int amf_cause) {
  int decode_rc;
  int rc = 1;
  SmfMsg reject_req;
  amf_smf_procedure_handler procedure_handler;
  amf_smf_t amf_smf_msg = {};
  char imsi[IMSI_BCD_DIGITS_MAX + 1];

  amf_cause =
      SMF_CAUSE_SUCCESS;  // TODO SMF_CAUSE_SUCCESS and AMF_CAUSE_SUCCESS values
                          // differ, make it inline
  if (amf_cause != SMF_CAUSE_SUCCESS) {
    rc = amf_send_pdusession_reject(
        &reject_req, msg->payload_container.smf_msg.header.pdu_session_id,
        msg->payload_container.smf_msg.header.procedure_transaction_id,
        amf_cause);
    return rc;
  }

  ue_m5gmm_context_s* ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  smf_context_t* smf_ctx;
  if (ue_context) {
    smf_ctx = &(ue_context->amf_context.smf_context);
    IMSI64_TO_STRING(
        ue_context->amf_context._imsi64, imsi,
        // ue_m5gmm_global_context.amf_context._imsi.length);
        15);
  } else {
    ue_context = &ue_m5gmm_global_context;
    smf_ctx    = &ue_m5gmm_global_context.amf_context
                   .smf_context;  // TODO AMF_TEST global var to temporarily
                                  // store context inserted to ht
    IMSI64_TO_STRING(
        ue_m5gmm_global_context.amf_context._imsi64, imsi,
        // ue_m5gmm_global_context.amf_context._imsi.length);
        15);
  }
  // Process initial NAS message
  switch (msg->payload_container.smf_msg.header.message_type) {
    case PDU_SESSION_ESTABLISHMENT_REQUEST: {
      amf_cause = procedure_handler.amf_smf_handle_pdu_establishment_request(
          &(msg->payload_container.smf_msg), &amf_smf_msg);

      if (amf_cause != SMF_CAUSE_SUCCESS) {
        rc = amf_send_pdusession_reject(
            &reject_req, msg->payload_container.smf_msg.header.pdu_session_id,
            msg->payload_container.smf_msg.header.procedure_transaction_id,
            amf_cause);
        return rc;
      }
      set_amf_smf_context(
          &(msg->payload_container.smf_msg.pdu_session_estab_request), smf_ctx);
#if 0
      strncpy(
          (char*) amf_smf_msg.u.establish.gnb_gtp_teid_ip_addr,
          (char*)smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr, 5);
      // memcpy(amf_smf_msg.u.establish.gnb_gtp_teid,
      // smf_ctx->gtp_tunnel_id.gnb_gtp_teid, 4);//TODO get the
      // gnb_gtp_teid_ip_addr and gnb_gtp_teid from
      // PDUSessionResourceSetupResponse
      strncpy(
          (char*) amf_smf_msg.u.establish.gnb_gtp_teid,
          (char*)smf_ctx->gtp_tunnel_id.gnb_gtp_teid,
          5);  // TODO get the gnb_gtp_teid_ip_addr and gnb_gtp_teid from
               // PDUSessionResourceSetupResponse
#endif
      memset(
          amf_smf_msg.u.establish.gnb_gtp_teid_ip_addr, '\0',
          sizeof(amf_smf_msg.u.establish.gnb_gtp_teid_ip_addr));
      memset(
          amf_smf_msg.u.establish.gnb_gtp_teid, '\0',
          sizeof(amf_smf_msg.u.establish.gnb_gtp_teid));
      memcpy(
          amf_smf_msg.u.establish.gnb_gtp_teid_ip_addr,
          smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr, 4);
      // memcpy(amf_smf_msg.u.establish.gnb_gtp_teid,
      // smf_ctx->gtp_tunnel_id.gnb_gtp_teid, 4);//TODO get the
      // gnb_gtp_teid_ip_addr and gnb_gtp_teid from
      // PDUSessionResourceSetupResponse
      // TODO fix string null_term for bytes
      memcpy(
          amf_smf_msg.u.establish.gnb_gtp_teid,
          smf_ctx->gtp_tunnel_id.gnb_gtp_teid,
          4);  // TODO get the gnb_gtp_teid_ip_addr and gnb_gtp_teid from
               // PDUSessionResourceSetupResponse
               // TODO fix string null_term for bytes

      // Invoke Grpc Send call
      rc = create_session_grpc_req(&amf_smf_msg.u.establish, imsi);
    } break;
    case PDU_SESSION_RELEASE_REQUEST: {
      amf_cause = procedure_handler.amf_smf_handle_pdu_release_request(
          &(msg->payload_container.smf_msg), &amf_smf_msg);
      OAILOG_INFO(
          LOG_AMF_APP,
          "sending PDU session resource release request to gNB \n");
      rc = pdu_session_resource_release_request(ue_context, ue_id);
      if (rc != RETURNok) {
        OAILOG_INFO(
            LOG_AMF_APP,
            "PDU session resource release request to gNB failed"
            "\n");
      }
      // Invoke Grpc Send call
      OAILOG_INFO(LOG_AMF_APP, "Prepare PDU session release request to SMF");
      rc = release_session_gprc_req(&amf_smf_msg.u.release, imsi);
      OAILOG_INFO(
          LOG_AMF_APP, "Releasing saved context associated with the PDU");
      clear_amf_smf_context(smf_ctx);
    } break;
    case PDU_SESSION_MODIFICATION_REQUEST:
      amf_cause = procedure_handler.amf_smf_handle_pdu_modif_request(
          &(msg->payload_container.smf_msg), &amf_smf_msg);

      // Invoke Grpc Send call
      OAILOG_INFO(
          LOG_AMF_APP,
          "Prepare PDU session modification request to SMF");
      rc = mod_sessionreq_grpc_req(&amf_smf_msg.u.modif, imsi);
      OAILOG_INFO(
          LOG_AMF_APP,
          "Releasing saved context associated with the PDU");
      clear_amf_smf_context(smf_ctx);   
      break;
    case PDU_SESSION_MODIFICATION_COMPLETE:
      amf_cause = procedure_handler.amf_smf_handle_pdu_modif_complete(
          &(msg->payload_container.smf_msg), &amf_smf_msg);

      // Invoke Grpc Send call
      OAILOG_INFO(
          LOG_AMF_APP,
          "Prepare PDU session modification complete to SMF");
      rc = mod_sessioncomp_grpc_req(&amf_smf_msg.u.modif, imsi);
      OAILOG_INFO(
          LOG_AMF_APP,
          "Releasing saved context associated with the PDU");
      clear_amf_smf_context(smf_ctx);
      break;
  case PDU_SESSION_MODIFICATION_COMMAND_REJECT:
      amf_cause = procedure_handler.amf_smf_handle_pdu_modif_cmd_reject(
          &(msg->payload_container.smf_msg), &amf_smf_msg);

      // Invoke Grpc Send call
      OAILOG_INFO(
          LOG_AMF_APP,
          "Prepare PDU session modification command reject to SMF");
      rc = mod_sessioncmd_reject_grpc_req(&amf_smf_msg.u.modif, imsi);
      OAILOG_INFO(
          LOG_AMF_APP,
          "Releasing saved context associated with the PDU");
      clear_amf_smf_context(smf_ctx);
      break;
    default:
      break;
  }
#if 0
  if (amf_cause != SMF_CAUSE_SUCCESS) {
    rc = amf_send_pdusession_reject(
        &reject_req, msg->payload_container.smf_msg.header.pdu_session_id,
        msg->payload_container.smf_msg.header.procedure_transaction_id,
        amf_cause);
    return rc;
  }

  // TODO: State Machine Handler

  ue_m5gmm_context_s* ue_context = amf_ue_context_exists_amf_ue_ngap_id(ue_id);
  if (ue_context) {
    smf_context_t* smf_ctx = &(ue_context->amf_context.smf_context);
    IMSI64_TO_STRING(
        ue_context->amf_context._imsi64, imsi,
        ue_context->amf_context._imsi.length);
    set_amf_smf_context(
        &(msg->payload_container.smf_msg.pdu_session_estab_request), smf_ctx);
  } else {
	  ue_context = &ue_m5gmm_global_context;
    smf_context_t* smf_ctx =
        &ue_m5gmm_global_context.amf_context
             .smf_context;  // TODO AMF_TEST global var to temporarily store
                            // context inserted to ht
    //IMSI64_TO_STRING(
      //  ue_m5gmm_global_context.amf_context._imsi64, imsi,
        //ue_m5gmm_global_context.amf_context._imsi.length);

    IMSI64_TO_STRING(ue_m5gmm_global_context.amf_context._imsi64, imsi, 15);
    //  strcpy(smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr, "10.32.116.1");
    uint8_t buff_ip[] = {0x0a, 0x20, 0x74, 0x01};
    memcpy(smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr, buff_ip, 4);
    uint8_t  buff_teid[] = {0x01, 0x00, 0x00, 0x09};
    memcpy(
        smf_ctx->gtp_tunnel_id.gnb_gtp_teid, buff_teid,
        4);  // TODO get the gnb_gtp_teid_ip_addr and gnb_gtp_teid from
             // PDUSessionResourceSetupResponse

    set_amf_smf_context(
        &(msg->payload_container.smf_msg.pdu_session_estab_request), smf_ctx);
memset(amf_smf_msg.u.establish.gnb_gtp_teid_ip_addr, '\0', sizeof(amf_smf_msg.u.establish.gnb_gtp_teid_ip_addr));
memset(amf_smf_msg.u.establish.gnb_gtp_teid, '\0', sizeof(amf_smf_msg.u.establish.gnb_gtp_teid));
    memcpy(
        amf_smf_msg.u.establish.gnb_gtp_teid_ip_addr,
        smf_ctx->gtp_tunnel_id.gnb_gtp_teid_ip_addr, 4);
    // memcpy(amf_smf_msg.u.establish.gnb_gtp_teid,
    // smf_ctx->gtp_tunnel_id.gnb_gtp_teid, 4);//TODO get the
    // gnb_gtp_teid_ip_addr and gnb_gtp_teid from
    // PDUSessionResourceSetupResponse
    //TODO fix string null_term for bytes
    memcpy(
        amf_smf_msg.u.establish.gnb_gtp_teid,
        smf_ctx->gtp_tunnel_id.gnb_gtp_teid,
        4);  // TODO get the gnb_gtp_teid_ip_addr and gnb_gtp_teid from
             // PDUSessionResourceSetupResponse
	     //TODO fix string null_term for bytes
  }
//  amf_smf_msg.u.establish.gnb_gtp_teid_ip_addr[4] = '\0';//TODO fix string null_term for bytes
//  amf_smf_msg.u.establish.gnb_gtp_teid[4] = '\0';//TODO fix string null_term for bytes
  /*
   * Before sending establishment request to SMF, AMF will send resource
   * setup request to NGAP/gNB with available data in smf pdu session
   * context or default values.
   * TODO Note: few values have been hard-coded
   */
  // rc = pdu_session_resource_setup_request(ue_context, ue_id);
  rc = RETURNok;
  if (rc != RETURNok) {
    OAILOG_INFO(
        LOG_AMF_APP,
        "PDU session resource request to gNB failed and no message sent to SMF "
        "\n");
    // OAILOG_FUNC_RETURN(LOG_NAS_AMF, rc);
    OAILOG_FUNC_RETURN(LOG_AMF_APP, rc);
  }
  // Invoke Grpc Send call
  rc = create_session_grpc_req(&amf_smf_msg.u.establish, imsi);
#endif
  // TODO
  // Release the buffer
  return rc;
}

/* This function for UE idle event notification to SMF or single PDU
 * session state change to Inactive state and notify to SMF.
 * 4 types of events are used in proto. 
 * PDU_SESSION_INACTIVE_NOTIFY => use for single PDU session notify
 * UE_IDLE_MODE_NOTIFY     => use for idle mode support
 * UE_PAGING_NOTIFY
 * UE_PERIODIC_REG_ACTIVE_MODE_NOTIFY
 */

int amf_procedure_handler::amf_smf_notification_send(
    amf_ue_ngap_id_t ue_id, ue_m5gmm_context_s* ue_context) {
    OAILOG_INFO(LOG_AMF_APP,
          "AMF_TEST: Preparing and sending idle notification to SMF \n");
   int rc = RETURNerror;
   /* Get gRPC structure of notification to be filled common and 
    * rat type elements. 
    * Only need  to be filled IMSI and ue_state_idle of UE
    */
   magma::lte::SetSmNotificationContext notify_req;
   auto* req_common = notify_req.mutable_common_context();
   auto* req_rat_specific = notify_req.mutable_rat_specific_notification();
   char imsi[IMSI_BCD_DIGITS_MAX + 1];
   IMSI64_TO_STRING(ue_context->amf_context._imsi64, imsi,15);

   req_common->mutable_sid()->mutable_id()->assign(imsi);
   //req_rat_specific->set_ue_state_idle(true);
   req_rat_specific->set_notify_ue_event(
		   magma::lte::NotifyUeEvents::UE_IDLE_MODE_NOTIFY);
   OAILOG_INFO(LOG_AMF_APP,
          "AMF_TEST: Notification gRPC filled with IMSI %s and "
	  "ue_state_idle is set to true \n", imsi);
   
   auto smf_srv_client = std::make_shared<magma5g::AsyncSmfServiceClient>();
   std::thread smf_srv_client_response_handling_thread(
       [&]() { smf_srv_client->rpc_response_loop(); });
   smf_srv_client_response_handling_thread.detach();

   OAILOG_INFO(LOG_AMF_APP,
          "AMF_TEST: Sending filled idle notification to SMF by gRPC \n");
   smf_srv_client->set_smf_notification(notify_req);
   return RETURNok;
}

// SMF procedure handler for PDU Modification Request
int amf_smf_procedure_handler::amf_smf_handle_pdu_modif_request(
    SmfMsg* msg, amf_smf_t* amf_smf_msg) {
  int rc        = RETURNerror;
  int smf_cause = SMF_CAUSE_SUCCESS;

  OAILOG_INFO(
      LOG_AMF_APP,
      "AMF SMF Handler- Received PDU Modification Request message ");

  // Procedure transaction identity checking
  if ((msg->header.procedure_transaction_id ==
       PROCEDURE_TRANSACTION_IDENTITY_UNASSIGNED_t) ||
      esm_pt_is_reserved(msg->header.procedure_transaction_id)) {
     OAILOG_ERROR(LOG_NAS_AMF,
     "AMF SMF Handler   - Invalid PTI value pti= %x",
     msg->header.procedure_transaction_id);
    amf_smf_msg->u.modif.cause_value = SMF_CAUSE_INVALID_PTI_VALUE;
    return (amf_smf_msg->u.modif.cause_value);
  } else {
    amf_smf_msg->u.modif.pti = msg->header.procedure_transaction_id;
  }

  amf_smf_msg->u.modif.cause_value    = smf_cause;

  // Return the ESM cause value
  return (smf_cause);
}

// SMF procedure handler for PDU Modification Complete
int amf_smf_procedure_handler::amf_smf_handle_pdu_modif_complete(
    SmfMsg* msg, amf_smf_t* amf_smf_msg) {
  int rc        = RETURNerror;
  int smf_cause = SMF_CAUSE_SUCCESS;

  OAILOG_INFO(
      LOG_AMF_APP,
      "AMF SMF Handler- Received PDU Modification Complete message ");

  // Procedure transaction identity checking
  if ((msg->header.procedure_transaction_id ==
       PROCEDURE_TRANSACTION_IDENTITY_UNASSIGNED_t) ||
      esm_pt_is_reserved(msg->header.procedure_transaction_id)) {
     OAILOG_ERROR(LOG_NAS_AMF,
     "AMF SMF Handler   - Invalid PTI value pti= %x",
     msg->header.procedure_transaction_id);
    amf_smf_msg->u.modif.cause_value = SMF_CAUSE_INVALID_PTI_VALUE;
    return (amf_smf_msg->u.modif.cause_value);
  } else {
    amf_smf_msg->u.modif.pti = msg->header.procedure_transaction_id;
  }

  amf_smf_msg->u.modif.cause_value    = smf_cause;

  // Return the ESM cause value
  return (smf_cause);
}

// SMF procedure handler for PDU Modification Command Reject
int amf_smf_procedure_handler::amf_smf_handle_pdu_modif_cmd_reject(
    SmfMsg* msg, amf_smf_t* amf_smf_msg) {
  int rc        = RETURNerror;
  int smf_cause = SMF_CAUSE_SUCCESS;

  OAILOG_INFO(
      LOG_AMF_APP,
      "AMF SMF Handler- Received PDU Modification Command Reject message ");

  // Procedure transaction identity checking
  if ((msg->header.procedure_transaction_id ==
       PROCEDURE_TRANSACTION_IDENTITY_UNASSIGNED_t) ||
      esm_pt_is_reserved(msg->header.procedure_transaction_id)) {
     OAILOG_ERROR(LOG_NAS_AMF,
     "AMF SMF Handler   - Invalid PTI value pti= %x",
     msg->header.procedure_transaction_id);
    amf_smf_msg->u.modif.cause_value = SMF_CAUSE_INVALID_PTI_VALUE;
    return (amf_smf_msg->u.modif.cause_value);
  } else {
    amf_smf_msg->u.modif.pti = msg->header.procedure_transaction_id;
  }

  amf_smf_msg->u.modif.cause_value    = smf_cause;

  // Return the ESM cause value
  return (smf_cause);
}
}  // namespace magma5g
