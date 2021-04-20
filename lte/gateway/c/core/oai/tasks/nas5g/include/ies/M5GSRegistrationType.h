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
// M5GSRegistrationType IE Class
class M5GSRegistrationTypeMsg {
 public:
  uint8_t iei;
  const int REGISTRATION_TYPE_MIN_LENGTH = 1;
  uint8_t FOR : 1;
  uint8_t type_val : 3;

  M5GSRegistrationTypeMsg();
  ~M5GSRegistrationTypeMsg();
  int EncodeM5GSRegistrationTypeMsg(
      M5GSRegistrationTypeMsg* m5gs_reg_type, uint8_t iei, uint8_t* buffer,
      uint32_t len);
  int DecodeM5GSRegistrationTypeMsg(
      M5GSRegistrationTypeMsg* m5gs_reg_type, uint8_t iei, uint8_t* buffer,
      uint32_t len);
};
}  // namespace magma5g
