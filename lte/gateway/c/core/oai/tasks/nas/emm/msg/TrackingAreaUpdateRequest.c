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

#include "lte/gateway/c/core/oai/common/TLVEncoder.h"
#include "lte/gateway/c/core/oai/common/TLVDecoder.h"
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/TrackingAreaUpdateRequest.h"
#include "lte/gateway/c/core/oai/tasks/nas/ies/UeNetworkCapability.h"
#include "lte/gateway/c/core/oai/common/common_defs.h"

int decode_tracking_area_update_request(
    tracking_area_update_request_msg* tracking_area_update_request,
    uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;

  /* Check if we got a NULL pointer and if buffer length is >=
   * minimum length expected for the message.
   */
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, TRACKING_AREA_UPDATE_REQUEST_MINIMUM_LENGTH, len);

  /*
   * Decoding mandatory fields
   */
  if ((decoded_result = decode_u8_eps_update_type(
           &tracking_area_update_request->epsupdatetype, 0,
           *(buffer + decoded) & 0x0f, len - decoded)) < 0)
    return decoded_result;

  if ((decoded_result = decode_u8_nas_key_set_identifier(
           &tracking_area_update_request->naskeysetidentifier, 0,
           *(buffer + decoded) >> 4, len - decoded)) < 0)
    return decoded_result;

  decoded++;

  if ((decoded_result = decode_eps_mobile_identity(
           &tracking_area_update_request->oldguti, 0, buffer + decoded,
           len - decoded)) < 0)
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
    if (ieiDecoded >= 0x80) ieiDecoded = ieiDecoded & 0xf0;

    switch (ieiDecoded) {
      case TRACKING_AREA_UPDATE_REQUEST_NONCURRENT_NATIVE_NAS_KEY_SET_IDENTIFIER_IEI:
        if ((decoded_result = decode_nas_key_set_identifier(
                 &tracking_area_update_request
                      ->noncurrentnativenaskeysetidentifier,
                 TRACKING_AREA_UPDATE_REQUEST_NONCURRENT_NATIVE_NAS_KEY_SET_IDENTIFIER_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_request->presencemask |=
            TRACKING_AREA_UPDATE_REQUEST_NONCURRENT_NATIVE_NAS_KEY_SET_IDENTIFIER_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_REQUEST_GPRS_CIPHERING_KEY_SEQUENCE_NUMBER_IEI:
        if ((decoded_result = decode_ciphering_key_sequence_number_ie(
                 &tracking_area_update_request->gprscipheringkeysequencenumber,
                 TRACKING_AREA_UPDATE_REQUEST_GPRS_CIPHERING_KEY_SEQUENCE_NUMBER_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_request->presencemask |=
            TRACKING_AREA_UPDATE_REQUEST_GPRS_CIPHERING_KEY_SEQUENCE_NUMBER_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_REQUEST_OLD_PTMSI_SIGNATURE_IEI:
        if ((decoded_result = decode_p_tmsi_signature_ie(
                 &tracking_area_update_request->oldptmsisignature,
                 TRACKING_AREA_UPDATE_REQUEST_OLD_PTMSI_SIGNATURE_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_request->presencemask |=
            TRACKING_AREA_UPDATE_REQUEST_OLD_PTMSI_SIGNATURE_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_REQUEST_ADDITIONAL_GUTI_IEI:
        if ((decoded_result = decode_eps_mobile_identity(
                 &tracking_area_update_request->additionalguti,
                 TRACKING_AREA_UPDATE_REQUEST_ADDITIONAL_GUTI_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_request->presencemask |=
            TRACKING_AREA_UPDATE_REQUEST_ADDITIONAL_GUTI_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_REQUEST_NONCEUE_IEI:
        if ((decoded_result = decode_nonce(
                 &tracking_area_update_request->nonceue,
                 TRACKING_AREA_UPDATE_REQUEST_NONCEUE_IEI, buffer + decoded,
                 len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_request->presencemask |=
            TRACKING_AREA_UPDATE_REQUEST_NONCEUE_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_REQUEST_UE_NETWORK_CAPABILITY_IEI:
        if ((decoded_result = decode_ue_network_capability(
                 &tracking_area_update_request->uenetworkcapability,
                 TRACKING_AREA_UPDATE_REQUEST_UE_NETWORK_CAPABILITY_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_request->presencemask |=
            TRACKING_AREA_UPDATE_REQUEST_UE_NETWORK_CAPABILITY_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_REQUEST_LAST_VISITED_REGISTERED_TAI_IEI:
        if ((decoded_result = decode_tracking_area_identity(
                 &tracking_area_update_request->lastvisitedregisteredtai,
                 TRACKING_AREA_UPDATE_REQUEST_LAST_VISITED_REGISTERED_TAI_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_request->presencemask |=
            TRACKING_AREA_UPDATE_REQUEST_LAST_VISITED_REGISTERED_TAI_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_REQUEST_DRX_PARAMETER_IEI:
        if ((decoded_result = decode_drx_parameter_ie(
                 &tracking_area_update_request->drxparameter,
                 TRACKING_AREA_UPDATE_REQUEST_DRX_PARAMETER_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_request->presencemask |=
            TRACKING_AREA_UPDATE_REQUEST_DRX_PARAMETER_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_REQUEST_UE_RADIO_CAPABILITY_INFORMATION_UPDATE_NEEDED_IEI:
        if ((decoded_result = decode_ue_radio_capability_information_update_needed(
                 &tracking_area_update_request
                      ->ueradiocapabilityinformationupdateneeded,
                 TRACKING_AREA_UPDATE_REQUEST_UE_RADIO_CAPABILITY_INFORMATION_UPDATE_NEEDED_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_request->presencemask |=
            TRACKING_AREA_UPDATE_REQUEST_UE_RADIO_CAPABILITY_INFORMATION_UPDATE_NEEDED_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_REQUEST_EPS_BEARER_CONTEXT_STATUS_IEI:
        if ((decoded_result = decode_eps_bearer_context_status(
                 &tracking_area_update_request->epsbearercontextstatus,
                 TRACKING_AREA_UPDATE_REQUEST_EPS_BEARER_CONTEXT_STATUS_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_request->presencemask |=
            TRACKING_AREA_UPDATE_REQUEST_EPS_BEARER_CONTEXT_STATUS_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_REQUEST_MS_NETWORK_CAPABILITY_IEI:
        if ((decoded_result = decode_ms_network_capability_ie(
                 &tracking_area_update_request->msnetworkcapability,
                 TRACKING_AREA_UPDATE_REQUEST_MS_NETWORK_CAPABILITY_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_request->presencemask |=
            TRACKING_AREA_UPDATE_REQUEST_MS_NETWORK_CAPABILITY_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_REQUEST_OLD_LOCATION_AREA_IDENTIFICATION_IEI:
        if ((decoded_result = decode_location_area_identification_ie(
                 &tracking_area_update_request->oldlocationareaidentification,
                 TRACKING_AREA_UPDATE_REQUEST_OLD_LOCATION_AREA_IDENTIFICATION_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_request->presencemask |=
            TRACKING_AREA_UPDATE_REQUEST_OLD_LOCATION_AREA_IDENTIFICATION_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_REQUEST_TMSI_STATUS_IEI:
        if ((decoded_result = decode_tmsi_status(
                 &tracking_area_update_request->tmsistatus,
                 TRACKING_AREA_UPDATE_REQUEST_TMSI_STATUS_IEI, buffer + decoded,
                 len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_request->presencemask |=
            TRACKING_AREA_UPDATE_REQUEST_TMSI_STATUS_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_REQUEST_MOBILE_STATION_CLASSMARK_2_IEI:
        if ((decoded_result = decode_mobile_station_classmark_2_ie(
                 &tracking_area_update_request->mobilestationclassmark2,
                 TRACKING_AREA_UPDATE_REQUEST_MOBILE_STATION_CLASSMARK_2_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_request->presencemask |=
            TRACKING_AREA_UPDATE_REQUEST_MOBILE_STATION_CLASSMARK_2_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_REQUEST_MOBILE_STATION_CLASSMARK_3_IEI:
        if ((decoded_result = decode_mobile_station_classmark_3_ie(
                 &tracking_area_update_request->mobilestationclassmark3,
                 TRACKING_AREA_UPDATE_REQUEST_MOBILE_STATION_CLASSMARK_3_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_request->presencemask |=
            TRACKING_AREA_UPDATE_REQUEST_MOBILE_STATION_CLASSMARK_3_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_REQUEST_SUPPORTED_CODECS_IEI:
        if ((decoded_result = decode_supported_codec_list_ie(
                 &tracking_area_update_request->supportedcodecs,
                 TRACKING_AREA_UPDATE_REQUEST_SUPPORTED_CODECS_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_request->presencemask |=
            TRACKING_AREA_UPDATE_REQUEST_SUPPORTED_CODECS_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_REQUEST_ADDITIONAL_UPDATE_TYPE_IEI:
        if ((decoded_result = decode_additional_update_type(
                 &tracking_area_update_request->additionalupdatetype,
                 TRACKING_AREA_UPDATE_REQUEST_ADDITIONAL_UPDATE_TYPE_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_request->presencemask |=
            TRACKING_AREA_UPDATE_REQUEST_ADDITIONAL_UPDATE_TYPE_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_REQUEST_OLD_GUTI_TYPE_IEI:
        if ((decoded_result = decode_guti_type(
                 &tracking_area_update_request->oldgutitype,
                 TRACKING_AREA_UPDATE_REQUEST_OLD_GUTI_TYPE_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_request->presencemask |=
            TRACKING_AREA_UPDATE_REQUEST_OLD_GUTI_TYPE_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_REQUEST_VOICE_DOMAIN_PREFERENCE_IEI:
        if ((decoded_result =
                 decode_voice_domain_preference_and_ue_usage_setting(
                     &tracking_area_update_request
                          ->voicedomainpreferenceandueusagesetting,
                     true, buffer + decoded, len - decoded)) <= 0) {
          OAILOG_FUNC_RETURN(LOG_NAS_EMM, decoded_result);
        }

        decoded += decoded_result;

        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_request->presencemask |=
            TRACKING_AREA_UPDATE_REQUEST_VOICE_DOMAIN_PREFERENCE_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_REQUEST_MS_NETWORK_FEATURE_SUPPORT_IEI:
        if ((decoded_result = decode_ms_network_feature_support_ie(
                 &tracking_area_update_request->msnetworkfeaturesupport,
                 TRACKING_AREA_UPDATE_REQUEST_MS_NETWORK_FEATURE_SUPPORT_IEI,
                 buffer + decoded, len - decoded)) <= 0) {
          // return decoded_result;
          OAILOG_FUNC_RETURN(LOG_NAS_EMM, decoded_result);
        }

        decoded += decoded_result;
        /* Set corresponding mask to 1 in presencemask */
        tracking_area_update_request->presencemask |=
            TRACKING_AREA_UPDATE_REQUEST_MS_NETWORK_FEATURE_SUPPORT_PRESENT;
        break;

      case TRACKING_AREA_UPDATE_REQUEST_UE_ADDITIONAL_SECURITY_CAPABILITY_IEI:
        if ((decoded_result = decode_ue_additional_security_capability(
                 &tracking_area_update_request->ueadditionalsecuritycapability,
                 TRACKING_AREA_UPDATE_REQUEST_UE_ADDITIONAL_SECURITY_CAPABILITY_IEI,
                 buffer + decoded, len - decoded)) <= 0) {
          OAILOG_FUNC_RETURN(LOG_NAS_EMM, decoded_result);
        }

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        tracking_area_update_request->presencemask |=
            TRACKING_AREA_UPDATE_REQUEST_UE_ADDITIONAL_SECURITY_CAPABILITY_PRESENT;
        break;

      default:
        errorCodeDecoder = TLV_UNEXPECTED_IEI;
        return TLV_UNEXPECTED_IEI;
    }
  }

  return decoded;
}

int encode_tracking_area_update_request(
    tracking_area_update_request_msg* tracking_area_update_request,
    uint8_t* buffer, uint32_t len) {
  int encoded       = 0;
  int encode_result = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, TRACKING_AREA_UPDATE_REQUEST_MINIMUM_LENGTH, len);
  *(buffer + encoded) =
      ((encode_u8_nas_key_set_identifier(
            &tracking_area_update_request->naskeysetidentifier) &
        0x0f)
       << 4) |
      (encode_u8_eps_update_type(&tracking_area_update_request->epsupdatetype) &
       0x0f);
  encoded++;

  if ((encode_result = encode_eps_mobile_identity(
           &tracking_area_update_request->oldguti, 0, buffer + encoded,
           len - encoded)) < 0) {  // Return in case of error
    return encode_result;
  } else {
    encoded += encode_result;
  }

  if ((tracking_area_update_request->presencemask &
       TRACKING_AREA_UPDATE_REQUEST_NONCURRENT_NATIVE_NAS_KEY_SET_IDENTIFIER_PRESENT) ==
      TRACKING_AREA_UPDATE_REQUEST_NONCURRENT_NATIVE_NAS_KEY_SET_IDENTIFIER_PRESENT) {
    if ((encode_result = encode_nas_key_set_identifier(
             &tracking_area_update_request->noncurrentnativenaskeysetidentifier,
             TRACKING_AREA_UPDATE_REQUEST_NONCURRENT_NATIVE_NAS_KEY_SET_IDENTIFIER_IEI,
             buffer + encoded, len - encoded)) < 0) {
      // Return in case of error
      return encode_result;
    } else {
      encoded += encode_result;
    }
  }

  if ((tracking_area_update_request->presencemask &
       TRACKING_AREA_UPDATE_REQUEST_GPRS_CIPHERING_KEY_SEQUENCE_NUMBER_PRESENT) ==
      TRACKING_AREA_UPDATE_REQUEST_GPRS_CIPHERING_KEY_SEQUENCE_NUMBER_PRESENT) {
    if ((encode_result = encode_ciphering_key_sequence_number_ie(
             &tracking_area_update_request->gprscipheringkeysequencenumber,
             TRACKING_AREA_UPDATE_REQUEST_GPRS_CIPHERING_KEY_SEQUENCE_NUMBER_IEI,
             buffer + encoded, len - encoded)) < 0) {
      // Return in case of error
      return encode_result;
    } else {
      encoded += encode_result;
    }
  }

  if ((tracking_area_update_request->presencemask &
       TRACKING_AREA_UPDATE_REQUEST_OLD_PTMSI_SIGNATURE_PRESENT) ==
      TRACKING_AREA_UPDATE_REQUEST_OLD_PTMSI_SIGNATURE_PRESENT) {
    if ((encode_result = encode_p_tmsi_signature_ie(
             tracking_area_update_request->oldptmsisignature,
             TRACKING_AREA_UPDATE_REQUEST_OLD_PTMSI_SIGNATURE_IEI,
             buffer + encoded, len - encoded)) < 0) {
      // Return in case of error
      return encode_result;
    } else {
      encoded += encode_result;
    }
  }

  if ((tracking_area_update_request->presencemask &
       TRACKING_AREA_UPDATE_REQUEST_ADDITIONAL_GUTI_PRESENT) ==
      TRACKING_AREA_UPDATE_REQUEST_ADDITIONAL_GUTI_PRESENT) {
    if ((encode_result = encode_eps_mobile_identity(
             &tracking_area_update_request->additionalguti,
             TRACKING_AREA_UPDATE_REQUEST_ADDITIONAL_GUTI_IEI, buffer + encoded,
             len - encoded)) < 0) {
      // Return in case of error
      return encode_result;
    } else {
      encoded += encode_result;
    }
  }

  if ((tracking_area_update_request->presencemask &
       TRACKING_AREA_UPDATE_REQUEST_NONCEUE_PRESENT) ==
      TRACKING_AREA_UPDATE_REQUEST_NONCEUE_PRESENT) {
    if ((encode_result = encode_nonce(
             &tracking_area_update_request->nonceue,
             TRACKING_AREA_UPDATE_REQUEST_NONCEUE_IEI, buffer + encoded,
             len - encoded)) < 0) {
      // Return in case of error
      return encode_result;
    } else {
      encoded += encode_result;
    }
  }

  if ((tracking_area_update_request->presencemask &
       TRACKING_AREA_UPDATE_REQUEST_UE_NETWORK_CAPABILITY_PRESENT) ==
      TRACKING_AREA_UPDATE_REQUEST_UE_NETWORK_CAPABILITY_PRESENT) {
    if ((encode_result = encode_ue_network_capability(
             &tracking_area_update_request->uenetworkcapability,
             TRACKING_AREA_UPDATE_REQUEST_UE_NETWORK_CAPABILITY_IEI,
             buffer + encoded, len - encoded)) < 0) {
      // Return in case of error
      return encode_result;
    } else {
      encoded += encode_result;
    }
  }

  if ((tracking_area_update_request->presencemask &
       TRACKING_AREA_UPDATE_REQUEST_LAST_VISITED_REGISTERED_TAI_PRESENT) ==
      TRACKING_AREA_UPDATE_REQUEST_LAST_VISITED_REGISTERED_TAI_PRESENT) {
    if ((encode_result = encode_tracking_area_identity(
             &tracking_area_update_request->lastvisitedregisteredtai,
             TRACKING_AREA_UPDATE_REQUEST_LAST_VISITED_REGISTERED_TAI_IEI,
             buffer + encoded, len - encoded)) < 0) {
      // Return in case of error
      return encode_result;
    } else {
      encoded += encode_result;
    }
  }

  if ((tracking_area_update_request->presencemask &
       TRACKING_AREA_UPDATE_REQUEST_DRX_PARAMETER_PRESENT) ==
      TRACKING_AREA_UPDATE_REQUEST_DRX_PARAMETER_PRESENT) {
    if ((encode_result = encode_drx_parameter_ie(
             &tracking_area_update_request->drxparameter,
             TRACKING_AREA_UPDATE_REQUEST_DRX_PARAMETER_IEI, buffer + encoded,
             len - encoded)) < 0) {
      // Return in case of error
      return encode_result;
    } else {
      encoded += encode_result;
    }
  }

  if ((tracking_area_update_request->presencemask &
       TRACKING_AREA_UPDATE_REQUEST_UE_RADIO_CAPABILITY_INFORMATION_UPDATE_NEEDED_PRESENT) ==
      TRACKING_AREA_UPDATE_REQUEST_UE_RADIO_CAPABILITY_INFORMATION_UPDATE_NEEDED_PRESENT) {
    if ((encode_result = encode_ue_radio_capability_information_update_needed(
             &tracking_area_update_request
                  ->ueradiocapabilityinformationupdateneeded,
             TRACKING_AREA_UPDATE_REQUEST_UE_RADIO_CAPABILITY_INFORMATION_UPDATE_NEEDED_IEI,
             buffer + encoded, len - encoded)) < 0) {
      // Return in case of error
      return encode_result;
    } else {
      encoded += encode_result;
    }
  }

  if ((tracking_area_update_request->presencemask &
       TRACKING_AREA_UPDATE_REQUEST_EPS_BEARER_CONTEXT_STATUS_PRESENT) ==
      TRACKING_AREA_UPDATE_REQUEST_EPS_BEARER_CONTEXT_STATUS_PRESENT) {
    if ((encode_result = encode_eps_bearer_context_status(
             &tracking_area_update_request->epsbearercontextstatus,
             TRACKING_AREA_UPDATE_REQUEST_EPS_BEARER_CONTEXT_STATUS_IEI,
             buffer + encoded, len - encoded)) < 0) {
      // Return in case of error
      return encode_result;
    } else {
      encoded += encode_result;
    }
  }

  if ((tracking_area_update_request->presencemask &
       TRACKING_AREA_UPDATE_REQUEST_MS_NETWORK_CAPABILITY_PRESENT) ==
      TRACKING_AREA_UPDATE_REQUEST_MS_NETWORK_CAPABILITY_PRESENT) {
    if ((encode_result = encode_ms_network_capability_ie(
             &tracking_area_update_request->msnetworkcapability,
             TRACKING_AREA_UPDATE_REQUEST_MS_NETWORK_CAPABILITY_IEI,
             buffer + encoded, len - encoded)) < 0) {
      // Return in case of error
      return encode_result;
    } else {
      encoded += encode_result;
    }
  }

  if ((tracking_area_update_request->presencemask &
       TRACKING_AREA_UPDATE_REQUEST_OLD_LOCATION_AREA_IDENTIFICATION_PRESENT) ==
      TRACKING_AREA_UPDATE_REQUEST_OLD_LOCATION_AREA_IDENTIFICATION_PRESENT) {
    if ((encode_result = encode_location_area_identification_ie(
             &tracking_area_update_request->oldlocationareaidentification,
             TRACKING_AREA_UPDATE_REQUEST_OLD_LOCATION_AREA_IDENTIFICATION_IEI,
             buffer + encoded, len - encoded)) < 0) {
      // Return in case of error
      return encode_result;
    } else {
      encoded += encode_result;
    }
  }

  if ((tracking_area_update_request->presencemask &
       TRACKING_AREA_UPDATE_REQUEST_TMSI_STATUS_PRESENT) ==
      TRACKING_AREA_UPDATE_REQUEST_TMSI_STATUS_PRESENT) {
    if ((encode_result = encode_tmsi_status(
             &tracking_area_update_request->tmsistatus,
             TRACKING_AREA_UPDATE_REQUEST_TMSI_STATUS_IEI, buffer + encoded,
             len - encoded)) < 0) {
      // Return in case of error
      return encode_result;
    } else {
      encoded += encode_result;
    }
  }

  if ((tracking_area_update_request->presencemask &
       TRACKING_AREA_UPDATE_REQUEST_MOBILE_STATION_CLASSMARK_2_PRESENT) ==
      TRACKING_AREA_UPDATE_REQUEST_MOBILE_STATION_CLASSMARK_2_PRESENT) {
    if ((encode_result = encode_mobile_station_classmark_2_ie(
             &tracking_area_update_request->mobilestationclassmark2,
             TRACKING_AREA_UPDATE_REQUEST_MOBILE_STATION_CLASSMARK_2_IEI,
             buffer + encoded, len - encoded)) < 0) {
      // Return in case of error
      return encode_result;
    } else {
      encoded += encode_result;
    }
  }

  if ((tracking_area_update_request->presencemask &
       TRACKING_AREA_UPDATE_REQUEST_MOBILE_STATION_CLASSMARK_3_PRESENT) ==
      TRACKING_AREA_UPDATE_REQUEST_MOBILE_STATION_CLASSMARK_3_PRESENT) {
    if ((encode_result = encode_mobile_station_classmark_3_ie(
             &tracking_area_update_request->mobilestationclassmark3,
             TRACKING_AREA_UPDATE_REQUEST_MOBILE_STATION_CLASSMARK_3_IEI,
             buffer + encoded, len - encoded)) < 0) {
      // Return in case of error
      return encode_result;
    } else {
      encoded += encode_result;
    }
  }

  if ((tracking_area_update_request->presencemask &
       TRACKING_AREA_UPDATE_REQUEST_SUPPORTED_CODECS_PRESENT) ==
      TRACKING_AREA_UPDATE_REQUEST_SUPPORTED_CODECS_PRESENT) {
    if ((encode_result = encode_supported_codec_list_ie(
             &tracking_area_update_request->supportedcodecs,
             TRACKING_AREA_UPDATE_REQUEST_SUPPORTED_CODECS_IEI,
             buffer + encoded, len - encoded)) < 0) {
      // Return in case of error
      return encode_result;
    } else {
      encoded += encode_result;
    }
  }

  if ((tracking_area_update_request->presencemask &
       TRACKING_AREA_UPDATE_REQUEST_ADDITIONAL_UPDATE_TYPE_PRESENT) ==
      TRACKING_AREA_UPDATE_REQUEST_ADDITIONAL_UPDATE_TYPE_PRESENT) {
    if ((encode_result = encode_additional_update_type(
             &tracking_area_update_request->additionalupdatetype,
             TRACKING_AREA_UPDATE_REQUEST_ADDITIONAL_UPDATE_TYPE_IEI,
             buffer + encoded, len - encoded)) < 0) {
      // Return in case of error
      return encode_result;
    } else {
      encoded += encode_result;
    }
  }

  if ((tracking_area_update_request->presencemask &
       TRACKING_AREA_UPDATE_REQUEST_OLD_GUTI_TYPE_PRESENT) ==
      TRACKING_AREA_UPDATE_REQUEST_OLD_GUTI_TYPE_PRESENT) {
    if ((encode_result = encode_guti_type(
             &tracking_area_update_request->oldgutitype,
             TRACKING_AREA_UPDATE_REQUEST_OLD_GUTI_TYPE_IEI, buffer + encoded,
             len - encoded)) < 0) {
      // Return in case of error
      return encode_result;
    } else {
      encoded += encode_result;
    }
  }

  return encoded;
}
