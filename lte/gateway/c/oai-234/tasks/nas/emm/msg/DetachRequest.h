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

#ifndef FILE_DETACH_REQUEST_SEEN
#define FILE_DETACH_REQUEST_SEEN

#include <stdint.h>

#include "SecurityHeaderType.h"
#include "MessageType.h"
#include "DetachType.h"
#include "NasKeySetIdentifier.h"
#include "EpsMobileIdentity.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define DETACH_REQUEST_MINIMUM_LENGTH                                          \
  (DETACH_TYPE_MINIMUM_LENGTH + NAS_KEY_SET_IDENTIFIER_MINIMUM_LENGTH +        \
   EPS_MOBILE_IDENTITY_MINIMUM_LENGTH)

/* Maximum length macro. Formed by maximum length of each field */
#define DETACH_REQUEST_MAXIMUM_LENGTH                                          \
  (DETACH_TYPE_MAXIMUM_LENGTH + NAS_KEY_SET_IDENTIFIER_MAXIMUM_LENGTH +        \
   EPS_MOBILE_IDENTITY_MAXIMUM_LENGTH)

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define NW_DETACH_REQUEST_MINIMUM_LENGTH DETACH_TYPE_MINIMUM_LENGTH

/*
 * Use a reserved message type to differentiate NW initiated detach and UE
 * initiated detach
 */
#define NW_DETACH_REQUEST 0b01000111                  // 71 = 0x47
#define NW_DETACH_TYPE_RE_ATTACH_REQUIRED 0b00000001  // Re-attach required
#define NW_DETACH_TYPE_RE_ATTACH_NOT_REQUIRED                                  \
  0b00000010                                   // Re-attach not required
#define NW_DETACH_TYPE_IMSI_DETACH 0b00000011  // IMSI Detach
#define NW_DETACH_REQ_EMM_CAUSE_PRESENCE 0b00000001
#define NW_DETACH_REQ_EMM_CAUSE_IEI 0x53

/*
 * Message name: Detach request
 * Description: This message is sent by the UE to request the release of an EMM
 * context. Significance: dual Direction: UE to network
 */

typedef struct detach_request_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  security_header_type_t securityheadertype : 4;
  message_type_t messagetype;
  detach_type_t detachtype;
  NasKeySetIdentifier naskeysetidentifier;
  eps_mobile_identity_t gutiorimsi;
} detach_request_msg;

int decode_detach_request(
    detach_request_msg* detachrequest, uint8_t* buffer, uint32_t len);

int encode_detach_request(
    detach_request_msg* detachrequest, uint8_t* buffer, uint32_t len);

typedef struct nw_detach_request_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  security_header_type_t securityheadertype : 4;
  message_type_t messagetype;
  uint8_t nw_detachtype;
  uint8_t presenceMask;
  uint8_t emm_cause;
} nw_detach_request_msg;

int encode_nw_detach_request(
    nw_detach_request_msg* nw_detachrequest, uint8_t* buffer, uint32_t len);

#endif /* ! defined(FILE_DETACH_REQUEST_SEEN) */
