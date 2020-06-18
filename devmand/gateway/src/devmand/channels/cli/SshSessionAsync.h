// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <boost/lockfree/policies.hpp>
#include <boost/lockfree/spsc_queue.hpp>
#include <boost/thread/mutex.hpp>
#include <devmand/channels/cli/Command.h>
#include <devmand/channels/cli/SshSession.h>
#include <event2/event.h>
#include <folly/executors/IOExecutor.h>
#include <folly/executors/SerialExecutor.h>
#include <folly/futures/Future.h>

namespace devmand {
namespace channels {
namespace cli {
namespace sshsession {

using boost::lockfree::capacity;
using boost::lockfree::spsc_queue;
using devmand::channels::cli::sshsession::SshSession;
using folly::Future;
using folly::IOExecutor;
using folly::makeFuture;
using folly::Promise;
using folly::SerialExecutor;
using folly::Timekeeper;
using folly::Unit;
using folly::via;
using std::condition_variable;
using std::shared_ptr;
using std::string;
using std::unique_ptr;

void readCallback(evutil_socket_t fd, short, void* ptr);

class SessionAsync {
 public:
  virtual Future<Unit> write(const string& command) = 0;
  virtual Future<string> read() = 0;
  virtual Future<string> readUntilOutput(const string& lastOutput) = 0;
  virtual folly::SemiFuture<folly::Unit> destroy() = 0;
};

class SshSessionAsync : public std::enable_shared_from_this<SshSessionAsync>,
                        public SessionAsync {
 private:
  string id;
  shared_ptr<Timekeeper> timekeeper;
  shared_ptr<folly::Executor> executor;
  folly::Executor::KeepAlive<SerialExecutor> serialExecutor;
  unique_ptr<SshSession> session;
  event* sessionEvent = nullptr;
  spsc_queue<string, capacity<200>> readQueue;
  boost::mutex readingMutex;
  std::atomic_bool shutdown;
  std::atomic_bool callbackFinished;
  std::atomic_bool matchingExpectedOutput;
  struct ReadingState {
    shared_ptr<Promise<string>> promise;
    string currentLastOutput;
    string outputSoFar;
  } readingState;

  Future<Unit> waitForCallbacks();

  Future<Unit> setReadingFlag();

 public:
  explicit SshSessionAsync(
      string _id,
      shared_ptr<folly::Executor> _executor,
      shared_ptr<Timekeeper> timekeeper);

  folly::SemiFuture<Unit> destroy() override;

  ~SshSessionAsync();

  Future<Unit> write(const string& command) override;
  Future<string> read() override;
  Future<string> readUntilOutput(const string& lastOutput) override;

  Future<Unit> openShell(
      const string& ip,
      int port,
      const string& username,
      const string& password,
      const long sshConnectionTimeout = 10);

  void setEvent(event*);
  void readSshDataToBuffer();
  void processDataInBuffer();
  socket_t getSshFd();
  void matchExpectedOutput();
  static void failCurrentRead(runtime_error e, shared_ptr<Promise<string>> ptr);
  void unregisterEvent();
};

} // namespace sshsession
} // namespace cli
} // namespace channels
} // namespace devmand
