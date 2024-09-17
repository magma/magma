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

#include <sstream>
#include <cstdint>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSessionAMBR.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
SessionAMBRMsg::SessionAMBRMsg() {}
SessionAMBRMsg::~SessionAMBRMsg() {}

// Decode SessionAMBR IE
int SessionAMBRMsg::DecodeSessionAMBRMsg(SessionAMBRMsg* session_ambr,
                                         uint8_t iei, uint8_t* buffer,
                                         uint32_t len) {
  int decoded = 0;
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(buffer, AMBR_MIN_LEN, len);

  if (iei > 0) {
    session_ambr->iei = *buffer;
    CHECK_IEI_DECODER((unsigned char)iei, session_ambr->iei);
    decoded++;
  }

  IES_DECODE_U8(buffer, decoded, session_ambr->length);
  IES_DECODE_U8(buffer, decoded, session_ambr->dl_unit);
  IES_DECODE_U16(buffer, decoded, session_ambr->dl_session_ambr);
  IES_DECODE_U8(buffer, decoded, session_ambr->ul_unit);
  IES_DECODE_U16(buffer, decoded, session_ambr->ul_session_ambr);
  return decoded;
}

// Encode SessionAMBR IE
int SessionAMBRMsg::EncodeSessionAMBRMsg(SessionAMBRMsg* session_ambr,
                                         uint8_t iei, uint8_t* buffer,
                                         uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;

  // Checking IEI and pointer
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, AMBR_MIN_LEN, len);

  if (iei > 0) {
    CHECK_IEI_ENCODER((unsigned char)iei, session_ambr->iei);
    *buffer = iei;
    encoded++;
  }

  lenPtr = reinterpret_cast<uint8_t*>(buffer + encoded);
  *(buffer + encoded) = session_ambr->length;
  encoded++;

  *(buffer + encoded) = session_ambr->dl_unit;
  encoded++;

  IES_ENCODE_U16(buffer, encoded, session_ambr->dl_session_ambr);

  *(buffer + encoded) = session_ambr->ul_unit;
  encoded++;

  IES_ENCODE_U16(buffer, encoded, session_ambr->ul_session_ambr);
  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);

  return encoded;
}
}  // namespace magma5g
