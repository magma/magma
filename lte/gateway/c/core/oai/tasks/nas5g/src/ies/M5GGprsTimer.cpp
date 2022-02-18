/*
 * Copyright 020 The Magma Authors.
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 * */

#include <sstream>
#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GGprsTimer.h"

namespace magma5g {
GPRSTimerMsg::GPRSTimerMsg() {}
GPRSTimerMsg::~GPRSTimerMsg() {}

int GPRSTimerMsg::DecodeGPRSTimerMsg(GPRSTimerMsg* gprstimer, uint8_t iei,
                                     uint8_t* buffer, uint32_t len) {
  int decoded = 0;

  if (iei > 0) {
    gprstimer->iei = *buffer;
    MLOG(MDEBUG) << "DecodeGPRSTimerMsg: iei = " << std::hex
                 << int(gprstimer->iei);
    decoded++;

    gprstimer->timervalue = *(buffer + decoded);
    MLOG(MDEBUG) << "DecodeGPRSTimerMsg: timervalue = " << std::hex
                 << int(gprstimer->timervalue);
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
