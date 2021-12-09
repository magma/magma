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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GMessageType.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
MessageTypeMsg::MessageTypeMsg(){};
MessageTypeMsg::~MessageTypeMsg(){};

// Decode MessageType IE
int MessageTypeMsg::DecodeMessageTypeMsg(MessageTypeMsg* message_type,
                                         uint8_t iei, uint8_t* buffer,
                                         uint32_t len) {
  uint8_t decoded = 0;

  OAILOG_DEBUG(LOG_NAS5G, " Decoding MessageType");
  message_type->msg_type = *(buffer + decoded);
  decoded++;
  OAILOG_DEBUG(
      LOG_NAS5G, "Message Type : 0x%X",
      static_cast<int>(message_type->msg_type));
  return (decoded);
};

// Encode MessageType IE
int MessageTypeMsg::EncodeMessageTypeMsg(MessageTypeMsg* message_type,
                                         uint8_t iei, uint8_t* buffer,
                                         uint32_t len) {
  uint8_t encoded = 0;

  OAILOG_DEBUG(LOG_NAS5G, " Encoding MessageType");
  *(buffer + encoded) = message_type->msg_type;
  OAILOG_DEBUG(
      LOG_NAS5G, "Message type = 0x%X", static_cast<int>(*(buffer + encoded)));
  encoded++;
  return (encoded);
};
}  // namespace magma5g
