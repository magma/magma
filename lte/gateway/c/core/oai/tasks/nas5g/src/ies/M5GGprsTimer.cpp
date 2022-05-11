/*
 * Copyright 2022 The Magma Authors.
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
#include <sstream>

#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GGprsTimer.hpp"

#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif

namespace magma5g {
GPRSTimerMsg::GPRSTimerMsg() {}
GPRSTimerMsg::~GPRSTimerMsg() {}

int GPRSTimerMsg::DecodeGPRSTimerMsg(GPRSTimerMsg* gprstimer, uint8_t iei,
                                     uint8_t* buffer, uint32_t len) {
  int decoded = 0;

  if (iei > 0) {
    gprstimer->iei = *buffer;
    OAILOG_DEBUG(LOG_NAS5G, "DecodeGPRSTimerMsg: iei 0x%x",
                 static_cast<int>(gprstimer->iei));
    decoded++;

    gprstimer->timervalue = *(buffer + decoded);
    OAILOG_DEBUG(LOG_NAS5G, "DecodeGPRSTimerMsg: timervalue = 0X%x",
                 static_cast<int>(gprstimer->timervalue));
    decoded++;
  }

  return decoded;
}

int GPRSTimerMsg::EncodeGPRSTimerMsg(GPRSTimerMsg* gprstimer, uint8_t iei,
                                     uint8_t* buffer, uint32_t len) {
  uint8_t encoded = 0;

  if (iei > 0) {
    *buffer = iei;
    encoded++;

    *(buffer + encoded) = gprstimer->timervalue;
    encoded++;
  }

  return encoded;
}

}  // namespace magma5g
