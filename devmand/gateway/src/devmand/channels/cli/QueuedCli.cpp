// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <boost/range/size_type.hpp>
#include <devmand/channels/cli/QueuedCli.h>

namespace devmand {
namespace channels {
namespace cli {

using namespace std;
using namespace folly;

shared_ptr<QueuedCli> QueuedCli::make(
    string id,
    shared_ptr<Cli> cli,
    shared_ptr<Executor> parentExecutor,
    const long capacity) {
  return std::make_shared<QueuedCli>(id, cli, parentExecutor, capacity);
}

QueuedCli::QueuedCli(
    string _id,
    shared_ptr<Cli> _cli,
    shared_ptr<Executor> _parentExecutor,
    const long capacity)
    : id(_id),
      sharedCli(_cli),
      sharedParentExecutor(_parentExecutor),
      sharedSerialExecutorKeepAlive(
          make_shared<Executor::KeepAlive<SerialExecutor>>(
              SerialExecutor::create(
                  Executor::getKeepAliveToken(_parentExecutor.get())))) {
  queuedParameters = std::make_shared<QueuedParameters>(
      id, sharedCli, sharedSerialExecutorKeepAlive, capacity);
}

SemiFuture<Unit> QueuedCli::destroy() {
  MLOG(MDEBUG) << "[" << id << "] "
               << "destroy: started";
  // call underlying destroy()
  SemiFuture<Unit> innerDestroy = sharedCli->destroy();

  MLOG(MDEBUG) << "[" << id << "] "
               << "destroy: dequeuing " << queuedParameters->queue->size()
               << " items";
  {
    QueueEntry queueEntry;
    while (queuedParameters->queue->try_dequeue(queueEntry)) {
      MLOG(MDEBUG) << "[" << id << "] (" << queueEntry.command << ") "
                   << "destroy: fulfilling promise with exception";
      queueEntry.promise->setException(runtime_error("QCli: Shutting down"));
    }
  } // drop queueEntry to release queuedParameters from its obtain.. function
  MLOG(MDEBUG) << "[" << id << "] "
               << "destroy: done";
  return innerDestroy;
}

QueuedCli::~QueuedCli() {
  MLOG(MDEBUG) << "[" << id << "] "
               << "~QCli: started";
  destroy().get();
  queuedParameters = nullptr;
  sharedSerialExecutorKeepAlive = nullptr;
  sharedParentExecutor = nullptr;
  sharedCli = nullptr;
  MLOG(MDEBUG) << "[" << id << "] "
               << "~QCli: done";
}

SemiFuture<string> QueuedCli::executeRead(const ReadCommand cmd) {
  boost::recursive_mutex::scoped_lock scoped_lock(writeCommandMutex);
  MLOG(MDEBUG) << "[" << queuedParameters->id << "] "
               << "Invoking read command: \"" << cmd << "\"";
  return executeSomething(cmd, "QCli.executeRead", [cmd](shared_ptr<Cli> cli) {
    return cli->executeRead(cmd);
  });
}

SemiFuture<string> QueuedCli::executeWrite(const WriteCommand cmd) {
  boost::recursive_mutex::scoped_lock scoped_lock(writeCommandMutex);
  MLOG(MDEBUG) << "[" << queuedParameters->id << "] "
               << "Invoking write command: \"" << cmd << "\"";
  Command command = cmd;
  if (!command.isMultiCommand()) {
    // Single line config command, execute with read
    return executeRead(ReadCommand::create(cmd.raw(), true)); //skip cache
  }

  const vector<Command>& commands = command.splitMultiCommand();
  vector<Future<string>> commmandsFutures;

  for (unsigned long i = 0; i < (commands.size() - 1); i++) {
    commmandsFutures.emplace_back(
        executeSomething(
            commands.at(i),
            "QCli.executeWrite",
            [cmd = commands.at(i)](shared_ptr<Cli> cli) {
              return cli->executeWrite(WriteCommand::create(cmd));
            })
            .via(sharedSerialExecutorKeepAlive->get()));
  }

  commmandsFutures.emplace_back(
      executeRead(ReadCommand::create(commands.back().raw(), true)) //skip cache
          .via(sharedSerialExecutorKeepAlive->get()));

  return reduce(
             commmandsFutures.begin(),
             commmandsFutures.end(),
             string(""),
             [](string s1, string s2) { return s1 + s2; })
      .semi();
}

SemiFuture<string> QueuedCli::executeSomething(
    const Command& cmd,
    const string& prefix,
    function<SemiFuture<string>(shared_ptr<Cli>)> innerFunc) {
  shared_ptr<Promise<string>> promise = make_shared<Promise<string>>();
  QueueEntry queueEntry = QueueEntry{move(innerFunc), promise, cmd, prefix};
  MLOG(MDEBUG) << "1[" << queuedParameters->id << "] (" << queueEntry.command
               << ") " << prefix << " adding to queue ('" << cmd << "')";

  bool success = queuedParameters->queue->try_enqueue(move(queueEntry));

  if (!success) {
    promise->setException(runtime_error("QCli: queue is full"));
    return promise->getFuture().semi();
  }

  if (!queuedParameters->isProcessing) {
    triggerDequeue(queuedParameters);
  }
  return promise->getFuture().semi();
}

/*
 * Start queue reading on consumer thread if queue contains new items.
 * It is safe to call this method anytime, it is thread safe.
 */
void QueuedCli::triggerDequeue(shared_ptr<QueuedParameters> queuedParameters) {
  if (shared_ptr<Executor::KeepAlive<SerialExecutor>> serialExecutor =
          queuedParameters->serialExecutorKeepAlive.lock()) {
    // switch to consumer thread
    via(serialExecutor->get(), [params = queuedParameters]() {
      MLOG(MDEBUG) << "[" << params->id << "] "
                   << "QCli.isProcessing:" << params->isProcessing
                   << ", queue size:" << params->queue->size();
      // do nothing if still waiting for remote device to respond
      if (params->isProcessing) {
        // waiting for response
        return;
      }
      QueueEntry queueEntry;
      // try to obtain entry from queue
      if (!params->queue->try_dequeue(queueEntry)) {
        // nothing in the queue
        return;
      }

      if (shared_ptr<Cli> cli = params->cli.lock()) {
        if (shared_ptr<Executor::KeepAlive<SerialExecutor>> serialExecutor2 =
                params->serialExecutorKeepAlive.lock()) {
          startProcessing(move(queueEntry), params, cli, serialExecutor2);
          return;
        }
      }
      // fail fast if qcli is destroyed
      MLOG(MDEBUG) << "[" << params->id << "] (" << queueEntry.command << ") "
                   << queueEntry.loggingPrefix << " Shutting down";
      queueEntry.promise->setException(
          DisconnectedException("SSH session expired"));
    });
  } else {
    MLOG(MDEBUG) << "[" << queuedParameters->id
                 << "] Cannot obtain serial executor";
    // locking weak pointers failed, do nothing as object was destroyed
  }
}

void QueuedCli::startProcessing(
    QueueEntry queueEntry,
    shared_ptr<QueuedParameters> params,
    shared_ptr<Cli> cli,
    shared_ptr<Executor::KeepAlive<SerialExecutor>> executor) {
  params->isProcessing = true;
  SemiFuture<string> cliFuture = queueEntry.obtainFutureFromCli(cli);
  MLOG(MDEBUG) << "[" << params->id << "] (" << queueEntry.command << ") "
               << queueEntry.loggingPrefix
               << " dequeued and cli future obtained";

  move(cliFuture)
      .via(executor->get())
      .thenValue([params, queueEntry](string result) -> Future<Unit> {
        // after cliFuture completes, finish processing on consumer thread
        onDequeueSuccess(params, queueEntry, result);
        return unit;
      })
      .thenError([params, queueEntry](exception_wrapper e) -> Future<Unit> {
        onDequeueError(params, queueEntry, e);
        return unit;
      });
}

void QueuedCli::onDequeueSuccess(
    const shared_ptr<QueuedParameters>& params,
    const QueuedCli::QueueEntry& queueEntry,
    const string& result) {
  MLOG(MDEBUG) << "[" << params->id << "] (" << queueEntry.command << ") "
               << queueEntry.loggingPrefix << " succeeded";
  params->isProcessing = false;
  queueEntry.promise->setValue(result);
  triggerDequeue(params);
}

void QueuedCli::onDequeueError(
    const shared_ptr<QueuedParameters>& params,
    const QueuedCli::QueueEntry& queueEntry,
    const exception_wrapper& e) {
  MLOG(MDEBUG) << "[" << params->id << "] (" << queueEntry.command << ") "
               << queueEntry.loggingPrefix << " failed with exception '"
               << e.what() << "'";
  params->isProcessing = false;
  queueEntry.promise->setException(e);
  triggerDequeue(params);
}

QueuedCli::QueuedParameters::QueuedParameters(
    const string& _id,
    weak_ptr<Cli> _cli,
    weak_ptr<Executor::KeepAlive<SerialExecutor>> _serialExecutorKeepAlive,
    const long capacity)
    : id(_id),
      cli(_cli),
      serialExecutorKeepAlive(_serialExecutorKeepAlive),
      queue(std::make_unique<CliQueue>(capacity)) {}
} // namespace cli
} // namespace channels
} // namespace devmand
