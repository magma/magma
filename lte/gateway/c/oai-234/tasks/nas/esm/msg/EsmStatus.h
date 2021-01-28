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

#ifndef ESM_STATUS_H_
#define ESM_STATUS_H_

#include <stdint.h>

#include "3gpp_24.007.h"
#include "EsmCause.h"
#include "MessageType.h"

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define ESM_STATUS_MINIMUM_LENGTH (ESM_CAUSE_MINIMUM_LENGTH)

/* Maximum length macro. Formed by maximum length of each field */
#define ESM_STATUS_MAXIMUM_LENGTH (ESM_CAUSE_MAXIMUM_LENGTH)

/*
 * Message name: ESM status
 * Description: This message is sent by the network or the UE to pass
 * information on the status of the indicated EPS bearer context and report
 * certain error conditions (e.g. as listed in clause 7). See table 8.3.15.1.
 * Significance: dual
 * Direction: both
 */

typedef struct esm_status_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  ebi_t epsbeareridentity : 4;
  pti_t proceduretransactionidentity;
  message_type_t messagetype;
  esm_cause_t esmcause;
} esm_status_msg;

int decode_esm_status(esm_status_msg* esmstatus, uint8_t* buffer, uint32_t len);

int encode_esm_status(esm_status_msg* esmstatus, uint8_t* buffer, uint32_t len);

#endif /* ! defined(ESM_STATUS_H_) */
