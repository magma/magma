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

#include "lte/gateway/c/core/oai/common/log.h"
#include "lte/gateway/c/core/oai/common/TLVEncoder.h"
#include "lte/gateway/c/core/oai/common/TLVDecoder.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/TrackingAreaUpdateAccept.h"
#include "lte/gateway/c/core/oai/common/common_defs.h"

int decode_tracking_area_update_accept(
    tracking_area_update_accept_msg* tracking_area_update_accept,
    uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, TRACKING_AREA_UPDATE_ACCEPT_MINIMUM_LENGTH, len);

  /*
   * Decoding mandatory fields
   */
  if ((decoded_result = decode_u8_eps_update_result(
           &tracking_area_update_accept->epsupdateresult, 0,
           *(buffer + decoded), len - decoded)) < 0)
    return decoded_result;

  decoded++;

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
      case TRACKING_AREA_UPDATE_ACCEPT_T3412_VALUE_IEI:
        if ((decoded_result = decode_gprs_timer_ie(
                 &tracking_area_update_accept->t3412value,
                 TRACKING_AREA_UPDATE_ACCEPT_T3412_VALUE_IEI, buffer + decoded,
                 len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_accept->presencemask |=
            TRACKING_AREA_UPDATE_ACCEPT_T3412_VALUE_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_ACCEPT_GUTI_IEI:
        if ((decoded_result = decode_eps_mobile_identity(
                 &tracking_area_update_accept->guti,
                 TRACKING_AREA_UPDATE_ACCEPT_GUTI_IEI, buffer + decoded,
                 len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_accept->presencemask |=
            TRACKING_AREA_UPDATE_ACCEPT_GUTI_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_ACCEPT_TAI_LIST_IEI:
        if ((decoded_result = decode_tracking_area_identity_list(
                 &tracking_area_update_accept->tailist,
                 TRACKING_AREA_UPDATE_ACCEPT_TAI_LIST_IEI, buffer + decoded,
                 len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_accept->presencemask |=
            TRACKING_AREA_UPDATE_ACCEPT_TAI_LIST_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_ACCEPT_EPS_BEARER_CONTEXT_STATUS_IEI:
        if ((decoded_result = decode_eps_bearer_context_status(
                 &tracking_area_update_accept->epsbearercontextstatus,
                 TRACKING_AREA_UPDATE_ACCEPT_EPS_BEARER_CONTEXT_STATUS_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_accept->presencemask |=
            TRACKING_AREA_UPDATE_ACCEPT_EPS_BEARER_CONTEXT_STATUS_PRESENT;
        break;

      case C_LOCATION_AREA_IDENTIFICATION_IEI:
        if ((decoded_result = decode_location_area_identification_ie(
                 &tracking_area_update_accept->locationareaidentification,
                 C_LOCATION_AREA_IDENTIFICATION_IEI, buffer + decoded,
                 len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_accept->presencemask |=
            TRACKING_AREA_UPDATE_ACCEPT_LOCATION_AREA_IDENTIFICATION_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_ACCEPT_MS_IDENTITY_IEI:
        if ((decoded_result = decode_mobile_identity_ie(
                 &tracking_area_update_accept->msidentity,
                 TRACKING_AREA_UPDATE_ACCEPT_MS_IDENTITY_IEI, buffer + decoded,
                 len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_accept->presencemask |=
            TRACKING_AREA_UPDATE_ACCEPT_MS_IDENTITY_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_ACCEPT_EMM_CAUSE_IEI:
        if ((decoded_result = decode_emm_cause(
                 &tracking_area_update_accept->emmcause,
                 TRACKING_AREA_UPDATE_ACCEPT_EMM_CAUSE_IEI, buffer + decoded,
                 len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_accept->presencemask |=
            TRACKING_AREA_UPDATE_ACCEPT_EMM_CAUSE_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_ACCEPT_T3402_VALUE_IEI:
        if ((decoded_result = decode_gprs_timer_ie(
                 &tracking_area_update_accept->t3402value,
                 TRACKING_AREA_UPDATE_ACCEPT_T3402_VALUE_IEI, buffer + decoded,
                 len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_accept->presencemask |=
            TRACKING_AREA_UPDATE_ACCEPT_T3402_VALUE_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_ACCEPT_T3423_VALUE_IEI:
        if ((decoded_result = decode_gprs_timer_ie(
                 &tracking_area_update_accept->t3423value,
                 TRACKING_AREA_UPDATE_ACCEPT_T3423_VALUE_IEI, buffer + decoded,
                 len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_accept->presencemask |=
            TRACKING_AREA_UPDATE_ACCEPT_T3423_VALUE_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_ACCEPT_EQUIVALENT_PLMNS_IEI:
        if ((decoded_result = decode_plmn_list_ie(
                 &tracking_area_update_accept->equivalentplmns,
                 TRACKING_AREA_UPDATE_ACCEPT_EQUIVALENT_PLMNS_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_accept->presencemask |=
            TRACKING_AREA_UPDATE_ACCEPT_EQUIVALENT_PLMNS_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_ACCEPT_EMERGENCY_NUMBER_LIST_IEI:
        if ((decoded_result = decode_emergency_number_list_ie(
                 &tracking_area_update_accept->emergencynumberlist,
                 TRACKING_AREA_UPDATE_ACCEPT_EMERGENCY_NUMBER_LIST_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_accept->presencemask |=
            TRACKING_AREA_UPDATE_ACCEPT_EMERGENCY_NUMBER_LIST_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_ACCEPT_EPS_NETWORK_FEATURE_SUPPORT_IEI:
        if ((decoded_result = decode_eps_network_feature_support(
                 &tracking_area_update_accept->epsnetworkfeaturesupport,
                 TRACKING_AREA_UPDATE_ACCEPT_EPS_NETWORK_FEATURE_SUPPORT_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_accept->presencemask |=
            TRACKING_AREA_UPDATE_ACCEPT_EPS_NETWORK_FEATURE_SUPPORT_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_ACCEPT_ADDITIONAL_UPDATE_RESULT_IEI:
        if ((decoded_result = decode_additional_update_result(
                 &tracking_area_update_accept->additionalupdateresult,
                 TRACKING_AREA_UPDATE_ACCEPT_ADDITIONAL_UPDATE_RESULT_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_accept->presencemask |=
            TRACKING_AREA_UPDATE_ACCEPT_ADDITIONAL_UPDATE_RESULT_PRESENT;
        break;

      default:
        errorCodeDecoder = TLV_UNEXPECTED_IEI;
        return TLV_UNEXPECTED_IEI;
    }
  }

  return decoded;
}

int encode_tracking_area_update_accept(
    tracking_area_update_accept_msg* tracking_area_update_accept,
    uint8_t* buffer, uint32_t len) {
  int encoded       = 0;
  int encode_result = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, TRACKING_AREA_UPDATE_ACCEPT_MINIMUM_LENGTH, len);
  // This is incoorect. No need to <<4 as it is already taken care in
  // encode_u8_eps_update_result(). Hence commenting it
  //*(buffer + encoded) = ((encode_u8_eps_update_result
  //(&tracking_area_update_accept->epsupdateresult) & 0x0f) << 4) | 0x00;
  *(buffer + encoded) = (encode_u8_eps_update_result(
      &tracking_area_update_accept->epsupdateresult));
  OAILOG_INFO(
      LOG_MME_APP, "epsupdateresult in encode_tracking_area_update_accept %x\n",
      *(buffer + encoded));
  encoded++;

  if ((tracking_area_update_accept->presencemask &
       TRACKING_AREA_UPDATE_ACCEPT_T3412_VALUE_PRESENT) ==
      TRACKING_AREA_UPDATE_ACCEPT_T3412_VALUE_PRESENT) {
    if ((encode_result = encode_gprs_timer_ie(
             &tracking_area_update_accept->t3412value,
             TRACKING_AREA_UPDATE_ACCEPT_T3412_VALUE_IEI, buffer + encoded,
             len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((tracking_area_update_accept->presencemask &
       TRACKING_AREA_UPDATE_ACCEPT_GUTI_PRESENT) ==
      TRACKING_AREA_UPDATE_ACCEPT_GUTI_PRESENT) {
    if ((encode_result = encode_eps_mobile_identity(
             &tracking_area_update_accept->guti,
             TRACKING_AREA_UPDATE_ACCEPT_GUTI_IEI, buffer + encoded,
             len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((tracking_area_update_accept->presencemask &
       TRACKING_AREA_UPDATE_ACCEPT_TAI_LIST_PRESENT) ==
      TRACKING_AREA_UPDATE_ACCEPT_TAI_LIST_PRESENT) {
    if ((encode_result = encode_tracking_area_identity_list(
             &tracking_area_update_accept->tailist,
             TRACKING_AREA_UPDATE_ACCEPT_TAI_LIST_IEI, buffer + encoded,
             len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((tracking_area_update_accept->presencemask &
       TRACKING_AREA_UPDATE_ACCEPT_EPS_BEARER_CONTEXT_STATUS_PRESENT) ==
      TRACKING_AREA_UPDATE_ACCEPT_EPS_BEARER_CONTEXT_STATUS_PRESENT) {
    if ((encode_result = encode_eps_bearer_context_status(
             &tracking_area_update_accept->epsbearercontextstatus,
             TRACKING_AREA_UPDATE_ACCEPT_EPS_BEARER_CONTEXT_STATUS_IEI,
             buffer + encoded, len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((tracking_area_update_accept->presencemask &
       TRACKING_AREA_UPDATE_ACCEPT_LOCATION_AREA_IDENTIFICATION_PRESENT) ==
      TRACKING_AREA_UPDATE_ACCEPT_LOCATION_AREA_IDENTIFICATION_PRESENT) {
    if ((encode_result = encode_location_area_identification_ie(
             &tracking_area_update_accept->locationareaidentification,
             C_LOCATION_AREA_IDENTIFICATION_IEI, buffer + encoded,
             len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((tracking_area_update_accept->presencemask &
       TRACKING_AREA_UPDATE_ACCEPT_MS_IDENTITY_PRESENT) ==
      TRACKING_AREA_UPDATE_ACCEPT_MS_IDENTITY_PRESENT) {
    if ((encode_result = encode_mobile_identity_ie(
             &tracking_area_update_accept->msidentity,
             TRACKING_AREA_UPDATE_ACCEPT_MS_IDENTITY_IEI, buffer + encoded,
             len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((tracking_area_update_accept->presencemask &
       TRACKING_AREA_UPDATE_ACCEPT_EMM_CAUSE_PRESENT) ==
      TRACKING_AREA_UPDATE_ACCEPT_EMM_CAUSE_PRESENT) {
    if ((encode_result = encode_emm_cause(
             &tracking_area_update_accept->emmcause,
             TRACKING_AREA_UPDATE_ACCEPT_EMM_CAUSE_IEI, buffer + encoded,
             len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((tracking_area_update_accept->presencemask &
       TRACKING_AREA_UPDATE_ACCEPT_T3402_VALUE_PRESENT) ==
      TRACKING_AREA_UPDATE_ACCEPT_T3402_VALUE_PRESENT) {
    if ((encode_result = encode_gprs_timer_ie(
             &tracking_area_update_accept->t3402value,
             TRACKING_AREA_UPDATE_ACCEPT_T3402_VALUE_IEI, buffer + encoded,
             len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((tracking_area_update_accept->presencemask &
       TRACKING_AREA_UPDATE_ACCEPT_T3423_VALUE_PRESENT) ==
      TRACKING_AREA_UPDATE_ACCEPT_T3423_VALUE_PRESENT) {
    if ((encode_result = encode_gprs_timer_ie(
             &tracking_area_update_accept->t3423value,
             TRACKING_AREA_UPDATE_ACCEPT_T3423_VALUE_IEI, buffer + encoded,
             len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((tracking_area_update_accept->presencemask &
       TRACKING_AREA_UPDATE_ACCEPT_EQUIVALENT_PLMNS_PRESENT) ==
      TRACKING_AREA_UPDATE_ACCEPT_EQUIVALENT_PLMNS_PRESENT) {
    if ((encode_result = encode_plmn_list_ie(
             &tracking_area_update_accept->equivalentplmns,
             TRACKING_AREA_UPDATE_ACCEPT_EQUIVALENT_PLMNS_IEI, buffer + encoded,
             len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((tracking_area_update_accept->presencemask &
       TRACKING_AREA_UPDATE_ACCEPT_EMERGENCY_NUMBER_LIST_PRESENT) ==
      TRACKING_AREA_UPDATE_ACCEPT_EMERGENCY_NUMBER_LIST_PRESENT) {
    if ((encode_result = encode_emergency_number_list_ie(
             &tracking_area_update_accept->emergencynumberlist,
             TRACKING_AREA_UPDATE_ACCEPT_EMERGENCY_NUMBER_LIST_IEI,
             buffer + encoded, len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((tracking_area_update_accept->presencemask &
       TRACKING_AREA_UPDATE_ACCEPT_EPS_NETWORK_FEATURE_SUPPORT_PRESENT) ==
      TRACKING_AREA_UPDATE_ACCEPT_EPS_NETWORK_FEATURE_SUPPORT_PRESENT) {
    if ((encode_result = encode_eps_network_feature_support(
             &tracking_area_update_accept->epsnetworkfeaturesupport,
             TRACKING_AREA_UPDATE_ACCEPT_EPS_NETWORK_FEATURE_SUPPORT_IEI,
             buffer + encoded, len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((tracking_area_update_accept->presencemask &
       TRACKING_AREA_UPDATE_ACCEPT_ADDITIONAL_UPDATE_RESULT_PRESENT) ==
      TRACKING_AREA_UPDATE_ACCEPT_ADDITIONAL_UPDATE_RESULT_PRESENT) {
    if ((encode_result = encode_additional_update_result(
             &tracking_area_update_accept->additionalupdateresult,
             TRACKING_AREA_UPDATE_ACCEPT_ADDITIONAL_UPDATE_RESULT_IEI,
             buffer + encoded, len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  return encoded;
}
