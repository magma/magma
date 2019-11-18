// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/channels/cli/Cli.h>
#include <folly/executors/CPUThreadPoolExecutor.h>

namespace devmand {
namespace test {
namespace utils {
namespace cli {

using namespace devmand::channels::cli;
using namespace std;

class EchoCli : public Cli {
 public:
  folly::Future<string> executeAndRead(const Command& cmd) override {
    return folly::Future<string>(cmd.toString());
  }

  folly::Future<string> execute(const Command& cmd) override {
    return folly::Future<string>(cmd.toString());
  }
};

class ErrCli : public Cli {
 public:
  folly::Future<string> executeAndRead(const Command& cmd) override {
    throw runtime_error(cmd.toString());
    return folly::Future<string>(runtime_error(cmd.toString()));
  }

  folly::Future<string> execute(const Command& cmd) override {
    throw runtime_error(cmd.toString());
    return folly::Future<string>(runtime_error(cmd.toString()));
  }
};

class AsyncCli : public Cli {
 public:
  AsyncCli(
      shared_ptr<Cli> _cli,
      shared_ptr<folly::CPUThreadPoolExecutor> _executor,
      vector<unsigned int> _durations)
      : cli(_cli), executor(_executor), durations(_durations), index(0) {}

  folly::Future<string> executeAndRead(const Command& cmd) override {
    folly::Future<string> f = via(executor.get()).thenValue([=](...) {
      unsigned int tis = durations[(index++) % durations.size()];
      this_thread::sleep_for(chrono::seconds(tis));
      return cli->executeAndRead(cmd);
    });
    return f;
  }

  folly::Future<string> execute(const Command& cmd) override {
    (void)cmd;
    return folly::Future<string>(runtime_error("Unsupported"));
  }

 protected:
  shared_ptr<Cli> cli; // underlying cli layers
  shared_ptr<folly::CPUThreadPoolExecutor> executor;
  vector<unsigned int> durations;
  unsigned int index;
  bool quit;
};

} // namespace cli
} // namespace utils
} // namespace test
} // namespace devmand
