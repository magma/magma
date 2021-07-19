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

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#include "TLVEncoder.h"
#include "TLVDecoder.h"
#include "AuthenticationParameterRand.h"

int decode_authentication_parameter_rand(
    AuthenticationParameterRand* authenticationparameterrand, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  int decoded   = 0;
  uint8_t ielen = 16;
  int decode_result;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  if ((decode_result = decode_bstring(
           authenticationparameterrand, ielen, buffer + decoded,
           len - decoded)) < 0)
    return decode_result;
  else
    decoded += decode_result;

  return decoded;
}

int encode_authentication_parameter_rand(
    AuthenticationParameterRand authenticationparameterrand, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  uint32_t encode_result;
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, AUTHENTICATION_PARAMETER_RAND_MINIMUM_LENGTH, len);

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  if ((encode_result = encode_bstring(
           authenticationparameterrand, buffer + encoded, len - encoded)) < 0) {
    return encode_result;
  } else {
    encoded += encode_result;
  }

  return encoded;
}

void dump_authentication_parameter_rand_xml(
    AuthenticationParameterRand authenticationparameterrand, uint8_t iei) {
  OAILOG_DEBUG(LOG_NAS, "<Authentication Parameter Rand>\n");

  if (iei > 0)
    /*
     * Don't display IEI if = 0
     */
    OAILOG_DEBUG(LOG_NAS, "    <IEI>0x%X</IEI>\n", iei);

  bstring b = dump_bstring_xml(authenticationparameterrand);
  OAILOG_DEBUG(LOG_NAS, "%s", bdata(b));
  bdestroy(b);
  OAILOG_DEBUG(LOG_NAS, "</Authentication Parameter Rand>\n");
}
