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
#include "M5GSIdentityType.h"
#include "M5GCommonDefs.h"

using namespace std;
namespace magma5g {
M5GSIdentityTypeMsg::M5GSIdentityTypeMsg(){};
M5GSIdentityTypeMsg::~M5GSIdentityTypeMsg(){};

// Decode M5GSIdentityType IE
int M5GSIdentityTypeMsg::DecodeM5GSIdentityTypeMsg(
    M5GSIdentityTypeMsg* m5gs_identity_type, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint8_t decoded = 0;

  MLOG(MDEBUG) << "   DecodeM5GSIdentityTypeMsg : ";
  m5gs_identity_type->toi = *(buffer + decoded) & 0x7;
  decoded++;
  MLOG(MDEBUG) << " Type of Identity = " << dec << int(m5gs_identity_type->toi);
  return (decoded);
};

// Encode M5GSIdentityType IE
int M5GSIdentityTypeMsg::EncodeM5GSIdentityTypeMsg(
    M5GSIdentityTypeMsg* m5gs_identity_type, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int encoded = 0;

  MLOG(MDEBUG) << " EncodeM5GSIdentityTypeMsg : ";
  *(buffer + encoded) = (m5gs_identity_type->toi) & 0x7;
  MLOG(MDEBUG) << " Type of identity = 0x" << hex << int(*(buffer + encoded));
  encoded++;
  return (encoded);
};
}  // namespace magma5g
