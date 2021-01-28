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
#include "EsmInformationRequest.h"

int decode_esm_information_request(
    esm_information_request_msg* esm_information_request, uint8_t* buffer,
    uint32_t len) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  uint32_t decoded = 0;

  // Check if we got a NULL pointer and if buffer length is >= minimum length
  // expected for the message.
  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, ESM_INFORMATION_REQUEST_MINIMUM_LENGTH, len);
  /*
   * Decoding mandatory fields
   */
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, decoded);
}

int encode_esm_information_request(
    esm_information_request_msg* esm_information_request, uint8_t* buffer,
    uint32_t len) {
  OAILOG_FUNC_IN(LOG_NAS_ESM);
  int encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, ESM_INFORMATION_REQUEST_MINIMUM_LENGTH, len);
  OAILOG_FUNC_RETURN(LOG_NAS_ESM, encoded);
}
