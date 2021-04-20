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
#include "EmmInformation.h"
#include "common_defs.h"

int decode_emm_information(
    emm_information_msg* emm_information, uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, EMM_INFORMATION_MINIMUM_LENGTH, len);

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
      case EMM_INFORMATION_FULL_NAME_FOR_NETWORK_IEI:
        if ((decoded_result = decode_network_name_ie(
                 &emm_information->fullnamefornetwork,
                 EMM_INFORMATION_FULL_NAME_FOR_NETWORK_IEI, buffer + decoded,
                 len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        emm_information->presencemask |=
            EMM_INFORMATION_FULL_NAME_FOR_NETWORK_PRESENT;
        break;

      case EMM_INFORMATION_SHORT_NAME_FOR_NETWORK_IEI:
        if ((decoded_result = decode_network_name_ie(
                 &emm_information->shortnamefornetwork,
                 EMM_INFORMATION_SHORT_NAME_FOR_NETWORK_IEI, buffer + decoded,
                 len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        emm_information->presencemask |=
            EMM_INFORMATION_SHORT_NAME_FOR_NETWORK_PRESENT;
        break;

      case EMM_INFORMATION_LOCAL_TIME_ZONE_IEI:
        if ((decoded_result = decode_time_zone(
                 &emm_information->localtimezone,
                 EMM_INFORMATION_LOCAL_TIME_ZONE_IEI, buffer + decoded,
                 len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        emm_information->presencemask |=
            EMM_INFORMATION_LOCAL_TIME_ZONE_PRESENT;
        break;

      case EMM_INFORMATION_UNIVERSAL_TIME_AND_LOCAL_TIME_ZONE_IEI:
        if ((decoded_result = decode_time_zone_and_time(
                 &emm_information->universaltimeandlocaltimezone,
                 EMM_INFORMATION_UNIVERSAL_TIME_AND_LOCAL_TIME_ZONE_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        emm_information->presencemask |=
            EMM_INFORMATION_UNIVERSAL_TIME_AND_LOCAL_TIME_ZONE_PRESENT;
        break;

      case MM_DAYLIGHT_SAVING_TIME_IEI:
        if ((decoded_result = decode_daylight_saving_time_ie(
                 &emm_information->networkdaylightsavingtime,
                 MM_DAYLIGHT_SAVING_TIME_IEI, buffer + decoded,
                 len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        emm_information->presencemask |=
            EMM_INFORMATION_NETWORK_DAYLIGHT_SAVING_TIME_PRESENT;
        break;

      default:
        errorCodeDecoder = TLV_UNEXPECTED_IEI;
        return TLV_UNEXPECTED_IEI;
    }
  }

  return decoded;
}

int encode_emm_information(
    emm_information_msg* emm_information, uint8_t* buffer, uint32_t len) {
  int encoded       = 0;
  int encode_result = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, EMM_INFORMATION_MINIMUM_LENGTH, len);

  if ((emm_information->presencemask &
       EMM_INFORMATION_FULL_NAME_FOR_NETWORK_PRESENT) ==
      EMM_INFORMATION_FULL_NAME_FOR_NETWORK_PRESENT) {
    if ((encode_result = encode_network_name_ie(
             &emm_information->fullnamefornetwork,
             EMM_INFORMATION_FULL_NAME_FOR_NETWORK_IEI, buffer + encoded,
             len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((emm_information->presencemask &
       EMM_INFORMATION_SHORT_NAME_FOR_NETWORK_PRESENT) ==
      EMM_INFORMATION_SHORT_NAME_FOR_NETWORK_PRESENT) {
    if ((encode_result = encode_network_name_ie(
             &emm_information->shortnamefornetwork,
             EMM_INFORMATION_SHORT_NAME_FOR_NETWORK_IEI, buffer + encoded,
             len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((emm_information->presencemask &
       EMM_INFORMATION_LOCAL_TIME_ZONE_PRESENT) ==
      EMM_INFORMATION_LOCAL_TIME_ZONE_PRESENT) {
    if ((encode_result = encode_time_zone(
             &emm_information->localtimezone,
             EMM_INFORMATION_LOCAL_TIME_ZONE_IEI, buffer + encoded,
             len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((emm_information->presencemask &
       EMM_INFORMATION_UNIVERSAL_TIME_AND_LOCAL_TIME_ZONE_PRESENT) ==
      EMM_INFORMATION_UNIVERSAL_TIME_AND_LOCAL_TIME_ZONE_PRESENT) {
    if ((encode_result = encode_time_zone_and_time(
             &emm_information->universaltimeandlocaltimezone,
             EMM_INFORMATION_UNIVERSAL_TIME_AND_LOCAL_TIME_ZONE_IEI,
             buffer + encoded, len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((emm_information->presencemask &
       EMM_INFORMATION_NETWORK_DAYLIGHT_SAVING_TIME_PRESENT) ==
      EMM_INFORMATION_NETWORK_DAYLIGHT_SAVING_TIME_PRESENT) {
    if ((encode_result = encode_daylight_saving_time_ie(
             &emm_information->networkdaylightsavingtime,
             MM_DAYLIGHT_SAVING_TIME_IEI, buffer + encoded, len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  return encoded;
}
