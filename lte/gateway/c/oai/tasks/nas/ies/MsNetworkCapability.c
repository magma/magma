/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#include "TLVEncoder.h"
#include "TLVDecoder.h"
#include "MsNetworkCapability.h"

int decode_ms_network_capability(
    MsNetworkCapability* msnetworkcapability, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded   = 0;
  uint8_t ielen = 0;
  uint8_t b     = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  DECODE_U8(buffer + decoded, ielen, decoded);
  memset(msnetworkcapability, 0, sizeof(MsNetworkCapability));
  OAILOG_INFO(LOG_NAS_EMM, "decode_ms_network_capability len = %d\n", ielen);
  CHECK_LENGTH_DECODER(len - decoded, ielen);

  b                         = *(buffer + decoded);
  msnetworkcapability->gea1 = (b & MS_NETWORK_CAPABILITY_GEA1) >> 7;
  msnetworkcapability->smdc =
      (b & MS_NETWORK_CAPABILITY_SM_CAP_VIA_DEDICATED_CHANNELS) >> 6;
  msnetworkcapability->smgc =
      (b & MS_NETWORK_CAPABILITY_SM_CAP_VIA_GPRS_CHANNELS) >> 5;
  msnetworkcapability->ucs2 = (b & MS_NETWORK_CAPABILITY_UCS2_SUPPORT) >> 4;
  msnetworkcapability->sssi =
      (b & MS_NETWORK_CAPABILITY_SS_SCREENING_INDICATOR) >> 2;
  msnetworkcapability->solsa = (b & MS_NETWORK_CAPABILITY_SOLSA) >> 1;
  msnetworkcapability->revli =
      (b & MS_NETWORK_CAPABILITY_REVISION_LEVEL_INDICATOR);
  decoded++;

  if (ielen > 1) {
    b = *(buffer + decoded);
    msnetworkcapability->pfc =
        (b & MS_NETWORK_CAPABILITY_PFC_FEATURE_MODE) >> 7;
    msnetworkcapability->egea =
        (b & (MS_NETWORK_CAPABILITY_GEA2 | MS_NETWORK_CAPABILITY_GEA3 |
              MS_NETWORK_CAPABILITY_GEA4 | MS_NETWORK_CAPABILITY_GEA5 |
              MS_NETWORK_CAPABILITY_GEA6 | MS_NETWORK_CAPABILITY_GEA7)) >>
        1;
    msnetworkcapability->lcs = (b & MS_NETWORK_CAPABILITY_LCS_VA);
    decoded++;

    if (ielen > 2) {
      b = *(buffer + decoded);
      msnetworkcapability->ps_ho_utran =
          (b & MS_NETWORK_CAPABILITY_PS_INTER_RAT_HO_GERAN_TO_UTRAN_IU) >> 7;
      msnetworkcapability->ps_ho_eutran =
          (b & MS_NETWORK_CAPABILITY_PS_INTER_RAT_HO_GERAN_TO_EUTRAN_S1) >> 6;
      msnetworkcapability->emm_cpc =
          (b & MS_NETWORK_CAPABILITY_EMM_COMBINED_PROCEDURE) >> 5;
      msnetworkcapability->isr     = (b & MS_NETWORK_CAPABILITY_ISR) >> 4;
      msnetworkcapability->srvcc   = (b & MS_NETWORK_CAPABILITY_SRVCC) >> 3;
      msnetworkcapability->epc_cap = (b & MS_NETWORK_CAPABILITY_EPC) >> 2;
      msnetworkcapability->nf_cap =
          (b & MS_NETWORK_CAPABILITY_NOTIFICATION) >> 1;
      msnetworkcapability->geran_ns =
          (b & MS_NETWORK_CAPABILITY_GERAN_NETWORK_SHARING);
      decoded++;
    }
    if (ielen > 3) {
      b = *(buffer + decoded);
      msnetworkcapability->up_integ_prot_support =
          (b & MS_NETWORK_CAPABILITY_USER_PLANE_INTEGRITY_PROTECTION_SUPPORT) >>
          7;
      msnetworkcapability->gia4 = (b & MS_NETWORK_CAPABILITY_GIA4) >> 6;
      msnetworkcapability->gia5 = (b & MS_NETWORK_CAPABILITY_GIA5) >> 5;
      msnetworkcapability->gia6 = (b & MS_NETWORK_CAPABILITY_GIA6) >> 4;
      msnetworkcapability->gia7 = (b & MS_NETWORK_CAPABILITY_GIA7) >> 3;
      msnetworkcapability->epco_ie_ind =
          (b & MS_NETWORK_CAPABILITY_EPCO_IE_INDICATOR) >> 2;
      msnetworkcapability->rest_use_enhanc_cov_cap =
          (b &
           MS_NETWORK_CAPABILITY_RESTRICTION_ON_USE_OF_ENHANCED_COVERAGE_CAPABILITY) >>
          1;
      msnetworkcapability->en_dc =
          (b & MS_NETWORK_CAPABILITY_DUAL_CONNECTIVITY_EUTRA_NR_CAPABILITY);
      decoded++;
    }
  }

#if NAS_DEBUG
  dump_ms_network_capability_xml(msnetworkcapability, iei);
#endif
  return decoded;
}

