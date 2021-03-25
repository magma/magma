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

#include "log.h"
#include "TLVEncoder.h"
#include "TLVDecoder.h"
#include "NASSecurityModeCommand.h"
#include "3gpp_24.008.h"
#include "NasKeySetIdentifier.h"
#include "NasSecurityAlgorithms.h"
#include "Nonce.h"
#include "UeSecurityCapability.h"
#include "common_defs.h"
#include "UeAdditionalSecurityCapability.h"

int decode_security_mode_command(
    security_mode_command_msg* security_mode_command, uint8_t* buffer,
    uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, SECURITY_MODE_COMMAND_MINIMUM_LENGTH, len);

  /*
   * Decoding mandatory fields
   */
  if ((decoded_result = decode_nas_security_algorithms(
           &security_mode_command->selectednassecurityalgorithms, 0,
           buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  if ((decoded_result = decode_u8_nas_key_set_identifier(
           &security_mode_command->naskeysetidentifier, 0,
           *(buffer + decoded) & 0x0f, len - decoded)) < 0)
    return decoded_result;

  decoded++;

  if ((decoded_result = decode_ue_security_capability(
           &security_mode_command->replayeduesecuritycapabilities, 0,
           buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  /*
   * Decoding optional fields
   */
  while (len > decoded) {
    uint8_t ieiDecoded = *(buffer + decoded);

    /*
     * Type | value iei are below 0x80 so just return the first 4 bits
     */
    if (ieiDecoded >= SECURITY_MODE_COMMAND_IMEISV_REQUEST_IEI)
      ieiDecoded = ieiDecoded & 0xf0;

    switch (ieiDecoded) {
      case SECURITY_MODE_COMMAND_IMEISV_REQUEST_IEI:
        if ((decoded_result = decode_imeisv_request_ie(
                 &security_mode_command->imeisvrequest,
                 SECURITY_MODE_COMMAND_IMEISV_REQUEST_IEI, buffer + decoded,
                 len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        security_mode_command->presencemask |=
            SECURITY_MODE_COMMAND_IMEISV_REQUEST_PRESENT;
        break;

      case SECURITY_MODE_COMMAND_REPLAYED_NONCEUE_IEI:
        if ((decoded_result = decode_nonce(
                 &security_mode_command->replayednonceue,
                 SECURITY_MODE_COMMAND_REPLAYED_NONCEUE_IEI, buffer + decoded,
                 len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        security_mode_command->presencemask |=
            SECURITY_MODE_COMMAND_REPLAYED_NONCEUE_PRESENT;
        break;

      case SECURITY_MODE_COMMAND_NONCEMME_IEI:
        if ((decoded_result = decode_nonce(
                 &security_mode_command->noncemme,
                 SECURITY_MODE_COMMAND_NONCEMME_IEI, buffer + decoded,
                 len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        security_mode_command->presencemask |=
            SECURITY_MODE_COMMAND_NONCEMME_PRESENT;
        break;

      case SECURITY_MODE_COMMAND_REPLAYED_UE_ADDITIONAL_SECURITY_CAPABILITY_IEI:
        if ((decoded_result = decode_ue_additional_security_capability(
                 &security_mode_command
                      ->replayedueadditionalsecuritycapabilities,
                 SECURITY_MODE_COMMAND_REPLAYED_UE_ADDITIONAL_SECURITY_CAPABILITY_IEI,
                 buffer + decoded, len - decoded)) <= 0) {
          return decoded_result;
        }
        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        security_mode_command->presencemask |=
            SECURITY_MODE_COMMAND_REPLAYED_UE_ADDITIONAL_SECU_CAPABILITY_PRESENT;
        break;

      default:
        errorCodeDecoder = TLV_UNEXPECTED_IEI;
        OAILOG_ERROR(
            LOG_NAS_EMM,
            "Failed to decode SECURITY_MODE_COMMAND unexpected IEI 0x%02x\n",
            ieiDecoded);
        return TLV_UNEXPECTED_IEI;
    }
  }

  return decoded;
}

int encode_security_mode_command(
    security_mode_command_msg* security_mode_command, uint8_t* buffer,
    uint32_t len) {
  int encoded       = 0;
  int encode_result = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, SECURITY_MODE_COMMAND_MINIMUM_LENGTH, len);

  if ((encode_result = encode_nas_security_algorithms(
           &security_mode_command->selectednassecurityalgorithms, 0,
           buffer + encoded,
           len - encoded)) < 0)  // Return in case of error
    return encode_result;
  else
    encoded += encode_result;

  *(buffer + encoded) =
      (encode_u8_nas_key_set_identifier(
           &security_mode_command->naskeysetidentifier) &
       0x0f);
  encoded++;

  if ((encode_result = encode_ue_security_capability(
           &security_mode_command->replayeduesecuritycapabilities, 0,
           buffer + encoded,
           len - encoded)) < 0)  // Return in case of error
    return encode_result;
  else
    encoded += encode_result;

  if ((security_mode_command->presencemask &
       SECURITY_MODE_COMMAND_IMEISV_REQUEST_PRESENT) ==
      SECURITY_MODE_COMMAND_IMEISV_REQUEST_PRESENT) {
    if ((encode_result = encode_imeisv_request_ie(
             &security_mode_command->imeisvrequest,
             SECURITY_MODE_COMMAND_IMEISV_REQUEST_IEI, buffer + encoded,
             len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((security_mode_command->presencemask &
       SECURITY_MODE_COMMAND_REPLAYED_NONCEUE_PRESENT) ==
      SECURITY_MODE_COMMAND_REPLAYED_NONCEUE_PRESENT) {
    if ((encode_result = encode_nonce(
             &security_mode_command->replayednonceue,
             SECURITY_MODE_COMMAND_REPLAYED_NONCEUE_IEI, buffer + encoded,
             len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((security_mode_command->presencemask &
       SECURITY_MODE_COMMAND_NONCEMME_PRESENT) ==
      SECURITY_MODE_COMMAND_NONCEMME_PRESENT) {
    if ((encode_result = encode_nonce(
             &security_mode_command->noncemme,
             SECURITY_MODE_COMMAND_NONCEMME_IEI, buffer + encoded,
             len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((security_mode_command->presencemask &
       SECURITY_MODE_COMMAND_REPLAYED_UE_ADDITIONAL_SECU_CAPABILITY_PRESENT) ==
      SECURITY_MODE_COMMAND_REPLAYED_UE_ADDITIONAL_SECU_CAPABILITY_PRESENT) {
    if ((encode_result = encode_ue_additional_security_capability(
             &security_mode_command->replayedueadditionalsecuritycapabilities,
             SECURITY_MODE_COMMAND_REPLAYED_UE_ADDITIONAL_SECURITY_CAPABILITY_IEI,
             buffer + encoded, len - encoded)) < 0) {
      // Return in case of error
      return encode_result;
    } else {
      encoded += encode_result;
    }
  }
  return encoded;
}