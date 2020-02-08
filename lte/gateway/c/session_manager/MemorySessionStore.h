/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <memory>

#include <lte/protos/session_manager.grpc.pb.h>
#include <folly/io/async/EventBaseManager.h>

#include "SessionState.h"
#include "SessionStore.h"

namespace magma {
namespace lte {

/**
 * In memory session store
 * Only allows a single data grant at a time
 * This is only used to emulate a persistent storage interface.
 *
 * TODO: Modify to allow concurrent data access
 * TODO: Deep-copy of stored objects so that retrieving a grant does not clear
 *       out the retrieved entries in MemorySessionStore
 */
class MemorySessionStore final : public SessionStore {
 public:
  MemorySessionStore();
  MemorySessionStore(MemorySessionStore const&) = delete;
  MemorySessionStore(MemorySessionStore&&) = default;

  void attach_event_base(folly::EventBase *evb);

  void operate_on_sessions(
    const std::vector<std::string>& requested_subscriber_ids,
    CallBackOnAccess callback) final;

  void commit_sessions (
    SessionMap session_map,
    std::function<void(bool)> callback) final;

 private:
  ~MemorySessionStore() = default;

 private:
  SessionMap session_map_;
  std::vector<std::string> granted_subscribers_;
  folly::EventBase* evb_;
};

/**
 * TODO: Add RedisSessionStore for persistent storage in Redis resistant to
 *       sessiond restarts
 */

} // namespace lte
} // namespace magma
