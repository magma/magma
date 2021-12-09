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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GNASKeySetIdentifier.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
NASKeySetIdentifierMsg::NASKeySetIdentifierMsg(){};
NASKeySetIdentifierMsg::~NASKeySetIdentifierMsg(){};

// Decode NASKeySetIdentifier IE
int NASKeySetIdentifierMsg::DecodeNASKeySetIdentifierMsg(
    NASKeySetIdentifierMsg* nas_key_set_identifier, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  int decoded = 0;

  OAILOG_DEBUG(LOG_NAS5G, "Decoding NASKeySetIdentifier");

  // Checking IEI and pointer
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(buffer,
                                       NAS_KEY_SET_IDENTIFIER_MIN_LENGTH, len);

  if (iei > 0) {
    CHECK_IEI_DECODER((unsigned char)(*buffer & 0xf0), iei);
  }

  nas_key_set_identifier->tsc = (*(buffer + decoded) >> 7) & 0x1;
  nas_key_set_identifier->nas_key_set_identifier =
      (*(buffer + decoded) >> 4) & 0x7;
  decoded++;
  OAILOG_DEBUG(
      LOG_NAS5G, "Type of Security Context : %d",
      static_cast<int>(nas_key_set_identifier->tsc));
  OAILOG_DEBUG(
      LOG_NAS5G, "NASkeysetidentifier : %d",
      static_cast<int>(nas_key_set_identifier->nas_key_set_identifier));
  return decoded;
};

// Encode NASKeySetIdentifier IE
int NASKeySetIdentifierMsg::EncodeNASKeySetIdentifierMsg(
    NASKeySetIdentifierMsg* nas_key_set_identifier, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

  // Checking IEI and pointer
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer,
                                       NAS_KEY_SET_IDENTIFIER_MIN_LENGTH, len);

  OAILOG_DEBUG(LOG_NAS5G, "Encoding NASKeySetIdentifier");
  if (iei > 0) {
    CHECK_IEI_ENCODER((unsigned char)iei, nas_key_set_identifier->iei);
    *buffer = iei;
    OAILOG_DEBUG(LOG_NAS5G, "IEI : %X", static_cast<int>(*buffer));
    encoded++;
  }

  *(buffer + encoded) = 0x00 | (nas_key_set_identifier->tsc & 0x1) << 3 |
                        (nas_key_set_identifier->nas_key_set_identifier & 0x7);
  OAILOG_DEBUG(
      LOG_NAS5G, "Type of Security Context  = 0x%X",
      static_cast<int>(nas_key_set_identifier->tsc));
  OAILOG_DEBUG(
      LOG_NAS5G, "NASkeysetidentifier = 0x%X",
      static_cast<int>(*(buffer + encoded)));
  encoded++;

  return encoded;
};
}  // namespace magma5g
