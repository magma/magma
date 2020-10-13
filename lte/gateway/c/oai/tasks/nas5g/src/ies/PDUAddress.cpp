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
#include "PDUAddress.h"
#include "CommonDefs.h"

using namespace std;
namespace magma5g {
PDUAddressMsg::PDUAddressMsg(){};
PDUAddressMsg::~PDUAddressMsg(){};

// Decode PDUAddress IE
int PDUAddressMsg::DecodePDUAddressMsg(
    PDUAddressMsg* pduaddress, uint8_t iei, uint8_t* buffer, uint32_t len) {
  // Not yet Implemented, will be supported POST MVC
  return 0;
};

// Encode PDUAddress IE
int PDUAddressMsg::EncodePDUAddressMsg(
    PDUAddressMsg* pduaddress, uint8_t iei, uint8_t* buffer, uint32_t len) {
  int encoded = 0;
  uint8_t* lenPtr;

  if (iei > 0) {
    pduaddress->iei = (*buffer & 0xf0) >> 4;
    CHECK_IEI_DECODER(iei, (unsigned char) pduaddress->iei);
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;
  *(buffer + encoded) = 0x00 | (pduaddress->typeval & 0x07);
  MLOG(MDEBUG) << "EncodePDUAddressMsg__: typeval = " << hex
               << int(pduaddress->typeval) << endl;
  encoded++;
  memcpy(buffer + encoded, pduaddress->addressinfo, 12);
  encoded = encoded + 12;

  return (encoded);
};
}  // namespace magma5g
