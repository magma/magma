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
#include "M5GRegistrationRequest.h"
#include "M5GCommonDefs.h"

namespace magma5g {
RegistrationRequestMsg::RegistrationRequestMsg(){};
RegistrationRequestMsg::~RegistrationRequestMsg(){};

// Decode RegistrationRequest Message and its IEs
int RegistrationRequestMsg::DecodeRegistrationRequestMsg(
    RegistrationRequestMsg* reg_request, uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, REGISTRATION_REQUEST_MINIMUM_LENGTH, len);

  MLOG(MDEBUG) << "DecodeRegistrationRequestMsg : \n";
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
  if ((decoded_result =
           reg_request->ue_sec_capability.DecodeUESecurityCapabilityMsg(
               &reg_request->ue_sec_capability, 0x2e, buffer + decoded,
               len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

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
