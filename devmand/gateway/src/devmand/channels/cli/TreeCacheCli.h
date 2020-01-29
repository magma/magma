// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <devmand/channels/cli/Cli.h>
#include <devmand/channels/cli/CliFlavour.h>
#include <devmand/channels/cli/TreeCache.h>
#include <folly/Executor.h>
#include <folly/futures/Future.h>

namespace devmand::channels::cli {

using namespace std;
using namespace folly;

class TreeCacheCli : public Cli {
 private:
  string id;
  shared_ptr<Cli> cli;
  shared_ptr<folly::Executor> executor;
  shared_ptr<CliFlavour> sharedCliFlavour;
  shared_ptr<TreeCache> cache;

 public:
  TreeCacheCli(
      string _id,
      shared_ptr<Cli> _cli,
      shared_ptr<folly::Executor> executor,
      shared_ptr<CliFlavour> sharedCliFlavour);

  // Visible for testing
  TreeCacheCli(
      string _id,
      shared_ptr<Cli> _cli,
      shared_ptr<folly::Executor> _executor,
      shared_ptr<CliFlavour> _sharedCliFlavour,
      shared_ptr<TreeCache> cache);

  SemiFuture<folly::Unit> destroy() override;

  ~TreeCacheCli() override;

  /*
   * Clear cache.
   */
  void clear();

  folly::SemiFuture<std::string> executeRead(const ReadCommand cmd) override;
  folly::SemiFuture<std::string> executeWrite(const WriteCommand cmd) override;
};
} // namespace devmand::channels::cli
