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

/*! \file ngap_common.c
   \brief ngap procedures for both eNB and MME
   \author Sebastien ROUX <sebastien.roux@eurecom.fr>
   \date 2012
   \version 0.1
*/

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
    uint8_t** buffer, uint32_t* length,Ngap_ProcedureCode_t procedureCode,
    Ngap_Criticality_t criticality, asn_TYPE_descriptor_t* td, void* sptr) {
  Ngap_NGAP_PDU_t pdu;
  ssize_t encoded;

  memset(&pdu, 0, sizeof(Ngap_NGAP_PDU_t));
  pdu.present                                = Ngap_NGAP_PDU_PR_initiatingMessage;
  pdu.choice.initiatingMessage.procedureCode = procedureCode;
  pdu.choice.initiatingMessage.criticality   = criticality;
  //ANY_fromType_aper(&(pdu.choice.initiatingMessage.value), td, sptr);

  if (asn1_xer_print) {
    xer_fprint(stdout, &asn_DEF_Ngap_PDUSessionType, (void*) &pdu);
  }

  /*
   * We can safely free list of IE from sptr
   */
  ASN_STRUCT_FREE_CONTENTS_ONLY(*td, sptr);

  //if ((encoded = aper_encode_to_new_buffer(
  if ((encoded = aper_encode_to_new_buffer(
           &asn_DEF_Ngap_PDUSessionType, 0, &pdu, (void**) buffer)) < 0) {
    OAILOG_ERROR(LOG_NGAP, "Encoding of %s failed\n", td->name);
    return -1;
  }

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_PDUSessionType, &pdu);

  *length = encoded;
  return encoded;
}

ssize_t ngap_generate_successfull_outcome(
    uint8_t** buffer, uint32_t* length, Ngap_ProcedureCode_t procedureCode,
    Ngap_Criticality_t criticality, asn_TYPE_descriptor_t* td, void* sptr) {
  Ngap_NGAP_PDU_t pdu;
  ssize_t encoded;

  memset(&pdu, 0, sizeof(Ngap_NGAP_PDU_t));
  pdu.present                                = Ngap_NGAP_PDU_PR_successfulOutcome;
  pdu.choice.successfulOutcome.procedureCode = procedureCode;
  pdu.choice.successfulOutcome.criticality   = criticality;
  //ANY_fromType_aper(&pdu.choice.successfulOutcome.value, td, sptr);

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

ssize_t ngap_generate_unsuccessfull_outcome(
    uint8_t** buffer, uint32_t* length, Ngap_ProcedureCode_t procedureCode,
    Ngap_Criticality_t criticality, asn_TYPE_descriptor_t* td, void* sptr) {
  Ngap_NGAP_PDU_t pdu;
  ssize_t encoded;

  memset(&pdu, 0, sizeof(Ngap_NGAP_PDU_t));
  pdu.present = Ngap_NGAP_PDU_PR_unsuccessfulOutcome;
  pdu.choice.unsuccessfulOutcome.procedureCode = procedureCode;
  pdu.choice.unsuccessfulOutcome.criticality   = criticality;
  //ANY_fromType_aper(pdu.choice.unsuccessfulOutcome.value, td, sptr);

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

/* commented by chandhu
Ngap_IE_t* ngap_new_ie(
    Ngap_ProtocolIE_ID_t id, Ngap_Criticality_t criticality,
    asn_TYPE_descriptor_t* type, void* sptr) {
  Ngap_IE_t* buff;

  if ((buff = malloc(sizeof(Ngap_IE_t))) == NULL) {
    // Possible error on malloc
    return NULL;
  }

  memset((void*) buff, 0, sizeof(Ngap_IE_t));
  buff->id          = id;
  buff->criticality = criticality;

  if (ANY_fromType_aper(&buff->value, type, sptr) < 0) {
    OAILOG_ERROR(LOG_NGAP, "Encoding of %s failed\n", type->name);
    free_wrapper((void**) &buff);
    return NULL;
  }

  if (asn1_xer_print)
    if (xer_fprint(stdout, &asn_DEF_Ngap_IE, buff) < 0) {
      free_wrapper((void**) &buff);
      return NULL;
    }

  return buff;
}
*/

// TODO: (amar) Unused function check with OAI
void ngap_handle_criticality(Ngap_Criticality_t criticality) {}
