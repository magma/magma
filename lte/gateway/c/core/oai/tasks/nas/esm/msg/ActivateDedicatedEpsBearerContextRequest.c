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
#include <stdbool.h>

#include "TLVEncoder.h"
#include "TLVDecoder.h"
#include "ActivateDedicatedEpsBearerContextRequest.h"
#include "common_defs.h"

int decode_activate_dedicated_eps_bearer_context_request(
    activate_dedicated_eps_bearer_context_request_msg*
        activate_dedicated_eps_bearer_context_request,
    uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_MINIMUM_LENGTH,
      len);

  /*
   * Decoding mandatory fields
   */
  if ((decoded_result = decode_u8_linked_eps_bearer_identity(
           &activate_dedicated_eps_bearer_context_request
                ->linkedepsbeareridentity,
           0, *(buffer + decoded), len - decoded)) < 0)
    return decoded_result;

  decoded++;

  if ((decoded_result = decode_eps_quality_of_service(
           &activate_dedicated_eps_bearer_context_request->epsqos, 0,
           buffer + decoded, len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  if ((decoded_result = decode_traffic_flow_template_ie(
           &activate_dedicated_eps_bearer_context_request->tft, 0,
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
    if (ieiDecoded >= 0x80) ieiDecoded = ieiDecoded & 0xf0;

    switch (ieiDecoded) {
      case ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_TRANSACTION_IDENTIFIER_IEI:
        if ((decoded_result = decode_linked_ti_ie(
                 &activate_dedicated_eps_bearer_context_request
                      ->transactionidentifier,
                 true, buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        activate_dedicated_eps_bearer_context_request->presencemask |=
            ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_TRANSACTION_IDENTIFIER_PRESENT;
        break;

      case ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_NEGOTIATED_QOS_IEI:
        if ((decoded_result = decode_quality_of_service_ie(
                 &activate_dedicated_eps_bearer_context_request->negotiatedqos,
                 true, buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        activate_dedicated_eps_bearer_context_request->presencemask |=
            ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_NEGOTIATED_QOS_PRESENT;
        break;

      case ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_NEGOTIATED_LLC_SAPI_IEI:
        if ((decoded_result = decode_llc_service_access_point_identifier_ie(
                 &activate_dedicated_eps_bearer_context_request
                      ->negotiatedllcsapi,
                 true, buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        activate_dedicated_eps_bearer_context_request->presencemask |=
            ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_NEGOTIATED_LLC_SAPI_PRESENT;
        break;

      case ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_RADIO_PRIORITY_IEI:
        if ((decoded_result = decode_radio_priority(
                 &activate_dedicated_eps_bearer_context_request->radiopriority,
                 ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_RADIO_PRIORITY_IEI,
                 buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        activate_dedicated_eps_bearer_context_request->presencemask |=
            ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_RADIO_PRIORITY_PRESENT;
        break;

      case ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_PACKET_FLOW_IDENTIFIER_IEI:
        if ((decoded_result = decode_packet_flow_identifier_ie(
                 &activate_dedicated_eps_bearer_context_request
                      ->packetflowidentifier,
                 true, buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        activate_dedicated_eps_bearer_context_request->presencemask |=
            ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_PACKET_FLOW_IDENTIFIER_PRESENT;
        break;

      case ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_IEI:
        if ((decoded_result = decode_protocol_configuration_options_ie(
                 &activate_dedicated_eps_bearer_context_request
                      ->protocolconfigurationoptions,
                 true, buffer + decoded, len - decoded)) <= 0)
          return decoded_result;

        decoded += decoded_result;
        /*
         * Set corresponding mask to 1 in presencemask
         */
        activate_dedicated_eps_bearer_context_request->presencemask |=
            ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT;
        break;

      default:
        errorCodeDecoder = TLV_UNEXPECTED_IEI;
        return TLV_UNEXPECTED_IEI;
    }
  }

  return decoded;
}

int encode_activate_dedicated_eps_bearer_context_request(
    activate_dedicated_eps_bearer_context_request_msg*
        activate_dedicated_eps_bearer_context_request,
    uint8_t* buffer, uint32_t len) {
  int encoded       = 0;
  int encode_result = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_MINIMUM_LENGTH,
      len);
  *(buffer + encoded) = encode_u8_linked_eps_bearer_identity(
                            &activate_dedicated_eps_bearer_context_request
                                 ->linkedepsbeareridentity) &
                        0x0f;
  encoded++;

  if ((encode_result = encode_eps_quality_of_service(
           &activate_dedicated_eps_bearer_context_request->epsqos, 0,
           buffer + encoded,
           len - encoded)) < 0)  // Return in case of error
    return encode_result;
  else
    encoded += encode_result;

  if ((encode_result = encode_traffic_flow_template_ie(
           &activate_dedicated_eps_bearer_context_request->tft,
           TFT_ENCODE_IEI_FALSE, buffer + encoded,
           len - encoded)) < 0)  // Return in case of error
    return encode_result;
  else
    encoded += encode_result;

  if ((activate_dedicated_eps_bearer_context_request->presencemask &
       ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_TRANSACTION_IDENTIFIER_PRESENT) ==
      ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_TRANSACTION_IDENTIFIER_PRESENT) {
    if ((encode_result = encode_linked_ti_ie(
             &activate_dedicated_eps_bearer_context_request
                  ->transactionidentifier,
             ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_TRANSACTION_IDENTIFIER_IEI,
             buffer + encoded, len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((activate_dedicated_eps_bearer_context_request->presencemask &
       ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_NEGOTIATED_QOS_PRESENT) ==
      ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_NEGOTIATED_QOS_PRESENT) {
    if ((encode_result = encode_quality_of_service_ie(
             &activate_dedicated_eps_bearer_context_request->negotiatedqos,
             true, buffer + encoded, len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((activate_dedicated_eps_bearer_context_request->presencemask &
       ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_NEGOTIATED_LLC_SAPI_PRESENT) ==
      ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_NEGOTIATED_LLC_SAPI_PRESENT) {
    if ((encode_result = encode_llc_service_access_point_identifier_ie(
             &activate_dedicated_eps_bearer_context_request->negotiatedllcsapi,
             true, buffer + encoded, len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((activate_dedicated_eps_bearer_context_request->presencemask &
       ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_RADIO_PRIORITY_PRESENT) ==
      ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_RADIO_PRIORITY_PRESENT) {
    if ((encode_result = encode_radio_priority(
             &activate_dedicated_eps_bearer_context_request->radiopriority,
             ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_RADIO_PRIORITY_IEI,
             buffer + encoded, len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((activate_dedicated_eps_bearer_context_request->presencemask &
       ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_PACKET_FLOW_IDENTIFIER_PRESENT) ==
      ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_PACKET_FLOW_IDENTIFIER_PRESENT) {
    if ((encode_result = encode_packet_flow_identifier_ie(
             &activate_dedicated_eps_bearer_context_request
                  ->packetflowidentifier,
             true, buffer + encoded, len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  if ((activate_dedicated_eps_bearer_context_request->presencemask &
       ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT) ==
      ACTIVATE_DEDICATED_EPS_BEARER_CONTEXT_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT) {
    if ((encode_result = encode_protocol_configuration_options_ie(
             &activate_dedicated_eps_bearer_context_request
                  ->protocolconfigurationoptions,
             true, buffer + encoded, len - encoded)) < 0)
      // Return in case of error
      return encode_result;
    else
      encoded += encode_result;
  }

  return encoded;
}
