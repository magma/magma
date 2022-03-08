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
namespace magma5g {
class M5GUplinkDataStatus {
 public:
  uint8_t iei;
  uint8_t len;
  uint16_t uplinkDataStatus;

  M5GUplinkDataStatus();
  ~M5GUplinkDataStatus();
  int EncodeUplinkDataStatus(M5GUplinkDataStatus* uplinkDataStatus, uint8_t iei,
                             uint8_t* buffer, uint32_t len);
  int DecodeUplinkDataStatus(M5GUplinkDataStatus* uplinkDataStatus, uint8_t iei,
                             uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g
