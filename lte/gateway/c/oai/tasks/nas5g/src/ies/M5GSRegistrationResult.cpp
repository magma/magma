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

#include <iostream>
#include <sstream>
#include <bitset>
#include <cstdint>
#include "M5GSRegistrationResult.h"
#include "M5GCommonDefs.h"

using namespace std;
namespace magma5g {
M5GSRegistrationResultMsg::M5GSRegistrationResultMsg(){};
M5GSRegistrationResultMsg::~M5GSRegistrationResultMsg(){};

// Decode 5GS Registration Result message and its IEs
int M5GSRegistrationResultMsg::DecodeM5GSRegistrationResultMsg(
    M5GSRegistrationResultMsg* m5gs_reg_result, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  // Not yet implemented, Will be supported POST MVC
  return 0;
};

// Encode 5GS Registration Result message and its IEs
int M5GSRegistrationResultMsg::EncodeM5GSRegistrationResultMsg(
    M5GSRegistrationResultMsg* m5gs_reg_result, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, (uint32_t) REGISTRATION_RESULT_MIN_LENGTH, len);

  if (iei > 0) {
    CHECK_IEI_ENCODER(iei, (unsigned char) m5gs_reg_result->iei);
    *buffer = iei;
    MLOG(MDEBUG) << "In EncodeM5GSRegistrationResultMsg___: iei  = " << hex
                 << int(*buffer);
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;
  *(buffer + encoded) = ((m5gs_reg_result->spare & 0xf) << 0x4) |
                        ((m5gs_reg_result->sms_allowed & 0x1) << 3) |
                        (m5gs_reg_result->reg_result_val & 0x7);
  MLOG(MDEBUG) << " EncodeM5GSRegistrationResultMsg : 0x0"
               << int(*(buffer + encoded)) << " ("
               << "spare = " << ((m5gs_reg_result->spare & 0xf) << 0x4)
               << ", sms allowed = "
               << ((m5gs_reg_result->sms_allowed & 0x1) << 3)
               << ", result val = " << hex
               << int(m5gs_reg_result->reg_result_val & 0x7) << ")";
  encoded++;
  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);
  MLOG(MDEBUG) << " EncodeM5GSRegistrationResultMsg : length  = 0x0" << hex
               << int(*lenPtr);
  return encoded;
};
}  // namespace magma5g
