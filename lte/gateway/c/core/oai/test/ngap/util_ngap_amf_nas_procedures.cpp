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

#include <iostream>
#include "lte/gateway/c/core/oai/test/ngap/util_ngap_pkt.hpp"

bool generator_ngap_pdusession_resource_setup_req(bstring& stream) {
  uint8_t* buffer_p = NULL;
  uint32_t length = 0;
  Ngap_NGAP_PDU_t pdu;
  Ngap_PDUSessionResourceSetupRequest_t* out = NULL;
  Ngap_PDUSessionResourceSetupRequestIEs_t* ie = NULL;
  Ngap_PDUSessionResourceSetupRequestTransferIEs_t* tx_ie = NULL;

  memset(&pdu, 0, sizeof(pdu));

  pdu.choice.initiatingMessage.procedureCode =
      Ngap_ProcedureCode_id_PDUSessionResourceSetup;
  pdu.choice.initiatingMessage.criticality = Ngap_Criticality_reject;
  pdu.present = Ngap_NGAP_PDU_PR_initiatingMessage;
  pdu.choice.initiatingMessage.value.present =
      Ngap_InitiatingMessage__value_PR_PDUSessionResourceSetupRequest;
  out =
      &pdu.choice.initiatingMessage.value.choice.PDUSessionResourceSetupRequest;

  /*
   * Setting UE information with the ones found in ue_ref
   */
  ie = (Ngap_PDUSessionResourceSetupRequestIEs_t*)calloc(
      1, sizeof(Ngap_PDUSessionResourceSetupRequestIEs_t));
  ie->id = Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID;
  ie->criticality = Ngap_Criticality_reject;
  ie->value.present =
      Ngap_PDUSessionResourceSetupRequestIEs__value_PR_AMF_UE_NGAP_ID;
  asn_uint642INTEGER(&ie->value.choice.AMF_UE_NGAP_ID, 256);
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* mandatory */
  ie = (Ngap_PDUSessionResourceSetupRequestIEs_t*)calloc(
      1, sizeof(Ngap_PDUSessionResourceSetupRequestIEs_t));
  ie->id = Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID;
  ie->criticality = Ngap_Criticality_reject;
  ie->value.present =
      Ngap_PDUSessionResourceSetupRequestIEs__value_PR_RAN_UE_NGAP_ID;
  ie->value.choice.RAN_UE_NGAP_ID = 1;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* mandatory */
  ie = (Ngap_PDUSessionResourceSetupRequestIEs_t*)calloc(
      1, sizeof(Ngap_PDUSessionResourceSetupRequestIEs_t));
  ie->id = Ngap_ProtocolIE_ID_id_PDUSessionResourceSetupListSUReq;
  ie->criticality = Ngap_Criticality_reject;
  ie->value.present =
      Ngap_PDUSessionResourceSetupRequestIEs__value_PR_PDUSessionResourceSetupListSUReq;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  // Resource Setup Request
  Ngap_PDUSessionResourceSetupItemSUReq_t* ngap_pdusession_setup_item_ies =
      (Ngap_PDUSessionResourceSetupItemSUReq_t*)calloc(
          1, sizeof(Ngap_PDUSessionResourceSetupItemSUReq_t));

  ngap_pdusession_setup_item_ies->pDUSessionID = 100;

  // NSSAI
  ngap_pdusession_setup_item_ies->s_NSSAI.sST.size = 1;
  ngap_pdusession_setup_item_ies->s_NSSAI.sST.buf =
      (uint8_t*)calloc(1, sizeof(uint8_t));
  ngap_pdusession_setup_item_ies->s_NSSAI.sST.buf[0] = 0x11;

  // Filling PDU TX Structure
  Ngap_PDUSessionResourceSetupRequestTransfer_t*
      pduSessionResourceSetupRequestTransferIEs =
          (Ngap_PDUSessionResourceSetupRequestTransfer_t*)calloc(
              1, sizeof(Ngap_PDUSessionResourceSetupRequestTransfer_t));

  tx_ie = (Ngap_PDUSessionResourceSetupRequestTransferIEs_t*)calloc(
      1, sizeof(Ngap_PDUSessionResourceSetupRequestTransferIEs_t));
  tx_ie->id = Ngap_ProtocolIE_ID_id_UL_NGU_UP_TNLInformation;
  tx_ie->criticality = Ngap_Criticality_reject;
  tx_ie->value.present =
      Ngap_PDUSessionResourceSetupRequestTransferIEs__value_PR_UPTransportLayerInformation;
  tx_ie->value.choice.UPTransportLayerInformation.present =
      Ngap_UPTransportLayerInformation_PR_gTPTunnel;

  /*transportLayerAddress*/
  tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel
      .transportLayerAddress.size = sizeof(uint32_t);
  tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel
      .transportLayerAddress.bits_unused = 0;
  tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel
      .transportLayerAddress.buf =
      (uint8_t*)calloc(1, tx_ie->value.choice.UPTransportLayerInformation.choice
                              .gTPTunnel.transportLayerAddress.size);
  tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel
      .transportLayerAddress.buf[0] = 0xc0;
  tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel
      .transportLayerAddress.buf[1] = 0xa8;
  tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel
      .transportLayerAddress.buf[2] = 0x3c;
  tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel
      .transportLayerAddress.buf[3] = 0x9b;

  tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel.gTP_TEID
      .size = sizeof(uint32_t);
  tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel.gTP_TEID
      .buf = (uint8_t*)calloc(1, tx_ie->value.choice.UPTransportLayerInformation
                                     .choice.gTPTunnel.gTP_TEID.size);
  tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel.gTP_TEID
      .buf[0] = 0x0;
  tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel.gTP_TEID
      .buf[1] = 0x0;
  tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel.gTP_TEID
      .buf[2] = 0xa;
  tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel.gTP_TEID
      .buf[3] = 0x0;

  int ret = ASN_SEQUENCE_ADD(
      &pduSessionResourceSetupRequestTransferIEs->protocolIEs.list, tx_ie);

  /*PDUSessionType*/
  tx_ie = (Ngap_PDUSessionResourceSetupRequestTransferIEs_t*)calloc(
      1, sizeof(Ngap_PDUSessionResourceSetupRequestTransferIEs_t));
  tx_ie->id = Ngap_ProtocolIE_ID_id_PDUSessionType;
  tx_ie->criticality = Ngap_Criticality_reject;
  tx_ie->value.present =
      Ngap_PDUSessionResourceSetupRequestTransferIEs__value_PR_PDUSessionType;
  tx_ie->value.choice.PDUSessionType = Ngap_PDUSessionType_ipv4;
  ret = ASN_SEQUENCE_ADD(
      &pduSessionResourceSetupRequestTransferIEs->protocolIEs.list, tx_ie);
  assert(ret == 0);

  /*Qos*/
  tx_ie = (Ngap_PDUSessionResourceSetupRequestTransferIEs_t*)calloc(
      1, sizeof(Ngap_PDUSessionResourceSetupRequestTransferIEs_t));
  tx_ie->id = Ngap_ProtocolIE_ID_id_QosFlowSetupRequestList;
  tx_ie->criticality = Ngap_Criticality_reject;
  tx_ie->value.present =
      Ngap_PDUSessionResourceSetupRequestTransferIEs__value_PR_QosFlowSetupRequestList;

  Ngap_QosFlowSetupRequestItem_t* qos_item =
      (Ngap_QosFlowSetupRequestItem_t*)calloc(
          1, sizeof(Ngap_QosFlowSetupRequestItem_t));
  qos_item->qosFlowIdentifier = 1;
  qos_item->qosFlowLevelQosParameters.qosCharacteristics.present =
      Ngap_QosCharacteristics_PR_nonDynamic5QI;
  qos_item->qosFlowLevelQosParameters.qosCharacteristics.choice.nonDynamic5QI
      .fiveQI = 9;
  qos_item->qosFlowLevelQosParameters.allocationAndRetentionPriority
      .priorityLevelARP = 8;
  qos_item->qosFlowLevelQosParameters.allocationAndRetentionPriority
      .pre_emptionCapability = 0;
  qos_item->qosFlowLevelQosParameters.allocationAndRetentionPriority
      .pre_emptionVulnerability = 0;
  asn_set_empty(&tx_ie->value.choice.QosFlowSetupRequestList.list);
  ASN_SEQUENCE_ADD(&tx_ie->value.choice.QosFlowSetupRequestList.list, qos_item);

  ret = ASN_SEQUENCE_ADD(
      &pduSessionResourceSetupRequestTransferIEs->protocolIEs.list, tx_ie);
  assert(ret == 0);

  uint8_t buffer[1024];
  asn_enc_rval_t er = aper_encode_to_buffer(
      &asn_DEF_Ngap_PDUSessionResourceSetupRequestTransfer, NULL,
      pduSessionResourceSetupRequestTransferIEs, buffer, 1024);
  if (er.encoded < 0) {
    assert(0);
    return false;
  }

  bstring transfer = blk2bstr(buffer, er.encoded);
  ngap_pdusession_setup_item_ies->pDUSessionResourceSetupRequestTransfer.size =
      blength(transfer);
  ngap_pdusession_setup_item_ies->pDUSessionResourceSetupRequestTransfer.buf =
      (uint8_t*)calloc(er.encoded, sizeof(uint8_t));

  memcpy((void*)ngap_pdusession_setup_item_ies
             ->pDUSessionResourceSetupRequestTransfer.buf,
         (void*)transfer->data, er.encoded);

  bdestroy(transfer);

  ASN_SEQUENCE_ADD(&ie->value.choice.PDUSessionResourceSetupListSUReq.list,
                   ngap_pdusession_setup_item_ies);

  if (ngap_amf_encode_pdu(&pdu, &buffer_p, &length) < 0) {
    assert(0);
    return false;
  }

  ASN_STRUCT_FREE_CONTENTS_ONLY(
      asn_DEF_Ngap_PDUSessionResourceSetupRequestTransfer,
      pduSessionResourceSetupRequestTransferIEs);
  stream = blk2bstr(buffer_p, length);
  free(buffer_p);
  free(pduSessionResourceSetupRequestTransferIEs);

  return true;
}

