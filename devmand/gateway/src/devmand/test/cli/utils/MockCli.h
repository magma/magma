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
using namespace folly;

class EchoCli : public Cli {
 public:
  ~EchoCli() {
    MLOG(MDEBUG) << "~EchoCli";
  }

  SemiFuture<Unit> destroy() override {
    return makeSemiFuture(unit);
  }

  SemiFuture<std::string> executeRead(const ReadCommand cmd) override {
    return Future<string>(cmd.raw());
  }

  SemiFuture<std::string> executeWrite(const WriteCommand cmd) override {
    return Future<string>(cmd.raw());
  }
};

class ErrCli : public Cli {
 public:
  ~ErrCli() {
    MLOG(MDEBUG) << "~ErrCli";
  }

  SemiFuture<Unit> destroy() override {
    return makeSemiFuture(unit);
  }

  SemiFuture<std::string> executeRead(const ReadCommand cmd) override {
    return Future<string>(runtime_error(cmd.raw()));
  }

  SemiFuture<std::string> executeWrite(const WriteCommand cmd) override {
    return Future<string>(runtime_error(cmd.raw()));
  }
};

class AsyncCli : public Cli {
 public:
  AsyncCli(
      shared_ptr<Cli> _cli,
      shared_ptr<CPUThreadPoolExecutor> _executor,
      vector<unsigned int> _durations)
      : cli(_cli), executor(_executor), durations(_durations), index(0) {}

  ~AsyncCli() {
    MLOG(MDEBUG) << "~AsyncCli";
  }

  SemiFuture<Unit> destroy() override {
    return makeSemiFuture(unit);
  }

  SemiFuture<std::string> executeRead(const ReadCommand cmd) override {
    Future<string> f = via(executor.get()).thenValue([=](...) {
      unsigned int tis = durations[(index++) % durations.size()];
      MLOG(MDEBUG) << "Sleeping for " << tis << "s";
      this_thread::sleep_for(chrono::seconds(tis));
      MLOG(MDEBUG) << "Sleeping done";
      return cli->executeRead(cmd);
    });
    return f;
  }

  SemiFuture<std::string> executeWrite(const WriteCommand cmd) override {
    (void)cmd;
    return Future<string>(runtime_error("Unsupported"));
  }

 protected:
  shared_ptr<Cli> cli; // underlying cli layers
  shared_ptr<CPUThreadPoolExecutor> executor;
  vector<unsigned int> durations;
  unsigned int index;
  bool quit;
};

template <typename NESTED>
shared_ptr<AsyncCli> getMockCli(
    uint delay,
    shared_ptr<CPUThreadPoolExecutor> exec) {
  vector<unsigned int> durations = {delay};
  return make_shared<AsyncCli>(make_shared<NESTED>(), exec, durations);
}

} // namespace cli
} // namespace utils
} // namespace test
} // namespace devmand
