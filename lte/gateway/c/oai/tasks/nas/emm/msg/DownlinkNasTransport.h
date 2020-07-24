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

#ifndef FILE_DOWNLINK_NAS_TRANSPORT_SEEN
#define FILE_DOWNLINK_NAS_TRANSPORT_SEEN

#include <stdint.h>

#include "SecurityHeaderType.h"
#include "MessageType.h"
#include "NasMessageContainer.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define DOWNLINK_NAS_TRANSPORT_MINIMUM_LENGTH                                  \
  (NAS_MESSAGE_CONTAINER_MINIMUM_LENGTH)

/* Maximum length macro. Formed by maximum length of each field */
#define DOWNLINK_NAS_TRANSPORT_MAXIMUM_LENGTH                                  \
  (NAS_MESSAGE_CONTAINER_MAXIMUM_LENGTH)

/*
 * Message name: Downlink NAS Transport
 * Description: This message is sent by the network to the UE in order to carry
 * an SMS message in encapsulated format. See table 8.2.12.1. Significance: dual
 * Direction: network to UE
 */

typedef struct downlink_nas_transport_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  security_header_type_t securityheadertype : 4;
  message_type_t messagetype;
  NasMessageContainer nasmessagecontainer;
} downlink_nas_transport_msg;

int decode_downlink_nas_transport(
    downlink_nas_transport_msg* downlinknastransport, uint8_t* buffer,
    uint32_t len);

int encode_downlink_nas_transport(
    downlink_nas_transport_msg* downlinknastransport, uint8_t* buffer,
    uint32_t len);

#endif /* ! defined(FILE_DOWNLINK_NAS_TRANSPORT_SEEN) */
