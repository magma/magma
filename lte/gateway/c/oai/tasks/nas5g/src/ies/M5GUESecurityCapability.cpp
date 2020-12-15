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
#include "M5GUESecurityCapability.h"
#include "M5GCommonDefs.h"

using namespace std;
namespace magma5g {
UESecurityCapabilityMsg::UESecurityCapabilityMsg(){};
UESecurityCapabilityMsg::~UESecurityCapabilityMsg(){};

int UESecurityCapabilityMsg::DecodeUESecurityCapabilityMsg(
    UESecurityCapabilityMsg* ue_sec_capability, uint8_t iei, uint8_t* buffer,
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

  ue_sec_capability->length = *(buffer + decoded);
  decoded++;
  MLOG(MDEBUG) << " length = " << hex << int(ue_sec_capability->length) << endl;

  // 5GS encryption algorithms
  ue_sec_capability->ea0 = (*(buffer + decoded) >> 7) & 0x1;
  ue_sec_capability->ea1 = (*(buffer + decoded) >> 6) & 0x1;
  ue_sec_capability->ea2 = (*(buffer + decoded) >> 5) & 0x1;
  ue_sec_capability->ea3 = (*(buffer + decoded) >> 4) & 0x1;
  ue_sec_capability->ea4 = (*(buffer + decoded) >> 3) & 0x1;
  ue_sec_capability->ea5 = (*(buffer + decoded) >> 2) & 0x1;
  ue_sec_capability->ea6 = (*(buffer + decoded) >> 1) & 0x1;
  ue_sec_capability->ea7 = *(buffer + decoded) & 0x1;
  decoded++;

  // 5GS integrity algorithm
  ue_sec_capability->ia0 = (*(buffer + decoded) >> 7) & 0x1;
  ue_sec_capability->ia1 = (*(buffer + decoded) >> 6) & 0x1;
  ue_sec_capability->ia2 = (*(buffer + decoded) >> 5) & 0x1;
  ue_sec_capability->ia3 = (*(buffer + decoded) >> 4) & 0x1;
  ue_sec_capability->ia4 = (*(buffer + decoded) >> 3) & 0x1;
  ue_sec_capability->ia5 = (*(buffer + decoded) >> 2) & 0x1;
  ue_sec_capability->ia6 = (*(buffer + decoded) >> 1) & 0x1;
  ue_sec_capability->ia7 = *(buffer + decoded) & 0x1;
  decoded++;

  // Decoded 5GS encryption algorithms
  MLOG(MDEBUG) << " ea0 = " << hex << int(ue_sec_capability->ea0) << endl;
  MLOG(MDEBUG) << " ea1 = " << hex << int(ue_sec_capability->ea1) << endl;
  MLOG(MDEBUG) << " ea2 = " << hex << int(ue_sec_capability->ea2) << endl;
  MLOG(MDEBUG) << " ea3 = " << hex << int(ue_sec_capability->ea3) << endl;
  MLOG(MDEBUG) << " ea4 = " << hex << int(ue_sec_capability->ea4) << endl;
  MLOG(MDEBUG) << " ea5 = " << hex << int(ue_sec_capability->ea5) << endl;
  MLOG(MDEBUG) << " ea6 = " << hex << int(ue_sec_capability->ea6) << endl;
  MLOG(MDEBUG) << " ea7 = " << hex << int(ue_sec_capability->ea7) << endl;
  // Decoded 5GS integrity algorithm
  MLOG(MDEBUG) << " ia0 = " << hex << int(ue_sec_capability->ia0) << endl;
  MLOG(MDEBUG) << " ia1 = " << hex << int(ue_sec_capability->ia1) << endl;
  MLOG(MDEBUG) << " ia2 = " << hex << int(ue_sec_capability->ia2) << endl;
  MLOG(MDEBUG) << " ia3 = " << hex << int(ue_sec_capability->ia3) << endl;
  MLOG(MDEBUG) << " ia4 = " << hex << int(ue_sec_capability->ia4) << endl;
  MLOG(MDEBUG) << " ia5 = " << hex << int(ue_sec_capability->ia5) << endl;
  MLOG(MDEBUG) << " ia6 = " << hex << int(ue_sec_capability->ia6) << endl;
  MLOG(MDEBUG) << " ia7 = " << hex << int(ue_sec_capability->ia7) << endl;

#ifdef HANDLE_POST_MVC
  // EPS encryption algorithms
  ue_sec_capability->eea0 = (*(buffer + decoded) >> 7) & 0x1;
  ue_sec_capability->eea1 = (*(buffer + decoded) >> 6) & 0x1;
  ue_sec_capability->eea2 = (*(buffer + decoded) >> 5) & 0x1;
  ue_sec_capability->eea3 = (*(buffer + decoded) >> 4) & 0x1;
  ue_sec_capability->eea4 = (*(buffer + decoded) >> 3) & 0x1;
  ue_sec_capability->eea5 = (*(buffer + decoded) >> 2) & 0x1;
  ue_sec_capability->eea6 = (*(buffer + decoded) >> 1) & 0x1;
  ue_sec_capability->eea7 = *(buffer + decoded) & 0x1;
  decoded++;
  // EPS integrity algorithms
  ue_sec_capability->eia0 = (*(buffer + decoded) >> 7) & 0x1;
  ue_sec_capability->eia1 = (*(buffer + decoded) >> 6) & 0x1;
  ue_sec_capability->eia2 = (*(buffer + decoded) >> 5) & 0x1;
  ue_sec_capability->eia3 = (*(buffer + decoded) >> 4) & 0x1;
  ue_sec_capability->eia4 = (*(buffer + decoded) >> 3) & 0x1;
  ue_sec_capability->eia5 = (*(buffer + decoded) >> 2) & 0x1;
  ue_sec_capability->eia6 = (*(buffer + decoded) >> 1) & 0x1;
  ue_sec_capability->eia7 = *(buffer + decoded) & 0x1;
  decoded++;

  // Decoded EPS encryption algorithms
  MLOG(MDEBUG) << " eea0 = " << hex << int(ue_sec_capability->eea0) << endl;
  MLOG(MDEBUG) << " eea1 = " << hex << int(ue_sec_capability->eea1) << endl;
  MLOG(MDEBUG) << " eea2 = " << hex << int(ue_sec_capability->eea2) << endl;
  MLOG(MDEBUG) << " eea3 = " << hex << int(ue_sec_capability->eea3) << endl;
  MLOG(MDEBUG) << " eea4 = " << hex << int(ue_sec_capability->eea4) << endl;
  MLOG(MDEBUG) << " eea5 = " << hex << int(ue_sec_capability->eea5) << endl;
  MLOG(MDEBUG) << " eea6 = " << hex << int(ue_sec_capability->eea6) << endl;
  MLOG(MDEBUG) << " eea7 = " << hex << int(ue_sec_capability->eea7) << endl;
  // Decoded EPS integrity algorithms
  MLOG(MDEBUG) << " eia0 = " << hex << int(ue_sec_capability->eia0) << endl;
  MLOG(MDEBUG) << " eia1 = " << hex << int(ue_sec_capability->eia1) << endl;
  MLOG(MDEBUG) << " eia2 = " << hex << int(ue_sec_capability->eia2) << endl;
  MLOG(MDEBUG) << " eia3 = " << hex << int(ue_sec_capability->eia3) << endl;
  MLOG(MDEBUG) << " eia4 = " << hex << int(ue_sec_capability->eia4) << endl;
  MLOG(MDEBUG) << " eia5 = " << hex << int(ue_sec_capability->eia5) << endl;
  MLOG(MDEBUG) << " eia6 = " << hex << int(ue_sec_capability->eia6) << endl;
  MLOG(MDEBUG) << " eia7 = " << hex << int(ue_sec_capability->eia7) << endl;
#endif

  return (decoded);
};

// Encode UE Security Capability
int UESecurityCapabilityMsg::EncodeUESecurityCapabilityMsg(
    UESecurityCapabilityMsg* ue_sec_capability, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int encoded = 0;
  MLOG(DEBUG) << " Encoding UE Security Capability : ";

  // Checking IEI and pointer
  if (iei > 0) {
    CHECK_IEI_ENCODER((unsigned char) iei, ue_sec_capability->iei);
    *buffer = iei;
    encoded++;
  }

  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, UE_SECURITY_CAPABILITY_MIN_LENGTH, len);

