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

/* using this stub code we are going to test Decoding functionality of Security
 * Mode Complete Message */

#include <iostream>
#include <M5GSecurityModeComplete.h>
#include <M5GCommonDefs.h>

using namespace std;
using namespace magma5g;

namespace magma5g {
// Testing Decoding functionality
int decode(void) {
  int ret = 0;

  // Message to be decoded
  uint8_t buffer[] = {0x7E, 0x00, 0x5E};
  int len          = 3;
  SecurityModeCompleteMsg msg;

  // Decoding Security Mode Complete Message
  ret = msg.DecodeSecurityModeCompleteMsg(&msg, buffer, len);

  // Decoded Message
  MLOG(MDEBUG) << " ---Decoded Message---\n";
  MLOG(MDEBUG)
      << " Extended Protocol Discriminator :" << dec
      << int(msg.extended_protocol_discriminator.extended_proto_discriminator);
  MLOG(MDEBUG) << " Spare Half Octet : " << dec
               << int(msg.spare_half_octet.spare);
  MLOG(MDEBUG) << " Security Header Type : " << dec
               << int(msg.sec_header_type.sec_hdr);
  MLOG(MDEBUG) << " Message Type : 0x" << hex << int(msg.message_type.msg_type);

  return 0;
}
}  // namespace magma5g

// Main function to call test decode function
int main(void) {
  int ret;
  ret = magma5g::decode();
  return 0;
}
