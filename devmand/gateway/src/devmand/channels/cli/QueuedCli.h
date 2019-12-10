// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <boost/thread/recursive_mutex.hpp>
#include <devmand/channels/cli/Cli.h>
#include <folly/Executor.h>
#include <folly/executors/SerialExecutor.h>
#include <folly/futures/Future.h>

namespace devmand::channels::cli {

using namespace std;
using namespace folly;
using boost::recursive_mutex;

/*
 * TODO: throw exception when queue is full
 */
class QueuedCli : public Cli {
 private:
  struct QueueEntry {
    function<SemiFuture<string>()> obtainFutureFromCli;
    shared_ptr<Promise<string>> promise;
    Command command = ReadCommand::create("dummy");
    string loggingPrefix;
  };

  struct QueuedParameters {
    string id;
    shared_ptr<Cli> cli;

    shared_ptr<Executor> parentExecutor;

    Executor::KeepAlive<SerialExecutor>
        serialExecutorKeepAlive; // maintain consumer thread

    /**
     * Unbounded multi producer single consumer queue where consumer is not
     * blocked on dequeue.
     */
    UnboundedQueue<QueueEntry, false, true, false>
        queue; // TODO: investigate priority queue for keepalive commands

    atomic<bool> isProcessing = ATOMIC_VAR_INIT(false);

    atomic<bool> shutdown = ATOMIC_VAR_INIT(false);

    recursive_mutex mutex;

    QueuedParameters(
        const string& id,
        const shared_ptr<Cli>& cli,
        const shared_ptr<Executor>& parentExecutor,
        const Executor::KeepAlive<SerialExecutor>& serialExecutorKeepAlive);
  };
  shared_ptr<QueuedParameters> queuedParameters;

  SemiFuture<string> executeSomething(
      const Command& cmd,
      const string& prefix,
      function<SemiFuture<string>()> innerFunc);

  static void triggerDequeue(shared_ptr<QueuedParameters> queuedParameters);
  static void onDequeueSuccess(
      const shared_ptr<QueuedParameters>& queuedParameters,
      const QueueEntry& queueEntry,
      const string& result);
  static void onDequeueError(
      const shared_ptr<QueuedParameters>& queuedParameters,
      const QueueEntry& queueEntry,
      const exception_wrapper& e);

 public:
  QueuedCli(
      string id,
      shared_ptr<Cli> cli,
      shared_ptr<Executor> parentExecutor);

  static std::shared_ptr<QueuedCli>
  make(string id, shared_ptr<Cli> cli, shared_ptr<Executor> parentExecutor);

  ~QueuedCli() override;

  folly::SemiFuture<std::string> executeRead(const ReadCommand cmd) override;

  folly::SemiFuture<std::string> executeWrite(const WriteCommand cmd) override;
};
} // namespace devmand::channels::cli
