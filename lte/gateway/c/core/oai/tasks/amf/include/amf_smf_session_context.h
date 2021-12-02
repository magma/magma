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

#include "lte/gateway/c/core/oai/tasks/amf/amf_app_ue_context_and_proc.h"

#define PDU_SESSION_DEFAULT_AMBR 1
namespace magma5g {
void amf_app_free_smf_context(struct ue_m5gmm_context_s *ue_context);

// Set the pdu sesision state
status_code_e amf_smf_set_pdu_session_state(
  std::shared_ptr<smf_context_t> &smf_ctx, SMSessionFSMState sess_fsm_state);

// Get the pdu sesision state
status_code_e amf_smf_get_pdu_session_state(
  std::shared_ptr<smf_context_t> &smf_ctx, SMSessionFSMState *sess_fsm_state);
}  // namespace magma5g

