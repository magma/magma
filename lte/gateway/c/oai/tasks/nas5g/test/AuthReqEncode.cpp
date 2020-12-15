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

/*** Using this stub code we are going to test Encoding Functionality of
 * Authentication Request Message ***/

#include <iostream>
#include <iomanip>
#include <cstring>
#include "M5GAuthenticationRequest.h"
#include "M5GCommonDefs.h"

using namespace std;
using namespace magma5g;

namespace magma5g {
int Encode(void) {
  int ret            = 0;
  uint8_t buffer[50] = {};
  int len            = 50;

  AuthenticationRequestMsg AuthReq;
  AuthReq.extended_protocol_discriminator.extended_proto_discriminator = 126;
  AuthReq.sec_header_type.sec_hdr                                      = 0;
  AuthReq.spare_half_octet.spare                                       = 0;
  AuthReq.message_type.msg_type                                        = 0x56;
  AuthReq.nas_key_set_identifier.tsc                                   = 0;
  AuthReq.nas_key_set_identifier.nas_key_set_identifier                = 0;
  uint8_t abba_buff[] = {0x71, 0x00, 0x0d, 0x01};
  AuthReq.abba.contents.assign((const char*) abba_buff, 4);
  AuthReq.auth_rand.iei = 0x21;
  uint8_t rand_buff[]   = {0x4f, 0x42, 0x84, 0xea, 0x63, 0xd6, 0x5a, 0xff,
                         0xa7, 0xbf, 0xe0, 0xae, 0xc6, 0xb2, 0xc8, 0x6b};
  AuthReq.auth_rand.rand_val.assign((const char*) rand_buff, 16);
  AuthReq.auth_autn.iei = 0x20;
  uint8_t autn_buff[]   = {0x6e, 0x3d, 0xa7, 0x99, 0x46, 0x2b, 0x80, 0x00,
                         0xc0, 0x66, 0x08, 0x12, 0xcd, 0xf1, 0x41, 0x3b};
  AuthReq.auth_autn.AUTN.assign((const char*) autn_buff, 16);

  // Encoding the Authentication Message
  MLOG(MDEBUG) << "\n\n---Encoding Authentication request Message---\n\n";
  ret = AuthReq.EncodeAuthenticationRequestMsg(&AuthReq, buffer, len);

  // Printing Encoded Message
  MLOG(MDEBUG)
      << "\n\n    ENCODED MESSAGE : " << setfill('0') << hex << int(buffer[0])
      << hex << setw(2) << int(buffer[1]) << hex << int(buffer[2]) << hex
      << setw(2) << int(buffer[3]) << hex << setw(2) << int(buffer[4]) << hex
      << setw(2) << int(buffer[5]) << hex << setw(2) << int(buffer[6]) << hex
      << setw(2) << int(buffer[7]) << hex << setw(2) << int(buffer[8]) << hex
      << setw(2) << int(buffer[9]) << hex << setw(2) << int(buffer[10]) << hex
      << setw(2) << int(buffer[11]) << hex << setw(2) << int(buffer[12]) << hex
      << setw(2) << int(buffer[13]) << hex << setw(2) << int(buffer[14]) << hex
      << setw(2) << int(buffer[15]) << hex << setw(2) << int(buffer[16]) << hex
      << setw(2) << int(buffer[17]) << hex << setw(2) << int(buffer[18]) << hex
      << setw(2) << int(buffer[19]) << hex << setw(2) << int(buffer[20]) << hex
      << setw(2) << int(buffer[21]) << hex << setw(2) << int(buffer[22]) << hex
      << setw(2) << int(buffer[23]) << hex << setw(2) << int(buffer[24]) << hex
      << setw(2) << int(buffer[25]) << hex << setw(2) << int(buffer[26]) << hex
      << setw(2) << int(buffer[27]) << hex << setw(2) << int(buffer[28]) << hex
      << setw(2) << int(buffer[29]) << hex << setw(2) << int(buffer[30]) << hex
      << setw(2) << int(buffer[31]) << hex << setw(2) << int(buffer[32]) << hex
      << setw(2) << int(buffer[33]) << hex << setw(2) << int(buffer[34]) << hex
      << setw(2) << int(buffer[35]) << hex << setw(2) << int(buffer[36]) << hex
      << setw(2) << int(buffer[37]) << hex << setw(2) << int(buffer[38]) << hex
      << setw(2) << int(buffer[39]) << hex << setw(2) << int(buffer[40]) << hex
      << setw(2) << int(buffer[41]) << hex << setw(2) << int(buffer[42]) << hex
      << setw(2) << int(buffer[43]) << "\n\n";
  return 0;
}
}  // namespace magma5g

// Main Function to call Test Encode Function
int main(void) {
  int ret;
  ret = magma5g::Encode();
  return 0;
}
