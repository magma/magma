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
/*****************************************************************************
  Source      amf_app_pdu_resource_setup_req_rel.cpp
  Version     0.1
  Date        2020/12/10
  Product     AMF Core
  Subsystem   Access and PDU session Management Function
  Author      Sanjay Kumar Ojha
  Description Defination of PDU session resource setup or release request
              or response. Follow 38-413(9.2.1.1) and 24-501
              This is non-NAS message
*****************************************************************************/
#ifndef _SEEN
#define DEREGISTRATION_REQUEST_SEEN

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
#include "amf_app_ue_context_and_proc.h"
#include "ngap_messages_types.h"

using namespace std;

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
    /* AMBR has not been populated till now and default assigned
     * TODO rechaeck the default values 64kbps & 32mbps
     */
    *dl_pdu_ambr = (64 * 32768);
    *ul_pdu_ambr = (64 * 32768);

  } else {
    // refer 24-501 9.11.4.14
    switch (smf_context->dl_ambr_unit) {
      case 1:
        *dl_pdu_ambr = (1 * smf_context->dl_session_ambr);
        break;
      case 2:
        *dl_pdu_ambr = (4 * smf_context->dl_session_ambr);
        break;
      case 3:
        *dl_pdu_ambr = (16 * smf_context->dl_session_ambr);
        break;
      case 4:
        *dl_pdu_ambr = (64 * smf_context->dl_session_ambr);
        break;
        // many more to be coded
      default:
        *dl_pdu_ambr = (256 * smf_context->dl_session_ambr);
        break;
    }
    switch (smf_context->ul_ambr_unit) {
      case 1:
        *ul_pdu_ambr = (1 * smf_context->ul_session_ambr);
        break;
      case 2:
        *ul_pdu_ambr = (4 * smf_context->ul_session_ambr);
        break;
      case 3:
        *ul_pdu_ambr = (16 * smf_context->ul_session_ambr);
        break;
      case 4:
        *ul_pdu_ambr = (64 * smf_context->ul_session_ambr);
        break;
        // many more to be coded
      default:
        *ul_pdu_ambr = (256 * smf_context->ul_session_ambr);
        break;
    }
  }
}

/*
 * As per Baicel capture and message flow, the function to be called
 * before sending gRPC message to SMF with avilable information in
 * smf_context or default.
 * No NAS message is involved and direct itti_message is sent to NGAP.
 * The transfer message is bstring and heavily nested structure is
 * converted to bstring before sending to NGAP
 */
