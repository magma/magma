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
/****************************************************************************
  Source      ngap_amf_nas_procedures.c
  Date        2020/07/28
  Subsystem   Access and Mobility Management Function
  Description Defines NG Application Protocol Messages

*****************************************************************************/

#include <stdio.h>
#include <stdint.h>
#include <stdlib.h>
#include <string.h>

#include "bstrlib.h"
#include "dynamic_memory_check.h"
#include "assertions.h"
#include "hashtable.h"
#include "log.h"
#include "conversions.h"
#include "asn1_conversions.h"
#include "ngap_amf_encoder.h"
#include "ngap_amf.h"
#include "ngap_amf_nas_procedures.h"
#include "ngap_amf_itti_messaging.h"
#include "service303.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_38.413.h"
#include "INTEGER.h"
#include "OCTET_STRING.h"
#include "Ngap_NGAP-PDU.h"
#include "Ngap_EUTRA-CGI.h"
#include "Ngap_GBR-QosInformation.h"
#include "Ngap_GUAMI.h"
#include "Ngap_NAS-PDU.h"
#include "Ngap_PLMNIdentity.h"
#include "Ngap_ProcedureCode.h"
#include "Ngap_ProtocolIE-Field.h"
#include "Ngap_SecurityKey.h"
#include "Ngap_TAI.h"
#include "Ngap_TransportLayerAddress.h"
#include "Ngap_UEAggregateMaximumBitRate.h"
#include "Ngap_UESecurityCapabilities.h"
#include "TrackingAreaIdentity.h"
#include "asn_SEQUENCE_OF.h"
#include "ngap_state.h"
#include "Ngap_CauseMisc.h"
#include "Ngap_CauseNas.h"
#include "Ngap_CauseProtocol.h"
#include "Ngap_CauseRadioNetwork.h"
#include "Ngap_CauseTransport.h"
#include "Ngap_InitialUEMessage.h"
#include "ngap_amf_handlers.h"
#include "ngap_common.h"


