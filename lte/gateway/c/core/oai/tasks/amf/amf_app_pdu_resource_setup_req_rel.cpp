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
#include <thread>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface_types.h"
#include "lte/gateway/c/core/oai/lib/itti/intertask_interface.h"
#include "lte/gateway/c/core/oai/lib/directoryd/directoryd.h"
#include "lte/gateway/c/core/oai/common/conversions.h"
#include "lte/gateway/c/core/oai/lib/bstr/bstrlib.h"
#include "lte/gateway/c/core/oai/common/dynamic_memory_check.h"
#ifdef __cplusplus
};
#endif
#include "lte/gateway/c/core/oai/common/common_defs.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"
#include "lte/gateway/c/core/oai/include/ngap_messages_types.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_common.h"
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_defs.h"
#include "lte/gateway/c/core/oai/tasks/amf/include/amf_smf_session_context.h"

namespace magma5g {
extern task_zmq_ctx_t amf_app_task_zmq_ctx;

uint64_t get_bit_rate(uint8_t ambr_unit) {
  if (ambr_unit < 6) {
    return (1000);
  } else if (ambr_unit < 11) {
    return (1000 * 1000);
  } else if (ambr_unit < 16) {
    return (1000 * 1000 * 1000);
  }
  return (0);
}

/*
 * AMBR calculation based on 9.11.4.14 of 24-501
 */
void ambr_calculation_pdu_session(
    uint16_t* dl_session_ambr, uint8_t* dl_ambr_unit, uint16_t* ul_session_ambr,
    uint8_t* ul_ambr_unit, uint64_t* dl_pdu_ambr, uint64_t* ul_pdu_ambr) {
  *dl_pdu_ambr = (*dl_session_ambr) * get_bit_rate(*dl_ambr_unit);

  *ul_pdu_ambr = (*ul_session_ambr) * get_bit_rate(*ul_ambr_unit);
}

/*
 * the function to be called before sending gRPC message to SMF with available
 * information in smf_context or default. No NAS message is involved and direct
 * itti_message is sent to NGAP.
 */
int pdu_session_resource_setup_request(
    ue_m5gmm_context_s* ue_context, amf_ue_ngap_id_t amf_ue_ngap_id,
    std::shared_ptr<smf_context_t> smf_context, bstring nas_msg) {
  pdu_session_resource_setup_request_transfer_t*
      amf_pdu_ses_setup_transfer_req                                = nullptr;
  itti_ngap_pdusession_resource_setup_req_t* ngap_pdu_ses_setup_req = nullptr;
  MessageDef* message_p                                             = nullptr;
  uint64_t dl_pdu_ambr;
  uint64_t ul_pdu_ambr;

  message_p =
      itti_alloc_new_message(TASK_AMF_APP, NGAP_PDUSESSION_RESOURCE_SETUP_REQ);
  ngap_pdu_ses_setup_req =
      &message_p->ittiMsg.ngap_pdusession_resource_setup_req;
  memset(
      ngap_pdu_ses_setup_req, 0,
      sizeof(itti_ngap_pdusession_resource_setup_req_t));

  // start filling message in DL to NGAP
  ngap_pdu_ses_setup_req->gnb_ue_ngap_id = ue_context->gnb_ue_ngap_id;
  ngap_pdu_ses_setup_req->amf_ue_ngap_id = amf_ue_ngap_id;

  /*
   * by this time amf and smf context available but ambr for pdu not available
   * considering default or max bit rate.
   * leveraged ambr calculation from qos_params_to_eps_qos and 24-501 spec used
   */
  ambr_calculation_pdu_session(
      &(smf_context->selected_ambr.dl_session_ambr),
      &(smf_context->selected_ambr.dl_ambr_unit),
      &(smf_context->selected_ambr.ul_session_ambr),
      &(smf_context->selected_ambr.ul_ambr_unit), &dl_pdu_ambr, &ul_pdu_ambr);

  amf_smf_context_ue_aggregate_max_bit_rate_get(
      &(ue_context->amf_context),
      &(ngap_pdu_ses_setup_req->ue_aggregate_maximum_bit_rate.dl),
      &(ngap_pdu_ses_setup_req->ue_aggregate_maximum_bit_rate.ul));

  // Hardcoded number of pdu sessions as 1
  ngap_pdu_ses_setup_req->pduSessionResource_setup_list.no_of_items = 1;
  ngap_pdu_ses_setup_req->pduSessionResource_setup_list.item[0].Pdu_Session_ID =
      (Ngap_PDUSessionID_t) smf_context->smf_proc_data.pdu_session_id;

  ngap_pdu_ses_setup_req->pduSessionResource_setup_list.no_of_items = 1;
  // Adding respective header to amf_pdu_ses_setup_transfer_request
  amf_pdu_ses_setup_transfer_req =
      &ngap_pdu_ses_setup_req->pduSessionResource_setup_list.item[0]
           .PDU_Session_Resource_Setup_Request_Transfer;

  /* preparing for PDU_Session_Resource_Setup_Transfer.
   * amf_pdu_ses_setup_transfer_req is the structure to be filled.
   */
  amf_pdu_ses_setup_transfer_req->pdu_aggregate_max_bit_rate.dl = dl_pdu_ambr;
  amf_pdu_ses_setup_transfer_req->pdu_aggregate_max_bit_rate.ul = ul_pdu_ambr;

  // UPF teid 4 octet and respective ip address are from SMF context
  memcpy(
      &amf_pdu_ses_setup_transfer_req->up_transport_layer_info.gtp_tnl.gtp_tied,
      smf_context->gtp_tunnel_id.upf_gtp_teid, GNB_TEID_LEN);
  amf_pdu_ses_setup_transfer_req->up_transport_layer_info.gtp_tnl
      .endpoint_ip_address = blk2bstr(
      &smf_context->gtp_tunnel_id.upf_gtp_teid_ip_addr, GNB_IPV4_ADDR_LEN);
  amf_pdu_ses_setup_transfer_req->pdu_ip_type.pdn_type = IPv4;

  memcpy(
      &amf_pdu_ses_setup_transfer_req->qos_flow_setup_request_list
           .qos_flow_req_item,
      &smf_context->subscribebed_qos_profile.qos_flow_req_item,
      sizeof(qos_flow_setup_request_item));

  ngap_pdu_ses_setup_req->nas_pdu = nas_msg;

  // Send message to NGAP task
  amf_send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);

