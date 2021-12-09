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
#include <cstring>
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GServiceType.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
ServiceTypeMsg::ServiceTypeMsg(){};
ServiceTypeMsg::~ServiceTypeMsg(){};

// Decode ServiceType IE
int ServiceTypeMsg::DecodeServiceTypeMsg(ServiceTypeMsg* svc_type, uint8_t iei,
                                         uint8_t* buffer, uint32_t len) {
  int decoded = 0;

  OAILOG_DEBUG(LOG_NAS5G, "Decoding Service Type");
  svc_type->service_type_value = ((*buffer & 0xf0) >> 4);
  OAILOG_DEBUG(
      LOG_NAS5G, "Service Type Value : %X",
      static_cast<int>(svc_type->service_type_value));

  return (decoded);
};

// Encode ServiceType IE
int ServiceTypeMsg::EncodeServiceTypeMsg(ServiceTypeMsg* svc_type, uint8_t iei,
                                         uint8_t* buffer, uint32_t len) {
  int encoded = 0;

  OAILOG_DEBUG(LOG_NAS5G, "Encoding Service Type");
  *buffer = svc_type->service_type_value & 0x0f;
  OAILOG_DEBUG(LOG_NAS5G, "Service Type Value : %X", static_cast<int>(*buffer));
  encoded++;

  return (encoded);
};
}  // namespace magma5g
