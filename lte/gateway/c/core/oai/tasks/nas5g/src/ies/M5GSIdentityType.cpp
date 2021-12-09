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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSIdentityType.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
M5GSIdentityTypeMsg::M5GSIdentityTypeMsg(){};
M5GSIdentityTypeMsg::~M5GSIdentityTypeMsg(){};

// Decode M5GSIdentityType IE
int M5GSIdentityTypeMsg::DecodeM5GSIdentityTypeMsg(
    M5GSIdentityTypeMsg* m5gs_identity_type, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint8_t decoded = 0;

  OAILOG_DEBUG(LOG_NAS5G, "Decoding 5GSIdentityType");
  m5gs_identity_type->toi = *(buffer + decoded) & 0x7;
  decoded++;
  OAILOG_DEBUG(
      LOG_NAS5G, "Type of Identity : %d",
      static_cast<int>(m5gs_identity_type->toi));
  return (decoded);
};

// Encode M5GSIdentityType IE
int M5GSIdentityTypeMsg::EncodeM5GSIdentityTypeMsg(
    M5GSIdentityTypeMsg* m5gs_identity_type, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int encoded = 0;

  OAILOG_DEBUG(LOG_NAS5G, "Encoding 5GSIdentityType");
  *(buffer + encoded) = (m5gs_identity_type->toi) & 0x7;
  OAILOG_DEBUG(
      LOG_NAS5G, "Type of Identity : %X",
      static_cast<int>(*(buffer + encoded)));
  encoded++;
  return (encoded);
};
}  // namespace magma5g
