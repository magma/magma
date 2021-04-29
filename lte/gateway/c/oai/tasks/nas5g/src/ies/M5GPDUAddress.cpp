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
#include "M5GPDUAddress.h"
#include "M5GCommonDefs.h"

using namespace std;
namespace magma5g {
PDUAddressMsg::PDUAddressMsg(){};
PDUAddressMsg::~PDUAddressMsg(){};

// Decode PDUAddress IE
int PDUAddressMsg::DecodePDUAddressMsg(
    PDUAddressMsg* pdu_address, uint8_t iei, uint8_t* buffer, uint32_t len) {
  // Not yet Implemented, will be supported POST MVC
  return 0;
};

// Encode PDUAddress IE
int PDUAddressMsg::EncodePDUAddressMsg(
    PDUAddressMsg* pdu_address, uint8_t iei, uint8_t* buffer, uint32_t len) {
  int encoded = 0;

  // CHECKING IEI
  if (iei > 0) {
    pdu_address->iei = iei;
    CHECK_IEI_DECODER(iei, (unsigned char) pdu_address->iei);
    *(buffer + encoded) = iei;
    encoded++;
  }

  if (pdu_address->type_val == TYPE_VAL_IPV4) {
    *(buffer + encoded) = pdu_address->length;
    encoded++;
    *(buffer + encoded) = 0x00 | (pdu_address->type_val & 0x07);
    encoded++;
    memcpy(buffer + encoded, pdu_address->address_info, IPV4_ADDRESS_LENGTH);
    encoded = encoded + IPV4_ADDRESS_LENGTH;
  }
  return (encoded);
};
}  // namespace magma5g
