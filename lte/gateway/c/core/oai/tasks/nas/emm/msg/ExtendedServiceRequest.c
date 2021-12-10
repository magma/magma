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
#include "lte/gateway/c/core/oai/tasks/nas/emm/msg/ExtendedServiceRequest.h"

int decode_extended_service_request(
    extended_service_request_msg* extended_service_request, uint8_t* buffer,
    uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, EXTENDED_SERVICE_REQUEST_MINIMUM_LENGTH, len);

  /*
   * Decoding mandatory fields
   */
  if ((decoded_result = decode_u8_service_type(
           &extended_service_request->servicetype, 0,
           *(buffer + decoded) & 0x0f, len - decoded)) < 0)
    return decoded_result;

  if ((decoded_result = decode_u8_nas_key_set_identifier(
           &extended_service_request->naskeysetidentifier, 0,
           *(buffer + decoded) >> 4, len - decoded)) < 0)
    return decoded_result;

  decoded++;

  if ((decoded_result = decode_mobile_identity_ie(
           &extended_service_request->mtmsi, 0, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;
  else
    decoded += decoded_result;

  /* CSFB Response is conditional IE,
   * The UE shall include this IE only if the Service type information element
   * indicates "mobile terminating CS fallback or 1xCS fallback"
   */
  if (extended_service_request->servicetype == MT_CS_FB) {
    if ((decoded_result = decode_csfb_response(
             &extended_service_request->csfbresponse, 0, (buffer + decoded),
             len - decoded)) < 0) {
      return decoded_result;
    } else {
      decoded += decoded_result;
      extended_service_request->presencemask |= EMM_CSFB_RSP_PRESENT;
    }
  }
  return decoded;
}

int encode_extended_service_request(
    extended_service_request_msg* extended_service_request, uint8_t* buffer,
    uint32_t len) {
  int encoded       = 0;
  int encode_result = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, EXTENDED_SERVICE_REQUEST_MINIMUM_LENGTH, len);
  *(buffer + encoded) =
      ((encode_u8_nas_key_set_identifier(
            &extended_service_request->naskeysetidentifier) &
        0x0f)
       << 4) |
      (encode_u8_service_type(&extended_service_request->servicetype) & 0x0f);
  encoded++;

  if ((encode_result = encode_mobile_identity_ie(
           &extended_service_request->mtmsi, 0, buffer + encoded,
           len - encoded)) < 0)  // Return in case of error
    return encode_result;
  else
    encoded += encode_result;

  if ((encode_result = encode_csfb_response(
           &extended_service_request->csfbresponse, 0, buffer + encoded,
           len - encoded)) < 0)  // Return in case of error
    return encode_result;
  else
    encoded += encode_result;

  return encoded;
}
