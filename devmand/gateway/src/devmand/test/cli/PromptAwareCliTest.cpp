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

  void SetUp() override {
    devmand::test::utils::log::initLog();
    devmand::test::utils::ssh::initSsh();
    testExec = make_shared<CPUThreadPoolExecutor>(1);
  }

  void TearDown() override {
    MLOG(MDEBUG) << "Waiting for test executor to finish";
    testExec->join();
  }
};

class MockSession : public SessionAsync {
 private:
  int counter = 0;
  shared_ptr<CPUThreadPoolExecutor> testExec;

 public:
  MockSession(shared_ptr<CPUThreadPoolExecutor> _testExec)
      : testExec(_testExec){};

  virtual Future<Unit> write(const string& command) {
    (void)command;
    return via(testExec.get(), [command]() {
      if (command == "error") {
        throw std::runtime_error("error");
      }
      return unit;
    });
  };

  virtual Future<string> read(int timeoutMillis) {
    (void)timeoutMillis;
    counter++;
    return via(testExec.get(), [c = counter]() {
      if (c == 1) {
        return "PROMPT>\nPROMPT>";
      }
      return "";
    });
  };

  virtual Future<string> readUntilOutput(const string& lastOutput) {
    (void)lastOutput;
    return via(testExec.get(), []() { return ""; });
  };
};

static shared_ptr<PromptAwareCli> getCli(
    shared_ptr<CPUThreadPoolExecutor> testExec) {
  return PromptAwareCli::make(
      "test",
      make_shared<MockSession>(testExec),
      CliFlavour::create(""),
      std::make_shared<folly::CPUThreadPoolExecutor>(1),
      make_shared<ThreadWheelTimekeeper>());
}

TEST_F(PromptAwareCliTest, cleanDestructOnSuccess) {
  auto testedCli = getCli(testExec);

  SemiFuture<string> future =
      testedCli->executeRead(ReadCommand::create("returning"));

  // Destruct cli
  testedCli.reset();

  ASSERT_EQ(move(future).via(testExec.get()).get(10s), "");
}

TEST_F(PromptAwareCliTest, cleanDestructOnWriteSuccess) {
  auto testedCli = getCli(testExec);

  SemiFuture<string> future =
      testedCli->executeWrite(WriteCommand::create("returning"));

  // Destruct cli
  testedCli.reset();

  ASSERT_EQ(move(future).via(testExec.get()).get(10s), "returning");
}

TEST_F(PromptAwareCliTest, cleanDestructOnError) {
  auto testedCli = getCli(testExec);

  SemiFuture<string> future =
      testedCli->executeRead(ReadCommand::create("error"));

  // Destruct cli
  testedCli.reset();

  EXPECT_THROW(move(future).via(testExec.get()).get(10s), runtime_error);
}

} // namespace cli
} // namespace test
} // namespace devmand
