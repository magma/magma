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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GIMEISVRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
ImeisvRequestMsg::ImeisvRequestMsg(){};
ImeisvRequestMsg::~ImeisvRequestMsg(){};

int ImeisvRequestMsg::DecodeImeisvRequestMsg(ImeisvRequestMsg* imeisv_request,
                                             uint8_t iei, uint8_t* buffer,
                                             uint32_t len) {
  int decoded = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, IMEISV_REQUEST_MINIMUM_LENGTH, len);

  OAILOG_DEBUG(LOG_NAS5G, "Decoding ImeisvRequest");
  if (iei > 0) {
    CHECK_IEI_DECODER((unsigned char) (*buffer & 0xf0), iei);
    OAILOG_DEBUG(LOG_NAS5G, "IEI %X", static_cast<int>(iei));
  }

  imeisv_request->spare = (*(buffer + decoded) >> 7) & 0x1;
  imeisv_request->imeisv_request = (*(buffer + decoded) >> 4) & 0x7;
  decoded++;
  OAILOG_DEBUG(
      LOG_NAS5G, "Spare : %d", static_cast<int>(imeisv_request->spare));
  OAILOG_DEBUG(
      LOG_NAS5G, "IMEISV request : %d",
      static_cast<int>(imeisv_request->imeisv_request));
  return decoded;
};

int ImeisvRequestMsg::EncodeImeisvRequestMsg(ImeisvRequestMsg* imeisv_request,
                                             uint8_t iei, uint8_t* buffer,
                                             uint32_t len) {
  uint32_t encoded = 0;

  OAILOG_DEBUG(LOG_NAS5G, "Encoding ImeisvRequest");
  *(buffer + encoded) = 0xe0 | (imeisv_request->spare & 0x1) << 3 |
                        (imeisv_request->imeisv_request & 0x7);
  OAILOG_DEBUG(
      LOG_NAS5G, "[Spare, IMEISV request] : %X",
      static_cast<int>(*(buffer + encoded)));
  encoded++;

  return encoded;
};
}  // namespace magma5g
