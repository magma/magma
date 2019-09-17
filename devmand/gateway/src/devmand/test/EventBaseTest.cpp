// Copyright (c) Facebook, Inc. and its affiliates. All Rights Reserved.

#include <devmand/test/EventBaseTest.h>

namespace devmand {
namespace test {

EventBaseTest::EventBaseTest() {
  start();
}

EventBaseTest::~EventBaseTest() {
  if (started) {
    stop();
  }
}
void EventBaseTest::start() {
  assert(not started);
  eventBaseThread =
      std::async(std::launch::async, [this] { eventBase.loopForever(); });
  started = true;
}

void EventBaseTest::stop() {
  assert(started);
  eventBase.terminateLoopSoon();
  eventBaseThread.wait();
  started = false;
}

} // namespace test
} // namespace devmand
