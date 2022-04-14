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
#include "lte/gateway/c/core/oai/test/ngap/util_ngap_pkt.hpp"

uint8_t intialUeGuti[157] = {
    0x00, 0x0f, 0x40, 0x80, 0x98, 0x00, 0x00, 0x06, 0x00, 0x55, 0x00, 0x02,
    0x00, 0x04, 0x00, 0x26, 0x00, 0x63, 0x62, 0x7e, 0x01, 0x93, 0xa0, 0x67,
    0x05, 0x03, 0x7e, 0x00, 0x41, 0x01, 0x00, 0x0b, 0xf2, 0x22, 0xf2, 0x54,
    0x00, 0x00, 0x00, 0x70, 0xe5, 0xc5, 0x00, 0x2e, 0x04, 0x80, 0xe0, 0x80,
    0xe0, 0x71, 0x00, 0x41, 0x7e, 0x00, 0x41, 0x01, 0x00, 0x0b, 0xf2, 0x22,
    0xf2, 0x54, 0x00, 0x00, 0x00, 0x70, 0xe5, 0xc5, 0x00, 0x10, 0x01, 0x03,
    0x2e, 0x04, 0x80, 0xe0, 0x80, 0xe0, 0x2f, 0x02, 0x01, 0x01, 0x52, 0x22,
    0x62, 0x54, 0x00, 0x00, 0x01, 0x17, 0x07, 0x80, 0xe0, 0xe0, 0x60, 0x00,
    0x1c, 0x30, 0x18, 0x01, 0x00, 0x74, 0x00, 0x0a, 0x09, 0x08, 0x69, 0x6e,
    0x74, 0x65, 0x72, 0x6e, 0x65, 0x74, 0x53, 0x01, 0x01, 0x00, 0x79, 0x00,
    0x0f, 0x40, 0x22, 0x42, 0x65, 0x00, 0x00, 0x00, 0x01, 0x00, 0x22, 0x42,
    0x65, 0x00, 0x00, 0x01, 0x00, 0x5a, 0x40, 0x01, 0x18, 0x00, 0x1a, 0x00,
    0x07, 0x00, 0x00, 0x00, 0x70, 0xe5, 0xc5, 0x00, 0x00, 0x70, 0x40, 0x01,
    0x00};

uint8_t ngapSetupRequestWithSD[] = {
    0x00, 0x15, 0x00, 0x28, 0x00, 0x00, 0x03, 0x00, 0x1b, 0x00, 0x08,
    0x00, 0x27, 0xf4, 0x77, 0x00, 0x00, 0x00, 0x08, 0x00, 0x66, 0x00,
    0x10, 0x00, 0x00, 0x00, 0x00, 0x01, 0x00, 0x27, 0xf4, 0x77, 0x00,
    0x00, 0x10, 0x08, 0xd1, 0x43, 0xa5, 0x00, 0x15, 0x40, 0x01, 0x60};

void fill_nR_CGI_cell_identity(Ngap_NRCellIdentity_t& nRCellIdentity) {
  uint64_t nr_cell_id; /* 36 bit */

  nr_cell_id = 0x0000000100;
  nRCellIdentity.size = 5;

  nRCellIdentity.buf = (uint8_t*)calloc(nRCellIdentity.size, sizeof(uint8_t));
  memset(nRCellIdentity.buf, 0, (nRCellIdentity.size * sizeof(uint8_t)));

  nRCellIdentity.buf[0] = (nr_cell_id >> 32);
  nRCellIdentity.buf[1] = (nr_cell_id >> 24);
  nRCellIdentity.buf[2] = (nr_cell_id >> 16);
  nRCellIdentity.buf[3] = (nr_cell_id >> 8);
  nRCellIdentity.buf[4] = (nr_cell_id);
  nRCellIdentity.bits_unused = 4;
}

void fill_nR_CGI_pLMNIdentity(Ngap_PLMNIdentity_t& pLMNIdentity) {
  pLMNIdentity.size = 3;
  pLMNIdentity.buf = (uint8_t*)calloc(1, pLMNIdentity.size * sizeof(uint8_t));
  pLMNIdentity.buf[0] = 0x9;
  pLMNIdentity.buf[1] = 0xf1;
  pLMNIdentity.buf[2] = 0x7;
}

