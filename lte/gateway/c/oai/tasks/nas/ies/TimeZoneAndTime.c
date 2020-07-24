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
#include "TimeZoneAndTime.h"

int decode_time_zone_and_time(
    TimeZoneAndTime* timezoneandtime, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  int decoded = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  timezoneandtime->year = *(buffer + decoded);
  decoded++;
  timezoneandtime->month = *(buffer + decoded);
  decoded++;
  timezoneandtime->day = *(buffer + decoded);
  decoded++;
  timezoneandtime->hour = *(buffer + decoded);
  decoded++;
  timezoneandtime->minute = *(buffer + decoded);
  decoded++;
  timezoneandtime->second = *(buffer + decoded);
  decoded++;
  timezoneandtime->timezone = *(buffer + decoded);
  decoded++;
#if NAS_DEBUG
  dump_time_zone_and_time_xml(timezoneandtime, iei);
#endif
  return decoded;
}

int encode_time_zone_and_time(
    TimeZoneAndTime* timezoneandtime, uint8_t iei, uint8_t* buffer,
    uint32_t len) {
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, TIME_ZONE_AND_TIME_MINIMUM_LENGTH, len);
#if NAS_DEBUG
  dump_time_zone_and_time_xml(timezoneandtime, iei);
#endif

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  *(buffer + encoded) = timezoneandtime->year;
  encoded++;
  *(buffer + encoded) = timezoneandtime->month;
  encoded++;
  *(buffer + encoded) = timezoneandtime->day;
  encoded++;
  *(buffer + encoded) = timezoneandtime->hour;
  encoded++;
  *(buffer + encoded) = timezoneandtime->minute;
  encoded++;
  *(buffer + encoded) = timezoneandtime->second;
  encoded++;
  *(buffer + encoded) = timezoneandtime->timezone;
  encoded++;
  return encoded;
}

void dump_time_zone_and_time_xml(
    TimeZoneAndTime* timezoneandtime, uint8_t iei) {
  OAILOG_DEBUG(LOG_NAS, "<Time Zone And Time>\n");

  if (iei > 0)
    /*
     * Don't display IEI if = 0
     */
    OAILOG_DEBUG(LOG_NAS, "    <IEI>0x%X</IEI>\n", iei);

  OAILOG_DEBUG(LOG_NAS, "    <Year>%u</Year>\n", timezoneandtime->year);
  OAILOG_DEBUG(LOG_NAS, "    <Month>%u</Month>\n", timezoneandtime->month);
  OAILOG_DEBUG(LOG_NAS, "    <Day>%u</Day>\n", timezoneandtime->day);
  OAILOG_DEBUG(LOG_NAS, "    <Hour>%u</Hour>\n", timezoneandtime->hour);
  OAILOG_DEBUG(LOG_NAS, "    <Minute>%u</Minute>\n", timezoneandtime->minute);
  OAILOG_DEBUG(LOG_NAS, "    <Second>%u</Second>\n", timezoneandtime->second);
  OAILOG_DEBUG(
      LOG_NAS, "    <Time Zone>%u</Time Zone>\n", timezoneandtime->timezone);
  OAILOG_DEBUG(LOG_NAS, "</Time Zone And Time>\n");
}