int encode_ms_network_capability(
    MsNetworkCapability* msnetworkcapability, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, MS_NETWORK_CAPABILITY_MINIMUM_LENGTH, len);
#if NAS_DEBUG
  dump_ms_network_capability_xml(msnetworkcapability, iei);
#endif

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;

  *(buffer + encoded) = ((msnetworkcapability->gea1 & 0x1) << 7) |
                        ((msnetworkcapability->smdc & 0x1) << 6) |
                        ((msnetworkcapability->smgc & 0x1) << 5) |
                        ((msnetworkcapability->ucs2 & 0x1) << 4) |
                        ((msnetworkcapability->sssi & 0x3) << 2) |
                        ((msnetworkcapability->solsa & 0x1) << 1) |
                        (msnetworkcapability->revli & 0x1);
  encoded++;

  *(buffer + encoded) = ((msnetworkcapability->pfc & 0x1) << 7) |
                        ((msnetworkcapability->egea & 0x3F) << 1) |
                        (msnetworkcapability->lcs & 0x1);
  encoded++;

  *(buffer + encoded) = ((msnetworkcapability->ps_ho_utran & 0x1) << 7) |
                        ((msnetworkcapability->ps_ho_eutran & 0x1) << 6) |
                        ((msnetworkcapability->emm_cpc & 0x1) << 5) |
                        ((msnetworkcapability->isr & 0x1) << 4) |
                        ((msnetworkcapability->srvcc & 0x1) << 3) |
                        ((msnetworkcapability->epc_cap & 0x1) << 2) |
                        ((msnetworkcapability->nf_cap & 0x1) << 1) |
                        (msnetworkcapability->geran_ns & 0x1);
  encoded++;

  *(buffer + encoded) =
      ((msnetworkcapability->up_integ_prot_support & 0x1) << 7) |
      ((msnetworkcapability->gia4 & 0x1) << 6) |
      ((msnetworkcapability->gia5 & 0x1) << 5) |
      ((msnetworkcapability->gia6 & 0x1) << 4) |
      ((msnetworkcapability->gia7 & 0x1) << 3) |
      ((msnetworkcapability->epco_ie_ind & 0x1) << 2) |
      ((msnetworkcapability->rest_use_enhanc_cov_cap & 0x1) << 1) |
      (msnetworkcapability->en_dc & 0x1);
  encoded++;

  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);
  return encoded;
}

