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

//=====================================================================================================
// NGSetupFailureIEs NGAP-PROTOCOL-IES ::= {
//   { ID id-Cause                  CRITICALITY ignore   TYPE Cause PRESENCE
//   mandatory           }| { ID id-TimeToWait             CRITICALITY ignore
//   TYPE TimeToWait  PRESENCE optional            }| { ID
//   id-CriticalityDiagnostics CRITICALITY ignore   TYPE CriticalityDiagnostics
//   PRESENCE optional },
//        ...
//   }
//=====================================================================================================

#include <iostream>
#include "util_ngap_pkt.h"

void fill_nR_CGI_cell_identity(Ngap_NRCellIdentity_t& nRCellIdentity) {
  uint64_t nr_cell_id; /* 36 bit */

  nr_cell_id          = 0x0000000100;
  nRCellIdentity.size = 5;

  nRCellIdentity.buf = (uint8_t*) calloc(nRCellIdentity.size, sizeof(uint8_t));
  memset(nRCellIdentity.buf, 0, (nRCellIdentity.size * sizeof(uint8_t)));

  nRCellIdentity.buf[0]      = (nr_cell_id >> 32);
  nRCellIdentity.buf[1]      = (nr_cell_id >> 24);
  nRCellIdentity.buf[2]      = (nr_cell_id >> 16);
  nRCellIdentity.buf[3]      = (nr_cell_id >> 8);
  nRCellIdentity.buf[4]      = (nr_cell_id);
  nRCellIdentity.bits_unused = 4;
}

void fill_nR_CGI_pLMNIdentity(Ngap_PLMNIdentity_t& pLMNIdentity) {
  pLMNIdentity.size = 3;
  pLMNIdentity.buf  = (uint8_t*) calloc(1, pLMNIdentity.size * sizeof(uint8_t));
  pLMNIdentity.buf[0] = 0x9;
  pLMNIdentity.buf[1] = 0xf1;
  pLMNIdentity.buf[2] = 0x7;
}

void fill_tAI_pLMNIdentity(Ngap_PLMNIdentity_t& pLMNIdentity) {
  pLMNIdentity.size = 3;
  pLMNIdentity.buf = (uint8_t*) calloc(1, sizeof(uint8_t*) * pLMNIdentity.size);
  pLMNIdentity.buf[0] = 0x9;
  pLMNIdentity.buf[1] = 0xf1;
  pLMNIdentity.buf[2] = 0x7;
}

void fill_tAI_tAC(Ngap_TAC_t& tAC) {
  tAC.size   = 3;
  tAC.buf    = (uint8_t*) calloc(1, sizeof(uint8_t*) * tAC.size);
  tAC.buf[0] = 0;
  tAC.buf[1] = 0;
  tAC.buf[2] = 0x1;
}

int encode_initate_ue_message(
    Ngap_NGAP_PDU_t* pdu, uint8_t** buffer, uint32_t* len) {
  asn_encode_to_new_buffer_result_t res = {NULL, {0, NULL, NULL}};

  memset(&res, 0, sizeof(res));

  res = asn_encode_to_new_buffer(
      NULL, ATS_ALIGNED_CANONICAL_PER, &asn_DEF_Ngap_NGAP_PDU, pdu);

  *buffer = (uint8_t*) res.buffer;
  *len    = res.result.encoded;

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, pdu);
  return (0);
}

bool ng_setup_initiate_ue_message_decode(
    const_bstring const raw, Ngap_NGAP_PDU_t* pdu) {
  asn_dec_rval_t dec_ret;

  dec_ret = aper_decode(
      NULL, &asn_DEF_Ngap_NGAP_PDU, (void**) &pdu, bdata(raw), blength(raw), 0,
      0);

  if (dec_ret.code != RC_OK) {
    return false;
  }

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, pdu);
  return true;
}

