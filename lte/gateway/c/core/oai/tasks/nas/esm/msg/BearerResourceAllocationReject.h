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

#ifndef BEARER_RESOURCE_ALLOCATION_REJECT_H_
#define BEARER_RESOURCE_ALLOCATION_REJECT_H_

#include <stdint.h>

#include "MessageType.h"
#include "EsmCause.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define BEARER_RESOURCE_ALLOCATION_REJECT_MINIMUM_LENGTH                       \
  (ESM_CAUSE_MINIMUM_LENGTH)

/* Maximum length macro. Formed by maximum length of each field */
#define BEARER_RESOURCE_ALLOCATION_REJECT_MAXIMUM_LENGTH                       \
  (ESM_CAUSE_MAXIMUM_LENGTH + PROTOCOL_CONFIGURATION_OPTIONS_IE_MAX_LENGTH)

/* If an optional value is present and should be encoded, the corresponding
 * Bit mask should be set to 1.
 */
#define BEARER_RESOURCE_ALLOCATION_REJECT_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT \
  (1 << 0)

typedef enum bearer_resource_allocation_reject_iei_tag {
  BEARER_RESOURCE_ALLOCATION_REJECT_PROTOCOL_CONFIGURATION_OPTIONS_IEI =
      SM_PROTOCOL_CONFIGURATION_OPTIONS_IEI,
} bearer_resource_allocation_reject_iei;

/*
 * Message name: Bearer resource allocation reject
 * Description: This message is sent by the network to the UE to reject the
 * allocation of a dedicated bearer resource. See table 8.3.7.1. Significance:
 * dual Direction: network to UE
 */

typedef struct bearer_resource_allocation_reject_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  ebi_t epsbeareridentity : 4;
  pti_t proceduretransactionidentity;
  message_type_t messagetype;
  esm_cause_t esmcause;
  /* Optional fields */
  uint32_t presencemask;
  protocol_configuration_options_t protocolconfigurationoptions;
} bearer_resource_allocation_reject_msg;

int decode_bearer_resource_allocation_reject(
    bearer_resource_allocation_reject_msg* bearerresourceallocationreject,
    uint8_t* buffer, uint32_t len);

int encode_bearer_resource_allocation_reject(
    bearer_resource_allocation_reject_msg* bearerresourceallocationreject,
    uint8_t* buffer, uint32_t len);

#endif /* ! defined(BEARER_RESOURCE_ALLOCATION_REJECT_H_) */
