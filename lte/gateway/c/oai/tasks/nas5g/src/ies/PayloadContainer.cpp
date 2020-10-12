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
#include "PayloadContainer.h"
#include "SmfMessage.h"
#include "CommonDefs.h"

using namespace std;
namespace magma5g {
PayloadContainerMsg::PayloadContainerMsg(){};
PayloadContainerMsg::~PayloadContainerMsg(){};

int PayloadContainerMsg::DecodePayloadContainerMsg(
    PayloadContainerMsg* payloadcontainer, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded    = 0;
  uint32_t ielen = 0;
  IES_DECODE_U16(buffer, decoded, ielen);
  payloadcontainer->len = ielen;
  MLOG(MDEBUG) << "DecodePayloadContainerMsg__: len = " << dec << int(payloadcontainer->len)
               << endl;
  memcpy(&payloadcontainer->contents, buffer + decoded, int(ielen));
  BUFFER_PRINT_LOG(payloadcontainer->contents, int(ielen));

  // SMF NAS Message Decode
  decoded = payloadcontainer->smfmsg.SmfMsgDecodeMsg(
      &payloadcontainer->smfmsg, payloadcontainer->contents, int(ielen));

  return (decoded);
};

int PayloadContainerMsg::EncodePayloadContainerMsg(
    PayloadContainerMsg* payloadcontainer, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int encoded    = 0;
  uint32_t ielen = 0;
  int tmp = 0;
  ielen          = payloadcontainer->len;
  IES_ENCODE_U16(buffer, encoded, ielen);
  MLOG(MDEBUG) << "DecodePayloadContainerMsg__: len = " << hex << int(ielen)
               << endl;
  tmp = encoded;

  // SMF NAS Message Decode
  encoded = payloadcontainer->smfmsg.SmfMsgEncodeMsg(
      &payloadcontainer->smfmsg, payloadcontainer->contents, payloadcontainer->len);
  BUFFER_PRINT_LOG(payloadcontainer->contents, payloadcontainer->len);
  memcpy(buffer + tmp, payloadcontainer->contents, payloadcontainer->len);

  return (encoded);
};
}  // namespace magma5g