int pdu_session_resource_setup_request(
    ue_m5gmm_context_s* ue_context, amf_ue_ngap_id_t amf_ue_ngap_id) {
  if (!ue_context) {
    // TODO ue_context = GLOBAL VARIABL coded by Ashish;
  }
  int rc = RETURNerror;
  smf_context_t* smf_context;
  amf_context_t* amf_context;
  amf_context = &ue_context->amf_context;
  smf_context = &ue_context->amf_context.smf_context;
  pdu_session_resource_setup_request_transfer_t
      amf_pdu_ses_setup_transfer_req;  // pdu_res_set_change
  // pdu_session_resource_setup_req_t amf_pdu_ses_setup_req;
  itti_ngap_pdusession_resource_setup_req_t* ngap_pdu_ses_setup_req = nullptr;
  MessageDef* message_p                                             = nullptr;
  uint64_t dl_pdu_ambr;
  uint64_t ul_pdu_ambr;

  OAILOG_INFO(
      LOG_AMF_APP,
      "PDU session resource setup request message construction to NGAP\n");

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

  // ngap_pdu_ses_setup_req->ue_aggregate_maximum_bit_rate_present = true;
  ngap_pdu_ses_setup_req->ue_aggregate_maximum_bit_rate.dl = dl_pdu_ambr;
  ngap_pdu_ses_setup_req->ue_aggregate_maximum_bit_rate.ul = ul_pdu_ambr;

  // Hardcored number of pdu sessions as 1
  ngap_pdu_ses_setup_req->pduSessionResource_setup_list.no_of_items = 1;
  ngap_pdu_ses_setup_req->pduSessionResource_setup_list.item[0].Pdu_Session_ID =
      (Ngap_PDUSessionID_t)
          smf_context->smf_proc_data.pdu_session_identity.pdu_session_id;

  /* preparing for  bstring PDU_Session_Resource_Setup_Transfer
   * amf_pdu_ses_setup_transfer_req is the structure to be filed
   * and make a bundle and converted to bstring
   */
#if 1
  amf_pdu_ses_setup_transfer_req.pdu_aggregate_max_bit_rate.dl = dl_pdu_ambr;
  amf_pdu_ses_setup_transfer_req.pdu_aggregate_max_bit_rate.ul = ul_pdu_ambr;
  // UPF tied 4 octet and respective ip address are from
  // SMF date that written to context
  // char buf[] = {0x00, 0x00, 0x03, 0x3a};
  memcpy(
      &amf_pdu_ses_setup_transfer_req.up_transport_layer_info.gtp_tnl.gtp_tied,
      smf_context->gtp_tunnel_id.upf_gtp_teid, 4);
  amf_pdu_ses_setup_transfer_req.up_transport_layer_info.gtp_tnl
      .endpoint_ip_address =
      blk2bstr(&smf_context->gtp_tunnel_id.upf_gtp_teid_ip_addr, 4);
  //.endpoint_ip_address = smf_context->smf_proc_data.pdn_addr;

  amf_pdu_ses_setup_transfer_req.pdu_ip_type.pdn_type = IPv4;
  amf_pdu_ses_setup_transfer_req.qos_flow_setup_request_list.qos_flow_req_item
      .qos_flow_identifier =
      smf_context->pdu_resource_setup_req
          .pdu_session_resource_setup_request_transfer
          .qos_flow_setup_request_list.qos_flow_req_item.qos_flow_identifier;
  amf_pdu_ses_setup_transfer_req.qos_flow_setup_request_list.qos_flow_req_item
      .qos_flow_level_qos_param.qos_characteristic.non_dynamic_5QI_desc.fiveQI =
      smf_context->pdu_resource_setup_req
          .pdu_session_resource_setup_request_transfer
          .qos_flow_setup_request_list.qos_flow_req_item
          .qos_flow_level_qos_param.qos_characteristic.non_dynamic_5QI_desc
          .fiveQI;

  amf_pdu_ses_setup_transfer_req.qos_flow_setup_request_list.qos_flow_req_item
      .qos_flow_level_qos_param.alloc_reten_priority.priority_level =
      smf_context->pdu_resource_setup_req
          .pdu_session_resource_setup_request_transfer
          .qos_flow_setup_request_list.qos_flow_req_item
          .qos_flow_level_qos_param.alloc_reten_priority.priority_level;
  amf_pdu_ses_setup_transfer_req.qos_flow_setup_request_list.qos_flow_req_item
      .qos_flow_level_qos_param.alloc_reten_priority.pre_emption_cap =
      // MAY_TRIGGER_PRE_EMPTION;
      smf_context->pdu_resource_setup_req
          .pdu_session_resource_setup_request_transfer
          .qos_flow_setup_request_list.qos_flow_req_item
          .qos_flow_level_qos_param.alloc_reten_priority.pre_emption_cap;
  amf_pdu_ses_setup_transfer_req.qos_flow_setup_request_list.qos_flow_req_item
      .qos_flow_level_qos_param.alloc_reten_priority.pre_emption_vul =
      // PRE_EMPTABLE;
      smf_context->pdu_resource_setup_req
          .pdu_session_resource_setup_request_transfer
          .qos_flow_setup_request_list.qos_flow_req_item
          .qos_flow_level_qos_param.alloc_reten_priority.pre_emption_vul;
  // Adding respective header to amf_pdu_ses_setup_transfer_request
#if 0
  amf_pdu_ses_setup_transfer_req.pdu_aggregate_max_bit_rate.hdr.protocol_id = (uint16_t) 130;
  amf_pdu_ses_setup_transfer_req.pdu_aggregate_max_bit_rate.hdr.criticality = 0x0;
  amf_pdu_ses_setup_transfer_req.pdu_aggregate_max_bit_rate.hdr.pad = 0x0a;

  amf_pdu_ses_setup_transfer_req.up_transport_layer_info.hdr.protocol_id = (uint16_t) 139;
  amf_pdu_ses_setup_transfer_req.up_transport_layer_info.hdr.criticality = 0x0;
  amf_pdu_ses_setup_transfer_req.up_transport_layer_info.hdr.pad = 0x0a;

  amf_pdu_ses_setup_transfer_req.pdu_ip_type.hdr.protocol_id = (uint16_t) 134;
  amf_pdu_ses_setup_transfer_req.pdu_ip_type.hdr.criticality = 0x0;
  amf_pdu_ses_setup_transfer_req.pdu_ip_type.hdr.pad = 0x01;

  amf_pdu_ses_setup_transfer_req.qos_flow_setup_request_list.hdr.protocol_id = (uint16_t) 136;
  amf_pdu_ses_setup_transfer_req.qos_flow_setup_request_list.hdr.criticality = 0x0;
  amf_pdu_ses_setup_transfer_req.qos_flow_setup_request_list.hdr.pad = 0x07;
#endif
  // Convert amf_pdu_ses_setup_transfer_req to bstring and assign to message_p

  OAILOG_INFO(
      LOG_AMF_APP,
      "Converting pdu_session_resource_setup_request_transfer_t to bstring \n");
  ngap_pdu_ses_setup_req->pduSessionResource_setup_list.item[0]
      .PDU_Session_Resource_Setup_Request_Transfer =
      amf_pdu_ses_setup_transfer_req;
#endif  // pdu_res_set_chang
  // Send message to NGAP task
  send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);

  return RETURNok;
}

