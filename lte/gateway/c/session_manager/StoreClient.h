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

#include "SessionState.h"

namespace magma {
namespace lte {

typedef std::unordered_map<std::string, std::vector<std::unique_ptr<SessionState>>> SessionMap;

/**
 * StoreClient is responsible for reading/writing sessions to/from storage.
 */
class StoreClient {
 public:
  /**
   * Directly read the subscriber's sessions from storage
   *
   * If one or more of the subscribers have no sessions, empty entries will be
   * returned.
   * @param subscriber_ids typically in IMSI
   * @return All sessions for the subscribers
   */
  virtual SessionMap read_sessions(std::vector<std::string> subscriber_ids) = 0;

  /**
   * Directly write the subscriber sessions into storage, overwriting previous
   * values.
   *
   * @param sessions Sessions to write into storage
   * @return True if writes have completed successfully for all sessions.
   */
  virtual bool write_sessions(SessionMap sessions) = 0;
};

} // namespace lte
} // namespace magma
