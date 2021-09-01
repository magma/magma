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

#ifndef MOBILE_STATION_CLASSMARK_2_H_
#define MOBILE_STATION_CLASSMARK_2_H_
#include <stdint.h>

#define MOBILE_STATION_CLASSMARK_2_MINIMUM_LENGTH 5
#define MOBILE_STATION_CLASSMARK_2_MAXIMUM_LENGTH 5

typedef struct MobileStationClassmark2_tag {
  uint8_t revisionlevel : 2;
  uint8_t esind : 1;
  uint8_t a51 : 1;
  uint8_t rfpowercapability : 3;
  uint8_t pscapability : 1;
  uint8_t ssscreenindicator : 2;
  uint8_t smcapability : 1;
  uint8_t vbs : 1;
  uint8_t vgcs : 1;
  uint8_t fc : 1;
  uint8_t cm3 : 1;
  uint8_t lcsvacap : 1;
  uint8_t ucs2 : 1;
  uint8_t solsa : 1;
  uint8_t cmsp : 1;
  uint8_t a53 : 1;
  uint8_t a52 : 1;
} MobileStationClassmark2;

int encode_mobile_station_classmark_2(
    MobileStationClassmark2* mobilestationclassmark2, uint8_t iei,
    uint8_t* buffer, uint32_t len);

int decode_mobile_station_classmark_2(
    MobileStationClassmark2* mobilestationclassmark2, uint8_t iei,
    uint8_t* buffer, uint32_t len);

void dump_mobile_station_classmark_2_xml(
    MobileStationClassmark2* mobilestationclassmark2, uint8_t iei);

#endif /* MOBILE STATION CLASSMARK 2_H_ */
