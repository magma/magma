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
  Source      ngap_common.c
  Date        2020/07/28
  Subsystem   Access and Mobility Management Function
  Author      Ashish Prajapati
  Description Defines NG Application Protocol Messages

*****************************************************************************/

#include <stdint.h>
#include <stdio.h>
#include <stdlib.h>
#include <string.h>

#include "ngap_common.h"
#include "dynamic_memory_check.h"
#include "log.h"
#include "ANY.h"
#include "Ngap_NGAP-PDU.h"
#include "Ngap_InitiatingMessage.h"
#include "Ngap_SuccessfulOutcome.h"
#include "Ngap_UnsuccessfulOutcome.h"
#include "per_encoder.h"
#include "xer_encoder.h"

int asn_debug      = 0;
int asn1_xer_print = 0;

ssize_t ngap_generate_initiating_message(
    uint8_t** buffer, uint32_t* length, Ngap_ProcedureCode_t procedureCode,
    Ngap_Criticality_t criticality, asn_TYPE_descriptor_t* td, void* sptr) {
  Ngap_NGAP_PDU_t pdu;
  ssize_t encoded;

  memset(&pdu, 0, sizeof(Ngap_NGAP_PDU_t));
  pdu.present = Ngap_NGAP_PDU_PR_initiatingMessage;
  pdu.choice.initiatingMessage.procedureCode = procedureCode;
  pdu.choice.initiatingMessage.criticality   = criticality;

  if (asn1_xer_print) {
    xer_fprint(stdout, &asn_DEF_Ngap_PDUSessionType, (void*) &pdu);
  }

  /*
   * We can safely free list of IE from sptr
   */
  ASN_STRUCT_FREE_CONTENTS_ONLY(*td, sptr);

  if ((encoded = aper_encode_to_new_buffer(
           &asn_DEF_Ngap_PDUSessionType, 0, &pdu, (void**) buffer)) < 0) {
    OAILOG_ERROR(LOG_NGAP, "Encoding of %s failed\n", td->name);
    return -1;
  }

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_PDUSessionType, &pdu);

  *length = encoded;
  return encoded;
}

ssize_t ngap_generate_successful_outcome(
    uint8_t** buffer, uint32_t* length, Ngap_ProcedureCode_t procedureCode,
    Ngap_Criticality_t criticality, asn_TYPE_descriptor_t* td, void* sptr) {
  Ngap_NGAP_PDU_t pdu;
  ssize_t encoded;

  memset(&pdu, 0, sizeof(Ngap_NGAP_PDU_t));
  pdu.present = Ngap_NGAP_PDU_PR_successfulOutcome;
  pdu.choice.successfulOutcome.procedureCode = procedureCode;
  pdu.choice.successfulOutcome.criticality   = criticality;
  // ANY_fromType_aper(&pdu.choice.successfulOutcome.value, td, sptr);

  if (asn1_xer_print) {
    xer_fprint(stdout, &asn_DEF_Ngap_PDUSessionType, (void*) &pdu);
  }

  /*
   * We can safely free list of IE from sptr
   */
  ASN_STRUCT_FREE_CONTENTS_ONLY(*td, sptr);

  if ((encoded = aper_encode_to_new_buffer(
           &asn_DEF_Ngap_PDUSessionType, 0, &pdu, (void**) buffer)) < 0) {
    OAILOG_ERROR(LOG_NGAP, "Encoding of %s failed\n", td->name);
    return -1;
  }

  // Might need this if there is a leak here
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_PDUSessionType, &pdu);

  *length = encoded;
  return encoded;
}

ssize_t ngap_generate_unsuccessful_outcome(
    uint8_t** buffer, uint32_t* length, Ngap_ProcedureCode_t procedureCode,
    Ngap_Criticality_t criticality, asn_TYPE_descriptor_t* td, void* sptr) {
  Ngap_NGAP_PDU_t pdu;
  ssize_t encoded;

  memset(&pdu, 0, sizeof(Ngap_NGAP_PDU_t));
  pdu.present = Ngap_NGAP_PDU_PR_unsuccessfulOutcome;
  pdu.choice.unsuccessfulOutcome.procedureCode = procedureCode;
  pdu.choice.unsuccessfulOutcome.criticality   = criticality;
  // ANY_fromType_aper(pdu.choice.unsuccessfulOutcome.value, td, sptr);

  if (asn1_xer_print) {
    xer_fprint(stdout, &asn_DEF_Ngap_PDUSessionType, (void*) &pdu);
  }

  /*
   * We can safely free list of IE from sptr
   */
  ASN_STRUCT_FREE_CONTENTS_ONLY(*td, sptr);

  if ((encoded = aper_encode_to_new_buffer(
           &asn_DEF_Ngap_PDUSessionType, 0, &pdu, (void**) buffer)) < 0) {
    OAILOG_ERROR(LOG_NGAP, "Encoding of %s failed\n", td->name);
    return -1;
  }

  // Might need this if there is a leak here
  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_PDUSessionType, &pdu);

  *length = encoded;
  return encoded;
}

void ngap_handle_criticality(Ngap_Criticality_t criticality) {}
