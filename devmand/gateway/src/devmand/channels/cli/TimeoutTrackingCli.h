// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <devmand/channels/cli/Cli.h>
#include <folly/Executor.h>
#include <folly/executors/SerialExecutor.h>
#include <folly/futures/Future.h>
#include <folly/futures/ThreadWheelTimekeeper.h>

namespace devmand::channels::cli {

using namespace std;
using namespace folly;
using devmand::channels::cli::Command;

static constexpr std::chrono::seconds defaultCommandTimeout =
    std::chrono::seconds(5);

// CLI layer that should be instanciated below QueuedCli. It throws
// FutureTimeout if Future returned by
// underlying layer does not return result within specified time period
// (timeout).
class TimeoutTrackingCli : public Cli {
 public:
  static shared_ptr<TimeoutTrackingCli> make(
      string id,
      shared_ptr<Cli> cli,
      shared_ptr<folly::Timekeeper> timekeeper,
      shared_ptr<folly::Executor> executor,
      std::chrono::milliseconds _timeoutInterval = defaultCommandTimeout);

  SemiFuture<Unit> destroy() override;

  ~TimeoutTrackingCli() override;

  folly::SemiFuture<std::string> executeRead(const ReadCommand cmd) override;

  folly::SemiFuture<std::string> executeWrite(const WriteCommand cmd) override;

  TimeoutTrackingCli(
      string id,
      shared_ptr<Cli> _cli,
      shared_ptr<folly::Timekeeper> _timekeeper,
      shared_ptr<folly::Executor> _executor,
      std::chrono::milliseconds _timeoutInterval);

 private:
  string id;
  shared_ptr<Cli> sharedCli; // underlying cli layer
  shared_ptr<folly::Timekeeper> sharedTimekeeper;
  shared_ptr<folly::Executor> sharedExecutor;
  const std::chrono::milliseconds timeoutInterval;

  Future<string> executeSomething(
      const Command& cmd,
      const string&& loggingPrefix,
      const function<SemiFuture<string>(shared_ptr<Cli> cli)>& innerFunc);
};
} // namespace devmand::channels::cli
