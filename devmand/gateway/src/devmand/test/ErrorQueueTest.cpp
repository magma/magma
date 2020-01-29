// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#include <devmand/error/ErrorQueue.h>
#include <folly/json.h>
#include <gtest/gtest.h>
#include <future>

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

static int
doManyAdds(int addCount, int getInterval, std::shared_ptr<ErrorQueue> errors) {
  // "add" lots of errors
  for (int i = 0; i < addCount; ++i) {
    errors->add(std::move(folly::to<std::string>(i)));
    // every n-th add, do a "get"
    if (i % getInterval == 0) {
      errors.get();
    }
  }
  return 1;
}

TEST(ErrorQueueTest, MultithreadAdd) {
  std::shared_ptr<ErrorQueue> errors = std::make_shared<ErrorQueue>(20);
  std::vector<std::future<int>> futures;
  auto numThreads = 3;
  // start 3 threads adding to ErrorQueue
  for (int i = 0; i < numThreads; ++i) {
    // each thread will "add" 10,000 strings and "get" every 50
    futures.emplace_back(std::async(doManyAdds, 10000, 50, errors));
  }
  // wait for all three threads to finish
  auto threadCount = 0;
  for (auto& f : futures) {
    threadCount += f.get();
  }
  // no segfaults, so test passed
  EXPECT_EQ(threadCount, numThreads);
}

} // namespace test
} // namespace devmand
