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

#ifndef UE_RADIO_CAPABILITY_INFORMATION_UPDATE_NEEDED_SEEN
#define UE_RADIO_CAPABILITY_INFORMATION_UPDATE_NEEDED_SEEN

#include <stdint.h>

#define UE_RADIO_CAPABILITY_INFORMATION_UPDATE_NEEDED_MINIMUM_LENGTH 1
#define UE_RADIO_CAPABILITY_INFORMATION_UPDATE_NEEDED_MAXIMUM_LENGTH 1

typedef uint8_t ue_radio_capability_information_update_needed_t;

int encode_ue_radio_capability_information_update_needed(
    ue_radio_capability_information_update_needed_t*
        ueradiocapabilityinformationupdateneeded,
    uint8_t iei, uint8_t* buffer, uint32_t len);

uint8_t encode_u8_ue_radio_capability_information_update_needed(
    ue_radio_capability_information_update_needed_t*
        ueradiocapabilityinformationupdateneeded);

int decode_ue_radio_capability_information_update_needed(
    ue_radio_capability_information_update_needed_t*
        ueradiocapabilityinformationupdateneeded,
    uint8_t iei, uint8_t* buffer, uint32_t len);

int decode_u8_ue_radio_capability_information_update_needed(
    ue_radio_capability_information_update_needed_t*
        ueradiocapabilityinformationupdateneeded,
    uint8_t iei, uint8_t value, uint32_t len);

#endif /* UE RADIO CAPABILITY INFORMATION UPDATE NEEDED_SEEN */
