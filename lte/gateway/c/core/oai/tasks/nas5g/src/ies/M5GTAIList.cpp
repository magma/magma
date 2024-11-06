/*
copyright 2020 The Magma Authors.
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
#include <string.h>
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GTAIList.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
namespace magma5g {
TAIListMsg::TAIListMsg(){};

TAIListMsg::~TAIListMsg(){};

int TAIListMsg::EncodeTAIListMsg(uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint8_t encoded = 0;

  if (iei > 0) {
    CHECK_IEI_ENCODER(iei, (unsigned char)this->iei);
    *(buffer + encoded++) = iei;
  }
  *(buffer + encoded++) = this->len;

  *(buffer + encoded++) =
      0x00 | ((this->list_type & 0x03) << 5) | (this->num_elements & 0x1f);
  *(buffer + encoded++) =
      0x00 | ((this->mcc_digit2 & 0x0f) << 4) | (this->mcc_digit1 & 0x0f);
  *(buffer + encoded++) =
      0x00 | ((this->mnc_digit3 & 0x0f) << 4) | (this->mcc_digit3 & 0x0f);
  *(buffer + encoded++) =
      0x00 | ((this->mnc_digit2 & 0x0f) << 4) | (this->mnc_digit1 & 0x0f);

  *(buffer + encoded++) = this->tac[0];
  *(buffer + encoded++) = this->tac[1];
  *(buffer + encoded++) = this->tac[2];

  return (encoded);
}

int TAIListMsg::DecodeTAIListMsg(uint8_t iei, uint8_t* buffer, uint32_t len) {
  return 0;
}

}  // namespace magma5g
