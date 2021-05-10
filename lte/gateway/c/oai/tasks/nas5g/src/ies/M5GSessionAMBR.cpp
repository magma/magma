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
#include "M5GSessionAMBR.h"
#include "M5GCommonDefs.h"

using namespace std;
namespace magma5g {
SessionAMBRMsg::SessionAMBRMsg(){};
SessionAMBRMsg::~SessionAMBRMsg(){};

// Decode SessionAMBR IE
int SessionAMBRMsg::DecodeSessionAMBRMsg(
    SessionAMBRMsg* session_ambr, uint8_t iei, uint8_t* buffer, uint32_t len) {
  int decoded = 0;
  // Not yet Implemented, will be supported POST MVC
  return (decoded);
};

// Encode SessionAMBR IE
int SessionAMBRMsg::EncodeSessionAMBRMsg(
    SessionAMBRMsg* session_ambr, uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint16_t* lenPtr;
  uint32_t encoded = 0;

  // Checking IEI and pointer
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, AMBR_MIN_LEN, len);

  if (iei > 0) {
    CHECK_IEI_ENCODER((unsigned char) iei, session_ambr->iei);
    *buffer = iei;
    MLOG(MDEBUG) << "In EncodeSessionAMBRMsg: iei" << hex << int(*buffer);
    encoded++;
  }

  lenPtr              = (uint16_t*) (buffer + encoded);
  *(buffer + encoded) = session_ambr->length;
  encoded++;
  *(buffer + encoded) = session_ambr->dl_unit;
  encoded++;
  IES_ENCODE_U16(buffer, encoded, session_ambr->dl_session_ambr);
  *(buffer + encoded) = session_ambr->ul_unit;
  encoded++;
  IES_ENCODE_U16(buffer, encoded, session_ambr->ul_session_ambr);
  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);

  return (encoded);
};
}  // namespace magma5g
