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
#include "SecurityModeCommand.h"
#include "CommonDefs.h"

using namespace std;
using namespace magma5g;

namespace magma5g {
int encode(void) {
  int ret           = 0;
  uint8_t buffer[8] = {};
  int len           = 8;
  SecurityModeCommandMsg msg;

  // Message to be Encoded
  msg.extendedprotocoldiscriminator.extendedprotodiscriminator = 126;
  msg.securityheadertype.securityhdr                           = 0;
  msg.sparehalfoctet.spare                                     = 0;
  msg.messagetype.msgtype                                      = 0x5D;
  msg.nassecurityalgorithms.tca                                = 0;
  msg.nassecurityalgorithms.tia                                = 1;
  msg.naskeysetidentifier.tsc                                  = 0;
  msg.naskeysetidentifier.naskeysetidentifier                  = 0;
  msg.sparehalfoctet.spare                                     = 0;
  msg.uesecuritycapability.length                              = 2;
  msg.uesecuritycapability.ea0                                 = 1;
  msg.uesecuritycapability.ea1                                 = 0;
  msg.uesecuritycapability.ea2                                 = 0;
  msg.uesecuritycapability.ea3                                 = 0;
  msg.uesecuritycapability.ea4                                 = 0;
  msg.uesecuritycapability.ea5                                 = 0;
  msg.uesecuritycapability.ea6                                 = 0;
  msg.uesecuritycapability.ea7                                 = 0;
  msg.uesecuritycapability.ia0                                 = 0;
  msg.uesecuritycapability.ia1                                 = 1;
  msg.uesecuritycapability.ia2                                 = 0;
  msg.uesecuritycapability.ia3                                 = 0;
  msg.uesecuritycapability.ia4                                 = 0;
  msg.uesecuritycapability.ia5                                 = 0;
  msg.uesecuritycapability.ia6                                 = 0;
  msg.uesecuritycapability.ia7                                 = 0;

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
