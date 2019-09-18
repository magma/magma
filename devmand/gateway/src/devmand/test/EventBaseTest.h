// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <future>

#include <gtest/gtest.h>

#include <folly/io/async/EventBase.h>

namespace devmand {
namespace test {

class EventBaseTest : public ::testing::Test {
 public:
  EventBaseTest();
  virtual ~EventBaseTest();
  EventBaseTest(const EventBaseTest&) = delete;
  EventBaseTest& operator=(const EventBaseTest&) = delete;
  EventBaseTest(EventBaseTest&&) = delete;
  EventBaseTest& operator=(EventBaseTest&&) = delete;

 protected:
  void start();
  void stop();

 protected:
  folly::EventBase eventBase;

 private:
  bool started{false};
  std::future<void> eventBaseThread;
};

} // namespace test
} // namespace devmand
