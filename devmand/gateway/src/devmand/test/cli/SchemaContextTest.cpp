// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/devices/cli/schema/Path.h>
#include <devmand/devices/cli/schema/SchemaContext.h>
#include <devmand/test/cli/utils/Log.h>
#include <gtest/gtest.h>

namespace devmand::test::cli {

using namespace devmand::devices::cli;
using devmand::devices::cli::SchemaContext;
using namespace folly;
using std::string;

class SchemaContextTest : public ::testing::Test {
 public:
  SchemaContextTest() : context(Model::OPENCONFIG_0_1_6){};

 protected:
  void SetUp() override {
    devmand::test::utils::log::initLog();
  }

  SchemaContext context;
};

TEST_F(SchemaContextTest, getKeys) {
  Path pathToNode = Path(
      "/openconfig-interfaces:interfaces/interface/subinterfaces/subinterface/openconfig-if-ip:ipv4/addresses/address");
  EXPECT_TRUE(context.isList(pathToNode));
  const vector<string>& vector = context.getKeys(pathToNode);
  EXPECT_EQ(1, vector.size());
  EXPECT_EQ("ip", vector.at(0));
}

TEST_F(SchemaContextTest, isList) {
  Path pathToList = Path("/openconfig-interfaces:interfaces/interface");
  Path pathToNonList =
      Path("/openconfig-interfaces:interfaces/interface/state");
  EXPECT_TRUE(context.isList(pathToList));
  EXPECT_FALSE(context.isList(pathToNonList));
}

TEST_F(SchemaContextTest, isConfig) {
  Path cfg = Path("/openconfig-interfaces:interfaces/interface/config");
  Path state = Path("/openconfig-interfaces:interfaces/interface/state");
  EXPECT_TRUE(context.isConfig(cfg));
  EXPECT_FALSE(context.isConfig(state));
}

TEST_F(SchemaContextTest, checkValidPath) {
  Path pathToRealNode = Path(
      "/openconfig-interfaces:interfaces/interface/subinterfaces/subinterface");
  Path pathToNonExistentNode = Path(
      "/openconfig-interfaces:interfaces/interface/sssubinterfaces/subinterface");
  EXPECT_TRUE(context.isPathValid(pathToRealNode));
  EXPECT_FALSE(context.isPathValid(pathToNonExistentNode));
}

} // namespace devmand::test::cli
