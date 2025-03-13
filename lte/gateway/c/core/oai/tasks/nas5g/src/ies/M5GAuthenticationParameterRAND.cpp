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
#include <cstring>
#include <cstdint>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GAuthenticationParameterRAND.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
AuthenticationParameterRANDMsg::AuthenticationParameterRANDMsg() {};
AuthenticationParameterRANDMsg::~AuthenticationParameterRANDMsg() {};

// Decode AuthenticationParameterRAND IE
int AuthenticationParameterRANDMsg::DecodeAuthenticationParameterRANDMsg(
    AuthenticationParameterRANDMsg* rand, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint8_t decoded = 0;
  /*** Not Implemented, Will be supported POST MVC ***/
  return (decoded);
};

// Encode AuthenticationParameterRAND IE
int AuthenticationParameterRANDMsg::EncodeAuthenticationParameterRANDMsg(
    AuthenticationParameterRANDMsg* rand, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint32_t encoded = 0;

  // Checking IEI and pointer
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, RAND_MIN_LEN, len);

  if (iei > 0) {
    CHECK_IEI_ENCODER((unsigned char)iei, rand->iei);
    *buffer = iei;
    encoded++;
  }

  memcpy(buffer + encoded, rand->rand_val, RAND_MAX_LEN);
  encoded = encoded + RAND_MAX_LEN;

  return (encoded);
};
}  // namespace magma5g
