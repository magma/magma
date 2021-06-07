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

#ifndef FILE_AUTHENTICATION_RESPONSE_SEEN
#define FILE_AUTHENTICATION_RESPONSE_SEEN

#include <stdint.h>

#include "SecurityHeaderType.h"
#include "MessageType.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define AUTHENTICATION_RESPONSE_MINIMUM_LENGTH                                 \
  (AUTHENTICATION_RESPONSE_PARAMETER_IE_MIN_LENGTH)

/* Maximum length macro. Formed by maximum length of each field */
#define AUTHENTICATION_RESPONSE_MAXIMUM_LENGTH                                 \
  (AUTHENTICATION_RESPONSE_PARAMETER_IE_MAX_LENGTH)

/*
 * Message name: Authentication response
 * Description: This message is sent by the UE to the network to deliver a
 * calculated authentication response to the network. See table 8.2.8.1.
 * Significance: dual
 * Direction: UE to network
 */

typedef struct authentication_response_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  security_header_type_t securityheadertype : 4;
  message_type_t messagetype;
  authentication_response_parameter_t authenticationresponseparameter;
} authentication_response_msg;

int decode_authentication_response(
    authentication_response_msg* authenticationresponse, uint8_t* buffer,
    uint32_t len);

int encode_authentication_response(
    authentication_response_msg* authenticationresponse, uint8_t* buffer,
    uint32_t len);

#endif /* ! defined(FILE_AUTHENTICATION_RESPONSE_SEEN) */
