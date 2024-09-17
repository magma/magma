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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSMCause.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
M5GSMCauseMsg::M5GSMCauseMsg() {}
M5GSMCauseMsg::~M5GSMCauseMsg() {}

// Decode M5GSMCause IE
int M5GSMCauseMsg::DecodeM5GSMCauseMsg(M5GSMCauseMsg* m5gsm_cause, uint8_t iei,
                                       uint8_t* buffer, uint32_t len) {
  int decoded = 0;

  // CHECKING IEI
  if (iei > 0) {
    m5gsm_cause->iei = *buffer;
    CHECK_IEI_DECODER(iei, (unsigned char)m5gsm_cause->iei);
    decoded++;
  }

  m5gsm_cause->cause_value = *(buffer + decoded);
  decoded++;

  return decoded;
}

// Encode M5GSMCause IE
int M5GSMCauseMsg::EncodeM5GSMCauseMsg(M5GSMCauseMsg* m5gsm_cause, uint8_t iei,
                                       uint8_t* buffer, uint32_t len) {
  int encoded = 0;

  // CHECKING IEI
  if (iei > 0) {
    *buffer = m5gsm_cause->iei;
    CHECK_IEI_DECODER(iei, (unsigned char)m5gsm_cause->iei);
    encoded++;
  }

  *(buffer + encoded) = m5gsm_cause->cause_value;
  encoded++;

  return encoded;
}
}  // namespace magma5g
