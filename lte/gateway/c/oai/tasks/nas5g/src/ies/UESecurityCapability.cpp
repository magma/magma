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

#include <iostream>
#include <iomanip>
#include <sstream>
#include <cstdint>
#include <cstring>
#include "UESecurityCapability.h"
#include "CommonDefs.h"

using namespace std;
namespace magma5g {
UESecurityCapabilityMsg::UESecurityCapabilityMsg(){};
UESecurityCapabilityMsg::~UESecurityCapabilityMsg(){};

int UESecurityCapabilityMsg::DecodeUESecurityCapabilityMsg(
    UESecurityCapabilityMsg* uesecuritycapability, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded = 0;
  MLOG(MDEBUG) << " Decoding UE Security Capability : ";

  // Checking IEI and pointer
  if (iei > 0) {
    CHECK_IEI_DECODER(iei, (unsigned char) *buffer);
    decoded++;
  }

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, UE_SECURITY_CAPABILITY_MIN_LENGTH, len);

  uesecuritycapability->length = *(buffer + decoded);
  decoded++;
  MLOG(MDEBUG) << " length = " << hex << int(uesecuritycapability->length)
               << endl;

  // 5GS encryption algorithms
  uesecuritycapability->ea0 = (*(buffer + decoded) >> 7) & 0x1;
  uesecuritycapability->ea1 = (*(buffer + decoded) >> 6) & 0x1;
  uesecuritycapability->ea2 = (*(buffer + decoded) >> 5) & 0x1;
  uesecuritycapability->ea3 = (*(buffer + decoded) >> 4) & 0x1;
  uesecuritycapability->ea4 = (*(buffer + decoded) >> 3) & 0x1;
  uesecuritycapability->ea5 = (*(buffer + decoded) >> 2) & 0x1;
  uesecuritycapability->ea6 = (*(buffer + decoded) >> 1) & 0x1;
  uesecuritycapability->ea7 = *(buffer + decoded) & 0x1;
  decoded++;

  // 5GS integrity algorithm
  uesecuritycapability->ia0 = (*(buffer + decoded) >> 7) & 0x1;
  uesecuritycapability->ia1 = (*(buffer + decoded) >> 6) & 0x1;
  uesecuritycapability->ia2 = (*(buffer + decoded) >> 5) & 0x1;
  uesecuritycapability->ia3 = (*(buffer + decoded) >> 4) & 0x1;
  uesecuritycapability->ia4 = (*(buffer + decoded) >> 3) & 0x1;
  uesecuritycapability->ia5 = (*(buffer + decoded) >> 2) & 0x1;
  uesecuritycapability->ia6 = (*(buffer + decoded) >> 1) & 0x1;
  uesecuritycapability->ia7 = *(buffer + decoded) & 0x1;
  decoded++;

