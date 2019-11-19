// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <gtest/gtest.h>

#include <devmand/DeviceIdentifiers.h>

namespace devmand {
namespace test {

TEST(DeviceIdentifiersTest, addAndLookup) {
  DeviceIdentifiers dis;
  dis.addIdentifier("1.1.1.1", "foo");
  EXPECT_EQ("foo", dis.lookup("1.1.1.1"));
}

TEST(DeviceIdentifiersTest, addRemAndLookup) {
  DeviceIdentifiers dis;
  dis.addIdentifier("1.1.1.1", "foo");
  dis.removeIdentifier("1.1.1.1", "foo");
  EXPECT_EQ("", dis.lookup("1.1.1.1"));
}

TEST(DeviceIdentifiersTest, addDupAndLookup) {
  DeviceIdentifiers dis;
  dis.addIdentifier("1.1.1.1", "foo");
  dis.addIdentifier("2.2.2.2", "foo");
  EXPECT_EQ("foo", dis.lookup("1.1.1.1"));
}

TEST(DeviceIdentifiersTest, addDupDiffAndLookup) {
  DeviceIdentifiers dis;
  dis.addIdentifier("1.1.1.1", "foo");
  dis.addIdentifier("2.2.2.2", "bar");
  EXPECT_EQ("foo", dis.lookup("1.1.1.1"));
}

TEST(DeviceIdentifiersTest, addDupIdenAndLookup) {
  DeviceIdentifiers dis;
  dis.addIdentifier("1.1.1.1", "foo");
  dis.addIdentifier("1.1.1.1", "foo");
  EXPECT_EQ("foo", dis.lookup("1.1.1.1"));
}

TEST(DeviceIdentifiersTest, addDupIdenDiffAndLookup) {
  DeviceIdentifiers dis;
  dis.addIdentifier("1.1.1.1", "foo");
  dis.addIdentifier("1.1.1.1", "bar");
  EXPECT_EQ("foo", dis.lookup("1.1.1.1"));
}

} // namespace test
} // namespace devmand
