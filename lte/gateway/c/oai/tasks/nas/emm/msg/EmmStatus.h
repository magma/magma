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

#ifndef FILE_EMM_STATUS_SEEN
#define FILE_EMM_STATUS_SEEN

#include <stdint.h>

#include "SecurityHeaderType.h"
#include "MessageType.h"
#include "EmmCause.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define EMM_STATUS_MINIMUM_LENGTH (EMM_CAUSE_MINIMUM_LENGTH)

/* Maximum length macro. Formed by maximum length of each field */
#define EMM_STATUS_MAXIMUM_LENGTH (EMM_CAUSE_MAXIMUM_LENGTH)

/*
 * Message name: EMM status
 * Description: This message is sent by the UE or by the network at any time to
 * report certain error conditions listed in clause 7. See table 8.2.14.1.
 * Significance: local
 * Direction: both
 */

typedef struct emm_status_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  security_header_type_t securityheadertype : 4;
  message_type_t messagetype;
  emm_cause_t emmcause;
} emm_status_msg;

int decode_emm_status(emm_status_msg* emmstatus, uint8_t* buffer, uint32_t len);

int encode_emm_status(emm_status_msg* emmstatus, uint8_t* buffer, uint32_t len);

#endif /* ! defined(FILE_EMM_STATUS_SEEN) */
