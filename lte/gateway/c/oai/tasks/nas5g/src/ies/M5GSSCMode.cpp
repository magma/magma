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
#include <cstdint>
#include <cstring>
#include "M5GSSCMode.h"
#include "M5GCommonDefs.h"

using namespace std;
namespace magma5g {
SSCModeMsg::SSCModeMsg(){};
SSCModeMsg::~SSCModeMsg(){};

// Decode SSCMode IE
int SSCModeMsg::DecodeSSCModeMsg(
    SSCModeMsg* ssc_mode, uint8_t iei, uint8_t* buffer, uint32_t len) {
  int decoded = 0;

  // CHECKING IEI
  if (iei > 0) {
    ssc_mode->iei = (*buffer & 0xf0) >> 4;
    CHECK_IEI_ENCODER((unsigned char) iei, ssc_mode->iei);
    MLOG(MDEBUG) << "In DecodeSSCModeMsg: iei = " << hex << int(ssc_mode->iei);
    decoded++;
  }

  ssc_mode->mode_val = (*buffer & 0x07);
  MLOG(MDEBUG) << "DecodeSSCModeMsg__: mode_val = " << hex
               << int(ssc_mode->mode_val);

  return decoded;
};

// Encode SSCMode IE
int SSCModeMsg::EncodeSSCModeMsg(
    SSCModeMsg* ssc_mode, uint8_t iei, uint8_t* buffer, uint32_t len) {
  int encoded = 0;

  // CHECKING IEI
  if (iei > 0) {
    CHECK_IEI_ENCODER((unsigned char) iei, ssc_mode->iei);
    *buffer = 0x00 | (ssc_mode->iei & 0x0f) << 4;
    MLOG(MDEBUG) << "In EncodeSSCModeMsg: iei" << hex << int(*buffer);
    encoded++;
  }

  *buffer = (ssc_mode->mode_val << 4) & 0xf0;
  MLOG(MDEBUG) << "EncodeSSCModeMsg__: mode_val = " << hex << int(*buffer);

  return (encoded);
};
}  // namespace magma5g
