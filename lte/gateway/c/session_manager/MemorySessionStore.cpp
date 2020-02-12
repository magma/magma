/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */

#include "SessionState.h"
#include "MemorySessionStore.h"
#include "magma_logging.h"

namespace magma {
namespace lte {

MemorySessionStore::MemorySessionStore(): session_map_(), granted_subscribers_(), evb_(nullptr) {}

void MemorySessionStore::attach_event_base(folly::EventBase* evb)
{
  evb_ = evb;
}

void MemorySessionStore::operate_on_sessions(
  const std::vector<std::string> &requested_subscriber_ids,
  CallBackOnAccess callback)
{
  if (!granted_subscribers_.empty()) {
    throw std::runtime_error("Cannot hold two SessionStore grants at once.");
  }
  granted_subscribers_ = requested_subscriber_ids;
  SessionMap session_map;
  for (const auto& imsi : requested_subscriber_ids) {
    if (session_map_.find(imsi) != session_map_.end()) {
      session_map[imsi] = std::move(session_map_[imsi]);
    }
  }
  if (evb_) {
    evb_->runInEventBaseThread([session_map = std::move(session_map), callback] () mutable -> void {
      callback(std::move(session_map));
    });
  } else {
    MLOG(MERROR) << "No EventBase attached to SessionStore. "
                 << "Running synchronously.";
    callback(std::move(session_map));
  }
}

void MemorySessionStore::commit_sessions(
  SessionMap session_map,
  std::function<void(bool)> callback)
{
  for (auto& it : session_map) {
    if (
      std::find(
        granted_subscribers_.begin(), granted_subscribers_.end(), it.first) !=
      granted_subscribers_.end()) {
      session_map_[it.first] = std::move(it.second);
    }
  }
  granted_subscribers_.clear();
  if (evb_) {
    evb_->runInEventBaseThread([=] {
      callback(true);
    });
  } else {
    MLOG(MERROR) << "No EventBase attached to SessionStore. "
                 << "Running synchronously.";
    callback(true);
  }
}

} // namespace lte
} // namespace magma