  return RETURNok;
}

/* Resource release request to gNB through NGAP */
int pdu_session_resource_release_request(
    ue_m5gmm_context_s* ue_context, amf_ue_ngap_id_t amf_ue_ngap_id,
    std::shared_ptr<smf_context_t> smf_ctx, bool retransmit) {
  bstring buffer;
  uint32_t bytes                = 0;
  DLNASTransportMsg* encode_msg = NULL;
  SmfMsg* smf_msg               = NULL;
  uint32_t len                  = 0;
  uint32_t container_len        = 0;
  amf_nas_message_t msg;
  nas5g_error_code_t rc = M5G_AS_SUCCESS;

  memset(&msg, 0, sizeof(amf_nas_message_t));

  msg.security_protected.plain.amf.header.extended_protocol_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  msg.security_protected.plain.amf.header.message_type = DLNASTRANSPORT;
  msg.header.security_header_type = SECURITY_HEADER_TYPE_INTEGRITY_PROTECTED;
  msg.header.extended_protocol_discriminator = M5G_MOBILITY_MANAGEMENT_MESSAGES;
  msg.header.sequence_number =
      ue_context->amf_context._security.dl_count.seq_num;

  encode_msg = &msg.security_protected.plain.amf.msg.downlinknas5gtransport;
  smf_msg    = &encode_msg->payload_container.smf_msg;

  // NAS5g AmfHeader
  encode_msg->extended_protocol_discriminator.extended_proto_discriminator =
      M5G_MOBILITY_MANAGEMENT_MESSAGES;
  len++;
  encode_msg->spare_half_octet.spare  = 0x00;
  encode_msg->sec_header_type.sec_hdr = 0x00;
  len++;
  encode_msg->message_type.msg_type = DLNASTRANSPORT;
  len++;
  encode_msg->payload_container.iei = PAYLOAD_CONTAINER;
  // encode_msg->payload_container_type.iei      = PAYLOAD_CONTAINER_TYPE;
  encode_msg->payload_container_type.iei      = 0;
  encode_msg->payload_container_type.type_val = N1_SM_INFO;
  len++;
  encode_msg->pdu_session_identity.iei = 0x12;
  len++;
  encode_msg->pdu_session_identity.pdu_session_id =
      smf_ctx->smf_proc_data.pdu_session_id;
  len++;

  // NAS SmfMsg
  smf_msg->header.extended_protocol_discriminator =
      M5G_SESSION_MANAGEMENT_MESSAGES;
  smf_msg->header.pdu_session_id = smf_ctx->smf_proc_data.pdu_session_id;
  smf_msg->header.message_type   = PDU_SESSION_RELEASE_COMMAND;
  smf_msg->header.procedure_transaction_id = smf_ctx->smf_proc_data.pti;
  smf_msg->msg.pdu_session_release_command.extended_protocol_discriminator
      .extended_proto_discriminator = M5G_SESSION_MANAGEMENT_MESSAGES;
  container_len++;
  smf_msg->msg.pdu_session_release_command.pdu_session_identity.pdu_session_id =
      smf_ctx->smf_proc_data.pdu_session_id;
  container_len++;
  smf_msg->msg.pdu_session_release_command.pti.pti = smf_ctx->smf_proc_data.pti;
  container_len++;
  smf_msg->msg.pdu_session_release_command.message_type.msg_type =
      PDU_SESSION_RELEASE_COMMAND;
  container_len++;
  smf_msg->msg.pdu_session_release_command.m5gsm_cause.cause_value =
      0x24;  // Regular deactivation
  container_len++;

  encode_msg->payload_container.len = container_len;
  len += 2;  // 2 bytes for container.len
  len += container_len;

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
    header->extended_protocol_discriminator = M5GS_MOBILITY_MANAGEMENT_MESSAGE;
    header->security_header_type = SECURITY_HEADER_TYPE_NOT_PROTECTED;
  }
  buffer = bfromcstralloc(len, "\0");
  bytes  = nas5g_message_encode(
      buffer->data, &msg, len, &ue_context->amf_context._security);

  if (retransmit) {
    if (bytes > 0) {
      buffer->slen = bytes;
      amf_app_handle_nas_dl_req(amf_ue_ngap_id, buffer, rc);

    } else {
      OAILOG_WARNING(LOG_AMF_APP, "NAS encode failed \n");
      bdestroy_wrapper(&buffer);
    }
    return rc;
  }

  itti_ngap_pdusessionresource_rel_req_t* ngap_pdu_ses_release_req = nullptr;
  MessageDef* message_p                                            = nullptr;
  pdu_session_resource_release_command_transfer amf_pdu_ses_rel_transfer_req;

  message_p =
      itti_alloc_new_message(TASK_AMF_APP, NGAP_PDUSESSIONRESOURCE_REL_REQ);
  ngap_pdu_ses_release_req =
      &message_p->ittiMsg.ngap_pdusessionresource_rel_req;
  memset(
      ngap_pdu_ses_release_req, 0,
      sizeof(itti_ngap_pdusessionresource_rel_req_t));

  // start filling message in DL to NGAP
  ngap_pdu_ses_release_req->gnb_ue_ngap_id = ue_context->gnb_ue_ngap_id;
  ngap_pdu_ses_release_req->amf_ue_ngap_id = amf_ue_ngap_id;

  /* Setting the cause of release as per OAI PCAP and sending to gNB
   * As it is UE initiated PDU session release, the cause would be
   * NAS & normal release
   */
  ngap_pdu_ses_release_req->pduSessionResourceToRelReqList.no_of_items = 1;
  ngap_pdu_ses_release_req->pduSessionResourceToRelReqList.item[0]
      .Pdu_Session_ID =
      (Ngap_PDUSessionID_t) smf_ctx->smf_proc_data.pdu_session_id;
  amf_pdu_ses_rel_transfer_req.cause.cause_group.u_group.nas.cause =
      NORMAL_RELEASE;
  amf_pdu_ses_rel_transfer_req.cause.cause_group.cause_group_type = NAS_GROUP;

  // Convert amf_pdu_ses_setup_transfer_req to bstring and assign to message_p
  ngap_pdu_ses_release_req->pduSessionResourceToRelReqList.item[0]
      .PDU_Session_Resource_TO_Release_Command_Transfer =
      amf_pdu_ses_rel_transfer_req;

  // Send message to NGAP task
  if (bytes > 0) {
    //    buffer->slen = bytes;
    buffer->slen                      = bytes;
    ngap_pdu_ses_release_req->nas_msg = bstrcpy(buffer);
    bdestroy(buffer);
  } else {
    bdestroy(buffer);
    OAILOG_ERROR(LOG_AMF_APP, "NAS encode failed for PDU Release Command\n");
    return RETURNerror;
  }
  amf_send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);
  return RETURNok;
}
}  // namespace magma5g
