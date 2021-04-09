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
#include "M5GNSSAI.h"
#include "M5GCommonDefs.h"
using namespace std;
namespace magma5g {
NSSAIMsg::NSSAIMsg(){};

NSSAIMsg::~NSSAIMsg(){};

int NSSAIMsg::EncodeNSSAIMsg(
    NSSAIMsg* NSSAI, uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint8_t encoded = 0;
  uint8_t ielen   = 0;
  uint8_t* lenPtr;
  int i;

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
  uint8_t decoded = 0;
  uint8_t ielen   = 0;
#if 0
  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(buffer, NSSAI_MIN_LENGTH, len);

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, (unsigned char) *buffer);
    *buffer = iei;
    decoded++;
  }

  IES_DECODE_U8(buffer, decoded, ielen);

  if (!memcpy(&NSSAI->nssaival, (buffer + decoded), ielen)) return decoded;

  decoded = decoded + ielen;
#endif
  return (decoded);
};
}  // namespace magma5g
