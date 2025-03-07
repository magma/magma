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
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPDUSessionType.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
PDUSessionTypeMsg::PDUSessionTypeMsg() {};
PDUSessionTypeMsg::~PDUSessionTypeMsg() {};

// Decode PDUSessionType IE
int PDUSessionTypeMsg::DecodePDUSessionTypeMsg(
    PDUSessionTypeMsg* pdu_session_type, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded = 0;

  // Store the IEI Information
  if (iei > 0) {
    pdu_session_type->iei = (*buffer & 0xf0) >> 4;
    decoded++;
  }

  pdu_session_type->type_val = (*buffer & 0x07);

  return (decoded);
};

// Encode PDUSessionType IE
int PDUSessionTypeMsg::EncodePDUSessionTypeMsg(
    PDUSessionTypeMsg* pdu_session_type, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int encoded = 0;

  // CHECKING IEI
  if (iei > 0) {
    *buffer = (pdu_session_type->iei & 0x0f) << 4;
    CHECK_IEI_ENCODER((uint8_t)iei,
                      (uint8_t)((pdu_session_type->iei & 0x0f) << 4));
  }

  *buffer = 0x00 | (*buffer & 0xf0) | (pdu_session_type->type_val & 0x07);
  encoded++;

  return (encoded);
};
}  // namespace magma5g
