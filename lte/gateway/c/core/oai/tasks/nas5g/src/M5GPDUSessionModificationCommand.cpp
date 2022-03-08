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
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GPDUSessionModificationCommand.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"

namespace magma5g {
PDUSessionModificationCommand::PDUSessionModificationCommand() {}
PDUSessionModificationCommand::~PDUSessionModificationCommand() {}

int PDUSessionModificationCommand::EncodePDUSessionModificationCommand(
    PDUSessionModificationCommand* pdu_sess_mod_comd, uint8_t* buffer,
    uint32_t len) {
  uint32_t encoded = 0;
  uint32_t encoded_result = 0;
  uint16_t qos_flow_des_encoded = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, PDU_SESSION_MODIFICATION_COMMAND_MIN_LEN, len);

  OAILOG_DEBUG(LOG_NAS5G, "EncodePDUSessionModificationCommand");
  if ((encoded_result =
           pdu_sess_mod_comd->extended_protocol_discriminator
               .EncodeExtendedProtocolDiscriminatorMsg(
                   &pdu_sess_mod_comd->extended_protocol_discriminator, 0,
                   buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result =
           pdu_sess_mod_comd->pdu_session_identity.EncodePDUSessionIdentityMsg(
               &pdu_sess_mod_comd->pdu_session_identity, 0, buffer + encoded,
               len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = pdu_sess_mod_comd->pti.EncodePTIMsg(
           &pdu_sess_mod_comd->pti, 0, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = pdu_sess_mod_comd->message_type.EncodeMessageTypeMsg(
           &pdu_sess_mod_comd->message_type, 0, buffer + encoded,
           len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if (pdu_sess_mod_comd->cause.iei) {
    if ((encoded_result = pdu_sess_mod_comd->cause.EncodeM5GSMCauseMsg(
             &pdu_sess_mod_comd->cause, 0, buffer + encoded, len - encoded)) <
        0) {
      return encoded_result;
    } else {
      encoded += encoded_result;
    }
  }
  if (pdu_sess_mod_comd->sessionambr.iei) {
    if ((encoded_result = pdu_sess_mod_comd->sessionambr.EncodeSessionAMBRMsg(
             &pdu_sess_mod_comd->sessionambr, PDU_SESSION_AMBR_IE_TYPE,
             buffer + encoded, len - encoded)) < 0) {
      return encoded_result;
    } else {
      encoded += encoded_result;
    }
  }
  for (uint8_t i = 0; i < pdu_sess_mod_comd->authqosrules.size(); i++) {
    if ((encoded_result = pdu_sess_mod_comd->authqosrules[i].EncodeQOSRulesMsg(
             &pdu_sess_mod_comd->authqosrules[i],
             PDU_SESSION_AUTH_QOS_RULES_IE_TYPE, buffer + encoded,
             len - encoded)) < 0) {
      return encoded_result;
    } else {
      encoded += encoded_result;
    }
  }
  pdu_sess_mod_comd->authqosrules.clear();
  for (uint8_t i = 0; i < pdu_sess_mod_comd->authqosflowdescriptors.size();
       i++) {
    if ((encoded_result = pdu_sess_mod_comd->authqosflowdescriptors[i]
                              .EncodeM5GQosFlowDescription(
                                  &pdu_sess_mod_comd->authqosflowdescriptors[i],
                                  buffer + encoded + 3, len - encoded)) < 0) {
      return encoded_result;
    } else {
      qos_flow_des_encoded += encoded_result;
    }
  }

  if (qos_flow_des_encoded) {
    // iei
    *(buffer + encoded) = 0x79;
    encoded++;
    IES_ENCODE_U16(buffer, encoded, qos_flow_des_encoded);
    encoded += qos_flow_des_encoded;
    pdu_sess_mod_comd->authqosflowdescriptors.clear();
  }

  return encoded;
}

int PDUSessionModificationCommand::DecodePDUSessionModificationCommand(
    PDUSessionModificationCommand* pdu_sess_mod_comd, uint8_t* buffer,
    uint32_t len) {
  uint32_t decoded = 0;
  uint32_t decoded_result = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, PDU_SESSION_MODIFICATION_COMMAND_MIN_LEN, len);

  OAILOG_DEBUG(LOG_NAS5G, "DecodePDUSessionModificationCommandMsg");
  if ((decoded_result =
           pdu_sess_mod_comd->extended_protocol_discriminator
               .DecodeExtendedProtocolDiscriminatorMsg(
                   &pdu_sess_mod_comd->extended_protocol_discriminator, 0,
                   buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           pdu_sess_mod_comd->pdu_session_identity.DecodePDUSessionIdentityMsg(
               &pdu_sess_mod_comd->pdu_session_identity, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = pdu_sess_mod_comd->pti.DecodePTIMsg(
           &pdu_sess_mod_comd->pti, 0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = pdu_sess_mod_comd->message_type.DecodeMessageTypeMsg(
           &pdu_sess_mod_comd->message_type, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  while (decoded < len) {
    uint8_t ie_type = *(buffer + decoded);

    switch (ie_type) {
      case M5GSM_CAUSE: {
        if ((decoded_result = pdu_sess_mod_comd->cause.DecodeM5GSMCauseMsg(
                 &pdu_sess_mod_comd->cause, M5GSM_CAUSE, buffer + decoded,
                 len - decoded)) < 0)
          return decoded_result;
        else
          decoded += decoded_result;
      } break;
      case PDU_SESSION_AMBR_IE_TYPE: {
        if ((decoded_result =
                 pdu_sess_mod_comd->sessionambr.DecodeSessionAMBRMsg(
                     &pdu_sess_mod_comd->sessionambr, PDU_SESSION_AMBR_IE_TYPE,
                     buffer + decoded, len - decoded)) < 0)
          return decoded_result;
        else
          decoded += decoded_result;
      } break;
      case PDU_SESSION_AUTH_QOS_RULES_IE_TYPE: {
        QOSRulesMsg qosRules;
        if ((decoded_result = qosRules.DecodeQOSRulesMsg(
                 &qosRules, PDU_SESSION_AUTH_QOS_RULES_IE_TYPE,
                 buffer + decoded, len - decoded)) < 0)
          return decoded_result;
        else
          pdu_sess_mod_comd->authqosrules.push_back(qosRules);
        decoded += decoded_result;
      } break;
      case PDU_SESSION_AUTH_QOS_FLOW_DESC_IE_TYPE: {
        // iei + length
        decoded += 3;
        M5GQosFlowDescription flowDes;
        if ((decoded_result = flowDes.DecodeM5GQosFlowDescription(
                 &flowDes, buffer + decoded, len - decoded)) < 0)
          return decoded_result;
        else
          pdu_sess_mod_comd->authqosflowdescriptors.push_back(flowDes);
        decoded += decoded_result;
      } break;
      case REQUEST_EXTENDED_PROTOCOL_CONFIGURATION_OPTIONS_TYPE: {
        if ((decoded_result =
                 pdu_sess_mod_comd->extprotocolconfigurationoptions
                     .DecodeProtocolConfigurationOptions(
                         &pdu_sess_mod_comd->extprotocolconfigurationoptions,
                         REQUEST_EXTENDED_PROTOCOL_CONFIGURATION_OPTIONS_TYPE,
                         buffer + decoded, len - decoded)) < 0) {
          return decoded_result;
        } else {
          decoded += decoded_result;
        }
      } break;
      default: {
      }
    }
  }
  return decoded;
}
}  // namespace magma5g
