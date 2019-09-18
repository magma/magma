// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

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
