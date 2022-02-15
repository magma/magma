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

#pragma once
#include <sstream>
#include <cstdint>

namespace magma5g {
// MessageType IE Class
class MessageTypeMsg {
 public:
  uint8_t msg_type;

  MessageTypeMsg();
  ~MessageTypeMsg();
  int EncodeMessageTypeMsg(MessageTypeMsg* message_type, uint8_t iei,
                           uint8_t* buffer, uint32_t len);
  int DecodeMessageTypeMsg(MessageTypeMsg* message_type, uint8_t iei,
                           uint8_t* buffer, uint32_t len);

  void copy(const MessageTypeMsg& m) { msg_type = m.msg_type; }
  bool isEqual(const MessageTypeMsg& m) { return (msg_type == m.msg_type); }
};
}  // namespace magma5g
