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
#include "PDUSessionType.h"
#include "CommonDefs.h"

using namespace std;
namespace magma5g {
PDUSessionTypeMsg::PDUSessionTypeMsg(){};
PDUSessionTypeMsg::~PDUSessionTypeMsg(){};

// Decode PDUSessionType IE
int PDUSessionTypeMsg::DecodePDUSessionTypeMsg(
    PDUSessionTypeMsg* pdusessiontype, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded = 0;

  if (iei > 0) {
    pdusessiontype->iei = (*buffer & 0xf0) >> 4;
    CHECK_IEI_DECODER((unsigned char) iei, pdusessiontype->iei);
    MLOG(MDEBUG) << "In DecodePDUSessionTypeMsg: iei" << hex
                 << int(pdusessiontype->iei) << endl;
    decoded++;
  }

  pdusessiontype->typeval = (*buffer & 0x07);
  MLOG(MDEBUG) << "DecodePDUSessionTypeMsg: typeval = " << hex
               << int(pdusessiontype->typeval) << endl;

  return (decoded);
};

// Encode PDUSessionType IE
int PDUSessionTypeMsg::EncodePDUSessionTypeMsg(
    PDUSessionTypeMsg* pdusessiontype, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int encoded = 0;

  if (iei > 0) {
    *buffer = (pdusessiontype->iei & 0x0f) << 4;
    CHECK_IEI_ENCODER((unsigned char) iei, pdusessiontype->iei);
    MLOG(MDEBUG) << "In EncodePDUSessionTypeMsg: iei" << hex << int(*buffer)
                 << endl;
  }

  *buffer = 0x00 | (*buffer & 0xf0) | (pdusessiontype->typeval & 0x07);
  MLOG(MDEBUG) << "EncodePDUSessionTypeMsg: typeval = " << hex << int(*(buffer))
               << endl;
  encoded++;

  return (encoded);
};
}  // namespace magma5g