  // Decoded 5GS encryption algorithms
  MLOG(MDEBUG) << " ea0 = " << hex << int(uesecuritycapability->ea0) << endl;
  MLOG(MDEBUG) << " ea1 = " << hex << int(uesecuritycapability->ea1) << endl;
  MLOG(MDEBUG) << " ea2 = " << hex << int(uesecuritycapability->ea2) << endl;
  MLOG(MDEBUG) << " ea3 = " << hex << int(uesecuritycapability->ea3) << endl;
  MLOG(MDEBUG) << " ea4 = " << hex << int(uesecuritycapability->ea4) << endl;
  MLOG(MDEBUG) << " ea5 = " << hex << int(uesecuritycapability->ea5) << endl;
  MLOG(MDEBUG) << " ea6 = " << hex << int(uesecuritycapability->ea6) << endl;
  MLOG(MDEBUG) << " ea7 = " << hex << int(uesecuritycapability->ea7) << endl;
  // Decoded 5GS integrity algorithm
  MLOG(MDEBUG) << " ia0 = " << hex << int(uesecuritycapability->ia0) << endl;
  MLOG(MDEBUG) << " ia1 = " << hex << int(uesecuritycapability->ia1) << endl;
  MLOG(MDEBUG) << " ia2 = " << hex << int(uesecuritycapability->ia2) << endl;
  MLOG(MDEBUG) << " ia3 = " << hex << int(uesecuritycapability->ia3) << endl;
  MLOG(MDEBUG) << " ia4 = " << hex << int(uesecuritycapability->ia4) << endl;
  MLOG(MDEBUG) << " ia5 = " << hex << int(uesecuritycapability->ia5) << endl;
  MLOG(MDEBUG) << " ia6 = " << hex << int(uesecuritycapability->ia6) << endl;
  MLOG(MDEBUG) << " ia7 = " << hex << int(uesecuritycapability->ia7) << endl;

#ifdef HANDLE_POST_MVC
  // EPS encryption algorithms
  uesecuritycapability->eea0 = (*(buffer + decoded) >> 7) & 0x1;
  uesecuritycapability->eea1 = (*(buffer + decoded) >> 6) & 0x1;
  uesecuritycapability->eea2 = (*(buffer + decoded) >> 5) & 0x1;
  uesecuritycapability->eea3 = (*(buffer + decoded) >> 4) & 0x1;
  uesecuritycapability->eea4 = (*(buffer + decoded) >> 3) & 0x1;
  uesecuritycapability->eea5 = (*(buffer + decoded) >> 2) & 0x1;
  uesecuritycapability->eea6 = (*(buffer + decoded) >> 1) & 0x1;
  uesecuritycapability->eea7 = *(buffer + decoded) & 0x1;
  decoded++;
  // EPS integrity algorithms
  uesecuritycapability->eia0 = (*(buffer + decoded) >> 7) & 0x1;
  uesecuritycapability->eia1 = (*(buffer + decoded) >> 6) & 0x1;
  uesecuritycapability->eia2 = (*(buffer + decoded) >> 5) & 0x1;
  uesecuritycapability->eia3 = (*(buffer + decoded) >> 4) & 0x1;
  uesecuritycapability->eia4 = (*(buffer + decoded) >> 3) & 0x1;
  uesecuritycapability->eia5 = (*(buffer + decoded) >> 2) & 0x1;
  uesecuritycapability->eia6 = (*(buffer + decoded) >> 1) & 0x1;
  uesecuritycapability->eia7 = *(buffer + decoded) & 0x1;
  decoded++;

  // Decoded EPS encryption algorithms
  MLOG(MDEBUG) << " eea0 = " << hex << int(uesecuritycapability->eea0) << endl;
  MLOG(MDEBUG) << " eea1 = " << hex << int(uesecuritycapability->eea1) << endl;
  MLOG(MDEBUG) << " eea2 = " << hex << int(uesecuritycapability->eea2) << endl;
  MLOG(MDEBUG) << " eea3 = " << hex << int(uesecuritycapability->eea3) << endl;
  MLOG(MDEBUG) << " eea4 = " << hex << int(uesecuritycapability->eea4) << endl;
  MLOG(MDEBUG) << " eea5 = " << hex << int(uesecuritycapability->eea5) << endl;
  MLOG(MDEBUG) << " eea6 = " << hex << int(uesecuritycapability->eea6) << endl;
  MLOG(MDEBUG) << " eea7 = " << hex << int(uesecuritycapability->eea7) << endl;
  // Decoded EPS integrity algorithms
  MLOG(MDEBUG) << " eia0 = " << hex << int(uesecuritycapability->eia0) << endl;
  MLOG(MDEBUG) << " eia1 = " << hex << int(uesecuritycapability->eia1) << endl;
  MLOG(MDEBUG) << " eia2 = " << hex << int(uesecuritycapability->eia2) << endl;
  MLOG(MDEBUG) << " eia3 = " << hex << int(uesecuritycapability->eia3) << endl;
  MLOG(MDEBUG) << " eia4 = " << hex << int(uesecuritycapability->eia4) << endl;
  MLOG(MDEBUG) << " eia5 = " << hex << int(uesecuritycapability->eia5) << endl;
  MLOG(MDEBUG) << " eia6 = " << hex << int(uesecuritycapability->eia6) << endl;
  MLOG(MDEBUG) << " eia7 = " << hex << int(uesecuritycapability->eia7) << endl;
#endif

