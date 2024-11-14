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
#include <iomanip>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GPDUSessionEstablishmentRequest.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
PDUSessionEstablishmentRequestMsg::PDUSessionEstablishmentRequestMsg(){};
PDUSessionEstablishmentRequestMsg::~PDUSessionEstablishmentRequestMsg(){};

// Decode PDUSessionEstablishmentRequest Message and its IEs
int PDUSessionEstablishmentRequestMsg::DecodePDUSessionEstablishmentRequestMsg(
    uint8_t* buffer, uint32_t len) {
  uint32_t decoded = 0;
  int decoded_result = 0;
  uint8_t type_len = sizeof(uint8_t);
  uint8_t length_len = sizeof(uint8_t);

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(buffer,
                                       PDU_SESSION_ESTABLISH_REQ_MIN_LEN, len);

  if ((decoded_result = this->extended_protocol_discriminator
                            .DecodeExtendedProtocolDiscriminatorMsg(
                                0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = this->pdu_session_identity.DecodePDUSessionIdentityMsg(
           0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           this->pti.DecodePTIMsg(0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = this->message_type.DecodeMessageTypeMsg(
           0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           this->integrity_prot_max_data_rate.DecodeIntegrityProtMaxDataRateMsg(
               0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  while (decoded < len) {
    // Size is incremented for the unhandled types by 1 byte
    uint32_t type = *(buffer + decoded) >= 0x80 ? ((*(buffer + decoded)) & 0xf0)
                                                : (*(buffer + decoded));
    decoded_result = 0;

    switch (type) {
      case REQUEST_PDU_SESSION_TYPE_TYPE:
        if ((decoded_result = this->pdu_session_type.DecodePDUSessionTypeMsg(
                 REQUEST_PDU_SESSION_TYPE_TYPE, buffer + decoded,
                 len - decoded)) < 0) {
          return decoded_result;
        } else {
          decoded += decoded_result;
        }
        break;

      case REQUEST_SSC_MODE_TYPE:
        if ((decoded_result = this->ssc_mode.DecodeSSCModeMsg(
                 REQUEST_SSC_MODE_TYPE, buffer + decoded, len - decoded)) < 0) {
          return decoded_result;
        } else {
          decoded += decoded_result;
        }
        break;
      case REQUEST_EXTENDED_PROTOCOL_CONFIGURATION_OPTIONS_TYPE:
        if ((decoded_result =
                 this->protocolconfigurationoptions
                     .DecodeProtocolConfigurationOptions(
                         REQUEST_EXTENDED_PROTOCOL_CONFIGURATION_OPTIONS_TYPE,
                         buffer + decoded, len - decoded)) < 0) {
          return decoded_result;
        } else {
          decoded += decoded_result;
        }
        break;
      case MAXIMUM_NUMBER_OF_SUPPORTED_PACKET_FILTERS_TYPE:
        if ((decoded_result =
                 this->maxNumOfSuppPacketFilters
                     .DecodeMaxNumOfSupportedPacketFilters(
                         MAXIMUM_NUMBER_OF_SUPPORTED_PACKET_FILTERS_TYPE,
                         buffer + decoded, len - decoded)) < 0) {
          return decoded_result;
        } else {
          decoded += decoded_result;
        }
        break;
      case REQUEST_5GSM_CAPABILITY_TYPE:
      case REQUEST_ALWAYS_ON_PDU_SESSION_REQUESTED_TYPE:
      case REQUEST_SM_PDU_DN_REQUEST_CONTAINER_TYPE:
      case REQUEST_HEADER_COMPRESSION_CONFIGURATION_TYPE:
      case REQUEST_DS_TT_ETHERNET_PORT_MAC_ADDRESS_TYPE:
      case REQUEST_UE_DS_TT_RESIDENCE_TIME_TYPE:
      case REQUEST_PORT_MANAGEMENT_INFORMATION_CONTAINER_TYPE:

        // TLV Types. 1 byte for Type and 1 Byte for size
        type_len = sizeof(uint8_t);
        length_len = sizeof(uint8_t);
        DECODE_U8(buffer + decoded + type_len, decoded_result, decoded);

        decoded += (length_len + decoded_result);
        break;

      default:
        decoded_result = -1;
        break;
    }

    if (decoded_result < 0) {
      return decoded_result;
    }
  }

  return decoded;
}

// Encode PDUSessionEstablishmentRequest Message and its IEs
int PDUSessionEstablishmentRequestMsg::EncodePDUSessionEstablishmentRequestMsg(
    uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;
  uint32_t encoded_result = 0;

  if (!buffer || (0 == len)) {
    OAILOG_ERROR(LOG_NAS5G, "Input arguments are not valid");
    return -1;
  }

  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer,
                                       PDU_SESSION_ESTABLISH_REQ_MIN_LEN, len);

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
  if ((encoded_result =
           this->integrity_prot_max_data_rate.EncodeIntegrityProtMaxDataRateMsg(
               0, buffer + encoded, len - encoded)) < 0) {
    return encoded_result;
  } else {
    encoded += encoded_result;
  }

  if (static_cast<uint32_t>(this->pdu_session_type.type_val)) {
    if ((encoded_result = this->pdu_session_type.EncodePDUSessionTypeMsg(
             REQUEST_PDU_SESSION_TYPE_TYPE, buffer + encoded, len - encoded)) <
        0) {
      return encoded_result;
    } else {
      encoded += encoded_result;
    }
  }

  if (static_cast<uint32_t>(this->ssc_mode.mode_val)) {
    if ((encoded_result = this->ssc_mode.EncodeSSCModeMsg(
             REQUEST_SSC_MODE_TYPE, buffer + encoded, len - encoded)) < 0) {
      return encoded_result;
    } else {
      encoded += encoded_result;
    }
  }
  return encoded;
}

}  // namespace magma5g
