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

#ifndef EMERGENCY_NUMBER_LIST_H_
#define EMERGENCY_NUMBER_LIST_H_
#include <stdint.h>

#define EMERGENCY_NUMBER_LIST_MINIMUM_LENGTH 5
#define EMERGENCY_NUMBER_LIST_MAXIMUM_LENGTH 50

typedef struct EmergencyNumberList_tag {
  uint8_t lengthofemergency;
  uint8_t emergencyservicecategoryvalue : 5;
} EmergencyNumberList;

int encode_emergency_number_list(
    EmergencyNumberList* emergencynumberlist, uint8_t iei, uint8_t* buffer,
    uint32_t len);

int decode_emergency_number_list(
    EmergencyNumberList* emergencynumberlist, uint8_t iei, uint8_t* buffer,
    uint32_t len);

void dump_emergency_number_list_xml(
    EmergencyNumberList* emergencynumberlist, uint8_t iei);

#endif /* EMERGENCY NUMBER LIST_H_ */
