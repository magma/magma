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
#include "IdentityRequest.h"

int decode_identity_request(
    identity_request_msg* identity_request, uint8_t* buffer, uint32_t len) {
  uint32_t decoded   = 0;
  int decoded_result = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, IDENTITY_REQUEST_MINIMUM_LENGTH, len);

  /*
   * Decoding mandatory fields
   */
  if ((decoded_result = decode_identity_type_2_ie(
           &identity_request->identitytype, false, buffer + decoded,
           len - decoded)) < 0)
    return decoded_result;

  decoded++;
  return decoded;
}

int encode_identity_request(
    identity_request_msg* identity_request, uint8_t* buffer, uint32_t len) {
  int encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, IDENTITY_REQUEST_MINIMUM_LENGTH, len);
  encoded += encode_identity_type_2_ie(
      &identity_request->identitytype, false, buffer,
      IDENTITY_TYPE_2_IE_MIN_LENGTH);
  return encoded;
}
