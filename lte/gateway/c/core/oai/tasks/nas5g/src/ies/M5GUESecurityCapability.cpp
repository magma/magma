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
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GUESecurityCapability.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
UESecurityCapabilityMsg::UESecurityCapabilityMsg(){};
UESecurityCapabilityMsg::~UESecurityCapabilityMsg(){};

int UESecurityCapabilityMsg::DecodeUESecurityCapabilityMsg(uint8_t iei,
                                                           uint8_t* buffer,
                                                           uint32_t len) {
  int decoded = 0;
  uint8_t type_len = sizeof(uint8_t);
  uint8_t length_len = sizeof(uint8_t);

  // Checking IEI and pointer
  if (iei > 0) {
    CHECK_IEI_DECODER(iei, (unsigned char)*buffer);
    this->iei = iei;
    decoded++;
  }

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(buffer,
                                       UE_SECURITY_CAPABILITY_MIN_LENGTH, len);

  this->length = *(buffer + decoded++);

  if (this->length <= 0) return (decoded);

  // 5GS encryption algorithms
  ea = *(buffer + decoded);
  this->ea0 = (ea >> 7) & 0x1;
  this->ea1 = (ea >> 6) & 0x1;
  this->ea2 = (ea >> 5) & 0x1;
  this->ea3 = (ea >> 4) & 0x1;
  this->ea4 = (ea >> 3) & 0x1;
  this->ea5 = (ea >> 2) & 0x1;
  this->ea6 = (ea >> 1) & 0x1;
  this->ea7 = ea & 0x1;
  decoded++;

  // 5GS integrity algorithm
  ia = *(buffer + decoded);
  this->ia0 = (ia >> 7) & 0x1;
  this->ia1 = (ia >> 6) & 0x1;
  this->ia2 = (ia >> 5) & 0x1;
  this->ia3 = (ia >> 4) & 0x1;
  this->ia4 = (ia >> 3) & 0x1;
  this->ia5 = (ia >> 2) & 0x1;
  this->ia6 = (ia >> 1) & 0x1;
  this->ia7 = ia & 0x1;
  decoded++;

  // If any optional buffers are present skip it.
  // 2 = 1 Byte for type + 1 Byte for length

  if (this->length > (decoded - (type_len + length_len))) {
    // 5GS encryption algorithms
    this->eea0 = (*(buffer + decoded) >> 7) & 0x1;
    this->ea1_128 = (*(buffer + decoded) >> 6) & 0x1;
    this->ea2_128 = (*(buffer + decoded) >> 5) & 0x1;
    this->ea3_128 = (*(buffer + decoded) >> 4) & 0x1;
    this->eea4 = (*(buffer + decoded) >> 3) & 0x1;
    this->eea5 = (*(buffer + decoded) >> 2) & 0x1;
    this->eea6 = (*(buffer + decoded) >> 1) & 0x1;
    this->eea7 = *(buffer + decoded) & 0x1;
    decoded++;

    // 5GS integrity algorithm
    this->eia0 = (*(buffer + decoded) >> 7) & 0x1;
    this->eia1_128 = (*(buffer + decoded) >> 6) & 0x1;
    this->eia2_128 = (*(buffer + decoded) >> 5) & 0x1;
    this->eia3_128 = (*(buffer + decoded) >> 4) & 0x1;
    this->eia4 = (*(buffer + decoded) >> 3) & 0x1;
    this->eia5 = (*(buffer + decoded) >> 2) & 0x1;
    this->eia6 = (*(buffer + decoded) >> 1) & 0x1;
    this->eia7 = *(buffer + decoded) & 0x1;
    decoded++;

    // Skipping the remaining bytes as not supported
    if (this->length > (decoded - (type_len + length_len))) {
      // The length of the bytes skipped
      uint8_t additional_ue_sec_cap_length =
          this->length - (decoded - (type_len + length_len));
      decoded += additional_ue_sec_cap_length;
      this->length -= additional_ue_sec_cap_length;
    }
  }

  return (decoded);
};

// Encode UE Security Capability
int UESecurityCapabilityMsg::EncodeUESecurityCapabilityMsg(uint8_t iei,
                                                           uint8_t* buffer,
                                                           uint32_t len) {
  int encoded = 0;

  // Checking IEI and pointer
  if (iei > 0) {
    CHECK_IEI_ENCODER((unsigned char)iei, this->iei);
    *(buffer + encoded++) = iei;
  }

  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer,
                                       UE_SECURITY_CAPABILITY_MIN_LENGTH, len);

  *(buffer + encoded++) = this->length;

  // 5GS encryption algorithms
  *(buffer + encoded++) = 0x00 | ((this->ea0 & 0x1) << 7) |
                          ((this->ea1 & 0x1) << 6) | ((this->ea2 & 0x1) << 5) |
                          ((this->ea3 & 0x1) << 4) | ((this->ea4 & 0x1) << 3) |
                          ((this->ea5 & 0x1) << 2) | ((this->ea6 & 0x1) << 1) |
                          ((this->ea7) & 0x1);

  // 5GS integrity algorithms
  *(buffer + encoded++) = 0x00 | ((this->ia0 & 0x1) << 7) |
                          ((this->ia1 & 0x1) << 6) | ((this->ia2 & 0x1) << 5) |
                          ((this->ia3 & 0x1) << 4) | ((this->ia4 & 0x1) << 3) |
                          ((this->ia5 & 0x1) << 2) | ((this->ia6 & 0x1) << 1) |
                          ((this->ia7) & 0x1);

  if (this->length > 2) {
    // 5GS encryption algorithms
    *(buffer + encoded++) =
        0x00 | ((this->eea0 & 0x1) << 7) | ((this->ea1_128 & 0x1) << 6) |
        ((this->ea2_128 & 0x1) << 5) | ((this->ea3_128 & 0x1) << 4) |
        ((this->eea4 & 0x1) << 3) | ((this->eea5 & 0x1) << 2) |
        ((this->eea6 & 0x1) << 1) | ((this->eea7) & 0x1);

    // 5GS integrity algorithms
    *(buffer + encoded++) =
        0x00 | ((this->eia0 & 0x1) << 7) | ((this->eia1_128 & 0x1) << 6) |
        ((this->eia2_128 & 0x1) << 5) | ((this->eia3_128 & 0x1) << 4) |
        ((this->eia4 & 0x1) << 3) | ((this->eia5 & 0x1) << 2) |
        ((this->eia6 & 0x1) << 1) | ((this->eia7) & 0x1);
  }

  return encoded;
};
}  // namespace magma5g
