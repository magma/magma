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

#ifndef ESM_INFORMATION_RESPONSE_H_
#define ESM_INFORMATION_RESPONSE_H_

#include <stdint.h>

#include "MessageType.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define ESM_INFORMATION_RESPONSE_MINIMUM_LENGTH (0)

/* Maximum length macro. Formed by maximum length of each field */
#define ESM_INFORMATION_RESPONSE_MAXIMUM_LENGTH                                \
  (ACCESS_POINT_NAME_MAXIMUM_LENGTH +                                          \
   PROTOCOL_CONFIGURATION_OPTIONS_MAXIMUM_LENGTH)

/* If an optional value is present and should be encoded, the corresponding
 * Bit mask should be set to 1.
 */
#define ESM_INFORMATION_RESPONSE_ACCESS_POINT_NAME_PRESENT (1 << 0)
#define ESM_INFORMATION_RESPONSE_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT (1 << 1)

typedef enum esm_information_response_iei_tag {
  ESM_INFORMATION_RESPONSE_ACCESS_POINT_NAME_IEI = SM_ACCESS_POINT_NAME_IEI,
  ESM_INFORMATION_RESPONSE_PROTOCOL_CONFIGURATION_OPTIONS_IEI =
      SM_PROTOCOL_CONFIGURATION_OPTIONS_IEI,
} esm_information_response_iei;

/*
 * Message name: ESM information response
 * Description: This message is sent by the UE to the network in response to an
 * ESM INFORMATION REQUEST message and provides the requested ESM information.
 * See table 8.3.14.1. Significance: dual Direction: UE to network
 */

typedef struct esm_information_response_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  ebi_t epsbeareridentity : 4;
  pti_t proceduretransactionidentity;
  message_type_t messagetype;
  /* Optional fields */
  uint32_t presencemask;
  access_point_name_t accesspointname;
  protocol_configuration_options_t protocolconfigurationoptions;
} esm_information_response_msg;

int decode_esm_information_response(
    esm_information_response_msg* esminformationresponse, uint8_t* buffer,
    uint32_t len);

int encode_esm_information_response(
    esm_information_response_msg* esminformationresponse, uint8_t* buffer,
    uint32_t len);

#endif /* ! defined(ESM_INFORMATION_RESPONSE_H_) */
