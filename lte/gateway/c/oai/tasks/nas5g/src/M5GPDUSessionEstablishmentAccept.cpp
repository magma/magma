/*
   Copyright 2020 The Magma Authors.
   This source code is licensed under the BSD-style license found in the
   LICENSE file in the root directory of this source tree.
   Unless required by applicable law or agreed to in writing, software
   distributed under the License is distributed on an "AS IS" BASIS,
   WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
   See the License for the specific language governing permissions and
   limitations under the License.
 */

#include <sstream>
#include "M5GPDUSessionEstablishmentAccept.h"
#include "M5GCommonDefs.h"
#include "M5gNasMessage.h"

namespace magma5g {
PDUSessionEstablishmentAcceptMsg::PDUSessionEstablishmentAcceptMsg(){};
PDUSessionEstablishmentAcceptMsg::~PDUSessionEstablishmentAcceptMsg(){};

// Decode PDUSessionEstablishmentAccept Message and its IEs
int PDUSessionEstablishmentAcceptMsg::DecodePDUSessionEstablishmentAcceptMsg(
    PDUSessionEstablishmentAcceptMsg* pdu_session_estab_accept, uint8_t* buffer,
    uint32_t len) {
  // Not yet implemented, will be supported POST MVC
  return 0;
}

// Encode PDUSessionEstablishmentAccept Message and its IEs
int PDUSessionEstablishmentAcceptMsg::EncodePDUSessionEstablishmentAcceptMsg(
    PDUSessionEstablishmentAcceptMsg* pdu_session_estab_accept, uint8_t* buffer,
    uint32_t len) {
  uint32_t encoded        = 0;
  uint32_t encoded_result = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, PDU_SESSION_ESTABLISH_ACPT_MIN_LEN, len);

  MLOG(MDEBUG) << "EncodePDUSessionEstablishmentAcceptMsg : \n";
  if ((encoded_result =
           pdu_session_estab_accept->extended_protocol_discriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &pdu_session_estab_accept->extended_protocol_discriminator,
                   0, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = pdu_session_estab_accept->pdu_session_identity
                            .EncodePDUSessionIdentityMsg(
                                &pdu_session_estab_accept->pdu_session_identity,
                                0, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = pdu_session_estab_accept->pti.EncodePTIMsg(
           &pdu_session_estab_accept->pti, 0, buffer + encoded,
           len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result =
           pdu_session_estab_accept->message_type.EncodeMessageTypeMsg(
               &pdu_session_estab_accept->message_type, 0, buffer + encoded,
               len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = pdu_session_estab_accept->ssc_mode.EncodeSSCModeMsg(
           &pdu_session_estab_accept->ssc_mode, 0, buffer + encoded,
           len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result =
           pdu_session_estab_accept->pdu_session_type.EncodePDUSessionTypeMsg(
               &pdu_session_estab_accept->pdu_session_type, 0, buffer + encoded,
               len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = pdu_session_estab_accept->qos_rules.EncodeQOSRulesMsg(
           &pdu_session_estab_accept->qos_rules, 0, buffer + encoded,
           len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result =
           pdu_session_estab_accept->session_ambr.EncodeSessionAMBRMsg(
               &pdu_session_estab_accept->session_ambr, 0, buffer + encoded,
               len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result =
           pdu_session_estab_accept->pdu_address.EncodePDUAddressMsg(
               &pdu_session_estab_accept->pdu_address, PDU_ADDRESS,
               buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = pdu_session_estab_accept->dnn.EncodeDNNMsg(
           &pdu_session_estab_accept->dnn, DNN, buffer + encoded,
           len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  return encoded;
}
}  // namespace magma5g
