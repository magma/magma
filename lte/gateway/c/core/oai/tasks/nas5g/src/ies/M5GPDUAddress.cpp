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
#include <sstream>
#include <cstdint>
#include <cstring>
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPDUAddress.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/include/nas/networkDef.h"
#ifdef __cplusplus
}
#endif

namespace magma5g {
PDUAddressMsg::PDUAddressMsg() {};
PDUAddressMsg::~PDUAddressMsg() {};

// Decode PDUAddress IE
int PDUAddressMsg::DecodePDUAddressMsg(PDUAddressMsg* pdu_address, uint8_t iei,
                                       uint8_t* buffer, uint32_t len) {
  uint8_t decoded = 0;
  // CHECKING IEI
  if (iei > 0) {
    IES_DECODE_U8(buffer, decoded, pdu_address->iei);
    CHECK_IEI_DECODER(iei, (unsigned char)pdu_address->iei);
  }

  IES_DECODE_U8(buffer, decoded, pdu_address->length);

  pdu_address->type_val = *(buffer + decoded) && 0x07;
  memset(pdu_address->address_info, 0, sizeof(pdu_address->address_info));
  decoded++;
  memcpy(buffer + decoded, pdu_address->address_info, pdu_address->length - 1);
  decoded += pdu_address->length - 1;

  return (decoded);
};

// Encode PDUAddress IE
int PDUAddressMsg::EncodePDUAddressMsg(PDUAddressMsg* pdu_address, uint8_t iei,
                                       uint8_t* buffer, uint32_t len) {
  int encoded = 0;

  // CHECKING IEI
  if (iei > 0) {
    pdu_address->iei = iei;
    CHECK_IEI_DECODER(iei, (unsigned char)pdu_address->iei);
    *(buffer + encoded) = iei;
    encoded++;
  }

  // Sizeof type_val + address length
  IES_ENCODE_U8(buffer, encoded, sizeof(uint8_t) + pdu_address->length);
  IES_ENCODE_U8(buffer, encoded, (0x00 | (pdu_address->type_val & 0x07)));
  memcpy(buffer + encoded, pdu_address->address_info, pdu_address->length);
  encoded = encoded + pdu_address->length;

  return (encoded);
};
}  // namespace magma5g
