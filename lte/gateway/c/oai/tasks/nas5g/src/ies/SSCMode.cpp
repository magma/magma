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
#include "SSCMode.h"
#include "CommonDefs.h"

using namespace std;
namespace magma5g {
SSCModeMsg::SSCModeMsg(){};
SSCModeMsg::~SSCModeMsg(){};

// Decode SSCMode IE
int SSCModeMsg::DecodeSSCModeMsg(
    SSCModeMsg* sscmode, uint8_t iei, uint8_t* buffer, uint32_t len) {
  int decoded = 0;

  if (iei > 0) {
    sscmode->iei = (*buffer & 0xf0) >> 4;
    CHECK_IEI_ENCODER((unsigned char) iei, sscmode->iei);
    MLOG(MDEBUG) << "In DecodeSSCModeMsg: iei = " << hex << int(sscmode->iei)
                 << endl;
    decoded++;
  }

  sscmode->modeval = (*buffer & 0x07);
  MLOG(MDEBUG) << "DecodeSSCModeMsg__: modeval = " << hex
               << int(sscmode->modeval) << endl;

  return (decoded);
};

// Encode SSCMode IE
int SSCModeMsg::EncodeSSCModeMsg(
    SSCModeMsg* sscmode, uint8_t iei, uint8_t* buffer, uint32_t len) {
  int encoded = 0;

  if (iei > 0) {
    CHECK_IEI_ENCODER((unsigned char) iei, sscmode->iei);
    *buffer = 0x00 | (sscmode->iei & 0x0f) << 4;
    MLOG(MDEBUG) << "In EncodeSSCModeMsg: iei" << hex << int(*buffer) << endl;
    encoded++;
  }

  *buffer = (sscmode->modeval << 4) & 0xf0;
  MLOG(MDEBUG) << "EncodeSSCModeMsg__: modeval = " << hex << int(*buffer)
               << endl;

  return (encoded);
};
}  // namespace magma5g