void fill_tAI_pLMNIdentity(Ngap_PLMNIdentity_t& pLMNIdentity) {
  pLMNIdentity.size = 3;
  pLMNIdentity.buf = (uint8_t*)calloc(1, sizeof(uint8_t*) * pLMNIdentity.size);
  pLMNIdentity.buf[0] = 0x9;
  pLMNIdentity.buf[1] = 0xf1;
  pLMNIdentity.buf[2] = 0x7;
}

void fill_tAI_tAC(Ngap_TAC_t& tAC) {
  tAC.size = 3;
  tAC.buf = (uint8_t*)calloc(1, sizeof(uint8_t*) * tAC.size);
  tAC.buf[0] = 0;
  tAC.buf[1] = 0;
  tAC.buf[2] = 0x1;
}

int encode_initate_ue_message(Ngap_NGAP_PDU_t* pdu, uint8_t** buffer,
                              uint32_t* len) {
  asn_encode_to_new_buffer_result_t res = {NULL, {0, NULL, NULL}};

  memset(&res, 0, sizeof(res));

  res = asn_encode_to_new_buffer(NULL, ATS_ALIGNED_CANONICAL_PER,
                                 &asn_DEF_Ngap_NGAP_PDU, pdu);

  *buffer = (uint8_t*)res.buffer;
  *len = res.result.encoded;

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, pdu);
  return (0);
}

bool ng_setup_initiate_ue_message_decode(const_bstring const raw,
                                         Ngap_NGAP_PDU_t* pdu) {
  asn_dec_rval_t dec_ret;

  dec_ret = aper_decode(NULL, &asn_DEF_Ngap_NGAP_PDU, (void**)&pdu, bdata(raw),
                        blength(raw), 0, 0);

  if (dec_ret.code != RC_OK) {
    return false;
  }

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, pdu);
  return true;
}

