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
#include <string.h>
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GNSSAI.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
namespace magma5g {
NSSAIMsg::NSSAIMsg(){};

NSSAIMsg::~NSSAIMsg(){};

int NSSAIMsg::EncodeNSSAIMsg(
    NSSAIMsg* NSSAI, uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint8_t encoded = 0;
  int i           = 0;

  /*
   * Checking IEI and pointer
   */
  /*CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, NSSAI_MIN_LENGTH, len);*/

  if (iei > 0) {
    CHECK_IEI_ENCODER(iei, (unsigned char) NSSAI->iei);
    *buffer = iei;
    encoded++;
  }
  // lenPtr = (buffer + encoded);
  *(buffer + encoded) = NSSAI->len;
  encoded++;

  for (i = 0; i < NSSAI->len; i++) {
    *(buffer + encoded) = NSSAI->nssaival[i];
    encoded++;
  }

  //*lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);

  return (encoded);
};

int NSSAIMsg::DecodeNSSAIMsg(
    NSSAIMsg* NSSAI, uint8_t iei, uint8_t* buffer, uint32_t len) {
  // will be implemented post MVC
  return (0);
};
}  // namespace magma5g
