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

#ifndef ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_ACCEPT_H_
#define ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_ACCEPT_H_

#include <stdint.h>

#include "MessageType.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_ACCEPT_MINIMUM_LENGTH (0)

/* Maximum length macro. Formed by maximum length of each field */
#define ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_ACCEPT_MAXIMUM_LENGTH              \
  (PROTOCOL_CONFIGURATION_OPTIONS_IE_MAX_LENGTH)

/* If an optional value is present and should be encoded, the corresponding
 * Bit mask should be set to 1.
 */
#define ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_ACCEPT_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT \
  (1 << 0)

typedef enum activate_default_eps_bearer_context_accept_iei_tag {
  ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_ACCEPT_PROTOCOL_CONFIGURATION_OPTIONS_IEI =
      SM_PROTOCOL_CONFIGURATION_OPTIONS_IEI,
} activate_default_eps_bearer_context_accept_iei;

/*
 * Message name: Activate default EPS bearer context accept
 * Description: This message is sent by the UE to the network to acknowledge
 * activation of a default EPS bearer context. See table 8.3.4.1. Significance:
 * dual Direction: UE to network
 */

typedef struct activate_default_eps_bearer_context_accept_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  ebi_t epsbeareridentity : 4;
  pti_t proceduretransactionidentity;
  message_type_t messagetype;
  /* Optional fields */
  uint32_t presencemask;
  protocol_configuration_options_t protocolconfigurationoptions;
} activate_default_eps_bearer_context_accept_msg;

int decode_activate_default_eps_bearer_context_accept(
    activate_default_eps_bearer_context_accept_msg*
        activatedefaultepsbearercontextaccept,
    uint8_t* buffer, uint32_t len);

int encode_activate_default_eps_bearer_context_accept(
    activate_default_eps_bearer_context_accept_msg*
        activatedefaultepsbearercontextaccept,
    uint8_t* buffer, uint32_t len);

#endif /* ! defined(ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_ACCEPT_H_) */
