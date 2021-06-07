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

#include <stdint.h>
#include <string.h>

#include "s1ap_common.h"
#include "s1ap_mme_encoder.h"
#include "assertions.h"
#include "log.h"
#include "common_defs.h"

static inline int s1ap_mme_encode_initiating(
    S1ap_S1AP_PDU_t* pdu, uint8_t** buffer, uint32_t* length);
static inline int s1ap_mme_encode_successfull_outcome(
    S1ap_S1AP_PDU_t* pdu, uint8_t** buffer, uint32_t* len);
static inline int s1ap_mme_encode_unsuccessfull_outcome(
    S1ap_S1AP_PDU_t* pdu, uint8_t** buffer, uint32_t* len);
//------------------------------------------------------------------------------
int s1ap_mme_encode_pdu(
    S1ap_S1AP_PDU_t* pdu, uint8_t** buffer, uint32_t* length) {
  int ret = -1;

  if (pdu == NULL) {
    OAILOG_DEBUG(LOG_S1AP, "PDU is NULL\n");
    return RETURNerror;
  }
  if (buffer == NULL) {
    OAILOG_DEBUG(LOG_S1AP, "Buffer is NULL\n");
    return RETURNerror;
  }
  if (length == NULL) {
    OAILOG_DEBUG(LOG_S1AP, "Length is NULL\n");
    return RETURNerror;
  }

  switch (pdu->present) {
    case S1ap_S1AP_PDU_PR_initiatingMessage:
      ret = s1ap_mme_encode_initiating(pdu, buffer, length);
      break;

    case S1ap_S1AP_PDU_PR_successfulOutcome:
      ret = s1ap_mme_encode_successfull_outcome(pdu, buffer, length);
      break;

    case S1ap_S1AP_PDU_PR_unsuccessfulOutcome:
      ret = s1ap_mme_encode_unsuccessfull_outcome(pdu, buffer, length);
      break;

    default:
      OAILOG_DEBUG(
          LOG_S1AP, "Unknown message outcome (%d) or not implemented",
          (int) pdu->present);
      break;
  }
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_S1ap_S1AP_PDU, pdu);
  return ret;
}

//------------------------------------------------------------------------------
static inline int s1ap_mme_encode_initiating(
    S1ap_S1AP_PDU_t* pdu, uint8_t** buffer, uint32_t* length) {
  asn_encode_to_new_buffer_result_t res = {NULL, {0, NULL, NULL}};

  if (pdu == NULL) {
    OAILOG_ERROR(LOG_S1AP, "PDU is NULL\n");
    return RETURNerror;
  }

  switch (pdu->choice.initiatingMessage.procedureCode) {
    case S1ap_ProcedureCode_id_downlinkNASTransport:
    case S1ap_ProcedureCode_id_InitialContextSetup:
    case S1ap_ProcedureCode_id_UEContextRelease:
    case S1ap_ProcedureCode_id_E_RABSetup:
    case S1ap_ProcedureCode_id_E_RABModify:
    case S1ap_ProcedureCode_id_E_RABRelease:
    case S1ap_ProcedureCode_id_HandoverResourceAllocation:
    case S1ap_ProcedureCode_id_MMEStatusTransfer:
    case S1ap_ProcedureCode_id_Paging:
    case S1ap_ProcedureCode_id_MMEConfigurationTransfer:
    case S1ap_ProcedureCode_id_HandoverPreparation:
      break;

    default:
      OAILOG_NOTICE(
          LOG_S1AP, "Unknown procedure ID (%d) for initiating message_p\n",
          (int) pdu->choice.initiatingMessage.procedureCode);
      *buffer = NULL;
      *length = 0;
      return RETURNerror;
  }

  memset(&res, 0, sizeof(res));
  res = asn_encode_to_new_buffer(
      NULL, ATS_ALIGNED_CANONICAL_PER, &asn_DEF_S1ap_S1AP_PDU, pdu);
  *buffer = res.buffer;
  *length = res.result.encoded;
  return RETURNok;
}

//------------------------------------------------------------------------------
static inline int s1ap_mme_encode_successfull_outcome(
    S1ap_S1AP_PDU_t* pdu, uint8_t** buffer, uint32_t* length) {
  asn_encode_to_new_buffer_result_t res = {NULL, {0, NULL, NULL}};

  if (pdu == NULL) {
    OAILOG_ERROR(LOG_S1AP, "PDU is NULL\n");
    return RETURNerror;
  }
  switch (pdu->choice.successfulOutcome.procedureCode) {
    case S1ap_ProcedureCode_id_S1Setup:
    case S1ap_ProcedureCode_id_PathSwitchRequest:
    case S1ap_ProcedureCode_id_HandoverPreparation:
    case S1ap_ProcedureCode_id_HandoverCancel:
    case S1ap_ProcedureCode_id_Reset:
    case S1ap_ProcedureCode_id_E_RABModificationIndication:
      break;

    default:
      OAILOG_DEBUG(
          LOG_S1AP,
          "Unknown procedure ID (%d) for successfull outcome message\n",
          (int) pdu->choice.successfulOutcome.procedureCode);
      *buffer = NULL;
      *length = 0;
      return RETURNerror;
  }
  res = asn_encode_to_new_buffer(
      NULL, ATS_ALIGNED_CANONICAL_PER, &asn_DEF_S1ap_S1AP_PDU, pdu);
  *buffer = res.buffer;
  *length = res.result.encoded;
  return RETURNok;
}

//------------------------------------------------------------------------------
static inline int s1ap_mme_encode_unsuccessfull_outcome(
    S1ap_S1AP_PDU_t* pdu, uint8_t** buffer, uint32_t* length) {
  asn_encode_to_new_buffer_result_t res = {NULL, {0, NULL, NULL}};

  if (pdu == NULL) {
    OAILOG_ERROR(LOG_S1AP, "PDU is NULL\n");
    return RETURNerror;
  }

  switch (pdu->choice.unsuccessfulOutcome.procedureCode) {
    case S1ap_ProcedureCode_id_S1Setup:
    case S1ap_ProcedureCode_id_PathSwitchRequest:
    case S1ap_ProcedureCode_id_HandoverPreparation:
      break;

    default:
      OAILOG_DEBUG(
          LOG_S1AP,
          "Unknown procedure ID (%d) for unsuccessfull outcome message\n",
          (int) pdu->choice.unsuccessfulOutcome.procedureCode);
      *buffer = NULL;
      *length = 0;
      return RETURNerror;
  }
  res = asn_encode_to_new_buffer(
      NULL, ATS_ALIGNED_CANONICAL_PER, &asn_DEF_S1ap_S1AP_PDU, pdu);
  *buffer = res.buffer;
  *length = res.result.encoded;
  return RETURNok;
}
