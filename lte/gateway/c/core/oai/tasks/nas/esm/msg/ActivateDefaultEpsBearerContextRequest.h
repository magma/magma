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

#ifndef ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_H_
#define ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_H_

#include <stdint.h>

#include "MessageType.h"
#include "EpsQualityOfService.h"
#include "PdnAddress.h"
#include "RadioPriority.h"
#include "ApnAggregateMaximumBitRate.h"
#include "EsmCause.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_MINIMUM_LENGTH             \
  (EPS_QUALITY_OF_SERVICE_MINIMUM_LENGTH + ACCESS_POINT_NAME_IE_MIN_LENGTH +   \
   PDN_ADDRESS_MINIMUM_LENGTH)

/* Maximum length macro. Formed by maximum length of each field */
#define ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_MAXIMUM_LENGTH             \
  (EPS_QUALITY_OF_SERVICE_MAXIMUM_LENGTH + ACCESS_POINT_NAME_IE_MAX_LENGTH +   \
   PDN_ADDRESS_MAXIMUM_LENGTH + TRANSACTION_IDENTIFIER_MAXIMUM_LENGTH +        \
   QUALITY_OF_SERVICE_IE_MAX_LENGTH +                                          \
   LLC_SERVICE_ACCESS_POINT_IDENTIFIER_MAXIMUM_LENGTH +                        \
   RADIO_PRIORITY_MAXIMUM_LENGTH + PACKET_FLOW_IDENTIFIER_MAXIMUM_LENGTH +     \
   APN_AGGREGATE_MAXIMUM_BIT_RATE_MAXIMUM_LENGTH + ESM_CAUSE_MAXIMUM_LENGTH +  \
   PROTOCOL_CONFIGURATION_OPTIONS_IE_MAX_LENGTH)

/* If an optional value is present and should be encoded, the corresponding
 * Bit mask should be set to 1.
 */
#define ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_TRANSACTION_IDENTIFIER_PRESENT \
  (1 << 0)
#define ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_NEGOTIATED_QOS_PRESENT     \
  (1 << 1)
#define ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_NEGOTIATED_LLC_SAPI_PRESENT \
  (1 << 2)
#define ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_RADIO_PRIORITY_PRESENT     \
  (1 << 3)
#define ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_PACKET_FLOW_IDENTIFIER_PRESENT \
  (1 << 4)
#define ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_APNAMBR_PRESENT (1 << 5)
#define ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_ESM_CAUSE_PRESENT (1 << 6)
#define ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT \
  (1 << 7)

typedef enum activate_default_eps_bearer_context_request_iei_tag {
  ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_TRANSACTION_IDENTIFIER_IEI =
      SM_LINKED_TI_IEI,
  ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_NEGOTIATED_QOS_IEI =
      SM_QUALITY_OF_SERVICE_IEI,
  ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_NEGOTIATED_LLC_SAPI_IEI =
      SM_LLC_SERVICE_ACCESS_POINT_IDENTIFIER_IEI,
  ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_RADIO_PRIORITY_IEI =
      0x80, /* 0x80 = 128 */
  ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_PACKET_FLOW_IDENTIFIER_IEI =
      SM_PACKET_FLOW_IDENTIFIER_IEI,
  ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_APNAMBR_IEI =
      0x5E, /* 0x5E = 94 */
  ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_ESM_CAUSE_IEI =
      0x58, /* 0x58 = 88 */
  ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_IEI =
      SM_PROTOCOL_CONFIGURATION_OPTIONS_IEI,
} activate_default_eps_bearer_context_request_iei;

/*
 * Message name: Activate default EPS bearer context request
 * Description: This message is sent by the network to the UE to request
 * activation of a default EPS bearer context. See table 8.3.6.1. Significance:
 * dual Direction: network to UE
 */

typedef struct activate_default_eps_bearer_context_request_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  ebi_t epsbeareridentity : 4;
  pti_t proceduretransactionidentity;
  message_type_t messagetype;
  EpsQualityOfService epsqos;
  access_point_name_t accesspointname;
  PdnAddress pdnaddress;
  /* Optional fields */
  uint32_t presencemask;
  linked_ti_t transactionidentifier;
  quality_of_service_t negotiatedqos;
  llc_service_access_point_identifier_t negotiatedllcsapi;
  radio_priority_t radiopriority;
  packet_flow_identifier_t packetflowidentifier;
  ApnAggregateMaximumBitRate apnambr;
  esm_cause_t esmcause;
  protocol_configuration_options_t protocolconfigurationoptions;
} activate_default_eps_bearer_context_request_msg;

int decode_activate_default_eps_bearer_context_request(
    activate_default_eps_bearer_context_request_msg*
        activatedefaultepsbearercontextrequest,
    uint8_t* buffer, uint32_t len);

int encode_activate_default_eps_bearer_context_request(
    activate_default_eps_bearer_context_request_msg*
        activatedefaultepsbearercontextrequest,
    uint8_t* buffer, uint32_t len);

#endif /* ! defined(ACTIVATE_DEFAULT_EPS_BEARER_CONTEXT_REQUEST_H_) */
