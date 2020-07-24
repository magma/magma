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

#ifndef FILE_EMM_INFORMATION_SEEN
#define FILE_EMM_INFORMATION_SEEN

#include <stdint.h>

#include "SecurityHeaderType.h"
#include "MessageType.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define EMM_INFORMATION_MINIMUM_LENGTH (0)

/* Maximum length macro. Formed by maximum length of each field */
#define EMM_INFORMATION_MAXIMUM_LENGTH                                         \
  (NETWORK_NAME_MAXIMUM_LENGTH + NETWORK_NAME_MAXIMUM_LENGTH +                 \
   TIME_ZONE_MAXIMUM_LENGTH + TIME_ZONE_AND_TIME_MAXIMUM_LENGTH +              \
   DAYLIGHT_SAVING_TIME_MAXIMUM_LENGTH)

/* If an optional value is present and should be encoded, the corresponding
 * Bit mask should be set to 1.
 */
#define EMM_INFORMATION_FULL_NAME_FOR_NETWORK_PRESENT (1 << 0)
#define EMM_INFORMATION_SHORT_NAME_FOR_NETWORK_PRESENT (1 << 1)
#define EMM_INFORMATION_LOCAL_TIME_ZONE_PRESENT (1 << 2)
#define EMM_INFORMATION_UNIVERSAL_TIME_AND_LOCAL_TIME_ZONE_PRESENT (1 << 3)
#define EMM_INFORMATION_NETWORK_DAYLIGHT_SAVING_TIME_PRESENT (1 << 4)

typedef enum emm_information_iei_tag {
  EMM_INFORMATION_FULL_NAME_FOR_NETWORK_IEI  = MM_FULL_NETWORK_NAME_IEI,
  EMM_INFORMATION_SHORT_NAME_FOR_NETWORK_IEI = MM_SHORT_NETWORK_NAME_IEI,
  EMM_INFORMATION_LOCAL_TIME_ZONE_IEI        = MM_TIME_ZONE_IEI,
  EMM_INFORMATION_UNIVERSAL_TIME_AND_LOCAL_TIME_ZONE_IEI =
      MM_TIME_ZONE_AND_TIME_IEI,
  EMM_INFORMATION_NETWORK_DAYLIGHT_SAVING_TIME_IEI = 0x49, /* 0x49 = 73 */
} emm_information_iei;

/*
 * Message name: EMM information
 * Description: This message is sent by the network at any time during EMM
 * context is established to send certain information to the UE. See
 * table 8.2.13.1. Significance: local Direction: network to UE
 */

typedef struct emm_information_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  security_header_type_t securityheadertype : 4;
  message_type_t messagetype;
  /* Optional fields */
  uint32_t presencemask;
  network_name_t fullnamefornetwork;
  network_name_t shortnamefornetwork;
  time_zone_t localtimezone;
  time_zone_and_time_t universaltimeandlocaltimezone;
  daylight_saving_time_t networkdaylightsavingtime;
} emm_information_msg;

int decode_emm_information(
    emm_information_msg* emminformation, uint8_t* buffer, uint32_t len);

int encode_emm_information(
    emm_information_msg* emminformation, uint8_t* buffer, uint32_t len);

#endif /* ! defined(FILE_EMM_INFORMATION_SEEN) */
