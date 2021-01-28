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

/* using this stub code we are going to test Encoding functionality of
 * Registration Accept Message */

#include <iostream>
#include <iomanip>
#include <M5GRegistrationAccept.h>
#include <AmfMessage.h>
#include <M5GCommonDefs.h>

using namespace std;
using namespace magma5g;
namespace magma5g {
int encode(void) {
  int enc_r          = 0;
  uint8_t buffer[19] = {};
  int len            = 19;
  RegistrationAcceptMsg msg;

  // Message to be Encoded
  msg.extended_protocol_discriminator.extended_proto_discriminator = 126;
  msg.sec_header_type.sec_hdr                                      = 0;
  msg.spare_half_octet.spare                                       = 0;
  msg.message_type.msg_type                                        = 0x42;
  msg.m5gs_reg_result.spare                                        = 0;
  msg.m5gs_reg_result.sms_allowed                                  = 0;
  msg.m5gs_reg_result.reg_result_val                               = 1;
  msg.mobile_id.mobile_identity.guti.odd_even                      = 0;
  msg.mobile_id.iei                                                = 0x77;
  msg.mobile_id.len                                                = 11;
  msg.mobile_id.mobile_identity.guti.type_of_identity              = 2;
  msg.mobile_id.mobile_identity.guti.mcc_digit1                    = 2;
  msg.mobile_id.mobile_identity.guti.mcc_digit2                    = 3;
  msg.mobile_id.mobile_identity.guti.mcc_digit3                    = 4;
  msg.mobile_id.mobile_identity.guti.mnc_digit1                    = 6;
  msg.mobile_id.mobile_identity.guti.mnc_digit2                    = 7;
  msg.mobile_id.mobile_identity.guti.mnc_digit3                    = 15;
  msg.mobile_id.mobile_identity.guti.amf_regionid                  = 68;
  msg.mobile_id.mobile_identity.guti.amf_setid                     = 204;
  msg.mobile_id.mobile_identity.guti.amf_pointer                   = 18;
  msg.mobile_id.mobile_identity.guti.tmsi1                         = 0;
  msg.mobile_id.mobile_identity.guti.tmsi2                         = 0;
  msg.mobile_id.mobile_identity.guti.tmsi3                         = 0;
  msg.mobile_id.mobile_identity.guti.tmsi4                         = 1;

  MLOG(MDEBUG) << "---Encoding Registration Accept Message---";
  // Encoding Message
  enc_r = msg.EncodeRegistrationAcceptMsg(&msg, buffer, len);

  MLOG(MDEBUG) << " Encoded Message : ";
  for (size_t i = 0; i < sizeof(buffer); i++) {
    MLOG(MDEBUG) << setfill('0') << hex << setw(2) << int(buffer[i]);
  }

  return 0;
}
}  // namespace magma5g

// Main Function to call test encode function
int main(void) {
  int ret;
  ret = magma5g::encode();
  return 0;
}
