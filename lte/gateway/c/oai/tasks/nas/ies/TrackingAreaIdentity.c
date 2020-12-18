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
#include <string.h>

#include "TLVEncoder.h"
#include "TLVDecoder.h"
#include "TrackingAreaIdentity.h"

//------------------------------------------------------------------------------
int decode_tracking_area_identity(
    tai_t* tai, uint8_t iei, uint8_t* buffer, uint32_t len) {
  int decoded = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  tai->plmn.mcc_digit2 = (*(buffer + decoded) >> 4) & 0xf;
  tai->plmn.mcc_digit1 = *(buffer + decoded) & 0xf;
  decoded++;
  tai->plmn.mnc_digit3 = (*(buffer + decoded) >> 4) & 0xf;
  tai->plmn.mcc_digit3 = *(buffer + decoded) & 0xf;
  decoded++;
  tai->plmn.mnc_digit2 = (*(buffer + decoded) >> 4) & 0xf;
  tai->plmn.mnc_digit1 = *(buffer + decoded) & 0xf;
  decoded++;
  // IES_DECODE_U16(tai->tac, *(buffer + decoded));
  IES_DECODE_U16(buffer, decoded, tai->tac);
  return decoded;
}

//------------------------------------------------------------------------------
int encode_tracking_area_identity(
    tai_t* tai, uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, TRACKING_AREA_IDENTITY_MINIMUM_LENGTH, len);

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  *(buffer + encoded) =
      0x00 | ((tai->plmn.mcc_digit2 & 0xf) << 4) | (tai->plmn.mcc_digit1 & 0xf);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((tai->plmn.mnc_digit3 & 0xf) << 4) | (tai->plmn.mcc_digit3 & 0xf);
  encoded++;
  *(buffer + encoded) =
      0x00 | ((tai->plmn.mnc_digit2 & 0xf) << 4) | (tai->plmn.mnc_digit1 & 0xf);
  encoded++;
  IES_ENCODE_U16(buffer, encoded, tai->tac);
  return encoded;
}

//------------------------------------------------------------------------------
/* Clear TAI without free it */
void clear_tai(tai_t* const tai) {
  memset(tai, 0, sizeof(tai_t));
}
