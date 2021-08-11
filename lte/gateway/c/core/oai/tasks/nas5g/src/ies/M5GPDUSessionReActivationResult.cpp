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
#include "M5GCommonDefs.h"
#include "M5GPDUSessionReActivationResult.h"

using namespace std;
namespace magma5g {
M5GPDUSessionReActivationResult::M5GPDUSessionReActivationResult(){};
M5GPDUSessionReActivationResult::~M5GPDUSessionReActivationResult(){};

int M5GPDUSessionReActivationResult::EncodePDUSessionReActivationResult(
    M5GPDUSessionReActivationResult* pduSessionReActivationStatus, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  int encoded = 0;
  if (pduSessionReActivationStatus->iei) {
    *(buffer + encoded) = pduSessionReActivationStatus->iei;
    encoded++;
    *(buffer + encoded) = pduSessionReActivationStatus->len;
    encoded++;
    *(buffer + encoded) =
        (pduSessionReActivationStatus->pduSessionReActivationResult & 0xFF);
    encoded++;
    *(buffer + encoded) =
        ((pduSessionReActivationStatus->pduSessionReActivationResult >> 8) &
         0xFF);
    encoded++;
  }

  return encoded;
}

int M5GPDUSessionReActivationResult::DecodePDUSessionReActivationResult(
    M5GPDUSessionReActivationResult* pduSessionReActivationStatus, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  int decoded = 0;

  if (iei > 0) {
    pduSessionReActivationStatus->iei = *buffer;
    MLOG(MDEBUG) << "In DecodeM5GPDUSessionReActivationResult: iei = " << hex
                 << int(pduSessionReActivationStatus->iei);
    decoded++;

    pduSessionReActivationStatus->len = *(buffer + decoded);
    MLOG(MDEBUG) << "In DecodeM5GPDUSessionReActivationResult: len = " << hex
                 << int(pduSessionReActivationStatus->len);
    decoded++;

    pduSessionReActivationStatus->pduSessionReActivationResult =
        *(buffer + decoded);
    decoded++;
    pduSessionReActivationStatus->pduSessionReActivationResult |=
        (*(buffer + decoded) << 8);
    MLOG(MDEBUG)
        << "In DecodeM5GPDUSessionReActivationResult: "
           "pduSessionReActivationResult = "
        << hex
        << int(pduSessionReActivationStatus->pduSessionReActivationResult);
    decoded++;
  }

  return decoded;
}
}  // namespace magma5g
