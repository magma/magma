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
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GPayloadContainer.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
PayloadContainerMsg::PayloadContainerMsg(){};
PayloadContainerMsg::~PayloadContainerMsg(){};

int PayloadContainerMsg::DecodePayloadContainerMsg(uint8_t iei, uint8_t* buffer,
                                                   uint32_t len) {
  int decoded = 0;
  uint32_t ielen = 0;
  IES_DECODE_U16(buffer, decoded, ielen);
  this->len = ielen;
  memcpy(&this->contents, buffer + decoded, static_cast<int>(ielen));

  // SMF NAS Message Decode
  decoded +=
      this->smf_msg.SmfMsgDecodeMsg(this->contents, static_cast<int>(ielen));

  return (decoded);
};

int PayloadContainerMsg::EncodePayloadContainerMsg(uint8_t iei, uint8_t* buffer,
                                                   uint32_t len) {
  int encoded = 0;
  uint32_t ielen = 0;
  int tmp = 0;

  /** TODO: consider removing the `this->len` variable */
  ielen = this->len;

  // SMF NAS Message Decode
  encoded +=
      this->smf_msg.SmfMsgEncodeMsg(this->contents, sizeof(this->contents));
  if (encoded <= 0) {
    OAILOG_ERROR(LOG_NAS5G, "SmfMsg Encoding Failed");
    return (RETURNerror);
  }

  if (static_cast<int>(ielen) != encoded) {
    OAILOG_WARNING(
        LOG_NAS5G,
        "Length missmatch : IE length : %d, Encoded SMF message length : %d",
        ielen, encoded);
  }

  IES_ENCODE_U16(buffer, tmp, encoded);
  memcpy(buffer + tmp, this->contents, encoded);

  return (encoded + tmp);
}
}  // namespace magma5g
