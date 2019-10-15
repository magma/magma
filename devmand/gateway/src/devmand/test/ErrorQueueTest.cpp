// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/test/EventBaseTest.h>
#include <devmand/ErrorQueue.h>
#include <folly/json.h>

namespace devmand {
namespace test {

TEST(ErrorQueueTest, SimpleEnqueueDequeue) {
  // create some errors (queue max length 5)
  ErrorQueue errors{5};
  errors.add("str one");
  errors.add("str two");
  errors.add("three");

  // get back the errors in order
  auto errs1 = folly::toJson(errors.get());
  auto expected1 =
      folly::toJson(folly::dynamic::array("str one", "str two", "three"));
  EXPECT_EQ(errs1, expected1);

  // add an error and "get" again
  errors.add("four!");
  auto errs2 = folly::toJson(errors.get());
  auto expected2 = folly::toJson(
      folly::dynamic::array("str one", "str two", "three", "four!"));
  EXPECT_EQ(errs2, expected2);

  // add two more errors, expect the oldest error to be discarded
  errors.add("FIVE");
  errors.add("sixth");
  auto errs3 = folly::toJson(errors.get());
  auto expected3 = folly::toJson(
      folly::dynamic::array("str two", "three", "four!", "FIVE", "sixth"));
  EXPECT_EQ(errs3, expected3);
}

} // namespace test
} // namespace devmand
