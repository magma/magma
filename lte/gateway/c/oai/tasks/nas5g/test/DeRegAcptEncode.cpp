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
 * De-Registration Accept Message */

#include <iostream>
#include <iomanip>
#include "M5GDeRegistrationAcceptUEInit.h"
#include "M5GCommonDefs.h"

using namespace std;
using namespace magma5g;
namespace magma5g {
int encode(void) {
  int ret           = 0;
  uint8_t buffer[3] = {};
  int len           = 5;
  DeRegistrationAcceptUEInitMsg msg;

  // Message to be Encoded
  msg.extended_protocol_discriminator.extended_proto_discriminator = 126;
  msg.sec_header_type.sec_hdr                                      = 0;
  msg.message_type.msg_type                                        = 0x46;

  MLOG(MDEBUG) << "---Encoding message--- \n";
  ret = msg.EncodeDeRegistrationAcceptUEInitMsg(&msg, buffer, len);

  MLOG(MDEBUG) << " ENCODED MESSAGE : " << setfill('0') << hex << int(buffer[0])
               << setw(2) << hex << int(buffer[1]) << hex << int(buffer[2])
               << "\n";

  MLOG(MDEBUG) << "---Decoding encoded message--- ";
  int ret2 = 0;
  ret2     = msg.DecodeDeRegistrationAcceptUEInitMsg(&msg, buffer, len);

  MLOG(BEBUG) << " ---DECODED MESSAGE ---\n";
  MLOG(MDEBUG)
      << " Extended Protocol Discriminator :" << dec
      << int(msg.extended_protocol_discriminator.extended_proto_discriminator);
  MLOG(MDEBUG) << " Spare half octet : 0";
  MLOG(MDEBUG) << " Security Header Type : " << dec
               << int(msg.sec_header_type.sec_hdr);
  MLOG(MDEBUG) << " Message Type : 0x" << hex << int(msg.message_type.msg_type);
  MLOG(MDEBUG) << "\n\n";
  return 0;
}
}  // namespace magma5g

int main(void) {
  int ret;
  ret = magma5g::encode();
  return 0;
}
