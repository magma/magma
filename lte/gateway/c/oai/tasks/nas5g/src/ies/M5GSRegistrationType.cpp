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
#include "5GSRegistrationType.h"
#include "CommonDefs.h"
#include <bitset>

using namespace std;
namespace magma5g
{
  M5GSRegistrationTypeMsg::M5GSRegistrationTypeMsg()
  {
  };

  M5GSRegistrationTypeMsg::~M5GSRegistrationTypeMsg()
  {
  };
  // Decode M5GSRegistrationType Message 
  int M5GSRegistrationTypeMsg::DecodeM5GSRegistrationTypeMsg(M5GSRegistrationTypeMsg *m5gsregistrationtype, uint8_t iei, uint8_t *buffer, uint32_t len) 
  {
    int decoded = 0;

    MLOG(MDEBUG) << "   DecodeM5GSRegistrationTypeMsg : "<<"\n";
    CHECK_PDU_POINTER_AND_LENGTH_DECODER(
        buffer, REGISTRATION_TYPE_MIN_LENGTH, len);

    if (iei > 0) {
      CHECK_IEI_DECODER((*buffer & 0xf0), iei);
    }

    m5gsregistrationtype->FOR = (*(buffer + decoded) >> 3) & 0x1;
    m5gsregistrationtype->typeval = *(buffer + decoded) & 0x7;
    MLOG(MDEBUG) << "      FOR = 0x" << hex << bitset<4>(int(m5gsregistrationtype->FOR))<<"\n";
    MLOG(MDEBUG) << "      typeval = 0x" << hex << bitset<3>(int(m5gsregistrationtype->typeval))<<"\n";
    return decoded;
   };

   // Encode M5GSRegistrationType Message 
   int M5GSRegistrationTypeMsg::EncodeM5GSRegistrationTypeMsg(M5GSRegistrationTypeMsg *m5gsregistrationtype, uint8_t iei, uint8_t * buffer, uint32_t len)
   {
      // TBD
      return 0;
   };
}

