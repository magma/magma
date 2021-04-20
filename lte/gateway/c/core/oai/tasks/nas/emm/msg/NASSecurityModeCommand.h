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

#ifndef FILE_NAS_SECURITY_MODE_COMMAND_SEEN
#define FILE_NAS_SECURITY_MODE_COMMAND_SEEN

#include "SecurityHeaderType.h"
#include "MessageType.h"
#include "NasSecurityAlgorithms.h"
#include "NasKeySetIdentifier.h"
#include "UeAdditionalSecurityCapability.h"
#include "UeSecurityCapability.h"
#include "Nonce.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define SECURITY_MODE_COMMAND_MINIMUM_LENGTH                                   \
  (NAS_SECURITY_ALGORITHMS_MINIMUM_LENGTH +                                    \
   NAS_KEY_SET_IDENTIFIER_MINIMUM_LENGTH +                                     \
   UE_SECURITY_CAPABILITY_MINIMUM_LENGTH)

/* Maximum length macro. Formed by maximum length of each field */
#define SECURITY_MODE_COMMAND_MAXIMUM_LENGTH                                   \
  (NAS_SECURITY_ALGORITHMS_MAXIMUM_LENGTH +                                    \
   NAS_KEY_SET_IDENTIFIER_MAXIMUM_LENGTH +                                     \
   UE_SECURITY_CAPABILITY_MAXIMUM_LENGTH + IMEISV_REQUEST_IE_MAX_LENGTH +      \
   NONCE_MAXIMUM_LENGTH + NONCE_MAXIMUM_LENGTH +                               \
   UE_ADDITIONAL_SECURITY_CAPABILITY_MAXIMUM_LENGTH)

/* If an optional value is present and should be encoded, the corresponding
 * Bit mask should be set to 1.
 */
#define SECURITY_MODE_COMMAND_IMEISV_REQUEST_PRESENT (1 << 0)
#define SECURITY_MODE_COMMAND_REPLAYED_NONCEUE_PRESENT (1 << 1)
#define SECURITY_MODE_COMMAND_NONCEMME_PRESENT (1 << 2)
#define SECURITY_MODE_COMMAND_REPLAYED_UE_ADDITIONAL_SECU_CAPABILITY_PRESENT   \
  (1 << 3)

typedef enum security_mode_command_iei_tag {
  SECURITY_MODE_COMMAND_IMEISV_REQUEST_IEI =
      GMM_IMEISV_REQUEST_IEI,                        /* 0xC0 = 192 */
  SECURITY_MODE_COMMAND_REPLAYED_NONCEUE_IEI = 0x55, /* 0x55 = 85 */
  SECURITY_MODE_COMMAND_NONCEMME_IEI         = 0x56, /* 0x56 = 86 */
  SECURITY_MODE_COMMAND_REPLAYED_UE_ADDITIONAL_SECURITY_CAPABILITY_IEI = 0x6F,
} security_mode_command_iei;

/*
 * Message name: Security mode command
 * Description: This message is sent by the network to the UE to establish NAS
 * signalling security. See table 8.2.20.1. Significance: dual Direction:
 * network to UE
 */

typedef struct security_mode_command_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  security_header_type_t securityheadertype : 4;
  message_type_t messagetype;
  NasSecurityAlgorithms selectednassecurityalgorithms;
  NasKeySetIdentifier naskeysetidentifier;
  ue_additional_security_capability_t replayedueadditionalsecuritycapabilities;
  ue_security_capability_t replayeduesecuritycapabilities;
  /* Optional fields */
  uint32_t presencemask;
  imeisv_request_t imeisvrequest;
  nonce_t replayednonceue;
  nonce_t noncemme;
} security_mode_command_msg;

int decode_security_mode_command(
    security_mode_command_msg* securitymodecommand, uint8_t* buffer,
    uint32_t len);

int encode_security_mode_command(
    security_mode_command_msg* securitymodecommand, uint8_t* buffer,
    uint32_t len);

#endif /* ! defined(FILE_NAS_SECURITY_MODE_COMMAND_SEEN) */
