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
  Source      ngap_amf_encoder.c
  Date        2020/07/28
  Subsystem   Access and Mobility Management Function
  Author      Ashish Prajapati
  Description Defines NG Application Protocol Messages

*****************************************************************************/

#include <stdint.h>
#include <string.h>

#include "ngap_amf_encoder.h"
#include "ngap_common.h"
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
static inline int ngap_amf_encode_successful_outcome(
    Ngap_NGAP_PDU_t* pdu, uint8_t** buffer, uint32_t* len);
static inline int ngap_amf_encode_unsuccessful_outcome(
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
      ret = ngap_amf_encode_initiating(pdu, buffer, length);
      break;

    case Ngap_NGAP_PDU_PR_successfulOutcome:
      ret = ngap_amf_encode_successful_outcome(pdu, buffer, length);
      break;

    case Ngap_NGAP_PDU_PR_unsuccessfulOutcome:
      ret = ngap_amf_encode_unsuccessful_outcome(pdu, buffer, length);
      break;

    default:
      OAILOG_NOTICE(
          LOG_NGAP, "Unknown message outcome (%d) or not implemented",
          (int) pdu->present);
      break;
  }
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, pdu);
  return ret;
}

//------------------------------------------------------------------------------
static inline int ngap_amf_encode_initiating(
    Ngap_NGAP_PDU_t* pdu, uint8_t** buffer, uint32_t* length) {
  asn_encode_to_new_buffer_result_t res = {NULL, {0, NULL, NULL}};
  DevAssert(pdu != NULL);

  switch (pdu->choice.initiatingMessage.procedureCode) {
    case Ngap_ProcedureCode_id_DownlinkNASTransport:
    case Ngap_ProcedureCode_id_InitialContextSetup:
    case Ngap_ProcedureCode_id_UEContextRelease:
    case Ngap_ProcedureCode_id_Paging:
      break;

    default:
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
  return 0;
}

//------------------------------------------------------------------------------
static inline int ngap_amf_encode_successful_outcome(
    Ngap_NGAP_PDU_t* pdu, uint8_t** buffer, uint32_t* length) {
  asn_encode_to_new_buffer_result_t res = {NULL, {0, NULL, NULL}};
  DevAssert(pdu != NULL);

  switch (pdu->choice.successfulOutcome.procedureCode) {
    case Ngap_ProcedureCode_id_NGSetup:
      break;

    default:
      OAILOG_DEBUG(
          LOG_NGAP,
          "Unknown procedure ID (%d) for successful outcome message\n",
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
  return 0;
}

//------------------------------------------------------------------------------
static inline int ngap_amf_encode_unsuccessful_outcome(
    Ngap_NGAP_PDU_t* pdu, uint8_t** buffer, uint32_t* length) {
  asn_encode_to_new_buffer_result_t res = {NULL, {0, NULL, NULL}};
  DevAssert(pdu != NULL);

  switch (pdu->choice.unsuccessfulOutcome.procedureCode) {
    case Ngap_ProcedureCode_id_NGSetup:
      break;

    default:
      OAILOG_DEBUG(
          LOG_NGAP,
          "Unknown procedure ID (%d) for unsuccessful outcome message\n",
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
