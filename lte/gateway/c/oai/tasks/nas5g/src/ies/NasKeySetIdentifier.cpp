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
#include <bitset>
#include <cstdint>
#include "NASKeySetIdentifier.h"
#include "CommonDefs.h"
using namespace std;

namespace magma5g
{
  NASKeySetIdentifierMsg::NASKeySetIdentifierMsg()
  {
  };

  NASKeySetIdentifierMsg::~NASKeySetIdentifierMsg()
  {
  };

  // Decode NASKeySetIdentifier IE
  int NASKeySetIdentifierMsg::DecodeNASKeySetIdentifierMsg(NASKeySetIdentifierMsg *naskeysetidentifier, uint8_t iei, uint8_t *buffer, uint32_t len) 
  {
    int decoded = 0;

    MLOG(MDEBUG) << "   DecodeNASKeySetIdentifierMsg : "<<"\n";

    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, NAS_KEY_SET_IDENTIFIER_MIN_LENGTH, len);

    if (iei > 0) {
      CHECK_IEI_DECODER((unsigned char)(*buffer & 0xf0), iei);
    }

    naskeysetidentifier->tsc                 = (*(buffer + decoded) >> 7) & 0x1;
    naskeysetidentifier->naskeysetidentifier = (*(buffer + decoded) >> 4) & 0x7;
    decoded++;
    MLOG(MDEBUG) << "      tsc = 0x" << hex << bitset<4>(int(naskeysetidentifier->tsc))<<"\n";
    MLOG(MDEBUG) << "      naskeysetidentifier = 0x" << hex << bitset<3>(int(naskeysetidentifier->naskeysetidentifier))<<"\n";
    return decoded;
  };

  // Encode NASKeySetIdentifier IE
  int NASKeySetIdentifierMsg::EncodeNASKeySetIdentifierMsg(NASKeySetIdentifierMsg *naskeysetidentifier, uint8_t iei, uint8_t * buffer, uint32_t len)
  {
    uint32_t encoded = 0;

    MLOG(MDEBUG) << "EncodeNASKeySetIdentifierMsg:";
    // TBD 
    CHECK_PDU_POINTER_AND_LENGTH_ENCODER (buffer,NAS_KEY_SET_IDENTIFIER_MIN_LENGTH , len);

    return encoded;
  };
}

