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
#include "M5GNASKeySetIdentifier.h"
#include "M5GCommonDefs.h"

using namespace std;
namespace magma5g {
NASKeySetIdentifierMsg::NASKeySetIdentifierMsg(){};
NASKeySetIdentifierMsg::~NASKeySetIdentifierMsg(){};

// Decode NASKeySetIdentifier IE
int NASKeySetIdentifierMsg::DecodeNASKeySetIdentifierMsg(
    NASKeySetIdentifierMsg* nas_key_set_identifier, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  int decoded = 0;

  MLOG(MDEBUG) << "DecoseNASKeySetIdentifierMsg : ";

  // Checking IEI and pointer
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, NAS_KEY_SET_IDENTIFIER_MIN_LENGTH, len);

  if (iei > 0) {
    CHECK_IEI_DECODER((unsigned char) (*buffer & 0xf0), iei);
  }

  nas_key_set_identifier->tsc = (*(buffer + decoded) >> 7) & 0x1;
  nas_key_set_identifier->nas_key_set_identifier =
      (*(buffer + decoded) >> 4) & 0x7;
  decoded++;
  MLOG(MDEBUG) << "   tsc = " << dec << int(nas_key_set_identifier->tsc);
  MLOG(MDEBUG) << "   NASkeysetidentifier = " << dec
               << int(nas_key_set_identifier->nas_key_set_identifier);
  return decoded;
};

// Encode NASKeySetIdentifier IE
int NASKeySetIdentifierMsg::EncodeNASKeySetIdentifierMsg(
    NASKeySetIdentifierMsg* nas_key_set_identifier, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

  // Checking IEI and pointer
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, NAS_KEY_SET_IDENTIFIER_MIN_LENGTH, len);

  if (iei > 0) {
    CHECK_IEI_ENCODER((unsigned char) iei, nas_key_set_identifier->iei);
    *buffer = iei;
    MLOG(MDEBUG) << "In EncodeNASKeySetIdentifierMsg: iei" << hex
                 << int(*buffer) << endl;
    encoded++;
  }

  MLOG(MDEBUG) << " EncodeNASKeySetIdentifierMsg : " << endl;
  *(buffer + encoded) = 0x00 | (nas_key_set_identifier->tsc & 0x1) << 3 |
                        (nas_key_set_identifier->nas_key_set_identifier & 0x7);
  MLOG(MDEBUG) << "   Type of Security Context  = 0x" << hex
               << int(nas_key_set_identifier->tsc) << "\n";
  MLOG(MDEBUG) << "   NAS key set identifier = 0x" << hex
               << int(*(buffer + encoded)) << "\n";
  encoded++;

  return encoded;
};
}  // namespace magma5g
