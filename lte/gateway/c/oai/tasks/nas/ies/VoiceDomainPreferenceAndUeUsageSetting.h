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

#ifndef VOICE_DOMAIN_PREFERENCE_AND_UE_USAGE_SETTING_H_
#define VOICE_DOMAIN_PREFERENCE_AND_UE_USAGE_SETTING_H_
#include <stdint.h>

#define VOICE_DOMAIN_PREFERENCE_AND_UE_USAGE_SETTING_MINIMUM_LENGTH 1
#define VOICE_DOMAIN_PREFERENCE_AND_UE_USAGE_SETTING_MAXIMUM_LENGTH 1

typedef struct VoiceDomainPreferenceAndUeUsageSetting_tag {
  uint8_t spare : 5;
#define UE_USAGE_SETTING_VOICE_CENTRIC 0b0
#define UE_USAGE_SETTING_DATA_CENTRIC 0b1
  uint8_t ue_usage_setting : 1;
#define VOICE_DOMAIN_PREFERENCE_CS_VOICE_ONLY 0b00
#define VOICE_DOMAIN_PREFERENCE_IMS_PS_VOICE_ONLY 0b01
#define VOICE_DOMAIN_PREFERENCE_CS_VOICE_PREFERRED_IMS_PS_VOICE_AS_SECONDARY   \
  0b10
#define VOICE_DOMAIN_PREFERENCE_IMS_PS_VOICE_PREFERRED_CS_VOICE_AS_SECONDARY   \
  0b11
  uint8_t voice_domain_for_eutran : 2;
} VoiceDomainPreferenceAndUeUsageSetting;

int encode_voice_domain_preference_and_ue_usage_setting(
    VoiceDomainPreferenceAndUeUsageSetting*
        voicedomainpreferenceandueusagesetting,
    uint8_t iei, uint8_t* buffer, uint32_t len);

int decode_voice_domain_preference_and_ue_usage_setting(
    VoiceDomainPreferenceAndUeUsageSetting*
        voicedomainpreferenceandueusagesetting,
    uint8_t iei, uint8_t* buffer, uint32_t len);

void dump_voice_domain_preference_and_ue_usage_setting_xml(
    VoiceDomainPreferenceAndUeUsageSetting*
        voicedomainpreferenceandueusagesetting,
    uint8_t iei);

#endif /* VOICE DOMAIN PREFERENCE AND UE USAGE SETTING_H_ */