bool ngap_initiate_ue_message(bstring& stream_initate_ue) {
  Ngap_NGAP_PDU_t pdu;
  Ngap_InitialUEMessage_t* out;
  Ngap_InitialUEMessage_IEs_t* ie;
  Ngap_UserLocationInformationNR_t* userinfo_nr_p = NULL;
  uint8_t* buffer = NULL;
  uint32_t length = 0;
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
  ie = (Ngap_InitialUEMessage_IEs_t*)calloc(
      1, sizeof(Ngap_InitialUEMessage_IEs_t));
  ie->id = Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID;
  ie->criticality = Ngap_Criticality_reject;
  ie->value.present = Ngap_InitialUEMessage_IEs__value_PR_RAN_UE_NGAP_ID;
  ie->value.choice.RAN_UE_NGAP_ID = 1;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* mandatory */
  ie = (Ngap_InitialUEMessage_IEs_t*)calloc(
      1, sizeof(Ngap_InitialUEMessage_IEs_t));
  ie->id = Ngap_ProtocolIE_ID_id_NAS_PDU;
  ie->criticality = Ngap_Criticality_reject;
  ie->value.present = Ngap_InitialUEMessage_IEs__value_PR_NAS_PDU;

  ie->value.choice.NAS_PDU.size = 38;
  ie->value.choice.NAS_PDU.buf =
      (uint8_t*)calloc(1, ie->value.choice.NAS_PDU.size * sizeof(uint8_t));
  memset(ie->value.choice.NAS_PDU.buf, 0,
         sizeof(ie->value.choice.NAS_PDU.size * sizeof(uint8_t)));
  for (uint32_t i = 0; i < ie->value.choice.NAS_PDU.size; i++) {
    ie->value.choice.NAS_PDU.buf[i] = hexbuf[i];
  }

  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* mandatory */
  ie = (Ngap_InitialUEMessage_IEs_t*)calloc(
      1, sizeof(Ngap_InitialUEMessage_IEs_t));
  ie->id = Ngap_ProtocolIE_ID_id_UserLocationInformation;
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
      (Ngap_TimeStamp_t*)calloc(1, sizeof(Ngap_TimeStamp_t));
  userinfo_nr_p->timeStamp->size = 4;
  userinfo_nr_p->timeStamp->buf =
      (uint8_t*)calloc(1, sizeof(uint8_t*) * userinfo_nr_p->timeStamp->size);
  userinfo_nr_p->timeStamp->buf[0] = 0xe4;
  userinfo_nr_p->timeStamp->buf[1] = 0x31;
  userinfo_nr_p->timeStamp->buf[2] = 0x20;
  userinfo_nr_p->timeStamp->buf[3] = 0x41;

  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* mandatory */
  ie = (Ngap_InitialUEMessage_IEs_t*)calloc(
      1, sizeof(Ngap_InitialUEMessage_IEs_t));
  ie->id = Ngap_ProtocolIE_ID_id_RRCEstablishmentCause;
  ie->criticality = Ngap_Criticality_ignore;
  ie->value.present = Ngap_InitialUEMessage_IEs__value_PR_RRCEstablishmentCause;
  ie->value.choice.RRCEstablishmentCause = Ngap_RRCEstablishmentCause_mo_Data;
  ASN_SEQUENCE_ADD(&out->protocolIEs.list, ie);

  /* optional */
  ie = (Ngap_InitialUEMessage_IEs_t*)calloc(
      1, sizeof(Ngap_InitialUEMessage_IEs_t));
  ie->id = Ngap_ProtocolIE_ID_id_UEContextRequest;
  ie->criticality = Ngap_Criticality_ignore;
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

bool generate_guti_ngap_pdu(Ngap_NGAP_PDU_t* pdu) {
  asn_dec_rval_t dec_ret;
  uint32_t guti_len = 157;

  dec_ret = aper_decode(NULL, &asn_DEF_Ngap_NGAP_PDU, (void**)&pdu,
                        intialUeGuti, guti_len, 0, 0);

  if (dec_ret.code != RC_OK) {
    return false;
  }

  return true;
}

bool validate_handle_initial_ue_message(gnb_description_t* gNB_ref,
                                        m5g_ue_description_t* ue_ref,
                                        Ngap_NGAP_PDU_t* pdu) {
  tai_t tai = {0};
  guamfi_t guamfi = {{0}, 0, 0};
  s_tmsi_m5_t s_tmsi = {0, 0, INVALID_M_TMSI};
  gnb_ue_ngap_id_t gnb_ue_ngap_id = 0;
  ecgi_t ecgi = {{0}, {0}};
  csg_id_t csg_id = 0;

  Ngap_InitialUEMessage_t* container;
  Ngap_InitialUEMessage_IEs_t *ie = NULL, *ie_e_tmsi, *ie_csg_id = NULL,
                              *ie_guamfi, *ie_cause;

  container = &pdu->choice.initiatingMessage.value.choice.InitialUEMessage;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_InitialUEMessage_IEs_t, ie,
                                      container,
                                      Ngap_ProtocolIE_ID_id_RAN_UE_NGAP_ID);

  // gNB UE NGAP ID is limited to 24 bits
  gnb_ue_ngap_id = (gnb_ue_ngap_id_t)(ie->value.choice.RAN_UE_NGAP_ID);

  ue_ref->ng_ue_state = NGAP_UE_WAITING_CSR;

  ue_ref->gnb_ue_ngap_id = gnb_ue_ngap_id;

  // Will be allocated by NAS
  ue_ref->amf_ue_ngap_id = INVALID_AMF_UE_NGAP_ID;

  ue_ref->ngap_ue_context_rel_timer.id = NGAP_TIMER_INACTIVE_ID;
  ue_ref->ngap_ue_context_rel_timer.msec =
      1000 * NGAP_UE_CONTEXT_REL_COMP_TIMER;

  // On which stream we received the message
  ue_ref->sctp_stream_recv = 2;
  ue_ref->sctp_stream_send = gNB_ref->next_sctp_stream;

  gNB_ref->next_sctp_stream += 1;
  if (gNB_ref->next_sctp_stream >= gNB_ref->instreams) {
    gNB_ref->next_sctp_stream = 1;
  }

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(
      Ngap_InitialUEMessage_IEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_UserLocationInformation);

  OCTET_STRING_TO_TAC_5G(&ie->value.choice.UserLocationInformation.choice
                              .userLocationInformationNR.tAI.tAC,
                         tai.tac);

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(
      Ngap_InitialUEMessage_IEs_t, ie, container,
      Ngap_ProtocolIE_ID_id_UserLocationInformation);

  TBCD_TO_PLMN_T(&ie->value.choice.UserLocationInformation.choice
                      .userLocationInformationEUTRA.eUTRA_CGI.pLMNIdentity,
                 &ecgi.plmn);

  BIT_STRING_TO_CELL_IDENTITY(
      &ie->value.choice.UserLocationInformation.choice
           .userLocationInformationEUTRA.eUTRA_CGI.eUTRACellIdentity,
      ecgi.cell_identity);

  /** Set the GNB Id. */
  ecgi.cell_identity.enb_id = gNB_ref->gnb_id;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_InitialUEMessage_IEs_t, ie_e_tmsi,
                                      container,
                                      Ngap_ProtocolIE_ID_id_FiveG_S_TMSI);

  if (ie_e_tmsi) {
    NGAP_TEST_PDU_FETCH_AMF_SET_ID_FROM_PDU(
        ie_e_tmsi->value.choice.FiveG_S_TMSI.aMFSetID, s_tmsi.amf_set_id);
    OCTET_STRING_TO_M_TMSI(&ie_e_tmsi->value.choice.FiveG_S_TMSI.fiveG_TMSI,
                           s_tmsi.m_tmsi);
  }

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_InitialUEMessage_IEs_t, ie_guamfi,
                                      container, Ngap_ProtocolIE_ID_id_GUAMI);

  memset(&guamfi, 0, sizeof(guamfi));

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_InitialUEMessage_IEs_t, ie,
                                      container, Ngap_ProtocolIE_ID_id_NAS_PDU);

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(
      Ngap_InitialUEMessage_IEs_t, ie_cause, container,
      Ngap_ProtocolIE_ID_id_RRCEstablishmentCause);

  Ngap_InitialUEMessage_IEs_t* ie_uecontextrequest = NULL;

  NGAP_TEST_PDU_FIND_PROTOCOLIE_BY_ID(Ngap_InitialUEMessage_IEs_t,
                                      ie_uecontextrequest, container,
                                      Ngap_ProtocolIE_ID_id_UEContextRequest);

  long ue_context_request = 0;
  if (ie_uecontextrequest) {
    ue_context_request = ie_uecontextrequest->value.choice.UEContextRequest + 1;
  }

  return (true);
}

