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

#include <set>
#include <string>
#include <utility>
#include <vector>

#include "magma_logging.h"
#include "MemoryStoreClient.h"
#include "SessionState.h"

namespace magma {
namespace lte {

MemoryStoreClient::MemoryStoreClient(
    std::shared_ptr<StaticRuleStore> rule_store)
    : session_map_({}), rule_store_(rule_store) {}

SessionMap MemoryStoreClient::read_sessions(
    std::set<std::string> subscriber_ids) {
  auto session_map = SessionMap{};
  for (const auto& subscriber_id : subscriber_ids) {
    auto sessions = SessionVector{};
    if (session_map_.find(subscriber_id) != session_map_.end()) {
      for (auto& stored_session : session_map_[subscriber_id]) {
        auto session = SessionState::unmarshal(stored_session, *rule_store_);
        sessions.push_back(std::move(session));
      }
    }
    session_map[subscriber_id] = std::move(sessions);
  }
  return session_map;
}

SessionMap MemoryStoreClient::read_all_sessions() {
  auto session_map = SessionMap{};
  for (auto& it : session_map_) {
    auto sessions = SessionVector{};
    for (auto& stored_session : it.second) {
      auto session = SessionState::unmarshal(stored_session, *rule_store_);
      sessions.push_back(std::move(session));
    }
    session_map[it.first] = std::move(sessions);
  }
  return session_map;
}

bool MemoryStoreClient::write_sessions(SessionMap session_map) {
  for (auto& it : session_map) {
    auto sessions = std::vector<StoredSessionState>{};
    for (auto const& session : it.second) {
      auto stored_session = session->marshal();
      sessions.push_back(stored_session);
    }
    if (sessions.empty()) {
      // if session is empty that means subs should be deleted from the map
      session_map_.erase(it.first);
      continue;
    }
    session_map_[it.first] = std::move(sessions);
  }
  return true;
}

}  // namespace lte
}  // namespace magma
