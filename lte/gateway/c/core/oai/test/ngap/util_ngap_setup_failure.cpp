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

int encode_setup_failure_pdu(Ngap_NGAP_PDU_t* pdu, uint8_t** buffer,
                             uint32_t* length) {
  asn_encode_to_new_buffer_result_t res = {NULL, {0, NULL, NULL}};

  res = asn_encode_to_new_buffer(NULL, ATS_ALIGNED_CANONICAL_PER,
                                 &asn_DEF_Ngap_NGAP_PDU, pdu);

  *buffer = (unsigned char*)res.buffer;
  *length = res.result.encoded;

  ASN_STRUCT_FREE_CONTENTS_ONLY(asn_DEF_Ngap_NGAP_PDU, pdu);

  return (0);
}

/*
 * Failure cause and type. Return the NGAPPDU Failure
 */
int ngap_ng_setup_failure_stream(const Ngap_Cause_PR cause_type,
                                 const long cause_value, bstring& stream) {
  uint8_t* buffer_p;
  uint32_t length = 0;
  Ngap_NGAP_PDU_t pdu;
  Ngap_NGSetupFailure_t* out;
  Ngap_NGSetupFailureIEs_t* ie = NULL;
  Ngap_Cause_t* cause_p = NULL;

  memset(&pdu, 0, sizeof(pdu));
  pdu.present = Ngap_NGAP_PDU_PR_unsuccessfulOutcome;
  pdu.choice.unsuccessfulOutcome.procedureCode = Ngap_ProcedureCode_id_NGSetup;
  pdu.choice.unsuccessfulOutcome.criticality = Ngap_Criticality_reject;
  pdu.choice.unsuccessfulOutcome.value.present =
      Ngap_UnsuccessfulOutcome__value_PR_NGSetupFailure;

  out = &pdu.choice.unsuccessfulOutcome.value.choice.NGSetupFailure;

  ie = (Ngap_NGSetupFailureIEs_t*)calloc(1, sizeof(Ngap_NGSetupFailureIEs_t));
  ie->id = Ngap_ProtocolIE_ID_id_Cause;
  ie->criticality = Ngap_Criticality_ignore;
  ie->value.present = Ngap_NGSetupFailureIEs__value_PR_Cause;

  cause_p = &ie->value.choice.Cause;
  cause_p->present = cause_type;
  cause_p->choice.nas = cause_value;

  ASN_SEQUENCE_ADD(&out->protocolIEs, ie);

  if (encode_setup_failure_pdu(&pdu, &buffer_p, &length) < 0) {
    return (EXIT_FAILURE);
  }

  stream = blk2bstr(buffer_p, length);
  free(buffer_p);

  return (EXIT_SUCCESS);
}

int ngap_ng_setup_failure_pdu(const Ngap_Cause_PR cause_type,
                              const long cause_value,
                              Ngap_NGAP_PDU_t& encode_pdu) {
  Ngap_NGSetupFailure_t* out;
  Ngap_NGSetupFailureIEs_t* ie = NULL;
  Ngap_Cause_t* cause_p = NULL;

  encode_pdu.present = Ngap_NGAP_PDU_PR_unsuccessfulOutcome;
  encode_pdu.choice.unsuccessfulOutcome.procedureCode =
      Ngap_ProcedureCode_id_NGSetup;
  encode_pdu.choice.unsuccessfulOutcome.criticality = Ngap_Criticality_reject;
  encode_pdu.choice.unsuccessfulOutcome.value.present =
      Ngap_UnsuccessfulOutcome__value_PR_NGSetupFailure;

  out = &encode_pdu.choice.unsuccessfulOutcome.value.choice.NGSetupFailure;

  ie = (Ngap_NGSetupFailureIEs_t*)calloc(1, sizeof(Ngap_NGSetupFailureIEs_t));
  ie->id = Ngap_ProtocolIE_ID_id_Cause;
  ie->criticality = Ngap_Criticality_ignore;
  ie->value.present = Ngap_NGSetupFailureIEs__value_PR_Cause;

  cause_p = &ie->value.choice.Cause;
  cause_p->present = cause_type;
  cause_p->choice.nas = cause_value;

  ASN_SEQUENCE_ADD(&out->protocolIEs, ie);

  return (EXIT_SUCCESS);
}

bool ng_setup_failure_decode(const_bstring const raw, Ngap_NGAP_PDU_t* pdu) {
  asn_dec_rval_t dec_ret;

  memset(pdu, 0, sizeof(Ngap_NGAP_PDU_t));

  dec_ret = aper_decode(NULL, &asn_DEF_Ngap_NGAP_PDU, (void**)&pdu, bdata(raw),
                        blength(raw), 0, 0);
  if (dec_ret.code != RC_OK) {
    return false;
  }

  return true;
}
