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

#include <sstream>
#include <string>

#include "EnumToString.h"

namespace magma {
std::string reauth_state_to_str(ReAuthState state) {
  switch (state) {
    case REAUTH_NOT_NEEDED:
      return "REAUTH_NOT_NEEDED";
    case REAUTH_REQUIRED:
      return "REAUTH_REQUIRED";
    case REAUTH_PROCESSING:
      return "REAUTH_PROCESSING";
    default:
      return "INVALID REAUTH STATE";
  }
}

std::string service_state_to_str(ServiceState state) {
  switch (state) {
    case SERVICE_ENABLED:
      return "SERVICE_ENABLED";
    case SERVICE_NEEDS_DEACTIVATION:
      return "SERVICE_NEEDS_DEACTIVATION";
    case SERVICE_NEEDS_SUSPENSION:
      return "SERVICE_NEEDS_SUSPENSION";
    case SERVICE_DISABLED:
      return "SERVICE_DISABLED";
    case SERVICE_NEEDS_ACTIVATION:
      return "SERVICE_NEEDS_ACTIVATION";
    case SERVICE_REDIRECTED:
      return "SERVICE_REDIRECTED";
    case SERVICE_RESTRICTED:
      return "SERVICE_RESTRICTED";
    default:
      return "INVALID SERVICE STATE";
  }
}

std::string final_action_to_str(ChargingCredit_FinalAction final_action) {
  switch (final_action) {
    case ChargingCredit_FinalAction_TERMINATE:
      return "TERMINATE";
    case ChargingCredit_FinalAction_REDIRECT:
      return "REDIRECT";
    case ChargingCredit_FinalAction_RESTRICT_ACCESS:
      return "RESTRICT_ACCESS";
    default:
      return "INVALID FINAL ACTION";
  }
}

std::string grant_type_to_str(GrantTrackingType grant_type) {
  switch (grant_type) {
    case TRACKING_UNSET:
      return "TRACKING_UNSET";
    case ALL_TOTAL_TX_RX:
      return "ALL_TOTAL_TX_RX";
    case TOTAL_ONLY:
      return "TOTAL_ONLY";
    case TX_ONLY:
      return "TX_ONLY";
    case RX_ONLY:
      return "RX_ONLY";
    case TX_AND_RX:
      return "TX_AND_RX";
    default:
      return "INVALID GRANT TRACKING TYPE";
  }
}

std::string session_fsm_state_to_str(SessionFsmState state) {
  switch (state) {
    case SESSION_ACTIVE:
    case ACTIVE:
      return "SESSION_ACTIVE";
    case SESSION_TERMINATED:
      return "SESSION_TERMINATED";
    case SESSION_RELEASED:
    case RELEASE:
      return "SESSION_RELEASED";
    case CREATING:
      return "SESSION_CREATING";
    case CREATED:
      return "SESSION_CREATED";
    default:
      return "INVALID SESSION FSM STATE";
  }
}

std::string credit_update_type_to_str(CreditUsage::UpdateType update) {
  switch (update) {
    case CreditUsage::THRESHOLD:
      return "THRESHOLD";
    case CreditUsage::QHT:
      return "QHT";
    case CreditUsage::TERMINATED:
      return "TERMINATED";
    case CreditUsage::QUOTA_EXHAUSTED:
      return "QUOTA_EXHAUSTED";
    case CreditUsage::VALIDITY_TIMER_EXPIRED:
      return "VALIDITY_TIMER_EXPIRED";
    case CreditUsage::OTHER_QUOTA_TYPE:
      return "OTHER_QUOTA_TYPE";
    case CreditUsage::RATING_CONDITION_CHANGE:
      return "RATING_CONDITION_CHANGE";
    case CreditUsage::REAUTH_REQUIRED:
      return "REAUTH_REQUIRED";
    case CreditUsage::POOL_EXHAUSTED:
      return "POOL_EXHAUSTED";
    default:
      return "INVALID CREDIT UPDATE TYPE";
  }
}

std::string raa_result_to_str(ReAuthResult res) {
  switch (res) {
    case UPDATE_INITIATED:
      return "UPDATE_INITIATED";
    case UPDATE_NOT_NEEDED:
      return "UPDATE_NOT_NEEDED";
    case SESSION_NOT_FOUND:
      return "SESSION_NOT_FOUND";
    case OTHER_FAILURE:
      return "OTHER_FAILURE";
    default:
      return "UNKNOWN_RESULT";
  }
}

std::string asr_result_to_str(AbortSessionResult_Code res) {
  switch (res) {
    case AbortSessionResult_Code_SESSION_REMOVED:
      return "SESSION_REMOVED";
    case AbortSessionResult_Code_SESSION_NOT_FOUND:
      return "SESSION_NOT_FOUND";
    case AbortSessionResult_Code_USER_NOT_FOUND:
      return "USER_NOT_FOUND";
    case AbortSessionResult_Code_GATEWAY_NOT_FOUND:
      return "GATEWAY_NOT_FOUND";
    case AbortSessionResult_Code_RADIUS_SERVER_ERROR:
      return "RADIUS_SERVER_ERROR";
    default:
      return "UNKNOWN_RESULT";
  }
}

std::string wallet_state_to_str(SubscriberQuotaUpdate_Type state) {
  switch (state) {
    case SubscriberQuotaUpdate_Type_VALID_QUOTA:
      return "VALID_QUOTA";
    case SubscriberQuotaUpdate_Type_NO_QUOTA:
      return "NO_QUOTA";
    case SubscriberQuotaUpdate_Type_TERMINATE:
      return "TERMINATE";
    default:
      return "INVALID";
  }
}

std::string service_action_type_to_str(ServiceActionType action) {
  switch (action) {
    case CONTINUE_SERVICE:
      return "CONTINUE_SERVICE";
    case TERMINATE_SERVICE:
      return "TERMINATE_SERVICE";
    case ACTIVATE_SERVICE:
      return "ACTIVATE_SERVICE";
    case REDIRECT:
      return "REDIRECT";
    case RESTRICT_ACCESS:
      return "RESTRICT_ACCESS";
    default:
      return "INVALID ACTION TYPE";
  }
}

std::string credit_validity_to_str(CreditValidity validity) {
  switch (validity) {
    case VALID_CREDIT:
      return "VALID_CREDIT";
    case INVALID_CREDIT:
      return "INVALID_CREDIT";
    case TRANSIENT_ERROR:
      return "TRANSIENT_ERROR";
    default:
      return "INVALID CREDIT TYPE";
  }
}

std::string event_trigger_to_str(EventTrigger event_trigger) {
  switch (event_trigger) {
    case USAGE_REPORT:
      return "USAGE_REPORT";
    case REVALIDATION_TIMEOUT:
      return "REVALIDATION_TIMEOUT";
    default:
      std::ostringstream message;
      message << "UNIMPLEMENTED EVENT TRIGGER: " << event_trigger;
      return message.str();
  }
}

std::string request_origin_type_to_str(
    RequestOriginType_OriginType request_type) {
  switch (request_type) {
    case RequestOriginType_OriginType_GX:
      return "GX";
    case RequestOriginType_OriginType_GY:
      return "GY";
    case RequestOriginType_OriginType_N4:
      return "N4";
    case RequestOriginType_OriginType_WILDCARD:
      return "WILDCARD";
    default:
      std::ostringstream message;
      message << "Unimplemented RequestOriginType: " << request_type;
      return message.str();
  }
}

}  // namespace magma
