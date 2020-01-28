// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <devmand/channels/cli/Cli.h>
#include <devmand/channels/cli/CliThreadWheelTimekeeper.h>
#include <devmand/channels/cli/CliTimekeeperWrapper.h>
#include <folly/Executor.h>
#include <folly/executors/SerialExecutor.h>
#include <folly/futures/Future.h>
#include <folly/futures/ThreadWheelTimekeeper.h>

namespace devmand::channels::cli {
using namespace std;
using devmand::channels::cli::Command;

static constexpr chrono::seconds defaultKeepaliveInterval = chrono::seconds(60);

// CLI layer that should be above QueuedCli. Periodically schedules keepalive
// command to prevent dropping
// of inactive connection.
class KeepaliveCli : public Cli {
 public:
  static shared_ptr<KeepaliveCli> make(
      string id,
      shared_ptr<Cli> _cli,
      shared_ptr<folly::Executor> parentExecutor,
      shared_ptr<CliThreadWheelTimekeeper> _timekeeper,
      chrono::milliseconds heartbeatInterval = defaultKeepaliveInterval,
      string keepAliveCommand = "",
      chrono::milliseconds backoffAfterKeepaliveTimeout = chrono::seconds(5));

  folly::SemiFuture<folly::Unit> destroy() override;

  KeepaliveCli(
      string id,
      shared_ptr<Cli> _cli,
      shared_ptr<folly::Executor> parentExecutor,
      shared_ptr<CliTimekeeperWrapper> _timekeeper,
      chrono::milliseconds heartbeatInterval,
      string keepAliveCommand,
      chrono::milliseconds backoffAfterKeepaliveTimeout);

  ~KeepaliveCli() override;

  folly::SemiFuture<std::string> executeRead(const ReadCommand cmd) override;

  folly::SemiFuture<std::string> executeWrite(const WriteCommand cmd) override;

 private:
  shared_ptr<Cli> sharedCli; // underlying cli layer
  shared_ptr<CliTimekeeperWrapper> sharedTimekeeper;
  shared_ptr<folly::Executor> parentExecutor;
  shared_ptr<folly::Executor::KeepAlive<folly::SerialExecutor>>
      sharedSerialExecutorKeepAlive;
  struct KeepaliveParameters {
    string id;
    weak_ptr<Cli> cli;
    weak_ptr<CliTimekeeperWrapper> timekeeper;
    weak_ptr<folly::Executor::KeepAlive<folly::SerialExecutor>>
        serialExecutorKeepAlive;
    string keepAliveCommand;
    chrono::milliseconds heartbeatInterval;
    chrono::milliseconds backoffAfterKeepaliveTimeout;

    KeepaliveParameters(
        const string& _id,
        weak_ptr<Cli> _cli,
        weak_ptr<CliTimekeeperWrapper> _timekeeper,
        weak_ptr<folly::Executor::KeepAlive<folly::SerialExecutor>>
            _serialExecutorKeepAlive,
        const string& _keepAliveCommand,
        const chrono::milliseconds& _heartbeatInterval,
        const chrono::milliseconds& _backoffAfterKeepaliveTimeout);

    KeepaliveParameters(KeepaliveParameters&&) = default;
  };

  shared_ptr<KeepaliveParameters> keepaliveParameters;

  static void triggerSendKeepAliveCommand(
      shared_ptr<KeepaliveParameters> keepaliveParameters);
};
} // namespace devmand::channels::cli
