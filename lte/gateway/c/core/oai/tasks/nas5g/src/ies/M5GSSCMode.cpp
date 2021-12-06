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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GSSCMode.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
SSCModeMsg::SSCModeMsg(){};
SSCModeMsg::~SSCModeMsg(){};

// Decode SSCMode IE
int SSCModeMsg::DecodeSSCModeMsg(
    SSCModeMsg* ssc_mode, uint8_t iei, uint8_t* buffer, uint32_t len) {
  int decoded = 0;

  // Storing the IEI Information
  if (iei > 0) {
    ssc_mode->iei = (*buffer & 0xf0) >> 4;
    MLOG(MDEBUG) << "In DecodeSSCModeMsg: iei = " << std::hex
                 << int(ssc_mode->iei);
    decoded++;
  }

  if (iei > 0) {
    ssc_mode->mode_val = (*buffer & 0x07);
  } else {
    ssc_mode->mode_val = (*buffer >> 4) & 0x07;
  }
  MLOG(MDEBUG) << "DecodeSSCModeMsg__: mode_val = " << std::hex
               << int(ssc_mode->mode_val);

  return decoded;
};

// Encode SSCMode IE
int SSCModeMsg::EncodeSSCModeMsg(
    SSCModeMsg* ssc_mode, uint8_t iei, uint8_t* buffer, uint32_t len) {
  int encoded = 0;

  // CHECKING IEI
  if (iei > 0) {
    CHECK_IEI_ENCODER(
        (uint8_t) iei, (uint8_t)(0x00 | (ssc_mode->iei & 0x0f) << 4));
    *buffer = (ssc_mode->iei & 0x0f) << 4;
    MLOG(MDEBUG) << "In EncodeSSCModeMsg: iei" << std::hex << int(*buffer);
    encoded++;
  }
  if (iei > 0) {
    *buffer = 0x00 | (*buffer & 0xf0) | (ssc_mode->mode_val & 0x07);
  } else {
    *buffer = 0x00 | (*buffer & 0x0f) | ((ssc_mode->mode_val & 0x07) << 4);
  }
  MLOG(MDEBUG) << "EncodeSSCModeMsg__: mode_val = " << std::hex << int(*buffer);

  return (encoded);
};
}  // namespace magma5g