//------------------------------------------------------------------------------
int ngap_generate_ngap_pdusession_resource_setup_req(
    ngap_state_t* state, itti_ngap_pdusession_resource_setup_req_t* const
                             pdusession_resource_setup_req) {
  OAILOG_FUNC_IN(LOG_NGAP);
  m5g_ue_description_t* ue_ref = NULL;
  uint8_t* buffer_p            = NULL;
  uint32_t length              = 0;
  void* id                     = NULL;
  const gnb_ue_ngap_id_t gnb_ue_ngap_id =
      pdusession_resource_setup_req->gnb_ue_ngap_id;
  const amf_ue_ngap_id_t amf_ue_ngap_id =
      pdusession_resource_setup_req->amf_ue_ngap_id;
  pdu_session_resource_setup_request_transfer_t amf_pdu_ses_setup_transfer_req;
  hashtable_ts_get(
      &state->amfid2associd, (const hash_key_t) amf_ue_ngap_id, (void**) &id);
  if (id) {
    sctp_assoc_id_t sctp_assoc_id = (sctp_assoc_id_t)(uintptr_t) id;
    gnb_description_t* gnb_ref    = ngap_state_get_gnb(state, sctp_assoc_id);
    if (gnb_ref) {
      ue_ref = ngap_state_get_ue_gnbid(gnb_ref->sctp_assoc_id, gnb_ue_ngap_id);
    }
  }
  // TODO remove soon:
  if (!ue_ref) {
    ue_ref = ngap_state_get_ue_amfid(amf_ue_ngap_id);
  }
  // finally!
  // if (!ue_ref) {
  if (ue_ref) {  // TODO tmp for testing
                 /*
                  * If the UE-associated logical NG-connection is not established,
                  * * * * the AMF shall allocate a unique AMF UE NGAP ID to be used for the
                  * UE.
                  */
    OAILOG_ERROR(
        LOG_NGAP,
        "Unknown UE AMF ID " AMF_UE_NGAP_ID_FMT
        ", This case is not handled right now\n",
        amf_ue_ngap_id);
    OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
  } else {
    /*
     * We have found the UE in the list.
     * Create new IE list message and encode it.
     */
    Ngap_NGAP_PDU_t pdu                                     = {0};
    Ngap_PDUSessionResourceSetupRequest_t* out              = NULL;
    Ngap_PDUSessionResourceSetupRequestIEs_t* ie            = NULL;
    Ngap_PDUSessionResourceSetupRequestTransferIEs_t* tx_ie = NULL;
    memset(&pdu, 0, sizeof(pdu));

    pdu.choice.initiatingMessage.procedureCode =
        Ngap_ProcedureCode_id_PDUSessionResourceSetup;
    pdu.choice.initiatingMessage.criticality = Ngap_Criticality_reject;
    pdu.present = Ngap_NGAP_PDU_PR_initiatingMessage;
    pdu.choice.initiatingMessage.value.present =
        Ngap_InitiatingMessage__value_PR_PDUSessionResourceSetupRequest;
    out = &pdu.choice.initiatingMessage.value.choice
               .PDUSessionResourceSetupRequest;
    ue_ref->ng_ue_state = NGAP_UE_CONNECTED;
    /*
     * Setting UE information with the ones found in ue_ref
     */
    ie = (Ngap_PDUSessionResourceSetupRequestIEs_t*) calloc(
        1, sizeof(Ngap_PDUSessionResourceSetupRequestIEs_t));
    ie->id          = Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID;
    ie->criticality = Ngap_Criticality_reject;
    ie->value.present =
        Ngap_PDUSessionResourceSetupRequestIEs__value_PR_AMF_UE_NGAP_ID;
    ie->value.choice.AMF_UE_NGAP_ID = ue_ref->amf_ue_ngap_id;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

    /* mandatory */
    ie = (Ngap_PDUSessionResourceSetupRequestIEs_t*) calloc(
        1, sizeof(Ngap_PDUSessionResourceSetupRequestIEs_t));
    ie->id          = Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID;
    ie->criticality = Ngap_Criticality_reject;
    ie->value.present =
        Ngap_PDUSessionResourceSetupRequestIEs__value_PR_RAN_UE_NGAP_ID;
    ie->value.choice.RAN_UE_NGAP_ID = ue_ref->gnb_ue_ngap_id;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);


    /* mandatory */
    ie = (Ngap_PDUSessionResourceSetupRequestIEs_t*) calloc(
        1, sizeof(Ngap_PDUSessionResourceSetupRequestIEs_t));
    ie->id          = Ngap_ProtocolIE_ID_id_PDUSessionResourceSetupListSUReq;
    ie->criticality = Ngap_Criticality_reject;
    ie->value.present =
        Ngap_PDUSessionResourceSetupRequestIEs__value_PR_PDUSessionResourceSetupListSUReq;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

    for (int i = 0; i < pdusession_resource_setup_req
                            ->pduSessionResource_setup_list.no_of_items;
         i++) {
      Ngap_PDUSessionResourceSetupItemSUReq_t* ngap_pdusession_setup_item_ies =
          calloc(1, sizeof(Ngap_PDUSessionResourceSetupItemSUReq_t));

      ngap_pdusession_setup_item_ies->pDUSessionID =
          pdusession_resource_setup_req->pduSessionResource_setup_list.item[i]
              .Pdu_Session_ID;


      /*NSSAI: TODO remove hardcoded value*/
      ngap_pdusession_setup_item_ies->s_NSSAI.sST.size = 1;
      ngap_pdusession_setup_item_ies->s_NSSAI.sST.buf =
          (uint8_t*) calloc(1, sizeof(uint8_t));
      ngap_pdusession_setup_item_ies->s_NSSAI.sST.buf[0] = 0x11;

      // filling PDU TX Structure
      amf_pdu_ses_setup_transfer_req =
          pdusession_resource_setup_req->pduSessionResource_setup_list.item[i]
              .PDU_Session_Resource_Setup_Request_Transfer;

      /*tx_out*/
      Ngap_PDUSessionResourceSetupRequestTransfer_t*
          pduSessionResourceSetupRequestTransferIEs =
              (Ngap_PDUSessionResourceSetupRequestTransfer_t*) calloc(
                  1, sizeof(Ngap_PDUSessionResourceSetupRequestTransfer_t));

      tx_ie = (Ngap_PDUSessionResourceSetupRequestTransferIEs_t*) calloc(
          1, sizeof(Ngap_PDUSessionResourceSetupRequestTransferIEs_t));
      tx_ie->id          = Ngap_ProtocolIE_ID_id_UL_NGU_UP_TNLInformation;
      tx_ie->criticality = Ngap_Criticality_reject;
      tx_ie->value.present =
          Ngap_PDUSessionResourceSetupRequestTransferIEs__value_PR_UPTransportLayerInformation;

      tx_ie->value.choice.UPTransportLayerInformation.present =
          Ngap_UPTransportLayerInformation_PR_gTPTunnel;

     /*transportLayerAddress*/
      tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel
          .transportLayerAddress.buf =
          calloc( blength(amf_pdu_ses_setup_transfer_req.up_transport_layer_info .gtp_tnl.endpoint_ip_address),
              sizeof(uint8_t));

      memcpy(
          tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel
              .transportLayerAddress.buf,
          amf_pdu_ses_setup_transfer_req.up_transport_layer_info.gtp_tnl
              .endpoint_ip_address->data,
          blength(amf_pdu_ses_setup_transfer_req.up_transport_layer_info.gtp_tnl
                      .endpoint_ip_address));

      tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel
          .transportLayerAddress.size =
          blength(amf_pdu_ses_setup_transfer_req.up_transport_layer_info.gtp_tnl
                      .endpoint_ip_address);

      tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel
          .transportLayerAddress.bits_unused = 0;

      /*gTP_TEID*/
      OCTET_STRING_fromBuf(
          &tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel
               .gTP_TEID,
          (char*) amf_pdu_ses_setup_transfer_req.up_transport_layer_info.gtp_tnl
              .gtp_tied,
          4);

      int ret = ASN_SEQUENCE_ADD(
          &pduSessionResourceSetupRequestTransferIEs->protocolIEs.list,
          tx_ie);
      if (ret != 0) OAILOG_ERROR(LOG_NGAP, " encode error ");

      /*PDUSessionType*/
      tx_ie = (Ngap_PDUSessionResourceSetupRequestTransferIEs_t*) calloc(
          1, sizeof(Ngap_PDUSessionResourceSetupRequestTransferIEs_t));
      tx_ie->id          = Ngap_ProtocolIE_ID_id_PDUSessionType;
      tx_ie->criticality = Ngap_Criticality_reject;
      tx_ie->value.present =
          Ngap_PDUSessionResourceSetupRequestTransferIEs__value_PR_PDUSessionType;

      tx_ie->value.choice.PDUSessionType =
          amf_pdu_ses_setup_transfer_req.pdu_ip_type.pdn_type;

      ret = ASN_SEQUENCE_ADD(
          &pduSessionResourceSetupRequestTransferIEs->protocolIEs.list,
          tx_ie);
      if (ret != 0) OAILOG_ERROR(LOG_NGAP, " encode  error ");


      /*Qos*/
      tx_ie = (Ngap_PDUSessionResourceSetupRequestTransferIEs_t*) calloc(
          1, sizeof(Ngap_PDUSessionResourceSetupRequestTransferIEs_t));
      tx_ie->id          = Ngap_ProtocolIE_ID_id_QosFlowSetupRequestList;
      tx_ie->criticality = Ngap_Criticality_reject;
      tx_ie->value.present =
          Ngap_PDUSessionResourceSetupRequestTransferIEs__value_PR_QosFlowSetupRequestList;

      for (int i = 0; i < /*no_of_qos_items*/ 1; i++) {
        Ngap_QosFlowSetupRequestItem_t* qos_item =
            (Ngap_QosFlowSetupRequestItem_t*) calloc(
                1, sizeof(Ngap_QosFlowSetupRequestItem_t));

        qos_item->qosFlowIdentifier =
            amf_pdu_ses_setup_transfer_req.qos_flow_setup_request_list
                .qos_flow_req_item.qos_flow_identifier;

        /* Ngap_QosCharacteristics */
        {
          qos_item->qosFlowLevelQosParameters.qosCharacteristics.present =
              Ngap_QosCharacteristics_PR_nonDynamic5QI;
          qos_item->qosFlowLevelQosParameters.qosCharacteristics.choice
              .nonDynamic5QI.fiveQI =
              amf_pdu_ses_setup_transfer_req.qos_flow_setup_request_list
                  .qos_flow_req_item.qos_flow_level_qos_param.qos_characteristic
                  .non_dynamic_5QI_desc.fiveQI;
        }
        /* Ngap_AllocationAndRetentionPriority */
        qos_item->qosFlowLevelQosParameters.allocationAndRetentionPriority
            .priorityLevelARP =
            amf_pdu_ses_setup_transfer_req.qos_flow_setup_request_list
                .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
                .priority_level;

        qos_item->qosFlowLevelQosParameters.allocationAndRetentionPriority
            .pre_emptionCapability =
            amf_pdu_ses_setup_transfer_req.qos_flow_setup_request_list
                .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
                .pre_emption_cap;

        qos_item->qosFlowLevelQosParameters.allocationAndRetentionPriority
            .pre_emptionVulnerability =
            amf_pdu_ses_setup_transfer_req.qos_flow_setup_request_list
                .qos_flow_req_item.qos_flow_level_qos_param.alloc_reten_priority
                .pre_emption_vul;

        asn_set_empty(&tx_ie->value.choice.QosFlowSetupRequestList.list);
        ASN_SEQUENCE_ADD(
            &tx_ie->value.choice.QosFlowSetupRequestList.list, qos_item);

      }

      ret = ASN_SEQUENCE_ADD(
          &pduSessionResourceSetupRequestTransferIEs->protocolIEs.list,
          tx_ie);
      if (ret != 0) OAILOG_ERROR(LOG_NGAP, " encode error \n");

      uint32_t buffer_size = 512;
      char* buffer         = (char*) calloc(1, buffer_size);

      asn_enc_rval_t er = aper_encode_to_buffer(
          &asn_DEF_Ngap_PDUSessionResourceSetupRequestTransfer, NULL,
          pduSessionResourceSetupRequestTransferIEs, buffer, buffer_size);

      asn_fprint(
          stderr, &asn_DEF_Ngap_PDUSessionResourceSetupRequestTransfer,
          pduSessionResourceSetupRequestTransferIEs);

      bstring transfer = blk2bstr(buffer, er.encoded);
      ngap_pdusession_setup_item_ies->pDUSessionResourceSetupRequestTransfer
          .buf = (uint8_t*) calloc(er.encoded, sizeof(uint8_t));

      memcpy(
          (void*) ngap_pdusession_setup_item_ies
              ->pDUSessionResourceSetupRequestTransfer.buf,
          (void*) transfer->data, er.encoded);

      ngap_pdusession_setup_item_ies->pDUSessionResourceSetupRequestTransfer
          .size = blength(transfer);

      ASN_SEQUENCE_ADD(
          &ie->value.choice.PDUSessionResourceSetupListSUReq.list,
          ngap_pdusession_setup_item_ies);

    } /*for loop*/

    if (ngap_amf_encode_pdu(&pdu, &buffer_p, &length) < 0) {
      // TODO: handle something
      OAILOG_ERROR(LOG_NGAP, "Encoding of ngap_PDUSessionResourceSetup failed \n");
      OAILOG_FUNC_RETURN(LOG_NGAP, RETURNerror);
    }
    OAILOG_NOTICE(
        LOG_NGAP,
        "Send NGAP PDUSessionResourceSetup message AMF_UE_NGAP_ID = " AMF_UE_NGAP_ID_FMT
        " gNB_UE_NGAP_ID = " GNB_UE_NGAP_ID_FMT "\n",
        (amf_ue_ngap_id_t) ue_ref->amf_ue_ngap_id,
        (gnb_ue_ngap_id_t) ue_ref->gnb_ue_ngap_id);
    bstring b = blk2bstr(buffer_p, length);
    free(buffer_p);
     ngap_amf_itti_send_sctp_request( &b, ue_ref->sctp_assoc_id,
     ue_ref->sctp_stream_send, ue_ref->amf_ue_ngap_id);
  }
  OAILOG_FUNC_RETURN(LOG_NGAP, RETURNok);
}

