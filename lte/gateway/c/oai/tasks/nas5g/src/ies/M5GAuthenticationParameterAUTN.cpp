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
#include "M5GAuthenticationParameterAUTN.h"
#include "M5GCommonDefs.h"

using namespace std;
namespace magma5g {
AuthenticationParameterAUTNMsg::AuthenticationParameterAUTNMsg(){};
AuthenticationParameterAUTNMsg::~AuthenticationParameterAUTNMsg(){};

// Decode AuthenticationParameterAUTN IE
int AuthenticationParameterAUTNMsg::DecodeAuthenticationParameterAUTNMsg(
    AuthenticationParameterAUTNMsg* autn, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint8_t decoded = 0;
  /*** Not Implemented, Will be supported POST MVC ***/
  return (decoded);
};

// Encode AuthenticationParameterAUTN IE
int AuthenticationParameterAUTNMsg::EncodeAuthenticationParameterAUTNMsg(
    AuthenticationParameterAUTNMsg* autn, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;

  // Checking IEI and pointer
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, AUTN_MIN_LEN, len);

  if (iei > 0) {
    CHECK_IEI_ENCODER((unsigned char) iei, autn->iei);
    *buffer = iei;
    MLOG(MDEBUG) << "In EncodeAuthenticationParameterAUTNMsg: iei" << hex
                 << int(*buffer);
    encoded++;
  }

  lenPtr = (uint8_t*) (buffer + encoded);
  encoded++;
  std::copy(autn->AUTN.begin(), autn->AUTN.end(), buffer + encoded);
  BUFFER_PRINT_LOG(buffer + encoded, autn->AUTN.length());
  encoded = encoded + autn->AUTN.length();
  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);

  return (encoded);
};
}  // namespace magma5g
