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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GUESecurityCapability.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
UESecurityCapabilityMsg::UESecurityCapabilityMsg(){};
UESecurityCapabilityMsg::~UESecurityCapabilityMsg(){};

int UESecurityCapabilityMsg::DecodeUESecurityCapabilityMsg(
    UESecurityCapabilityMsg* ue_sec_capability, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded        = 0;
  uint8_t type_len   = sizeof(uint8_t);
  uint8_t length_len = sizeof(uint8_t);

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
  MLOG(MDEBUG) << " length = " << std::hex << int(ue_sec_capability->length)
               << std::endl;

  // 5GS encryption algorithms
  ea                     = *(buffer + decoded);
  ue_sec_capability->ea0 = (ea >> 7) & 0x1;
  ue_sec_capability->ea1 = (ea >> 6) & 0x1;
  ue_sec_capability->ea2 = (ea >> 5) & 0x1;
  ue_sec_capability->ea3 = (ea >> 4) & 0x1;
  ue_sec_capability->ea4 = (ea >> 3) & 0x1;
  ue_sec_capability->ea5 = (ea >> 2) & 0x1;
  ue_sec_capability->ea6 = (ea >> 1) & 0x1;
  ue_sec_capability->ea7 = ea & 0x1;
  decoded++;

  // 5GS integrity algorithm
  ia                     = *(buffer + decoded);
  ue_sec_capability->ia0 = (ia >> 7) & 0x1;
  ue_sec_capability->ia1 = (ia >> 6) & 0x1;
  ue_sec_capability->ia2 = (ia >> 5) & 0x1;
  ue_sec_capability->ia3 = (ia >> 4) & 0x1;
  ue_sec_capability->ia4 = (ia >> 3) & 0x1;
  ue_sec_capability->ia5 = (ia >> 2) & 0x1;
  ue_sec_capability->ia6 = (ia >> 1) & 0x1;
  ue_sec_capability->ia7 = ia & 0x1;
  decoded++;

  // If any optional buffers are present skip it.
  // 2 = 1 Byte for type + 1 Byte for length
  type_len   = sizeof(uint8_t);
  length_len = sizeof(uint8_t);
  if (ue_sec_capability->length > (decoded - (type_len + length_len))) {
    // 5GS encryption algorithms
    ue_sec_capability->eea0    = (*(buffer + decoded) >> 7) & 0x1;
    ue_sec_capability->ea1_128 = (*(buffer + decoded) >> 6) & 0x1;
    ue_sec_capability->ea2_128 = (*(buffer + decoded) >> 5) & 0x1;
    ue_sec_capability->ea3_128 = (*(buffer + decoded) >> 4) & 0x1;
    ue_sec_capability->eea4    = (*(buffer + decoded) >> 3) & 0x1;
    ue_sec_capability->eea5    = (*(buffer + decoded) >> 2) & 0x1;
    ue_sec_capability->eea6    = (*(buffer + decoded) >> 1) & 0x1;
    ue_sec_capability->eea7    = *(buffer + decoded) & 0x1;
    decoded++;

    // 5GS integrity algorithm
    ue_sec_capability->eia0     = (*(buffer + decoded) >> 7) & 0x1;
    ue_sec_capability->eia1_128 = (*(buffer + decoded) >> 6) & 0x1;
    ue_sec_capability->eia2_128 = (*(buffer + decoded) >> 5) & 0x1;
    ue_sec_capability->eia3_128 = (*(buffer + decoded) >> 4) & 0x1;
    ue_sec_capability->eia4     = (*(buffer + decoded) >> 3) & 0x1;
    ue_sec_capability->eia5     = (*(buffer + decoded) >> 2) & 0x1;
    ue_sec_capability->eia6     = (*(buffer + decoded) >> 1) & 0x1;
    ue_sec_capability->eia7     = *(buffer + decoded) & 0x1;
    decoded++;

    // decoded = type_len + length_len + ue_sec_capability->length;
  }

  // Decoded 5GS encryption algorithms
  MLOG(MDEBUG) << " ea0 = " << std::hex << int(ue_sec_capability->ea0)
               << std::endl;
  MLOG(MDEBUG) << " ea1 = " << std::hex << int(ue_sec_capability->ea1)
               << std::endl;
  MLOG(MDEBUG) << " ea2 = " << std::hex << int(ue_sec_capability->ea2)
               << std::endl;
  MLOG(MDEBUG) << " ea3 = " << std::hex << int(ue_sec_capability->ea3)
               << std::endl;
  MLOG(MDEBUG) << " ea4 = " << std::hex << int(ue_sec_capability->ea4)
               << std::endl;
  MLOG(MDEBUG) << " ea5 = " << std::hex << int(ue_sec_capability->ea5)
               << std::endl;
  MLOG(MDEBUG) << " ea6 = " << std::hex << int(ue_sec_capability->ea6)
               << std::endl;
  MLOG(MDEBUG) << " ea7 = " << std::hex << int(ue_sec_capability->ea7)
               << std::endl;
  // Decoded 5GS integrity algorithm
  MLOG(MDEBUG) << " ia0 = " << std::hex << int(ue_sec_capability->ia0)
               << std::endl;
  MLOG(MDEBUG) << " ia1 = " << std::hex << int(ue_sec_capability->ia1)
               << std::endl;
  MLOG(MDEBUG) << " ia2 = " << std::hex << int(ue_sec_capability->ia2)
               << std::endl;
  MLOG(MDEBUG) << " ia3 = " << std::hex << int(ue_sec_capability->ia3)
               << std::endl;
  MLOG(MDEBUG) << " ia4 = " << std::hex << int(ue_sec_capability->ia4)
               << std::endl;
  MLOG(MDEBUG) << " ia5 = " << std::hex << int(ue_sec_capability->ia5)
               << std::endl;
  MLOG(MDEBUG) << " ia6 = " << std::hex << int(ue_sec_capability->ia6)
               << std::endl;
  MLOG(MDEBUG) << " ia7 = " << std::hex << int(ue_sec_capability->ia7)
               << std::endl;

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
  MLOG(MDEBUG) << "Length : " << std::setfill('0') << std::hex << std::setw(2)
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
  MLOG(MDEBUG) << " 5GS Encryption Algorithms Supported : " << std::hex
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
  MLOG(MDEBUG) << " 5GS Integrity Algorithms Supported : " << std::hex
               << int(*(buffer + encoded));
  encoded++;

  if (ue_sec_capability->length > 2) {
    // 5GS encryption algorithms
    *(buffer + encoded) = 0x00 | ((ue_sec_capability->eea0 & 0x1) << 7) |
                          ((ue_sec_capability->ea1_128 & 0x1) << 6) |
                          ((ue_sec_capability->ea2_128 & 0x1) << 5) |
                          ((ue_sec_capability->ea3_128 & 0x1) << 4) |
                          ((ue_sec_capability->eea4 & 0x1) << 3) |
                          ((ue_sec_capability->eea5 & 0x1) << 2) |
                          ((ue_sec_capability->eea6 & 0x1) << 1) |
                          ((ue_sec_capability->eea7) & 0x1);

    encoded++;

    // 5GS integrity algorithms
    *(buffer + encoded) = 0x00 | ((ue_sec_capability->eia0 & 0x1) << 7) |
                          ((ue_sec_capability->eia1_128 & 0x1) << 6) |
                          ((ue_sec_capability->eia2_128 & 0x1) << 5) |
                          ((ue_sec_capability->eia3_128 & 0x1) << 4) |
                          ((ue_sec_capability->eia4 & 0x1) << 3) |
                          ((ue_sec_capability->eia5 & 0x1) << 2) |
                          ((ue_sec_capability->eia6 & 0x1) << 1) |
                          ((ue_sec_capability->eia7) & 0x1);

    encoded++;
  }

  MLOG(DEBUG) << " Encoding UE Security Capability : encoded  " << encoded
              << std::endl;
  return encoded;
};
}  // namespace magma5g