  *(buffer + encoded) = ue_sec_capability->length;
  MLOG(MDEBUG) << "Length : " << setfill('0') << hex << setw(2)
               << int(*(buffer + encoded));
  encoded++;

  // 5GS encryption algorithms
  *(buffer + encoded) = 0x00 | ((ue_sec_capability->ea0 & 0x1) << 7) |
                        ((ue_sec_capability->ea1 & 0x1) << 6) |
                        ((ue_sec_capability->ea2 & 0x1) << 5) |
                        ((ue_sec_capability->ea3 & 0x1) << 4) |
                        ((ue_sec_capability->ea4 & 0x1) << 3) |
                        ((ue_sec_capability->ea5 & 0x1) << 2) |
                        ((ue_sec_capability->ea6 & 0x1) << 1) |
                        ((ue_sec_capability->ea7) & 0x1);
  MLOG(MDEBUG) << " 5GS Encryption Algorithms Supported : " << hex
               << int(*(buffer + encoded));
  encoded++;

  // 5GS integrity algorithms
  *(buffer + encoded) = 0x00 | ((ue_sec_capability->ia0 & 0x1) << 7) |
                        ((ue_sec_capability->ia1 & 0x1) << 6) |
                        ((ue_sec_capability->ia2 & 0x1) << 5) |
                        ((ue_sec_capability->ia3 & 0x1) << 4) |
                        ((ue_sec_capability->ia4 & 0x1) << 3) |
                        ((ue_sec_capability->ia5 & 0x1) << 2) |
                        ((ue_sec_capability->ia6 & 0x1) << 1) |
                        ((ue_sec_capability->ia7) & 0x1);
  MLOG(MDEBUG) << " 5GS Integrity Algorithms Supported : " << hex
               << int(*(buffer + encoded));
  encoded++;

#ifdef HANDLE_POST_MVC
  // EPS encryption algorithms
  *(buffer + encoded) = 0x00 | ((ue_sec_capability->eea0 & 0x1) << 7) |
                        ((ue_sec_capability->eea1 & 0x1) << 6) |
                        ((ue_sec_capability->eea2 & 0x1) << 5) |
                        ((ue_sec_capability->eea3 & 0x1) << 4) |
                        ((ue_sec_capability->eea4 & 0x1) << 3) |
                        ((ue_sec_capability->eea5 & 0x1) << 2) |
                        ((ue_sec_capability->eea6 & 0x1) << 1) |
                        ((ue_sec_capability->eea7) & 0x1);
  MLOG(MDEBUG) << " EPS Encryption Algorithms Supported : " << hex
               << int(*(buffer + encoded));
  encoded++;

  // EPS integrity algorithms
  *(buffer + encoded) = 0x00 | ((ue_sec_capability->eia0 & 0x1) << 7) |
                        ((ue_sec_capability->eia1 & 0x1) << 6) |
                        ((ue_sec_capability->eia2 & 0x1) << 5) |
                        ((ue_sec_capability->eia3 & 0x1) << 4) |
                        ((ue_sec_capability->eia4 & 0x1) << 3) |
                        ((ue_sec_capability->eia5 & 0x1) << 2) |
                        ((ue_sec_capability->eia6 & 0x1) << 1) |
                        ((ue_sec_capability->eia7) & 0x1);
  MLOG(MDEBUG) << " EPS Integrity Algorithms Supported : " << hex
               << int(*(buffer + encoded));
  encoded++;
#endif
  return encoded;
};
}  // namespace magma5g