bool generator_itti_ngap_pdusession_resource_setup_req(bstring& stream) {
  itti_ngap_pdusession_resource_setup_req_t resource_setup_req;
  m5g_ue_description_t ue_ref;
  int ret = RETURNok;

  memset(&resource_setup_req, 0,
         sizeof(itti_ngap_pdusession_resource_setup_req_t));
  memset(&ue_ref, 0, sizeof(m5g_ue_description_t));

  resource_setup_req.gnb_ue_ngap_id = 10;
  resource_setup_req.amf_ue_ngap_id = 100;

  // maximum_bit_rate
  resource_setup_req.ue_aggregate_maximum_bit_rate.dl = 10000;
  resource_setup_req.ue_aggregate_maximum_bit_rate.ul = 10000;

  // Resource Setup List
  resource_setup_req.pduSessionResource_setup_list.no_of_items = 1;
  pdusession_setup_item_t* item =
      &(resource_setup_req.pduSessionResource_setup_list.item[0]);
  item->Pdu_Session_ID = 99;

  pdu_session_resource_setup_request_transfer_t* transfer_req =
      &(item->PDU_Session_Resource_Setup_Request_Transfer);

  transfer_req->pdu_aggregate_max_bit_rate.dl = 1000;
  transfer_req->pdu_aggregate_max_bit_rate.ul = 2000;

  unsigned char buf[sizeof(struct in6_addr)];
  if ((inet_pton(AF_INET, "192.168.1.11", buf)) < 0) {
    return false;
  }

  transfer_req->up_transport_layer_info.gtp_tnl.endpoint_ip_address =
      blk2bstr(buf, 4);
  transfer_req->up_transport_layer_info.gtp_tnl.gtp_tied[0] = 0x0;
  transfer_req->up_transport_layer_info.gtp_tnl.gtp_tied[1] = 0x0;
  transfer_req->up_transport_layer_info.gtp_tnl.gtp_tied[2] = 0xa;
  transfer_req->up_transport_layer_info.gtp_tnl.gtp_tied[3] = 0x0;
  transfer_req->pdu_ip_type.pdn_type = IPv4;

  /* QoS */
  transfer_req->qos_flow_add_or_mod_request_list.maxNumOfQosFlows = 1;
  qos_flow_setup_request_item* qos_flow = &(
      transfer_req->qos_flow_add_or_mod_request_list.item[0].qos_flow_req_item);
  qos_flow->qos_flow_identifier = 1;
  qos_flow->qos_flow_level_qos_param.qos_characteristic.non_dynamic_5QI_desc
      .fiveQI = 9;
  qos_flow->qos_flow_level_qos_param.alloc_reten_priority.priority_level = 8;
  qos_flow->qos_flow_level_qos_param.alloc_reten_priority.pre_emption_cap =
      SHALL_NOT_TRIGGER_PRE_EMPTION;
  qos_flow->qos_flow_level_qos_param.alloc_reten_priority.pre_emption_vul =
      PRE_EMPTABLE;

  ue_ref.gnb_ue_ngap_id = 1001;
  ue_ref.amf_ue_ngap_id = 2001;

  ret = ngap_amf_nas_pdusession_resource_setup_stream(&resource_setup_req,
                                                      &ue_ref, &stream);

  if (ret != 0) {
    return false;
  }

  bdestroy(transfer_req->up_transport_layer_info.gtp_tnl.endpoint_ip_address);

  return (true);
}

