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
#include "NASSecurityModeComplete.h"
#include "3gpp_24.008.h"
#include "common_defs.h"

int decode_security_mode_complete(
    security_mode_complete_msg* security_mode_complete, uint8_t* buffer,
    uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, SECURITY_MODE_COMPLETE_MINIMUM_LENGTH, len);

  /*
   * Decoding mandatory fields
   */
  /*
   * Decoding optional fields
   */
  while (len > decoded) {
    uint8_t ieiDecoded = *(buffer + decoded);

    /*
     * Type | value iei are below 0x80 so just return the first 4 bits
     */
    if (ieiDecoded >= 0x80) ieiDecoded = ieiDecoded & 0xf0;

    switch (ieiDecoded) {
      case SECURITY_MODE_COMPLETE_IMEISV_IEI:
        if ((decoded_result = decode_mobile_identity_ie(
                 &security_mode_complete->imeisv,
                 SECURITY_MODE_COMPLETE_IMEISV_IEI, buffer + decoded,
                 len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        security_mode_complete->presencemask |=
            SECURITY_MODE_COMPLETE_IMEISV_PRESENT;
        break;

      default:
        errorCodeDecoder = TLV_UNEXPECTED_IEI;
        return TLV_UNEXPECTED_IEI;
    }
  }

  return decoded;
}

int encode_security_mode_complete(
    security_mode_complete_msg* security_mode_complete, uint8_t* buffer,
    uint32_t len) {
  int encoded       = 0;
  int encode_result = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, SECURITY_MODE_COMPLETE_MINIMUM_LENGTH, len);

  if ((security_mode_complete->presencemask &
       SECURITY_MODE_COMPLETE_IMEISV_PRESENT) ==
      SECURITY_MODE_COMPLETE_IMEISV_PRESENT) {
    if ((encode_result = encode_mobile_identity_ie(
             &security_mode_complete->imeisv, SECURITY_MODE_COMPLETE_IMEISV_IEI,
             buffer + encoded, len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  return encoded;
}
