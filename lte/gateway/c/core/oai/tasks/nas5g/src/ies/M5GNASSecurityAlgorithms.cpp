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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GNASSecurityAlgorithms.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
NASSecurityAlgorithmsMsg::NASSecurityAlgorithmsMsg(){};
NASSecurityAlgorithmsMsg::~NASSecurityAlgorithmsMsg(){};

// Decode NASSecurityAlgorithms IE
int NASSecurityAlgorithmsMsg::DecodeNASSecurityAlgorithmsMsg(
    NASSecurityAlgorithmsMsg* nas_sec_algorithms, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint8_t decoded = 0;

  // Checking IEI
  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  OAILOG_DEBUG(LOG_NAS5G, "Decoding NASSecurityAlgorithms");
  nas_sec_algorithms->tca = (*(buffer + decoded) >> 4) & 0x7;
  nas_sec_algorithms->tia = *(buffer + decoded) & 0x7;
  decoded++;
  OAILOG_DEBUG(
      LOG_NAS5G, "Type of ciphering algorithm : %X",
      static_cast<int>(nas_sec_algorithms->tca));
  OAILOG_DEBUG(
      LOG_NAS5G, "Type of integrity protection algorithm : %X",
      static_cast<int>(nas_sec_algorithms->tia));
  return (decoded);
};

// Encode NASSecurityAlgorithms IE
int NASSecurityAlgorithmsMsg::EncodeNASSecurityAlgorithmsMsg(
    NASSecurityAlgorithmsMsg* nas_sec_algorithms, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int encoded = 0;

  // Checking IEI and pointer
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, NAS_SECURITY_ALGORITHMS_MINIMUM_LENGTH, len);

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  OAILOG_DEBUG(LOG_NAS5G, "Encoding NASSecurityAlgorithms");
  *(buffer + encoded) = 0x00 | ((nas_sec_algorithms->tca & 0x7) << 4) |
                        (nas_sec_algorithms->tia & 0x7);

  OAILOG_DEBUG(
      LOG_NAS5G, "Type of ciphering algorithm : %X",
      static_cast<int>(nas_sec_algorithms->tca));
  OAILOG_DEBUG(
      LOG_NAS5G, "Type of integrity protection algorithm : %X",
      static_cast<int>(nas_sec_algorithms->tia));
  encoded++;
  return (encoded);
};
}  // namespace magma5g
