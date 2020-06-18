// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <gtest/gtest.h>

#include <devmand/utils/Diff.h>

#include <list>
#include <set>

namespace devmand {
namespace test {

// genericInsert() either inserts an element into a set or pushes it onto the
// back of some list. It's used in diffAddDelete depending on the type of
// collection being used.

template <class Element>
void genericInsert(std::set<Element>& collection, const Element& element) {
  collection.insert(element);
}

template <class Element>
void genericInsert(std::list<Element>& collection, const Element& element) {
  collection.push_back(element);
}

// Populate `added` and `deleted` with the relevant elements between s1 and s2
// using the diff function. Useful for unit testing.
// template <template <class> class Container, class Type>
template <class Container, class Type>
void diffAddDelete(
    Container& s1,
    Container& s2,
    Container& added,
    Container& deleted) {
  DiffEventHandler<Type> deh = [&added, &deleted](
                                   DiffEvent de, const Type& val) {
    switch (de) {
      case DiffEvent::Add:
        genericInsert(added, val);
        break;
      case DiffEvent::Delete:
        genericInsert(deleted, val);
        break;
      case DiffEvent::Modify:
        FAIL() << "This case can't be reached.";
        break;
    }
  };
  diff(s1, s2, deh);
}

TEST(DiffTest, IntegerSets) {
  std::set<int> added;
  std::set<int> deleted;
  std::set<int> s1 = {2, 5, 8, 4, 6};
  std::set<int> s2 = {9, 7, 4, 2, 8, 1};
  std::set<int> expectedAdded = {9, 7, 1};
  std::set<int> expectedRemoved = {5, 6};
  diffAddDelete<std::set<int>, int>(s1, s2, added, deleted);
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
  diffAddDelete<std::set<std::string>, std::string>(s1, s2, added, deleted);
  EXPECT_EQ(added, expectedAdded);
  EXPECT_EQ(deleted, expectedRemoved);
}

// the key here is that lists aren't sorted and contain duplicates, unlike sets
TEST(DiffTest, IntegerList) {
  std::list<int> l1 = {4, 5, 5, 5, 3, 7, 7, 2, 2, 2};
  std::list<int> l2 = {4, 5, 3, 3, 7, 2, 2, 2, 2};
  std::list<int> expectedAdded = {2, 3};
  std::list<int> expectedDeleted = {5, 5, 7};
  std::list<int> added;
  std::list<int> deleted;
  diffAddDelete<std::list<int>, int>(l1, l2, added, deleted);
  EXPECT_EQ(added, expectedAdded);
  EXPECT_EQ(deleted, expectedDeleted);
}

TEST(DiffTest, StringList) {
  std::list<std::string> l1 = {"hi", "my", "name", "is", "ken"};
  std::list<std::string> l2 = {"hi", "name", "dog", "my"};
  std::list<std::string> expectedAdded = {"dog"};
  std::list<std::string> expectedRemoved = {"is", "ken"};
  std::list<std::string> added;
  std::list<std::string> deleted;
  diffAddDelete<std::list<std::string>, std::string>(l1, l2, added, deleted);
  EXPECT_EQ(added, expectedAdded);
  EXPECT_EQ(deleted, expectedRemoved);
}

// TODO: write tests for diff on DeviceConfig

} // namespace test
} // namespace devmand