bool generator_ngap_pdusession_resource_rel_cmd_stream(bstring& stream) {
  uint8_t* buffer_p = NULL;
  uint32_t length = 0;
  Ngap_NGAP_PDU_t pdu;
  Ngap_PDUSessionResourceReleaseCommand_t* out = NULL;
  Ngap_PDUSessionResourceReleaseCommandIEs_t* ie = NULL;
  int hexbuf[] = {0x7e, 0x03, 0x00, 0x00, 0x00, 0x00, 0x03, 0x7e, 0x00, 0x68,
                  0x01, 0x00, 0x05, 0x2e, 0x01, 0x01, 0xd3, 0x24, 0x12, 0x01};

  memset(&pdu, 0, sizeof(pdu));
  pdu.present = Ngap_NGAP_PDU_PR_initiatingMessage;
  pdu.choice.initiatingMessage.procedureCode =
      Ngap_ProcedureCode_id_PDUSessionResourceRelease;
  pdu.choice.initiatingMessage.criticality = Ngap_Criticality_ignore;
  pdu.choice.initiatingMessage.value.present =
      Ngap_InitiatingMessage__value_PR_PDUSessionResourceReleaseCommand;
  out = &pdu.choice.initiatingMessage.value.choice
             .PDUSessionResourceReleaseCommand;

  /* mandatory */
  ie = (Ngap_PDUSessionResourceReleaseCommandIEs_t*)calloc(
      1, sizeof(Ngap_PDUSessionResourceReleaseCommandIEs_t));
  ie->id = Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID;
  ie->criticality = Ngap_Criticality_reject;
  ie->value.present =
      Ngap_PDUSessionResourceReleaseCommandIEs__value_PR_AMF_UE_NGAP_ID;
  asn_uint642INTEGER(&ie->value.choice.AMF_UE_NGAP_ID, 1);
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* mandatory */
  ie = (Ngap_PDUSessionResourceReleaseCommandIEs_t*)calloc(
      1, sizeof(Ngap_PDUSessionResourceReleaseCommandIEs_t));
  ie->id = Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID;
  ie->criticality = Ngap_Criticality_reject;
  ie->value.present =
      Ngap_PDUSessionResourceReleaseCommandIEs__value_PR_RAN_UE_NGAP_ID;
  ie->value.choice.RAN_UE_NGAP_ID = 1;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* optional NAS pdu */
  ie = (Ngap_PDUSessionResourceReleaseCommandIEs_t*)calloc(
      1, sizeof(Ngap_PDUSessionResourceReleaseCommandIEs_t));
  ie->id = Ngap_ProtocolIE_ID_id_NAS_PDU;
  ie->criticality = Ngap_Criticality_reject;
  ie->value.present =
      Ngap_PDUSessionResourceReleaseCommandIEs__value_PR_NAS_PDU;

  ie->value.choice.NAS_PDU.size = 20;
  ie->value.choice.NAS_PDU.buf =
      (uint8_t*)calloc(1, ie->value.choice.NAS_PDU.size * sizeof(uint8_t));
  memset(ie->value.choice.NAS_PDU.buf, 0,
         sizeof(ie->value.choice.NAS_PDU.size * sizeof(uint8_t)));

  for (uint32_t i = 0; i < ie->value.choice.NAS_PDU.size; i++) {
    ie->value.choice.NAS_PDU.buf[i] = hexbuf[i];
  }
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  ie = (Ngap_PDUSessionResourceReleaseCommandIEs_t*)calloc(
      1, sizeof(Ngap_PDUSessionResourceReleaseCommandIEs_t));
  ie->id = Ngap_ProtocolIE_ID_id_PDUSessionResourceToReleaseListRelCmd;
  ie->criticality = Ngap_Criticality_reject;
  ie->value.present =
      Ngap_PDUSessionResourceReleaseCommandIEs__value_PR_PDUSessionResourceToReleaseListRelCmd;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  Ngap_PDUSessionResourceToReleaseItemRelCmd_t* RelItem =
      (Ngap_PDUSessionResourceToReleaseItemRelCmd_t*)calloc(
          1, sizeof(Ngap_PDUSessionResourceToReleaseItemRelCmd_t));
  RelItem->pDUSessionID = 1;

  Ngap_PDUSessionResourceReleaseCommandTransfer_t*
      PDUSessionResourceReleaseCommandTransferIEs =
          (Ngap_PDUSessionResourceReleaseCommandTransfer_t*)calloc(
              1, sizeof(Ngap_PDUSessionResourceReleaseCommandTransfer_t));

  PDUSessionResourceReleaseCommandTransferIEs->cause.present =
      Ngap_Cause_PR_nas;

  PDUSessionResourceReleaseCommandTransferIEs->cause.choice.nas =
      Ngap_CauseNas_normal_release;

  uint8_t* buffer = NULL;

  ssize_t encoded_size = aper_encode_to_new_buffer(
      &asn_DEF_Ngap_PDUSessionResourceReleaseCommandTransfer, NULL,
      PDUSessionResourceReleaseCommandTransferIEs, (void**)&buffer);

  RelItem->pDUSessionResourceReleaseCommandTransfer.size = encoded_size;
  RelItem->pDUSessionResourceReleaseCommandTransfer.buf = (uint8_t*)calloc(
      1,
      RelItem->pDUSessionResourceReleaseCommandTransfer.size * sizeof(uint8_t));
  memcpy((void*)RelItem->pDUSessionResourceReleaseCommandTransfer.buf,
         (void*)buffer, encoded_size);

  ASN_SEQUENCE_ADD(&ie->value.choice.PDUSessionResourceToReleaseListRelCmd.list,
                   RelItem);
  free(buffer);

  if (ngap_amf_encode_pdu(&pdu, &buffer_p, &length) < 0) {
    return (false);
  }

  ASN_STRUCT_FREE_CONTENTS_ONLY(
      asn_DEF_Ngap_PDUSessionResourceReleaseCommandTransfer,
      PDUSessionResourceReleaseCommandTransferIEs);
  free(PDUSessionResourceReleaseCommandTransferIEs);
  stream = blk2bstr(buffer_p, length);
  free(buffer_p);

  return (true);
}

status_code_e send_ngap_gnb_reset_ack() {
  status_code_e rc = RETURNok;
  itti_ngap_gnb_initiated_reset_ack_t gnb_reset_ack_msg = {};

  gnb_reset_ack_msg.sctp_assoc_id = 1;
  gnb_reset_ack_msg.sctp_stream_id = 1;
  gnb_reset_ack_msg.ngap_reset_type = M5G_RESET_ALL;
  gnb_reset_ack_msg.num_ue = 1;
  gnb_reset_ack_msg.ue_to_reset_list =
      reinterpret_cast<ng_sig_conn_id_t*>(calloc(2, sizeof(ng_sig_conn_id_t)));

  gnb_reset_ack_msg.ue_to_reset_list[0].amf_ue_ngap_id = 1;
  gnb_reset_ack_msg.ue_to_reset_list[0].gnb_ue_ngap_id = 1;

  rc = ngap_handle_gnb_initiated_reset_ack(&gnb_reset_ack_msg);

  free(gnb_reset_ack_msg.ue_to_reset_list);
  return rc;
}
