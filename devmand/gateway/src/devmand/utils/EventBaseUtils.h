// Copyright (c) 2016-present, Facebook, Inc.
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree. An additional grant
// of patent rights can be found in the PATENTS file in the same directory.

#pragma once

#include <folly/io/async/EventBase.h>

#include <devmand/error/ErrorHandler.h>

namespace devmand {

class EventBaseUtils final {
 public:
  EventBaseUtils() = delete;
  ~EventBaseUtils() = delete;
  EventBaseUtils(const EventBaseUtils&) = delete;
  EventBaseUtils& operator=(const EventBaseUtils&) = delete;
  EventBaseUtils(EventBaseUtils&&) = delete;
  EventBaseUtils& operator=(EventBaseUtils&&) = delete;

 public:
  static inline void scheduleEvery(
      folly::EventBase& eventBase,
      std::function<void()> event,
      const std::chrono::seconds& seconds) {
    scheduleEvery(
        eventBase,
        event,
        std::chrono::duration_cast<std::chrono::milliseconds>(seconds));
  }

  static inline void scheduleEvery(
      folly::EventBase& eventBase,
      std::function<void()> event,
      const std::chrono::milliseconds& milliseconds) {
    if (eventBase.isRunning()) {
      eventBase.runInEventBaseThread([&eventBase, event, milliseconds]() {
        ErrorHandler::executeWithCatch([&event]() { event(); });

        std::function<void()> recurse = [&eventBase, event, milliseconds]() {
          scheduleEvery(eventBase, event, milliseconds);
        };

        eventBase.scheduleAt(recurse, eventBase.now() + milliseconds);
      });
    }
  }

  static inline void scheduleIn(
      folly::EventBase& eventBase,
      std::function<void()> event,
      const std::chrono::seconds& seconds) {
    eventBase.runInEventBaseThread([&eventBase, event, seconds]() {
      std::function<void()> recurse = [event]() {
        ErrorHandler::executeWithCatch([event]() { event(); });
      };

      eventBase.scheduleAt(recurse, eventBase.now() + seconds);
    });
  }
};

} // namespace devmand
