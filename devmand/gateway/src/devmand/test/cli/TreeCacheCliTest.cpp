#define LOG_WITH_GLOG

#include <magma_logging.h>

#include <devmand/cartography/DeviceConfig.h>
#include <devmand/channels/cli/Cli.h>
#include <devmand/channels/cli/TreeCacheCli.h>
#include <devmand/test/cli/TreeCacheTestData.h>
#include <devmand/test/cli/utils/Log.h>
#include <devmand/test/cli/utils/MockCli.h>
#include <folly/executors/CPUThreadPoolExecutor.h>
#include <folly/executors/ThreadedExecutor.h>
#include <folly/futures/ThreadWheelTimekeeper.h>
#include <gtest/gtest.h>
#include <chrono>

namespace devmand::channels::cli {

using namespace std;
using namespace std::chrono_literals;
using devmand::test::utils::cli::MockedCli;
using folly::CPUThreadPoolExecutor;

static const shared_ptr<CliFlavour> ubiquitiFlavour = CliFlavour::getUbiquiti();
static const char* showRunningCommand = "show running-config";

class TreeCacheCliTest : public ::testing::Test {
 protected:
  shared_ptr<CPUThreadPoolExecutor> testExec =
      make_shared<CPUThreadPoolExecutor>(1);
  shared_ptr<MockedCli> mockedCli;
  shared_ptr<TreeCacheCli> tested_ubiquiti;
  shared_ptr<TreeCache> treeCache;

  void SetUp() override {
    devmand::test::utils::log::initLog();

    map<string, string> mockedResponses;
    mockedResponses.insert(
        make_pair(showRunningCommand, testdata::SH_RUN_UBIQUITI));
    mockedCli = make_shared<MockedCli>(mockedResponses);

    treeCache = make_shared<TreeCache>(ubiquitiFlavour);

    tested_ubiquiti = make_shared<TreeCacheCli>(
        "id", mockedCli, testExec, ubiquitiFlavour, treeCache);
    EXPECT_TRUE(treeCache->isEmpty());
  }

  void TearDown() override {
    MLOG(MDEBUG) << "Waiting for test executor to finish";
    testExec->join();
  }
};

TEST_F(TreeCacheCliTest, checkMockWorks) {
  // just load show running config, should be handled by MockedCli
  string shResult =
      tested_ubiquiti
          ->executeRead(ReadCommand::create(showRunningCommand, false))
          .get();
  EXPECT_EQ(testdata::SH_RUN_UBIQUITI, shResult);
  mockedCli->clear();
  EXPECT_THROW(
      tested_ubiquiti
          ->executeRead(ReadCommand::create(showRunningCommand, false))
          .get(),
      runtime_error);
}

TEST_F(TreeCacheCliTest, getParticularIfc_samePass) {
  // this should execute base command on MockedCli, output will be parsed
  string result = tested_ubiquiti
                      ->executeRead(ReadCommand::create(
                          "show running-config interface 0/14", false))
                      .get();
  EXPECT_EQ(testdata::SH_RUN_INT_GI4, result);
}

TEST_F(
    TreeCacheCliTest,
    getParticularIfc_populateCache_thenGoStraightToTreeCache) {
  // first request will be passed to MockedCli
  string result = tested_ubiquiti
                      ->executeRead(ReadCommand::create(
                          "show running-config interface 0/14", false))
                      .get();
  EXPECT_FALSE(treeCache->isEmpty());
  EXPECT_EQ(testdata::SH_RUN_INT_GI4, result);
  // second request should end in tree cache
  mockedCli->clear();
  result = tested_ubiquiti
               ->executeRead(ReadCommand::create(
                   "show running-config interface 0/14", false))
               .get();
  EXPECT_EQ(testdata::SH_RUN_INT_GI4, result);
}
} // namespace devmand::channels::cli
