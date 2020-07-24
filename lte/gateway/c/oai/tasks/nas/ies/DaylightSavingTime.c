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
#include "DaylightSavingTime.h"

int decode_daylight_saving_time(
    DaylightSavingTime* daylightsavingtime, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded   = 0;
  uint8_t ielen = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  *daylightsavingtime = *buffer & 0x3;
  decoded++;
  return decoded;
}

int encode_daylight_saving_time(
    DaylightSavingTime* daylightsavingtime, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, DAYLIGHT_SAVING_TIME_MINIMUM_LENGTH, len);

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;
  *(buffer + encoded) = 0x00 | (*daylightsavingtime & 0x3);
  encoded++;
  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);
  return encoded;
}

void dump_daylight_saving_time_xml(
    DaylightSavingTime* daylightsavingtime, uint8_t iei) {
  OAILOG_DEBUG(LOG_NAS, "<Daylight Saving Time>\n");

  if (iei > 0)
    /*
     * Don't display IEI if = 0
     */
    OAILOG_DEBUG(LOG_NAS, "    <IEI>0x%X</IEI>\n", iei);

  OAILOG_DEBUG(LOG_NAS, "    <Value>%u</Value>\n", *daylightsavingtime);
  OAILOG_DEBUG(LOG_NAS, "</Daylight Saving Time>\n");
}
