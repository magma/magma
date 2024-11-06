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

#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GPDUSessionModificationRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"

namespace magma5g {
PDUSessionModificationRequestMsg::PDUSessionModificationRequestMsg(){};
PDUSessionModificationRequestMsg::~PDUSessionModificationRequestMsg(){};

// Decode PDUSessionModificationRequest Message and its IEs
int PDUSessionModificationRequestMsg::DecodePDUSessionModificationRequestMsg(
    uint8_t* buffer, uint32_t len) {
  uint32_t decoded = 0;
  int decoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, PDU_SESSION_MODIFICATION_REQ_MIN_LEN, len);

  OAILOG_DEBUG(LOG_NAS5G, "DecodePDUSessionModificationRequestMessage");

  if ((decoded_result = this->extended_protocol_discriminator
                            .DecodeExtendedProtocolDiscriminatorMsg(
                                0, buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result = this->pdu_session_identity.DecodePDUSessionIdentityMsg(
           0, buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result =
           this->pti.DecodePTIMsg(0, buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  if ((decoded_result = this->message_type.DecodeMessageTypeMsg(
           0, buffer + decoded, len - decoded)) < 0) {
    return decoded_result;
  } else {
    decoded += decoded_result;
  }
  while (decoded < len) {
    uint8_t ie_type = *(buffer + decoded);

    switch (ie_type) {
      case M5GSM_CAUSE: {
        if ((decoded_result = this->cause.DecodeM5GSMCauseMsg(
                 M5GSM_CAUSE, buffer + decoded, len - decoded)) < 0) {
          return decoded_result;
        } else {
          decoded += decoded_result;
        }
      } break;
      case PDU_SESSION_QOS_RULES_IE_TYPE: {
        QOSRulesMsg qosRules;
        if ((decoded_result = qosRules.DecodeQOSRulesMsg(
                 PDU_SESSION_QOS_RULES_IE_TYPE, buffer + decoded,
                 len - decoded)) < 0) {
          return decoded_result;
        } else {
          this->authqosrules.push_back(qosRules);
          decoded += decoded_result;
        }
      } break;
      case PDU_SESSION_QOS_FLOW_DESC_IE_TYPE: {
        // iei + length
        decoded += 3;
        M5GQosFlowDescription flowDes;
        if ((decoded_result = flowDes.DecodeM5GQosFlowDescription(
                 buffer + decoded, len - decoded)) < 0) {
          return decoded_result;
        } else {
          this->authqosflowdescriptors.push_back(flowDes);
          decoded += decoded_result;
        }
      } break;
      case REQUEST_EXTENDED_PROTOCOL_CONFIGURATION_OPTIONS_TYPE: {
        if ((decoded_result =
                 this->extprotocolconfigurationoptions
                     .DecodeProtocolConfigurationOptions(
                         REQUEST_EXTENDED_PROTOCOL_CONFIGURATION_OPTIONS_TYPE,
                         buffer + decoded, len - decoded)) < 0) {
          return decoded_result;
        } else {
          decoded += decoded_result;
        }
      } break;
      default: {
      } break;
    }
  }
  return decoded;
}

// Encode PDUSessionModificationRequest Message and its IEs
int PDUSessionModificationRequestMsg::EncodePDUSessionModificationRequestMsg(
    uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;
  uint32_t encoded_result = 0;
  uint16_t qos_flow_des_encoded = 0;
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, PDU_SESSION_MODIFICATION_REQ_MIN_LEN, len);

  OAILOG_DEBUG(LOG_NAS5G, "EncodePDUSessionModificationRequestMessage");
  if ((encoded_result = this->extended_protocol_discriminator
                            .EncodeExtendedProtocolDiscriminatorMsg(
                                0, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = this->pdu_session_identity.EncodePDUSessionIdentityMsg(
           0, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result =
           this->pti.EncodePTIMsg(0, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if ((encoded_result = this->message_type.EncodeMessageTypeMsg(
           0, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }
  if (this->cause.iei) {
    if ((encoded_result = this->cause.EncodeM5GSMCauseMsg(0, buffer + encoded,
                                                          len - encoded)) < 0) {
      return encoded_result;
    } else {
      encoded += encoded_result;
    }
  }
  for (uint8_t i = 0; i < this->authqosrules.size(); i++) {
    if ((encoded_result = this->authqosrules[i].EncodeQOSRulesMsg(
             PDU_SESSION_QOS_RULES_IE_TYPE, buffer + encoded, len - encoded)) <
        0) {
      return encoded_result;
    } else {
      encoded += encoded_result;
    }
  }
  this->authqosrules.clear();
  for (uint8_t i = 0; i < this->authqosflowdescriptors.size(); i++) {
    if ((encoded_result =
             this->authqosflowdescriptors[i].EncodeM5GQosFlowDescription(
                 buffer + encoded + 3, len - encoded)) < 0) {
      return encoded_result;
    } else {
      qos_flow_des_encoded += encoded_result;
    }
  }

  if (qos_flow_des_encoded) {
    // iei
    *(buffer + encoded++) = 0x79;
    IES_ENCODE_U16(buffer, encoded, qos_flow_des_encoded);
    encoded += qos_flow_des_encoded;
    this->authqosflowdescriptors.clear();
  }

  return encoded;
}

}  // namespace magma5g
