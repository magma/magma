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

#ifndef TIME_ZONE_AND_TIME_H_
#define TIME_ZONE_AND_TIME_H_
#include <stdint.h>

#define TIME_ZONE_AND_TIME_MINIMUM_LENGTH 8
#define TIME_ZONE_AND_TIME_MAXIMUM_LENGTH 8

typedef struct TimeZoneAndTime_tag {
  uint8_t year;
  uint8_t month;
  uint8_t day;
  uint8_t hour;
  uint8_t minute;
  uint8_t second;
  uint8_t timezone;
} TimeZoneAndTime;

int encode_time_zone_and_time(
    TimeZoneAndTime* timezoneandtime, uint8_t iei, uint8_t* buffer,
    uint32_t len);

void dump_time_zone_and_time_xml(TimeZoneAndTime* timezoneandtime, uint8_t iei);

int decode_time_zone_and_time(
    TimeZoneAndTime* timezoneandtime, uint8_t iei, uint8_t* buffer,
    uint32_t len);

#endif /* TIME ZONE AND TIME_H_ */
