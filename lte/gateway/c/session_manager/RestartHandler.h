/**
 * Copyright (c) 2016-present, Facebook, Inc.
 * All rights reserved.
 *
 * This source code is licensed under the BSD-style license found in the
 * LICENSE file in the root directory of this source tree. An additional grant
 * of patent rights can be found in the PATENTS file in the same directory.
 */
#pragma once

#include <future>

#include "AAAClient.h"
#include "SessionReporter.h"
#include "DirectorydClient.h"
#include "LocalEnforcer.h"

namespace magma {
namespace sessiond {

/**
 * Restart handler cleans up previous sessions after a SessionD service restart
 */
class RestartHandler {
 public:
  RestartHandler(
      std::shared_ptr<AsyncDirectorydClient> directoryd_client,
      std::shared_ptr<aaa::AsyncAAAClient> aaa_client,
      std::shared_ptr<LocalEnforcer> enforcer,
      SessionReporter* reporter, SessionStore& session_store);

  /**
   * Cleanup previous sessions stored in directoryD
   */
  void cleanup_previous_sessions();

  /**
   * Re-create AAA sessions stored in sessiond
   */
  void setup_aaa_sessions();
 private:
  void terminate_previous_session(
      const std::string& sid, const std::string& session_id);
  bool populate_sessions_to_terminate_with_retries();
  bool launch_threads_to_terminate_with_retries();

 private:
  SessionStore& session_store_;
  std::shared_ptr<LocalEnforcer> enforcer_;
  std::shared_ptr<AsyncDirectorydClient> directoryd_client_;
  std::shared_ptr<aaa::AsyncAAAClient> aaa_client_;
  SessionReporter* reporter_;
  std::mutex sessions_to_terminate_lock_; // mutex to guard add/remove access to sessions_to_terminate
  std::unordered_map<std::string, std::string> sessions_to_terminate_;
  static const uint max_cleanup_retries_;
  static const uint rpc_retry_interval_s_;
};
}  // namespace sessiond
}  // namespace magma
