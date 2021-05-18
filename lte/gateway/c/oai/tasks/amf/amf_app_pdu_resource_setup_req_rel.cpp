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
#include "log.h"
#include "intertask_interface_types.h"
#include "intertask_interface.h"
#include "directoryd.h"
#include "conversions.h"
#include "bstrlib.h"
#ifdef __cplusplus
};
#endif
#include "common_defs.h"
#include "amf_app_ue_context_and_proc.h"
#include "ngap_messages_types.h"

namespace magma5g {
extern task_zmq_ctx_t amf_app_task_zmq_ctx;
/*
 * AMBR calculation based on 9.11.4.14 of 24-501
 */
void ambr_calculation_pdu_session(
    smf_context_t* smf_context, uint64_t* dl_pdu_ambr, uint64_t* ul_pdu_ambr) {
  if ((smf_context->dl_ambr_unit == 0) || (smf_context->ul_ambr_unit == 0) ||
      (smf_context->dl_session_ambr == 0) ||
      (smf_context->dl_session_ambr == 0)) {
    // AMBR has not been populated till now and default assigned
    *dl_pdu_ambr = (64 * 32768);
    *ul_pdu_ambr = (64 * 32768);
  } else {
    // refer 24-501 9.11.4.14
    if ((smf_context->dl_ambr_unit) < 4) {
      *dl_pdu_ambr =
          4 ^ (smf_context->dl_ambr_unit - 1) * (smf_context->dl_session_ambr);
    } else {
      *dl_pdu_ambr = 256 * (smf_context->dl_session_ambr);
    }

    if ((smf_context->ul_ambr_unit) < 4) {
      *ul_pdu_ambr =
          4 ^ (smf_context->ul_ambr_unit - 1) * (smf_context->ul_session_ambr);
    } else {
      *ul_pdu_ambr = 256 * (smf_context->ul_session_ambr);
    }
  }
}

/*
 * the function to be called before sending gRPC message to SMF with available
 * information in smf_context or default. No NAS message is involved and direct
 * itti_message is sent to NGAP.
 */
int pdu_session_resource_setup_request(
    ue_m5gmm_context_s* ue_context, amf_ue_ngap_id_t amf_ue_ngap_id,
    smf_context_t* smf_context) {
  pdu_session_resource_setup_request_transfer_t amf_pdu_ses_setup_transfer_req;
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
  ambr_calculation_pdu_session(smf_context, &dl_pdu_ambr, &ul_pdu_ambr);
  ngap_pdu_ses_setup_req->ue_aggregate_maximum_bit_rate.dl = dl_pdu_ambr;
  ngap_pdu_ses_setup_req->ue_aggregate_maximum_bit_rate.ul = ul_pdu_ambr;

  // Hardcoded number of pdu sessions as 1
  ngap_pdu_ses_setup_req->pduSessionResource_setup_list.no_of_items = 1;
  ngap_pdu_ses_setup_req->pduSessionResource_setup_list.item[0].Pdu_Session_ID =
      (Ngap_PDUSessionID_t)
          smf_context->smf_proc_data.pdu_session_identity.pdu_session_id;

  /* preparing for PDU_Session_Resource_Setup_Transfer.
   * amf_pdu_ses_setup_transfer_req is the structure to be filled.
   */
  amf_pdu_ses_setup_transfer_req.pdu_aggregate_max_bit_rate.dl = dl_pdu_ambr;
  amf_pdu_ses_setup_transfer_req.pdu_aggregate_max_bit_rate.ul = ul_pdu_ambr;

  // UPF teid 4 octet and respective ip address are from SMF context
  memcpy(
      &amf_pdu_ses_setup_transfer_req.up_transport_layer_info.gtp_tnl.gtp_tied,
      smf_context->gtp_tunnel_id.upf_gtp_teid, GNB_TEID_LEN);
  amf_pdu_ses_setup_transfer_req.up_transport_layer_info.gtp_tnl
      .endpoint_ip_address = blk2bstr(
      &smf_context->gtp_tunnel_id.upf_gtp_teid_ip_addr, GNB_IPV4_ADDR_LEN);
  amf_pdu_ses_setup_transfer_req.pdu_ip_type.pdn_type = IPv4;

  memcpy(
      &amf_pdu_ses_setup_transfer_req.qos_flow_setup_request_list
           .qos_flow_req_item,
      &smf_context->pdu_resource_setup_req
           .pdu_session_resource_setup_request_transfer
           .qos_flow_setup_request_list.qos_flow_req_item,
      sizeof(qos_flow_setup_request_item));
  // Adding respective header to amf_pdu_ses_setup_transfer_request
  ngap_pdu_ses_setup_req->pduSessionResource_setup_list.item[0]
      .PDU_Session_Resource_Setup_Request_Transfer =
      amf_pdu_ses_setup_transfer_req;

  // Send message to NGAP task
  send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);

  return RETURNok;
}

/* Resource release request to gNB through NGAP */
int pdu_session_resource_release_request(
    ue_m5gmm_context_s* ue_context, amf_ue_ngap_id_t amf_ue_ngap_id,
    smf_context_t* smf_ctx) {
  itti_ngap_pdusessionresource_rel_req_t* ngap_pdu_ses_release_req = nullptr;
  MessageDef* message_p                                            = nullptr;
  pdu_session_resource_release_command_transfer amf_pdu_ses_rel_transfer_req;

  OAILOG_DEBUG(
      LOG_AMF_APP,
      "PDU session resource release request message construction to NGAP\n");

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

  /* Setting the cause of release and sending to gNB
   * As it is UE initiated PDU session release, the cause would be
   * NAS & normal release
   */
  ngap_pdu_ses_release_req->pduSessionResourceToRelReqList.no_of_items = 1;
  ngap_pdu_ses_release_req->pduSessionResourceToRelReqList.item[0]
      .Pdu_Session_ID =
      (Ngap_PDUSessionID_t)
          smf_ctx->smf_proc_data.pdu_session_identity.pdu_session_id;
  amf_pdu_ses_rel_transfer_req.cause.cause_group.u_group.nas.cause =
      NORMAL_RELEASE;
  amf_pdu_ses_rel_transfer_req.cause.cause_group.cause_group_type = NAS_GROUP;

  OAILOG_DEBUG(
      LOG_AMF_APP, "filling pdu_session_resource_release_command_transfer\n");
  ngap_pdu_ses_release_req->pduSessionResourceToRelReqList.item[0]
      .PDU_Session_Resource_TO_Release_Command_Transfer =
      amf_pdu_ses_rel_transfer_req;

  // Send message to NGAP task
  send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);

  return RETURNok;
}
}  // namespace magma5g
