/*
 * Licensed to the OpenAirInterface (OAI) Software Alliance under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The OpenAirInterface Software Alliance licenses this file to You under
 * the terms found in the LICENSE file in the root of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *-------------------------------------------------------------------------------
 * For more information about the OpenAirInterface (OAI) Software Alliance:
 *      contact@openairinterface.org
 */

#include <stdint.h>
#include <sys/types.h>

#include "TLVEncoder.h"
#include "TLVDecoder.h"
#include "EsmMessageContainer.h"
#include "log.h"
#include "common_defs.h"

//------------------------------------------------------------------------------
int decode_esm_message_container(
    EsmMessageContainer* esmmessagecontainer, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded = 0;
  int decode_result;
  uint16_t ielen;

  OAILOG_FUNC_IN(LOG_NAS_ESM);

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  DECODE_LENGTH_U16(buffer + decoded, ielen, decoded);
  CHECK_LENGTH_DECODER(len - decoded, ielen);

  if ((decode_result = decode_bstring(
           esmmessagecontainer, ielen, buffer + decoded, len - decoded)) < 0) {
    OAILOG_FUNC_RETURN(LOG_NAS_ESM, decode_result);
  } else {
    decoded += decode_result;
  }

  OAILOG_FUNC_RETURN(LOG_NAS_ESM, decoded);
}

//------------------------------------------------------------------------------
int encode_esm_message_container(
    EsmMessageContainer esmmessagecontainer, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;
  int32_t encode_result;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, ESM_MESSAGE_CONTAINER_MINIMUM_LENGTH, len);

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  lenPtr = (buffer + encoded);

  if ((encode_result = encode_bstring(
           esmmessagecontainer, lenPtr + sizeof(uint16_t),
           len - sizeof(uint16_t))) < 0)
    return encode_result;
  else
    encoded += encode_result;

  ENCODE_U16(lenPtr, encode_result, encoded);
#if 0
  lenPtr[1] = (((encoded - 2 - ((iei > 0) ? 1 : 0))) & 0x0000ff00) >> 8;
  lenPtr[0] = ((encoded - 2 - ((iei > 0) ? 1 : 0))) & 0x000000ff;
#endif
  return encoded;
}
