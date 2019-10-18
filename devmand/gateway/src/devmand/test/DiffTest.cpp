// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <gtest/gtest.h>

#include <devmand/Diff.h>

namespace devmand {
namespace test {

template <class Type>
static void diffAddDelete(
    const std::set<Type>& s1,
    const std::set<Type>& s2,
    std::set<Type>& added,
    std::set<Type>& deleted) {
  DiffEventHandler<Type> deh =
    [&added, &deleted](DiffEvent de, const Type& val) {
      switch (de) {
        case DiffEvent::Add:
          added.insert(val);
          break;
        case DiffEvent::Delete:
          deleted.insert(val);
          break;
        case DiffEvent::Modify:
          FAIL() << "This case can't be reached.";
          break;
      }
    };
  diff(s1, s2, deh);
}

TEST(DiffTest, IntegerSets) {
  std::set<int> s1 = {2,5,8,4,6};
  std::set<int> s2 = {9,7,4,2,8,1};
  std::set<int> expectedAdded = {9,7,1};
  std::set<int> expectedRemoved = {5,6};
  std::set<int> added;
  std::set<int> deleted;
  diffAddDelete(s1, s2, added, deleted);
  EXPECT_EQ(added, expectedAdded);
  EXPECT_EQ(deleted, expectedRemoved);
}

TEST(DiffTest, StringSet) {
  std::set<std::string> s1 = {"hi", "my", "name", "is", "ken"};
  std::set<std::string> s2 = {"hi", "name", "dog", "my"};
  std::set<std::string> expectedAdded = {"dog"};
  std::set<std::string> expectedRemoved = {"is", "ken"};
  std::set<std::string> added;
  std::set<std::string> deleted;
  diffAddDelete(s1, s2, added, deleted);
  EXPECT_EQ(added, expectedAdded);
  EXPECT_EQ(deleted, expectedRemoved);
}

// TODO: write tests for diff on DeviceConfig

// TODO: write tests for diff on std::list, std::vec, and std::map
// (once you implement that - right now only std::set supported)

} // namespace test
} // namespace devmand
