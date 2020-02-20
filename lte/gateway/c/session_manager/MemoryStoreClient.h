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

#include "StoreClient.h"
#include "StoredState.h"

namespace magma {
namespace lte {

/**
 * Non-persistent StoreClient used as intermediate stage in development
 */
class MemoryStoreClient final : public StoreClient {
 public:
  MemoryStoreClient(std::shared_ptr<StaticRuleStore> rule_store);
  MemoryStoreClient(MemoryStoreClient const&) = delete;
  MemoryStoreClient(MemoryStoreClient&&) = default;

  SessionMap read_sessions(std::vector<std::string> subscriber_ids);

  bool write_sessions(SessionMap session_map);

 private:
  ~MemoryStoreClient() = default;

 private:
  std::unordered_map<std::string, std::vector<StoredSessionState>> session_map_;
  std::shared_ptr<StaticRuleStore> rule_store_;
};

} // namespace lte
} // namespace magma
