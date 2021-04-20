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

#ifndef ESM_INFORMATION_REQUEST_H_
#define ESM_INFORMATION_REQUEST_H_

#include <stdint.h>

#include "MessageType.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define ESM_INFORMATION_REQUEST_MINIMUM_LENGTH (0)

/* Maximum length macro. Formed by maximum length of each field */
#define ESM_INFORMATION_REQUEST_MAXIMUM_LENGTH (0)

/*
 * Message name: ESM information request
 * Description: This message is sent by the network to the UE to request the UE
 * to provide ESM information, i.e. protocol configuration options or APN or
 * both. See table 8.3.13.1. Significance: dual Direction: network to UE
 */

typedef struct esm_information_request_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  ebi_t epsbeareridentity : 4;
  pti_t proceduretransactionidentity;
  message_type_t messagetype;
} esm_information_request_msg;

int decode_esm_information_request(
    esm_information_request_msg* esminformationrequest, uint8_t* buffer,
    uint32_t len);

int encode_esm_information_request(
    esm_information_request_msg* esminformationrequest, uint8_t* buffer,
    uint32_t len);

#endif /* ! defined(ESM_INFORMATION_REQUEST_H_) */
