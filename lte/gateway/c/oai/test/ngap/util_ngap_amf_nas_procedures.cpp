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
#include "util_ngap_pkt.h"

bool generator_ngap_pdusession_resource_setup_req(bstring &stream) {
    uint8_t* buffer_p                                       = NULL;
    uint32_t length                                         = 0;
    Ngap_NGAP_PDU_t pdu;
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

    /*
     * Setting UE information with the ones found in ue_ref
     */
    ie = (Ngap_PDUSessionResourceSetupRequestIEs_t*) calloc(
        1, sizeof(Ngap_PDUSessionResourceSetupRequestIEs_t));
    ie->id          = Ngap_ProtocolIE_ID_id_AMF_UE_NGAP_ID;
    ie->criticality = Ngap_Criticality_reject;
    ie->value.present =
        Ngap_PDUSessionResourceSetupRequestIEs__value_PR_AMF_UE_NGAP_ID;
    ie->value.choice.AMF_UE_NGAP_ID = 2;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

    /* mandatory */
    ie = (Ngap_PDUSessionResourceSetupRequestIEs_t*) calloc(
        1, sizeof(Ngap_PDUSessionResourceSetupRequestIEs_t));
    ie->id          = Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID;
    ie->criticality = Ngap_Criticality_reject;
    ie->value.present =
        Ngap_PDUSessionResourceSetupRequestIEs__value_PR_RAN_UE_NGAP_ID;
    ie->value.choice.RAN_UE_NGAP_ID = 1;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

    /* mandatory */
    ie = (Ngap_PDUSessionResourceSetupRequestIEs_t*) calloc(
        1, sizeof(Ngap_PDUSessionResourceSetupRequestIEs_t));
    ie->id          = Ngap_ProtocolIE_ID_id_PDUSessionResourceSetupListSUReq;
    ie->criticality = Ngap_Criticality_reject;
    ie->value.present =
        Ngap_PDUSessionResourceSetupRequestIEs__value_PR_PDUSessionResourceSetupListSUReq;
    ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

    // Resource Setup Request
    Ngap_PDUSessionResourceSetupItemSUReq_t* ngap_pdusession_setup_item_ies =
           (Ngap_PDUSessionResourceSetupItemSUReq_t *)calloc(1, sizeof(Ngap_PDUSessionResourceSetupItemSUReq_t));

    ngap_pdusession_setup_item_ies->pDUSessionID = 100;

    //NSSAI
    ngap_pdusession_setup_item_ies->s_NSSAI.sST.size = 1;
    ngap_pdusession_setup_item_ies->s_NSSAI.sST.buf = (uint8_t*) calloc(1, sizeof(uint8_t));
    ngap_pdusession_setup_item_ies->s_NSSAI.sST.buf[0] = 0x11;

    //Filling PDU TX Structure
    Ngap_PDUSessionResourceSetupRequestTransfer_t*
      pduSessionResourceSetupRequestTransferIEs =
          (Ngap_PDUSessionResourceSetupRequestTransfer_t*) calloc(
              1, sizeof(Ngap_PDUSessionResourceSetupRequestTransfer_t));

    tx_ie = (Ngap_PDUSessionResourceSetupRequestTransferIEs_t*) calloc(
                   1, sizeof(Ngap_PDUSessionResourceSetupRequestTransferIEs_t));
    tx_ie->id          = Ngap_ProtocolIE_ID_id_UL_NGU_UP_TNLInformation;
    tx_ie->criticality = Ngap_Criticality_reject;
    tx_ie->value.present = Ngap_PDUSessionResourceSetupRequestTransferIEs__value_PR_UPTransportLayerInformation;
    tx_ie->value.choice.UPTransportLayerInformation.present = Ngap_UPTransportLayerInformation_PR_gTPTunnel;

    /*transportLayerAddress*/
    tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel.transportLayerAddress.size=sizeof(uint32_t);
    tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel.transportLayerAddress.bits_unused=0;
    tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel.transportLayerAddress.buf=
        (uint8_t*) calloc(1, tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel.transportLayerAddress.size);
    tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel.transportLayerAddress.buf[0]=0xc0;
    tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel.transportLayerAddress.buf[1]=0xa8;
    tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel.transportLayerAddress.buf[2]=0x3c;
    tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel.transportLayerAddress.buf[3]=0x9b;

    tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel.gTP_TEID.size=sizeof(uint32_t);
    tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel.gTP_TEID.buf=
         (uint8_t*) calloc(1, tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel.gTP_TEID.size);
    tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel.gTP_TEID.buf[0]=0x0;
    tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel.gTP_TEID.buf[1]=0x0;
    tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel.gTP_TEID.buf[2]=0xa;
    tx_ie->value.choice.UPTransportLayerInformation.choice.gTPTunnel.gTP_TEID.buf[3]=0x0;

    int ret = ASN_SEQUENCE_ADD(
                 &pduSessionResourceSetupRequestTransferIEs->protocolIEs.list, tx_ie);

    /*PDUSessionType*/
    tx_ie = (Ngap_PDUSessionResourceSetupRequestTransferIEs_t*) calloc(
                      1, sizeof(Ngap_PDUSessionResourceSetupRequestTransferIEs_t));
    tx_ie->id          = Ngap_ProtocolIE_ID_id_PDUSessionType;
    tx_ie->criticality = Ngap_Criticality_reject;
    tx_ie->value.present = Ngap_PDUSessionResourceSetupRequestTransferIEs__value_PR_PDUSessionType;
    tx_ie->value.choice.PDUSessionType = Ngap_PDUSessionType_ipv4;
    ret = ASN_SEQUENCE_ADD(&pduSessionResourceSetupRequestTransferIEs->protocolIEs.list, tx_ie);
    assert(ret == 0);

    /*Qos*/
    tx_ie = (Ngap_PDUSessionResourceSetupRequestTransferIEs_t*) calloc(
                  1, sizeof(Ngap_PDUSessionResourceSetupRequestTransferIEs_t));
    tx_ie->id          = Ngap_ProtocolIE_ID_id_QosFlowSetupRequestList;
    tx_ie->criticality = Ngap_Criticality_reject;
    tx_ie->value.present = Ngap_PDUSessionResourceSetupRequestTransferIEs__value_PR_QosFlowSetupRequestList;

    Ngap_QosFlowSetupRequestItem_t* qos_item = (Ngap_QosFlowSetupRequestItem_t*) calloc(1, sizeof(Ngap_QosFlowSetupRequestItem_t)); 
    qos_item->qosFlowIdentifier = 1;
    qos_item->qosFlowLevelQosParameters.qosCharacteristics.present = Ngap_QosCharacteristics_PR_nonDynamic5QI;
    qos_item->qosFlowLevelQosParameters.qosCharacteristics.choice.nonDynamic5QI.fiveQI = 9;
    qos_item->qosFlowLevelQosParameters.allocationAndRetentionPriority.priorityLevelARP = 8;
    qos_item->qosFlowLevelQosParameters.allocationAndRetentionPriority.pre_emptionCapability = 0;
    qos_item->qosFlowLevelQosParameters.allocationAndRetentionPriority.pre_emptionVulnerability = 0;
    asn_set_empty(&tx_ie->value.choice.QosFlowSetupRequestList.list);
    ASN_SEQUENCE_ADD(&tx_ie->value.choice.QosFlowSetupRequestList.list, qos_item);

    ret = ASN_SEQUENCE_ADD(&pduSessionResourceSetupRequestTransferIEs->protocolIEs.list, tx_ie);
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
    ngap_pdusession_setup_item_ies->pDUSessionResourceSetupRequestTransfer.size =  blength(transfer);
    ngap_pdusession_setup_item_ies->pDUSessionResourceSetupRequestTransfer.buf =
           (uint8_t*) calloc(er.encoded, sizeof(uint8_t));

    memcpy((void*) ngap_pdusession_setup_item_ies->pDUSessionResourceSetupRequestTransfer.buf,
           (void*) transfer->data, er.encoded);

    bdestroy(transfer);

    ASN_SEQUENCE_ADD(&ie->value.choice.PDUSessionResourceSetupListSUReq.list, ngap_pdusession_setup_item_ies);

    if (ngap_amf_encode_pdu(&pdu, &buffer_p, &length) < 0) {
	assert(0);
        return false;
    }

    ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_PDUSessionResourceSetupRequestTransfer, pduSessionResourceSetupRequestTransferIEs);
    stream = blk2bstr(buffer_p, length);
    free(buffer_p);
    free(pduSessionResourceSetupRequestTransferIEs);

    return true;
}
