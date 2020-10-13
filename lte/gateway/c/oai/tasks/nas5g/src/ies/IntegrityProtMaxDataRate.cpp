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
#include "IntegrityProtMaxDataRate.h"
#include "CommonDefs.h"

using namespace std;
namespace magma5g {
IntegrityProtMaxDataRateMsg::IntegrityProtMaxDataRateMsg(){};
IntegrityProtMaxDataRateMsg::~IntegrityProtMaxDataRateMsg(){};

// Decode IntegrityProtMaxDataRate IE
int IntegrityProtMaxDataRateMsg::DecodeIntegrityProtMaxDataRateMsg(
    IntegrityProtMaxDataRateMsg* integrityprotmaxdatarate, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  uint8_t decoded = 0;

  MLOG(MDEBUG) << "   DecodeIntegrityProtMaxDataRateMsg : ";
  integrityprotmaxdatarate->maxuplink = *(buffer + decoded);
  decoded++;
  integrityprotmaxdatarate->maxdownlink = *(buffer + decoded);
  decoded++;
  MLOG(MDEBUG) << " maxuplink = " << dec
               << int(integrityprotmaxdatarate->maxuplink);
  MLOG(MDEBUG) << " maxdownlink = " << dec
               << int(integrityprotmaxdatarate->maxdownlink);
  return (decoded);
};

// Encode IntegrityProtMaxDataRate IE
int IntegrityProtMaxDataRateMsg::EncodeIntegrityProtMaxDataRateMsg(
    IntegrityProtMaxDataRateMsg* integrityprotmaxdatarate, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  int encoded = 0;

  MLOG(MDEBUG) << " EncodeIntegrityProtMaxDataRateMsg : ";
  *(buffer + encoded) = integrityprotmaxdatarate->maxuplink;
  MLOG(MDEBUG) << " maxuplink =0x" << hex
               << int(integrityprotmaxdatarate->maxuplink);
  encoded++;
  *(buffer + encoded) = integrityprotmaxdatarate->maxdownlink;
  MLOG(MDEBUG) << " maxdownlink =0x" << hex
               << int(integrityprotmaxdatarate->maxdownlink);
  encoded++;
  return (encoded);
};
}  // namespace magma5g
