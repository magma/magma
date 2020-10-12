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
#include "AuthenticationResult.h"
#include "CommonDefs.h"
using namespace std;
using namespace magma5g;

namespace magma5g {
int encode(void) {
  int ret           = 0;
  uint8_t buffer[13] = {};
  int len           = 13;
  AuthenticationResultMsg msg;
  // Message to be Encoded
  msg.extendedprotocoldiscriminator.extendedprotodiscriminator = 126;
  msg.sparehalfoctet.spare                                     = 0;
  msg.securityheadertype.securityhdr                           = 0;
  msg.messagetype.msgtype                                      = 0x46;
  msg.naskeysetidentifier.tsc                                  = 0;
  msg.naskeysetidentifier.naskeysetidentifier                  = 0;
  msg.eapmessage.len                                           = 7;
  uint8_t eap_buff[]                                           = {0x71, 0x00, 0x0d, 0x01, 0x77, 0x93, 0x11};
  msg.eapmessage.eap.assign((const char*) eap_buff, 7);
 
  MLOG(MDEBUG) << "---Encoding message--- \n";
  ret = msg.EncodeAuthenticationResultMsg(&msg, buffer, len);

  MLOG(MDEBUG) << " ENCODED MESSAGE : ";
  for(int i=0; i < 13; i++) {
    MLOG(MDEBUG) << setfill('0') << hex << setw(2) << int(buffer[i]);
  }
}
}  // namespace magma5g

int main(void) {
  int ret;
  ret = magma5g::encode();
  return 0;
}
