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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GIntegrityProtMaxDataRate.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
IntegrityProtMaxDataRateMsg::IntegrityProtMaxDataRateMsg() {};
IntegrityProtMaxDataRateMsg::~IntegrityProtMaxDataRateMsg() {};

// Decode IntegrityProtMaxDataRate IE
int IntegrityProtMaxDataRateMsg::DecodeIntegrityProtMaxDataRateMsg(
    IntegrityProtMaxDataRateMsg* integrity_prot_max_data_rate, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  uint8_t decoded = 0;

  integrity_prot_max_data_rate->max_uplink = *(buffer + decoded);
  decoded++;
  integrity_prot_max_data_rate->max_downlink = *(buffer + decoded);
  decoded++;
  return (decoded);
};

// Encode IntegrityProtMaxDataRate IE
int IntegrityProtMaxDataRateMsg::EncodeIntegrityProtMaxDataRateMsg(
    IntegrityProtMaxDataRateMsg* integrity_prot_max_data_rate, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  int encoded = 0;

  *(buffer + encoded) = integrity_prot_max_data_rate->max_uplink;
  encoded++;
  *(buffer + encoded) = integrity_prot_max_data_rate->max_downlink;
  encoded++;
  return (encoded);
};
}  // namespace magma5g
