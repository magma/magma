// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <devmand/channels/cli/Cli.h>
#include <folly/Synchronized.h>
#include <folly/container/EvictingCacheMap.h>

namespace devmand::channels::cli {

using folly::EvictingCacheMap;
using folly::Future;
using folly::Synchronized;
using std::shared_ptr;
using std::string;
using CliCache = Synchronized<EvictingCacheMap<string, string>>;

class ReadCachingCli : public Cli {
 private:
  string id;
  shared_ptr<Cli> cli{};
  shared_ptr<CliCache> cache;

 public:
  ReadCachingCli(
      string _id,
      const shared_ptr<Cli>& _cli,
      const shared_ptr<CliCache>& _cache);

  static shared_ptr<CliCache> createCache();
  Future<string> executeAndRead(const Command& cmd) override;
  Future<string> execute(const Command& cmd) override;
};
} // namespace devmand::channels::cli