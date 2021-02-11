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

using namespace std;
namespace magma5g {
class M5GSDeRegistrationTypeMsg {
 public:
  uint8_t switchoff : 1;
  uint8_t re_reg_required : 1;
  uint8_t access_type : 2;

  M5GSDeRegistrationTypeMsg();
  ~M5GSDeRegistrationTypeMsg();
  int DecodeM5GSDeRegistrationTypeMsg(
      M5GSDeRegistrationTypeMsg* m5gs_de_reg_type, uint8_t iei, uint8_t* buffer,
      uint32_t len);
  int EncodeM5GSDeRegistrationTypeMsg(
      M5GSDeRegistrationTypeMsg* m5gs_de_reg_type, uint8_t iei, uint8_t* buffer,
      uint32_t len);
};
}  // namespace magma5g
