/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "SessionState.h"
#include "MemoryStoreClient.h"
#include "magma_logging.h"

namespace magma {
namespace lte {

MemoryStoreClient::MemoryStoreClient(
  std::shared_ptr<StaticRuleStore> rule_store):
  session_map_({}),
  rule_store_(rule_store) {}

SessionMap MemoryStoreClient::read_sessions(std::vector<std::string> subscriber_ids)
{
  auto session_map = SessionMap{};
  for (const auto& subscriber_id : subscriber_ids) {
    auto sessions = std::vector<std::unique_ptr<SessionState>>{};
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

bool MemoryStoreClient::write_sessions(SessionMap session_map)
{
  for (auto& it : session_map) {
    auto sessions = std::vector<StoredSessionState>{};
    for (auto const& session : it.second) {
      auto stored_session = session->marshal();
      sessions.push_back(stored_session);
    }
    session_map_[it.first] = std::move(sessions);
  }
  return true;
}

} // namespace lte
} // namespace magma
