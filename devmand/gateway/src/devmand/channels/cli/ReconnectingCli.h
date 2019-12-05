// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <boost/thread/mutex.hpp>
#include <devmand/channels/cli/Cli.h>
#include <devmand/channels/cli/CliThreadWheelTimekeeper.h>
#include <devmand/channels/cli/CliTimekeeperWrapper.h>
#include <devmand/channels/cli/Command.h>
#include <folly/Executor.h>
#include <folly/executors/SerialExecutor.h>
#include <folly/futures/Future.h>
#include <folly/futures/ThreadWheelTimekeeper.h>

namespace devmand {
namespace channels {
namespace cli {

using namespace std;
using namespace folly;
using boost::mutex;
using devmand::channels::cli::Cli;
using devmand::channels::cli::Command;

class ReconnectingCli : public Cli {
 public:
  static shared_ptr<ReconnectingCli> make(
      string id,
      shared_ptr<Executor> executor,
      function<SemiFuture<shared_ptr<Cli>>()>&& createCliStack,
      shared_ptr<CliThreadWheelTimekeeper> timekeeper,
      chrono::milliseconds quietPeriod);

  SemiFuture<Unit> destroy() override;

  ~ReconnectingCli() override;

  folly::SemiFuture<std::string> executeRead(const ReadCommand cmd) override;

  folly::SemiFuture<std::string> executeWrite(const WriteCommand cmd) override;

  ReconnectingCli(
      string id,
      shared_ptr<Executor> executor,
      function<SemiFuture<shared_ptr<Cli>>()>&& createCliStack,
      shared_ptr<CliTimekeeperWrapper> timekeeper,
      chrono::milliseconds quietPeriod);

 private:
  struct ReconnectParameters {
    string id;
    shared_ptr<CliTimekeeperWrapper> timekeeper;
    shared_ptr<folly::Executor> executor;
    function<SemiFuture<shared_ptr<Cli>>()> createCliStack;
    mutex cliMutex;
    shared_ptr<Cli> maybeCli;
    std::chrono::milliseconds quietPeriod;
    atomic<bool> isReconnecting; // TODO: merge with maybeCli
    atomic<bool> shutdown;
  };

  shared_ptr<ReconnectParameters> reconnectParameters;

  SemiFuture<string> executeSomething(
      const string&& loggingPrefix,
      const function<SemiFuture<string>(shared_ptr<Cli>)>& innerFunc,
      const Command cmd);

  static void triggerReconnect(shared_ptr<ReconnectParameters> params);

  static Future<Unit> setReconnecting(shared_ptr<ReconnectParameters> params);
};
} // namespace cli
} // namespace channels
} // namespace devmand
