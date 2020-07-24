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

#ifndef FILE_AUTHENTICATION_REQUEST_SEEN
#define FILE_AUTHENTICATION_REQUEST_SEEN

#include <stdint.h>

#include "SecurityHeaderType.h"
#include "MessageType.h"
#include "NasKeySetIdentifier.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define AUTHENTICATION_REQUEST_MINIMUM_LENGTH                                  \
  (NAS_KEY_SET_IDENTIFIER_MINIMUM_LENGTH +                                     \
   AUTHENTICATION_PARAMETER_RAND_IE_MIN_LENGTH - 1 +                           \
   AUTHENTICATION_PARAMETER_AUTN_IE_MIN_LENGTH - 2)

/* Maximum length macro. Formed by maximum length of each field */
#define AUTHENTICATION_REQUEST_MAXIMUM_LENGTH                                  \
  (NAS_KEY_SET_IDENTIFIER_MAXIMUM_LENGTH +                                     \
   AUTHENTICATION_PARAMETER_RAND_IE_MAX_LENGTH +                               \
   AUTHENTICATION_PARAMETER_AUTN_IE_MAX_LENGTH)

/*
 * Message name: Authentication request
 * Description: This message is sent by the network to the UE to initiate
 * authentication of the UE identity. See table 8.2.7.1. Significance: dual
 * Direction: network to UE
 */

typedef struct authentication_request_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  security_header_type_t securityheadertype : 4;
  message_type_t messagetype;
  NasKeySetIdentifier naskeysetidentifierasme;
  authentication_parameter_rand_t authenticationparameterrand;
  authentication_parameter_autn_t authenticationparameterautn;
} authentication_request_msg;

int decode_authentication_request(
    authentication_request_msg* authenticationrequest, uint8_t* buffer,
    uint32_t len);

int encode_authentication_request(
    authentication_request_msg* authenticationrequest, uint8_t* buffer,
    uint32_t len);

#endif /* ! defined(FILE_AUTHENTICATION_REQUEST_SEEN) */
