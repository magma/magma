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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GGprsTimer2.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
GPRSTimer2Msg::GPRSTimer2Msg() {};
GPRSTimer2Msg::~GPRSTimer2Msg() {};

int GPRSTimer2Msg::DecodeGPRSTimer2Msg(GPRSTimer2Msg* gprstimer, uint8_t iei,
                                       uint8_t* buffer, uint32_t len) {
  int decoded = 0;
  if (iei > 0) {
    gprstimer->iei = *buffer;
    decoded++;

    gprstimer->len = *(buffer + decoded);
    decoded++;

    gprstimer->timervalue = *(buffer + decoded);
    decoded++;
  }

  return decoded;
};

int GPRSTimer2Msg::EncodeGPRSTimer2Msg(GPRSTimer2Msg* gprstimer, uint8_t iei,
                                       uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

  if (iei > 0) {
    *buffer = iei;
    encoded++;

    *(buffer + encoded) = gprstimer->len;
    encoded++;
    *(buffer + encoded) = gprstimer->timervalue;
    encoded++;
  }

  return encoded;
};

}  // namespace magma5g
