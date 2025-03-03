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
#include "lte/gateway/c/core/common/common_defs.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GRegistrationRequest.hpp"

namespace magma5g {
RegistrationRequestMsg::RegistrationRequestMsg() {};
RegistrationRequestMsg::~RegistrationRequestMsg() {};

// Decode RegistrationRequest Message and its IEs
int RegistrationRequestMsg::DecodeRegistrationRequestMsg(
    RegistrationRequestMsg* reg_request, uint8_t* buffer, uint32_t len) {
  uint32_t decoded = 0;
  int decoded_result = 0;
  uint8_t type_len = 0;
  uint8_t length_len = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, REGISTRATION_REQUEST_MINIMUM_LENGTH, len);

  if ((decoded_result = reg_request->extended_protocol_discriminator
                            .DecodeExtendedProtocolDiscriminatorMsg(
                                &reg_request->extended_protocol_discriminator,
                                0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = reg_request->spare_half_octet.DecodeSpareHalfOctetMsg(
           &reg_request->spare_half_octet, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           reg_request->sec_header_type.DecodeSecurityHeaderTypeMsg(
               &reg_request->sec_header_type, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result = reg_request->message_type.DecodeMessageTypeMsg(
           &reg_request->message_type, 0, buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           reg_request->m5gs_reg_type.DecodeM5GSRegistrationTypeMsg(
               &reg_request->m5gs_reg_type, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           reg_request->nas_key_set_identifier.DecodeNASKeySetIdentifierMsg(
               &reg_request->nas_key_set_identifier, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;
  if ((decoded_result =
           reg_request->m5gs_mobile_identity.DecodeM5GSMobileIdentityMsg(
               &reg_request->m5gs_mobile_identity, 0, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  while (decoded < len) {
    // Size is incremented for the unhandled types by 1 byte
    uint32_t type = *(buffer + decoded) >= 0x80 ? ((*(buffer + decoded)) & 0xf0)
                                                : (*(buffer + decoded));
    decoded_result = 0;

    switch (type) {
      case REGISTRATION_REQUEST_UE_SECURITY_CAPABILITY_TYPE:
        decoded_result =
            reg_request->ue_sec_capability.DecodeUESecurityCapabilityMsg(
                &reg_request->ue_sec_capability,
                REGISTRATION_REQUEST_UE_SECURITY_CAPABILITY_TYPE,
                buffer + decoded, len - decoded);
        if (decoded_result < 0) {
          return decoded_result;
        }

        decoded += decoded_result;
        break;

      case REGISTRATION_REQUEST_5GMM_CAPABILITY_TYPE:
      case REGISTRATION_REQUEST_REQUESTED_NSSAI_TYPE:
      case REGISTRATION_REQUEST_LAST_VISITED_REGISTERED_TAI_TYPE:
      case REGISTRATION_REQUEST_S1_UE_NETWORK_CAPABILITY_TYPE:
      case REGISTRATION_REQUEST_UPLINK_DATA_STATUS_TYPE:
      case REGISTRATION_REQUEST_PDU_SESSION_STATUS_TYPE:
      case REGISTRATION_REQUEST_MICO_INDICATION_TYPE:
      case REGISTRATION_REQUEST_UE_STATUS_TYPE:
      case REGISTRATION_REQUEST_ADDITIONAL_GUTI_TYPE:
      case REGISTRATION_REQUEST_ALLOWED_PDU_SESSION_STATUS_TYPE:
      case REGISTRATION_REQUEST_UE_USAGE_SETTING_TYPE:
      case REGISTRATION_REQUEST_REQUESTED_DRX_PARAMETERS_TYPE:
      case REGISTRATION_REQUEST_LADN_INDICATION_TYPE:
      case REGISTRATION_REQUEST_PAYLOAD_CONTAINER_TYPE_TYPE:
      case REGISTRATION_REQUEST_PAYLOAD_CONTAINER_TYPE:
      case REGISTRATION_REQUEST_NETWORK_SLICING_INDICATION_TYPE:
      case REGISTRATION_REQUEST_5GS_UPDATE_TYPE_TYPE:
      case REGISTRATION_REQUEST_EPS_BEARER_CONTEXT_STATUS_TYPE:
        // TLV Types. 1 byte for Type and 1 Byte for size

        type_len = sizeof(uint8_t);
        length_len = sizeof(uint8_t);
        DECODE_U8(buffer + decoded + type_len, decoded_result, decoded);
        decoded += (length_len + decoded_result);
        break;

      case REGISTRATION_REQUEST_EPS_NAS_MESSAGE_CONTAINER_TYPE:
      case REGISTRATION_REQUEST_NAS_MESSAGE_CONTAINER_TYPE:
        // TLV Types. 1 byte for Type and 2 Byte for size
        type_len = sizeof(uint8_t);
        length_len = 2 * sizeof(uint8_t);
        DECODE_U16(buffer + decoded + type_len, decoded_result, decoded);
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

// Will be supported POST MVC
// Encode Registration Request Message and its IEs
int RegistrationRequestMsg::EncodeRegistrationRequestMsg(
    RegistrationRequestMsg* reg_request, uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;
  // Will be supported POST MVC
  return encoded;
}
}  // namespace magma5g
