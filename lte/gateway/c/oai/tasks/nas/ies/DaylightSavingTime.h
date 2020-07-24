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

#ifndef DAYLIGHT_SAVING_TIME_H_
#define DAYLIGHT_SAVING_TIME_H_
#include <stdint.h>

#define DAYLIGHT_SAVING_TIME_MINIMUM_LENGTH 3
#define DAYLIGHT_SAVING_TIME_MAXIMUM_LENGTH 3

typedef uint8_t DaylightSavingTime;

int encode_daylight_saving_time(
    DaylightSavingTime* daylightsavingtime, uint8_t iei, uint8_t* buffer,
    uint32_t len);

int decode_daylight_saving_time(
    DaylightSavingTime* daylightsavingtime, uint8_t iei, uint8_t* buffer,
    uint32_t len);

void dump_daylight_saving_time_xml(
    DaylightSavingTime* daylightsavingtime, uint8_t iei);

#endif /* DAYLIGHT SAVING TIME_H_ */
