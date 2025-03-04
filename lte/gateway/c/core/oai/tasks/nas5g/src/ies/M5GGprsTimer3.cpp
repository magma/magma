/*
 * Copyright 2020 The Magma Authors.
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * */

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GGprsTimer3.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
GPRSTimer3Msg::GPRSTimer3Msg() {};
GPRSTimer3Msg::~GPRSTimer3Msg() {};

int GPRSTimer3Msg::DecodeGPRSTimer3Msg(GPRSTimer3Msg* gprstimer, uint8_t iei,
                                       uint8_t* buffer, uint32_t len) {
  int decoded = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  gprstimer->unit = (*(buffer + decoded) >> 5) & 0x7;
  gprstimer->timervalue = *(buffer + decoded) & 0x1f;
  decoded++;
  return decoded;
};

int GPRSTimer3Msg::EncodeGPRSTimer3Msg(GPRSTimer3Msg* gprstimer, uint8_t iei,
                                       uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, GPRS_TIMER3_MINIMUM_LENGTH, len);

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  *(buffer + encoded) = gprstimer->len;
  encoded++;
  *(buffer + encoded) =
      0x00 | ((gprstimer->unit & 0x7) << 5) | (gprstimer->timervalue & 0x1f);
  encoded++;
  return encoded;
};

}  // namespace magma5g
