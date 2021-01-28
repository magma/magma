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

#include "ngap_common.h"
#include "ngap_amf_encoder.h"
#include "assertions.h"
#include "log.h"
#include "Ngap_NGAP-PDU.h"
#include "Ngap_Criticality.h"
#include "Ngap_DownlinkNASTransport.h"
#include "Ngap_InitialContextSetupRequest.h"
#include "Ngap_Paging.h"
#include "Ngap_ProcedureCode.h"
#include "Ngap_NGSetupFailure.h"
#include "Ngap_NGSetupResponse.h"
#include "Ngap_UEContextModificationRequest.h"
#include "Ngap_UEContextReleaseCommand.h"

static inline int ngap_amf_encode_initiating(
    Ngap_NGAP_PDU_t* pdu, uint8_t** buffer, uint32_t* length);
static inline int ngap_amf_encode_successfull_outcome(
    Ngap_NGAP_PDU_t* pdu, uint8_t** buffer, uint32_t* len);
static inline int ngap_amf_encode_unsuccessfull_outcome(
    Ngap_NGAP_PDU_t* pdu, uint8_t** buffer, uint32_t* len);
//------------------------------------------------------------------------------
int ngap_amf_encode_pdu(
    Ngap_NGAP_PDU_t* pdu, uint8_t** buffer, uint32_t* length) {
  int ret = -1;
  DevAssert(pdu != NULL);
  DevAssert(buffer != NULL);
  DevAssert(length != NULL);

  switch (pdu->present) {
    case Ngap_NGAP_PDU_PR_initiatingMessage:
      OAILOG_ERROR(LOG_NGAP, "1");
      ret = ngap_amf_encode_initiating(pdu, buffer, length);
      OAILOG_ERROR(LOG_NGAP, "####ACL_TAG");
      break;

    case Ngap_NGAP_PDU_PR_successfulOutcome:
      OAILOG_ERROR(LOG_NGAP, "2");
      ret = ngap_amf_encode_successfull_outcome(pdu, buffer, length);
      break;

    case Ngap_NGAP_PDU_PR_unsuccessfulOutcome:
      OAILOG_ERROR(LOG_NGAP, "3");
      ret = ngap_amf_encode_unsuccessfull_outcome(pdu, buffer, length);
      break;

    default:
      OAILOG_ERROR(LOG_NGAP, "4");
      OAILOG_NOTICE(
          LOG_NGAP, "Unknown message outcome (%d) or not implemented",
          (int) pdu->present);
      break;
  }
  OAILOG_ERROR(LOG_NGAP, "####ACL_TAG");
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, pdu);
  OAILOG_ERROR(LOG_NGAP, "####ACL_TAG");
  return ret;
}

//------------------------------------------------------------------------------
static inline int ngap_amf_encode_initiating(
    Ngap_NGAP_PDU_t* pdu, uint8_t** buffer, uint32_t* length) {
  asn_encode_to_new_buffer_result_t res = {NULL, {0, NULL, NULL}};

  OAILOG_ERROR(LOG_NGAP, "######ACL_TAG: %s, %d  ", __func__, __LINE__);
  DevAssert(pdu != NULL);
  OAILOG_ERROR(LOG_NGAP, "######ACL_TAG: %s, %d  ", __func__, __LINE__);

  //  int i;
  //  for (i = 0; i < sizeof(Ngap_NGAP_PDU_t); i++) {
  //    OAILOG_ERROR(LOG_NGAP, "%02x ", ((unsigned char*) pdu)[i]);
  //  }

  switch (pdu->choice.initiatingMessage.procedureCode) {
    case Ngap_ProcedureCode_id_DownlinkNASTransport:
    case Ngap_ProcedureCode_id_InitialContextSetup:
    case Ngap_ProcedureCode_id_UEContextRelease:
    case Ngap_ProcedureCode_id_HandoverResourceAllocation:
    case Ngap_ProcedureCode_id_Paging:
    case Ngap_ProcedureCode_id_PDUSessionResourceSetup:
    case Ngap_ProcedureCode_id_PDUSessionResourceRelease:
      OAILOG_ERROR(LOG_NGAP, "######ACL_TAG: %s, %d  ", __func__, __LINE__);

      break;

    default:
      OAILOG_ERROR(LOG_NGAP, "######ACL_TAG: %s, %d  ", __func__, __LINE__);

      OAILOG_NOTICE(
          LOG_NGAP, "Unknown procedure ID (%d) for initiating message_p\n",
          (int) pdu->choice.initiatingMessage.procedureCode);
      *buffer = NULL;
      *length = 0;
      return -1;
  }

  memset(&res, 0, sizeof(res));
  res = asn_encode_to_new_buffer(
      NULL, ATS_ALIGNED_CANONICAL_PER, &asn_DEF_Ngap_NGAP_PDU, pdu);
  *buffer = res.buffer;
  *length = res.result.encoded;

  OAILOG_ERROR(LOG_NGAP, "buf:%p", buffer);
  OAILOG_ERROR(LOG_NGAP, "*buf:%p", *buffer);
  OAILOG_ERROR(LOG_NGAP, "l:%d", *length);
  //  for (i = 0; i < *length; i++) {
  //    OAILOG_ERROR(LOG_NGAP, "%02x ", ((unsigned char*) res.buffer)[i]);
  //  }

  return 0;
}

