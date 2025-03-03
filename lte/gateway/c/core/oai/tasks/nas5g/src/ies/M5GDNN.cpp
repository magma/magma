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
#ifdef __cplusplus
extern "C" {
#endif
#include "lte/gateway/c/core/oai/common/log.h"
#ifdef __cplusplus
}
#endif
#include "lte/gateway/c/core/oai/tasks/nas5g/include/ies/M5GDNN.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5GCommonDefs.h"

namespace magma5g {
DNNMsg::DNNMsg() {};
DNNMsg::~DNNMsg() {};

// Decode DNN Message
int DNNMsg::DecodeDNNMsg(DNNMsg* dnn_message, uint8_t iei, uint8_t* buffer,
                         uint32_t len) {
  int decoded = 0;
  uint8_t ielen = 0;

  if (iei > 0) {
    DECODE_U8(buffer + decoded, dnn_message->iei, decoded);
    CHECK_IEI_DECODER(iei, (unsigned char)*buffer);
  }
  DECODE_U8(buffer + decoded, ielen, decoded);
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  dnn_message->len = ielen;
  uint8_t dnn_length = 0;
  uint8_t dnn_len = 0;
  DECODE_U8(buffer + decoded, dnn_len, decoded);

  if ((ielen <= dnn_len) || (dnn_len < DNN_MIN_LENGTH) ||
      (dnn_len > MAX_DNN_LENGTH)) {
    OAILOG_ERROR(LOG_NAS5G,
                 "Mismatch Length: IE length : %u, DNN String length: %u",
                 ielen, dnn_len);
    return -1;
  }

  memcpy(dnn_message->dnn, buffer + decoded, dnn_len);
  dnn_length += dnn_len;
  decoded = decoded + dnn_len;

  while (dnn_length + 1 < ielen) {
    dnn_len = 0;
    memcpy(dnn_message->dnn + dnn_length, ".", 1);
    DECODE_U8(buffer + decoded, dnn_len, decoded);
    dnn_length = dnn_length + 1;

    memcpy(dnn_message->dnn + dnn_length, buffer + decoded, dnn_len);

    decoded = decoded + dnn_len;
    dnn_length += dnn_len;
  }
  return decoded;
}

// Encode DNN Message
int DNNMsg::EncodeDNNMsg(DNNMsg* dnn_message, uint8_t iei, uint8_t* buffer,
                         uint32_t len) {
  uint32_t encoded = 0;

  // Checking IEI and pointer
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, DNN_MIN_LENGTH, len);

  if (iei > 0) {
    CHECK_IEI_ENCODER(iei, (unsigned char)dnn_message->iei);
    ENCODE_U8(buffer, iei, encoded);
  }

  ENCODE_U8(buffer + encoded, dnn_message->len, encoded);
  uint8_t dnn_length = 0;
  while (dnn_length < dnn_message->len - 1) {
    uint8_t dnn_len = 0;
    for (int i = dnn_length; (memcmp(dnn_message->dnn + i, ".", 1) != 0) &&
                             i < dnn_message->len - 1;
         i++) {
      ++dnn_len;
    }
    ENCODE_U8(buffer + encoded, dnn_len, encoded);

    memcpy(buffer + encoded, dnn_message->dnn + dnn_length, dnn_len);
    if (dnn_len + dnn_length < dnn_message->len - 1) {
      if (memcmp(dnn_message->dnn + dnn_length + dnn_len, ".", 1) == 0)
        dnn_length = dnn_length + 1;
    }
    dnn_length = dnn_length + dnn_len;

    encoded = encoded + dnn_len;
  }
  return encoded;
};

}  // namespace magma5g
