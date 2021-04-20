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
class M5GSRegistrationResultMsg {
 public:
  uint8_t iei;
  const int REGISTRATION_RESULT_MIN_LENGTH = 2;
  uint8_t spare : 4;
  uint8_t sms_allowed : 1;
  uint8_t reg_result_val : 3;

  M5GSRegistrationResultMsg();
  ~M5GSRegistrationResultMsg();
  int EncodeM5GSRegistrationResultMsg(
      M5GSRegistrationResultMsg* m5gs_reg_result, uint8_t iei, uint8_t* buffer,
      uint32_t len);
  int DecodeM5GSRegistrationResultMsg(
      M5GSRegistrationResultMsg* m5gs_reg_result, uint8_t iei, uint8_t* buffer,
      uint32_t len);
};
}  // namespace magma5g
