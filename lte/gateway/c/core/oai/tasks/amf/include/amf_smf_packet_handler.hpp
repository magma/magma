/**
 * Copyright 2020 The Magma Authors.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree.
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */
#pragma once

#include <sstream>
#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.hpp"
#include "lte/gateway/c/core/oai/tasks/nas5g/include/M5gNasMessage.h"

namespace magma5g {
#define PAYLOAD_CONTAINER_TAG_LENGTH 2
#define AMF_CAUSE_SUCCESS 1
#define MAX_UE_PDU_SESSION_LIMIT 15
#define MAX_UE_INITIAL_PDU_SESSION_ESTABLISHMENT_REQ_ALLOWED 5

status_code_e handle_sm_message_routing_failure(amf_ue_ngap_id_t ue_id,
                                                ULNASTransportMsg* msg,
                                                M5GMmCause m5gmmcause);
int amf_max_pdu_session_reject(amf_ue_ngap_id_t ue_id, ULNASTransportMsg* msg);
status_code_e amf_pdu_session_establishment_reject(amf_ue_ngap_id_t ue_id,
                                                   uint8_t session_id,
                                                   uint8_t pti, uint8_t cause);
int construct_pdu_session_reject_dl_req(uint8_t sequence_number,
                                        uint8_t session_id, uint8_t pti,
                                        uint8_t cause, bool is_security_enabled,
                                        amf_nas_message_t* msg);
M5GSmCause amf_smf_get_smcause(amf_ue_ngap_id_t ue_id, ULNASTransportMsg* msg);
M5GMmCause amf_smf_validate_context(amf_ue_ngap_id_t ue_id,
                                    ULNASTransportMsg* msg);
}  // namespace magma5g
