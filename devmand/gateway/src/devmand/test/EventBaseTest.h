// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

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
