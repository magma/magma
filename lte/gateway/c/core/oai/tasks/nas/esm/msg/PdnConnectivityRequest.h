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

#ifndef PDN_CONNECTIVITY_REQUEST_H_
#define PDN_CONNECTIVITY_REQUEST_H_

#include <stdint.h>

#include "MessageType.h"
#include "NasRequestType.h"
#include "EsmInformationTransferFlag.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"
#include "3gpp_33.401.h"
#include "security_types.h"
#include "common_types.h"
#include "PdnType.h"

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define PDN_CONNECTIVITY_REQUEST_MINIMUM_LENGTH (PDN_TYPE_MINIMUM_LENGTH)

/* Maximum length macro. Formed by maximum length of each field */
#define PDN_CONNECTIVITY_REQUEST_MAXIMUM_LENGTH                                \
  (PDN_TYPE_MAXIMUM_LENGTH + ESM_INFORMATION_TRANSFER_FLAG_MAXIMUM_LENGTH +    \
   ACCESS_POINT_NAME_IE_MAX_LENGTH +                                           \
   PROTOCOL_CONFIGURATION_OPTIONS_IE_MAX_LENGTH)

/* If an optional value is present and should be encoded, the corresponding
 * Bit mask should be set to 1.
 */
#define PDN_CONNECTIVITY_REQUEST_ESM_INFORMATION_TRANSFER_FLAG_PRESENT (1 << 0)
#define PDN_CONNECTIVITY_REQUEST_ACCESS_POINT_NAME_PRESENT (1 << 1)
#define PDN_CONNECTIVITY_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT (1 << 2)

typedef enum pdn_connectivity_request_iei_tag {
  PDN_CONNECTIVITY_REQUEST_ESM_INFORMATION_TRANSFER_FLAG_IEI =
      0xD0, /* 0xD0 = 208 */
  PDN_CONNECTIVITY_REQUEST_ACCESS_POINT_NAME_IEI = SM_ACCESS_POINT_NAME_IEI,
  PDN_CONNECTIVITY_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_IEI =
      SM_PROTOCOL_CONFIGURATION_OPTIONS_IEI,
  PDN_CONNECTIVITY_REQUEST_DEVICE_PROPERTIES_IEI          = 0xC0,
  PDN_CONNECTIVITY_REQUEST_DEVICE_PROPERTIES_LOW_PRIO_IEI = 0xC1,
} pdn_connectivity_request_iei;

/*
 * Message name: PDN connectivity request
 * Description: This message is sent by the UE to the network to initiate
 * establishment of a PDN connection. See table 8.3.20.1. Significance: dual
 * Direction: UE to network
 */

typedef struct pdn_connectivity_request_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  ebi_t epsbeareridentity : 4;
  pti_t proceduretransactionidentity;
  message_type_t messagetype;
  request_type_t requesttype;
  pdn_type_t pdntype;
  /* Optional fields */
  uint32_t presencemask;
  esm_information_transfer_flag_t esminformationtransferflag;
  access_point_name_t accesspointname;
  protocol_configuration_options_t protocolconfigurationoptions;
} pdn_connectivity_request_msg;

int decode_pdn_connectivity_request(
    pdn_connectivity_request_msg* pdnconnectivityrequest, uint8_t* buffer,
    uint32_t len);

int encode_pdn_connectivity_request(
    pdn_connectivity_request_msg* pdnconnectivityrequest, uint8_t* buffer,
    uint32_t len);

#endif /* ! defined(PDN_CONNECTIVITY_REQUEST_H_) */
