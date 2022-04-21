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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPDUSessionStatus.hpp"

namespace magma5g {
M5GPDUSessionStatus::M5GPDUSessionStatus(){};
M5GPDUSessionStatus::~M5GPDUSessionStatus(){};

int M5GPDUSessionStatus::EncodePDUSessionStatus(
    M5GPDUSessionStatus* pduSessionStatus, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int encoded = 0;
  if (pduSessionStatus->iei) {
    *(buffer + encoded) = pduSessionStatus->iei;
    encoded++;
    *(buffer + encoded) = pduSessionStatus->len;
    encoded++;
    *(buffer + encoded) = (uint8_t)(pduSessionStatus->pduSessionStatus & 0xFF);
    encoded++;
    *(buffer + encoded) =
        (uint8_t)(((pduSessionStatus->pduSessionStatus >> 8) & 0xFF));
    encoded++;
  }

  return encoded;
}

int M5GPDUSessionStatus::DecodePDUSessionStatus(
    M5GPDUSessionStatus* pduSessionStatus, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded = 0;

  if (iei > 0) {
    pduSessionStatus->iei = *buffer;
    decoded++;

    pduSessionStatus->len = *(buffer + decoded);
    decoded++;

    pduSessionStatus->pduSessionStatus = *(buffer + decoded);
    decoded++;
    pduSessionStatus->pduSessionStatus |= (*(buffer + decoded) << 8);
    decoded++;
  }

  return decoded;
}
}  // namespace magma5g
