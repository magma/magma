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

/* using this stub code we are going to test Decoding functionality of
 * De-Registration Request Message */

#include <iostream>
#include "M5GDeRegistrationRequestUEInit.h"
#include "M5GCommonDefs.h"

using namespace std;
using namespace magma5g;
namespace magma5g {
int decode(void) {
  int ret = 0;
  // Message to be Decoded
  uint8_t buffer[] = {0x7E, 0x00, 0x45, 0x01, 0x00, 0x0B, 0xF2, 0x13, 0x00,
                      0x14, 0x44, 0x33, 0x12, 0x00, 0x00, 0x00, 0x01};
  int len          = 17;
  DeRegistrationRequestUEInitMsg De_Req;
  MLOG(MDEBUG) << "\n\n---Decoding De-registration request (UE originating) "
                  "Message---\n\n";
  ret = De_Req.DecodeDeRegistrationRequestUEInitMsg(&De_Req, buffer, len);

  MLOG(BEBUG) << " ---DECODED MESSAGE ---\n";
  MLOG(MDEBUG) << " Extended Protocol Discriminator :" << dec
               << int(De_Req.extended_protocol_discriminator
                          .extended_proto_discriminator);
  MLOG(MDEBUG) << " Spare half octet : " << dec
               << int(De_Req.spare_half_octet.spare);
  MLOG(MDEBUG) << " Security Header Type : " << dec
               << int(De_Req.sec_header_type.sec_hdr);
  MLOG(MDEBUG) << " Message Type : 0x" << hex
               << int(De_Req.message_type.msg_type);
  MLOG(MDEBUG) << " M5GS De-Registration Type :";
  MLOG(DEBUG) << "   Switch off = " << dec
              << int(De_Req.m5gs_de_reg_type.switchoff);
  MLOG(MDEBUG) << "   Re-registration required = " << dec
               << int(De_Req.m5gs_de_reg_type.re_reg_required);
  MLOG(MDEBUG) << "   Access Type = " << dec
               << int(De_Req.m5gs_de_reg_type.access_type);
  MLOG(MDEBUG) << " NAS key set identifier : ";
  MLOG(MDEBUG) << "   Type of security context flag = " << dec
               << int(De_Req.nas_key_set_identifier.tsc);
  MLOG(MDEBUG) << "   NAS key set identifier = " << dec
               << int(De_Req.nas_key_set_identifier.nas_key_set_identifier);
  MLOG(MDEBUG) << " M5GS mobile identity : ";
  MLOG(MDEBUG)
      << "   Odd/even Indication = " << dec
      << int(De_Req.m5gs_mobile_identity.mobile_identity.guti.odd_even);
  MLOG(MDEBUG)
      << "   Type of identity = " << dec
      << int(De_Req.m5gs_mobile_identity.mobile_identity.guti.type_of_identity);
  MLOG(MDEBUG)
      << "   Mobile Country Code (MCC) = " << dec
      << int(De_Req.m5gs_mobile_identity.mobile_identity.guti.mcc_digit1) << dec
      << int(De_Req.m5gs_mobile_identity.mobile_identity.guti.mcc_digit2) << dec
      << int(De_Req.m5gs_mobile_identity.mobile_identity.guti.mcc_digit3);
  MLOG(MDEBUG)
      << "   Mobile NetWork Code (MNC) = " << dec
      << int(De_Req.m5gs_mobile_identity.mobile_identity.guti.mnc_digit1) << dec
      << int(De_Req.m5gs_mobile_identity.mobile_identity.guti.mnc_digit2) << dec
      << int(De_Req.m5gs_mobile_identity.mobile_identity.guti.mnc_digit3);
  MLOG(MDEBUG)
      << " AMF Region ID = " << dec
      << int(De_Req.m5gs_mobile_identity.mobile_identity.guti.amf_regionid);
  MLOG(MDEBUG)
      << " AMF Set ID = " << dec
      << int(De_Req.m5gs_mobile_identity.mobile_identity.guti.amf_setid);
  MLOG(MDEBUG)
      << " AMF Pointer = " << dec
      << int(De_Req.m5gs_mobile_identity.mobile_identity.guti.amf_pointer);
  MLOG(MDEBUG) << " 5G-TMSI = 0x0" << hex
               << int(De_Req.m5gs_mobile_identity.mobile_identity.guti.tmsi1)
               << "0" << hex
               << int(De_Req.m5gs_mobile_identity.mobile_identity.guti.tmsi2)
               << "0" << hex
               << int(De_Req.m5gs_mobile_identity.mobile_identity.guti.tmsi3)
               << "0" << hex
               << int(De_Req.m5gs_mobile_identity.mobile_identity.guti.tmsi4)
               << "\n\n";

  return 0;
}
}  // namespace magma5g

// Main Function to call test decode function
int main(void) {
  int ret;
  ret = magma5g::decode();
  return 0;
}
