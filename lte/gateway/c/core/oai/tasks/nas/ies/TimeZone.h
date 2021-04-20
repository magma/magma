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

#ifndef TIME_ZONE_H_
#define TIME_ZONE_H_
#include <stdint.h>

#define TIME_ZONE_MINIMUM_LENGTH 2
#define TIME_ZONE_MAXIMUM_LENGTH 2

typedef uint8_t TimeZone;

int encode_time_zone(
    TimeZone* timezone, uint8_t iei, uint8_t* buffer, uint32_t len);

void dump_time_zone_xml(TimeZone* timezone, uint8_t iei);

int decode_time_zone(
    TimeZone* timezone, uint8_t iei, uint8_t* buffer, uint32_t len);

#endif /* TIME ZONE_H_ */
