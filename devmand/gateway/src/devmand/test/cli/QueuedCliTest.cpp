// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <boost/algorithm/string/trim.hpp>
#include <devmand/cartography/DeviceConfig.h>
#include <devmand/channels/cli/Cli.h>
#include <devmand/channels/cli/QueuedCli.h>
#include <devmand/test/cli/utils/Log.h>
#include <devmand/test/cli/utils/MockCli.h>
#include <devmand/test/cli/utils/Ssh.h>
#include <folly/Executor.h>
#include <folly/executors/CPUThreadPoolExecutor.h>
#include <folly/futures/Future.h>
#include <folly/futures/ThreadWheelTimekeeper.h>
#include <gtest/gtest.h>
#include <chrono>

namespace devmand {
namespace test {
namespace cli {

using namespace devmand::channels::cli;
using namespace devmand::test::utils::cli;
using namespace std;

shared_ptr<CPUThreadPoolExecutor> executor =
    make_shared<CPUThreadPoolExecutor>(8);

class QueuedCliTest : public ::testing::Test {
 protected:
  void SetUp() override {
    devmand::test::utils::log::initLog(MDEBUG);
    devmand::test::utils::ssh::initSsh();
  }
};

TEST_F(QueuedCliTest, queuedCli) {
  vector<unsigned int> durations = {2};
  shared_ptr<QueuedCli> cli = QueuedCli::make(
      "testConnection",
      make_shared<AsyncCli>(make_shared<EchoCli>(), executor, durations),
      executor);

  vector<string> results{"one", "two", "three", "four", "five", "six"};

  // create requests
  vector<Command> cmds;
  for (const auto& str : results) {
    cmds.push_back(ReadCommand::create(str));
  }

  // send requests
  vector<folly::SemiFuture<string>> futures;
  for (const auto& cmd : cmds) {
    MLOG(MDEBUG) << "Executing command '" << cmd;
    futures.push_back(cli->executeRead(ReadCommand::create(cmd)));
  }

  // collect values
  const vector<folly::Try<string>>& values =
      collectAll(futures.begin(), futures.end()).get();

  // check values
  EXPECT_EQ(values.size(), results.size());
  for (unsigned int i = 0; i < values.size(); ++i) {
    EXPECT_EQ(boost::algorithm::trim_copy(values[i].value()), results[i]);
  }
}

TEST_F(QueuedCliTest, queueFullTest) {
  unsigned int parallelThreads = 1;
  shared_ptr<CPUThreadPoolExecutor> queuedCliParallelExecutor =
      make_shared<CPUThreadPoolExecutor>(parallelThreads);
  WriteCommand cmd = WriteCommand::create("a\nb\nc", true);
  const shared_ptr<QueuedCli>& cli = QueuedCli::make(
      "testFullCapacity", make_shared<EchoCli>(), executor, 1); // capacity at 1

  vector<Future<string>> queuedFutures;

  queuedFutures.emplace_back(folly::via(
      queuedCliParallelExecutor.get(),
      [cli, cmd]() { return cli->executeWrite(cmd); }));

  EXPECT_THROW(
      collect(queuedFutures.begin(), queuedFutures.end()).get(), runtime_error);
}

TEST_F(QueuedCliTest, queueOrderingTest) {
  unsigned int iterations = 200;
  unsigned int parallelThreads = 165;
  shared_ptr<CPUThreadPoolExecutor> queuedCliParallelExecutor =
      make_shared<CPUThreadPoolExecutor>(parallelThreads);
  WriteCommand cmd = WriteCommand::create("1\n2\n3\n4\n5\n6\n7\n8\n9", true);
  const shared_ptr<QueuedCli>& cli =
      QueuedCli::make("testOrder", make_shared<EchoCli>(), executor);

  vector<Future<string>> queuedFutures;

  for (unsigned long i = 0; i < iterations; i++) {
    queuedFutures.emplace_back(folly::via(
        queuedCliParallelExecutor.get(),
        [cli, cmd]() { return cli->executeWrite(cmd); }));
  }

  vector<string> output =
      collect(queuedFutures.begin(), queuedFutures.end()).get();

  for (unsigned long i = 0; i < iterations; i++) {
    EXPECT_EQ(output[i], "123456789");
  }
}

TEST_F(QueuedCliTest, queuedCliMT) {
  const int loopcount = 10;
  vector<unsigned int> durations = {1};

  shared_ptr<QueuedCli> cli = QueuedCli::make(
      "testConnection",
      make_shared<AsyncCli>(make_shared<EchoCli>(), executor, durations),
      executor);

  // create requests
  vector<folly::Future<string>> futures;
  ReadCommand cmd = ReadCommand::create("hello");
  for (int i = 0; i < loopcount; ++i) {
    MLOG(MDEBUG) << "Executing command '" << cmd;
    futures.push_back(
        folly::via(executor.get(), [&]() { return cli->executeRead(cmd); }));
  }

  // collect values
  const vector<folly::Try<string>>& values =
      collectAll(futures.begin(), futures.end()).get();

  // check values
  EXPECT_EQ(values.size(), loopcount);
  for (auto v : values) {
    EXPECT_EQ(boost::algorithm::trim_copy(v.value()), "hello");
  }
}

TEST_F(QueuedCliTest, threadSafety) {
  int iterations = 1000;
  shared_ptr<CPUThreadPoolExecutor> testExec =
      make_shared<CPUThreadPoolExecutor>(32, 1, iterations);

  shared_ptr<QueuedCli> cli =
      QueuedCli::make("testConnection", make_shared<EchoCli>(), executor);

  vector<Future<string>> execs;
  for (int i = 0; i < iterations; i++) {
    Future<string> future = via(testExec.get(), [cli, i]() {
      return cli->executeRead(ReadCommand::create(to_string(i))).get();
    });
    execs.push_back(std::move(future));
  }

  for (uint i = 0; i < execs.size(); i++) {
    const string& basicString = std::move(execs.at(i)).get();
  }
}

TEST_F(QueuedCliTest, cleanDestructOnSuccess) {
  auto testExec = make_shared<CPUThreadPoolExecutor>(1);
  auto delegate = getMockCli<EchoCli>(3, testExec);
  auto testedCli = QueuedCli::make(
      "testConnection", delegate, make_shared<CPUThreadPoolExecutor>(1));

  vector<SemiFuture<string>> futures;
  for (int i = 0; i < 10; i++) {
    futures.emplace_back(
        testedCli->executeRead(ReadCommand::create("command")));
  }

  // Destruct cli
  testedCli.reset();

  // First request can succeed, other are canceled
  try {
    ASSERT_EQ(move(futures.at(0)).via(testExec.get()).get(10s), "command");
  } catch (const runtime_error& e) {
    // The first one can succeed or be canceled, depends on timing
    // But since we are only testing clean destruct, we don't care
  }

  for (uint i = 1; i < 10; i++) {
    EXPECT_THROW(
        move(futures.at(i)).via(testExec.get()).get(10s), runtime_error);
  }

  MLOG(MDEBUG) << "Waiting for test executor to finish";
  testExec->join();
}

class CustomErr : public runtime_error {
 public:
  CustomErr(string msg) : runtime_error(msg){};
};

class CustomErrCli : public Cli {
 public:
  SemiFuture<folly::Unit> destroy() override {
    return folly::makeSemiFuture(unit);
  }

  folly::SemiFuture<std::string> executeRead(const ReadCommand cmd) override {
    return folly::Future<string>(CustomErr(Command::escape(cmd.raw())));
  }

  folly::SemiFuture<std::string> executeWrite(const WriteCommand cmd) override {
    return folly::Future<string>(CustomErr(Command::escape(cmd.raw())));
  }
};

TEST_F(QueuedCliTest, preserveExType) {
  auto testExec = make_shared<CPUThreadPoolExecutor>(1);
  auto delegate = getMockCli<CustomErrCli>(0, testExec);
  auto testedCli = QueuedCli::make(
      "testConnection", delegate, make_shared<CPUThreadPoolExecutor>(1));
  try {
    testedCli->executeRead(ReadCommand::create("command"))
        .via(testExec.get())
        .get(10s);
    FAIL() << "No exception thrown";
  } catch (const CustomErr& e) {
    // Proper ex type caught
  } catch (const exception& e) {
    FAIL() << "Wrong exception thrown " << e.what();
  }

  MLOG(MDEBUG) << "Waiting for test executor to finish";
  testExec->join();
}

} // namespace cli
} // namespace test
} // namespace devmand
