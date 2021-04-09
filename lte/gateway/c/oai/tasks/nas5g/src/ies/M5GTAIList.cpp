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
#include "M5GTAIList.h"
#include "M5GCommonDefs.h"
using namespace std;
namespace magma5g {
TAIListMsg::TAIListMsg(){};

TAIListMsg::~TAIListMsg(){};

int TAIListMsg::EncodeTAIListMsg(
    TAIListMsg* TAIList, uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint8_t encoded = 0;
  //  uint8_t ielen   = 0;
  //  uint8_t* lenPtr;
  // int i;
  if (iei > 0) {
    CHECK_IEI_ENCODER(iei, (unsigned char) TAIList->iei);
    *buffer = iei;
    MLOG(MDEBUG) << "iei = " << hex << int(*(buffer + encoded));
    encoded++;
  }
  MLOG(MDEBUG) << "encoded = " << encoded;
  MLOG(MDEBUG) << "iei-- =  " << hex << int(*(buffer - 2));
  MLOG(MDEBUG) << "iei- = " << hex << int(*(buffer - 1));
  MLOG(MDEBUG) << "iei = " << hex << int(*(buffer));
  MLOG(MDEBUG) << "iei+ = " << hex << int(*(buffer + 1));

  // lenPtr = (buffer + encoded);
  *(buffer + encoded) = TAIList->len;
  encoded++;
  MLOG(MDEBUG) << "encoded = " << encoded;
  MLOG(MDEBUG) << "iei = " << hex << int(*(buffer));
  MLOG(MDEBUG) << "iei+ = " << hex << int(*(buffer + 1));

  *(buffer + encoded) = 0x00 | ((TAIList->list_type & 0x03) << 5) |
                        (TAIList->num_elements & 0x1f);
  encoded++;
  MLOG(MDEBUG) << "encoded = " << encoded;
  MLOG(MDEBUG) << "iei = " << hex << int(*(buffer));
  MLOG(MDEBUG) << "iei+ = " << hex << int(*(buffer + 1));
  *(buffer + encoded) =
      0x00 | ((TAIList->mcc_digit2 & 0x0f) << 4) | (TAIList->mcc_digit1 & 0x0f);
  MLOG(MDEBUG) << "mcc_digit2 >mcc_digit1 type_of_identity = " << hex
               << int(*(buffer + encoded));
  encoded++;
  MLOG(MDEBUG) << "encoded = " << encoded;
  MLOG(MDEBUG) << "iei = " << hex << int(*(buffer));
  MLOG(MDEBUG) << "iei+ = " << hex << int(*(buffer + 1));
  *(buffer + encoded) =
      0x00 | ((TAIList->mnc_digit3 & 0x0f) << 4) | (TAIList->mcc_digit3 & 0x0f);
  MLOG(MDEBUG) << "mnc_digit3 >mcc_digit3 type_of_identity = " << hex
               << int(*(buffer + encoded));
  encoded++;
  MLOG(MDEBUG) << "encoded = " << encoded;
  MLOG(MDEBUG) << "iei = " << hex << int(*(buffer));
  MLOG(MDEBUG) << "iei+ = " << hex << int(*(buffer + 1));
  *(buffer + encoded) =
      0x00 | ((TAIList->mnc_digit2 & 0x0f) << 4) | (TAIList->mnc_digit1 & 0x0f);
  MLOG(MDEBUG) << "mnc_digit2 >mcc_digit1 type_of_identity = " << hex
               << int(*(buffer + encoded));
  encoded++;
  MLOG(MDEBUG) << "encoded = " << encoded;
  MLOG(MDEBUG) << "iei = " << hex << int(*(buffer));
  MLOG(MDEBUG) << "iei+ = " << hex << int(*(buffer + 1));

  *(buffer + encoded) = TAIList->tac[0];
  encoded++;
  *(buffer + encoded) = TAIList->tac[1];
  encoded++;
  *(buffer + encoded) = TAIList->tac[2];
  encoded++;

  //*lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);

  return (encoded);
}
int TAIListMsg::DecodeTAIListMsg(
    TAIListMsg* TAIList, uint8_t iei, uint8_t* buffer, uint32_t len) {
  return 0;
}

}  // namespace magma5g