bool generate_ngap_request_msg(Ngap_NGAP_PDU_t* pdu) {
  asn_dec_rval_t dec_ret;
  uint32_t pkt_len = sizeof(ngapSetupRequestWithSD) / sizeof(uint8_t);

  dec_ret = aper_decode(NULL, &asn_DEF_Ngap_NGAP_PDU, (void**)&pdu,
                        ngapSetupRequestWithSD, pkt_len, 0, 0);

  if (dec_ret.code != RC_OK) {
    return false;
  }

  return true;
}

bool validate_ngap_setup_request(Ngap_NGAP_PDU_t* pdu) {
  Ngap_NGSetupRequest_t* container = NULL;
  Ngap_NGSetupRequestIEs_t* ie = NULL;
  Ngap_NGSetupRequestIEs_t* ie_gnb_name = NULL;
  Ngap_NGSetupRequestIEs_t* ie_supported_tas = NULL;
  Ngap_NGSetupRequestIEs_t* ie_default_paging_drx = NULL;

  gnb_description_t* gnb_association = NULL;
  uint32_t gnb_id = 0;
  char* gnb_name = NULL;
  int ta_ret = 0;
  uint8_t bplmn_list_count = 0;  // Broadcast PLMN list count

  container = &pdu->choice.initiatingMessage.value.choice.NGSetupRequest;

  NGAP_FIND_PROTOCOLIE_BY_ID(Ngap_NGSetupRequestIEs_t, ie_gnb_name, container,
                             Ngap_ProtocolIE_ID_id_RANNodeName, false);
  if (ie_gnb_name) {
    gnb_name = (char*)ie_gnb_name->value.choice.RANNodeName.buf;
  }

  NGAP_FIND_PROTOCOLIE_BY_ID(Ngap_NGSetupRequestIEs_t, ie, container,
                             Ngap_ProtocolIE_ID_id_GlobalRANNodeID, true);
  if (ie->value.choice.GlobalRANNodeID.choice.globalGNB_ID.gNB_ID.present ==
      Ngap_GNB_ID_PR_gNB_ID) {
    // Home gNB ID = 28 bits
    uint8_t* gnb_id_buf = ie->value.choice.GlobalRANNodeID.choice.globalGNB_ID
                              .gNB_ID.choice.gNB_ID.buf;

    if (ie->value.choice.GlobalRANNodeID.choice.globalGNB_ID.gNB_ID.choice
            .gNB_ID.size != 28) {
      // TODO: handle case were size != 28 -> notify ? reject ?
    }

    gnb_id = (gnb_id_buf[0] << 20) + (gnb_id_buf[1] << 12) +
             (gnb_id_buf[2] << 4) + ((gnb_id_buf[3] & 0xf0) >> 4);
  }

  NGAP_FIND_PROTOCOLIE_BY_ID(Ngap_NGSetupRequestIEs_t, ie_supported_tas,
                             container, Ngap_ProtocolIE_ID_id_SupportedTAList,
                             true);

  Ngap_SupportedTAList_t* ta_list =
      &ie_supported_tas->value.choice.SupportedTAList;
  m5g_supported_ta_list_t gnb_association_supported_ta_list;
  memset(&gnb_association_supported_ta_list, 0,
         sizeof(m5g_supported_ta_list_t));

  m5g_supported_ta_list_t* supp_ta_list = &gnb_association_supported_ta_list;
  supp_ta_list->list_count = ta_list->list.count;

  /* Storing supported TAI lists received in Ng SETUP REQUEST message */
  for (int tai_idx = 0; tai_idx < supp_ta_list->list_count; tai_idx++) {
    Ngap_SupportedTAItem_t* tai = NULL;
    tai = ta_list->list.array[tai_idx];
    tai->tAC.size = 2;  // ACL_TAG temp to test remove later
    OCTET_STRING_TO_TAC(&tai->tAC,
                        supp_ta_list->supported_tai_items[tai_idx].tac);

    bplmn_list_count = tai->broadcastPLMNList.list.count;
    if (bplmn_list_count > NGAP_MAX_BROADCAST_PLMNS) {
      return false;
    }
    supp_ta_list->supported_tai_items[tai_idx].bplmnlist_count =
        bplmn_list_count;

    for (int plmn_idx = 0; plmn_idx < bplmn_list_count; plmn_idx++) {
      TBCD_TO_PLMN_T(&tai->broadcastPLMNList.list.array[plmn_idx]->pLMNIdentity,
                     &supp_ta_list->supported_tai_items[tai_idx]
                          .bplmn_list[plmn_idx]
                          .plmn_id);

      supp_ta_list->supported_tai_items[tai_idx]
          .bplmn_list[plmn_idx]
          .num_of_s_nssai = tai->broadcastPLMNList.list.array[plmn_idx]
                                ->tAISliceSupportList.list.count;

      for (int nssai_index = 0;
           nssai_index < supp_ta_list->supported_tai_items[tai_idx]
                             .bplmn_list[plmn_idx]
                             .num_of_s_nssai;
           nssai_index++) {
        supp_ta_list->supported_tai_items[tai_idx]
            .bplmn_list[plmn_idx]
            .s_nssai[nssai_index]
            .sst = tai->broadcastPLMNList.list.array[plmn_idx]
                       ->tAISliceSupportList.list.array[nssai_index]
                       ->s_NSSAI.sST.buf[0];

        if (tai->broadcastPLMNList.list.array[plmn_idx]
                ->tAISliceSupportList.list.array[nssai_index]
                ->s_NSSAI.sD) {
          memcpy(&(supp_ta_list->supported_tai_items[tai_idx]
                       .bplmn_list[plmn_idx]
                       .s_nssai[nssai_index]
                       .sd),
                 tai->broadcastPLMNList.list.array[plmn_idx]
                     ->tAISliceSupportList.list.array[nssai_index]
                     ->s_NSSAI.sD->buf,
                 sizeof(amf_s_nssai_t));
        }
      }
    }
  }

  return true;
}
