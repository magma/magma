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

#include <cstring>
#include <sstream>
#include <cstdint>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPDUSessionIdentity.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"

namespace magma5g {
PDUSessionIdentityMsg::PDUSessionIdentityMsg() {
  memset(this, 0, sizeof(PDUSessionIdentityMsg));
};
PDUSessionIdentityMsg::~PDUSessionIdentityMsg(){};

// Decode PDUSessionIdentity IE
int PDUSessionIdentityMsg::DecodePDUSessionIdentityMsg(uint8_t iei,
                                                       uint8_t* buffer,
                                                       uint32_t len) {
  uint8_t decoded = 0;

  if (iei > 0) {
    this->iei = *(buffer + decoded);
    CHECK_IEI_DECODER((unsigned char)iei, this->iei);
    decoded++;
  }

  this->pdu_session_id = *(buffer + decoded);
  decoded++;

  return (decoded);
};

// Encode PDUSessionIdentity IE
int PDUSessionIdentityMsg::EncodePDUSessionIdentityMsg(uint8_t iei,
                                                       uint8_t* buffer,
                                                       uint32_t len) {
  int encoded = 0;

  // Checking IEI and pointer
  if (iei > 0) {
    CHECK_IEI_ENCODER((unsigned char)iei,
                      static_cast<uint8_t>(M5GIei::PDU_SESSION_IDENTITY_2));
    *buffer = iei;
    encoded++;
  }

  *(buffer + encoded) = this->pdu_session_id;
  encoded++;

  return (encoded);
};
}  // namespace magma5g
