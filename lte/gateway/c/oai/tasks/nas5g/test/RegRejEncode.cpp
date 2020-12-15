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
 * Registration Reject Message */

#include <iostream>
#include <iomanip>
#include <M5GRegistrationReject.h>
#include <AmfMessage.h>
#include <M5GCommonDefs.h>

using namespace std;
using namespace magma5g;
namespace magma5g {
int encode(void) {
  int enc_r         = 0;
  uint8_t buffer[4] = {};
  int len           = 4;
  RegistrationRejectMsg msg;

  // Message to be Encoded
  msg.extended_protocol_discriminator.extended_proto_discriminator = 126;
  msg.sec_header_type.sec_hdr                                      = 0;
  msg.spare_half_octet.spare                                       = 0;
  msg.message_type.msg_type                                        = 0x44;
  msg.m5gmm_cause.m5gmm_cause                                      = 12;

  MLOG(MDEBUG) << "---Encoding Registration Reject Message---";
  // Encoding Message
  enc_r = msg.EncodeRegistrationRejectMsg(&msg, buffer, len);

  MLOG(MDEBUG) << " Encoded Message : " << setfill('0') << hex << int(buffer[0])
               << setw(2) << hex << int(buffer[1]) << hex << int(buffer[2])
               << setw(2) << hex << int(buffer[3]);

  MLOG(MDEBUG) << "---Decoding Encoded Registration Reject Message---";
  int dec_r = 0;
  // Decoding Message
  dec_r = msg.DecodeRegistrationRejectMsg(&msg, buffer, len);
  MLOG(BEBUG) << "\n\n ---DECODED MESSAGE ---\n\n";
  MLOG(MDEBUG)
      << " Extended Protocol Discriminator :" << dec
      << int(msg.extended_protocol_discriminator.extended_proto_discriminator)
      << endl;
  MLOG(MDEBUG) << " Spare half octet : " << dec
               << int(msg.spare_half_octet.spare) << endl;
  MLOG(MDEBUG) << " Security Header Type : " << dec
               << int(msg.sec_header_type.sec_hdr) << endl;
  MLOG(MDEBUG) << " Message Type : 0x" << hex << int(msg.message_type.msg_type)
               << endl;
  MLOG(MDEBUG) << " 5GMM Cause :" << dec << int(msg.m5gmm_cause.m5gmm_cause)
               << endl;

  return 0;
}
}  // namespace magma5g

// Main Function to call test encode function
int main(void) {
  int ret;
  ret = magma5g::encode();
  return 0;
}