  return (decoded);
};

// Encode UE Security Capability
int UESecurityCapabilityMsg::EncodeUESecurityCapabilityMsg(
    UESecurityCapabilityMsg* uesecuritycapability, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int encoded = 0;
  MLOG(DEBUG) << " Encoding UE Security Capability : ";

  // Checking IEI and pointer
  if (iei > 0) {
    CHECK_IEI_ENCODER((unsigned char) iei, uesecuritycapability->iei);
    *buffer = iei;
    encoded++;
  }

  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, UE_SECURITY_CAPABILITY_MIN_LENGTH, len);

  *(buffer + encoded) = uesecuritycapability->length;
  MLOG(MDEBUG) << "Length : " << setfill('0') << hex << setw(2)
               << int(*(buffer + encoded));
  encoded++;

  // 5GS encryption algorithms
  *(buffer + encoded) = 0x00 | ((uesecuritycapability->ea0 & 0x1) << 7) |
                        ((uesecuritycapability->ea1 & 0x1) << 6) |
                        ((uesecuritycapability->ea2 & 0x1) << 5) |
                        ((uesecuritycapability->ea3 & 0x1) << 4) |
                        ((uesecuritycapability->ea4 & 0x1) << 3) |
                        ((uesecuritycapability->ea5 & 0x1) << 2) |
                        ((uesecuritycapability->ea6 & 0x1) << 1) |
                        ((uesecuritycapability->ea7) & 0x1);
  MLOG(MDEBUG) << " 5GS Encryption Algorithms Supported : " << hex
               << int(*(buffer + encoded));
  encoded++;

  // 5GS integrity algorithms
  *(buffer + encoded) = 0x00 | ((uesecuritycapability->ia0 & 0x1) << 7) |
                        ((uesecuritycapability->ia1 & 0x1) << 6) |
                        ((uesecuritycapability->ia2 & 0x1) << 5) |
                        ((uesecuritycapability->ia3 & 0x1) << 4) |
                        ((uesecuritycapability->ia4 & 0x1) << 3) |
                        ((uesecuritycapability->ia5 & 0x1) << 2) |
                        ((uesecuritycapability->ia6 & 0x1) << 1) |
                        ((uesecuritycapability->ia7) & 0x1);
  MLOG(MDEBUG) << " 5GS Integrity Algorithms Supported : " << hex
               << int(*(buffer + encoded));
  encoded++;

#ifdef HANDLE_POST_MVC
  // EPS encryption algorithms
  *(buffer + encoded) = 0x00 | ((uesecuritycapability->eea0 & 0x1) << 7) |
                        ((uesecuritycapability->eea1 & 0x1) << 6) |
                        ((uesecuritycapability->eea2 & 0x1) << 5) |
                        ((uesecuritycapability->eea3 & 0x1) << 4) |
                        ((uesecuritycapability->eea4 & 0x1) << 3) |
                        ((uesecuritycapability->eea5 & 0x1) << 2) |
                        ((uesecuritycapability->eea6 & 0x1) << 1) |
                        ((uesecuritycapability->eea7) & 0x1);
  MLOG(MDEBUG) << " EPS Encryption Algorithms Supported : " << hex
               << int(*(buffer + encoded));
  encoded++;

  // EPS integrity algorithms
  *(buffer + encoded) = 0x00 | ((uesecuritycapability->eia0 & 0x1) << 7) |
                        ((uesecuritycapability->eia1 & 0x1) << 6) |
                        ((uesecuritycapability->eia2 & 0x1) << 5) |
                        ((uesecuritycapability->eia3 & 0x1) << 4) |
                        ((uesecuritycapability->eia4 & 0x1) << 3) |
                        ((uesecuritycapability->eia5 & 0x1) << 2) |
                        ((uesecuritycapability->eia6 & 0x1) << 1) |
                        ((uesecuritycapability->eia7) & 0x1);
  MLOG(MDEBUG) << " EPS Integrity Algorithms Supported : " << hex
               << int(*(buffer + encoded));
  encoded++;
#endif
  return 0;
};
}  // namespace magma5g
