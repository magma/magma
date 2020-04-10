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
#include "MemoryStoreClient.h"
#include "StoredState.h"
#include "RedisStoreClient.h"
#include "RuleStore.h"

namespace magma {
namespace lte {

typedef std::
  unordered_map<std::string, std::vector<std::unique_ptr<SessionState>>>
    SessionMap;
// Value int represents the request numbers needed for requests to PCRF
typedef std::set<std::string> SessionRead;
typedef std::unordered_map<
  std::string,
  std::unordered_map<std::string, SessionStateUpdateCriteria>>
  SessionUpdate;

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
  static SessionUpdate get_default_session_update(SessionMap& session_map);

  SessionStore(std::shared_ptr<StaticRuleStore> rule_store);

  SessionStore(
    std::shared_ptr<StaticRuleStore> rule_store,
    std::shared_ptr<RedisStoreClient> store_client);

  /**
   * Read the last written values for the requested sessions through the
   * storage interface.
   * @param req
   * @return Last written values for requested sessions. Returns an empty vector
   *         for subscribers that do not have active sessions.
   */
  SessionMap read_sessions(const SessionRead& req);

  /**
   * Read the last written values for the requested sessions through the
   * storage interface. This also modifies the request_numbers stored before
   * returning the SessionMap to the caller.
   * NOTE: It is assumed that the correct number of request_numbers are
   *       reserved on each read_sessions call. If more requests are made to
   *       the OCS/PCRF than are requested, this can cause undefined behavior.
   * @param req
   * @return Last written values for requested sessions. Returns an empty vector
   *         for subscribers that do not have active sessions.
   */
  SessionMap read_sessions_for_reporting(const SessionRead& req);

  /**
   * Read the last written values for the requested sessions through the
   * storage interface. This also modifies the request_numbers stored before
   * returning the SessionMap to the caller, incremented by one for each
   * session.
   * NOTE: It is assumed that the correct number of request_numbers are
   *       reserved on each read_sessions call. If more requests are made to
   *       the OCS/PCRF than are requested, this can cause undefined behavior.
   * NOTE: Here, it is expected that the caller will use one additional
   *       request_number for each session.
   * @param req
   * @return Last written values for requested sessions. Returns an empty vector
   *         for subscribers that do not have active sessions.
   */
  SessionMap read_sessions_for_deletion(const SessionRead& req);

  /**
   * Create sessions for a subscriber. Redundant creations will fail.
   * @param subscriber_id
   * @param sessions
   * @return true if successful, otherwise the update to storage is discarded.
   */
  bool create_sessions(
    const std::string& subscriber_id,
    std::vector<std::unique_ptr<SessionState>> sessions);

  /**
   * Attempt to update sessions with update criteria. If any update to any of
   * the sessions is invalid, the whole update request is assumed to be invalid,
   * and nothing in storage will be overwritten.
   * @param update_criteria
   * @return true if successful, otherwise the update to storage is discarded.
   */
  bool update_sessions(const SessionUpdate& update_criteria);


 private:
  static bool merge_into_session(
    std::unique_ptr<SessionState>& session,
    const SessionStateUpdateCriteria& update_criteria);

 private:
  std::shared_ptr<StoreClient> store_client_;
  std::shared_ptr<StaticRuleStore> rule_store_;
};

} // namespace lte
} // namespace magma
