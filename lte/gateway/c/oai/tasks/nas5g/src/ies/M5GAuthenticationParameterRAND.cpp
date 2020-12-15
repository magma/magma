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
#include "M5GAuthenticationParameterRAND.h"
#include "M5GCommonDefs.h"

using namespace std;
namespace magma5g {
AuthenticationParameterRANDMsg::AuthenticationParameterRANDMsg(){};
AuthenticationParameterRANDMsg::~AuthenticationParameterRANDMsg(){};

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
    CHECK_IEI_ENCODER((unsigned char) iei, rand->iei);
    *buffer = iei;
    MLOG(MDEBUG) << "In EncodeAuthenticationParameterRANDMsg: iei" << hex
                 << int(*buffer) << endl;
    encoded++;
  }

  std::copy(rand->rand_val.begin(), rand->rand_val.end(), buffer + encoded);
  BUFFER_PRINT_LOG(buffer + encoded, rand->rand_val.length());
  encoded = encoded + rand->rand_val.length();

  return (encoded);
};
}  // namespace magma5g
