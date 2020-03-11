// Copyright (c) 2019-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#define LOG_WITH_GLOG
#include <magma_logging.h>

#include <devmand/devices/cli/schema/Path.h>
#include <devmand/test/cli/utils/Log.h>
#include <gtest/gtest.h>

namespace devmand {
namespace test {
namespace cli {

using namespace devmand::devices::cli;
using namespace folly;

class PathTest : public ::testing::Test {
 protected:
  void SetUp() override {
    devmand::test::utils::log::initLog();
  }
};

TEST_F(PathTest, path) {
  Path root = "/";
  ASSERT_EQ(root.getDepth(), 0);
  ASSERT_EQ(root.getSegments(), vector<string>());
  ASSERT_EQ(root.getChild("abcd"), "/abcd");

  Path ifcs = "/openconfig-interfaces:interfaces";
  ASSERT_EQ(ifcs.getDepth(), 1);
  ASSERT_TRUE(ifcs.isChildOf(root));
  ASSERT_FALSE(ifcs.isChildOf("/something-else"));
  ASSERT_FALSE(root.isChildOf(ifcs));
  ASSERT_EQ(
      ifcs.getSegments(), vector<string>{"openconfig-interfaces:interfaces"});
  ASSERT_EQ(ifcs.getLastSegment(), "openconfig-interfaces:interfaces");
  ASSERT_EQ(ifcs.getParent(), root);

  Path ifc =
      "/openconfig-interfaces:interfaces/interface[id='ethernet 0/1']/config";
  ASSERT_EQ(ifc.getDepth(), 3);
  ASSERT_TRUE(ifc.isChildOf(root));
  ASSERT_TRUE(ifc.isChildOf(ifcs));
  ASSERT_EQ(
      ifc.unkeyed(), "/openconfig-interfaces:interfaces/interface/config");
  ASSERT_EQ(
      ifc.getParent(),
      "/openconfig-interfaces:interfaces/interface[id='ethernet 0/1']");
  vector<string> expectedSegments =
      vector<string>{"openconfig-interfaces:interfaces",
                     "interface[id='ethernet 0/1']",
                     "config"};
  ASSERT_EQ(ifc.getSegments(), expectedSegments);
  Path::Keys expectedKeys = dynamic::object("id", "ethernet 0/1");
  ASSERT_EQ(ifc.getKeysFromSegment("interface"), expectedKeys);
  Path::Keys expectedEmptyKeys = dynamic::object();
  ASSERT_EQ(ifc.getKeysFromSegment("config"), expectedEmptyKeys);

  Path ip =
      R"(/openconfig-interfaces:interfaces/interface[id='ethernet 0/1']/subinterfaces/subinterface[index=0]/openconfig-if-ip:ip/ipv4/address[ip='4:4:4:4'])";
  ASSERT_EQ(ip.unkeyed().getSegments().size(), 7);
  ASSERT_EQ(ip.getSegments().size(), ip.getDepth());
  ASSERT_EQ(
      ip.prefixAllSegments().str(),
      R"(/openconfig-interfaces:interfaces/openconfig-interfaces:interface[id='ethernet 0/1']/openconfig-interfaces:subinterfaces/openconfig-interfaces:subinterface[index=0]/openconfig-if-ip:ip/openconfig-if-ip:ipv4/openconfig-if-ip:address[ip='4:4:4:4'])");
}

TEST_F(PathTest, invalidPath) {
  EXPECT_THROW(Path(""), InvalidPathException);
  EXPECT_THROW(Path("openconfig-interfaces:interfaces"), InvalidPathException);
  EXPECT_THROW(
      Path("/openconfig-interfaces:interfaces").getChild(""),
      InvalidPathException);
  EXPECT_THROW(
      Path("/openconfig-interfaces:interfaces").getChild("/abcd"),
      InvalidPathException);
}

TEST_F(PathTest, segmentKeys) {
  EXPECT_EQ(
      Path("/openconfig-interfaces:interfaces/interface[name='0/85']/state")
          .getKeysFromSegment("interface")
          .size(),
      1);
}

} // namespace cli
} // namespace test
} // namespace devmand
