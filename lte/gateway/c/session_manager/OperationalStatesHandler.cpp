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
#include "lte/gateway/c/session_manager/OperationalStatesHandler.hpp"

#include <folly/dynamic.h>
#include <folly/json.h>
#include <glog/logging.h>
#include <google/protobuf/stubs/status.h>
#include <google/protobuf/stubs/stringpiece.h>
#include <google/protobuf/util/json_util.h>
#include <lte/protos/policydb.pb.h>
#include <lte/protos/session_manager.pb.h>
#include <list>
#include <map>
#include <memory>
#include <ostream>
#include <string>
#include <unordered_map>
#include <utility>
#include <vector>

#include "lte/gateway/c/session_manager/EnumToString.hpp"
#include "lte/gateway/c/session_manager/SessionState.hpp"
#include "lte/gateway/c/session_manager/SessionStore.hpp"
#include "lte/gateway/c/session_manager/Types.hpp"
#include "orc8r/gateway/c/common/logging/magma_logging.hpp"

namespace magma {

OpState get_operational_states(magma::SessionStore* session_store) {
  std::list<std::map<std::string, std::string>> states;
  auto session_map = session_store->read_all_sessions();
  for (auto& it : session_map) {
    std::map<std::string, std::string> state;
    state[TYPE] = SUBSCRIBER_STATE_TYPE;
    state[DEVICE_ID] = it.first;
    folly::dynamic sessions_by_apn = folly::dynamic::object;

    for (auto& session : it.second) {
      const auto apn = session->get_config().common_context.apn();
      if (sessions_by_apn[apn].empty()) {
        sessions_by_apn[apn] = folly::dynamic::array;
      }
      sessions_by_apn[apn].push_back(get_dynamic_session_state(session));
    }
    state[VALUE] = folly::toJson(sessions_by_apn);
    states.push_back(state);
  }
  return states;
}

folly::dynamic get_dynamic_session_state(
    const std::unique_ptr<SessionState>& session) {
  folly::dynamic state = folly::dynamic::object;
  const auto config = session->get_config().common_context;
  state[SESSION_ID] = session->get_session_id();
  state[MSISDN] = config.msisdn();
  state[magma::IPV4] = config.ue_ipv4();
  state[APN] = config.apn();
  state[SESSION_START_TIME] = session->get_pdp_start_time();
  state[LIFECYCLE_STATE] = session_fsm_state_to_str(session->get_state());
  state[ACTIVE_DURATION_SECOND] = session->get_active_duration_in_seconds();
  state[ACTIVE_POLICY_RULES] = get_dynamic_active_policies(session);
  return state;
}

folly::dynamic get_dynamic_active_policies(
    const std::unique_ptr<SessionState>& session) {
  google::protobuf::util::JsonPrintOptions options;
  options.add_whitespace = false;

  folly::dynamic policies = folly::dynamic::array;
  auto active_policies = session->get_all_active_policies();
  for (auto& policy : active_policies) {
    std::string json_policy;
    auto status = google::protobuf::util::MessageToJsonString(
        policy, &json_policy, options);
    if (!status.ok()) {
      MLOG(MERROR) << "Error serializing PolicyRule " << policy.id()
                   << " to JSON: " << status.ToString();
      continue;
    }
    policies.push_back(json_policy);
  }
  return policies;
}

}  // namespace magma
