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

#ifndef FILE_IDENTITY_REQUEST_SEEN
#define FILE_IDENTITY_REQUEST_SEEN

#include <stdint.h>

#include "SecurityHeaderType.h"
#include "MessageType.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define IDENTITY_REQUEST_MINIMUM_LENGTH (IDENTITY_TYPE_2_IE_MIN_LENGTH)

/* Maximum length macro. Formed by maximum length of each field */
#define IDENTITY_REQUEST_MAXIMUM_LENGTH (IDENTITY_TYPE_2_IE_MAX_LENGTH)

/*
 * Message name: Identity request
 * Description: This message is sent by the network to the UE to request the UE
 * to provide the specified identity. See table 8.2.18.1. Significance: dual
 * Direction: network to UE
 */

typedef struct identity_request_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  security_header_type_t securityheadertype : 4;
  message_type_t messagetype;
  identity_type2_t identitytype;
} identity_request_msg;

int decode_identity_request(
    identity_request_msg* identityrequest, uint8_t* buffer, uint32_t len);

int encode_identity_request(
    identity_request_msg* identityrequest, uint8_t* buffer, uint32_t len);

#endif /* ! defined(FILE_IDENTITY_REQUEST_SEEN) */