bool ngap_initiate_ue_message(bstring& stream_initate_ue) {
  Ngap_NGAP_PDU_t pdu;
  Ngap_NGAP_PDU_t dec_pdu;
  Ngap_InitialUEMessage_t* out;
  Ngap_InitialUEMessage_IEs_t* ie;
  Ngap_UserLocationInformationNR_t* userinfo_nr_p = NULL;
  uint8_t* buffer                                 = NULL;
  uint32_t length                                 = 0;
  int hexbuf[] = {0x7E, 0x00, 0x41, 0x79, 0x00, 0x0D, 0x01, 0x09, 0xF1, 0x07,
                  0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x10,
                  0x01, 0x00, 0x2E, 0x04, 0xF0, 0xF0, 0xF0, 0xF0, 0x2F, 0x05,
                  0x04, 0x01, 0x00, 0x00, 0x01, 0x53, 0x01, 0x00};

  memset(&pdu, 0, sizeof(pdu));
  pdu.present = Ngap_NGAP_PDU_PR_initiatingMessage;
  // pdu.choice.initiatingMessage = (NGAP_InitiatingMessage_t
  // *)calloc(1,sizeof(NGAP_InitiatingMessage_t));
  pdu.choice.initiatingMessage.procedureCode =
      Ngap_ProcedureCode_id_InitialUEMessage;
  pdu.choice.initiatingMessage.criticality = Ngap_Criticality_ignore;
  pdu.choice.initiatingMessage.value.present =
      Ngap_InitiatingMessage__value_PR_InitialUEMessage;
  out = &pdu.choice.initiatingMessage.value.choice.InitialUEMessage;

  /* mandatory */
  ie = (Ngap_InitialUEMessage_IEs_t*) calloc(
      1, sizeof(Ngap_InitialUEMessage_IEs_t));
  ie->id            = Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID;
  ie->criticality   = Ngap_Criticality_reject;
  ie->value.present = Ngap_InitialUEMessage_IEs__value_PR_RAN_UE_NGAP_ID;
  ie->value.choice.RAN_UE_NGAP_ID = 1;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* mandatory */
  ie = (Ngap_InitialUEMessage_IEs_t*) calloc(
      1, sizeof(Ngap_InitialUEMessage_IEs_t));
  ie->id            = Ngap_ProtocolIE_ID_id_NAS_PDU;
  ie->criticality   = Ngap_Criticality_reject;
  ie->value.present = Ngap_InitialUEMessage_IEs__value_PR_NAS_PDU;

  ie->value.choice.NAS_PDU.size = 38;
  ie->value.choice.NAS_PDU.buf =
      (uint8_t*) calloc(1, ie->value.choice.NAS_PDU.size * sizeof(uint8_t));
  memset(
      ie->value.choice.NAS_PDU.buf, 0,
      sizeof(ie->value.choice.NAS_PDU.size * sizeof(uint8_t)));
  for (uint32_t i = 0; i < ie->value.choice.NAS_PDU.size; i++) {
    ie->value.choice.NAS_PDU.buf[i] = hexbuf[i];
  }

  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* mandatory */
  ie = (Ngap_InitialUEMessage_IEs_t*) calloc(
      1, sizeof(Ngap_InitialUEMessage_IEs_t));
  ie->id          = Ngap_ProtocolIE_ID_id_UserLocationInformation;
  ie->criticality = Ngap_Criticality_reject;
  ie->value.present =
      Ngap_InitialUEMessage_IEs__value_PR_UserLocationInformation;

  ie->value.choice.UserLocationInformation.present =
      Ngap_UserLocationInformation_PR_userLocationInformationNR;

  userinfo_nr_p = &(ie->value.choice.UserLocationInformation.choice
                        .userLocationInformationNR);

  /* Set nRCellIdentity. default userLocationInformationNR */
  fill_nR_CGI_cell_identity(userinfo_nr_p->nR_CGI.nRCellIdentity);
  fill_nR_CGI_pLMNIdentity(userinfo_nr_p->nR_CGI.pLMNIdentity);

  fill_tAI_pLMNIdentity(userinfo_nr_p->tAI.pLMNIdentity);
  fill_tAI_tAC(userinfo_nr_p->tAI.tAC);

  userinfo_nr_p->timeStamp =
      (Ngap_TimeStamp_t*) calloc(1, sizeof(Ngap_TimeStamp_t));
  userinfo_nr_p->timeStamp->size = 4;
  userinfo_nr_p->timeStamp->buf =
      (uint8_t*) calloc(1, sizeof(uint8_t*) * userinfo_nr_p->timeStamp->size);
  userinfo_nr_p->timeStamp->buf[0] = 0xe4;
  userinfo_nr_p->timeStamp->buf[1] = 0x31;
  userinfo_nr_p->timeStamp->buf[2] = 0x20;
  userinfo_nr_p->timeStamp->buf[3] = 0x41;

  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* mandatory */
  ie = (Ngap_InitialUEMessage_IEs_t*) calloc(
      1, sizeof(Ngap_InitialUEMessage_IEs_t));
  ie->id            = Ngap_ProtocolIE_ID_id_RRCEstablishmentCause;
  ie->criticality   = Ngap_Criticality_ignore;
  ie->value.present = Ngap_InitialUEMessage_IEs__value_PR_RRCEstablishmentCause;
  ie->value.choice.RRCEstablishmentCause = Ngap_RRCEstablishmentCause_mo_Data;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* optional */
  ie = (Ngap_InitialUEMessage_IEs_t*) calloc(
      1, sizeof(Ngap_InitialUEMessage_IEs_t));
  ie->id            = Ngap_ProtocolIE_ID_id_UEContextRequest;
  ie->criticality   = Ngap_Criticality_ignore;
  ie->value.present = Ngap_InitialUEMessage_IEs__value_PR_UEContextRequest;
  ie->value.choice.UEContextRequest = Ngap_UEContextRequest_requested;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  if (encode_initate_ue_message(&pdu, &buffer, &length) != 0) {
    return false;
  }

  stream_initate_ue = blk2bstr(buffer, length);

  free(buffer);
  return (true);
}
