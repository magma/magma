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

#pragma once
#include <sstream>
#include <cstdint>

using namespace std;
namespace magma5g {
class UESecurityCapabilityMsg {
 public:
#define UE_SECURITY_CAPABILITY_MIN_LENGTH 1
  uint8_t length;
  uint8_t iei;
  uint8_t ea0 : 1;
  uint8_t ea1 : 1;
  uint8_t ea2 : 1;
  uint8_t ea3 : 1;
  uint8_t ea4 : 1;
  uint8_t ea5 : 1;
  uint8_t ea6 : 1;
  uint8_t ea7 : 1;
  uint8_t ia0 : 1;
  uint8_t ia1 : 1;
  uint8_t ia2 : 1;
  uint8_t ia3 : 1;
  uint8_t ia4 : 1;
  uint8_t ia5 : 1;
  uint8_t ia6 : 1;
  uint8_t ia7 : 1;
  uint8_t spare[3];

  UESecurityCapabilityMsg();
  ~UESecurityCapabilityMsg();
  int EncodeUESecurityCapabilityMsg(
      UESecurityCapabilityMsg* ue_sec_capability, uint8_t iei, uint8_t* buffer,
      uint32_t len);
  int DecodeUESecurityCapabilityMsg(
      UESecurityCapabilityMsg* ue_sec_capability, uint8_t iei, uint8_t* buffer,
      uint32_t len);
};
}  // namespace magma5g