void dump_ms_network_capability_xml(
    MsNetworkCapability* msnetworkcapability, uint8_t iei) {
  OAILOG_DEBUG(LOG_NAS, "<Ms Network Capability>\n");

  if (iei > 0)
    /*
     * Don't display IEI if = 0
     */
    OAILOG_DEBUG(LOG_NAS, "    <IEI>0x%X</IEI>\n", iei);

  OAILOG_DEBUG(LOG_NAS, "    <GEA1>%01x</GEA1>\n", msnetworkcapability->gea1);
  OAILOG_DEBUG(LOG_NAS, "    <SMDC>%01x</SMDC>\n", msnetworkcapability->smdc);
  OAILOG_DEBUG(LOG_NAS, "    <SMGC>%01x</SMGC>\n", msnetworkcapability->smgc);
  OAILOG_DEBUG(LOG_NAS, "    <UCS2>%01x</UCS2>\n", msnetworkcapability->ucs2);
  OAILOG_DEBUG(LOG_NAS, "    <SSSI>%02x</SSSI>\n", msnetworkcapability->sssi);
  OAILOG_DEBUG(
      LOG_NAS, "    <SOLSA>%01x</SOLSA>\n", msnetworkcapability->solsa);
  OAILOG_DEBUG(
      LOG_NAS, "    <REVLI>%01x</REVLI>\n", msnetworkcapability->revli);
  OAILOG_DEBUG(LOG_NAS, "    <PFC>%01x</PFC>\n", msnetworkcapability->pfc);
  OAILOG_DEBUG(LOG_NAS, "    <EGEA>%06x</EGEA>\n", msnetworkcapability->egea);
  OAILOG_DEBUG(LOG_NAS, "    <LCS>%01x</LCS>\n", msnetworkcapability->lcs);
  OAILOG_DEBUG(
      LOG_NAS, "    <PS_HO_UTRAN>%01x<PS_HO_UTRAN/>\n",
      msnetworkcapability->ps_ho_utran);
  OAILOG_DEBUG(
      LOG_NAS, "    <PS_HO_EUTRAN>%01x<PS_HO_EUTRAN/>\n",
      msnetworkcapability->ps_ho_eutran);
  OAILOG_DEBUG(
      LOG_NAS, "    <EMM_CPC>%01x<EMM_CPC/>\n", msnetworkcapability->emm_cpc);
  OAILOG_DEBUG(LOG_NAS, "    <ISR>%01x<ISR/>\n", msnetworkcapability->isr);
  OAILOG_DEBUG(
      LOG_NAS, "    <SRVCC>%01x<SRVCC/>\n", msnetworkcapability->srvcc);
  OAILOG_DEBUG(
      LOG_NAS, "    <EPC_CAP>%01x<EPC_CAP/>\n", msnetworkcapability->epc_cap);
  OAILOG_DEBUG(
      LOG_NAS, "    <NF_CAP>%01x<NF_CAP/>\n", msnetworkcapability->nf_cap);
  OAILOG_DEBUG(
      LOG_NAS, "    <GERAN_NS>%01x<GERAN_NS/>\n",
      msnetworkcapability->geran_ns);
  OAILOG_DEBUG(
      LOG_NAS, "    <UP_INTEG_PROT_SUPPORT>%01x<UP_INTEG_PROT_SUPPORT/>\n",
      msnetworkcapability->up_integ_prot_support);
  OAILOG_DEBUG(LOG_NAS, "    <GIA4>%01x<GIA4/>\n", msnetworkcapability->gia4);
  OAILOG_DEBUG(LOG_NAS, "    <GIA5>%01x<GIA5/>\n", msnetworkcapability->gia5);
  OAILOG_DEBUG(LOG_NAS, "    <GIA6>%01x<GIA6/>\n", msnetworkcapability->gia6);
  OAILOG_DEBUG(LOG_NAS, "    <GIA7>%01x<GIA7/>\n", msnetworkcapability->gia7);
  OAILOG_DEBUG(
      LOG_NAS, "    <EPCO_IE_IND>%01x<EPCO_IE_IND/>\n",
      msnetworkcapability->epco_ie_ind);
  OAILOG_DEBUG(
      LOG_NAS, "    <REST_USE_ENHANC_COV_CAP>%01x<REST_USE_ENHANC_COV_CAP/>\n",
      msnetworkcapability->rest_use_enhanc_cov_cap);
  OAILOG_DEBUG(
      LOG_NAS, "    <EN_DC>%01x<EN_DC/>\n", msnetworkcapability->en_dc);
  OAILOG_DEBUG(LOG_NAS, "</Ms Network Capability>\n");
}
