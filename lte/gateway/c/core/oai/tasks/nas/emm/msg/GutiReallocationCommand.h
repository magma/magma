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

#ifndef FILE_GUTI_REALLOCATION_COMMAND_SEEN
#define FILE_GUTI_REALLOCATION_COMMAND_SEEN

#include <stdint.h>

#include "SecurityHeaderType.h"
#include "MessageType.h"
#include "NasMessageContainer.h"
#include "EpsMobileIdentity.h"
#include "TrackingAreaIdentityList.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define GUTI_REALLOCATION_COMMAND_MINIMUM_LENGTH                               \
  (EPS_MOBILE_IDENTITY_MINIMUM_LENGTH)

/* Maximum length macro. Formed by maximum length of each field */
#define GUTI_REALLOCATION_COMMAND_MAXIMUM_LENGTH                               \
  (EPS_MOBILE_IDENTITY_MAXIMUM_LENGTH +                                        \
   TRACKING_AREA_IDENTITY_LIST_MAXIMUM_LENGTH)

/* If an optional value is present and should be encoded, the corresponding
 * Bit mask should be set to 1.
 */
#define GUTI_REALLOCATION_COMMAND_TAI_LIST_PRESENT (1 << 0)

typedef enum guti_reallocation_command_iei_tag {
  GUTI_REALLOCATION_COMMAND_TAI_LIST_IEI = 0x54, /* 0x54 = 84 */
} guti_reallocation_command_iei;

/*
 * Message name: GUTI reallocation command
 * Description: This message is sent by the network to the UE to reallocate a
 * GUTI and optionally to provide a new TAI list. See table 8.2.16.1.
 * Significance: dual
 * Direction: network to UE
 */

typedef struct guti_reallocation_command_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  security_header_type_t securityheadertype : 4;
  message_type_t messagetype;
  eps_mobile_identity_t guti;
  /* Optional fields */
  uint32_t presencemask;
  tai_list_t tailist;
} guti_reallocation_command_msg;

int decode_guti_reallocation_command(
    guti_reallocation_command_msg* gutireallocationcommand, uint8_t* buffer,
    uint32_t len);

int encode_guti_reallocation_command(
    guti_reallocation_command_msg* gutireallocationcommand, uint8_t* buffer,
    uint32_t len);

#endif /* ! defined(FILE_GUTI_REALLOCATION_COMMAND_SEEN) */
