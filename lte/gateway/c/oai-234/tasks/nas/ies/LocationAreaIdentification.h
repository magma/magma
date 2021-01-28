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

#ifndef LOCATION_AREA_IDENTIFICATION_H_
#define LOCATION_AREA_IDENTIFICATION_H_
#include <stdint.h>

#define LOCATION_AREA_IDENTIFICATION_MINIMUM_LENGTH 6
#define LOCATION_AREA_IDENTIFICATION_MAXIMUM_LENGTH 6

typedef struct LocationAreaIdentification_tag {
  uint8_t mccdigit2 : 4;
  uint8_t mccdigit1 : 4;
  uint8_t mncdigit3 : 4;
  uint8_t mccdigit3 : 4;
  uint8_t mncdigit2 : 4;
  uint8_t mncdigit1 : 4;
  uint16_t lac;
} LocationAreaIdentification;

int encode_location_area_identification(
    LocationAreaIdentification* locationareaidentification, uint8_t iei,
    uint8_t* buffer, uint32_t len);

void dump_location_area_identification_xml(
    LocationAreaIdentification* locationareaidentification, uint8_t iei);

int decode_location_area_identification(
    LocationAreaIdentification* locationareaidentification, uint8_t iei,
    uint8_t* buffer, uint32_t len);

#endif /* LOCATION AREA IDENTIFICATION_H_ */
