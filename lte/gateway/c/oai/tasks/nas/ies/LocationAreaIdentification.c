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
#include "LocationAreaIdentification.h"

int decode_location_area_identification(
    LocationAreaIdentification* locationareaidentification, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  int decoded = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  locationareaidentification->mccdigit2 = (*(buffer + decoded) >> 4) & 0xf;
  locationareaidentification->mccdigit1 = *(buffer + decoded) & 0xf;
  decoded++;
  locationareaidentification->mncdigit3 = (*(buffer + decoded) >> 4) & 0xf;
  locationareaidentification->mccdigit3 = *(buffer + decoded) & 0xf;
  decoded++;
  locationareaidentification->mncdigit2 = (*(buffer + decoded) >> 4) & 0xf;
  locationareaidentification->mncdigit1 = *(buffer + decoded) & 0xf;
  decoded++;
  // IES_DECODE_U16(locationareaidentification->lac, *(buffer + decoded));
  IES_DECODE_U16(buffer, decoded, locationareaidentification->lac);
#if NAS_DEBUG
  dump_location_area_identification_xml(locationareaidentification, iei);
#endif
  return decoded;
}

int encode_location_area_identification(
    LocationAreaIdentification* locationareaidentification, uint8_t iei,
    uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, LOCATION_AREA_IDENTIFICATION_MINIMUM_LENGTH, len);
#if NAS_DEBUG
  dump_location_area_identification_xml(locationareaidentification, iei);
#endif

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  *(buffer + encoded) = 0x00 |
                        ((locationareaidentification->mccdigit2 & 0xf) << 4) |
                        (locationareaidentification->mccdigit1 & 0xf);
  encoded++;
  *(buffer + encoded) = 0x00 |
                        ((locationareaidentification->mncdigit3 & 0xf) << 4) |
                        (locationareaidentification->mccdigit3 & 0xf);
  encoded++;
  *(buffer + encoded) = 0x00 |
                        ((locationareaidentification->mncdigit2 & 0xf) << 4) |
                        (locationareaidentification->mncdigit1 & 0xf);
  encoded++;
  IES_ENCODE_U16(buffer, encoded, locationareaidentification->lac);
  return encoded;
}

void dump_location_area_identification_xml(
    LocationAreaIdentification* locationareaidentification, uint8_t iei) {
  OAILOG_DEBUG(LOG_NAS, "<Location Area Identification>\n");

  if (iei > 0)
    /*
     * Don't display IEI if = 0
     */
    OAILOG_DEBUG(LOG_NAS, "    <IEI>0x%X</IEI>\n", iei);

  OAILOG_DEBUG(
      LOG_NAS, "    <MCC digit 2>%u</MCC digit 2>\n",
      locationareaidentification->mccdigit2);
  OAILOG_DEBUG(
      LOG_NAS, "    <MCC digit 1>%u</MCC digit 1>\n",
      locationareaidentification->mccdigit1);
  OAILOG_DEBUG(
      LOG_NAS, "    <MNC digit 3>%u</MNC digit 3>\n",
      locationareaidentification->mncdigit3);
  OAILOG_DEBUG(
      LOG_NAS, "    <MCC digit 3>%u</MCC digit 3>\n",
      locationareaidentification->mccdigit3);
  OAILOG_DEBUG(
      LOG_NAS, "    <MNC digit 2>%u</MNC digit 2>\n",
      locationareaidentification->mncdigit2);
  OAILOG_DEBUG(
      LOG_NAS, "    <MNC digit 1>%u</MNC digit 1>\n",
      locationareaidentification->mncdigit1);
  OAILOG_DEBUG(LOG_NAS, "    <LAC>%u</LAC>\n", locationareaidentification->lac);
  OAILOG_DEBUG(LOG_NAS, "</Location Area Identification>\n");
}
