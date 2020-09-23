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

/*** Using this stub code we are going to test Decoding Functionality of
 * Authentication Response Message ***/

#include <iostream>
#include <iomanip>
#include <cstring>
#include "AuthenticationResponse.h"
#include "CommonDefs.h"

using namespace std;
using namespace magma5g;

namespace magma5g {
int Decode(void) {
  int ret          = 0;
  uint8_t buffer[] = {0x7E, 0x00, 0x57, 0x2D, 0x10, 0x25, 0xE8,
                      0x7B, 0x06, 0x52, 0xC3, 0xC6, 0x3B, 0x36,
                      0x82, 0x8B, 0x54, 0x51, 0x7E, 0xBF, 0x15};
  int len          = 21;
  AuthenticationResponseMsg AuthRes;

  // Decoding Authentication Response Message
  MLOG(MDEBUG) << " ---Authentication response Message---\n";
  ret = AuthRes.DecodeAuthenticationResponseMsg(&AuthRes, buffer, len);

  // Printing Decoded Authentication Response Message
  MLOG(MDEBUG) << " ---Decoded Message---\n";
  MLOG(MDEBUG)
      << " Extended Protocol Discriminator :" << dec
      << int(AuthRes.extendedprotocoldiscriminator.extendedprotodiscriminator);
  MLOG(MDEBUG) << " Spare Half Octet : " << dec
               << int(AuthRes.sparehalfoctet.spare);
  MLOG(MDEBUG) << " Security Header Type : " << dec
               << int(AuthRes.securityheadertype.securityhdr);
  MLOG(MDEBUG) << " Message Type : 0x" << hex
               << int(AuthRes.messagetype.msgtype);
  MLOG(MDEBUG) << " Response Parameter : "
               << "ElementID = " << hex << int(AuthRes.responseparameter.iei)
               << " Length = " << dec << int(AuthRes.responseparameter.length);
  MLOG(MDEBUG)
      << " RES : 0x" << setfill('0') << hex
      << int(AuthRes.responseparameter.ResponseParameter[0]) << hex << setw(2)
      << int(AuthRes.responseparameter.ResponseParameter[1]) << hex
      << int(AuthRes.responseparameter.ResponseParameter[2]) << hex << setw(2)
      << int(AuthRes.responseparameter.ResponseParameter[3]) << hex << setw(2)
      << int(AuthRes.responseparameter.ResponseParameter[4]) << hex << setw(2)
      << int(AuthRes.responseparameter.ResponseParameter[5]) << hex << setw(2)
      << int(AuthRes.responseparameter.ResponseParameter[6]) << hex << setw(2)
      << int(AuthRes.responseparameter.ResponseParameter[7]) << hex << setw(2)
      << int(AuthRes.responseparameter.ResponseParameter[8]) << hex << setw(2)
      << int(AuthRes.responseparameter.ResponseParameter[9]) << hex << setw(2)
      << int(AuthRes.responseparameter.ResponseParameter[10]) << hex << setw(2)
      << int(AuthRes.responseparameter.ResponseParameter[11]) << hex << setw(2)
      << int(AuthRes.responseparameter.ResponseParameter[12]) << hex << setw(2)
      << int(AuthRes.responseparameter.ResponseParameter[13]) << hex << setw(2)
      << int(AuthRes.responseparameter.ResponseParameter[14]) << hex << setw(2)
      << int(AuthRes.responseparameter.ResponseParameter[15]) << endl;

  return 0;
}
}  // namespace magma5g

// Main Function to call test Decode function
int main(void) {
  int ret;
  ret = magma5g::Decode();
  return 0;
}
