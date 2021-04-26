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
#include "M5GDNN.h"
#include "M5GCommonDefs.h"

namespace magma5g {
DNNMsg::DNNMsg(){};
DNNMsg::~DNNMsg(){};

// Decode DNN Message
int DNNMsg::DecodeDNNMsg(
    DNNMsg* dnn_message, uint8_t iei, uint8_t* buffer, uint32_t len) {
  int decoded   = 0;
  uint8_t ielen = 0;

  /*** Will be supported POST MVC ***/
  if (iei > 0) {
    CHECK_IEI_DECODER(iei, (unsigned char) *buffer);
    IES_DECODE_U8(buffer, decoded, ielen);
  }
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  return decoded;
};

// Encode DNN Message
int DNNMsg::EncodeDNNMsg(
    DNNMsg* dnn_message, uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

  // Checking IEI and pointer
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, DNN_MIN_LENGTH, len);

  if (iei > 0) {
    CHECK_IEI_ENCODER(iei, (unsigned char) dnn_message->iei);
    *buffer = iei;
    encoded++;
  }

  MLOG(MDEBUG) << "EncodeDNN : ";
  IES_ENCODE_U8(buffer, encoded, dnn_message->len);
  MLOG(MDEBUG) << "Length = " << std::hex << int(dnn_message->len);
  std::copy(dnn_message->dnn.begin(), dnn_message->dnn.end(), buffer + encoded);
  BUFFER_PRINT_LOG(buffer + encoded, (int) dnn_message->dnn.length());
  encoded = encoded + dnn_message->dnn.length();

  return encoded;
};
}  // namespace magma5g
