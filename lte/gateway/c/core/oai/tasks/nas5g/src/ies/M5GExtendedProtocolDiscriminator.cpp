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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GExtendedProtocolDiscriminator.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
ExtendedProtocolDiscriminatorMsg::ExtendedProtocolDiscriminatorMsg() {};
ExtendedProtocolDiscriminatorMsg::~ExtendedProtocolDiscriminatorMsg() {};

// Decode ExtendedProtocolDiscriminator IE
int ExtendedProtocolDiscriminatorMsg::DecodeExtendedProtocolDiscriminatorMsg(
    ExtendedProtocolDiscriminatorMsg* extended_protocol_discriminator,
    uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint8_t decoded = 0;

  extended_protocol_discriminator->extended_proto_discriminator =
      *(buffer + decoded);
  decoded++;
  return (decoded);
};

// Encode ExtendedProtocolDiscriminator IE
int ExtendedProtocolDiscriminatorMsg::EncodeExtendedProtocolDiscriminatorMsg(
    ExtendedProtocolDiscriminatorMsg* extended_protocol_discriminator,
    uint8_t iei, uint8_t* buffer, uint32_t len) {
  int encoded = 0;

  *(buffer + encoded) =
      extended_protocol_discriminator->extended_proto_discriminator;
  encoded++;
  return (encoded);
};
}  // namespace magma5g
