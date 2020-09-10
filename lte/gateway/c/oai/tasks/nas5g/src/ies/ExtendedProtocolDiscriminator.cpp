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
#include "ExtendedProtocolDiscriminator.h"
#include "CommonDefs.h"
using namespace std;
namespace magma5g
{
  ExtendedProtocolDiscriminatorMsg::ExtendedProtocolDiscriminatorMsg()
  {
  };

  ExtendedProtocolDiscriminatorMsg::~ExtendedProtocolDiscriminatorMsg()
  {
  };

  // Decode ExtendedProtocolDiscriminator IE
  int ExtendedProtocolDiscriminatorMsg::DecodeExtendedProtocolDiscriminatorMsg(ExtendedProtocolDiscriminatorMsg *extendedprotocoldiscriminator, uint8_t iei, uint8_t *buffer, uint32_t len) 
  {
    uint8_t decoded = 0;

    MLOG(MDEBUG) << "   DecodeExtendedProtocolDiscriminatorMsg : ";
    extendedprotocoldiscriminator->extendedprotodiscriminator = *(buffer + decoded);
    decoded++;
    MLOG(MDEBUG) << " epd = " << hex << int(extendedprotocoldiscriminator->extendedprotodiscriminator)<<"\n"; 
    return (decoded);
  };

  // Encode ExtendedProtocolDiscriminator IE
  int ExtendedProtocolDiscriminatorMsg::EncodeExtendedProtocolDiscriminatorMsg(ExtendedProtocolDiscriminatorMsg *extendedprotocoldiscriminator, uint8_t iei, uint8_t * buffer, uint32_t len)
  {
    int encoded = 0;

    MLOG(MDEBUG) << "\n\n EncodeExtendedProtocolDiscriminatorMsg : ";
    *(buffer + encoded) = extendedprotocoldiscriminator->extendedprotodiscriminator;
    MLOG(MDEBUG) << "epd = 0x" << hex << int(*(buffer + encoded))<<"\n"; 
    encoded++;
    return (encoded);
  };
}

