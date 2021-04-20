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

#ifndef FILE_SECURITY_MODE_COMPLETE_SEEN
#define FILE_SECURITY_MODE_COMPLETE_SEEN

#include "SecurityHeaderType.h"
#include "MessageType.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define SECURITY_MODE_COMPLETE_MINIMUM_LENGTH (0)

/* Maximum length macro. Formed by maximum length of each field */
#define SECURITY_MODE_COMPLETE_MAXIMUM_LENGTH (MOBILE_IDENTITY_MAXIMUM_LENGTH)

/* If an optional value is present and should be encoded, the corresponding
 * Bit mask should be set to 1.
 */
#define SECURITY_MODE_COMPLETE_IMEISV_PRESENT (1 << 0)

typedef enum security_mode_complete_iei_tag {
  SECURITY_MODE_COMPLETE_IMEISV_IEI = C_MOBILE_IDENTITY_IEI,
} security_mode_complete_iei;

/*
 * Message name: Security mode complete
 * Description: This message is sent by the UE to the network in response to a
 * SECURITY MODE COMMAND message. See table 8.2.21.1. Significance: dual
 * Direction: UE to network
 */

typedef struct security_mode_complete_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  security_header_type_t securityheadertype : 4;
  message_type_t messagetype;
  /* Optional fields */
  uint32_t presencemask;
  mobile_identity_t imeisv;
} security_mode_complete_msg;

int decode_security_mode_complete(
    security_mode_complete_msg* securitymodecomplete, uint8_t* buffer,
    uint32_t len);

int encode_security_mode_complete(
    security_mode_complete_msg* securitymodecomplete, uint8_t* buffer,
    uint32_t len);

#endif /* ! defined(FILE_SECURITY_MODE_COMPLETE_SEEN) */
