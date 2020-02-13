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

namespace magma {
namespace lte {

typedef std::unordered_map<std::string, std::vector<std::unique_ptr<SessionState>>> SessionMap;
typedef std::function<void(SessionMap)> CallBackOnAccess;

/**
 * SessionStore acts as a broker to storage of sessiond state.
 *
 * This allows sessiond to service gRPC requests in a stateless manner.
 * Instead of keeping state in memory, sessiond uses the request parameters and
 * fetches state through SessionStore, handles the request, then writes back
 * to SessionStore, and responds to the gRPC request.
 *
 * SessionStore is intended to be a thread-safe singleton. Each gRPC request
 * should make a single read from SessionStore, and make a single write after
 * the request is serviced. The transactional nature of how requests should be
 * handled is intended to keep sessiond restartable in case of crashes.
 */
class SessionStore {
 public:
  virtual ~SessionStore() = default;

  /**
   * Requests R/W access to specified subscribers.
   *
   * It is intended that for each gRPC call that session_manager handles,
   * operate_on_sessions is called once.
   *
   *   gRPC request -> Request R/W on Sessions -> Read Sessions -> Update
   *   -> Commit Subscribers -> respond to gRPC call
   *
   * Only one caller can have R/W access to any given subscriber at a time, and
   * the subscriber is freed after commit.
   *
   * @param requested_subscriber_ids
   *    Subscribers for which R/W access is requested.
   * @param callback_on_access
   *    Callback for when R/W access on specified subscribers is granted.
   *    After updating the granted sessions, they should be returned to be
   *    written and committed to the SessionStore.
   *    @param commit
   *       Callback for when sessions are committed back to SessionStore.
   *       Input specifies whether the commit was successful.
   */
  virtual void operate_on_sessions(
    const std::vector<std::string>& requested_subscriber_ids,
    CallBackOnAccess callback) = 0;

  virtual void commit_sessions(
    SessionMap session_map,
    std::function<void(bool)> callback) = 0;
};

} // namespace lte
} // namespace magma