//------------------------------------------------------------------------------
static inline int ngap_amf_encode_successfull_outcome(
    Ngap_NGAP_PDU_t* pdu, uint8_t** buffer, uint32_t* length) {
  asn_encode_to_new_buffer_result_t res = {NULL, {0, NULL, NULL}};
  DevAssert(pdu != NULL);

  switch (pdu->choice.successfulOutcome.procedureCode) {
    OAILOG_ERROR(LOG_NGAP, "5");
    case Ngap_ProcedureCode_id_NGSetup:
    case Ngap_ProcedureCode_id_PathSwitchRequest:
    case Ngap_ProcedureCode_id_HandoverPreparation:
    case Ngap_ProcedureCode_id_HandoverCancel:
      // case Ngap_ProcedureCode_id_Reset:
      break;

    default:
      OAILOG_ERROR(LOG_NGAP, "6");
      OAILOG_DEBUG(
          LOG_NGAP,
          "Unknown procedure ID (%d) for successfull outcome message\n",
          (int) pdu->choice.successfulOutcome.procedureCode);
      *buffer = NULL;
      *length = 0;
      return -1;
  }

  memset(&res, 0, sizeof(res));
  res = asn_encode_to_new_buffer(
      NULL, ATS_ALIGNED_CANONICAL_PER, &asn_DEF_Ngap_NGAP_PDU, pdu);

  *buffer = res.buffer;
  *length = res.result.encoded;
  OAILOG_ERROR(LOG_NGAP, "buf:%p", buffer);
  OAILOG_ERROR(LOG_NGAP, "*buf:%p", *buffer);
  OAILOG_ERROR(LOG_NGAP, "l:%d", *length);
  return 0;
}

//------------------------------------------------------------------------------
static inline int ngap_amf_encode_unsuccessfull_outcome(
    Ngap_NGAP_PDU_t* pdu, uint8_t** buffer, uint32_t* length) {
  asn_encode_to_new_buffer_result_t res = {NULL, {0, NULL, NULL}};
  DevAssert(pdu != NULL);

  switch (pdu->choice.unsuccessfulOutcome.procedureCode) {
    case Ngap_ProcedureCode_id_NGSetup:
    case Ngap_ProcedureCode_id_PathSwitchRequest:
    case Ngap_ProcedureCode_id_HandoverPreparation:
      break;

    default:
      OAILOG_DEBUG(
          LOG_NGAP,
          "Unknown procedure ID (%d) for unsuccessfull outcome message\n",
          (int) pdu->choice.unsuccessfulOutcome.procedureCode);
      *buffer = NULL;
      *length = 0;
      return -1;
  }

  memset(&res, 0, sizeof(res));
  res = asn_encode_to_new_buffer(
      NULL, ATS_ALIGNED_CANONICAL_PER, &asn_DEF_Ngap_NGAP_PDU, pdu);
  *buffer = res.buffer;
  *length = res.result.encoded;
  return 0;
}
