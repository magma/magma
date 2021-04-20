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
#include "M5GIMEISVRequest.h"
#include "M5GCommonDefs.h"

using namespace std;
namespace magma5g {
ImeisvRequestMsg::ImeisvRequestMsg(){};
ImeisvRequestMsg::~ImeisvRequestMsg(){};

int ImeisvRequestMsg::DecodeImeisvRequestMsg(
    ImeisvRequestMsg* imeisv_request, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded = 0;

  MLOG(MDEBUG) << "DecoseImeisvRequestMsg : ";

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, IMEISV_REQUEST_MINIMUM_LENGTH, len);

  if (iei > 0) {
    CHECK_IEI_DECODER((unsigned char) (*buffer & 0xf0), iei);
  }

  imeisv_request->spare          = (*(buffer + decoded) >> 7) & 0x1;
  imeisv_request->imeisv_request = (*(buffer + decoded) >> 4) & 0x7;
  decoded++;
  MLOG(MDEBUG) << "   spare = " << dec << int(imeisv_request->spare);
  MLOG(MDEBUG) << "   imeisv request = " << dec
               << int(imeisv_request->imeisv_request);
  return decoded;
};

int ImeisvRequestMsg::EncodeImeisvRequestMsg(
    ImeisvRequestMsg* imeisv_request, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint32_t encoded = 0;

  MLOG(MDEBUG) << " EncodeImeisvRequestMsg : " << endl;
  *(buffer + encoded) = 0xe0 | (imeisv_request->spare & 0x1) << 3 |
                        (imeisv_request->imeisv_request & 0x7);
  encoded++;

  return encoded;
};
}  // namespace magma5g
