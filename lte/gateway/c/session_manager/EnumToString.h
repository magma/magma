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

#include <lte/protos/abort_session.pb.h>
#include <lte/protos/session_manager.pb.h>
#include <string>

#include "ChargingGrant.h"
#include "ServiceAction.h"
#include "StoredState.h"
#include "Types.h"
#include "lte/protos/pipelined.pb.h"

namespace magma {
std::string reauth_state_to_str(ReAuthState state);

std::string service_state_to_str(ServiceState state);

std::string final_action_to_str(ChargingCredit_FinalAction final_action);

std::string grant_type_to_str(GrantTrackingType grant_type);

std::string session_fsm_state_to_str(SessionFsmState state);

std::string credit_update_type_to_str(CreditUsage::UpdateType update);

std::string raa_result_to_str(ReAuthResult res);

std::string asr_result_to_str(AbortSessionResult_Code res);

std::string wallet_state_to_str(SubscriberQuotaUpdate_Type state);

std::string service_action_type_to_str(ServiceActionType action);

std::string credit_validity_to_str(CreditValidity validity);

std::string event_trigger_to_str(EventTrigger event_trigger);

std::string request_origin_type_to_str(
    RequestOriginType_OriginType request_type);
}  // namespace magma
