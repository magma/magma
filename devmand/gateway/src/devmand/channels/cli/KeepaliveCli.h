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
      shared_ptr<folly::Timekeeper> _timekeeper,
      chrono::milliseconds heartbeatInterval = defaultKeepaliveInterval,
      string keepAliveCommand = "",
      chrono::milliseconds backoffAfterKeepaliveTimeout = chrono::seconds(5));

  KeepaliveCli(
      string id,
      shared_ptr<Cli> _cli,
      shared_ptr<folly::Executor> parentExecutor,
      shared_ptr<folly::Timekeeper> _timekeeper,
      chrono::milliseconds heartbeatInterval,
      string keepAliveCommand,
      chrono::milliseconds backoffAfterKeepaliveTimeout);

  ~KeepaliveCli() override;

  folly::SemiFuture<std::string> executeRead(const ReadCommand cmd) override;

  folly::SemiFuture<std::string> executeWrite(const WriteCommand cmd) override;

 private:
  struct KeepaliveParameters {
    string id;
    shared_ptr<Cli> cli; // underlying cli layer
    shared_ptr<folly::Timekeeper> timekeeper;
    shared_ptr<folly::Executor> parentExecutor;
    folly::Executor::KeepAlive<folly::SerialExecutor> serialExecutorKeepAlive;
    string keepAliveCommand;
    chrono::milliseconds heartbeatInterval;
    chrono::milliseconds backoffAfterKeepaliveTimeout;
    atomic<bool> shutdown;

    KeepaliveParameters(
        const string& _id,
        const shared_ptr<Cli>& _cli,
        const shared_ptr<folly::Timekeeper>& _timekeeper,
        const shared_ptr<folly::Executor>& _parentExecutor,
        const folly::Executor::KeepAlive<folly::SerialExecutor>&
            _serialExecutorKeepAlive,
        const string& _keepAliveCommand,
        const chrono::milliseconds& _heartbeatInterval,
        const chrono::milliseconds& _backoffAfterKeepaliveTimeout,
        const bool _shutdown);

    KeepaliveParameters(KeepaliveParameters&&) = default;
  };

  shared_ptr<KeepaliveParameters> keepaliveParameters;

  static void triggerSendKeepAliveCommand(
      shared_ptr<KeepaliveParameters> keepaliveParameters);
};
} // namespace devmand::channels::cli
