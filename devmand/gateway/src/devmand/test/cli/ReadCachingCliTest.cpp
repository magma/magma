// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/cartography/DeviceConfig.h>
#include <devmand/channels/cli/ReadCachingCli.h>
#include <devmand/test/cli/utils/Log.h>
#include <devmand/test/cli/utils/MockCli.h>
#include <devmand/test/cli/utils/Ssh.h>
#include <folly/Executor.h>
#include <folly/executors/CPUThreadPoolExecutor.h>
#include <folly/executors/IOThreadPoolExecutor.h>
#include <gtest/gtest.h>
#include <chrono>

namespace devmand {
namespace test {
namespace cli {

using namespace devmand::channels::cli;
using namespace devmand::test::utils::cli;
using namespace std;
using namespace folly;

class ReadCachingCliTest : public ::testing::Test {
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

template <typename NESTED>
static shared_ptr<AsyncCli> getMockCli(
    uint delay,
    shared_ptr<CPUThreadPoolExecutor> exec) {
  vector<unsigned int> durations = {delay};
  return make_shared<AsyncCli>(make_shared<NESTED>(), exec, durations);
}

static shared_ptr<ReadCachingCli> getCli(shared_ptr<Cli> delegate) {
  return make_shared<ReadCachingCli>(
      "test",
      delegate,
      ReadCachingCli::createCache(),
      std::make_shared<folly::IOThreadPoolExecutor>(
          1, std::make_shared<folly::NamedThreadFactory>("rccli")));
}

TEST_F(ReadCachingCliTest, cleanDestructOnSuccess) {
  shared_ptr<AsyncCli> delegate = getMockCli<EchoCli>(0, testExec);
  auto testedCli = getCli(delegate);

  SemiFuture<string> future =
      testedCli->executeRead(ReadCommand::create("returning"));

  // Destruct cli
  testedCli.reset();

  ASSERT_EQ(move(future).via(testExec.get()).get(10s), "returning");
}

TEST_F(ReadCachingCliTest, cleanDestructOnError) {
  shared_ptr<AsyncCli> delegate = getMockCli<ErrCli>(0, testExec);
  auto testedCli = getCli(delegate);

  SemiFuture<string> future =
      testedCli->executeRead(ReadCommand::create("error"));

  // Destruct cli
  testedCli.reset();

  EXPECT_THROW(move(future).via(testExec.get()).get(10s), runtime_error);
}

} // namespace cli
} // namespace test
} // namespace devmand
