/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

#include <sstream>
#include <cstdint>
#include <cstring>
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GDNN.h"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
DNNMsg::DNNMsg(){};
DNNMsg::~DNNMsg(){};

// Decode DNN Message
int DNNMsg::DecodeDNNMsg(
    DNNMsg* dnn_message, uint8_t iei, uint8_t* buffer, uint32_t len) {
  int decoded   = 0;
  uint8_t ielen = 0;

  MLOG(MDEBUG) << "DecodeDNN : ";

  if (iei > 0) {
    DECODE_U8(buffer + decoded, dnn_message->iei, decoded);
    CHECK_IEI_DECODER(iei, (unsigned char) *buffer);
    MLOG(MDEBUG) << "iei : " << std::hex << static_cast<int>(dnn_message->iei);
  }
  DECODE_U8(buffer + decoded, ielen, decoded);
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  dnn_message->len = ielen;
  MLOG(MDEBUG) << "len : " << static_cast<int>(dnn_message->len);

  uint8_t dnn_len = 0;
  DECODE_U8(buffer + decoded, dnn_len, decoded);
  MLOG(MDEBUG) << "dnn_len : " << static_cast<int>(dnn_len);

  memcpy(dnn_message->dnn, buffer + decoded, dnn_len);

  decoded = decoded + dnn_len;
  MLOG(MDEBUG) << "dnn str : " << dnn_message->dnn;

  return decoded;
}

// Encode DNN Message
int DNNMsg::EncodeDNNMsg(
    DNNMsg* dnn_message, uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

  MLOG(MDEBUG) << "EncodeDNN : ";
  // Checking IEI and pointer
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, DNN_MIN_LENGTH, len);

  if (iei > 0) {
    CHECK_IEI_ENCODER(iei, (unsigned char) dnn_message->iei);
    ENCODE_U8(buffer, iei, encoded);
    MLOG(MDEBUG) << "iei : " << std::hex << static_cast<int>(dnn_message->iei);
  }

  ENCODE_U8(buffer + encoded, dnn_message->len, encoded);
  MLOG(MDEBUG) << "len : " << static_cast<int>(dnn_message->len);

  ENCODE_U8(buffer + encoded, dnn_message->len - 1, encoded);
  MLOG(MDEBUG) << "dnn len : " << dnn_message->len - 1;

  memcpy(buffer + encoded, dnn_message->dnn, dnn_message->len - 1);
  BUFFER_PRINT_LOG(buffer + encoded, dnn_message->len - 1);
  MLOG(MDEBUG) << "dnn str : " << dnn_message->dnn;
  encoded = encoded + dnn_message->len - 1;

  return encoded;
};

}  // namespace magma5g
