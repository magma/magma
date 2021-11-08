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
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPayloadContainer.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
PayloadContainerMsg::PayloadContainerMsg(){};
PayloadContainerMsg::~PayloadContainerMsg(){};

int PayloadContainerMsg::DecodePayloadContainerMsg(
    PayloadContainerMsg* payload_container, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded    = 0;
  uint32_t ielen = 0;
  IES_DECODE_U16(buffer, decoded, ielen);
  payload_container->len = ielen;
  MLOG(MDEBUG) << "DecodePayloadContainerMsg__: len = " << std::dec
               << int(payload_container->len) << std::endl;
  memcpy(&payload_container->contents, buffer + decoded, int(ielen));
  BUFFER_PRINT_LOG(payload_container->contents, int(ielen));

  // SMF NAS Message Decode
  decoded += payload_container->smf_msg.SmfMsgDecodeMsg(
      &payload_container->smf_msg, payload_container->contents, int(ielen));

  return (decoded);
};

int PayloadContainerMsg::EncodePayloadContainerMsg(
    PayloadContainerMsg* payload_container, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int encoded    = 0;
  uint32_t ielen = 0;
  int tmp        = 0;
  ielen          = payload_container->len;

  MLOG(MDEBUG) << "DecodePayloadContainerMsg__: len = " << std::hex
               << int(ielen) << std::endl;

  // SMF NAS Message Decode
  encoded += payload_container->smf_msg.SmfMsgEncodeMsg(
      &payload_container->smf_msg, payload_container->contents,
      payload_container->len);

  if (ielen != encoded) {
    MLOG(MDEBUG) << "WARNING: mismatch IE length :" << ielen
                 << " encoded SmfMsg length :" << encoded;
  }

  IES_ENCODE_U16(buffer, tmp, encoded);

  BUFFER_PRINT_LOG(payload_container->contents, encoded);
  memcpy(buffer + tmp, payload_container->contents, encoded);

  return (encoded + tmp);
}
}  // namespace magma5g
