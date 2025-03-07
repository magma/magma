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

#include <sstream>
#include <cstdint>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GMMCause.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
M5GMMCauseMsg::M5GMMCauseMsg() {};
M5GMMCauseMsg::~M5GMMCauseMsg() {};

// Decode 5GMMCause IE
int M5GMMCauseMsg::DecodeM5GMMCauseMsg(M5GMMCauseMsg* m5gmm_cause, uint8_t iei,
                                       uint8_t* buffer, uint32_t len) {
  uint8_t decoded = 0;

  if (iei > 0) {
    m5gmm_cause->iei = *(buffer + decoded);
    CHECK_IEI_DECODER((unsigned char)iei, m5gmm_cause->iei);
    decoded++;
  }

  m5gmm_cause->m5gmm_cause = *(buffer + decoded);
  decoded++;
  return (decoded);
};

// Encode 5GMMCause IE
int M5GMMCauseMsg::EncodeM5GMMCauseMsg(M5GMMCauseMsg* m5gmm_cause, uint8_t iei,
                                       uint8_t* buffer, uint32_t len) {
  int encoded = 0;

  if (iei > 0) {
    *(buffer + encoded) = m5gmm_cause->iei;
    CHECK_IEI_ENCODER((unsigned char)iei, m5gmm_cause->iei);
    encoded++;
  }

  *(buffer + encoded) = m5gmm_cause->m5gmm_cause;
  encoded++;
  return (encoded);
};
}  // namespace magma5g
