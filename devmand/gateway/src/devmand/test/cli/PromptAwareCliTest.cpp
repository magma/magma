// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/cartography/DeviceConfig.h>
#include <devmand/channels/cli/Cli.h>
#include <devmand/channels/cli/PromptAwareCli.h>
#include <devmand/test/cli/utils/Log.h>
#include <devmand/test/cli/utils/Ssh.h>
#include <folly/Executor.h>
#include <folly/executors/CPUThreadPoolExecutor.h>
#include <folly/futures/ThreadWheelTimekeeper.h>
#include <gtest/gtest.h>
#include <chrono>

namespace devmand {
namespace test {
namespace cli {

using namespace devmand::channels::cli;
using namespace std;
using namespace folly;
using namespace std::chrono_literals;

class PromptAwareCliTest : public ::testing::Test {
 protected:
  shared_ptr<CPUThreadPoolExecutor> testExec;
  shared_ptr<folly::ThreadWheelTimekeeper> timekeeper;

  void SetUp() override {
    devmand::test::utils::log::initLog();
    devmand::test::utils::ssh::initSsh();
    testExec = make_shared<CPUThreadPoolExecutor>(1);
    timekeeper = make_shared<ThreadWheelTimekeeper>();
  }

  void TearDown() override {
    MLOG(MDEBUG) << "Waiting for test executor to finish";
    testExec->join();
  }
};

class MockSession : public SessionAsync {
 private:
  int counter = 0;
  chrono::milliseconds delay = 1s;
  shared_ptr<Executor> testExec;
  shared_ptr<Timekeeper> timekeeper;

 public:
  MockSession(
      shared_ptr<Executor> _testExec,
      shared_ptr<Timekeeper> _timekeeper)
      : testExec(_testExec), timekeeper(_timekeeper){};

  SemiFuture<folly::Unit> destroy() override {
    return folly::makeSemiFuture(unit);
  }

  virtual Future<Unit> write(const string& command) {
    return via(testExec.get())
        .delayed(delay, timekeeper.get())
        .thenValue([command](auto) {
          if (command == "error") {
            throw std::runtime_error("error");
          }
          return unit;
        });
  };

  virtual Future<string> read() {
    counter++;
    return via(testExec.get())
        .delayed(delay, timekeeper.get())
        .thenValue([c = counter](auto) {
          if (c == 1) {
            return "PROMPT>\nPROMPT>";
          }
          return "";
        });
  };

  virtual Future<string> readUntilOutput(const string& lastOutput) {
    (void)lastOutput;
    return via(testExec.get())
        .delayed(delay, timekeeper.get())
        .thenValue([](auto) { return ""; });
  };
};

static shared_ptr<PromptAwareCli> getCli(
    shared_ptr<CPUThreadPoolExecutor> testExec,
    shared_ptr<folly::Timekeeper> timekeeper) {
  shared_ptr<MockSession> session =
      make_shared<MockSession>(testExec, timekeeper);
  return PromptAwareCli::make(
      "test",
      session,
      CliFlavour::create(""),
      std::make_shared<folly::CPUThreadPoolExecutor>(1),
      timekeeper);
}

TEST_F(PromptAwareCliTest, shortCircutReadAfterDestructing) {
  auto testedCli = getCli(testExec, timekeeper);

  SemiFuture<string> future =
      testedCli->executeRead(ReadCommand::create("returning"));

  // Destruct cli
  testedCli.reset();

  EXPECT_THROW(
      move(future).via(testExec.get()).get(10s), DisconnectedException);
}

TEST_F(PromptAwareCliTest, shortCircutWriteAfterDestructing) {
  auto testedCli = getCli(testExec, timekeeper);

  SemiFuture<string> future =
      testedCli->executeWrite(WriteCommand::create("returning"));

  // Destruct cli
  testedCli.reset();

  EXPECT_THROW(
      move(future).via(testExec.get()).get(10s), DisconnectedException);
}

TEST_F(PromptAwareCliTest, propagateErrorAfterDestructing) {
  auto testedCli = getCli(testExec, timekeeper);

  SemiFuture<string> future =
      testedCli->executeRead(ReadCommand::create("error"));

  // Destruct cli
  testedCli.reset();

  EXPECT_THROW(move(future).via(testExec.get()).get(10s), runtime_error);
}

} // namespace cli
} // namespace test
} // namespace devmand
