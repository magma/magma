// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <gtest/gtest.h>

#include <devmand/syslog/Manager.h>

namespace devmand {
namespace test {

TEST(SyslogManagerTest, addAndLookup) {
  syslog::Manager sylogManager;
  sylogManager.addIdentifier("1.1.1.1", "foo");
  EXPECT_EQ("foo", sylogManager.lookup("1.1.1.1"));
}

TEST(SyslogManagerTest, addRemAndLookup) {
  syslog::Manager sylogManager;
  sylogManager.addIdentifier("1.1.1.1", "foo");
  sylogManager.removeIdentifier("1.1.1.1", "foo");
  EXPECT_EQ("", sylogManager.lookup("1.1.1.1"));
}

TEST(SyslogManagerTest, addDupAndLookup) {
  syslog::Manager sylogManager;
  sylogManager.addIdentifier("1.1.1.1", "foo");
  sylogManager.addIdentifier("2.2.2.2", "foo");
  EXPECT_EQ("foo", sylogManager.lookup("1.1.1.1"));
}

TEST(SyslogManagerTest, addDupDiffAndLookup) {
  syslog::Manager sylogManager;
  sylogManager.addIdentifier("1.1.1.1", "foo");
  sylogManager.addIdentifier("2.2.2.2", "bar");
  EXPECT_EQ("foo", sylogManager.lookup("1.1.1.1"));
}

TEST(SyslogManagerTest, addDupIdenAndLookup) {
  syslog::Manager sylogManager;
  sylogManager.addIdentifier("1.1.1.1", "foo");
  sylogManager.addIdentifier("1.1.1.1", "foo");
  EXPECT_EQ("foo", sylogManager.lookup("1.1.1.1"));
}

TEST(SyslogManagerTest, addDupIdenDiffAndLookup) {
  syslog::Manager sylogManager;
  sylogManager.addIdentifier("1.1.1.1", "foo");
  sylogManager.addIdentifier("1.1.1.1", "bar");
  EXPECT_EQ("foo", sylogManager.lookup("1.1.1.1"));
}

} // namespace test
} // namespace devmand
