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
#include <folly/concurrency/DynamicBoundedQueue.h>
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
  string id;
  shared_ptr<Cli> sharedCli;
  shared_ptr<Executor> sharedParentExecutor;
  shared_ptr<Executor::KeepAlive<SerialExecutor>>
      sharedSerialExecutorKeepAlive; // maintain consumer thread
  recursive_mutex writeCommandMutex; // used to block executeRead while
  // unpacking commands from executeWrite

  struct QueueEntry {
    function<SemiFuture<string>(shared_ptr<Cli> cli)> obtainFutureFromCli;
    shared_ptr<Promise<string>> promise;
    Command command = ReadCommand::create("dummy");
    string loggingPrefix;
  };

  using CliQueue = DynamicBoundedQueue<QueueEntry, false, true, false>;

  struct QueuedParameters {
    string id;
    weak_ptr<Cli> cli;
    weak_ptr<Executor::KeepAlive<SerialExecutor>> serialExecutorKeepAlive;

    /**
     * Bounded multi producer single consumer queue where consumer is not
     * blocked on dequeue.
     */
    std::unique_ptr<CliQueue>
        queue; // TODO: investigate priority queue for keepalive commands

    atomic<bool> isProcessing = ATOMIC_VAR_INIT(false);

    QueuedParameters(
        const string& id,
        weak_ptr<Cli> cli,
        weak_ptr<Executor::KeepAlive<SerialExecutor>> serialExecutorKeepAlive,
        const long capacity);
  };

  shared_ptr<QueuedParameters> queuedParameters;

  SemiFuture<string> executeSomething(
      const Command& cmd,
      const string& prefix,
      function<SemiFuture<string>(shared_ptr<Cli>)> innerFunc);

  static void triggerDequeue(shared_ptr<QueuedParameters> queuedParameters);

  static void startProcessing(
      QueueEntry entry,
      shared_ptr<QueuedParameters> params,
      shared_ptr<Cli> cli,
      shared_ptr<Executor::KeepAlive<SerialExecutor>> executor);

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
      shared_ptr<Executor> parentExecutor,
      const long capacity);

  static std::shared_ptr<QueuedCli> make(
      string id,
      shared_ptr<Cli> cli,
      shared_ptr<Executor> parentExecutor,
      const long capacity = 100000);

  SemiFuture<Unit> destroy() override;

  ~QueuedCli() override;

  SemiFuture<std::string> executeRead(const ReadCommand cmd) override;

  SemiFuture<std::string> executeWrite(const WriteCommand cmd) override;
};
} // namespace devmand::channels::cli
