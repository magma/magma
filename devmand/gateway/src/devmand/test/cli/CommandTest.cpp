// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/channels/cli/Channel.h>
#include <devmand/channels/cli/Cli.h>
#include <devmand/test/cli/utils/Log.h>
#include <devmand/test/cli/utils/MockCli.h>
#include <gtest/gtest.h>

namespace devmand {
namespace test {
namespace cli {

using namespace devmand::test::utils::cli;

class CommandTest : public ::testing::Test {
 protected:
  void SetUp() override {
    devmand::test::utils::log::initLog();
  }
};

TEST_F(CommandTest, api) {
  std::string foo("foo");
  ReadCommand cmd = ReadCommand::create(foo);
  EXPECT_EQ("foo", cmd.raw());
  foo.clear();
  EXPECT_EQ("foo", cmd.raw());
  cmd.raw().clear();
  EXPECT_EQ("foo", cmd.raw());

  const auto mockCli = std::make_shared<EchoCli>();
  folly::SemiFuture<std::string> future = mockCli->executeRead(cmd);
  EXPECT_EQ("foo", std::move(future).get());

  Channel cliChannel("cmdTEst", std::make_shared<EchoCli>());
  folly::SemiFuture<std::string> futureFromChannel =
      cliChannel.executeRead(cmd);
  EXPECT_EQ("foo", std::move(futureFromChannel).get());
}

} // namespace cli
} // namespace test
} // namespace devmand
