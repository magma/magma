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

#include <folly/dynamic.h>
#include <list>
#include <map>
#include <memory>
#include <string>

#include "SessionStore.h"

namespace magma {
class SessionState;
namespace lte {
class SessionStore;
}  // namespace lte

const std::string TYPE = "type";
const std::string SUBSCRIBER_STATE_TYPE = "subscriber_state";
const std::string DEVICE_ID = "device_id";
const std::string VALUE = "value";
const std::string SESSION_ID = "session_id";
const std::string IMSI = "imsi";
const std::string MSISDN = "msisdn";
const std::string APN = "apn";
const std::string IPV4 = "ipv4";
const std::string ACTIVE_POLICY_RULES = "active_policy_rules";
const std::string SESSION_START_TIME = "session_start_time";
const std::string ACTIVE_DURATION_SECOND = "active_duration_sec";
const std::string LIFECYCLE_STATE = "lifecycle_state";

using OpState = std::list<std::map<std::string, std::string>>;

OpState get_operational_states(magma::SessionStore* session_store);

folly::dynamic get_dynamic_session_state(
    const std::unique_ptr<SessionState>& session);

folly::dynamic get_dynamic_active_policies(
    const std::unique_ptr<SessionState>& session);

}  // namespace magma
