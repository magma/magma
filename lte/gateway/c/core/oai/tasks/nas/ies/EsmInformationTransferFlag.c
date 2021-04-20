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

#include "TLVEncoder.h"
#include "TLVDecoder.h"
#include "EsmInformationTransferFlag.h"

//------------------------------------------------------------------------------
int decode_esm_information_transfer_flag(
    esm_information_transfer_flag_t* esminformationtransferflag, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  int decoded = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, ESM_INFORMATION_TRANSFER_FLAG_MINIMUM_LENGTH, len);

  if (iei > 0) {
    CHECK_IEI_DECODER((*buffer & 0xf0), iei);
  }

  *esminformationtransferflag = *buffer & 0x1;
  decoded++;
  return decoded;
}

//------------------------------------------------------------------------------
int decode_u8_esm_information_transfer_flag(
    esm_information_transfer_flag_t* esminformationtransferflag, uint8_t iei,
    uint8_t value, uint32_t len) {
  int decoded     = 0;
  uint8_t* buffer = &value;

  *esminformationtransferflag = *buffer & 0x1;
  decoded++;
  return decoded;
}

//------------------------------------------------------------------------------
int encode_esm_information_transfer_flag(
    esm_information_transfer_flag_t* esminformationtransferflag, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  uint8_t encoded = 0;

  /*
   * Checking length and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, ESM_INFORMATION_TRANSFER_FLAG_MINIMUM_LENGTH, len);
  *(buffer + encoded) =
      0x00 | (iei & 0xf0) | (*esminformationtransferflag & 0x1);
  encoded++;
  return encoded;
}

//------------------------------------------------------------------------------
uint8_t encode_u8_esm_information_transfer_flag(
    esm_information_transfer_flag_t* esminformationtransferflag) {
  uint8_t bufferReturn;
  uint8_t* buffer = &bufferReturn;
  uint8_t encoded = 0;
  uint8_t iei     = 0;

  *(buffer + encoded) =
      0x00 | (iei & 0xf0) | (*esminformationtransferflag & 0x1);
  encoded++;
  return bufferReturn;
}
