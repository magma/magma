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
#include "M5GAuthenticationFailureIE.h"
#include "M5GCommonDefs.h"

using namespace std;
namespace magma5g {
M5GAuthenticationFailureIE::M5GAuthenticationFailureIE(){};
M5GAuthenticationFailureIE::~M5GAuthenticationFailureIE(){};

// Decode 5GMMCause IE
int M5GAuthenticationFailureIE::DecodeM5GAuthenticationFailureIE(
    M5GAuthenticationFailureIE* m5g_auth_failure_ie, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  uint8_t decoded = 0;
  uint8_t ielen   = 0;
  int decode_result;

  if (iei > 0) {
    m5g_auth_failure_ie->iei = *(buffer + decoded);
    CHECK_IEI_DECODER((unsigned char) iei, m5g_auth_failure_ie->iei);
    MLOG(MDEBUG) << "In DecodeM5GAuthenticationFailureIE: iei = " << dec
                 << int(m5g_auth_failure_ie->iei) << endl;
    decoded++;
  }

  ielen = *(buffer + decoded);
  decoded++;

  m5g_auth_failure_ie->authentication_failure_info =
      blk2bstr(buffer + decoded, ielen);

  decoded += ielen;

  return (decoded);
};

}  // namespace magma5g
