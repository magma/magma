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

#ifndef FILE_TRACKING_AREA_UPDATE_COMPLETE_SEEN
#define FILE_TRACKING_AREA_UPDATE_COMPLETE_SEEN
#include <stdint.h>

#include "SecurityHeaderType.h"
#include "MessageType.h"
#include "3gpp_23.003.h"
#include "3gpp_24.007.h"
#include "3gpp_24.008.h"

/* Minimum length macro. Formed by minimum length of each mandatory field */
#define TRACKING_AREA_UPDATE_COMPLETE_MINIMUM_LENGTH (0)

/* Maximum length macro. Formed by maximum length of each field */
#define TRACKING_AREA_UPDATE_COMPLETE_MAXIMUM_LENGTH (0)

/*
 * Message name: Tracking area update complete
 * Description: This message shall be sent by the UE to the network in response
 * to a tracking area update accept message if a GUTI has been changed or a new
 * TMSI has been assigned. See table 8.2.27.1. Significance: dual Direction: UE
 * to network
 */

typedef struct tracking_area_update_complete_msg_tag {
  /* Mandatory fields */
  eps_protocol_discriminator_t protocoldiscriminator : 4;
  security_header_type_t securityheadertype : 4;
  message_type_t messagetype;
} tracking_area_update_complete_msg;

int decode_tracking_area_update_complete(
    tracking_area_update_complete_msg* trackingareaupdatecomplete,
    uint8_t* buffer, uint32_t len);

int encode_tracking_area_update_complete(
    tracking_area_update_complete_msg* trackingareaupdatecomplete,
    uint8_t* buffer, uint32_t len);

#endif /* ! defined(FILE_TRACKING_AREA_UPDATE_COMPLETE_SEEN) */
