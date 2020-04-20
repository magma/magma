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
#include "TimeZone.h"

int decode_time_zone(
    TimeZone* timezone, uint8_t iei, uint8_t* buffer, uint32_t len) {
  int decoded = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  *timezone = *(buffer + decoded);
  decoded++;
#if NAS_DEBUG
  dump_time_zone_xml(timezone, iei);
#endif
  return decoded;
}

int encode_time_zone(
    TimeZone* timezone, uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(buffer, TIME_ZONE_MINIMUM_LENGTH, len);
#if NAS_DEBUG
  dump_time_zone_xml(timezone, iei);
#endif

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  *(buffer + encoded) = *timezone;
  encoded++;
  return encoded;
}

void dump_time_zone_xml(TimeZone* timezone, uint8_t iei) {
  OAILOG_DEBUG(LOG_NAS, "<Time Zone>\n");

  if (iei > 0)
    /*
     * Don't display IEI if = 0
     */
    OAILOG_DEBUG(LOG_NAS, "    <IEI>0x%X</IEI>\n", iei);

  OAILOG_DEBUG(LOG_NAS, "    <Time zone>%u</Time zone>\n", *timezone);
  OAILOG_DEBUG(LOG_NAS, "</Time Zone>\n");
}
