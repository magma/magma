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
#include "M5GSRegistrationType.h"
#include "M5GCommonDefs.h"
#include <bitset>

using namespace std;
namespace magma5g {
M5GSRegistrationTypeMsg::M5GSRegistrationTypeMsg(){};
M5GSRegistrationTypeMsg::~M5GSRegistrationTypeMsg(){};

// Decode M5GSRegistrationType Message
int M5GSRegistrationTypeMsg::DecodeM5GSRegistrationTypeMsg(
    M5GSRegistrationTypeMsg* m5gs_reg_type, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded = 0;

  // CHECKING IEI
  if (iei > 0) {
    CHECK_IEI_DECODER((*buffer & 0xf0), iei);
  }

  m5gs_reg_type->FOR      = (*(buffer + decoded) >> 3) & 0x1;
  m5gs_reg_type->type_val = *(buffer + decoded) & 0x7;
  MLOG(MDEBUG) << " FOR = 0x" << hex << int(m5gs_reg_type->FOR);
  MLOG(MDEBUG) << " type_val = 0x" << hex << int(m5gs_reg_type->type_val);
  return decoded;
};

// Encode M5GSRegistrationType Message
int M5GSRegistrationTypeMsg::EncodeM5GSRegistrationTypeMsg(
    M5GSRegistrationTypeMsg* m5gs_reg_type, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  // Will be supported POST MVC
  return 0;
};
}  // namespace magma5g
