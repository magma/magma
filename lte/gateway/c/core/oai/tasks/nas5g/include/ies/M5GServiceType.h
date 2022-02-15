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
class ServiceTypeMsg {
 public:
  uint8_t iei : 4;

#define SERVICE_TYPE_SIGNALING 0X00
#define SERVICE_TYPE_DATA 0X01
#define SERVICE_TYPE_MOBILE_TERMINATED_SERVICES 0X02
#define SERVICE_TYPE_EMERGENCY_SERVICE 0X03
#define SERVICE_TYPE_EMERGENCY_SERVICE_FALL_BACK 0X04
#define SERVICE_TYPE_HIGH_PRIORITY_ACCESS 0X05
  uint8_t service_type_value : 4;

  ServiceTypeMsg();
  ~ServiceTypeMsg();
  int EncodeServiceTypeMsg(ServiceTypeMsg* service_type, uint8_t iei,
                           uint8_t* buffer, uint32_t len);
  int DecodeServiceTypeMsg(ServiceTypeMsg* service_type, uint8_t iei,
                           uint8_t* buffer, uint32_t len);
};
}  // namespace magma5g