/* Resourece release request to gNB through NGAP */
int pdu_session_resource_release_request(
    ue_m5gmm_context_s* ue_context, amf_ue_ngap_id_t amf_ue_ngap_id) {
  if (!ue_context) {
    // TODO ue_context = GLOBAL VARIABL coded by Ashish;
  }
  itti_ngap_pdusessionresource_rel_req_t* ngap_pdu_ses_release_req = nullptr;
  MessageDef* message_p                                            = nullptr;
  pdu_session_resource_release_command_transfer amf_pdu_ses_rel_transfer_req;

  OAILOG_INFO(
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

  /* Setting the cause of release as per OAI PCAP and sending to gNB
   * As it is UE initiated PDU session release, the cause would be
   * NAS & normal release
   */
  ngap_pdu_ses_release_req->pduSessionResourceToRelReqList.no_of_items = 1;
  ngap_pdu_ses_release_req->pduSessionResourceToRelReqList.item[0]
      .Pdu_Session_ID =
      (Ngap_PDUSessionID_t) ue_context->amf_context.smf_context.smf_proc_data
          .pdu_session_identity.pdu_session_id;
  amf_pdu_ses_rel_transfer_req.cause.cause_group.u_group.nas.cause =
      NORMAL_RELEASE;
  amf_pdu_ses_rel_transfer_req.cause.cause_group.cause_group_type = NAS_GROUP;

  // Convert amf_pdu_ses_setup_transfer_req to bstring and assign to message_p

  OAILOG_INFO(
      LOG_AMF_APP,
      //"Converting pdu_session_resource_release_command_transfer to
      // bstring\n");
      "filling pdu_session_resource_release_command_transfer\n");
  // ngap_pdu_ses_release_req->pduSessionResourceToRelReqList.item[0]
  //    .PDU_Session_Resource_TO_Release_Command_Transfer = blk2bstr(
  //    &amf_pdu_ses_rel_transfer_req,
  //    sizeof(
  //        pdu_session_resource_release_command_transfer));  //
  //        pdu_res_set_change
  ngap_pdu_ses_release_req->pduSessionResourceToRelReqList.item[0]
      .PDU_Session_Resource_TO_Release_Command_Transfer =
      amf_pdu_ses_rel_transfer_req;
  // Send message to NGAP task
  send_msg_to_task(&amf_app_task_zmq_ctx, TASK_NGAP, message_p);

  return RETURNok;
}

}  // end  namespace magma5g
#endif
