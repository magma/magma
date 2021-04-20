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
#include "M5GNASSecurityAlgorithms.h"
#include "M5GCommonDefs.h"

using namespace std;
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

  MLOG(MDEBUG) << " DecodeNASSecurityAlgorithmsMsg : ";
  nas_sec_algorithms->tca = (*(buffer + decoded) >> 4) & 0x7;
  nas_sec_algorithms->tia = *(buffer + decoded) & 0x7;
  decoded++;
  MLOG(MDEBUG) << " Type of ciphering algorithm  = " << hex
               << int(nas_sec_algorithms->tca);
  MLOG(MDEBUG) << " Type of integrity protection algorithm  = " << hex
               << int(nas_sec_algorithms->tia);
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

  MLOG(MDEBUG) << " EncodeNASSecurityAlgorithmsMsg : ";
  *(buffer + encoded) = 0x00 | ((nas_sec_algorithms->tca & 0x7) << 4) |
                        (nas_sec_algorithms->tia & 0x7);

  MLOG(MDEBUG) << " Type of ciphering algorithm  = " << hex
               << int(nas_sec_algorithms->tca);
  MLOG(MDEBUG) << " Type of integrity protection algorithm  = " << hex
               << int(nas_sec_algorithms->tia);
  encoded++;
  return (encoded);
};
}  // namespace magma5g
