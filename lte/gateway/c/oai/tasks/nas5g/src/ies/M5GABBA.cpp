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
#include "M5GABBA.h"
#include "M5GCommonDefs.h"

using namespace std;
namespace magma5g {
ABBAMsg::ABBAMsg(){};
ABBAMsg::~ABBAMsg(){};

// Decode ABBA Message IE
int ABBAMsg::DecodeABBAMsg(
    ABBAMsg* abba, uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint8_t decoded = 0;
  /*** Not Implemented, Will be supported POST MVC ***/
  return (decoded);
};

// Encode ABBA Message IE
int ABBAMsg::EncodeABBAMsg(
    ABBAMsg* abba, uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;

  // Checking IEI and pointer
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, ABBA_MIN_LEN, len);

  if (iei > 0) {
    CHECK_IEI_ENCODER((unsigned char) iei, abba->iei);
    *buffer = iei;
    MLOG(MDEBUG) << "In EncodeABBAMsg: iei" << hex << int(*buffer);
    encoded++;
  }

  MLOG(MDEBUG) << " EncodeABBAMsg : ";
  lenPtr = buffer + encoded;
  encoded++;
  std::copy(abba->contents.begin(), abba->contents.end(), buffer + encoded);
  MLOG(MDEBUG) << "   Length : " << dec << int(abba->contents.length());
  MLOG(MDEBUG) << "   ABBA Contents : ";
  BUFFER_PRINT_LOG(buffer + encoded, abba->contents.length());
  encoded = encoded + abba->contents.length();
  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);

  return (encoded);
};
}  // namespace magma5g
