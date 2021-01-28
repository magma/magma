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
 * Authentication Result Message */
#include <iostream>
#include <iomanip>
#include "M5GAuthenticationResult.h"
#include "M5GCommonDefs.h"
using namespace std;
using namespace magma5g;

namespace magma5g {
int encode(void) {
  int ret            = 0;
  uint8_t buffer[13] = {};
  int len            = 13;
  AuthenticationResultMsg msg;
  // Message to be Encoded
  msg.extended_protocol_discriminator.extended_proto_discriminator = 126;
  msg.spare_half_octet.spare                                       = 0;
  msg.sec_header_type.sec_hdr                                      = 0;
  msg.message_type.msg_type                                        = 0x46;
  msg.nas_key_set_identifier.tsc                                   = 0;
  msg.nas_key_set_identifier.nas_key_set_identifier                = 0;
  msg.eap_message.len                                              = 7;
  uint8_t eap_buff[] = {0x71, 0x00, 0x0d, 0x01, 0x77, 0x93, 0x11};
  msg.eap_message.eap.assign((const char*) eap_buff, 7);

  MLOG(MDEBUG) << "---Encoding message--- \n";
  ret = msg.EncodeAuthenticationResultMsg(&msg, buffer, len);

  MLOG(MDEBUG) << " ENCODED MESSAGE : ";
  for (int i = 0; i < 13; i++) {
    MLOG(MDEBUG) << setfill('0') << hex << setw(2) << int(buffer[i]);
  }
}
}  // namespace magma5g

int main(void) {
  int ret;
  ret = magma5g::encode();
  return 0;
}
