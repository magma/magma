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
#include "ImeisvRequest.h"

int decode_imeisv_request(
    ImeisvRequest* imeisvrequest, uint8_t iei, uint8_t* buffer, uint32_t len) {
  int decoded = 0;

  CHECK_PDU_POINTER_AND_LENGTH_DECODER(
      buffer, IMEISV_REQUEST_MINIMUM_LENGTH, len);

  if (iei > 0) {
    CHECK_IEI_DECODER((*buffer & 0xf0), iei);
  }

  *imeisvrequest = *buffer & 0x7;
  decoded++;
#if NAS_DEBUG
  dump_imeisv_request_xml(imeisvrequest, iei);
#endif
  return decoded;
}

int decode_u8_imeisv_request(
    ImeisvRequest* imeisvrequest, uint8_t iei, uint8_t value, uint32_t len) {
  int decoded     = 0;
  uint8_t* buffer = &value;

  *imeisvrequest = *buffer & 0x7;
  decoded++;
#if NAS_DEBUG
  dump_imeisv_request_xml(imeisvrequest, iei);
#endif
  return decoded;
}

int encode_imeisv_request(
    ImeisvRequest* imeisvrequest, uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint8_t encoded = 0;

  /*
   * Checking length and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, IMEISV_REQUEST_MINIMUM_LENGTH, len);
#if NAS_DEBUG
  dump_imeisv_request_xml(imeisvrequest, iei);
#endif
  *(buffer + encoded) = 0x00 | (iei & 0xf0) | (*imeisvrequest & 0x7);
  encoded++;
  return encoded;
}

uint8_t encode_u8_imeisv_request(ImeisvRequest* imeisvrequest) {
  uint8_t bufferReturn;
  uint8_t* buffer = &bufferReturn;
  uint8_t encoded = 0;
  uint8_t iei     = 0;

  dump_imeisv_request_xml(imeisvrequest, 0);
  *(buffer + encoded) = 0x00 | (iei & 0xf0) | (*imeisvrequest & 0x7);
  encoded++;
  return bufferReturn;
}

void dump_imeisv_request_xml(ImeisvRequest* imeisvrequest, uint8_t iei) {
  OAILOG_DEBUG(LOG_NAS, "<Imeisv Request>\n");

  if (iei > 0)
    /*
     * Don't display IEI if = 0
     */
    OAILOG_DEBUG(LOG_NAS, "    <IEI>0x%X</IEI>\n", iei);

  OAILOG_DEBUG(
      LOG_NAS, "    <IMEISV request value>%u</IMEISV request value>\n",
      *imeisvrequest);
  OAILOG_DEBUG(LOG_NAS, "</Imeisv Request>\n");
}
