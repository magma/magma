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

#ifndef MOBILE_STATION_CLASSMARK_3_H_
#define MOBILE_STATION_CLASSMARK_3_H_
#include <stdint.h>

#define MOBILE_STATION_CLASSMARK_3_MINIMUM_LENGTH 1
#define MOBILE_STATION_CLASSMARK_3_MAXIMUM_LENGTH 1

typedef struct {
  uint8_t field;
} MobileStationClassmark3;

int encode_mobile_station_classmark_3(
    MobileStationClassmark3* mobilestationclassmark3, uint8_t iei,
    uint8_t* buffer, uint32_t len);

void dump_mobile_station_classmark_3_xml(
    MobileStationClassmark3* mobilestationclassmark3, uint8_t iei);

int decode_mobile_station_classmark_3(
    MobileStationClassmark3* mobilestationclassmark3, uint8_t iei,
    uint8_t* buffer, uint32_t len);

#endif /* MOBILE STATION CLASSMARK 3_H_ */
