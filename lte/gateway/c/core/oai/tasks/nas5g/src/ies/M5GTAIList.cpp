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
#include "M5GTAIList.h"
#include <string.h>
#include <cstdint>
#include <sstream>
#include "M5GCommonDefs.h"
using namespace std;
namespace magma5g {
TAIListMsg::TAIListMsg(){};

TAIListMsg::~TAIListMsg(){};

int TAIListMsg::EncodeTAIListMsg(TAIListMsg* TAIList, uint8_t iei,
                                 uint8_t* buffer, uint32_t len) {
  uint8_t encoded = 0;
  if (iei > 0) {
    CHECK_IEI_ENCODER(iei, (unsigned char)TAIList->iei);
    *buffer = iei;
    MLOG(MDEBUG) << "iei = " << hex << int(*(buffer + encoded));
    encoded++;
  }
  *(buffer + encoded) = TAIList->len;
  encoded++;

  *(buffer + encoded) = 0x00 | ((TAIList->list_type & 0x03) << 5) |
                        (TAIList->num_elements & 0x1f);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((TAIList->mcc_digit2 & 0x0f) << 4) | (TAIList->mcc_digit1 & 0x0f);
  MLOG(MDEBUG) << "mcc_digit2 >mcc_digit1 type_of_identity = " << hex
               << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) =
      0x00 | ((TAIList->mnc_digit3 & 0x0f) << 4) | (TAIList->mcc_digit3 & 0x0f);
  MLOG(MDEBUG) << "mnc_digit3 >mcc_digit3 type_of_identity = " << hex
               << int(*(buffer + encoded));
  encoded++;
  *(buffer + encoded) =
      0x00 | ((TAIList->mnc_digit2 & 0x0f) << 4) | (TAIList->mnc_digit1 & 0x0f);
  MLOG(MDEBUG) << "mnc_digit2 >mcc_digit1 type_of_identity = " << hex
               << int(*(buffer + encoded));
  encoded++;

  *(buffer + encoded) = TAIList->tac[0];
  encoded++;
  *(buffer + encoded) = TAIList->tac[1];
  encoded++;
  *(buffer + encoded) = TAIList->tac[2];
  encoded++;

  return (encoded);
}
int TAIListMsg::DecodeTAIListMsg(TAIListMsg* TAIList, uint8_t iei,
                                 uint8_t* buffer, uint32_t len) {
  return 0;
}

}  // namespace magma5g
