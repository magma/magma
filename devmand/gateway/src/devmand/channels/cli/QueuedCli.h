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

namespace devmand::channels::cli {

using namespace std;
using namespace folly;

class QueuedCli : public Cli {
 private:
  string id;
  shared_ptr<Cli> cli;

  Executor::KeepAlive<SerialExecutor>
      serialExecutorKeepAlive; // maintain consumer thread

  struct QueueEntry {
    function<Future<string>()> obtainFutureFromCli;
    shared_ptr<Promise<string>> promise;
    string command;
    string loggingPrefix;
  };

  /**
   * Unbounded multi producer single consumer queue where consumer is not
   * blocked on dequeue.
   */
  UnboundedQueue<QueueEntry, false, true, false> queue;
  bool isProcessing = false; // only accessed from consumer thread

 public:
  QueuedCli(
      string _id,
      shared_ptr<Cli> _cli,
      const shared_ptr<Executor>& _parentExecutor);
  QueuedCli() = delete;
  QueuedCli(const QueuedCli&) = delete;

  Future<string> executeAndRead(const Command& cmd) override;
  Future<string> execute(const Command& cmd) override;

 private:
  Future<string> executeSomething(
      const Command& cmd,
      const string& prefix,
      function<Future<string>()> innerFunc);

  void triggerDequeue();
};
} // namespace devmand::channels::cli
