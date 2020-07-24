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

#ifndef MODIFY_EPS_BEARER_CONTEXT_REQUEST_H_
#define MODIFY_EPS_BEARER_CONTEXT_REQUEST_H_
#include <stdint.h>

#include "MessageType.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"
#include "ApnAggregateMaximumBitRate.h"
#include "EpsQualityOfService.h"
#include "RadioPriority.h"

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define MODIFY_EPS_BEARER_CONTEXT_REQUEST_MINIMUM_LENGTH (0)

/* Maximum length macro. Formed by maximum length of each field */
#define MODIFY_EPS_BEARER_CONTEXT_REQUEST_MAXIMUM_LENGTH                       \
  (EPS_QUALITY_OF_SERVICE_MAXIMUM_LENGTH +                                     \
   TRAFFIC_FLOW_TEMPLATE_MAXIMUM_LENGTH + QUALITY_OF_SERVICE_MAXIMUM_LENGTH +  \
   LLC_SERVICE_ACCESS_POINT_IDENTIFIER_MAXIMUM_LENGTH +                        \
   RADIO_PRIORITY_MAXIMUM_LENGTH + PACKET_FLOW_IDENTIFIER_MAXIMUM_LENGTH +     \
   APN_AGGREGATE_MAXIMUM_BIT_RATE_IE_MAX_LENGTH +                              \
   PROTOCOL_CONFIGURATION_OPTIONS_IE_MAX_LENGTH)

/* If an optional value is present and should be encoded, the corresponding
 * Bit mask should be set to 1.
 */
#define MODIFY_EPS_BEARER_CONTEXT_REQUEST_NEW_EPS_QOS_PRESENT (1 << 0)
#define MODIFY_EPS_BEARER_CONTEXT_REQUEST_TFT_PRESENT (1 << 1)
#define MODIFY_EPS_BEARER_CONTEXT_REQUEST_NEW_QOS_PRESENT (1 << 2)
#define MODIFY_EPS_BEARER_CONTEXT_REQUEST_NEGOTIATED_LLC_SAPI_PRESENT (1 << 3)
#define MODIFY_EPS_BEARER_CONTEXT_REQUEST_RADIO_PRIORITY_PRESENT (1 << 4)
#define MODIFY_EPS_BEARER_CONTEXT_REQUEST_PACKET_FLOW_IDENTIFIER_PRESENT       \
  (1 << 5)
#define MODIFY_EPS_BEARER_CONTEXT_REQUEST_APNAMBR_PRESENT (1 << 6)
#define MODIFY_EPS_BEARER_CONTEXT_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_PRESENT \
  (1 << 7)

typedef enum modify_eps_bearer_context_request_iei_tag {
  MODIFY_EPS_BEARER_CONTEXT_REQUEST_NEW_EPS_QOS_IEI = 0x5B, /* 0x5B = 91 */
  MODIFY_EPS_BEARER_CONTEXT_REQUEST_TFT_IEI     = SM_TRAFFIC_FLOW_TEMPLATE_IEI,
  MODIFY_EPS_BEARER_CONTEXT_REQUEST_NEW_QOS_IEI = SM_QUALITY_OF_SERVICE_IEI,
  MODIFY_EPS_BEARER_CONTEXT_REQUEST_NEGOTIATED_LLC_SAPI_IEI =
      SM_LLC_SERVICE_ACCESS_POINT_IDENTIFIER_IEI,
  MODIFY_EPS_BEARER_CONTEXT_REQUEST_RADIO_PRIORITY_IEI = 0x80, /* 0x80 = 128 */
  MODIFY_EPS_BEARER_CONTEXT_REQUEST_PACKET_FLOW_IDENTIFIER_IEI =
      SM_PACKET_FLOW_IDENTIFIER_IEI,
  MODIFY_EPS_BEARER_CONTEXT_REQUEST_APNAMBR_IEI = 0x5E, /* 0x5E = 94 */
  MODIFY_EPS_BEARER_CONTEXT_REQUEST_PROTOCOL_CONFIGURATION_OPTIONS_IEI =
      SM_PROTOCOL_CONFIGURATION_OPTIONS_IEI,
} modify_eps_bearer_context_request_iei;

/*
 * Message name: Modify EPS bearer context request
 * Description: This message is sent by the network to the UE to request
 * modification of an active EPS bearer context. See table 8.3.18.1.
 * Significance: dual
 * Direction: network to UE
 */

typedef struct modify_eps_bearer_context_request_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  ebi_t epsbeareridentity : 4;
  pti_t proceduretransactionidentity;
  message_type_t messagetype;
  /* Optional fields */
  uint32_t presencemask;
  EpsQualityOfService newepsqos;
  traffic_flow_template_t tft;
  quality_of_service_t newqos;
  llc_service_access_point_identifier_t negotiatedllcsapi;
  radio_priority_t radiopriority;
  packet_flow_identifier_t packetflowidentifier;
  ApnAggregateMaximumBitRate apnambr;
  protocol_configuration_options_t protocolconfigurationoptions;
} modify_eps_bearer_context_request_msg;

int decode_modify_eps_bearer_context_request(
    modify_eps_bearer_context_request_msg* modifyepsbearercontextrequest,
    uint8_t* buffer, uint32_t len);

int encode_modify_eps_bearer_context_request(
    modify_eps_bearer_context_request_msg* modifyepsbearercontextrequest,
    uint8_t* buffer, uint32_t len);

#endif /* ! defined(MODIFY_EPS_BEARER_CONTEXT_REQUEST_H_) */
