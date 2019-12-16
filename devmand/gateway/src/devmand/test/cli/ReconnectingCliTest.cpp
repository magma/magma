// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/cartography/DeviceConfig.h>
#include <devmand/channels/cli/ReconnectingCli.h>
#include <devmand/test/cli/utils/Log.h>
#include <devmand/test/cli/utils/MockCli.h>
#include <devmand/test/cli/utils/Ssh.h>
#include <folly/Executor.h>
#include <folly/executors/CPUThreadPoolExecutor.h>
#include <gtest/gtest.h>
#include <chrono>

namespace devmand {
namespace test {
namespace cli {

using namespace devmand::channels::cli;
using namespace devmand::test::utils::cli;
using namespace std;
using namespace folly;

class ReconnectingCliTest : public ::testing::Test {
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

static shared_ptr<ReconnectingCli> getCli(
    function<SemiFuture<shared_ptr<Cli>>()> createCliStack) {
  shared_ptr<ReconnectingCli> cli = ReconnectingCli::make(
      "test",
      make_shared<CPUThreadPoolExecutor>(1),
      move(createCliStack),
      make_shared<folly::ThreadWheelTimekeeper>(),
      2s);
  return cli;
}

static void waitTillCliRecreated() {
  std::this_thread::sleep_for(2s);
}

TEST_F(ReconnectingCliTest, cleanDestructNotConnected) {
  auto factory = [this]() {
    std::this_thread::sleep_for(2s);
    return SemiFuture<shared_ptr<Cli>>(getMockCli<EchoCli>(0, testExec));
  };
  auto testedCli = getCli(factory);

  SemiFuture<string> future =
      testedCli->executeRead(ReadCommand::create("not returning"));

  // Destruct cli
  testedCli.reset();

  EXPECT_THROW(move(future).via(testExec.get()).get(10s), runtime_error);
}

TEST_F(ReconnectingCliTest, cleanDestructConnected) {
  auto factory = [this]() {
    return SemiFuture<shared_ptr<Cli>>(getMockCli<EchoCli>(0, testExec));
  };
  auto testedCli = getCli(factory);
  waitTillCliRecreated();

  SemiFuture<string> future =
      testedCli->executeRead(ReadCommand::create("returning"));

  // Destruct cli
  testedCli.reset();

  ASSERT_EQ(move(future).via(testExec.get()).get(10s), "returning");
}

TEST_F(ReconnectingCliTest, cleanDestructConnectedWithDelay) {
  auto factory = [this]() {
    return SemiFuture<shared_ptr<Cli>>(getMockCli<EchoCli>(2, testExec));
  };
  auto testedCli = getCli(factory);
  waitTillCliRecreated();

  SemiFuture<string> future =
      testedCli->executeRead(ReadCommand::create("returning"));

  // Destruct cli
  testedCli.reset();

  ASSERT_EQ(move(future).via(testExec.get()).get(10s), "returning");
}

TEST_F(ReconnectingCliTest, cleanDestructOnError) {
  auto factory = [this]() {
    return SemiFuture<shared_ptr<Cli>>(getMockCli<ErrCli>(0, testExec));
  };
  auto testedCli = getCli(factory);
  waitTillCliRecreated();

  SemiFuture<string> future =
      testedCli->executeRead(ReadCommand::create("not returning"));

  // Destruct cli
  testedCli.reset();

  EXPECT_THROW(move(future).via(testExec.get()).get(10s), runtime_error);
}

TEST_F(ReconnectingCliTest, cleanDestructOnErrorWithDelay) {
  auto factory = [this]() {
    return SemiFuture<shared_ptr<Cli>>(getMockCli<ErrCli>(2, testExec));
  };
  auto testedCli = getCli(factory);
  waitTillCliRecreated();

  SemiFuture<string> future =
      testedCli->executeRead(ReadCommand::create("not returning"));

  // Destruct cli
  testedCli.reset();

  EXPECT_THROW(move(future).via(testExec.get()).get(10s), runtime_error);
}

} // namespace cli
} // namespace test
} // namespace devmand
