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
#include <string.h>

#include "TLVEncoder.h"
#include "TLVDecoder.h"
#include "VoiceDomainPreferenceAndUeUsageSetting.h"

int decode_voice_domain_preference_and_ue_usage_setting(
    VoiceDomainPreferenceAndUeUsageSetting*
        voicedomainpreferenceandueusagesetting,
    uint8_t iei, uint8_t* buffer, uint32_t len) {
  int decoded   = 0;
  uint8_t ielen = 0;

  if (iei > 0) {
    CHECK_IEI_DECODER(iei, *buffer);
    decoded++;
  }

  memset(
      voicedomainpreferenceandueusagesetting, 0,
      sizeof(VoiceDomainPreferenceAndUeUsageSetting));
  ielen = *(buffer + decoded);
  decoded++;
  CHECK_LENGTH_DECODER(len - decoded, ielen);
  voicedomainpreferenceandueusagesetting->ue_usage_setting =
      (*(buffer + decoded) >> 2) & 0x1;
  voicedomainpreferenceandueusagesetting->voice_domain_for_eutran =
      *(buffer + decoded) & 0x3;
  decoded++;
#if NAS_DEBUG
  dump_voice_domain_preference_and_ue_usage_setting_xml(
      voicedomainpreferenceandueusagesetting, iei);
#endif
  return decoded;
}

int encode_voice_domain_preference_and_ue_usage_setting(
    VoiceDomainPreferenceAndUeUsageSetting*
        voicedomainpreferenceandueusagesetting,
    uint8_t iei, uint8_t* buffer, uint32_t len) {
  uint8_t* lenPtr;
  uint32_t encoded = 0;

  /*
   * Checking IEI and pointer
   */
  CHECK_PDU_POINTER_AND_LENGTH_ENCODER(
      buffer, VOICE_DOMAIN_PREFERENCE_AND_UE_USAGE_SETTING_MINIMUM_LENGTH, len);
#if NAS_DEBUG
  dump_voice_domain_preference_and_ue_usage_setting_xml(
      voicedomainpreferenceandueusagesetting, iei);
#endif

  if (iei > 0) {
    *buffer = iei;
    encoded++;
  }

  lenPtr = (buffer + encoded);
  encoded++;
  *(buffer + encoded) =
      0x00 | (voicedomainpreferenceandueusagesetting->ue_usage_setting << 2) |
      voicedomainpreferenceandueusagesetting->voice_domain_for_eutran;
  encoded++;
  *lenPtr = encoded - 1 - ((iei > 0) ? 1 : 0);
  return encoded;
}

void dump_voice_domain_preference_and_ue_usage_setting_xml(
    VoiceDomainPreferenceAndUeUsageSetting*
        voicedomainpreferenceandueusagesetting,
    uint8_t iei) {
  OAILOG_DEBUG(LOG_NAS, "<Voice domain preference and UE usage setting>\n");

  if (iei > 0)
    /*
     * Don't display IEI if = 0
     */
    OAILOG_DEBUG(LOG_NAS, "    <IEI>0x%X</IEI>\n", iei);

  OAILOG_DEBUG(
      LOG_NAS, "    <UE_USAGE_SETTING>%u</UE_USAGE_SETTING>\n",
      voicedomainpreferenceandueusagesetting->ue_usage_setting);
  OAILOG_DEBUG(
      LOG_NAS, "    <VOICE_DOMAIN_FOR_EUTRAN>%u</VOICE_DOMAIN_FOR_EUTRAN>\n",
      voicedomainpreferenceandueusagesetting->voice_domain_for_eutran);
  OAILOG_DEBUG(LOG_NAS, "</Voice domain preference and UE usage setting>\n");
}
