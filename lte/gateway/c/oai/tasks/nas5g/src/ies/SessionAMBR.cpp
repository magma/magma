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
#include "SessionAMBR.h"
#include "CommonDefs.h"

using namespace std;
namespace magma5g {
SessionAMBRMsg::SessionAMBRMsg(){};
SessionAMBRMsg::~SessionAMBRMsg(){};

// Decode SessionAMBR IE
int SessionAMBRMsg::DecodeSessionAMBRMsg(
    SessionAMBRMsg* sessionambr, uint8_t iei, uint8_t* buffer, uint32_t len) {
  int decoded = 0;
  // Not yet Implemented, will be supported POST MVC
  return (decoded);
};

// Encode SessionAMBR IE
int SessionAMBRMsg::EncodeSessionAMBRMsg(
    SessionAMBRMsg* sessionambr, uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint16_t* lenPtr;
  uint32_t encoded = 0;

  // Checking IEI and pointer
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, AMBR_MIN_LEN, len);

  if (iei > 0) {
    CHECK_IEI_ENCODER((unsigned char) iei, sessionambr->iei);
    *buffer = iei;
    MLOG(MDEBUG) << "In EncodeSessionAMBRMsg: iei" << hex << int(*buffer)
                 << endl;
    encoded++;
  }

  lenPtr              = (uint16_t*) (buffer + encoded);
  *(buffer + encoded) = sessionambr->length;
  encoded++;
  *(buffer + encoded) = sessionambr->dlunit;
  encoded++;
  IES_ENCODE_U16(buffer, encoded, sessionambr->dlsessionambr);
  *(buffer + encoded) = sessionambr->ulunit;
  encoded++;
  IES_ENCODE_U16(buffer, encoded, sessionambr->ulsessionambr);
  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);

  return (encoded);
};
}  // namespace magma5g
