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

/* using this stub code we are going to test Encoding functionality of Secuirty
 * Mode Command Message */

#include <iostream>
#include <iomanip>
#include "M5GSecurityModeCommand.h"
#include "M5GCommonDefs.h"

using namespace std;
using namespace magma5g;

namespace magma5g {
int encode(void) {
  int ret           = 0;
  uint8_t buffer[8] = {};
  int len           = 8;
  SecurityModeCommandMsg msg;

  // Message to be Encoded
  msg.extended_protocol_discriminator.extended_proto_discriminator = 126;
  msg.sec_header_type.sec_hdr                                      = 0;
  msg.spare_half_octet.spare                                       = 0;
  msg.message_type.msg_type                                        = 0x5D;
  msg.nas_sec_algorithms.tca                                       = 0;
  msg.nas_sec_algorithms.tia                                       = 1;
  msg.nas_key_set_identifier.tsc                                   = 0;
  msg.nas_key_set_identifier.nas_key_set_identifier                = 0;
  msg.spare_half_octet.spare                                       = 0;
  msg.ue_sec_capability.length                                     = 2;
  msg.ue_sec_capability.ea0                                        = 1;
  msg.ue_sec_capability.ea1                                        = 0;
  msg.ue_sec_capability.ea2                                        = 0;
  msg.ue_sec_capability.ea3                                        = 0;
  msg.ue_sec_capability.ea4                                        = 0;
  msg.ue_sec_capability.ea5                                        = 0;
  msg.ue_sec_capability.ea6                                        = 0;
  msg.ue_sec_capability.ea7                                        = 0;
  msg.ue_sec_capability.ia0                                        = 0;
  msg.ue_sec_capability.ia1                                        = 1;
  msg.ue_sec_capability.ia2                                        = 0;
  msg.ue_sec_capability.ia3                                        = 0;
  msg.ue_sec_capability.ia4                                        = 0;
  msg.ue_sec_capability.ia5                                        = 0;
  msg.ue_sec_capability.ia6                                        = 0;
  msg.ue_sec_capability.ia7                                        = 0;

  MLOG(MDEBUG) << "---Encoding message--- \n";
  ret = msg.EncodeSecurityModeCommandMsg(&msg, buffer, len);

  MLOG(MDEBUG) << " ENCODED MESSAGE : " << setfill('0') << hex << setw(2)
               << int(buffer[0]) << hex << setw(2) << int(buffer[1]) << hex
               << setw(2) << int(buffer[2]) << hex << setw(2) << int(buffer[3])
               << hex << setw(2) << int(buffer[4]) << hex << setw(2)
               << int(buffer[5]) << hex << setw(2) << int(buffer[6]) << hex
               << setw(2) << int(buffer[7]) << "\n";

  MLOG(MDEBUG) << "---Decoding encoded message--- ";
  int ret2 = 0;
  ret2     = msg.DecodeSecurityModeCommandMsg(&msg, buffer, len);
  MLOG(MDEBUG) << "\n\n";
  return 0;
}
}  // namespace magma5g

int main(void) {
  int ret;
  ret = magma5g::encode();
  return 0;
}
